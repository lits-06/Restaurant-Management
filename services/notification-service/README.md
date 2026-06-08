# Notification Service

Notification Service handles multi-channel notifications (Email, SMS, Push, In-App) for the restaurant management system.

## Overview

- **Port**: 50058
- **Protocol**: gRPC
- **Architecture**: Clean Architecture
- **Repository**: In-memory (ready for PostgreSQL implementation)

## Features

### Multi-Channel Notifications
- **Email**: Full email support with title and message
- **SMS**: Text message notifications
- **Push**: Mobile push notifications with custom data
- **In-App**: Application notifications stored in database

### Notification Management
- Send individual notifications
- Send bulk notifications to multiple recipients
- Read/unread tracking
- Notification history
- Priority levels (Low, Normal, High, Urgent)
- Status tracking (Pending, Sent, Delivered, Failed, Read)

### Recipient Management
- Notifications linked to user IDs
- Email and phone number support
- Bulk recipient targeting
- Recipient filtering

### Advanced Features
- Priority-based notifications
- Status workflow tracking
- Read receipts
- Notification deletion
- Unread count tracking
- Type and status filtering

## Notification Types

1. **Email**: Sent via email service (SendGrid, AWS SES)
2. **SMS**: Sent via SMS service (Twilio, AWS SNS)
3. **Push**: Sent via push service (Firebase, OneSignal)
4. **In-App**: Stored notifications displayed within app

## Priority Levels

- **Low**: Regular updates, can be delayed
- **Normal**: Standard notifications (default)
- **High**: Important notifications, sent immediately
- **Urgent**: Critical alerts, highest priority

## Status Workflow

```
Pending → Sent → Delivered → Read
       ↓
     Failed
```

- **Pending**: Created but not yet sent
- **Sent**: Successfully sent to delivery service
- **Delivered**: Confirmed delivered to recipient
- **Failed**: Delivery failed
- **Read**: Recipient has read the notification

## Business Rules

### Notification Creation
- Recipient ID required
- Type must be specified
- Email required for email notifications
- Phone required for SMS notifications
- Title required (max 200 characters)
- Message required (max 1000 characters)
- Data optional (max 5000 characters)

### Sending
- Notifications sent immediately upon creation
- Failed notifications marked with failed status
- Successful sends marked as sent
- In-app notifications don't require external sending

### Read Tracking
- Only sent/delivered notifications can be marked as read
- Read timestamp recorded automatically
- Read status is permanent (cannot unread)

### Bulk Operations
- Multiple recipients supported
- Individual success/failure tracking
- Returns sent count and failed count

## gRPC Methods

### SendNotification
Sends a notification to a single recipient.

**Request:**
```protobuf
{
  "recipient_id": "user-123",
  "recipient_email": "user@example.com",
  "recipient_phone": "+1234567890",
  "type": "TYPE_EMAIL",
  "priority": "PRIORITY_HIGH",
  "title": "Order Ready",
  "message": "Your order #456 is ready for pickup",
  "data": "{\"order_id\":\"456\"}"
}
```

### SendBulkNotification
Sends notifications to multiple recipients.

**Request:**
```protobuf
{
  "recipient_ids": ["user-1", "user-2", "user-3"],
  "type": "TYPE_PUSH",
  "priority": "PRIORITY_NORMAL",
  "title": "New Menu Items",
  "message": "Check out our new summer menu!",
  "data": "{}"
}
```

**Response:**
```protobuf
{
  "sent_count": 2,
  "failed_count": 1,
  "success": true,
  "message": "bulk notifications processed"
}
```

### GetNotification
Retrieves a notification by ID.

**Request:**
```protobuf
{
  "notification_id": "notif-123"
}
```

### ListNotifications
Lists notifications for a recipient with filters.

**Request:**
```protobuf
{
  "recipient_id": "user-123",
  "page": 1,
  "page_size": 20,
  "type": "TYPE_EMAIL",
  "status": "STATUS_UNREAD",
  "unread_only": true
}
```

**Response:**
```protobuf
{
  "notifications": [...],
  "total": 50,
  "unread_count": 12,
  "page": 1,
  "page_size": 20,
  "success": true
}
```

### MarkAsRead
Marks a notification as read.

**Request:**
```protobuf
{
  "notification_id": "notif-123"
}
```

### DeleteNotification
Deletes a notification.

**Request:**
```protobuf
{
  "notification_id": "notif-123"
}
```

## Environment Variables

```bash
ENVIRONMENT=development
GRPC_PORT=50058
LOG_LEVEL=info
```

## Testing

### Using grpcurl

1. Send email notification:
```bash
grpcurl -plaintext -d '{
  "recipient_id": "user-123",
  "recipient_email": "user@example.com",
  "type": "TYPE_EMAIL",
  "priority": "PRIORITY_HIGH",
  "title": "Order Confirmation",
  "message": "Your order has been confirmed"
}' localhost:50058 notification.NotificationService/SendNotification
```

2. Send bulk push notification:
```bash
grpcurl -plaintext -d '{
  "recipient_ids": ["user-1", "user-2"],
  "type": "TYPE_PUSH",
  "priority": "PRIORITY_NORMAL",
  "title": "Special Offer",
  "message": "20% off today only!"
}' localhost:50058 notification.NotificationService/SendBulkNotification
```

3. List unread notifications:
```bash
grpcurl -plaintext -d '{
  "recipient_id": "user-123",
  "page": 1,
  "page_size": 10,
  "unread_only": true
}' localhost:50058 notification.NotificationService/ListNotifications
```

4. Mark as read:
```bash
grpcurl -plaintext -d '{
  "notification_id": "notif-id"
}' localhost:50058 notification.NotificationService/MarkAsRead
```

## Architecture

### Domain Layer (`internal/domain/`)
- **notification.go**: Notification entity with types, priorities, statuses
- **errors.go**: Domain-specific errors

### Repository Layer (`internal/repository/`)
- **repository.go**: NotificationRepository interface
- **notification_memory.go**: In-memory implementation with recipient index

### UseCase Layer (`internal/usecase/`)
- **notification_usecase.go**: Business logic for notification management
- **mock_sender.go**: Mock implementation of notification senders

### Delivery Layer (`internal/delivery/grpc/`)
- **notification_handler.go**: gRPC handlers, proto conversion, 6 methods

### Server (`cmd/server/`)
- **main.go**: Server setup, dependency injection

## Development

### Run locally
```bash
cd services/notification-service
go run cmd/server/main.go
```

### Run with Docker
```bash
docker build -t notification-service -f services/notification-service/Dockerfile .
docker run -p 50058:50058 notification-service
```

### Run with docker-compose
```bash
docker-compose up notification-service
```

## Integration Points

**Currently Standalone** - No external service dependencies

**Mock Senders:**
- Email: Prints to console (integrate with SendGrid, AWS SES)
- SMS: Prints to console (integrate with Twilio, AWS SNS)
- Push: Prints to console (integrate with Firebase, OneSignal)

**Can be integrated with:**
- Order Service: Send order status notifications
- Payment Service: Send payment confirmations
- Inventory Service: Send low stock alerts
- User Service: Fetch user contact details for bulk sends

## Future Enhancements

- [ ] Real email service integration (SendGrid, AWS SES)
- [ ] Real SMS service integration (Twilio, AWS SNS)
- [ ] Real push notification service (Firebase, OneSignal)
- [ ] Template system for notifications
- [ ] Multi-language support
- [ ] Notification preferences per user
- [ ] Scheduled notifications
- [ ] Retry mechanism for failed sends
- [ ] Delivery tracking and analytics
- [ ] Notification grouping/threading
- [ ] Rich content support (images, buttons)
- [ ] Webhook support for status updates
- [ ] Rate limiting per recipient
- [ ] PostgreSQL repository implementation
- [ ] Notification campaigns
