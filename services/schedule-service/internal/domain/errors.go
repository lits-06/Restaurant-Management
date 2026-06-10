package domain

import "errors"

var (
	ErrShiftNotFound      = errors.New("shift not found")
	ErrUserIDRequired     = errors.New("user_id is required")
	ErrInvalidDate        = errors.New("invalid date format, expected YYYY-MM-DD")
	ErrInvalidTime        = errors.New("invalid time format, expected HH:MM")
	ErrEndTimeBeforeStart = errors.New("end_time must be after start_time")
	ErrInvalidRole        = errors.New("role must be CHEF, WAITER, MANAGER, or ADMIN")
)
