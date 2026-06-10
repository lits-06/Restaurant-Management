package domain

import "errors"

// Domain errors for User entity
var (
	// User validation errors
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrInvalidUsername    = errors.New("invalid username")
	ErrUsernameTooShort   = errors.New("username must be at least 3 characters")
	ErrInvalidFullName    = errors.New("invalid full name")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidPhone       = errors.New("invalid phone number")

	// User status errors
	ErrInvalidStatus = errors.New("invalid user status")

	// Role errors
	ErrInvalidRole          = errors.New("invalid role")
	ErrNoRolesAssigned      = errors.New("at least one role must be assigned")
	ErrRoleAlreadyAssigned  = errors.New("role already assigned to user")
	ErrRoleNotFound         = errors.New("role not found for user")
	ErrCannotRemoveLastRole = errors.New("cannot remove last role from user")

	// User not found
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrAccountSuspended    = errors.New("account is suspended")

	// Email/Username uniqueness
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)
