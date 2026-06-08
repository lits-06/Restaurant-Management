package domain

import "errors"

// Domain errors for Staff entity
var (
	ErrInvalidStaffID      = errors.New("invalid staff id")
	ErrInvalidStaffName    = errors.New("invalid staff name")
	ErrStaffNameTooShort   = errors.New("staff name must be at least 2 characters")
	ErrStaffNameTooLong    = errors.New("staff name must be at most 120 characters")
	ErrInvalidStaffRole    = errors.New("invalid staff role")
	ErrInvalidStaffContact = errors.New("invalid staff contact")
	ErrStaffContactTooLong = errors.New("staff contact must be at most 100 characters")
	ErrStaffUntilTooLong   = errors.New("staff shift time must be at most 20 characters")
	ErrStaffAvatarTooLong  = errors.New("staff avatar URL must be at most 500 characters")
	ErrStaffNotFound       = errors.New("staff not found")
)
