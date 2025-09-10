package middleware

import (
	"net/http"
)

// SecurityHeaders implementa headers de seguridad avanzados
type SecurityHeaders struct {
	enableCSP bool
	cspPolicy string
}

// NewSecurityHeaders crea una nueva instancia de SecurityHeaders
func NewSecurityHeaders(enableCSP bool) *SecurityHeaders {
	return &SecurityHeaders{
		enableCSP: enableCSP,
		cspPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' ws: wss:; frame-ancestors 'none';",
	}
}

// Middleware retorna el middleware de headers de seguridad
func (sh *SecurityHeaders) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Headers básicos de seguridad
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Content Security Policy
		if sh.enableCSP {
			w.Header().Set("Content-Security-Policy", sh.cspPolicy)
		}

		// Headers adicionales de seguridad
		w.Header().Set("X-DNS-Prefetch-Control", "off")
		w.Header().Set("X-Download-Options", "noopen")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")

		// CORS restrictivo (debería ser configurado por dominio)
		w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: Restringir en producción
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		next.ServeHTTP(w, r)
	})
}

// SetCSPPolicy establece una política CSP personalizada
func (sh *SecurityHeaders) SetCSPPolicy(policy string) {
	sh.cspPolicy = policy
}

// GetCSPPolicy retorna la política CSP actual
func (sh *SecurityHeaders) GetCSPPolicy() string {
	return sh.cspPolicy
}
