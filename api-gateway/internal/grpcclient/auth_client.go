package grpcclient

import (
	"context"
	"fmt"
	"time"

	authpb "restaurant-management/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultRPCTimeout = 10 * time.Second

// AuthClient wraps gRPC connection and generated auth client.
type AuthClient struct {
	conn   *grpc.ClientConn
	client authpb.AuthServiceClient
}

// NewAuthClient creates a new AuthService gRPC client.
func NewAuthClient(host string, port int) (*AuthClient, error) {
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &AuthClient{
		conn:   conn,
		client: authpb.NewAuthServiceClient(conn),
	}, nil
}

func (c *AuthClient) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *AuthClient) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.Register(ctx, req)
}

func (c *AuthClient) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.Login(ctx, req)
}

func (c *AuthClient) RefreshToken(ctx context.Context, req *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.RefreshToken(ctx, req)
}

func (c *AuthClient) VerifyToken(ctx context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.VerifyToken(ctx, req)
}

func (c *AuthClient) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.Logout(ctx, req)
}

func (c *AuthClient) ChangePassword(ctx context.Context, req *authpb.ChangePasswordRequest) (*authpb.ChangePasswordResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.ChangePassword(ctx, req)
}
