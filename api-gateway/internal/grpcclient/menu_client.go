package grpcclient

import (
	"context"
	"fmt"

	menupb "restaurant-management/proto/menu"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// MenuClient wraps gRPC connection and generated menu client.
type MenuClient struct {
	conn   *grpc.ClientConn
	client menupb.MenuServiceClient
}

// NewMenuClient creates a new MenuService gRPC client.
func NewMenuClient(host string, port int) (*MenuClient, error) {
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to menu service: %w", err)
	}

	return &MenuClient{
		conn:   conn,
		client: menupb.NewMenuServiceClient(conn),
	}, nil
}

func (c *MenuClient) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *MenuClient) CreateMenuItem(ctx context.Context, req *menupb.CreateMenuItemRequest) (*menupb.CreateMenuItemResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.CreateMenuItem(ctx, req)
}

func (c *MenuClient) GetMenuItem(ctx context.Context, req *menupb.GetMenuItemRequest) (*menupb.GetMenuItemResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.GetMenuItem(ctx, req)
}

func (c *MenuClient) UpdateMenuItem(ctx context.Context, req *menupb.UpdateMenuItemRequest) (*menupb.UpdateMenuItemResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.UpdateMenuItem(ctx, req)
}

func (c *MenuClient) DeleteMenuItem(ctx context.Context, req *menupb.DeleteMenuItemRequest) (*menupb.DeleteMenuItemResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.DeleteMenuItem(ctx, req)
}

func (c *MenuClient) ListMenuItems(ctx context.Context, req *menupb.ListMenuItemsRequest) (*menupb.ListMenuItemsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.ListMenuItems(ctx, req)
}

func (c *MenuClient) CreateCategory(ctx context.Context, req *menupb.CreateCategoryRequest) (*menupb.CreateCategoryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.CreateCategory(ctx, req)
}

func (c *MenuClient) GetCategory(ctx context.Context, req *menupb.GetCategoryRequest) (*menupb.GetCategoryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.GetCategory(ctx, req)
}

func (c *MenuClient) UpdateCategory(ctx context.Context, req *menupb.UpdateCategoryRequest) (*menupb.UpdateCategoryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.UpdateCategory(ctx, req)
}

func (c *MenuClient) DeleteCategory(ctx context.Context, req *menupb.DeleteCategoryRequest) (*menupb.DeleteCategoryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.DeleteCategory(ctx, req)
}

func (c *MenuClient) ListCategories(ctx context.Context, req *menupb.ListCategoriesRequest) (*menupb.ListCategoriesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.ListCategories(ctx, req)
}
