package middleware

import (
	"net/http"
	"strconv"
)

// RequestValidator implementa validación de tamaño de request
type RequestValidator struct {
	maxBodySize int64
}

// NewRequestValidator crea un nuevo RequestValidator
func NewRequestValidator(maxBodySize int64) *RequestValidator {
	return &RequestValidator{
		maxBodySize: maxBodySize,
	}
}

// Middleware retorna el middleware de validación de request
func (rv *RequestValidator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validar tamaño del body
		if r.ContentLength > rv.maxBodySize {
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
			return
		}

		// Limitar el tamaño del body leído
		r.Body = http.MaxBytesReader(w, r.Body, rv.maxBodySize)

		// Validar headers importantes
		if !rv.validateHeaders(r) {
			http.Error(w, "Invalid headers", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// validateHeaders valida headers importantes
func (rv *RequestValidator) validateHeaders(r *http.Request) bool {
	// Validar Content-Length si está presente
	if contentLength := r.Header.Get("Content-Length"); contentLength != "" {
		if length, err := strconv.ParseInt(contentLength, 10, 64); err != nil || length > rv.maxBodySize {
			return false
		}
	}

	// Validar Content-Type para requests con body
	if r.Method != "GET" && r.Method != "HEAD" && r.Method != "DELETE" {
		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			return false
		}
	}

	return true
}

// GetMaxBodySize retorna el tamaño máximo del body
func (rv *RequestValidator) GetMaxBodySize() int64 {
	return rv.maxBodySize
}
