package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	authpb "restaurant-management/proto/auth"
	"restaurant-management/services/auth-service/internal/domain"
	"restaurant-management/services/auth-service/internal/usecase"
	"restaurant-management/shared/pkg/jwt"
)

// AuthHandler handles gRPC requests for authentication
type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	authUseCase *usecase.AuthUseCase
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	// Validate request
	if req.Email == "" || req.Password == "" {
		return &authpb.RegisterResponse{
			Success: false,
			Message: "Email and password are required",
		}, nil
	}

	// Register user
	userID, err := h.authUseCase.Register(ctx, req.Email, req.Password, req.Username, req.FullName, req.Phone)
	if err != nil {
		if err == domain.ErrUserAlreadyExists {
			return &authpb.RegisterResponse{
				Success: false,
				Message: "User with this email already exists",
			}, nil
		}
		if err == domain.ErrWeakPassword {
			return &authpb.RegisterResponse{
				Success: false,
				Message: "Password is too weak. Minimum 6 characters required",
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &authpb.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		UserId:  userID,
	}, nil
}

// Login handles user login
func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	// Validate request
	if req.Email == "" || req.Password == "" {
		return &authpb.LoginResponse{
			Success: false,
			Message: "Email and password are required",
		}, nil
	}

	// Login user
	session, err := h.authUseCase.Login(ctx, req.Email, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			return &authpb.LoginResponse{
				Success: false,
				Message: "Invalid email or password",
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to login: %v", err)
	}

	return &authpb.LoginResponse{
		Success:      true,
		Message:      "Login successful",
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		UserId:       session.UserID,
	}, nil
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	// Validate request
	if req.RefreshToken == "" {
		return &authpb.RefreshTokenResponse{
			Success: false,
			Message: "Refresh token is required",
		}, nil
	}

	// Refresh token
	session, err := h.authUseCase.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if err == domain.ErrInvalidToken {
			return &authpb.RefreshTokenResponse{
				Success: false,
				Message: "Invalid refresh token",
			}, nil
		}
		if err == domain.ErrTokenExpired {
			return &authpb.RefreshTokenResponse{
				Success: false,
				Message: "Refresh token has expired",
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to refresh token: %v", err)
	}

	return &authpb.RefreshTokenResponse{
		Success:     true,
		Message:     "Token refreshed successfully",
		AccessToken: session.AccessToken,
	}, nil
}

// VerifyToken handles token verification
func (h *AuthHandler) VerifyToken(ctx context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	// Validate request
	if req.AccessToken == "" {
		return &authpb.VerifyTokenResponse{
			Valid: false,
		}, nil
	}

	// Verify token
	claims, err := h.authUseCase.VerifyToken(ctx, req.AccessToken)
	if err != nil {
		if err == jwt.ErrExpiredToken || err == jwt.ErrInvalidToken {
			return &authpb.VerifyTokenResponse{
				Valid: false,
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to verify token: %v", err)
	}

	return &authpb.VerifyTokenResponse{
		Valid:     true,
		UserId:    claims.UserID,
		Email:     claims.Email,
		ExpiresAt: timestamppb.New(claims.ExpiresAt.Time),
	}, nil
}

// Logout handles user logout
func (h *AuthHandler) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	userID := req.UserId
	if userID == "" && req.AccessToken != "" {
		claims, err := h.authUseCase.VerifyToken(ctx, req.AccessToken)
		if err != nil {
			if err == jwt.ErrExpiredToken || err == jwt.ErrInvalidToken {
				return &authpb.LogoutResponse{
					Success: false,
					Message: "Invalid access token",
				}, nil
			}
			return nil, status.Errorf(codes.Internal, "failed to verify token: %v", err)
		}
		userID = claims.UserID
	}

	if userID == "" {
		return &authpb.LogoutResponse{
			Success: false,
			Message: "User ID or access token is required",
		}, nil
	}

	// Logout user (revoke refresh tokens)
	if err := h.authUseCase.Logout(ctx, userID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to logout: %v", err)
	}

	return &authpb.LogoutResponse{
		Success: true,
		Message: "Logout successful",
	}, nil
}
