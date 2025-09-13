package usecase

import (
	"context"

	"frontend-challenge/internal/domain/entity"
	"frontend-challenge/internal/domain/repository"
)

// NotificationUsecase defines the use cases for notifications
type NotificationUsecase struct {
	notificationRepo repository.NotificationRepository
	documentRepo     repository.DocumentRepository
	userRepo         repository.UserRepository
}

// NewNotificationUsecase creates a new instance of NotificationUsecase
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

// CreateNotification creates a new notification
func (u *NotificationUsecase) CreateNotification(ctx context.Context, notification *entity.Notification) error {
	if err := notification.Validate(); err != nil {
		return err
	}
	return u.notificationRepo.Create(ctx, notification)
}

// GetNotificationsByUserID gets notifications by user ID
func (u *NotificationUsecase) GetNotificationsByUserID(ctx context.Context, userID string) ([]*entity.Notification, error) {
	if userID == "" {
		return nil, entity.ErrInvalidUserID
	}
	return u.notificationRepo.GetByUserID(ctx, userID)
}

// GetAllNotifications gets all notifications
func (u *NotificationUsecase) GetAllNotifications(ctx context.Context) ([]*entity.Notification, error) {
	return u.notificationRepo.GetAll(ctx)
}

// NotifyDocumentCreated creates a notification when a document is created
func (u *NotificationUsecase) NotifyDocumentCreated(ctx context.Context, userID, userName, documentID, documentTitle string) error {
	notification := entity.NewNotification(userID, userName, documentID, documentTitle, "document.created")
	return u.CreateNotification(ctx, notification)
}
