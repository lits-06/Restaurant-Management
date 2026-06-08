# Payment Service

Payment Service handles payment processing, refunds, and payment tracking for the restaurant management system.

## Overview

- **Port**: 50056
- **Protocol**: gRPC
- **Architecture**: Clean Architecture
- **Repository**: In-memory (ready for PostgreSQL implementation)

## Features

### Payment Management
- Create payments for orders
- Process payments with transaction tracking
- Support multiple payment methods (Cash, Credit Card, Debit Card, Mobile Wallet, Bank Transfer)
- Track payment status workflow (Pending → Processing → Completed/Failed)
- Payment amount validation against order total

### Refund Processing
- Full refunds
- Partial refunds
- Refund tracking and history
- Automatic status updates (Partially Refunded → Fully Refunded)

### Payment Queries
- Get payment by ID
- List payments with pagination and filters (status, method, date range)
- Get all payments for an order
- Calculate total paid for an order

### Integration
- **Order Service**: Validates orders and updates order status when payment completes
- Ensures payment amount matches order total
- Prevents payment for cancelled orders

## Status Workflow

### Payment Status States
1. **Pending**: Payment created, awaiting processing
2. **Processing**: Payment is being processed
3. **Completed**: Payment successfully completed
4. **Failed**: Payment processing failed
5. **Partially Refunded**: Part of payment has been refunded
6. **Refunded**: Full payment has been refunded

### Allowed Transitions
```
Pending → Processing → Completed → Partially Refunded → Refunded
       ↓             ↓
     Failed        Failed
       ↓
    Pending (retry)
```

## Payment Methods

- **Cash**: Immediate completion
- **Credit Card**: Requires transaction ID, processed via payment gateway
- **Debit Card**: Similar to credit card
- **Mobile Wallet**: Digital wallet payments (e.g., Momo, ZaloPay)
- **Bank Transfer**: Direct bank transfer

## Business Rules

### Payment Creation
- Order must exist and not be cancelled
- Payment amount must match or exceed order total (can include tip)
- Payment method must be specified
- Customer name optional but recommended

### Payment Processing
- Requires transaction ID (except cash)
- Cash payments complete immediately
- Card/wallet payments go through processing state
- On completion, order status updated to "Completed"

### Refunds
- Only completed or partially refunded payments can be refunded
- Refund amount must not exceed remaining balance
- Full refund sets status to "Refunded"
- Partial refund sets status to "Partially Refunded"
- Refund reason stored in notes

## gRPC Methods

### CreatePayment
Creates a new payment for an order.

**Request:**
```protobuf
{
  "order_id": "order-123",
  "amount": 150000.00,
  "tip": 10000.00,
  "method": "METHOD_CREDIT_CARD",
  "customer_name": "John Doe",
  "notes": "Table 5"
}
```

**Response:**
```protobuf
{
  "payment": {...},
  "success": true,
  "message": "payment created successfully"
}
```

### GetPayment
Retrieves a payment by ID.

**Request:**
```protobuf
{
  "payment_id": "payment-123"
}
```

### ProcessPayment
Processes a pending payment.

**Request:**
```protobuf
{
  "payment_id": "payment-123",
  "transaction_id": "TXN-456789"
}
```

### RefundPayment
Processes a refund.

**Request:**
```protobuf
{
  "payment_id": "payment-123",
  "amount": 50000.00,
  "reason": "Customer complaint"
}
```

### ListPayments
Lists payments with filters.

**Request:**
```protobuf
{
  "page": 1,
  "page_size": 20,
  "status": "STATUS_COMPLETED",
  "method": "METHOD_CASH",
  "from_date": "2024-01-01T00:00:00Z",
  "to_date": "2024-01-31T23:59:59Z"
}
```

### GetPaymentsByOrder
Gets all payments for an order.

**Request:**
```protobuf
{
  "order_id": "order-123"
}
```

## Environment Variables

```bash
ENVIRONMENT=development
GRPC_PORT=50056
LOG_LEVEL=info
ORDER_SERVICE_URL=localhost:50055
```

## Testing

### Using grpcurl

1. Create a payment:
```bash
grpcurl -plaintext -d '{
  "order_id": "order-123",
  "amount": 150000.00,
  "tip": 10000.00,
  "method": "METHOD_CASH",
  "customer_name": "John Doe"
}' localhost:50056 payment.PaymentService/CreatePayment
```

2. Process payment:
```bash
grpcurl -plaintext -d '{
  "payment_id": "payment-id",
  "transaction_id": "TXN-12345"
}' localhost:50056 payment.PaymentService/ProcessPayment
```

3. Refund payment:
```bash
grpcurl -plaintext -d '{
  "payment_id": "payment-id",
  "amount": 50000.00,
  "reason": "Customer request"
}' localhost:50056 payment.PaymentService/RefundPayment
```

4. List payments:
```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10,
  "status": "STATUS_COMPLETED"
}' localhost:50056 payment.PaymentService/ListPayments
```

## Architecture

### Domain Layer (`internal/domain/`)
- **payment.go**: Payment entity with validation, status workflow, refund logic
- **errors.go**: Domain-specific errors

### Repository Layer (`internal/repository/`)
- **repository.go**: PaymentRepository interface
- **payment_memory.go**: In-memory implementation with order index

### UseCase Layer (`internal/usecase/`)
- **payment_usecase.go**: Business logic, order validation, payment processing
- Integrates with Order Service for validation and status updates

### Delivery Layer (`internal/delivery/grpc/`)
- **payment_handler.go**: gRPC handlers, proto conversion

### Server (`cmd/server/`)
- **main.go**: Server setup, dependency injection, gRPC client connections

## Development

### Run locally
```bash
cd services/payment-service
go run cmd/server/main.go
```

### Run with Docker
```bash
docker build -t payment-service -f services/payment-service/Dockerfile .
docker run -p 50056:50056 payment-service
```

### Run with dependencies (Order Service required)
```bash
docker-compose up payment-service
```

## Integration Requirements

**Requires:**
- Order Service (port 50055) - for order validation and status updates

**Used by:**
- API Gateway
- Report Service (analytics)

## Future Enhancements

- [ ] Real payment gateway integration (Stripe, PayPal, VNPay)
- [ ] Payment retry mechanism for failed payments
- [ ] Webhook support for async payment notifications
- [ ] Split payments (multiple payment methods for one order)
- [ ] Payment scheduling/installments
- [ ] Currency conversion support
- [ ] Payment receipt generation
- [ ] Fraud detection
- [ ] PostgreSQL repository implementation
- [ ] Payment analytics and reporting
