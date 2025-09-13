package repository

import (
	"context"

	"frontend-challenge/internal/domain/entity"
)

// NotificationRepository defines the interface for the notification repository
type NotificationRepository interface {
	// Create creates a new notification
	Create(ctx context.Context, notification *entity.Notification) error

	// GetByUserID gets notifications by user ID
	GetByUserID(ctx context.Context, userID string) ([]*entity.Notification, error)

	// GetAll gets all notifications
	GetAll(ctx context.Context) ([]*entity.Notification, error)
}
