package domain

import "errors"

// Domain errors for Table entity
var (
	// Table validation errors
	ErrInvalidTableNumber = errors.New("invalid table number")
	ErrInvalidCapacity    = errors.New("invalid capacity")
	ErrCapacityTooLarge   = errors.New("capacity cannot exceed 50 people")
	ErrInvalidLocation    = errors.New("invalid location")
	ErrInvalidStatus      = errors.New("invalid table status")
	ErrInvalidOrderID     = errors.New("invalid order ID")
	ErrInvalidTableID     = errors.New("invalid table ID")

	// Reservation validation errors
	ErrInvalidReservationTime     = errors.New("invalid reservation time range")
	ErrInvalidReservationStatus   = errors.New("invalid reservation status")
	ErrInvalidReservationItem     = errors.New("invalid reservation item")
	ErrInvalidReservationCustomer = errors.New("invalid reservation customer")

	// Table not found
	ErrTableNotFound      = errors.New("table not found")
	ErrTableAlreadyExists = errors.New("table already exists")

	// Table number uniqueness
	ErrTableNumberAlreadyExists = errors.New("table number already exists")

	// Table status errors
	ErrTableNotAvailable    = errors.New("table is not available")
	ErrTableOutOfService    = errors.New("table is out of service")
	ErrTableAlreadyOccupied = errors.New("table is already occupied")

	// Reservation errors
	ErrReservationNotFound         = errors.New("reservation not found")
	ErrReservationConflict         = errors.New("reservation time conflict")
	ErrReservationAlreadyCancelled = errors.New("reservation already cancelled")
	ErrReservationAlreadyCompleted = errors.New("reservation already completed")
)
