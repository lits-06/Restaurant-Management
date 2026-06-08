package domain

import "errors"

// Notification errors
var (
	ErrNotificationNotFound          = errors.New("notification not found")
	ErrNotificationRecipientRequired = errors.New("recipient ID is required")
	ErrNotificationTypeRequired      = errors.New("notification type is required")
	ErrNotificationEmailRequired     = errors.New("recipient email is required for email notifications")
	ErrNotificationEmailTooLong      = errors.New("recipient email is too long (max 200 characters)")
	ErrNotificationPhoneRequired     = errors.New("recipient phone is required for SMS notifications")
	ErrNotificationPhoneTooLong      = errors.New("recipient phone is too long (max 20 characters)")
	ErrNotificationTitleRequired     = errors.New("notification title is required")
	ErrNotificationTitleTooLong      = errors.New("notification title is too long (max 200 characters)")
	ErrNotificationMessageRequired   = errors.New("notification message is required")
	ErrNotificationMessageTooLong    = errors.New("notification message is too long (max 1000 characters)")
	ErrNotificationDataTooLong       = errors.New("notification data is too long (max 5000 characters)")
	ErrNotificationAlreadyRead       = errors.New("notification is already read")
	ErrNotificationCannotRead        = errors.New("notification cannot be marked as read")
	ErrNotificationSendFailed        = errors.New("failed to send notification")
)
