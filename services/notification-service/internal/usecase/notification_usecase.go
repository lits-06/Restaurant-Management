package usecase

import (
	"context"
	"fmt"

	"restaurant-management/services/notification-service/internal/domain"
	"restaurant-management/services/notification-service/internal/repository"
)

// NotificationSender defines the interface for sending notifications
type NotificationSender interface {
	SendEmail(recipientEmail, title, message string) error
	SendSMS(recipientPhone, message string) error
	SendPush(recipientID, title, message, data string) error
}

// NotificationUseCase handles notification business logic
type NotificationUseCase struct {
	notificationRepo repository.NotificationRepository
	sender           NotificationSender
}

// NewNotificationUseCase creates a new notification use case
func NewNotificationUseCase(
	notificationRepo repository.NotificationRepository,
	sender NotificationSender,
) *NotificationUseCase {
	return &NotificationUseCase{
		notificationRepo: notificationRepo,
		sender:           sender,
	}
}

// SendNotification sends a notification to a recipient
func (uc *NotificationUseCase) SendNotification(
	ctx context.Context,
	recipientID, recipientEmail, recipientPhone string,
	notifType domain.NotificationType,
	priority domain.NotificationPriority,
	title, message, data string,
) (*domain.Notification, error) {
	// Create notification
	notification, err := domain.NewNotification(
		recipientID,
		recipientEmail,
		recipientPhone,
		notifType,
		priority,
		title,
		message,
		data,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Save notification
	if err := uc.notificationRepo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to save notification: %w", err)
	}

	// Send notification
	if err := uc.sendNotification(notification); err != nil {
		notification.MarkAsFailed()
		_ = uc.notificationRepo.Update(ctx, notification)
		return notification, fmt.Errorf("failed to send notification: %w", err)
	}

	// Mark as sent
	notification.MarkAsSent()
	_ = uc.notificationRepo.Update(ctx, notification)

	return notification, nil
}

// SendBulkNotification sends notifications to multiple recipients
func (uc *NotificationUseCase) SendBulkNotification(
	ctx context.Context,
	recipientIDs []string,
	notifType domain.NotificationType,
	priority domain.NotificationPriority,
	title, message, data string,
) (sentCount, failedCount int, err error) {
	for _, recipientID := range recipientIDs {
		// For bulk notifications, we don't have individual email/phone
		// In real implementation, you'd fetch user details from User Service
		_, sendErr := uc.SendNotification(
			ctx,
			recipientID,
			"", // Would fetch from user service
			"", // Would fetch from user service
			notifType,
			priority,
			title,
			message,
			data,
		)

		if sendErr != nil {
			failedCount++
		} else {
			sentCount++
		}
	}

	return sentCount, failedCount, nil
}

// GetNotification retrieves a notification by ID
func (uc *NotificationUseCase) GetNotification(ctx context.Context, notificationID string) (*domain.Notification, error) {
	notification, err := uc.notificationRepo.GetByID(ctx, notificationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	return notification, nil
}

// ListNotifications retrieves notifications for a recipient with filters
func (uc *NotificationUseCase) ListNotifications(
	ctx context.Context,
	recipientID string,
	page, pageSize int,
	notifType domain.NotificationType,
	status domain.NotificationStatus,
	unreadOnly bool,
) ([]*domain.Notification, int, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	notifications, total, unreadCount, err := uc.notificationRepo.ListByRecipient(
		ctx,
		recipientID,
		page,
		pageSize,
		notifType,
		status,
		unreadOnly,
	)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to list notifications: %w", err)
	}

	return notifications, total, unreadCount, nil
}

// MarkAsRead marks a notification as read
func (uc *NotificationUseCase) MarkAsRead(ctx context.Context, notificationID string) error {
	notification, err := uc.notificationRepo.GetByID(ctx, notificationID)
	if err != nil {
		return fmt.Errorf("failed to get notification: %w", err)
	}

	if notification.IsRead() {
		return domain.ErrNotificationAlreadyRead
	}

	notification.MarkAsRead()

	if err := uc.notificationRepo.Update(ctx, notification); err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	return nil
}

// DeleteNotification deletes a notification
func (uc *NotificationUseCase) DeleteNotification(ctx context.Context, notificationID string) error {
	if err := uc.notificationRepo.Delete(ctx, notificationID); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

// sendNotification sends a notification via the appropriate channel
func (uc *NotificationUseCase) sendNotification(notification *domain.Notification) error {
	switch notification.Type {
	case domain.TypeEmail:
		return uc.sender.SendEmail(notification.RecipientEmail, notification.Title, notification.Message)
	case domain.TypeSMS:
		return uc.sender.SendSMS(notification.RecipientPhone, notification.Message)
	case domain.TypePush:
		return uc.sender.SendPush(notification.RecipientID, notification.Title, notification.Message, notification.Data)
	case domain.TypeInApp:
		// In-app notifications are just stored in DB, no external sending needed
		return nil
	default:
		return fmt.Errorf("unsupported notification type: %v", notification.Type)
	}
}
