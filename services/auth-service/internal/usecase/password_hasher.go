package usecase

import (
	"restaurant-management/shared/pkg/utils"
)

// BcryptPasswordHasher implements PasswordHasher using bcrypt
type BcryptPasswordHasher struct{}

// NewBcryptPasswordHasher creates a new bcrypt password hasher
func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{}
}

// Hash hashes a password using bcrypt
func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	return utils.HashPassword(password)
}

// Compare compares a hashed password with a plain password
func (h *BcryptPasswordHasher) Compare(hashedPassword, password string) error {
	return utils.ComparePassword(hashedPassword, password)
}
