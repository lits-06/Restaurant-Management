package domain

import "errors"

// Order errors
var (
	ErrOrderNotFound                = errors.New("Order not found")
	ErrOrderNameRequired            = errors.New("Guest name is required")
	ErrOrderPhoneRequired           = errors.New("Guest phone is required")
	ErrOrderTimeRequired            = errors.New("Reservation time is required")
	ErrOrderDateRequired            = errors.New("Reservation date is required")
	ErrOrderPartySizeInvalid        = errors.New("Party size must be greater than 0")
	ErrOrderEndTimeInvalid          = errors.New("End time must be after start time")
	ErrOrderTimeOutsideHours        = errors.New("Start time must be between 10:00 and 22:00")
	ErrOrderEndTimeOutsideHours     = errors.New("End time must be 22:00 or earlier")
	ErrOrderStatusInvalid           = errors.New("Invalid order status")
	ErrOrderInvalidStatusTransition = errors.New("Invalid order status transition")
	ErrOrderAlreadyCancelled        = errors.New("Order is already cancelled")
	ErrOrderCannotCancelCompleted   = errors.New("Cannot cancel completed order")
	ErrOrderCannotModify            = errors.New("Order cannot be modified in current status")
)

// OrderItem errors
var (
	ErrOrderItemNotFound                = errors.New("Order item not found")
	ErrOrderItemNameRequired            = errors.New("Order item name is required")
	ErrOrderItemQuantityInvalid         = errors.New("Quantity must be greater than 0")
	ErrOrderItemPriceInvalid            = errors.New("Price must be non-negative")
	ErrOrderItemStatusInvalid           = errors.New("Invalid item status")
	ErrOrderItemInvalidStatusTransition = errors.New("Invalid item status transition")
)

// Table assignment errors
var (
	ErrNoTableAvailable = errors.New("No available table for the requested time slot and party size")
	ErrTableRequired    = errors.New("Table ID is required when table service is unavailable")
)
