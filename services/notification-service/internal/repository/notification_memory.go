package repository

import (
	"context"
	"sort"
	"sync"

	"restaurant-management/services/notification-service/internal/domain"

	"github.com/google/uuid"
)

// InMemoryNotificationRepository is an in-memory implementation of NotificationRepository
type InMemoryNotificationRepository struct {
	mu              sync.RWMutex
	notifications   map[string]*domain.Notification
	recipientIndex  map[string][]string // recipientID -> []notificationID
}

// NewInMemoryNotificationRepository creates a new in-memory notification repository
func NewInMemoryNotificationRepository() *InMemoryNotificationRepository {
	return &InMemoryNotificationRepository{
		notifications:  make(map[string]*domain.Notification),
		recipientIndex: make(map[string][]string),
	}
}

// Create creates a new notification
func (r *InMemoryNotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Generate ID if not set
	if notification.NotificationID == "" {
		notification.NotificationID = uuid.New().String()
	}

	// Store notification
	r.notifications[notification.NotificationID] = notification

	// Update recipient index
	r.recipientIndex[notification.RecipientID] = append(r.recipientIndex[notification.RecipientID], notification.NotificationID)

	return nil
}

// GetByID retrieves a notification by ID
func (r *InMemoryNotificationRepository) GetByID(ctx context.Context, notificationID string) (*domain.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	notification, exists := r.notifications[notificationID]
	if !exists {
		return nil, domain.ErrNotificationNotFound
	}

	return notification, nil
}

// Update updates a notification
func (r *InMemoryNotificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.notifications[notification.NotificationID]
	if !exists {
		return domain.ErrNotificationNotFound
	}

	// Update recipient index if recipient changed (unlikely but handle it)
	if existing.RecipientID != notification.RecipientID {
		// Remove from old recipient
		oldNotifs := r.recipientIndex[existing.RecipientID]
		for i, id := range oldNotifs {
			if id == notification.NotificationID {
				r.recipientIndex[existing.RecipientID] = append(oldNotifs[:i], oldNotifs[i+1:]...)
				break
			}
		}
		// Add to new recipient
		r.recipientIndex[notification.RecipientID] = append(r.recipientIndex[notification.RecipientID], notification.NotificationID)
	}

	r.notifications[notification.NotificationID] = notification
	return nil
}

// Delete deletes a notification
func (r *InMemoryNotificationRepository) Delete(ctx context.Context, notificationID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	notification, exists := r.notifications[notificationID]
	if !exists {
		return domain.ErrNotificationNotFound
	}

	// Remove from recipient index
	recipientNotifs := r.recipientIndex[notification.RecipientID]
	for i, id := range recipientNotifs {
		if id == notificationID {
			r.recipientIndex[notification.RecipientID] = append(recipientNotifs[:i], recipientNotifs[i+1:]...)
			break
		}
	}

	delete(r.notifications, notificationID)
	return nil
}

// ListByRecipient retrieves notifications for a recipient with pagination and filters
// Returns: notifications, total count, unread count, error
func (r *InMemoryNotificationRepository) ListByRecipient(
	ctx context.Context,
	recipientID string,
	page, pageSize int,
	notifType domain.NotificationType,
	status domain.NotificationStatus,
	unreadOnly bool,
) ([]*domain.Notification, int, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	notificationIDs := r.recipientIndex[recipientID]
	var filtered []*domain.Notification
	unreadCount := 0

	for _, notifID := range notificationIDs {
		notification, exists := r.notifications[notifID]
		if !exists {
			continue
		}

		// Count unread
		if !notification.IsRead() {
			unreadCount++
		}

		// Filter by type if specified
		if notifType != domain.TypeUnknown && notification.Type != notifType {
			continue
		}

		// Filter by status if specified
		if status != domain.StatusUnknown && notification.Status != status {
			continue
		}

		// Filter unread only
		if unreadOnly && notification.IsRead() {
			continue
		}

		filtered = append(filtered, notification)
	}

	// Sort by created_at descending (newest first)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	total := len(filtered)

	// Apply pagination
	start := (page - 1) * pageSize
	if start >= total {
		return []*domain.Notification{}, total, unreadCount, nil
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	return filtered[start:end], total, unreadCount, nil
}
