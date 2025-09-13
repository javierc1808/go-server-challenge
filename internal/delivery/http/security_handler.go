package http

import (
	"encoding/json"
	"net/http"
	"time"

	"frontend-challenge/pkg/security"
)

// SecurityHandler handles security endpoints
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

// NewSecurityHandler creates a new instance of SecurityHandler
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

// GetSecurityStats handles the GET /security/stats request
func (h *SecurityHandler) GetSecurityStats(w http.ResponseWriter, r *http.Request) {
	// Add security headers
	h.addSecurityHeaders(w)

	// Get statistics
	threatStats := h.threatMonitor.GetThreatStats()
	rateLimitStats := h.rateLimiter.GetStats()
	logStats, err := h.logRotator.GetLogStats()
	if err != nil {
		http.Error(w, "Error getting log statistics", http.StatusInternalServerError)
		return
	}
	cacheStats := h.cache.GetStats()

	// Combine statistics
	stats := map[string]interface{}{
		"threats":     threatStats,
		"rate_limits": rateLimitStats,
		"logs":        logStats,
		"cache":       cacheStats,
		"timestamp":   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// addSecurityHeaders adds security headers
func (h *SecurityHandler) addSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: Restrict in production
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
