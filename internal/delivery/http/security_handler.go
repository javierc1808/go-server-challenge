package http

import (
	"encoding/json"
	"net/http"
	"time"

	"frontend-challenge/pkg/security"
)

// SecurityHandler maneja endpoints de seguridad
type SecurityHandler struct {
	threatMonitor *security.ThreatMonitor
	rateLimiter   interface {
		GetStats() map[string]interface{}
	}
	logRotator interface {
		GetLogStats() (map[string]interface{}, error)
	}
	cache interface {
		GetStats() map[string]interface{}
	}
}

// NewSecurityHandler crea una nueva instancia de SecurityHandler
func NewSecurityHandler(threatMonitor *security.ThreatMonitor, rateLimiter interface {
	GetStats() map[string]interface{}
}, logRotator interface {
	GetLogStats() (map[string]interface{}, error)
}, cache interface {
	GetStats() map[string]interface{}
}) *SecurityHandler {
	return &SecurityHandler{
		threatMonitor: threatMonitor,
		rateLimiter:   rateLimiter,
		logRotator:    logRotator,
		cache:         cache,
	}
}

// GetSecurityStats maneja la petición GET /security/stats
func (h *SecurityHandler) GetSecurityStats(w http.ResponseWriter, r *http.Request) {
	// Agregar headers de seguridad
	h.addSecurityHeaders(w)

	// Obtener estadísticas
	threatStats := h.threatMonitor.GetThreatStats()
	rateLimitStats := h.rateLimiter.GetStats()
	logStats, err := h.logRotator.GetLogStats()
	if err != nil {
		http.Error(w, "Error al obtener estadísticas de logs", http.StatusInternalServerError)
		return
	}
	cacheStats := h.cache.GetStats()

	// Combinar estadísticas
	stats := map[string]interface{}{
		"threats":     threatStats,
		"rate_limits": rateLimitStats,
		"logs":        logStats,
		"cache":       cacheStats,
		"timestamp":   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Error al codificar respuesta", http.StatusInternalServerError)
		return
	}
}

// addSecurityHeaders añade headers de seguridad
func (h *SecurityHandler) addSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: Restringir en producción
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
