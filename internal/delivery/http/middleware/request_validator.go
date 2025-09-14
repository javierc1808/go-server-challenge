package middleware

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// RequestValidator implements request size validation
type RequestValidator struct {
	maxBodySize int64
}

// NewRequestValidator creates a new RequestValidator
func NewRequestValidator(maxBodySize int64) *RequestValidator {
	return &RequestValidator{
		maxBodySize: maxBodySize,
	}
}

func isWebSocketHandshake(r *http.Request) bool {
	conn := strings.ToLower(r.Header.Get("Connection"))
	upg := strings.ToLower(r.Header.Get("Upgrade"))
	return strings.Contains(conn, "upgrade") && upg == "websocket"
}

// Middleware returns the request validation middleware
func (rv *RequestValidator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate body size
		if r.ContentLength > rv.maxBodySize {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		// Limit the size of the request body read
		r.Body = http.MaxBytesReader(w, r.Body, rv.maxBodySize)

		if isWebSocketHandshake(r) {
			next.ServeHTTP(w, r)
			return
		}

		// Validate important headers
		if !rv.validateHeaders(r) {
			http.Error(w, "Invalid headers", http.StatusBadRequest)
			return
		}

		// Authorization: Basic base64(user-name:user-id)
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Authorization header is required", http.StatusBadRequest)
			return
		}

		userName, userID, err := extractBasicUserHeaders(auth)
		if err != nil || userName == "" || userID == "" {
			http.Error(w, "Invalid Authorization", http.StatusBadRequest)
			return
		}

		r.Header.Set("user-name", userName)
		r.Header.Set("user-id", userID)

		next.ServeHTTP(w, r)
	})
}

// validateHeaders validates important headers
func (rv *RequestValidator) validateHeaders(r *http.Request) bool {
	// Validate Content-Length if present
	if contentLength := r.Header.Get("Content-Length"); contentLength != "" {
		if length, err := strconv.ParseInt(contentLength, 10, 64); err != nil || length > rv.maxBodySize {
			return false
		}
	}

	// Validate Content-Type for requests with a body
	if r.Method != "GET" && r.Method != "HEAD" && r.Method != "DELETE" {
		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			return false
		}
	}

	return true
}

// GetMaxBodySize returns the maximum body size
func (rv *RequestValidator) GetMaxBodySize() int64 {
	return rv.maxBodySize
}

// extractBasicUserHeaders parses Authorization: Basic base64(user-name:user-id)
// and returns (userName, userID).
func extractBasicUserHeaders(authorization string) (string, string, error) {
	const prefix = "Basic "
	if !strings.HasPrefix(authorization, prefix) {
		return "", "", errors.New("authorization must use basic scheme")
	}
	encoded := strings.TrimSpace(strings.TrimPrefix(authorization, prefix))
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", err
	}
	creds := string(decodedBytes)
	parts := strings.SplitN(creds, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", errors.New("authorization basic payload must be 'user-name:user-id'")
	}
	return parts[0], parts[1], nil
}
