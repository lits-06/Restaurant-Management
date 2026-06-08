package domain

import "errors"

// Domain errors
var (
	// User errors
	ErrInvalidUserID   = errors.New("invalid user id")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound    = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Token errors
	ErrInvalidTokenID    = errors.New("invalid token id")
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidExpiryTime = errors.New("invalid expiry time")
	ErrTokenExpired      = errors.New("token has expired")
	ErrTokenRevoked      = errors.New("token has been revoked")
	ErrTokenNotFound     = errors.New("token not found")

	// Authentication errors
	ErrUnauthorized        = errors.New("unauthorized")
	ErrPasswordMismatch    = errors.New("password does not match")
	ErrWeakPassword        = errors.New("password is too weak")
)
