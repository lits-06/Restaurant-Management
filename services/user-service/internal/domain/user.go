package domain

import (
	"strings"
	"time"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleUser    UserRole = "USER"    // customer
	RoleManager UserRole = "MANAGER"
	RoleChef    UserRole = "CHEF"
	RoleWaiter  UserRole = "WAITER"
	RoleAdmin   UserRole = "ADMIN"
)

// UserStatus represents the status of a user account
type UserStatus string

const (
	StatusActive    UserStatus = "ACTIVE"
	StatusInactive  UserStatus = "INACTIVE"
	StatusSuspended UserStatus = "SUSPENDED"
)

// User represents a user in the system
type User struct {
	ID        string
	Email     string
	Username  string
	FullName  string
	Phone     string
	Password  string // Hashed password
	Roles     []UserRole
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new User with validation
func NewUser(email, username, fullName, phone, hashedPassword string, roles []UserRole) (*User, error) {
	user := &User{
		Email:     strings.TrimSpace(email),
		Username:  strings.TrimSpace(username),
		FullName:  strings.TrimSpace(fullName),
		Phone:     strings.TrimSpace(phone),
		Password:  hashedPassword,
		Roles:     roles,
		Status:    StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate validates the user fields
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}

	if !isValidEmail(u.Email) {
		return ErrInvalidEmailFormat
	}

	if u.Username == "" {
		return ErrInvalidUsername
	}

	if len(u.Username) < 3 {
		return ErrUsernameTooShort
	}

	if u.FullName == "" {
		return ErrInvalidFullName
	}

	if u.Password == "" {
		return ErrInvalidPassword
	}

	if len(u.Roles) == 0 {
		return ErrNoRolesAssigned
	}

	// Validate roles
	for _, role := range u.Roles {
		if !isValidRole(role) {
			return ErrInvalidRole
		}
	}

	// Validate status
	if !isValidStatus(u.Status) {
		return ErrInvalidStatus
	}

	return nil
}

// HasRole checks if user has a specific role
func (u *User) HasRole(role UserRole) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// IsActive checks if user account is active
func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

// Activate activates the user account
func (u *User) Activate() {
	u.Status = StatusActive
	u.UpdatedAt = time.Now()
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.Status = StatusInactive
	u.UpdatedAt = time.Now()
}

// Suspend suspends the user account
func (u *User) Suspend() {
	u.Status = StatusSuspended
	u.UpdatedAt = time.Now()
}

// AddRole adds a role to the user if not already present
func (u *User) AddRole(role UserRole) error {
	if !isValidRole(role) {
		return ErrInvalidRole
	}

	if u.HasRole(role) {
		return ErrRoleAlreadyAssigned
	}

	u.Roles = append(u.Roles, role)
	u.UpdatedAt = time.Now()
	return nil
}

// RemoveRole removes a role from the user
func (u *User) RemoveRole(role UserRole) error {
	if !u.HasRole(role) {
		return ErrRoleNotFound
	}

	// Cannot remove last role
	if len(u.Roles) == 1 {
		return ErrCannotRemoveLastRole
	}

	newRoles := make([]UserRole, 0, len(u.Roles)-1)
	for _, r := range u.Roles {
		if r != role {
			newRoles = append(newRoles, r)
		}
	}

	u.Roles = newRoles
	u.UpdatedAt = time.Now()
	return nil
}

// Update updates user information
func (u *User) Update(email, username, fullName, phone string) error {
	if email != "" {
		email = strings.TrimSpace(email)
		if !isValidEmail(email) {
			return ErrInvalidEmailFormat
		}
		u.Email = email
	}

	if username != "" {
		username = strings.TrimSpace(username)
		if len(username) < 3 {
			return ErrUsernameTooShort
		}
		u.Username = username
	}

	if fullName != "" {
		u.FullName = strings.TrimSpace(fullName)
	}

	if phone != "" {
		u.Phone = strings.TrimSpace(phone)
	}

	u.UpdatedAt = time.Now()
	return nil
}

// Helper functions

func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}
	// Simple email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func isValidRole(role UserRole) bool {
	switch role {
	case RoleUser, RoleManager, RoleChef, RoleWaiter, RoleAdmin:
		return true
	default:
		return false
	}
}

func isValidStatus(status UserStatus) bool {
	switch status {
	case StatusActive, StatusInactive, StatusSuspended:
		return true
	default:
		return false
	}
}
