package usecase

import (
	"context"
	"time"

	"restaurant-management/services/auth-service/internal/domain"
	"restaurant-management/services/auth-service/internal/repository"
	"restaurant-management/shared/pkg/jwt"
	"restaurant-management/shared/pkg/utils"
)

// AuthUseCase handles authentication business logic
type AuthUseCase struct {
	userRepo       repository.UserRepository
	tokenRepo      repository.RefreshTokenRepository
	jwtManager     *jwt.Manager
	passwordHasher PasswordHasher
}

// PasswordHasher interface for password hashing
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) error
}

// NewAuthUseCase creates a new AuthUseCase
func NewAuthUseCase(
	userRepo repository.UserRepository,
	tokenRepo repository.RefreshTokenRepository,
	jwtManager *jwt.Manager,
	passwordHasher PasswordHasher,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		jwtManager:     jwtManager,
		passwordHasher: passwordHasher,
	}
}

// Register registers a new user
func (uc *AuthUseCase) Register(ctx context.Context, email, password, username, fullName, phone string) (string, error) {
	// Check if user already exists
	exists, err := uc.userRepo.Exists(ctx, email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", domain.ErrUserAlreadyExists
	}

	// Validate password strength
	if len(password) < 8 {
		return "", domain.ErrWeakPassword
	}

	// Hash password
	passwordHash, err := uc.passwordHasher.Hash(password)
	if err != nil {
		return "", err
	}

	// Create user
	userID := utils.GenerateUUID()
	user := domain.NewUser(userID, email, passwordHash)

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return "", err
	}

	return userID, nil
}

// Login authenticates a user and returns tokens
func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*domain.AuthSession, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if err := uc.passwordHasher.Compare(user.Password, password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate access token
	accessToken, err := uc.jwtManager.GenerateAccessToken(user.UserID, user.Email)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenStr, err := uc.jwtManager.GenerateRefreshToken(user.UserID)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	expiresAt := time.Now().Add(uc.jwtManager.GetRefreshTokenDuration())
	refreshToken := domain.NewRefreshToken(user.UserID, refreshTokenStr, expiresAt)

	if err := uc.tokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, err
	}

	// Create session
	session := domain.NewAuthSession(
		user.UserID,
		user.Email,
		accessToken,
		refreshTokenStr,
		int64(uc.jwtManager.GetAccessTokenDuration().Seconds()),
	)

	return session, nil
}

// RefreshToken generates new tokens using refresh token
func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshTokenStr string) (*domain.AuthSession, error) {
	// Verify refresh token format
	claims, err := uc.jwtManager.VerifyToken(refreshTokenStr)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Find refresh token in database
	refreshToken, err := uc.tokenRepo.FindByToken(ctx, refreshTokenStr)
	if err != nil {
		if err == domain.ErrTokenNotFound {
			return nil, domain.ErrInvalidToken
		}
		return nil, err
	}

	// Check if token is valid
	if !refreshToken.IsValid() {
		return nil, domain.ErrTokenExpired
	}
	if claims.UserID != "" && claims.UserID != refreshToken.UserID {
		return nil, domain.ErrInvalidToken
	}

	// Get user
	user, err := uc.userRepo.FindByID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new access token
	newAccessToken, err := uc.jwtManager.GenerateAccessToken(user.UserID, user.Email)
	if err != nil {
		return nil, err
	}

	// Create session
	session := domain.NewAuthSession(
		user.UserID,
		user.Email,
		newAccessToken,
		refreshTokenStr,
		int64(uc.jwtManager.GetAccessTokenDuration().Seconds()),
	)

	return session, nil
}

// VerifyToken verifies an access token
func (uc *AuthUseCase) VerifyToken(ctx context.Context, accessToken string) (*jwt.Claims, error) {
	claims, err := uc.jwtManager.VerifyToken(accessToken)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// TODO
// Logout revokes all refresh tokens for a user
func (uc *AuthUseCase) Logout(ctx context.Context, userID string) error {
	return nil
}
