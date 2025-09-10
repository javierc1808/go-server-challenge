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
	hub                 *Hub
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
		hub: NewHub(),
	}
}

// Hub expone el hub para inyectarlo si fuese necesario
func (h *NotificationHandler) Hub() *Hub { return h.hub }

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

	h.hub.Register(conn)
	defer func() {
		h.hub.Unregister(conn)
		conn.Close()
	}()

	// Si no vamos a emitir desde el servidor, mantener un fallback de fake
	// go h.fallbackFakeNotifications(conn)

	// Mantener la conexión abierta mientras el cliente no cierre
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

// fallbackFakeNotifications envía notificaciones simuladas
func (h *NotificationHandler) fallbackFakeNotifications(conn *websocket.Conn) {
	for {
		notification := entity.NewNotification(
			gofakeit.UUID(),
			gofakeit.Name(),
			gofakeit.UUID(),
			gofakeit.BeerName(),
		)

		if err := conn.WriteJSON(notification); err != nil {
			break
		}

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
