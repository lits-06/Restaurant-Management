package repository

import (
	"context"

	"restaurant-management/services/notification-service/internal/domain"
)

// NotificationRepository defines the interface for notification data access
type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
	GetByID(ctx context.Context, notificationID string) (*domain.Notification, error)
	Update(ctx context.Context, notification *domain.Notification) error
	Delete(ctx context.Context, notificationID string) error
	ListByRecipient(ctx context.Context, recipientID string, page, pageSize int, notifType domain.NotificationType, status domain.NotificationStatus, unreadOnly bool) ([]*domain.Notification, int, int, error)
}
