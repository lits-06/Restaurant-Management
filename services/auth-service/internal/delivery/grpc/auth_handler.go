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

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	if req.Email == "" || req.Password == "" {
		return &authpb.RegisterResponse{Success: false, Message: "email and password are required"}, nil
	}

	userID, err := h.authUseCase.Register(ctx, req.Email, req.Password, req.Username, req.FullName, req.Phone)
	if err != nil {
		switch err {
		case domain.ErrWeakPassword:
			return &authpb.RegisterResponse{Success: false, Message: "password must be at least 8 characters"}, nil
		}
		// user-service may return "user already exists" wrapped in the error message
		return &authpb.RegisterResponse{Success: false, Message: err.Error()}, nil
	}

	return &authpb.RegisterResponse{Success: true, Message: "User registered successfully", UserId: userID}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return &authpb.LoginResponse{Success: false, Message: "email and password are required"}, nil
	}

	accessToken, refreshToken, userID, err := h.authUseCase.Login(ctx, req.Email, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			return &authpb.LoginResponse{Success: false, Message: "invalid email or password"}, nil
		}
		return nil, status.Errorf(codes.Internal, "login failed: %v", err)
	}

	return &authpb.LoginResponse{
		Success:      true,
		Message:      "login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       userID,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return &authpb.RefreshTokenResponse{Success: false, Message: "refresh token is required"}, nil
	}

	accessToken, err := h.authUseCase.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		switch err {
		case domain.ErrInvalidToken:
			return &authpb.RefreshTokenResponse{Success: false, Message: "invalid refresh token"}, nil
		case domain.ErrTokenExpired:
			return &authpb.RefreshTokenResponse{Success: false, Message: "refresh token has expired"}, nil
		}
		return nil, status.Errorf(codes.Internal, "refresh failed: %v", err)
	}

	return &authpb.RefreshTokenResponse{Success: true, Message: "token refreshed", AccessToken: accessToken}, nil
}

func (h *AuthHandler) VerifyToken(ctx context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	if req.AccessToken == "" {
		return &authpb.VerifyTokenResponse{Valid: false}, nil
	}

	claims, err := h.authUseCase.VerifyToken(ctx, req.AccessToken)
	if err != nil {
		if err == jwt.ErrExpiredToken || err == jwt.ErrInvalidToken {
			return &authpb.VerifyTokenResponse{Valid: false}, nil
		}
		return nil, status.Errorf(codes.Internal, "verify failed: %v", err)
	}

	return &authpb.VerifyTokenResponse{
		Valid:     true,
		UserId:    claims.UserID,
		Email:     claims.Email,
		Roles:     claims.Roles,
		ExpiresAt: timestamppb.New(claims.ExpiresAt.Time),
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	if err := h.authUseCase.Logout(ctx, req.RefreshToken); err != nil {
		return nil, status.Errorf(codes.Internal, "logout failed: %v", err)
	}
	return &authpb.LogoutResponse{Success: true, Message: "logout successful"}, nil
}

func (h *AuthHandler) ChangePassword(ctx context.Context, req *authpb.ChangePasswordRequest) (*authpb.ChangePasswordResponse, error) {
	if req.UserId == "" || req.OldPassword == "" || req.NewPassword == "" {
		return &authpb.ChangePasswordResponse{Success: false, Message: "user_id, old_password, and new_password are required"}, nil
	}

	err := h.authUseCase.ChangePassword(ctx, req.UserId, req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case domain.ErrWeakPassword:
			return &authpb.ChangePasswordResponse{Success: false, Message: "new password must be at least 8 characters"}, nil
		}
		return &authpb.ChangePasswordResponse{Success: false, Message: err.Error()}, nil
	}

	return &authpb.ChangePasswordResponse{Success: true, Message: "password changed successfully"}, nil
}
