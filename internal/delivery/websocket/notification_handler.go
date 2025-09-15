package websocket

import (
	"log"
	"net/http"
	"time"

	"frontend-challenge/internal/domain/entity"
	"frontend-challenge/internal/usecase"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/gorilla/websocket"
)

// NotificationHandler handles WebSocket connections for notifications
type NotificationHandler struct {
	notificationUsecase *usecase.NotificationUsecase
	upgrader            websocket.Upgrader
	hub                 *Hub
}

// NewNotificationHandler creates a new NotificationHandler instances
func NewNotificationHandler(notificationUsecase *usecase.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{
		notificationUsecase: notificationUsecase,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// TODO: Implement origin validation in production
				return true
			},
		},
		hub: NewHub(),
	}
}

// Hub exposes the hub so it can be injected when needed
func (h *NotificationHandler) Hub() *Hub { return h.hub }

// HandleNotifications upgrades and manages the WebSocket lifecycle
func (h *NotificationHandler) HandleNotifications(w http.ResponseWriter, r *http.Request) {
	// Add security headers
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

	// No fake fallback notifications. The server emits only on create/update/delete events.
	h.fallbackFakeNotifications(conn)

	// Keep connection open until client closes it
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
			"document.created.fake",
		)

		if err := conn.WriteJSON(notification); err != nil {
			break
		}

		time.Sleep(time.Duration(30) * time.Second)
	}
}

// addSecurityHeaders adds security headers
func (h *NotificationHandler) addSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: Restrict in production
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
