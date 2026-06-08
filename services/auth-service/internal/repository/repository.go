package repository

import (
	"context"
	"restaurant-management/services/auth-service/internal/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error

	// FindByID finds a user by ID
	FindByID(ctx context.Context, userID string) (*domain.User, error)

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *domain.User) error

	// Delete deletes a user
	Delete(ctx context.Context, userID string) error

	// Exists checks if a user with the given email exists
	Exists(ctx context.Context, email string) (bool, error)
}

// RefreshTokenRepository defines the interface for refresh token data access
type RefreshTokenRepository interface {
	// Create creates a new refresh token
	Create(ctx context.Context, token *domain.RefreshToken) error

	// FindByToken finds a refresh token by its value
	FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error)

	// Update updates a refresh token
	Update(ctx context.Context, token *domain.RefreshToken) error

	// Delete deletes a refresh token
	Delete(ctx context.Context, token string) error
}
