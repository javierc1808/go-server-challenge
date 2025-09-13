package main

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	deliveryhttp "frontend-challenge/internal/delivery/http"
	"frontend-challenge/internal/delivery/http/middleware"
	"frontend-challenge/internal/delivery/websocket"
	"frontend-challenge/internal/infrastructure/repository"
	"frontend-challenge/internal/usecase"
	"frontend-challenge/pkg/config"
	"frontend-challenge/pkg/logger"
	"frontend-challenge/pkg/security"

	"github.com/brianvoe/gofakeit/v5"
)

// buildHTTPHandler wires middlewares and routes
func buildHTTPHandler(
	threatMonitor *security.ThreatMonitor,
	documentHandler *deliveryhttp.DocumentHandler,
	notificationHandler *websocket.NotificationHandler,
	securityHandler *deliveryhttp.SecurityHandler,
	rateLimiter *middleware.RateLimiter,
	requestValidator *middleware.RequestValidator,
	securityHeaders *middleware.SecurityHeaders,
) http.Handler {
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Threat monitoring
		if threatMonitor.AnalyzeRequest(r) {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		// Route requests
		switch r.URL.Path {
		case "/documents":
			switch r.Method {
			case "GET":
				documentHandler.GetDocuments(w, r)
			case "POST":
				documentHandler.CreateDocument(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		case "/notifications":
			notificationHandler.HandleNotifications(w, r)
		case "/security/stats":
			securityHandler.GetSecurityStats(w, r)
		case "/health":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		default:
			http.NotFound(w, r)
		}
	})

	// Apply middlewares in order
	return rateLimiter.Middleware(
		requestValidator.Middleware(
			securityHeaders.Middleware(base),
		),
	)
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := logger.NewSimpleLogger()

	// Initialize security logger
	securityLogger, err := security.NewSecurityLogger("logs/security.log")
	if err != nil {
		logger.Error("Error creating security logger", err)
		os.Exit(1)
	}
	defer securityLogger.Close()

	// Initialize threat monitoring
	threatMonitor := security.NewThreatMonitor(securityLogger)

	// Initialize log rotation
	logRotator := security.NewLogRotator("logs", 10, 10*1024*1024, true) // 10 files, 10MB, daily rotation
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := logRotator.RotateLogs(); err != nil {
				logger.Error("Error rotating logs", err)
			}
		}
	}()

	// Initialize random data generators
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	gofakeit.Seed(seed)

	// Initialize cache
	cache := repository.NewMemoryCache(24 * time.Hour) // 24 hour TTL

	// Initialize repositories
	documentRepo := repository.NewDocumentRepositoryImpl(cache)
	userRepo := repository.NewUserRepositoryImpl()
	notificationRepo := repository.NewNotificationRepositoryImpl()

	// Initialize use cases
	documentUsecase := usecase.NewDocumentUsecase(documentRepo, userRepo)
	notificationUsecase := usecase.NewNotificationUsecase(notificationRepo, documentRepo, userRepo)

	// Initialize handlers
	notificationHandler := websocket.NewNotificationHandler(notificationUsecase)
	documentHandler := deliveryhttp.NewDocumentHandler(documentUsecase).WithNotifier(notificationHandler.Hub())

	// Configure security middlewares
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)      // 100 requests per minute
	requestValidator := middleware.NewRequestValidator(1024 * 1024) // 1MB max
	securityHeaders := middleware.NewSecurityHeaders(true)          // Enable CSP
	securityHandler := deliveryhttp.NewSecurityHandler(threatMonitor, rateLimiter, logRotator, cache)

	// Configure routes with middlewares
	mux := http.NewServeMux()
	handler := buildHTTPHandler(threatMonitor, documentHandler, notificationHandler, securityHandler, rateLimiter, requestValidator, securityHeaders)
	mux.Handle("/", handler)

	// Configure server
	server := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Channel to handle system signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		logger.Info("Server started at " + cfg.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error starting server", err)
			os.Exit(1)
		}
	}()

	// Wait for termination signal
	<-done
	logger.Info("Server shutting down...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Gracefully close server
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Error shutting down server", err)
		os.Exit(1)
	}

	logger.Info("Server closed successfully")
}
