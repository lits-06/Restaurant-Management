package grpcclient

import (
	"context"
	"fmt"

	staffpb "restaurant-management/proto/staff"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// StaffClient wraps gRPC connection and generated staff client.
type StaffClient struct {
	conn   *grpc.ClientConn
	client staffpb.StaffServiceClient
}

// NewStaffClient creates a new StaffService gRPC client.
func NewStaffClient(host string, port int) (*StaffClient, error) {
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to staff service: %w", err)
	}

	return &StaffClient{
		conn:   conn,
		client: staffpb.NewStaffServiceClient(conn),
	}, nil
}

func (c *StaffClient) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *StaffClient) CreateStaff(ctx context.Context, req *staffpb.CreateStaffRequest) (*staffpb.CreateStaffResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.CreateStaff(ctx, req)
}

func (c *StaffClient) GetStaff(ctx context.Context, req *staffpb.GetStaffRequest) (*staffpb.GetStaffResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.GetStaff(ctx, req)
}

func (c *StaffClient) UpdateStaff(ctx context.Context, req *staffpb.UpdateStaffRequest) (*staffpb.UpdateStaffResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.UpdateStaff(ctx, req)
}

func (c *StaffClient) DeleteStaff(ctx context.Context, req *staffpb.DeleteStaffRequest) (*staffpb.DeleteStaffResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.DeleteStaff(ctx, req)
}

func (c *StaffClient) ListStaff(ctx context.Context, req *staffpb.ListStaffRequest) (*staffpb.ListStaffResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.ListStaff(ctx, req)
}
