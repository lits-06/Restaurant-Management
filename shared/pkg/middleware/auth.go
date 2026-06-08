package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"restaurant-management/shared/pkg/jwt"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
	RolesKey  contextKey = "roles"
)

// AuthInterceptor validates JWT tokens
func AuthInterceptor(jwtManager *jwt.Manager, publicMethods map[string]bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip authentication for public methods
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		accessToken := values[0]
		if !strings.HasPrefix(accessToken, "Bearer ") {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}

		accessToken = strings.TrimPrefix(accessToken, "Bearer ")

		// Verify token
		claims, err := jwtManager.VerifyToken(accessToken)
		if err != nil {
			if err == jwt.ErrExpiredToken {
				return nil, status.Error(codes.Unauthenticated, "token has expired")
			}
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Add claims to context
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)

		return handler(ctx, req)
	}
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetEmailFromContext extracts email from context
func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailKey).(string)
	return email, ok
}

// GetRolesFromContext extracts roles from context
func GetRolesFromContext(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(RolesKey).([]string)
	return roles, ok
}

// HasRole checks if user has a specific role
func HasRole(ctx context.Context, role string) bool {
	roles, ok := GetRolesFromContext(ctx)
	if !ok {
		return false
	}

	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole checks if user has any of the specified roles
func HasAnyRole(ctx context.Context, requiredRoles []string) bool {
	roles, ok := GetRolesFromContext(ctx)
	if !ok {
		return false
	}

	for _, required := range requiredRoles {
		for _, userRole := range roles {
			if userRole == required {
				return true
			}
		}
	}
	return false
}
