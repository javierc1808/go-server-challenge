package repository

import (
	"context"

	"frontend-challenge/internal/domain/entity"
)

// NotificationRepository define la interfaz para el repositorio de notificaciones
type NotificationRepository interface {
	// Create crea una nueva notificación
	Create(ctx context.Context, notification *entity.Notification) error

	// GetByUserID obtiene notificaciones por ID de usuario
	GetByUserID(ctx context.Context, userID string) ([]*entity.Notification, error)

	// GetAll obtiene todas las notificaciones
	GetAll(ctx context.Context) ([]*entity.Notification, error)
}
