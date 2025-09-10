package entity

import "time"

// Notification representa una notificación en el dominio
type Notification struct {
	Timestamp     time.Time `json:"timestamp"`
	UserID        string    `json:"user_id"`
	UserName      string    `json:"user_name"`
	DocumentID    string    `json:"document_id"`
	DocumentTitle string    `json:"document_title"`
}

// NewNotification crea una nueva instancia de Notification
func NewNotification(userID, userName, documentID, documentTitle string) *Notification {
	return &Notification{
		Timestamp:     time.Now(),
		UserID:        userID,
		UserName:      userName,
		DocumentID:    documentID,
		DocumentTitle: documentTitle,
	}
}

// Validate valida los datos de la notificación
func (n *Notification) Validate() error {
	if n.UserID == "" {
		return ErrInvalidUserID
	}
	if n.UserName == "" {
		return ErrInvalidUserName
	}
	if n.DocumentID == "" {
		return ErrInvalidDocumentID
	}
	if n.DocumentTitle == "" {
		return ErrInvalidDocumentTitle
	}
	return nil
}
