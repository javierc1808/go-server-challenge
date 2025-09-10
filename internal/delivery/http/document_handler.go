package http

import (
	"encoding/json"
	"net/http"
	"time"

	"frontend-challenge/internal/domain/entity"
	"frontend-challenge/internal/usecase"
	"frontend-challenge/pkg/security"
)

// NotificationBroadcaster interfaz mínima para emitir notificaciones (se implementa en websocket.Hub)
type NotificationBroadcaster interface {
	BroadcastNotification(notification *entity.Notification)
}

// DocumentHandler maneja las peticiones HTTP para documentos
type DocumentHandler struct {
	documentUsecase *usecase.DocumentUsecase
	sanitizer       *security.Sanitizer
	notifier        NotificationBroadcaster
}

// NewDocumentHandler crea una nueva instancia de DocumentHandler
func NewDocumentHandler(documentUsecase *usecase.DocumentUsecase) *DocumentHandler {
	return &DocumentHandler{
		documentUsecase: documentUsecase,
		sanitizer:       security.NewSanitizer(),
	}
}

// WithNotifier permite inyectar un emisor de notificaciones
func (h *DocumentHandler) WithNotifier(n NotificationBroadcaster) *DocumentHandler {
	h.notifier = n
	return h
}

// GetDocuments maneja la petición GET /documents
func (h *DocumentHandler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	// Agregar headers de seguridad
	h.addSecurityHeaders(w)

	// Manejar CORS
	if r.Method == "OPTIONS" {
		return
	}

	documents, err := h.documentUsecase.GetAllDocuments(r.Context())
	if err != nil {
		http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(documents); err != nil {
		http.Error(w, "Error al codificar respuesta", http.StatusInternalServerError)
		return
	}
}

// CreateDocument maneja la petición POST /documents
func (h *DocumentHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	// Agregar headers de seguridad
	h.addSecurityHeaders(w)

	// Manejar CORS
	if r.Method == "OPTIONS" {
		return
	}

	// Solo permitir POST
	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodificar el documento del body
	var document entity.Document
	if err := json.NewDecoder(r.Body).Decode(&document); err != nil {
		http.Error(w, "Error al decodificar documento: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validar que el documento tenga datos requeridos
	if document.ID == "" {
		http.Error(w, "ID del documento es requerido", http.StatusBadRequest)
		return
	}
	if document.Title == "" {
		http.Error(w, "Título del documento es requerido", http.StatusBadRequest)
		return
	}
	if document.Version == "" {
		http.Error(w, "Versión del documento es requerida", http.StatusBadRequest)
		return
	}

	// Validar el documento usando la función de validación
	if err := document.Validate(); err != nil {
		http.Error(w, "Documento inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validar user-name y user-id en el header
	userID := r.Header.Get("user-id")
	userName := r.Header.Get("user-name")
	if userID == "" {
		http.Error(w, "El header 'user-id' es requerido", http.StatusBadRequest)
		return
	}
	if userName == "" {
		http.Error(w, "El header 'user-name' es requerido", http.StatusBadRequest)
		return
	}

	// El servidor establece los timestamps automáticamente
	now := time.Now()
	document.CreatedAt = now
	document.UpdatedAt = now

	// Establecer timestamps para los contribuidores también
	for i := range document.Contributors {
		if document.Contributors[i].CreatedAt.IsZero() {
			document.Contributors[i].CreatedAt = now
		}
		if document.Contributors[i].UpdatedAt.IsZero() {
			document.Contributors[i].UpdatedAt = now
		}
	}

	// Crear el documento en el cache
	if err := h.documentUsecase.CreateDocument(r.Context(), &document); err != nil {
		http.Error(w, "Error al crear documento: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Emitir notificación si hay notifier
	if h.notifier != nil {
		n := entity.NewNotification(
			userID,
			userName,
			document.ID,
			document.Title,
		)
		h.notifier.BroadcastNotification(n)
	}

	// Responder con el documento creado (exactamente como se envió + timestamps del server)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(document); err != nil {
		http.Error(w, "Error al codificar respuesta", http.StatusInternalServerError)
		return
	}
}

// addSecurityHeaders añade headers de seguridad
func (h *DocumentHandler) addSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: Restringir en producción
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
