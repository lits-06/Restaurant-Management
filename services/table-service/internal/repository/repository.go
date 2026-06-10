package repository

import (
	"context"

	"restaurant-management/services/table-service/internal/domain"
)

// TableRepository defines the interface for table data access.
type TableRepository interface {
	Create(ctx context.Context, table *domain.Table) error
	GetByID(ctx context.Context, id string) (*domain.Table, error)
	GetByTableNumber(ctx context.Context, tableNumber int) (*domain.Table, error)
	Update(ctx context.Context, table *domain.Table) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters ListFilters) ([]*domain.Table, int, error)
	GetAvailableTables(ctx context.Context, minCapacity int) ([]*domain.Table, error)
	ExistsByTableNumber(ctx context.Context, tableNumber int) (bool, error)
}

// ListFilters holds filter and pagination options for listing tables.
type ListFilters struct {
	Page     int
	PageSize int
	Status   domain.TableStatus
}
