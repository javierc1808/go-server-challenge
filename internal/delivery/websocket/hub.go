package websocket

import (
	"sync"

	"frontend-challenge/internal/domain/entity"

	"github.com/gorilla/websocket"
)

// Notifier define la interfaz para difundir notificaciones a clientes WebSocket
type Notifier interface {
	BroadcastNotification(notification *entity.Notification)
}

// Hub gestiona conexiones WebSocket y difunde mensajes
type Hub struct {
	mu    sync.RWMutex
	conns map[*websocket.Conn]bool
}

// NewHub crea una nueva instancia de Hub
func NewHub() *Hub {
	return &Hub{
		conns: make(map[*websocket.Conn]bool),
	}
}

// Register añade una conexión al hub
func (h *Hub) Register(conn *websocket.Conn) {
	h.mu.Lock()
	h.conns[conn] = true
	h.mu.Unlock()
}

// Unregister elimina una conexión del hub
func (h *Hub) Unregister(conn *websocket.Conn) {
	h.mu.Lock()
	delete(h.conns, conn)
	h.mu.Unlock()
}

// BroadcastNotification envía la notificación a todas las conexiones activas
func (h *Hub) BroadcastNotification(notification *entity.Notification) {
	h.mu.RLock()
	for conn := range h.conns {
		if err := conn.WriteJSON(notification); err != nil {
			// Si hay error al escribir, desconectar
			h.mu.RUnlock()
			h.mu.Lock()
			conn.Close()
			delete(h.conns, conn)
			h.mu.Unlock()
			h.mu.RLock()
		}
	}
	h.mu.RUnlock()
}
