package domain

import (
	"time"
)

// User represents user credentials for authentication
type User struct {
	UserID    string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new User entity
func NewUser(userID, email, password string) *User {
	now := time.Now()
	return &User{
		UserID:    userID,
		Email:     email,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdatePassword updates the user's password
func (u *User) UpdatePassword(newPassword string) {
	u.Password = newPassword
	u.UpdatedAt = time.Now()
}

// Validate validates the user entity
func (u *User) Validate() error {
	if u.UserID == "" {
		return ErrInvalidUserID
	}
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if u.Password == "" {
		return ErrInvalidPassword
	}
	return nil
}
