package http

import (
	"encoding/json"
	"net/http"

	"frontend-challenge/internal/usecase"
)

// DocumentHandler maneja las peticiones HTTP para documentos
type DocumentHandler struct {
	documentUsecase *usecase.DocumentUsecase
}

// NewDocumentHandler crea una nueva instancia de DocumentHandler
func NewDocumentHandler(documentUsecase *usecase.DocumentUsecase) *DocumentHandler {
	return &DocumentHandler{
		documentUsecase: documentUsecase,
	}
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
