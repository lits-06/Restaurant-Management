package repository

import (
	"context"
	"time"

	"restaurant-management/services/payment-service/internal/domain"
)

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	GetByID(ctx context.Context, paymentID string) (*domain.Payment, error)
	Update(ctx context.Context, payment *domain.Payment) error
	Delete(ctx context.Context, paymentID string) error
	List(ctx context.Context, page, pageSize int, status domain.PaymentStatus, method domain.PaymentMethod, fromDate, toDate *time.Time) ([]*domain.Payment, int, error)
	ListByOrder(ctx context.Context, orderID string) ([]*domain.Payment, error)
}
