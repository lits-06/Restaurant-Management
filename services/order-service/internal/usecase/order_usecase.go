package usecase

import (
	"context"
	"fmt"

	"restaurant-management/proto/menu"
	"restaurant-management/services/order-service/internal/domain"
	"restaurant-management/services/order-service/internal/repository"
)

// MenuServiceClient defines the interface for menu service operations
type MenuServiceClient interface {
	ValidateMenuItem(ctx context.Context, itemID string) (*menu.MenuItem, error)
}

// OrderUseCase handles order business logic
type OrderUseCase struct {
	orderRepo  repository.OrderRepository
	menuClient MenuServiceClient
}

// NewOrderUseCase creates a new order use case
func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	menuClient MenuServiceClient,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:  orderRepo,
		menuClient: menuClient,
	}
}

// OrderItemRequest represents a request to add an order item
type OrderItemRequest struct {
	MenuItemID string
	Quantity   int32
	// Notes      string
}

// CreateOrder creates a new order
func (uc *OrderUseCase) CreateOrder(ctx context.Context, name, phone, time, date string, partySize int32, status string, itemRequests []OrderItemRequest) (*domain.Order, error) {
	// Create order items with menu validation
	items := make([]*domain.OrderItem, 0, len(itemRequests))
	for _, req := range itemRequests {
		// Validate menu item and get price
		item, err := uc.menuClient.ValidateMenuItem(ctx, req.MenuItemID)
		if err != nil {
			return nil, fmt.Errorf("menu item validation failed for %s: %w", req.MenuItemID, err)
		}

		oItem, err := domain.NewOrderItem(item.ItemId, item.Name, item.Price, req.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
		// item.Notes = req.Notes
		items = append(items, oItem)
	}

	// Create order
	order, err := domain.NewOrder(name, phone, time, date, partySize, domain.OrderStatus(status), items)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	// order.Notes = notes

	// Save order
	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}

// GetOrder retrieves an order by ID
func (uc *OrderUseCase) GetOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

// UpdateOrder updates an order
func (uc *OrderUseCase) UpdateOrder(ctx context.Context, orderID string, name, phone, time, date string, partySize int32, status string, itemRequests []OrderItemRequest) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Update basic info
	if name != "" {
		order.Name = name
	}
	if phone != "" {
		order.Phone = phone
	}
	if time != "" || date != "" {
		order.Time, _ = domain.ParseReservationTime(date, time)
	}
	if partySize > 0 {
		order.PartySize = partySize
	}
	if status != "" {
		order.Status = domain.OrderStatus(status)
	}

	// Update items if provided
	if itemRequests != nil {
		// Clear existing items
		order.Items = make([]*domain.OrderItem, 0)

		// Add new items
		for _, req := range itemRequests {
			item, err := uc.menuClient.ValidateMenuItem(ctx, req.MenuItemID)
			if err != nil {
				return nil, fmt.Errorf("menu item validation failed: %w", err)
			}

			oItem, err := domain.NewOrderItem(item.ItemId, item.Name, item.Price, req.Quantity)
			if err != nil {
				return nil, fmt.Errorf("failed to create order item: %w", err)
			}
			// oItem.Notes = req.Notes
			order.Items = append(order.Items, oItem)
		}
	}

	// Update notes if provided
	// if notes != "" {
	// 	order.Notes = notes
	// }

	// Recalculate total
	order.Total = order.TotalPrice()
	// Save
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	return order, nil
}

// CancelOrder cancels an order
func (uc *OrderUseCase) CancelOrder(ctx context.Context, orderID string) error {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Cancel order
	if err := order.Cancel(); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	// Save
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to save cancelled order: %w", err)
	}

	return nil
}

// ListOrders retrieves orders with pagination and filters
func (uc *OrderUseCase) ListOrders(ctx context.Context, page, pageSize int, status domain.OrderStatus, keyword string) ([]*domain.Order, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orders, total, err := uc.orderRepo.List(ctx, page, pageSize, status, keyword)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}

	return orders, total, nil
}

// UpdateOrderStatus updates the status of an order
func (uc *OrderUseCase) UpdateOrderStatus(ctx context.Context, orderID string, newStatus domain.OrderStatus) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Update status
	if err := order.UpdateStatus(newStatus); err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// Save
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}

// AddOrderItem adds an item to an order
func (uc *OrderUseCase) AddOrderItem(ctx context.Context, orderID string, req OrderItemRequest) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Validate menu item
	item, err := uc.menuClient.ValidateMenuItem(ctx, req.MenuItemID)
	if err != nil {
		return nil, fmt.Errorf("menu item validation failed: %w", err)
	}

	// Create item
	oItem, err := domain.NewOrderItem(item.ItemId, item.Name, item.Price, req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to create order item: %w", err)
	}
	// oItem.Notes = req.Notes

	// Add to order
	if err := order.AddItem(oItem); err != nil {
		return nil, fmt.Errorf("failed to add item: %w", err)
	}

	// Save
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}

// RemoveOrderItem removes an item from an order
func (uc *OrderUseCase) RemoveOrderItem(ctx context.Context, orderID, itemID string) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Remove item
	if err := order.RemoveItem(itemID); err != nil {
		return nil, fmt.Errorf("failed to remove item: %w", err)
	}

	// Check if order still has items
	if len(order.Items) == 0 {
		return nil, fmt.Errorf("cannot remove last item from order")
	}

	// Save
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}
