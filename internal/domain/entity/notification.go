package entity

import "time"

// Notification represents a notification in the domain
type Notification struct {
	Timestamp     time.Time `json:"timestamp"`
	UserID        string    `json:"userId"`
	UserName      string    `json:"userName"`
	DocumentID    string    `json:"documentId"`
	DocumentTitle string    `json:"documentTitle"`
	Type          string    `json:"type"`
}

// NewNotification creates a new instance of Notification
func NewNotification(userID, userName, documentID, documentTitle, notificationType string) *Notification {
	return &Notification{
		Timestamp:     time.Now(),
		UserID:        userID,
		UserName:      userName,
		DocumentID:    documentID,
		DocumentTitle: documentTitle,
		Type:          notificationType,
	}
}

// Validate validates the notification data
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
	if n.Type == "" {
		return ErrInvalidNotificationType
	}
	return nil
}
