package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Manager handles JWT operations
type Manager struct {
	secretKey            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	issuer               string
}

// NewManager creates a new JWT manager
func NewManager(secretKey string, accessTokenMinutes, refreshTokenHours int, issuer string) *Manager {
	return &Manager{
		secretKey:            secretKey,
		accessTokenDuration:  time.Duration(accessTokenMinutes) * time.Minute,
		refreshTokenDuration: time.Duration(refreshTokenHours) * time.Hour,
		issuer:               issuer,
	}
}

// GenerateAccessToken generates a new access token
func (m *Manager) GenerateAccessToken(userID, email string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// GenerateRefreshToken generates a new refresh token
func (m *Manager) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshTokenDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// VerifyToken verifies and parses a token
func (m *Manager) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetAccessTokenDuration returns access token duration
func (m *Manager) GetAccessTokenDuration() time.Duration {
	return m.accessTokenDuration
}

// GetRefreshTokenDuration returns refresh token duration
func (m *Manager) GetRefreshTokenDuration() time.Duration {
	return m.refreshTokenDuration
}
