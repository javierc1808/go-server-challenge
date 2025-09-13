package security

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// SecurityEvent represents a security event
type SecurityEvent struct {
	Timestamp   time.Time              `json:"timestamp"`
	EventType   string                 `json:"event_type"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	RequestPath string                 `json:"request_path"`
	Method      string                 `json:"method"`
	Severity    string                 `json:"severity"`
	Message     string                 `json:"message"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// SecurityLogger handles security event logging
type SecurityLogger struct {
	logger *log.Logger
	file   *os.File
}

// NewSecurityLogger creates a new instance of SecurityLogger
func NewSecurityLogger(logFile string) (*SecurityLogger, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	logger := log.New(file, "SECURITY: ", log.LstdFlags|log.LUTC)

	return &SecurityLogger{
		logger: logger,
		file:   file,
	}, nil
}

// LogEvent logs a security event
func (sl *SecurityLogger) LogEvent(event SecurityEvent) {
	event.Timestamp = time.Now().UTC()

	// Convert to JSON
	jsonData, err := json.Marshal(event)
	if err != nil {
		sl.logger.Printf("Error marshaling security event: %v", err)
		return
	}

	sl.logger.Println(string(jsonData))
}

// LogSuspiciousActivity logs suspicious activity
func (sl *SecurityLogger) LogSuspiciousActivity(r *http.Request, message string, details map[string]interface{}) {
	event := SecurityEvent{
		EventType:   "SUSPICIOUS_ACTIVITY",
		IPAddress:   sl.getClientIP(r),
		UserAgent:   r.UserAgent(),
		RequestPath: r.URL.Path,
		Method:      r.Method,
		Severity:    "HIGH",
		Message:     message,
		Details:     details,
	}

	sl.LogEvent(event)
}

// LogRateLimitExceeded logs rate limit exceeded
func (sl *SecurityLogger) LogRateLimitExceeded(r *http.Request, limit int) {
	event := SecurityEvent{
		EventType:   "RATE_LIMIT_EXCEEDED",
		IPAddress:   sl.getClientIP(r),
		UserAgent:   r.UserAgent(),
		RequestPath: r.URL.Path,
		Method:      r.Method,
		Severity:    "MEDIUM",
		Message:     fmt.Sprintf("Rate limit exceeded: %d requests", limit),
		Details: map[string]interface{}{
			"rate_limit": limit,
		},
	}

	sl.LogEvent(event)
}

// LogInvalidInput logs invalid input
func (sl *SecurityLogger) LogInvalidInput(r *http.Request, input string, error string) {
	event := SecurityEvent{
		EventType:   "INVALID_INPUT",
		IPAddress:   sl.getClientIP(r),
		UserAgent:   r.UserAgent(),
		RequestPath: r.URL.Path,
		Method:      r.Method,
		Severity:    "LOW",
		Message:     "Invalid input detected",
		Details: map[string]interface{}{
			"input": input,
			"error": error,
		},
	}

	sl.LogEvent(event)
}

// LogAuthenticationFailure logs authentication failure
func (sl *SecurityLogger) LogAuthenticationFailure(r *http.Request, username string) {
	event := SecurityEvent{
		EventType:   "AUTHENTICATION_FAILURE",
		IPAddress:   sl.getClientIP(r),
		UserAgent:   r.UserAgent(),
		RequestPath: r.URL.Path,
		Method:      r.Method,
		Severity:    "HIGH",
		Message:     "Authentication failure",
		Details: map[string]interface{}{
			"username": username,
		},
	}

	sl.LogEvent(event)
}

// getClientIP gets the real client IP
func (sl *SecurityLogger) getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	return r.RemoteAddr
}

// Close closes the logger
func (sl *SecurityLogger) Close() error {
	return sl.file.Close()
}
