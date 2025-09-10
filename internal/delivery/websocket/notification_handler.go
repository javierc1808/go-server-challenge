package websocket

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"frontend-challenge/internal/domain/entity"
	"frontend-challenge/internal/usecase"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/gorilla/websocket"
)

// NotificationHandler maneja las conexiones WebSocket para notificaciones
type NotificationHandler struct {
	notificationUsecase *usecase.NotificationUsecase
	upgrader            websocket.Upgrader
}

// NewNotificationHandler crea una nueva instancia de NotificationHandler
func NewNotificationHandler(notificationUsecase *usecase.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{
		notificationUsecase: notificationUsecase,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// TODO: Implementar validación de origen en producción
				return true
			},
		},
	}
}

// HandleNotifications maneja las conexiones WebSocket para notificaciones
func (h *NotificationHandler) HandleNotifications(w http.ResponseWriter, r *http.Request) {
	// Agregar headers de seguridad
	h.addSecurityHeaders(w)

	if r.Method == "OPTIONS" {
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}
	defer conn.Close()

	// Enviar notificaciones en tiempo real
	h.sendNotifications(conn)
}

// sendNotifications envía notificaciones simuladas
func (h *NotificationHandler) sendNotifications(conn *websocket.Conn) {
	for {
		// Crear notificación simulada
		notification := entity.NewNotification(
			gofakeit.UUID(),
			gofakeit.Name(),
			gofakeit.UUID(),
			gofakeit.BeerName(),
		)

		// Enviar notificación al cliente
		if err := conn.WriteJSON(notification); err != nil {
			log.Printf("Error writing to websocket: %v", err)
			break
		}

		// Esperar un tiempo aleatorio antes de la siguiente notificación
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	}
}

// addSecurityHeaders añade headers de seguridad
func (h *NotificationHandler) addSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: Restringir en producción
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
