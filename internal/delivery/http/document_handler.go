package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"frontend-challenge/internal/domain/entity"
	"frontend-challenge/internal/usecase"
	"frontend-challenge/pkg/security"
)

// NotificationBroadcaster minimal interface to emit notifications (implemented by websocket.Hub)
type NotificationBroadcaster interface {
	BroadcastNotification(notification *entity.Notification)
}

// DocumentHandler handles HTTP requests for documents
type DocumentHandler struct {
	documentUsecase *usecase.DocumentUsecase
	sanitizer       *security.Sanitizer
	notifier        NotificationBroadcaster
}

// NewDocumentHandler creates a new DocumentHandler instance
func NewDocumentHandler(documentUsecase *usecase.DocumentUsecase) *DocumentHandler {
	return &DocumentHandler{
		documentUsecase: documentUsecase,
		sanitizer:       security.NewSanitizer(),
	}
}

// WithNotifier allows injecting a notification broadcaster
func (h *DocumentHandler) WithNotifier(n NotificationBroadcaster) *DocumentHandler {
	h.notifier = n
	return h
}

// GetDocuments handles GET /documents
func (h *DocumentHandler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	// Add security headers
	h.addSecurityHeaders(w)

	// Handle CORS
	if r.Method == "OPTIONS" {
		return
	}

	documents, err := h.documentUsecase.GetAllDocuments(r.Context())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(documents); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// CreateDocument handles POST /documents
func (h *DocumentHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	// Add security headers
	h.addSecurityHeaders(w)

	// Handle CORS
	if r.Method == "OPTIONS" {
		return
	}

	// Only allow POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the document from the body
	var document entity.Document
	if err := json.NewDecoder(r.Body).Decode(&document); err != nil {
		http.Error(w, "Error decoding document: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields and domain rules
	if err := validateRequiredFields(&document); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// The server sets the timestamps automatically
	applyTimestamps(&document, time.Now())

	// Create the document in the cache
	if err := h.documentUsecase.CreateDocument(r.Context(), &document); err != nil {
		http.Error(w, "Error creating document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userID := r.Header.Get("user-id")
	userName := r.Header.Get("user-name")

	// Emit notification if notifier is present
	if h.notifier != nil {
		n := entity.NewNotification(
			userID,
			userName,
			document.ID,
			document.Title,
			"document.created",
		)
		h.notifier.BroadcastNotification(n)
	}

	// Respond with the created document (exactly as sent + server timestamps)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(document); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// addSecurityHeaders adds security headers
func (h *DocumentHandler) addSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: Restrict in production
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// --- helpers to reduce cognitive complexity ---

func validateRequiredFields(d *entity.Document) error {
	if d.ID == "" || d.Title == "" || d.Version == "" {
		return errors.New("missing required fields (id, title, version)")
	}
	return d.Validate()
}

func applyTimestamps(d *entity.Document, now time.Time) {
	d.CreatedAt = now
	d.UpdatedAt = now
	for i := range d.Contributors {
		if d.Contributors[i].CreatedAt.IsZero() {
			d.Contributors[i].CreatedAt = now
		}
		if d.Contributors[i].UpdatedAt.IsZero() {
			d.Contributors[i].UpdatedAt = now
		}
	}
}
