package usecase

import (
	"context"
	"fmt"
	"time"

	"restaurant-management/services/payment-service/internal/domain"
	"restaurant-management/services/payment-service/internal/repository"
)

// OrderServiceClient defines the interface for order service operations
type OrderServiceClient interface {
	GetOrder(ctx context.Context, orderID string) (total float64, status string, err error)
	UpdateOrderToCompleted(ctx context.Context, orderID string) error
}

// PaymentUseCase handles payment business logic
type PaymentUseCase struct {
	paymentRepo repository.PaymentRepository
	orderClient OrderServiceClient
}

// NewPaymentUseCase creates a new payment use case
func NewPaymentUseCase(
	paymentRepo repository.PaymentRepository,
	orderClient OrderServiceClient,
) *PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo: paymentRepo,
		orderClient: orderClient,
	}
}

// CreatePayment creates a new payment
func (uc *PaymentUseCase) CreatePayment(ctx context.Context, orderID string, amount, tip float64, method domain.PaymentMethod, customerName, notes string) (*domain.Payment, error) {
	// Validate order exists and get order total
	orderTotal, orderStatus, err := uc.orderClient.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order validation failed: %w", err)
	}

	// Check if order is in a valid state for payment
	if orderStatus == "cancelled" {
		return nil, domain.ErrPaymentOrderNotValid
	}

	// Validate payment amount matches order total (with some tolerance for tip)
	if amount < orderTotal-0.01 { // allow small floating point differences
		return nil, fmt.Errorf("%w: order total is %.2f, payment amount is %.2f",
			domain.ErrPaymentAmountMismatch, orderTotal, amount)
	}

	// Create payment
	payment, err := domain.NewPayment(orderID, amount, tip, method, customerName)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}
	payment.Notes = notes

	// Save payment
	if err := uc.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	return payment, nil
}

// GetPayment retrieves a payment by ID
func (uc *PaymentUseCase) GetPayment(ctx context.Context, paymentID string) (*domain.Payment, error) {
	payment, err := uc.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	return payment, nil
}

// ProcessPayment processes a payment
func (uc *PaymentUseCase) ProcessPayment(ctx context.Context, paymentID, transactionID string) (*domain.Payment, error) {
	payment, err := uc.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Mark as processing
	if err := payment.Process(transactionID); err != nil {
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	// In real implementation, this would call payment gateway
	// For now, we'll immediately mark as completed for cash payments
	if payment.Method == domain.MethodCash {
		if err := payment.Complete(); err != nil {
			return nil, fmt.Errorf("failed to complete payment: %w", err)
		}

		// Update order status to completed
		if err := uc.orderClient.UpdateOrderToCompleted(ctx, payment.OrderID); err != nil {
			// Log error but don't fail the payment
			fmt.Printf("Warning: failed to update order status: %v\n", err)
		}
	} else {
		// For card/wallet payments, complete after processing
		if err := payment.Complete(); err != nil {
			return nil, fmt.Errorf("failed to complete payment: %w", err)
		}

		// Update order status to completed
		if err := uc.orderClient.UpdateOrderToCompleted(ctx, payment.OrderID); err != nil {
			fmt.Printf("Warning: failed to update order status: %v\n", err)
		}
	}

	// Save
	if err := uc.paymentRepo.Update(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to save payment: %w", err)
	}

	return payment, nil
}

// RefundPayment processes a refund
func (uc *PaymentUseCase) RefundPayment(ctx context.Context, paymentID string, amount float64, reason string) (*domain.Payment, error) {
	payment, err := uc.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	// Check if payment can be refunded
	if !payment.CanBeRefunded() {
		return nil, domain.ErrPaymentNotCompleted
	}

	// Process refund
	if err := payment.Refund(amount); err != nil {
		return nil, fmt.Errorf("failed to refund payment: %w", err)
	}

	// Add reason to notes
	if reason != "" {
		if payment.Notes != "" {
			payment.Notes += "; "
		}
		payment.Notes += "Refund: " + reason
	}

	// Save
	if err := uc.paymentRepo.Update(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to save refund: %w", err)
	}

	return payment, nil
}

// ListPayments retrieves payments with pagination and filters
func (uc *PaymentUseCase) ListPayments(ctx context.Context, page, pageSize int, status domain.PaymentStatus, method domain.PaymentMethod, fromDate, toDate *time.Time) ([]*domain.Payment, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	payments, total, err := uc.paymentRepo.List(ctx, page, pageSize, status, method, fromDate, toDate)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list payments: %w", err)
	}

	return payments, total, nil
}

// GetPaymentsByOrder retrieves all payments for an order
func (uc *PaymentUseCase) GetPaymentsByOrder(ctx context.Context, orderID string) ([]*domain.Payment, error) {
	// Validate order exists
	_, _, err := uc.orderClient.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order validation failed: %w", err)
	}

	payments, err := uc.paymentRepo.ListByOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by order: %w", err)
	}

	return payments, nil
}

// GetTotalPaid calculates the total amount paid for an order
func (uc *PaymentUseCase) GetTotalPaid(ctx context.Context, orderID string) (float64, error) {
	payments, err := uc.paymentRepo.ListByOrder(ctx, orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to get payments: %w", err)
	}

	var total float64
	for _, payment := range payments {
		if payment.IsCompleted() {
			total += payment.Total - payment.RefundedAmount
		}
	}

	return total, nil
}
