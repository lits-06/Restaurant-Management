package grpcclient

import (
	"context"
	"fmt"
	"time"

	userpb "restaurant-management/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultTimeout = 5 * time.Second

type UserClient struct {
	conn   *grpc.ClientConn
	client userpb.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user-service at %s: %w", addr, err)
	}
	return &UserClient{conn: conn, client: userpb.NewUserServiceClient(conn)}, nil
}

func (c *UserClient) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// CreateUser creates a new user in user-service. Returns the created user_id.
func (c *UserClient) CreateUser(ctx context.Context, email, username, fullName, phone, password string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	resp, err := c.client.CreateUser(ctx, &userpb.CreateUserRequest{
		Email:    email,
		Username: username,
		FullName: fullName,
		Phone:    phone,
		Password: password,
		Roles:    []userpb.UserRole{userpb.UserRole_ROLE_USER},
	})
	if err != nil {
		return "", err
	}
	if !resp.Success || resp.User == nil {
		return "", fmt.Errorf("create user failed: %s", resp.Message)
	}
	return resp.User.UserId, nil
}

// VerifyCredentials checks email + password. Returns user_id, email, roles string slice on success.
func (c *UserClient) VerifyCredentials(ctx context.Context, email, password string) (userID, userEmail string, roles []string, err error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	resp, err := c.client.VerifyCredentials(ctx, &userpb.VerifyCredentialsRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", "", nil, err
	}
	if !resp.Success {
		return "", "", nil, fmt.Errorf("%s", resp.Message)
	}
	return resp.UserId, resp.Email, protoRolesToStrings(resp.Roles), nil
}

// GetUser fetches a user by ID. Returns email and roles.
func (c *UserClient) GetUser(ctx context.Context, userID string) (email string, roles []string, err error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	resp, err := c.client.GetUser(ctx, &userpb.GetUserRequest{UserId: userID})
	if err != nil {
		return "", nil, err
	}
	if !resp.Success || resp.User == nil {
		return "", nil, fmt.Errorf("user not found")
	}
	return resp.User.Email, protoRolesToStrings(resp.User.Roles), nil
}

// ChangePassword delegates password change to user-service.
func (c *UserClient) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	resp, err := c.client.ChangePassword(ctx, &userpb.ChangePasswordRequest{
		UserId:      userID,
		OldPassword: oldPassword,
		NewPassword: newPassword,
	})
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("%s", resp.Message)
	}
	return nil
}

func protoRolesToStrings(roles []userpb.UserRole) []string {
	out := make([]string, 0, len(roles))
	for _, r := range roles {
		switch r {
		case userpb.UserRole_ROLE_USER:
			out = append(out, "USER")
		case userpb.UserRole_ROLE_MANAGER:
			out = append(out, "MANAGER")
		case userpb.UserRole_ROLE_CHEF:
			out = append(out, "CHEF")
		case userpb.UserRole_ROLE_WAITER:
			out = append(out, "WAITER")
		case userpb.UserRole_ROLE_ADMIN:
			out = append(out, "ADMIN")
		}
	}
	return out
}
