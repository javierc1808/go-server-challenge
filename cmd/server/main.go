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
	"frontend-challenge/internal/delivery/websocket"
	"frontend-challenge/internal/infrastructure/repository"
	"frontend-challenge/internal/usecase"
	"frontend-challenge/pkg/config"
	"frontend-challenge/pkg/logger"

	"github.com/brianvoe/gofakeit/v5"
)

func main() {
	// Cargar configuración
	cfg := config.Load()

	// Inicializar logger
	logger := logger.NewSimpleLogger()

	// Inicializar generadores de datos aleatorios
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	gofakeit.Seed(seed)

	// Inicializar repositorios
	documentRepo := repository.NewDocumentRepositoryImpl()
	userRepo := repository.NewUserRepositoryImpl()
	notificationRepo := repository.NewNotificationRepositoryImpl()

	// Inicializar casos de uso
	documentUsecase := usecase.NewDocumentUsecase(documentRepo, userRepo)
	notificationUsecase := usecase.NewNotificationUsecase(notificationRepo, documentRepo, userRepo)

	// Inicializar handlers
	documentHandler := deliveryhttp.NewDocumentHandler(documentUsecase)
	notificationHandler := websocket.NewNotificationHandler(notificationUsecase)

	// Configurar rutas
	mux := http.NewServeMux()
	mux.HandleFunc("/documents", documentHandler.GetDocuments)
	mux.HandleFunc("/notifications", notificationHandler.HandleNotifications)

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
