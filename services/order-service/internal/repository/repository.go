package repository

import (
	"context"
	"time"

	"restaurant-management/services/order-service/internal/domain"
)

// OrderRepository defines the interface for order data access
type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, orderID string) (*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, orderID string) error
	List(ctx context.Context, page, pageSize int, status domain.OrderStatus, keyword, userID, sortOrder string) ([]*domain.Order, int, error)
	// UpdateItemStatus updates the item_status of a single order_items row.
	UpdateItemStatus(ctx context.Context, orderID, itemID string, status domain.ItemStatus) error
	// GetOccupiedTableIDs returns table IDs that have a non-cancelled order
	// whose time window overlaps [startTime, endTime).
	GetOccupiedTableIDs(ctx context.Context, startTime, endTime time.Time) ([]string, error)
}
