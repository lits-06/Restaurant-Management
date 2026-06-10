package repository

import (
	"context"
	"restaurant-management/services/user-service/internal/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id string) (*domain.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*domain.User, error)

	// Update updates a user
	Update(ctx context.Context, user *domain.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id string) error

	// List retrieves users with filters and pagination
	List(ctx context.Context, filters ListFilters) ([]*domain.User, int, error)

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUsername checks if a user with the given username exists
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

// ListFilters represents filters for listing users
type ListFilters struct {
	Page     int
	PageSize int
	Status   domain.UserStatus
	Role     domain.UserRole
	Keyword  string // search in email, username, full_name
}
