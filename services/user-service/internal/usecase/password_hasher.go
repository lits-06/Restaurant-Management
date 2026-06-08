package usecase

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptPasswordHasher implements PasswordHasher using bcrypt
type BcryptPasswordHasher struct {
	cost int
}

// NewBcryptPasswordHasher creates a new BcryptPasswordHasher
func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{
		cost: bcrypt.DefaultCost,
	}
}

// Hash hashes a password using bcrypt
func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Compare compares a hashed password with a plaintext password
func (h *BcryptPasswordHasher) Compare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
