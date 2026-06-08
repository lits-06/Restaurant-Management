package domain

import (
	"time"
)

// RefreshToken represents a refresh token for JWT authentication
type RefreshToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
	CreatedAt time.Time
	IsRevoked bool
}

// NewRefreshToken creates a new RefreshToken entity
func NewRefreshToken(userID, token string, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		IsRevoked: false,
	}
}

// IsValid checks if the refresh token is still valid
func (rt *RefreshToken) IsValid() bool {
	if rt.IsRevoked {
		return false
	}
	if time.Now().After(rt.ExpiresAt) {
		return false
	}
	return true
}

// Revoke marks the token as revoked
func (rt *RefreshToken) Revoke() {
	rt.IsRevoked = true
}

// Validate validates the refresh token entity
func (rt *RefreshToken) Validate() error {
	if rt.UserID == "" {
		return ErrInvalidUserID
	}
	if rt.Token == "" {
		return ErrInvalidToken
	}
	if rt.ExpiresAt.IsZero() {
		return ErrInvalidExpiryTime
	}
	return nil
}
