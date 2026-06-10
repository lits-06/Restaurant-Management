package repository

import (
	"context"

	"restaurant-management/services/schedule-service/internal/domain"
)

type ShiftRepository interface {
	Create(ctx context.Context, shift *domain.Shift) error
	GetByID(ctx context.Context, id string) (*domain.Shift, error)
	Update(ctx context.Context, shift *domain.Shift) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, f ListFilters) ([]*domain.Shift, int, error)
}

type ListFilters struct {
	Month    string // "YYYY-MM"
	UserID   string
	Role     string
	Page     int
	PageSize int
}
