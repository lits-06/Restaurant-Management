package usecase

import (
	"context"
	"fmt"
	"restaurant-management/services/user-service/internal/domain"
	"restaurant-management/services/user-service/internal/repository"

	"github.com/google/uuid"
)

// PasswordHasher defines the interface for password hashing
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

// UserUseCase handles user business logic
type UserUseCase struct {
	userRepo       repository.UserRepository
	passwordHasher PasswordHasher
}

// NewUserUseCase creates a new UserUseCase
func NewUserUseCase(userRepo repository.UserRepository, passwordHasher PasswordHasher) *UserUseCase {
	return &UserUseCase{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(ctx context.Context, email, username, fullName, phone, password string, roles []domain.UserRole) (*domain.User, error) {
	// Check if email already exists
	exists, err := uc.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, domain.ErrEmailAlreadyExists
	}

	// Check if username already exists
	exists, err = uc.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, domain.ErrUsernameAlreadyExists
	}

	// Hash password
	hashedPassword, err := uc.passwordHasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Set default role if none provided
	if len(roles) == 0 {
		roles = []domain.UserRole{domain.RoleWaiter} // Default role
	}

	// Create user entity
	user, err := domain.NewUser(email, username, fullName, phone, hashedPassword, roles)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	// Generate ID
	user.ID = uuid.New().String()

	// Save to repository
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (uc *UserUseCase) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	if userID == "" {
		return nil, domain.ErrInvalidUsername
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (uc *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if email == "" {
		return nil, domain.ErrInvalidEmail
	}

	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (uc *UserUseCase) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	if username == "" {
		return nil, domain.ErrInvalidUsername
	}

	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

// UpdateUser updates user information
func (uc *UserUseCase) UpdateUser(ctx context.Context, userID, email, username, fullName, phone string, status domain.UserStatus) (*domain.User, error) {
	// Get existing user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields
	if err := user.Update(email, username, fullName, phone); err != nil {
		return nil, fmt.Errorf("failed to update user fields: %w", err)
	}

	// Update status if provided
	if status != "" {
		user.Status = status
	}

	// Save changes
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser deletes a user
func (uc *UserUseCase) DeleteUser(ctx context.Context, userID string) error {
	if userID == "" {
		return domain.ErrInvalidUsername
	}

	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Delete user
	if err := uc.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers retrieves users with filters and pagination
func (uc *UserUseCase) ListUsers(ctx context.Context, page, pageSize int, status domain.UserStatus, role domain.UserRole) ([]*domain.User, int, error) {
	filters := repository.ListFilters{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		Role:     role,
	}

	users, total, err := uc.userRepo.List(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// AssignRoles assigns roles to a user (replaces existing roles)
func (uc *UserUseCase) AssignRoles(ctx context.Context, userID string, roles []domain.UserRole) error {
	if userID == "" {
		return domain.ErrInvalidUsername
	}

	if len(roles) == 0 {
		return domain.ErrNoRolesAssigned
	}

	// Validate roles
	for _, role := range roles {
		switch role {
		case domain.RoleAdmin, domain.RoleManager, domain.RoleWaiter, domain.RoleChef, domain.RoleCashier:
			// Valid role
		default:
			return domain.ErrInvalidRole
		}
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Assign roles (replace existing)
	user.Roles = roles

	// Save changes
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user roles: %w", err)
	}

	return nil
}

// GetUserRoles retrieves the roles of a user
func (uc *UserUseCase) GetUserRoles(ctx context.Context, userID string) ([]domain.UserRole, error) {
	if userID == "" {
		return nil, domain.ErrInvalidUsername
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user.Roles, nil
}

// ChangePassword changes a user's password
func (uc *UserUseCase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	if userID == "" {
		return domain.ErrInvalidUsername
	}

	// Get user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify old password
	if err := uc.passwordHasher.Compare(user.Password, oldPassword); err != nil {
		return fmt.Errorf("invalid old password")
	}

	// Hash new password
	hashedPassword, err := uc.passwordHasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	user.Password = hashedPassword

	// Save changes
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
