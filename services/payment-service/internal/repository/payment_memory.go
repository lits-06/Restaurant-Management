package repository

import (
	"context"
	"sort"
	"sync"
	"time"

	"restaurant-management/services/payment-service/internal/domain"

	"github.com/google/uuid"
)

// InMemoryPaymentRepository is an in-memory implementation of PaymentRepository
type InMemoryPaymentRepository struct {
	mu         sync.RWMutex
	payments   map[string]*domain.Payment
	orderIndex map[string][]string // orderID -> []paymentID
}

// NewInMemoryPaymentRepository creates a new in-memory payment repository
func NewInMemoryPaymentRepository() *InMemoryPaymentRepository {
	return &InMemoryPaymentRepository{
		payments:   make(map[string]*domain.Payment),
		orderIndex: make(map[string][]string),
	}
}

// Create creates a new payment
func (r *InMemoryPaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Generate ID if not set
	if payment.PaymentID == "" {
		payment.PaymentID = uuid.New().String()
	}

	// Store payment
	r.payments[payment.PaymentID] = payment

	// Update order index
	r.orderIndex[payment.OrderID] = append(r.orderIndex[payment.OrderID], payment.PaymentID)

	return nil
}

// GetByID retrieves a payment by ID
func (r *InMemoryPaymentRepository) GetByID(ctx context.Context, paymentID string) (*domain.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	payment, exists := r.payments[paymentID]
	if !exists {
		return nil, domain.ErrPaymentNotFound
	}

	return payment, nil
}

// Update updates a payment
func (r *InMemoryPaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.payments[payment.PaymentID]
	if !exists {
		return domain.ErrPaymentNotFound
	}

	// Update order index if order changed (unlikely but handle it)
	if existing.OrderID != payment.OrderID {
		// Remove from old order
		oldPayments := r.orderIndex[existing.OrderID]
		for i, id := range oldPayments {
			if id == payment.PaymentID {
				r.orderIndex[existing.OrderID] = append(oldPayments[:i], oldPayments[i+1:]...)
				break
			}
		}
		// Add to new order
		r.orderIndex[payment.OrderID] = append(r.orderIndex[payment.OrderID], payment.PaymentID)
	}

	r.payments[payment.PaymentID] = payment
	return nil
}

// Delete deletes a payment
func (r *InMemoryPaymentRepository) Delete(ctx context.Context, paymentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	payment, exists := r.payments[paymentID]
	if !exists {
		return domain.ErrPaymentNotFound
	}

	// Remove from order index
	orderPayments := r.orderIndex[payment.OrderID]
	for i, id := range orderPayments {
		if id == paymentID {
			r.orderIndex[payment.OrderID] = append(orderPayments[:i], orderPayments[i+1:]...)
			break
		}
	}

	delete(r.payments, paymentID)
	return nil
}

// List retrieves payments with pagination and filters
func (r *InMemoryPaymentRepository) List(ctx context.Context, page, pageSize int, status domain.PaymentStatus, method domain.PaymentMethod, fromDate, toDate *time.Time) ([]*domain.Payment, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Collect all payments that match filters
	var filtered []*domain.Payment
	for _, payment := range r.payments {
		// Filter by status if specified
		if status != domain.StatusUnknown && payment.Status != status {
			continue
		}
		// Filter by method if specified
		if method != domain.MethodUnknown && payment.Method != method {
			continue
		}
		// Filter by date range if specified
		if fromDate != nil && payment.CreatedAt.Before(*fromDate) {
			continue
		}
		if toDate != nil && payment.CreatedAt.After(*toDate) {
			continue
		}
		filtered = append(filtered, payment)
	}

	// Sort by created_at descending (newest first)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	total := len(filtered)

	// Apply pagination
	start := (page - 1) * pageSize
	if start >= total {
		return []*domain.Payment{}, total, nil
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

// ListByOrder retrieves all payments for an order
func (r *InMemoryPaymentRepository) ListByOrder(ctx context.Context, orderID string) ([]*domain.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	paymentIDs := r.orderIndex[orderID]
	payments := make([]*domain.Payment, 0)

	for _, paymentID := range paymentIDs {
		if payment, exists := r.payments[paymentID]; exists {
			payments = append(payments, payment)
		}
	}

	// Sort by created_at descending
	sort.Slice(payments, func(i, j int) bool {
		return payments[i].CreatedAt.After(payments[j].CreatedAt)
	})

	return payments, nil
}
