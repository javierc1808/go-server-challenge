package security

import (
	"html"
	"regexp"
	"strings"
)

// Sanitizer provides sanitization functions
type Sanitizer struct {
	// Regex patterns for validation
	emailRegex        *regexp.Regexp
	uuidRegex         *regexp.Regexp
	alphanumericRegex *regexp.Regexp
}

// NewSanitizer creates a new instance of Sanitizer
func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		emailRegex:        regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		uuidRegex:         regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`),
		alphanumericRegex: regexp.MustCompile(`^[a-zA-Z0-9\s\-_.,!?]+$`),
	}
}

// SanitizeString sanitizes a text string
func (s *Sanitizer) SanitizeString(input string) string {
	// Escape HTML characters
	sanitized := html.EscapeString(input)

	// Remove control characters
	sanitized = s.removeControlCharacters(sanitized)

	// Normalize whitespace
	sanitized = s.normalizeWhitespace(sanitized)

	return strings.TrimSpace(sanitized)
}

// SanitizeEmail sanitizes and validates an email
func (s *Sanitizer) SanitizeEmail(email string) (string, error) {
	sanitized := strings.ToLower(strings.TrimSpace(email))

	if !s.emailRegex.MatchString(sanitized) {
		return "", ErrInvalidEmail
	}

	return sanitized, nil
}

// SanitizeUUID sanitizes and validates a UUID
func (s *Sanitizer) SanitizeUUID(uuid string) (string, error) {
	sanitized := strings.ToLower(strings.TrimSpace(uuid))

	if !s.uuidRegex.MatchString(sanitized) {
		return "", ErrInvalidUUID
	}

	return sanitized, nil
}

// SanitizeAlphanumeric sanitizes alphanumeric text
func (s *Sanitizer) SanitizeAlphanumeric(input string) (string, error) {
	sanitized := s.SanitizeString(input)

	if !s.alphanumericRegex.MatchString(sanitized) {
		return "", ErrInvalidAlphanumeric
	}

	return sanitized, nil
}

// removeControlCharacters removes control characters
func (s *Sanitizer) removeControlCharacters(input string) string {
	var result strings.Builder
	for _, r := range input {
		if r >= 32 || r == '\t' || r == '\n' || r == '\r' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// normalizeWhitespace normalizes whitespace
func (s *Sanitizer) normalizeWhitespace(input string) string {
	// Replace multiple spaces with a single one
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(input, " ")
}

// ValidateInput validates and sanitizes generic input
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
