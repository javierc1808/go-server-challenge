package entity

import "errors"

// Errores del dominio
var (
	ErrInvalidUserID          = errors.New("invalid user ID")
	ErrInvalidUserName        = errors.New("invalid user name")
	ErrInvalidDocumentID      = errors.New("invalid document ID")
	ErrInvalidDocumentTitle   = errors.New("invalid document title")
	ErrInvalidDocumentVersion = errors.New("invalid document version")
	ErrUserNotFound           = errors.New("user not found")
	ErrDocumentNotFound       = errors.New("document not found")
	ErrInvalidNotification    = errors.New("invalid notification")
)
