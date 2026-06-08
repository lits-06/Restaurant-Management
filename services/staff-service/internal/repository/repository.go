package repository

import (
	"context"

	"restaurant-management/services/staff-service/internal/domain"
)

// StaffRepository defines the interface for staff data access.
type StaffRepository interface {
	Create(ctx context.Context, staff *domain.Staff) error
	GetByID(ctx context.Context, staffID string) (*domain.Staff, error)
	Update(ctx context.Context, staff *domain.Staff) error
	Delete(ctx context.Context, staffID string) error
	List(ctx context.Context, page, pageSize int, keyword string) ([]*domain.Staff, int, error)
}
