package grpc

import (
	"context"
	"errors"

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
		return &authpb.RegisterResponse{Success: false, Message: "Email and password are required"}, nil
	}

	userID, err := h.authUseCase.Register(ctx, req.Email, req.Password, req.Username, req.FullName, req.Phone)
	if err != nil {
		switch err {
		case domain.ErrWeakPassword:
			return &authpb.RegisterResponse{Success: false, Message: "Password must be at least 8 characters"}, nil
		}
		return &authpb.RegisterResponse{Success: false, Message: grpcDesc(err)}, nil
	}

	return &authpb.RegisterResponse{Success: true, Message: "User registered successfully", UserId: userID}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return &authpb.LoginResponse{Success: false, Message: "Email and password are required"}, nil
	}

	accessToken, refreshToken, userID, err := h.authUseCase.Login(ctx, req.Email, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			return &authpb.LoginResponse{Success: false, Message: "Invalid email or password"}, nil
		}
		return nil, status.Errorf(codes.Internal, "Login failed: %v", err)
	}

	return &authpb.LoginResponse{
		Success:      true,
		Message:      "Login successful",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserId:       userID,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return &authpb.RefreshTokenResponse{Success: false, Message: "Refresh token is required"}, nil
	}

	accessToken, err := h.authUseCase.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		switch err {
		case domain.ErrInvalidToken:
			return &authpb.RefreshTokenResponse{Success: false, Message: "Invalid refresh token"}, nil
		case domain.ErrTokenExpired:
			return &authpb.RefreshTokenResponse{Success: false, Message: "Refresh token has expired"}, nil
		}
		return nil, status.Errorf(codes.Internal, "Refresh failed: %v", err)
	}

	return &authpb.RefreshTokenResponse{Success: true, Message: "Token refreshed", AccessToken: accessToken}, nil
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
		return nil, status.Errorf(codes.Internal, "Verify failed: %v", err)
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
		return nil, status.Errorf(codes.Internal, "Logout failed: %v", err)
	}
	return &authpb.LogoutResponse{Success: true, Message: "Logout successful"}, nil
}

func (h *AuthHandler) ChangePassword(ctx context.Context, req *authpb.ChangePasswordRequest) (*authpb.ChangePasswordResponse, error) {
	if req.UserId == "" || req.OldPassword == "" || req.NewPassword == "" {
		return &authpb.ChangePasswordResponse{Success: false, Message: "user_id, old_password, and new_password are required"}, nil
	}

	err := h.authUseCase.ChangePassword(ctx, req.UserId, req.OldPassword, req.NewPassword)
	if err != nil {
		switch err {
		case domain.ErrWeakPassword:
			return &authpb.ChangePasswordResponse{Success: false, Message: "New password must be at least 8 characters"}, nil
		}
		return &authpb.ChangePasswordResponse{Success: false, Message: grpcDesc(err)}, nil
	}

	return &authpb.ChangePasswordResponse{Success: true, Message: "Password changed successfully"}, nil
}

// grpcDesc extracts just the description from a gRPC status error,
// walking the error chain in case the gRPC error was wrapped with fmt.Errorf("%w").
func grpcDesc(err error) string {
	for e := err; e != nil; e = errors.Unwrap(e) {
		if s, ok := status.FromError(e); ok {
			return s.Message()
		}
	}
	return err.Error()
}
