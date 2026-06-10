package domain

import "errors"

// Order errors
var (
	ErrOrderNotFound              = errors.New("order not found")
	ErrOrderNameRequired          = errors.New("guest name is required")
	ErrOrderPhoneRequired         = errors.New("guest phone is required")
	ErrOrderTimeRequired          = errors.New("reservation time is required")
	ErrOrderDateRequired          = errors.New("reservation date is required")
	ErrOrderPartySizeInvalid      = errors.New("party size must be greater than 0")
	ErrOrderEndTimeInvalid        = errors.New("end_time must be after start time")
	ErrOrderStatusInvalid         = errors.New("invalid order status")
	ErrOrderInvalidStatusTransition = errors.New("invalid order status transition")
	ErrOrderAlreadyCancelled      = errors.New("order is already cancelled")
	ErrOrderCannotCancelCompleted = errors.New("cannot cancel completed order")
	ErrOrderCannotModify          = errors.New("order cannot be modified in current status")
)

// OrderItem errors
var (
	ErrOrderItemNotFound              = errors.New("order item not found")
	ErrOrderItemNameRequired          = errors.New("order item name is required")
	ErrOrderItemQuantityInvalid       = errors.New("quantity must be greater than 0")
	ErrOrderItemPriceInvalid          = errors.New("price must be non-negative")
	ErrOrderItemStatusInvalid         = errors.New("invalid item status")
	ErrOrderItemInvalidStatusTransition = errors.New("invalid item status transition")
)

// Table assignment errors
var (
	ErrNoTableAvailable = errors.New("no available table for the requested time slot and party size")
	ErrTableRequired    = errors.New("table_id is required when table service is unavailable")
)
