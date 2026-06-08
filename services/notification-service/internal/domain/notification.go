package domain

import (
	"time"
)

// NotificationType represents the type of notification
type NotificationType int

const (
	TypeUnknown NotificationType = iota
	TypeEmail
	TypeSMS
	TypePush
	TypeInApp
)

// String returns the string representation of NotificationType
func (t NotificationType) String() string {
	switch t {
	case TypeEmail:
		return "email"
	case TypeSMS:
		return "sms"
	case TypePush:
		return "push"
	case TypeInApp:
		return "in_app"
	default:
		return "unknown"
	}
}

// NotificationPriority represents the priority level
type NotificationPriority int

const (
	PriorityLow NotificationPriority = iota
	PriorityNormal
	PriorityHigh
	PriorityUrgent
)

// String returns the string representation of NotificationPriority
func (p NotificationPriority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityNormal:
		return "normal"
	case PriorityHigh:
		return "high"
	case PriorityUrgent:
		return "urgent"
	default:
		return "low"
	}
}

// NotificationStatus represents the status of a notification
type NotificationStatus int

const (
	StatusUnknown NotificationStatus = iota
	StatusPending
	StatusSent
	StatusDelivered
	StatusFailed
	StatusRead
)

// String returns the string representation of NotificationStatus
func (s NotificationStatus) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusSent:
		return "sent"
	case StatusDelivered:
		return "delivered"
	case StatusFailed:
		return "failed"
	case StatusRead:
		return "read"
	default:
		return "unknown"
	}
}

// Notification represents a notification entity
type Notification struct {
	NotificationID string
	RecipientID    string
	RecipientEmail string
	RecipientPhone string
	Type           NotificationType
	Priority       NotificationPriority
	Title          string
	Message        string
	Data           string
	Status         NotificationStatus
	SentAt         *time.Time
	ReadAt         *time.Time
	CreatedAt      time.Time
}

// Validate validates the notification
func (n *Notification) Validate() error {
	if n.RecipientID == "" {
		return ErrNotificationRecipientRequired
	}

	if n.Type == TypeUnknown {
		return ErrNotificationTypeRequired
	}

	// Validate recipient info based on type
	switch n.Type {
	case TypeEmail:
		if n.RecipientEmail == "" {
			return ErrNotificationEmailRequired
		}
		if len(n.RecipientEmail) > 200 {
			return ErrNotificationEmailTooLong
		}
	case TypeSMS:
		if n.RecipientPhone == "" {
			return ErrNotificationPhoneRequired
		}
		if len(n.RecipientPhone) > 20 {
			return ErrNotificationPhoneTooLong
		}
	}

	if n.Title == "" {
		return ErrNotificationTitleRequired
	}

	if len(n.Title) > 200 {
		return ErrNotificationTitleTooLong
	}

	if n.Message == "" {
		return ErrNotificationMessageRequired
	}

	if len(n.Message) > 1000 {
		return ErrNotificationMessageTooLong
	}

	if len(n.Data) > 5000 {
		return ErrNotificationDataTooLong
	}

	return nil
}

// MarkAsSent marks the notification as sent
func (n *Notification) MarkAsSent() {
	n.Status = StatusSent
	now := time.Now()
	n.SentAt = &now
}

// MarkAsDelivered marks the notification as delivered
func (n *Notification) MarkAsDelivered() {
	n.Status = StatusDelivered
	if n.SentAt == nil {
		now := time.Now()
		n.SentAt = &now
	}
}

// MarkAsFailed marks the notification as failed
func (n *Notification) MarkAsFailed() {
	n.Status = StatusFailed
}

// MarkAsRead marks the notification as read
func (n *Notification) MarkAsRead() {
	if n.Status != StatusRead {
		n.Status = StatusRead
		now := time.Now()
		n.ReadAt = &now
	}
}

// IsRead checks if the notification has been read
func (n *Notification) IsRead() bool {
	return n.Status == StatusRead
}

// IsSent checks if the notification has been sent
func (n *Notification) IsSent() bool {
	return n.Status == StatusSent || n.Status == StatusDelivered || n.Status == StatusRead
}

// IsDelivered checks if the notification has been delivered
func (n *Notification) IsDelivered() bool {
	return n.Status == StatusDelivered || n.Status == StatusRead
}

// IsFailed checks if the notification has failed
func (n *Notification) IsFailed() bool {
	return n.Status == StatusFailed
}

// IsUrgent checks if the notification is urgent or high priority
func (n *Notification) IsUrgent() bool {
	return n.Priority == PriorityUrgent || n.Priority == PriorityHigh
}

// CanBeRead checks if the notification can be marked as read
func (n *Notification) CanBeRead() bool {
	return n.IsSent() && !n.IsRead()
}

// NewNotification creates a new notification
func NewNotification(
	recipientID, recipientEmail, recipientPhone string,
	notifType NotificationType,
	priority NotificationPriority,
	title, message, data string,
) (*Notification, error) {
	notification := &Notification{
		RecipientID:    recipientID,
		RecipientEmail: recipientEmail,
		RecipientPhone: recipientPhone,
		Type:           notifType,
		Priority:       priority,
		Title:          title,
		Message:        message,
		Data:           data,
		Status:         StatusPending,
		CreatedAt:      time.Now(),
	}

	if err := notification.Validate(); err != nil {
		return nil, err
	}

	return notification, nil
}
