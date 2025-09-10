package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implementa rate limiting por IP
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter crea un nuevo RateLimiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Limpiar requests antiguos cada minuto
	go rl.cleanup()

	return rl
}

// Middleware retorna el middleware de rate limiting
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr

		// Obtener IP real si está detrás de proxy
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			clientIP = forwarded
		} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			clientIP = realIP
		}

		if !rl.allowRequest(clientIP) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// allowRequest verifica si se permite la petición
func (rl *RateLimiter) allowRequest(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Obtener requests existentes para esta IP
	requests := rl.requests[clientIP]

	// Filtrar requests dentro de la ventana de tiempo
	var validRequests []time.Time
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Verificar si excede el límite
	if len(validRequests) >= rl.limit {
		return false
	}

	// Añadir la nueva petición
	validRequests = append(validRequests, now)
	rl.requests[clientIP] = validRequests

	return true
}

// cleanup limpia requests antiguos
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window)

		for ip, requests := range rl.requests {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if reqTime.After(cutoff) {
					validRequests = append(validRequests, reqTime)
				}
			}

			if len(validRequests) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = validRequests
			}
		}
		rl.mutex.Unlock()
	}
}

// GetStats retorna estadísticas del rate limiter
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	activeIPs := 0
	totalRequests := 0

	for _, requests := range rl.requests {
		validRequests := 0
		for _, reqTime := range requests {
			if reqTime.After(cutoff) {
				validRequests++
			}
		}
		if validRequests > 0 {
			activeIPs++
			totalRequests += validRequests
		}
	}

	return map[string]interface{}{
		"active_ips":     activeIPs,
		"total_requests": totalRequests,
		"limit":          rl.limit,
		"window_seconds": rl.window.Seconds(),
	}
}
