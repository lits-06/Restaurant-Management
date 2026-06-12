package grpc

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"restaurant-management/proto/order"
	"restaurant-management/services/order-service/internal/domain"
	"restaurant-management/services/order-service/internal/usecase"
	"restaurant-management/shared/pkg/logger"
)

// OrderHandler handles gRPC requests for order service.
type OrderHandler struct {
	order.UnimplementedOrderServiceServer
	orderUseCase *usecase.OrderUseCase
}

// NewOrderHandler creates a new order handler.
func NewOrderHandler(orderUseCase *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{orderUseCase: orderUseCase}
}

// CreateOrder creates a new order.
func (h *OrderHandler) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	logger.Info("CreateOrder request",
		zap.String("name", req.Name),
		zap.String("phone", req.Phone),
	)

	if req.Name == "" {
		return &order.CreateOrderResponse{Success: false, Message: "name is required"}, nil
	}
	if req.Phone == "" {
		return &order.CreateOrderResponse{Success: false, Message: "phone is required"}, nil
	}
	if req.Time == "" {
		return &order.CreateOrderResponse{Success: false, Message: "time is required"}, nil
	}
	if req.Date == "" {
		return &order.CreateOrderResponse{Success: false, Message: "date is required"}, nil
	}
	if req.PartySize <= 0 {
		return &order.CreateOrderResponse{Success: false, Message: "party_size must be greater than 0"}, nil
	}

	itemRequests := make([]usecase.OrderItemRequest, 0, len(req.Items))
	for _, item := range req.Items {
		itemRequests = append(itemRequests, usecase.OrderItemRequest{
			MenuItemID: item.ItemId,
			Quantity:   item.Quantity,
		})
	}

	ord, err := h.orderUseCase.CreateOrder(ctx, req.TableId, req.UserId, req.Name, req.Phone, req.Notes, req.Time, req.EndTime, req.Date, req.PartySize, req.Status, itemRequests)
	if err != nil {
		logger.Error("Failed to create order", zap.Error(err))
		return &order.CreateOrderResponse{Success: false, Message: err.Error()}, nil
	}

	return &order.CreateOrderResponse{Order: convertOrderToProto(ord), Success: true, Message: "order created successfully"}, nil
}

// GetOrder retrieves an order by ID.
func (h *OrderHandler) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	logger.Info("GetOrder request", zap.String("order_id", req.OrderId))

	if req.OrderId == "" {
		return &order.GetOrderResponse{Success: false, Message: "order_id is required"}, nil
	}

	ord, err := h.orderUseCase.GetOrder(ctx, req.OrderId)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return &order.GetOrderResponse{Success: false, Message: "order not found"}, nil
		}
		logger.Error("Failed to get order", zap.Error(err))
		return &order.GetOrderResponse{Success: false, Message: err.Error()}, nil
	}

	return &order.GetOrderResponse{Order: convertOrderToProto(ord), Success: true, Message: "order retrieved successfully"}, nil
}

// UpdateOrder updates an order.
func (h *OrderHandler) UpdateOrder(ctx context.Context, req *order.UpdateOrderRequest) (*order.UpdateOrderResponse, error) {
	logger.Info("UpdateOrder request", zap.String("order_id", req.OrderId))

	if req.OrderId == "" {
		return &order.UpdateOrderResponse{Success: false, Message: "order_id is required"}, nil
	}

	var itemRequests []usecase.OrderItemRequest
	if req.Items != nil {
		itemRequests = make([]usecase.OrderItemRequest, 0, len(req.Items))
		for _, item := range req.Items {
			itemRequests = append(itemRequests, usecase.OrderItemRequest{
				MenuItemID: item.ItemId,
				Quantity:   item.Quantity,
			})
		}
	}

	ord, err := h.orderUseCase.UpdateOrder(ctx, req.OrderId, req.TableId, req.Name, req.Phone, req.Notes, req.Time, req.EndTime, req.Date, req.PartySize, req.Status, itemRequests)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return &order.UpdateOrderResponse{Success: false, Message: "order not found"}, nil
		}
		logger.Error("Failed to update order", zap.Error(err))
		return &order.UpdateOrderResponse{Success: false, Message: err.Error()}, nil
	}

	return &order.UpdateOrderResponse{Order: convertOrderToProto(ord), Success: true, Message: "order updated successfully"}, nil
}

// DeleteOrder deletes an order.
func (h *OrderHandler) DeleteOrder(ctx context.Context, req *order.DeleteOrderRequest) (*order.DeleteOrderResponse, error) {
	logger.Info("DeleteOrder request", zap.String("order_id", req.OrderId))

	if req.OrderId == "" {
		return &order.DeleteOrderResponse{Success: false, Message: "order_id is required"}, nil
	}

	if err := h.orderUseCase.DeleteOrder(ctx, req.OrderId); err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return &order.DeleteOrderResponse{Success: false, Message: "order not found"}, nil
		}
		logger.Error("Failed to delete order", zap.Error(err))
		return &order.DeleteOrderResponse{Success: false, Message: err.Error()}, nil
	}

	return &order.DeleteOrderResponse{Success: true, Message: "order deleted successfully"}, nil
}

// CancelOrder cancels an order.
func (h *OrderHandler) CancelOrder(ctx context.Context, req *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	logger.Info("CancelOrder request", zap.String("order_id", req.OrderId))

	if req.OrderId == "" {
		return &order.CancelOrderResponse{Success: false, Message: "order_id is required"}, nil
	}

	if err := h.orderUseCase.CancelOrder(ctx, req.OrderId); err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return &order.CancelOrderResponse{Success: false, Message: "order not found"}, nil
		}
		logger.Error("Failed to cancel order", zap.Error(err))
		return &order.CancelOrderResponse{Success: false, Message: err.Error()}, nil
	}

	return &order.CancelOrderResponse{Success: true, Message: "order cancelled successfully"}, nil
}

// ListOrders retrieves orders with pagination and filters.
func (h *OrderHandler) ListOrders(ctx context.Context, req *order.ListOrdersRequest) (*order.ListOrdersResponse, error) {
	logger.Info("ListOrders request",
		zap.Int32("page", req.Page),
		zap.Int32("page_size", req.PageSize),
	)

	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	orders, total, err := h.orderUseCase.ListOrders(ctx, page, pageSize, domain.OrderStatus(req.Status), req.Keyword, req.UserId, req.SortOrder)
	if err != nil {
		logger.Error("Failed to list orders", zap.Error(err))
		return &order.ListOrdersResponse{Success: false, Message: err.Error()}, nil
	}

	protoOrders := make([]*order.Order, len(orders))
	for i, ord := range orders {
		protoOrders[i] = convertOrderToProto(ord)
	}

	return &order.ListOrdersResponse{Orders: protoOrders, Total: int32(total), Page: int32(page), PageSize: int32(pageSize), Success: true, Message: "orders retrieved successfully"}, nil
}

// UpdateOrderStatus updates the status of an order.
func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, req *order.UpdateOrderStatusRequest) (*order.UpdateOrderStatusResponse, error) {
	logger.Info("UpdateOrderStatus request", zap.String("order_id", req.OrderId), zap.String("status", req.Status))

	if req.OrderId == "" {
		return &order.UpdateOrderStatusResponse{Success: false, Message: "order_id is required"}, nil
	}

	ord, err := h.orderUseCase.UpdateOrderStatus(ctx, req.OrderId, domain.OrderStatus(req.Status))
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return &order.UpdateOrderStatusResponse{Success: false, Message: "order not found"}, nil
		}
		logger.Error("Failed to update order status", zap.Error(err))
		return &order.UpdateOrderStatusResponse{Success: false, Message: err.Error()}, nil
	}

	return &order.UpdateOrderStatusResponse{Order: convertOrderToProto(ord), Success: true, Message: "order status updated successfully"}, nil
}

// AddOrderItem adds an item to an order.
func (h *OrderHandler) AddOrderItem(ctx context.Context, req *order.AddOrderItemRequest) (*order.AddOrderItemResponse, error) {
	logger.Info("AddOrderItem request", zap.String("order_id", req.OrderId))

	if req.OrderId == "" {
		return &order.AddOrderItemResponse{Success: false, Message: "order_id is required"}, nil
	}
	if req.Item == nil {
		return &order.AddOrderItemResponse{Success: false, Message: "item is required"}, nil
	}

	ord, err := h.orderUseCase.AddOrderItem(ctx, req.OrderId, usecase.OrderItemRequest{MenuItemID: req.Item.ItemId, Quantity: req.Item.Quantity})
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return &order.AddOrderItemResponse{Success: false, Message: "order not found"}, nil
		}
		logger.Error("Failed to add order item", zap.Error(err))
		return &order.AddOrderItemResponse{Success: false, Message: err.Error()}, nil
	}

	return &order.AddOrderItemResponse{Order: convertOrderToProto(ord), Success: true, Message: "item added successfully"}, nil
}

// RemoveOrderItem removes an item from an order.
func (h *OrderHandler) RemoveOrderItem(ctx context.Context, req *order.RemoveOrderItemRequest) (*order.RemoveOrderItemResponse, error) {
	logger.Info("RemoveOrderItem request", zap.String("order_id", req.OrderId), zap.String("item_id", req.ItemId))

	if req.OrderId == "" {
		return &order.RemoveOrderItemResponse{Success: false, Message: "order_id is required"}, nil
	}
	if req.ItemId == "" {
		return &order.RemoveOrderItemResponse{Success: false, Message: "item_id is required"}, nil
	}

	ord, err := h.orderUseCase.RemoveOrderItem(ctx, req.OrderId, req.ItemId)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return &order.RemoveOrderItemResponse{Success: false, Message: "order not found"}, nil
		}
		logger.Error("Failed to remove order item", zap.Error(err))
		return &order.RemoveOrderItemResponse{Success: false, Message: err.Error()}, nil
	}

	return &order.RemoveOrderItemResponse{Order: convertOrderToProto(ord), Success: true, Message: "item removed successfully"}, nil
}

func convertOrderToProto(ord *domain.Order) *order.Order {
	if ord == nil {
		return nil
	}

	items := make([]*order.OrderItem, len(ord.Items))
	for i, item := range ord.Items {
		items[i] = convertOrderItemToProto(item)
	}

	proto := &order.Order{
		OrderId:    ord.OrderID,
		TableId:    ord.TableID,
		UserId:     ord.UserID,
		Name:       ord.Name,
		Phone:      ord.Phone,
		Notes:      ord.Notes,
		Time:       timestamppb.New(ord.Time),
		PartySize:  ord.PartySize,
		Status:     string(ord.Status),
		TotalPrice: int32(ord.Total),
		Items:      items,
	}
	if !ord.EndTime.IsZero() {
		proto.EndTime = timestamppb.New(ord.EndTime)
	}
	return proto
}

// UpdateOrderItemStatus advances a single item's kitchen status.
func (h *OrderHandler) UpdateOrderItemStatus(ctx context.Context, req *order.UpdateOrderItemStatusRequest) (*order.UpdateOrderItemStatusResponse, error) {
	logger.Info("UpdateOrderItemStatus request",
		zap.String("order_id", req.OrderId),
		zap.String("item_id", req.ItemId),
		zap.String("item_status", req.ItemStatus),
	)

	if req.OrderId == "" {
		return &order.UpdateOrderItemStatusResponse{Success: false, Message: "order_id is required"}, nil
	}
	if req.ItemId == "" {
		return &order.UpdateOrderItemStatusResponse{Success: false, Message: "item_id is required"}, nil
	}

	next := domain.ItemStatus(req.ItemStatus)
	ord, err := h.orderUseCase.UpdateOrderItemStatus(ctx, req.OrderId, req.ItemId, next)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrOrderNotFound):
			return &order.UpdateOrderItemStatusResponse{Success: false, Message: "order not found"}, nil
		case errors.Is(err, domain.ErrOrderItemNotFound):
			return &order.UpdateOrderItemStatusResponse{Success: false, Message: "order item not found"}, nil
		case errors.Is(err, domain.ErrOrderItemStatusInvalid):
			return &order.UpdateOrderItemStatusResponse{Success: false, Message: "invalid item status"}, nil
		case errors.Is(err, domain.ErrOrderItemInvalidStatusTransition):
			return &order.UpdateOrderItemStatusResponse{Success: false, Message: "invalid item status transition"}, nil
		}
		logger.Error("Failed to update order item status", zap.Error(err))
		return &order.UpdateOrderItemStatusResponse{Success: false, Message: err.Error()}, nil
	}

	return &order.UpdateOrderItemStatusResponse{Order: convertOrderToProto(ord), Success: true, Message: "item status updated"}, nil
}

func convertOrderItemToProto(item *domain.OrderItem) *order.OrderItem {
	if item == nil {
		return nil
	}

	status := string(item.ItemStatus)
	if status == "" {
		status = string(domain.ItemStatusPending)
	}

	return &order.OrderItem{
		ItemId:     item.ItemID,
		Name:       item.Name,
		Price:      item.Price,
		Quantity:   item.Quantity,
		ItemStatus: status,
	}
}
