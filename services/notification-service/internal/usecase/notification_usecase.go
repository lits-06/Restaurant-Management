package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"restaurant-management/services/notification-service/internal/domain"
	"restaurant-management/services/notification-service/internal/repository"
)

type NotificationUseCase struct {
	repo *repository.PubSubRepository
}

func NewNotificationUseCase(repo *repository.PubSubRepository) *NotificationUseCase {
	return &NotificationUseCase{repo: repo}
}

func (uc *NotificationUseCase) Send(ctx context.Context, n *domain.Notification) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	if n.CreatedAt == 0 {
		n.CreatedAt = time.Now().Unix()
	}
	if err := uc.repo.Publish(ctx, n); err != nil {
		return fmt.Errorf("failed to publish notification: %w", err)
	}
	return nil
}

func (uc *NotificationUseCase) Subscribe(ctx context.Context, role string) (<-chan *domain.Notification, error) {
	return uc.repo.Subscribe(ctx, role)
}
