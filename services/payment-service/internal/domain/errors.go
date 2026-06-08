package domain

import "errors"

// Payment errors
var (
	ErrPaymentNotFound                = errors.New("payment not found")
	ErrPaymentOrderRequired           = errors.New("order ID is required")
	ErrPaymentAmountInvalid           = errors.New("payment amount must be non-negative")
	ErrPaymentAmountTooHigh           = errors.New("payment amount exceeds maximum allowed")
	ErrPaymentTipInvalid              = errors.New("tip amount must be non-negative")
	ErrPaymentTotalInvalid            = errors.New("total amount must be non-negative")
	ErrPaymentMethodRequired          = errors.New("payment method is required")
	ErrPaymentRefundAmountInvalid     = errors.New("refund amount is invalid")
	ErrPaymentCustomerNameTooLong     = errors.New("customer name is too long (max 200 characters)")
	ErrPaymentNotesTooLong            = errors.New("notes are too long (max 500 characters)")
	ErrPaymentInvalidStatusTransition = errors.New("invalid payment status transition")
	ErrPaymentNotPending              = errors.New("payment is not in pending status")
	ErrPaymentNotProcessing           = errors.New("payment is not in processing status")
	ErrPaymentNotCompleted            = errors.New("payment is not completed")
	ErrPaymentCannotFail              = errors.New("payment cannot be failed in current status")
	ErrPaymentAlreadyCompleted        = errors.New("payment is already completed")
	ErrPaymentAlreadyRefunded         = errors.New("payment is already fully refunded")
	ErrPaymentOrderNotFound           = errors.New("order not found")
	ErrPaymentOrderNotValid           = errors.New("order is not valid for payment")
	ErrPaymentAmountMismatch          = errors.New("payment amount does not match order total")
)
