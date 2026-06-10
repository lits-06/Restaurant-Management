package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"restaurant-management/services/notification-service/internal/domain"
)

const channelPrefix = "notifications:"

type PubSubRepository struct {
	rdb *redis.Client
}

func NewPubSubRepository(rdb *redis.Client) *PubSubRepository {
	return &PubSubRepository{rdb: rdb}
}

func (r *PubSubRepository) Publish(ctx context.Context, n *domain.Notification) error {
	data, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}
	channel := channelPrefix + n.TargetRole
	return r.rdb.Publish(ctx, channel, string(data)).Err()
}

// Subscribe returns a channel that emits notifications for the given role.
// The caller must cancel ctx to stop the subscription.
func (r *PubSubRepository) Subscribe(ctx context.Context, role string) (<-chan *domain.Notification, error) {
	channel := channelPrefix + role
	sub := r.rdb.Subscribe(ctx, channel)

	// Verify subscription is up.
	if _, err := sub.Receive(ctx); err != nil {
		sub.Close()
		return nil, fmt.Errorf("failed to subscribe to %s: %w", channel, err)
	}

	out := make(chan *domain.Notification, 32)
	go func() {
		defer sub.Close()
		defer close(out)
		msgCh := sub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgCh:
				if !ok {
					return
				}
				var n domain.Notification
				if err := json.Unmarshal([]byte(msg.Payload), &n); err != nil {
					continue
				}
				select {
				case out <- &n:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out, nil
}
