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

func main() {
	// Cargar configuración
	cfg := config.Load()

	// Inicializar logger
	logger := logger.NewSimpleLogger()

	// Inicializar logger de seguridad
	securityLogger, err := security.NewSecurityLogger("logs/security.log")
	if err != nil {
		logger.Error("Error al crear logger de seguridad", err)
		os.Exit(1)
	}
	defer securityLogger.Close()

	// Inicializar monitoreo de amenazas
	threatMonitor := security.NewThreatMonitor(securityLogger)

	// Inicializar rotación de logs
	logRotator := security.NewLogRotator("logs", 10, 10*1024*1024, true) // 10 archivos, 10MB, rotación diaria
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			if err := logRotator.RotateLogs(); err != nil {
				logger.Error("Error al rotar logs", err)
			}
		}
	}()

	// Inicializar generadores de datos aleatorios
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	gofakeit.Seed(seed)

	// Inicializar cache
	cache := repository.NewMemoryCache(24 * time.Hour) // TTL de 24 horas

	// Inicializar repositorios
	documentRepo := repository.NewDocumentRepositoryImpl(cache)
	userRepo := repository.NewUserRepositoryImpl()
	notificationRepo := repository.NewNotificationRepositoryImpl()

	// Inicializar casos de uso
	documentUsecase := usecase.NewDocumentUsecase(documentRepo, userRepo)
	notificationUsecase := usecase.NewNotificationUsecase(notificationRepo, documentRepo, userRepo)

	// Inicializar handlers
	documentHandler := deliveryhttp.NewDocumentHandler(documentUsecase)
	notificationHandler := websocket.NewNotificationHandler(notificationUsecase)

	// Configurar middlewares de seguridad
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)      // 100 requests por minuto
	requestValidator := middleware.NewRequestValidator(1024 * 1024) // 1MB máximo
	securityHeaders := middleware.NewSecurityHeaders(true)          // Habilitar CSP
	securityHandler := deliveryhttp.NewSecurityHandler(threatMonitor, rateLimiter, logRotator, cache)

	// Configurar rutas con middlewares
	mux := http.NewServeMux()

	// Aplicar middlewares en orden
	handler := rateLimiter.Middleware(
		requestValidator.Middleware(
			securityHeaders.Middleware(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Monitoreo de amenazas
					if threatMonitor.AnalyzeRequest(r) {
						http.Error(w, "Access denied", http.StatusForbidden)
						return
					}

					// Enrutar peticiones
					switch r.URL.Path {
					case "/documents":
						switch r.Method {
						case "GET":
							documentHandler.GetDocuments(w, r)
						case "POST":
							documentHandler.CreateDocument(w, r)
						default:
							http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
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
				}),
			),
		),
	)

	mux.Handle("/", handler)

	// Configurar servidor
	server := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Canal para manejar señales del sistema
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en goroutine
	go func() {
		logger.Info("Servidor iniciado en " + cfg.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error al iniciar servidor", err)
			os.Exit(1)
		}
	}()

	// Esperar señal de terminación
	<-done
	logger.Info("Servidor deteniéndose...")

	// Crear contexto con timeout para el shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Cerrar servidor gracefully
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Error al cerrar servidor", err)
		os.Exit(1)
	}

	logger.Info("Servidor cerrado correctamente")
}
