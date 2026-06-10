package repository

import (
	"context"
	"restaurant-management/services/auth-service/internal/domain"
)

// RefreshTokenRepository defines the interface for refresh token storage (Redis).
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *domain.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	Update(ctx context.Context, token *domain.RefreshToken) error
	Delete(ctx context.Context, token string) error
}
