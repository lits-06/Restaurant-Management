package grpcclient

import (
	"context"
	"fmt"

	tablepb "restaurant-management/proto/table"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TableClient wraps the generated table-service gRPC client.
type TableClient struct {
	conn   *grpc.ClientConn
	client tablepb.TableServiceClient
}

// NewTableClient creates a new TableService gRPC client.
func NewTableClient(host string, port int) (*TableClient, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to table service: %w", err)
	}
	return &TableClient{conn: conn, client: tablepb.NewTableServiceClient(conn)}, nil
}

func (c *TableClient) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *TableClient) CreateTable(ctx context.Context, req *tablepb.CreateTableRequest) (*tablepb.CreateTableResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.CreateTable(ctx, req)
}

func (c *TableClient) GetTable(ctx context.Context, req *tablepb.GetTableRequest) (*tablepb.GetTableResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.GetTable(ctx, req)
}

func (c *TableClient) UpdateTable(ctx context.Context, req *tablepb.UpdateTableRequest) (*tablepb.UpdateTableResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.UpdateTable(ctx, req)
}

func (c *TableClient) DeleteTable(ctx context.Context, req *tablepb.DeleteTableRequest) (*tablepb.DeleteTableResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.DeleteTable(ctx, req)
}

func (c *TableClient) ListTables(ctx context.Context, req *tablepb.ListTablesRequest) (*tablepb.ListTablesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.ListTables(ctx, req)
}

func (c *TableClient) UpdateTableStatus(ctx context.Context, req *tablepb.UpdateTableStatusRequest) (*tablepb.UpdateTableStatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.UpdateTableStatus(ctx, req)
}

func (c *TableClient) GetAvailableTables(ctx context.Context, req *tablepb.GetAvailableTablesRequest) (*tablepb.GetAvailableTablesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.GetAvailableTables(ctx, req)
}
