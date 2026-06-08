package repository

import (
	"context"
	"restaurant-management/services/user-service/internal/domain"
	"sync"
)

// InMemoryUserRepository is an in-memory implementation of UserRepository
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User // key: user ID
	// Indexes for fast lookup
	emailIndex    map[string]string // email -> user ID
	usernameIndex map[string]string // username -> user ID
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:         make(map[string]*domain.User),
		emailIndex:    make(map[string]string),
		usernameIndex: make(map[string]string),
	}
}

// Create creates a new user
func (r *InMemoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user already exists
	if _, exists := r.users[user.ID]; exists {
		return domain.ErrUserAlreadyExists
	}

	// Check email uniqueness
	if _, exists := r.emailIndex[user.Email]; exists {
		return domain.ErrEmailAlreadyExists
	}

	// Check username uniqueness
	if _, exists := r.usernameIndex[user.Username]; exists {
		return domain.ErrUsernameAlreadyExists
	}

	// Store user
	r.users[user.ID] = user
	r.emailIndex[user.Email] = user.ID
	r.usernameIndex[user.Username] = user.ID

	return nil
}

// GetByID retrieves a user by ID
func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, exists := r.emailIndex[email]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	user, exists := r.users[userID]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (r *InMemoryUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userID, exists := r.usernameIndex[username]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	user, exists := r.users[userID]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

// Update updates a user
func (r *InMemoryUserRepository) Update(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user exists
	oldUser, exists := r.users[user.ID]
	if !exists {
		return domain.ErrUserNotFound
	}

	// If email changed, check uniqueness and update index
	if oldUser.Email != user.Email {
		if _, exists := r.emailIndex[user.Email]; exists {
			return domain.ErrEmailAlreadyExists
		}
		delete(r.emailIndex, oldUser.Email)
		r.emailIndex[user.Email] = user.ID
	}

	// If username changed, check uniqueness and update index
	if oldUser.Username != user.Username {
		if _, exists := r.usernameIndex[user.Username]; exists {
			return domain.ErrUsernameAlreadyExists
		}
		delete(r.usernameIndex, oldUser.Username)
		r.usernameIndex[user.Username] = user.ID
	}

	// Update user
	r.users[user.ID] = user

	return nil
}

// Delete deletes a user by ID
func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return domain.ErrUserNotFound
	}

	// Remove from indexes
	delete(r.emailIndex, user.Email)
	delete(r.usernameIndex, user.Username)

	// Remove user
	delete(r.users, id)

	return nil
}

// List retrieves users with filters and pagination
func (r *InMemoryUserRepository) List(ctx context.Context, filters ListFilters) ([]*domain.User, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Collect matching users
	var matchingUsers []*domain.User
	for _, user := range r.users {
		// Apply status filter
		if filters.Status != "" && user.Status != filters.Status {
			continue
		}

		// Apply role filter
		if filters.Role != "" {
			hasRole := false
			for _, role := range user.Roles {
				if role == filters.Role {
					hasRole = true
					break
				}
			}
			if !hasRole {
				continue
			}
		}

		matchingUsers = append(matchingUsers, user)
	}

	total := len(matchingUsers)

	// Apply pagination
	if filters.PageSize <= 0 {
		filters.PageSize = 10
	}
	if filters.Page <= 0 {
		filters.Page = 1
	}

	start := (filters.Page - 1) * filters.PageSize
	end := start + filters.PageSize

	if start >= total {
		return []*domain.User{}, total, nil
	}

	if end > total {
		end = total
	}

	return matchingUsers[start:end], total, nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *InMemoryUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.emailIndex[email]
	return exists, nil
}

// ExistsByUsername checks if a user with the given username exists
func (r *InMemoryUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.usernameIndex[username]
	return exists, nil
}
