package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements per-IP rate limiting
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new RateLimiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Cleanup old requests every minute
	go rl.cleanup()

	return rl
}

// Middleware returns the rate limiting middleware
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr

		// Get real IP if behind a proxy
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

// allowRequest checks if the request is allowed
func (rl *RateLimiter) allowRequest(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Get existing requests for this IP
	requests := rl.requests[clientIP]

	// Filter requests within the time window
	var validRequests []time.Time
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Check if it exceeds the limit
	if len(validRequests) >= rl.limit {
		return false
	}

	// Add the new request
	validRequests = append(validRequests, now)
	rl.requests[clientIP] = validRequests

	return true
}

// cleanup removes old requests
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

// GetStats returns rate limiter statistics
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
