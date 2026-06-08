package errors

import (
	"fmt"
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	// Authentication & Authorization
	ErrCodeUnauthorized     ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden        ErrorCode = "FORBIDDEN"
	ErrCodeInvalidToken     ErrorCode = "INVALID_TOKEN"
	ErrCodeExpiredToken     ErrorCode = "EXPIRED_TOKEN"
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"

	// Validation
	ErrCodeValidation       ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidInput     ErrorCode = "INVALID_INPUT"
	ErrCodeMissingField     ErrorCode = "MISSING_FIELD"

	// Resource
	ErrCodeNotFound         ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists    ErrorCode = "ALREADY_EXISTS"
	ErrCodeConflict         ErrorCode = "CONFLICT"

	// Business Logic
	ErrCodeInsufficientStock ErrorCode = "INSUFFICIENT_STOCK"
	ErrCodeTableOccupied    ErrorCode = "TABLE_OCCUPIED"
	ErrCodeOrderNotFound    ErrorCode = "ORDER_NOT_FOUND"
	ErrCodePaymentFailed    ErrorCode = "PAYMENT_FAILED"

	// System
	ErrCodeInternal         ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabase         ErrorCode = "DATABASE_ERROR"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
)

// AppError represents an application error
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with an AppError
func Wrap(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error constructors

func Unauthorized(message string) *AppError {
	return New(ErrCodeUnauthorized, message)
}

func Forbidden(message string) *AppError {
	return New(ErrCodeForbidden, message)
}

func NotFound(resource string) *AppError {
	return New(ErrCodeNotFound, fmt.Sprintf("%s not found", resource))
}

func AlreadyExists(resource string) *AppError {
	return New(ErrCodeAlreadyExists, fmt.Sprintf("%s already exists", resource))
}

func ValidationError(message string) *AppError {
	return New(ErrCodeValidation, message)
}

func InvalidInput(field string) *AppError {
	return New(ErrCodeInvalidInput, fmt.Sprintf("invalid input for field: %s", field))
}

func InternalError(message string, err error) *AppError {
	return Wrap(ErrCodeInternal, message, err)
}

func DatabaseError(operation string, err error) *AppError {
	return Wrap(ErrCodeDatabase, fmt.Sprintf("database error during %s", operation), err)
}
