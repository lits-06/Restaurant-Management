package domain

import "errors"

var (
	ErrInvalidTableNumber       = errors.New("table number must be a positive integer")
	ErrInvalidCapacity          = errors.New("capacity must be greater than 0")
	ErrCapacityTooLarge         = errors.New("capacity cannot exceed 50")
	ErrInvalidStatus            = errors.New("invalid table status")
	ErrInvalidTableID           = errors.New("table_id is required")

	ErrTableNotFound            = errors.New("table not found")
	ErrTableNumberAlreadyExists = errors.New("table number already exists")
	ErrTableNotAvailable        = errors.New("table is not available for this transition")
	ErrTableOutOfService        = errors.New("table is out of service")
)
