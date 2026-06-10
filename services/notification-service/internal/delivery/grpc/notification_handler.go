package grpc

import (
	"context"
	"time"

	"go.uber.org/zap"

	notifpb "restaurant-management/proto/notification"
	"restaurant-management/services/notification-service/internal/domain"
	"restaurant-management/services/notification-service/internal/usecase"
	"restaurant-management/shared/pkg/logger"
)

type NotificationHandler struct {
	notifpb.UnimplementedNotificationServiceServer
	uc *usecase.NotificationUseCase
}

func NewNotificationHandler(uc *usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{uc: uc}
}

func (h *NotificationHandler) SendNotification(ctx context.Context, req *notifpb.SendNotificationRequest) (*notifpb.SendNotificationResponse, error) {
	if req.TargetRole == "" {
		return &notifpb.SendNotificationResponse{Success: false, Message: "target_role is required"}, nil
	}
	if req.Type == "" {
		return &notifpb.SendNotificationResponse{Success: false, Message: "type is required"}, nil
	}

	items := make([]domain.OrderItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, domain.OrderItem{
			ItemID:   it.ItemId,
			ItemName: it.ItemName,
			Quantity: it.Quantity,
		})
	}

	n := &domain.Notification{
		Type:         req.Type,
		TargetRole:   req.TargetRole,
		OrderID:      req.OrderId,
		TableID:      req.TableId,
		ItemID:       req.ItemId,
		ItemName:     req.ItemName,
		Message:      req.Message,
		CustomerName: req.CustomerName,
		PartySize:    req.PartySize,
		Notes:        req.Notes,
		Items:        items,
		CreatedAt:    time.Now().Unix(),
	}

	if err := h.uc.Send(ctx, n); err != nil {
		logger.Error("failed to send notification", zap.Error(err))
		return &notifpb.SendNotificationResponse{Success: false, Message: err.Error()}, nil
	}

	logger.Info("notification sent",
		zap.String("type", req.Type),
		zap.String("target_role", req.TargetRole),
		zap.String("order_id", req.OrderId),
	)
	return &notifpb.SendNotificationResponse{Success: true, Message: "notification sent"}, nil
}

func (h *NotificationHandler) Subscribe(req *notifpb.SubscribeRequest, stream notifpb.NotificationService_SubscribeServer) error {
	if req.Role == "" {
		return nil
	}

	logger.Info("new subscriber", zap.String("role", req.Role))

	ctx := stream.Context()
	ch, err := h.uc.Subscribe(ctx, req.Role)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case n, ok := <-ch:
			if !ok {
				return nil
			}
			pbItems := make([]*notifpb.NotificationOrderItem, 0, len(n.Items))
			for _, it := range n.Items {
				pbItems = append(pbItems, &notifpb.NotificationOrderItem{
					ItemId:   it.ItemID,
					ItemName: it.ItemName,
					Quantity: it.Quantity,
				})
			}
			if err := stream.Send(&notifpb.Notification{
				Id:           n.ID,
				Type:         n.Type,
				TargetRole:   n.TargetRole,
				OrderId:      n.OrderID,
				TableId:      n.TableID,
				ItemId:       n.ItemID,
				ItemName:     n.ItemName,
				CreatedAt:    n.CreatedAt,
				Message:      n.Message,
				CustomerName: n.CustomerName,
				PartySize:    n.PartySize,
				Notes:        n.Notes,
				Items:        pbItems,
			}); err != nil {
				return err
			}
		}
	}
}
