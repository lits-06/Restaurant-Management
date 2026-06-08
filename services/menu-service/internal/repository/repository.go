package repository

import (
	"context"

	"restaurant-management/services/menu-service/internal/domain"
)

// MenuItemRepository defines the interface for menu item data access
type MenuItemRepository interface {
	Create(ctx context.Context, item *domain.MenuItem) error
	GetByID(ctx context.Context, itemID string) (*domain.MenuItem, error)
	GetByName(ctx context.Context, name string) (*domain.MenuItem, error)
	Update(ctx context.Context, item *domain.MenuItem) error
	Delete(ctx context.Context, itemID string) error
	List(ctx context.Context, page, pageSize int, categoryID string, keyword string) ([]*domain.MenuItem, int, error)
	ListByCategory(ctx context.Context, categoryID string) ([]*domain.MenuItem, error)
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	GetByID(ctx context.Context, categoryID string) (*domain.Category, error)
	GetByName(ctx context.Context, name string) (*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, categoryID string) error
	// List(ctx context.Context, page, pageSize int) ([]*domain.Category, int, error)
	ListAll(ctx context.Context) ([]*domain.Category, error)
}
