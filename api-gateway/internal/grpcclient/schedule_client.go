package grpcclient

import (
	"context"
	"fmt"

	schedulepb "restaurant-management/proto/schedule"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ScheduleClient struct {
	conn   *grpc.ClientConn
	client schedulepb.ScheduleServiceClient
}

func NewScheduleClient(host string, port int) (*ScheduleClient, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to schedule service: %w", err)
	}
	return &ScheduleClient{conn: conn, client: schedulepb.NewScheduleServiceClient(conn)}, nil
}

func (c *ScheduleClient) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *ScheduleClient) CreateShift(ctx context.Context, req *schedulepb.CreateShiftRequest) (*schedulepb.CreateShiftResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.CreateShift(ctx, req)
}

func (c *ScheduleClient) GetShift(ctx context.Context, req *schedulepb.GetShiftRequest) (*schedulepb.GetShiftResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.GetShift(ctx, req)
}

func (c *ScheduleClient) UpdateShift(ctx context.Context, req *schedulepb.UpdateShiftRequest) (*schedulepb.UpdateShiftResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.UpdateShift(ctx, req)
}

func (c *ScheduleClient) DeleteShift(ctx context.Context, req *schedulepb.DeleteShiftRequest) (*schedulepb.DeleteShiftResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.DeleteShift(ctx, req)
}

func (c *ScheduleClient) ListShifts(ctx context.Context, req *schedulepb.ListShiftsRequest) (*schedulepb.ListShiftsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.ListShifts(ctx, req)
}
