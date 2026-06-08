package grpcclient

import (
	"context"
	"fmt"

	orderpb "restaurant-management/proto/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// OrderClient wraps gRPC connection and generated order client.
type OrderClient struct {
	conn   *grpc.ClientConn
	client orderpb.OrderServiceClient
}

// NewOrderClient creates a new OrderService gRPC client.
func NewOrderClient(host string, port int) (*OrderClient, error) {
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}

	return &OrderClient{
		conn:   conn,
		client: orderpb.NewOrderServiceClient(conn),
	}, nil
}

func (c *OrderClient) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *OrderClient) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.CreateOrder(ctx, req)
}

func (c *OrderClient) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.GetOrder(ctx, req)
}

func (c *OrderClient) UpdateOrder(ctx context.Context, req *orderpb.UpdateOrderRequest) (*orderpb.UpdateOrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.UpdateOrder(ctx, req)
}

func (c *OrderClient) DeleteOrder(ctx context.Context, req *orderpb.DeleteOrderRequest) (*orderpb.DeleteOrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.DeleteOrder(ctx, req)
}

func (c *OrderClient) CancelOrder(ctx context.Context, req *orderpb.CancelOrderRequest) (*orderpb.CancelOrderResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.CancelOrder(ctx, req)
}

func (c *OrderClient) ListOrders(ctx context.Context, req *orderpb.ListOrdersRequest) (*orderpb.ListOrdersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.ListOrders(ctx, req)
}

func (c *OrderClient) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*orderpb.UpdateOrderStatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.UpdateOrderStatus(ctx, req)
}

func (c *OrderClient) AddOrderItem(ctx context.Context, req *orderpb.AddOrderItemRequest) (*orderpb.AddOrderItemResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.AddOrderItem(ctx, req)
}

func (c *OrderClient) RemoveOrderItem(ctx context.Context, req *orderpb.RemoveOrderItemRequest) (*orderpb.RemoveOrderItemResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
	defer cancel()
	return c.client.RemoveOrderItem(ctx, req)
}

// func (c *OrderClient) GetOrdersByTable(ctx context.Context, req *orderpb.GetOrdersByTableRequest) (*orderpb.GetOrdersByTableResponse, error) {
// 	ctx, cancel := context.WithTimeout(ctx, defaultRPCTimeout)
// 	defer cancel()
// 	return c.client.GetOrdersByTable(ctx, req)
// }
