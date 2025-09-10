package security

import "errors"

// Errores de seguridad
var (
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrInvalidUUID         = errors.New("invalid UUID format")
	ErrInvalidAlphanumeric = errors.New("invalid alphanumeric input")
	ErrInputTooLong        = errors.New("input exceeds maximum length")
	ErrEmptyInput          = errors.New("input cannot be empty")
	ErrSuspiciousInput     = errors.New("suspicious input detected")
)
