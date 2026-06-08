package usecase

import "fmt"

// MockNotificationSender is a mock implementation of NotificationSender for development
type MockNotificationSender struct{}

// NewMockNotificationSender creates a new mock notification sender
func NewMockNotificationSender() *MockNotificationSender {
	return &MockNotificationSender{}
}

// SendEmail mock sends an email
func (s *MockNotificationSender) SendEmail(recipientEmail, title, message string) error {
	// In production, integrate with email service (SendGrid, AWS SES, etc.)
	fmt.Printf("[MOCK EMAIL] To: %s, Title: %s, Message: %s\n", recipientEmail, title, message)
	return nil
}

// SendSMS mock sends an SMS
func (s *MockNotificationSender) SendSMS(recipientPhone, message string) error {
	// In production, integrate with SMS service (Twilio, AWS SNS, etc.)
	fmt.Printf("[MOCK SMS] To: %s, Message: %s\n", recipientPhone, message)
	return nil
}

// SendPush mock sends a push notification
func (s *MockNotificationSender) SendPush(recipientID, title, message, data string) error {
	// In production, integrate with push service (Firebase, OneSignal, etc.)
	fmt.Printf("[MOCK PUSH] To: %s, Title: %s, Message: %s\n", recipientID, title, message)
	return nil
}
