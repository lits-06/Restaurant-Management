package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	notifpb "restaurant-management/proto/notification"
)

type NotificationClient struct {
	conn   *grpc.ClientConn
	client notifpb.NotificationServiceClient
}

func NewNotificationClient(host string, port int) (*NotificationClient, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to notification service: %w", err)
	}
	return &NotificationClient{conn: conn, client: notifpb.NewNotificationServiceClient(conn)}, nil
}

func (c *NotificationClient) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

// Subscribe opens a server-side streaming RPC for the given role.
// The stream remains open until ctx is cancelled or the server closes it.
func (c *NotificationClient) Subscribe(ctx context.Context, role string) (notifpb.NotificationService_SubscribeClient, error) {
	return c.client.Subscribe(ctx, &notifpb.SubscribeRequest{Role: role})
}
