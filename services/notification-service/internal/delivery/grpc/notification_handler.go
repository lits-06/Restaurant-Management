package grpc

import (
	"context"
	"errors"

	"restaurant-management/proto/notification"
	"restaurant-management/services/notification-service/internal/domain"
	"restaurant-management/services/notification-service/internal/usecase"
	"restaurant-management/shared/pkg/logger"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// NotificationHandler handles gRPC requests for notification service
type NotificationHandler struct {
	notification.UnimplementedNotificationServiceServer
	notificationUseCase *usecase.NotificationUseCase
	logger              logger.Logger
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationUseCase *usecase.NotificationUseCase, logger logger.Logger) *NotificationHandler {
	return &NotificationHandler{
		notificationUseCase: notificationUseCase,
		logger:              logger,
	}
}

// SendNotification sends a notification
func (h *NotificationHandler) SendNotification(ctx context.Context, req *notification.SendNotificationRequest) (*notification.SendNotificationResponse, error) {
	h.logger.Info("SendNotification request", "recipient_id", req.RecipientId, "type", req.Type)

	if req.RecipientId == "" {
		return &notification.SendNotificationResponse{
			Success: false,
			Message: "recipient_id is required",
		}, nil
	}

	if req.Type == notification.NotificationType_TYPE_UNKNOWN {
		return &notification.SendNotificationResponse{
			Success: false,
			Message: "notification type is required",
		}, nil
	}

	if req.Title == "" {
		return &notification.SendNotificationResponse{
			Success: false,
			Message: "title is required",
		}, nil
	}

	if req.Message == "" {
		return &notification.SendNotificationResponse{
			Success: false,
			Message: "message is required",
		}, nil
	}

	// Convert proto types to domain types
	domainType := convertProtoTypeToDomain(req.Type)
	domainPriority := convertProtoPriorityToDomain(req.Priority)

	// Send notification
	notif, err := h.notificationUseCase.SendNotification(
		ctx,
		req.RecipientId,
		req.RecipientEmail,
		req.RecipientPhone,
		domainType,
		domainPriority,
		req.Title,
		req.Message,
		req.Data,
	)

	if err != nil {
		h.logger.Error("Failed to send notification", "error", err)
		return &notification.SendNotificationResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &notification.SendNotificationResponse{
		Notification: convertNotificationToProto(notif),
		Success:      true,
		Message:      "notification sent successfully",
	}, nil
}

// SendBulkNotification sends notifications to multiple recipients
func (h *NotificationHandler) SendBulkNotification(ctx context.Context, req *notification.SendBulkNotificationRequest) (*notification.SendBulkNotificationResponse, error) {
	h.logger.Info("SendBulkNotification request", "recipient_count", len(req.RecipientIds))

	if len(req.RecipientIds) == 0 {
		return &notification.SendBulkNotificationResponse{
			Success: false,
			Message: "at least one recipient is required",
		}, nil
	}

	if req.Type == notification.NotificationType_TYPE_UNKNOWN {
		return &notification.SendBulkNotificationResponse{
			Success: false,
			Message: "notification type is required",
		}, nil
	}

	if req.Title == "" {
		return &notification.SendBulkNotificationResponse{
			Success: false,
			Message: "title is required",
		}, nil
	}

	if req.Message == "" {
		return &notification.SendBulkNotificationResponse{
			Success: false,
			Message: "message is required",
		}, nil
	}

	// Convert proto types to domain types
	domainType := convertProtoTypeToDomain(req.Type)
	domainPriority := convertProtoPriorityToDomain(req.Priority)

	// Send bulk notifications
	sentCount, failedCount, err := h.notificationUseCase.SendBulkNotification(
		ctx,
		req.RecipientIds,
		domainType,
		domainPriority,
		req.Title,
		req.Message,
		req.Data,
	)

	if err != nil {
		h.logger.Error("Failed to send bulk notifications", "error", err)
		return &notification.SendBulkNotificationResponse{
			SentCount:   int32(sentCount),
			FailedCount: int32(failedCount),
			Success:     false,
			Message:     err.Error(),
		}, nil
	}

	return &notification.SendBulkNotificationResponse{
		SentCount:   int32(sentCount),
		FailedCount: int32(failedCount),
		Success:     true,
		Message:     "bulk notifications processed",
	}, nil
}

// GetNotification retrieves a notification by ID
func (h *NotificationHandler) GetNotification(ctx context.Context, req *notification.GetNotificationRequest) (*notification.GetNotificationResponse, error) {
	h.logger.Info("GetNotification request", "notification_id", req.NotificationId)

	if req.NotificationId == "" {
		return &notification.GetNotificationResponse{
			Success: false,
			Message: "notification_id is required",
		}, nil
	}

	notif, err := h.notificationUseCase.GetNotification(ctx, req.NotificationId)
	if err != nil {
		if errors.Is(err, domain.ErrNotificationNotFound) {
			return &notification.GetNotificationResponse{
				Success: false,
				Message: "notification not found",
			}, nil
		}
		h.logger.Error("Failed to get notification", "error", err)
		return &notification.GetNotificationResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &notification.GetNotificationResponse{
		Notification: convertNotificationToProto(notif),
		Success:      true,
		Message:      "notification retrieved successfully",
	}, nil
}

// ListNotifications retrieves notifications for a recipient
func (h *NotificationHandler) ListNotifications(ctx context.Context, req *notification.ListNotificationsRequest) (*notification.ListNotificationsResponse, error) {
	h.logger.Info("ListNotifications request", "recipient_id", req.RecipientId)

	if req.RecipientId == "" {
		return &notification.ListNotificationsResponse{
			Success: false,
			Message: "recipient_id is required",
		}, nil
	}

	page := int(req.Page)
	pageSize := int(req.PageSize)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	// Convert proto types to domain types
	domainType := convertProtoTypeToDomain(req.Type)
	domainStatus := convertProtoStatusToDomain(req.Status)

	notifications, total, unreadCount, err := h.notificationUseCase.ListNotifications(
		ctx,
		req.RecipientId,
		page,
		pageSize,
		domainType,
		domainStatus,
		req.UnreadOnly,
	)
	if err != nil {
		h.logger.Error("Failed to list notifications", "error", err)
		return &notification.ListNotificationsResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	protoNotifications := make([]*notification.Notification, len(notifications))
	for i, notif := range notifications {
		protoNotifications[i] = convertNotificationToProto(notif)
	}

	return &notification.ListNotificationsResponse{
		Notifications: protoNotifications,
		Total:         int32(total),
		UnreadCount:   int32(unreadCount),
		Page:          int32(page),
		PageSize:      int32(pageSize),
		Success:       true,
		Message:       "notifications retrieved successfully",
	}, nil
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(ctx context.Context, req *notification.MarkAsReadRequest) (*notification.MarkAsReadResponse, error) {
	h.logger.Info("MarkAsRead request", "notification_id", req.NotificationId)

	if req.NotificationId == "" {
		return &notification.MarkAsReadResponse{
			Success: false,
			Message: "notification_id is required",
		}, nil
	}

	err := h.notificationUseCase.MarkAsRead(ctx, req.NotificationId)
	if err != nil {
		if errors.Is(err, domain.ErrNotificationNotFound) {
			return &notification.MarkAsReadResponse{
				Success: false,
				Message: "notification not found",
			}, nil
		}
		if errors.Is(err, domain.ErrNotificationAlreadyRead) {
			return &notification.MarkAsReadResponse{
				Success: false,
				Message: "notification is already read",
			}, nil
		}
		h.logger.Error("Failed to mark notification as read", "error", err)
		return &notification.MarkAsReadResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &notification.MarkAsReadResponse{
		Success: true,
		Message: "notification marked as read",
	}, nil
}

// DeleteNotification deletes a notification
func (h *NotificationHandler) DeleteNotification(ctx context.Context, req *notification.DeleteNotificationRequest) (*notification.DeleteNotificationResponse, error) {
	h.logger.Info("DeleteNotification request", "notification_id", req.NotificationId)

	if req.NotificationId == "" {
		return &notification.DeleteNotificationResponse{
			Success: false,
			Message: "notification_id is required",
		}, nil
	}

	err := h.notificationUseCase.DeleteNotification(ctx, req.NotificationId)
	if err != nil {
		if errors.Is(err, domain.ErrNotificationNotFound) {
			return &notification.DeleteNotificationResponse{
				Success: false,
				Message: "notification not found",
			}, nil
		}
		h.logger.Error("Failed to delete notification", "error", err)
		return &notification.DeleteNotificationResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &notification.DeleteNotificationResponse{
		Success: true,
		Message: "notification deleted successfully",
	}, nil
}

// Helper functions

func convertNotificationToProto(notif *domain.Notification) *notification.Notification {
	if notif == nil {
		return nil
	}

	protoNotif := &notification.Notification{
		NotificationId: notif.NotificationID,
		RecipientId:    notif.RecipientID,
		RecipientEmail: notif.RecipientEmail,
		RecipientPhone: notif.RecipientPhone,
		Type:           convertDomainTypeToProto(notif.Type),
		Priority:       convertDomainPriorityToProto(notif.Priority),
		Title:          notif.Title,
		Message:        notif.Message,
		Data:           notif.Data,
		Status:         convertDomainStatusToProto(notif.Status),
		CreatedAt:      timestamppb.New(notif.CreatedAt),
	}

	if notif.SentAt != nil {
		protoNotif.SentAt = timestamppb.New(*notif.SentAt)
	}

	if notif.ReadAt != nil {
		protoNotif.ReadAt = timestamppb.New(*notif.ReadAt)
	}

	return protoNotif
}

func convertDomainTypeToProto(notifType domain.NotificationType) notification.NotificationType {
	switch notifType {
	case domain.TypeEmail:
		return notification.NotificationType_TYPE_EMAIL
	case domain.TypeSMS:
		return notification.NotificationType_TYPE_SMS
	case domain.TypePush:
		return notification.NotificationType_TYPE_PUSH
	case domain.TypeInApp:
		return notification.NotificationType_TYPE_IN_APP
	default:
		return notification.NotificationType_TYPE_UNKNOWN
	}
}

func convertProtoTypeToDomain(notifType notification.NotificationType) domain.NotificationType {
	switch notifType {
	case notification.NotificationType_TYPE_EMAIL:
		return domain.TypeEmail
	case notification.NotificationType_TYPE_SMS:
		return domain.TypeSMS
	case notification.NotificationType_TYPE_PUSH:
		return domain.TypePush
	case notification.NotificationType_TYPE_IN_APP:
		return domain.TypeInApp
	default:
		return domain.TypeUnknown
	}
}

func convertDomainPriorityToProto(priority domain.NotificationPriority) notification.NotificationPriority {
	switch priority {
	case domain.PriorityLow:
		return notification.NotificationPriority_PRIORITY_LOW
	case domain.PriorityNormal:
		return notification.NotificationPriority_PRIORITY_NORMAL
	case domain.PriorityHigh:
		return notification.NotificationPriority_PRIORITY_HIGH
	case domain.PriorityUrgent:
		return notification.NotificationPriority_PRIORITY_URGENT
	default:
		return notification.NotificationPriority_PRIORITY_NORMAL
	}
}

func convertProtoPriorityToDomain(priority notification.NotificationPriority) domain.NotificationPriority {
	switch priority {
	case notification.NotificationPriority_PRIORITY_LOW:
		return domain.PriorityLow
	case notification.NotificationPriority_PRIORITY_NORMAL:
		return domain.PriorityNormal
	case notification.NotificationPriority_PRIORITY_HIGH:
		return domain.PriorityHigh
	case notification.NotificationPriority_PRIORITY_URGENT:
		return domain.PriorityUrgent
	default:
		return domain.PriorityNormal
	}
}

func convertDomainStatusToProto(status domain.NotificationStatus) notification.NotificationStatus {
	switch status {
	case domain.StatusPending:
		return notification.NotificationStatus_STATUS_PENDING
	case domain.StatusSent:
		return notification.NotificationStatus_STATUS_SENT
	case domain.StatusDelivered:
		return notification.NotificationStatus_STATUS_DELIVERED
	case domain.StatusFailed:
		return notification.NotificationStatus_STATUS_FAILED
	case domain.StatusRead:
		return notification.NotificationStatus_STATUS_READ
	default:
		return notification.NotificationStatus_STATUS_UNKNOWN
	}
}

func convertProtoStatusToDomain(status notification.NotificationStatus) domain.NotificationStatus {
	switch status {
	case notification.NotificationStatus_STATUS_PENDING:
		return domain.StatusPending
	case notification.NotificationStatus_STATUS_SENT:
		return domain.StatusSent
	case notification.NotificationStatus_STATUS_DELIVERED:
		return domain.StatusDelivered
	case notification.NotificationStatus_STATUS_FAILED:
		return domain.StatusFailed
	case notification.NotificationStatus_STATUS_READ:
		return domain.StatusRead
	default:
		return domain.StatusUnknown
	}
}
