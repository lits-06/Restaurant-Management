package usecase

import (
	"context"
	"fmt"
	"time"

	"restaurant-management/services/auth-service/internal/domain"
	"restaurant-management/services/auth-service/internal/repository"
	"restaurant-management/shared/pkg/jwt"
)

// UserServiceClient is the interface auth-service uses to talk to user-service.
type UserServiceClient interface {
	CreateUser(ctx context.Context, email, username, fullName, phone, password string) (userID string, err error)
	VerifyCredentials(ctx context.Context, email, password string) (userID string, userEmail string, roles []string, err error)
	GetUser(ctx context.Context, userID string) (userEmail string, roles []string, err error)
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
}

type AuthUseCase struct {
	tokenRepo  repository.RefreshTokenRepository
	jwtManager *jwt.Manager
	userClient UserServiceClient
}

func NewAuthUseCase(
	tokenRepo repository.RefreshTokenRepository,
	jwtManager *jwt.Manager,
	userClient UserServiceClient,
) *AuthUseCase {
	return &AuthUseCase{
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
		userClient: userClient,
	}
}

// Register creates a new user via user-service. Returns the new user_id.
func (uc *AuthUseCase) Register(ctx context.Context, email, password, username, fullName, phone string) (string, error) {
	if len(password) < 8 {
		return "", domain.ErrWeakPassword
	}
	userID, err := uc.userClient.CreateUser(ctx, email, username, fullName, phone, password)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// Login verifies credentials via user-service, then issues JWT tokens.
func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (accessToken, refreshToken, userID string, err error) {
	uid, userEmail, roles, err := uc.userClient.VerifyCredentials(ctx, email, password)
	if err != nil {
		return "", "", "", domain.ErrInvalidCredentials
	}

	accessToken, err = uc.jwtManager.GenerateAccessToken(uid, userEmail, roles)
	if err != nil {
		return "", "", "", fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err = uc.jwtManager.GenerateRefreshToken(uid)
	if err != nil {
		return "", "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	expiresAt := time.Now().Add(uc.jwtManager.GetRefreshTokenDuration())
	rt := domain.NewRefreshToken(uid, refreshToken, expiresAt)
	if err = uc.tokenRepo.Create(ctx, rt); err != nil {
		return "", "", "", fmt.Errorf("store refresh token: %w", err)
	}

	return accessToken, refreshToken, uid, nil
}

// RefreshToken issues a new access token using an existing refresh token.
func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshTokenStr string) (string, error) {
	claims, err := uc.jwtManager.VerifyToken(refreshTokenStr)
	if err != nil {
		return "", domain.ErrInvalidToken
	}

	rt, err := uc.tokenRepo.FindByToken(ctx, refreshTokenStr)
	if err != nil {
		return "", domain.ErrInvalidToken
	}
	if !rt.IsValid() {
		return "", domain.ErrTokenExpired
	}
	if claims.UserID != "" && claims.UserID != rt.UserID {
		return "", domain.ErrInvalidToken
	}

	email, roles, err := uc.userClient.GetUser(ctx, rt.UserID)
	if err != nil {
		return "", fmt.Errorf("get user for refresh: %w", err)
	}

	newAccessToken, err := uc.jwtManager.GenerateAccessToken(rt.UserID, email, roles)
	if err != nil {
		return "", fmt.Errorf("generate access token: %w", err)
	}
	return newAccessToken, nil
}

// VerifyToken parses and validates an access token.
func (uc *AuthUseCase) VerifyToken(ctx context.Context, accessToken string) (*jwt.Claims, error) {
	return uc.jwtManager.VerifyToken(accessToken)
}

// Logout revokes the given refresh token.
func (uc *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	if err := uc.tokenRepo.Delete(ctx, refreshToken); err != nil && err != domain.ErrTokenNotFound {
		return err
	}
	return nil
}

// ChangePassword delegates to user-service.
func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	if len(newPassword) < 8 {
		return domain.ErrWeakPassword
	}
	return uc.userClient.ChangePassword(ctx, userID, oldPassword, newPassword)
}
