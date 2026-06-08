package domain

import (
	"fmt"
	"time"
)

// PaymentMethod represents the payment method
type PaymentMethod int

const (
	MethodUnknown PaymentMethod = iota
	MethodCash
	MethodCreditCard
	MethodDebitCard
	MethodMobileWallet
	MethodBankTransfer
)

// String returns the string representation of PaymentMethod
func (m PaymentMethod) String() string {
	switch m {
	case MethodCash:
		return "cash"
	case MethodCreditCard:
		return "credit_card"
	case MethodDebitCard:
		return "debit_card"
	case MethodMobileWallet:
		return "mobile_wallet"
	case MethodBankTransfer:
		return "bank_transfer"
	default:
		return "unknown"
	}
}

// PaymentStatus represents the status of a payment
type PaymentStatus int

const (
	StatusUnknown PaymentStatus = iota
	StatusPending
	StatusProcessing
	StatusCompleted
	StatusFailed
	StatusRefunded
	StatusPartiallyRefunded
)

// String returns the string representation of PaymentStatus
func (s PaymentStatus) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusProcessing:
		return "processing"
	case StatusCompleted:
		return "completed"
	case StatusFailed:
		return "failed"
	case StatusRefunded:
		return "refunded"
	case StatusPartiallyRefunded:
		return "partially_refunded"
	default:
		return "unknown"
	}
}

// Payment represents a payment transaction
type Payment struct {
	PaymentID     string
	OrderID       string
	Amount        float64
	Tip           float64
	Total         float64
	Method        PaymentMethod
	Status        PaymentStatus
	TransactionID string
	CustomerName  string
	Notes         string
	RefundedAmount float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Validate validates the payment
func (p *Payment) Validate() error {
	if p.OrderID == "" {
		return ErrPaymentOrderRequired
	}

	if p.Amount < 0 {
		return ErrPaymentAmountInvalid
	}

	if p.Amount > 1000000000 { // 1 billion max
		return ErrPaymentAmountTooHigh
	}

	if p.Tip < 0 {
		return ErrPaymentTipInvalid
	}

	if p.Total < 0 {
		return ErrPaymentTotalInvalid
	}

	if p.Method == MethodUnknown {
		return ErrPaymentMethodRequired
	}

	if p.RefundedAmount < 0 || p.RefundedAmount > p.Total {
		return ErrPaymentRefundAmountInvalid
	}

	if len(p.CustomerName) > 200 {
		return ErrPaymentCustomerNameTooLong
	}

	if len(p.Notes) > 500 {
		return ErrPaymentNotesTooLong
	}

	return nil
}

// CalculateTotal calculates the total payment amount
func (p *Payment) CalculateTotal() {
	p.Total = p.Amount + p.Tip
	p.UpdatedAt = time.Now()
}

// CanTransitionTo checks if the payment can transition to the given status
func (p *Payment) CanTransitionTo(newStatus PaymentStatus) bool {
	// Define allowed transitions
	allowedTransitions := map[PaymentStatus][]PaymentStatus{
		StatusPending:            {StatusProcessing, StatusFailed},
		StatusProcessing:         {StatusCompleted, StatusFailed},
		StatusCompleted:          {StatusRefunded, StatusPartiallyRefunded},
		StatusFailed:             {StatusPending}, // can retry
		StatusPartiallyRefunded:  {StatusRefunded},
		StatusRefunded:           {}, // final state
	}

	allowed, exists := allowedTransitions[p.Status]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == newStatus {
			return true
		}
	}

	return false
}

// UpdateStatus updates the payment status
func (p *Payment) UpdateStatus(newStatus PaymentStatus) error {
	if !p.CanTransitionTo(newStatus) {
		return fmt.Errorf("%w: cannot transition from %s to %s",
			ErrPaymentInvalidStatusTransition, p.Status.String(), newStatus.String())
	}

	p.Status = newStatus
	p.UpdatedAt = time.Now()
	return nil
}

// Process marks the payment as processing
func (p *Payment) Process(transactionID string) error {
	if p.Status != StatusPending {
		return ErrPaymentNotPending
	}

	p.Status = StatusProcessing
	p.TransactionID = transactionID
	p.UpdatedAt = time.Now()
	return nil
}

// Complete marks the payment as completed
func (p *Payment) Complete() error {
	if p.Status != StatusProcessing {
		return ErrPaymentNotProcessing
	}

	p.Status = StatusCompleted
	p.UpdatedAt = time.Now()
	return nil
}

// Fail marks the payment as failed
func (p *Payment) Fail() error {
	if p.Status != StatusPending && p.Status != StatusProcessing {
		return ErrPaymentCannotFail
	}

	p.Status = StatusFailed
	p.UpdatedAt = time.Now()
	return nil
}

// Refund processes a full or partial refund
func (p *Payment) Refund(amount float64) error {
	if p.Status != StatusCompleted && p.Status != StatusPartiallyRefunded {
		return ErrPaymentNotCompleted
	}

	if amount <= 0 {
		return ErrPaymentRefundAmountInvalid
	}

	maxRefundable := p.Total - p.RefundedAmount
	if amount > maxRefundable {
		return fmt.Errorf("%w: max refundable is %.2f", ErrPaymentRefundAmountInvalid, maxRefundable)
	}

	p.RefundedAmount += amount

	// Update status based on refunded amount
	if p.RefundedAmount >= p.Total {
		p.Status = StatusRefunded
	} else {
		p.Status = StatusPartiallyRefunded
	}

	p.UpdatedAt = time.Now()
	return nil
}

// IsCompleted checks if the payment is completed
func (p *Payment) IsCompleted() bool {
	return p.Status == StatusCompleted
}

// IsPending checks if the payment is pending
func (p *Payment) IsPending() bool {
	return p.Status == StatusPending
}

// CanBeRefunded checks if the payment can be refunded
func (p *Payment) CanBeRefunded() bool {
	return p.Status == StatusCompleted || p.Status == StatusPartiallyRefunded
}

// GetRefundableAmount returns the amount that can still be refunded
func (p *Payment) GetRefundableAmount() float64 {
	return p.Total - p.RefundedAmount
}

// NewPayment creates a new payment
func NewPayment(orderID string, amount, tip float64, method PaymentMethod, customerName string) (*Payment, error) {
	now := time.Now()

	payment := &Payment{
		OrderID:      orderID,
		Amount:       amount,
		Tip:          tip,
		Method:       method,
		Status:       StatusPending,
		CustomerName: customerName,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	payment.CalculateTotal()

	if err := payment.Validate(); err != nil {
		return nil, err
	}

	return payment, nil
}
