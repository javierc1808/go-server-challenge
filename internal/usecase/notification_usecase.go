package usecase

import (
	"context"

	"frontend-challenge/internal/domain/entity"
	"frontend-challenge/internal/domain/repository"
)

// NotificationUsecase define los casos de uso para notificaciones
type NotificationUsecase struct {
	notificationRepo repository.NotificationRepository
	documentRepo     repository.DocumentRepository
	userRepo         repository.UserRepository
}

// NewNotificationUsecase crea una nueva instancia de NotificationUsecase
func NewNotificationUsecase(
	notificationRepo repository.NotificationRepository,
	documentRepo repository.DocumentRepository,
	userRepo repository.UserRepository,
) *NotificationUsecase {
	return &NotificationUsecase{
		notificationRepo: notificationRepo,
		documentRepo:     documentRepo,
		userRepo:         userRepo,
	}
}

// CreateNotification crea una nueva notificación
func (u *NotificationUsecase) CreateNotification(ctx context.Context, notification *entity.Notification) error {
	if err := notification.Validate(); err != nil {
		return err
	}
	return u.notificationRepo.Create(ctx, notification)
}

// GetNotificationsByUserID obtiene notificaciones por ID de usuario
func (u *NotificationUsecase) GetNotificationsByUserID(ctx context.Context, userID string) ([]*entity.Notification, error) {
	if userID == "" {
		return nil, entity.ErrInvalidUserID
	}
	return u.notificationRepo.GetByUserID(ctx, userID)
}

// GetAllNotifications obtiene todas las notificaciones
func (u *NotificationUsecase) GetAllNotifications(ctx context.Context) ([]*entity.Notification, error) {
	return u.notificationRepo.GetAll(ctx)
}

// NotifyDocumentCreated crea una notificación cuando se crea un documento
func (u *NotificationUsecase) NotifyDocumentCreated(ctx context.Context, userID, userName, documentID, documentTitle string) error {
	notification := entity.NewNotification(userID, userName, documentID, documentTitle)
	return u.CreateNotification(ctx, notification)
}
