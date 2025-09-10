package security

import (
	"html"
	"regexp"
	"strings"
)

// Sanitizer proporciona funciones de sanitización
type Sanitizer struct {
	// Patrones de regex para validación
	emailRegex        *regexp.Regexp
	uuidRegex         *regexp.Regexp
	alphanumericRegex *regexp.Regexp
}

// NewSanitizer crea una nueva instancia de Sanitizer
func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		emailRegex:        regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		uuidRegex:         regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
		alphanumericRegex: regexp.MustCompile(`^[a-zA-Z0-9\s\-_.,!?]+$`),
	}
}

// SanitizeString sanitiza una cadena de texto
func (s *Sanitizer) SanitizeString(input string) string {
	// Escapar caracteres HTML
	sanitized := html.EscapeString(input)

	// Eliminar caracteres de control
	sanitized = s.removeControlCharacters(sanitized)

	// Normalizar espacios en blanco
	sanitized = s.normalizeWhitespace(sanitized)

	return strings.TrimSpace(sanitized)
}

// SanitizeEmail sanitiza y valida un email
func (s *Sanitizer) SanitizeEmail(email string) (string, error) {
	sanitized := strings.ToLower(strings.TrimSpace(email))

	if !s.emailRegex.MatchString(sanitized) {
		return "", ErrInvalidEmail
	}

	return sanitized, nil
}

// SanitizeUUID sanitiza y valida un UUID
func (s *Sanitizer) SanitizeUUID(uuid string) (string, error) {
	sanitized := strings.ToLower(strings.TrimSpace(uuid))

	if !s.uuidRegex.MatchString(sanitized) {
		return "", ErrInvalidUUID
	}

	return sanitized, nil
}

// SanitizeAlphanumeric sanitiza texto alfanumérico
func (s *Sanitizer) SanitizeAlphanumeric(input string) (string, error) {
	sanitized := s.SanitizeString(input)

	if !s.alphanumericRegex.MatchString(sanitized) {
		return "", ErrInvalidAlphanumeric
	}

	return sanitized, nil
}

// removeControlCharacters elimina caracteres de control
func (s *Sanitizer) removeControlCharacters(input string) string {
	var result strings.Builder
	for _, r := range input {
		if r >= 32 || r == '\t' || r == '\n' || r == '\r' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// normalizeWhitespace normaliza espacios en blanco
func (s *Sanitizer) normalizeWhitespace(input string) string {
	// Reemplazar múltiples espacios con uno solo
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(input, " ")
}

// ValidateInput valida y sanitiza input genérico
func (s *Sanitizer) ValidateInput(input string, maxLength int) (string, error) {
	if len(input) > maxLength {
		return "", ErrInputTooLong
	}

	sanitized := s.SanitizeString(input)

	if len(sanitized) == 0 {
		return "", ErrEmptyInput
	}

	return sanitized, nil
}
