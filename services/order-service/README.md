# Order Service

Order Service quản lý đơn hàng cho hệ thống nhà hàng với inter-service communication.

## 🎯 Features

### Order Management
- ✅ Create orders with menu items
- ✅ Update order details (items, discount, notes)
- ✅ Cancel orders
- ✅ List orders with filters (status, table, date range)
- ✅ Update order status (7-state workflow)
- ✅ Add/remove items from orders
- ✅ Get orders by table
- ✅ Automatic total calculation (subtotal + tax - discount)

### Order Status Workflow
- **PENDING** → CONFIRMED → PREPARING → READY → SERVED → COMPLETED
- **Any status** → CANCELLED (except COMPLETED)

### Inter-Service Integration
- ✅ Menu Service - Validate menu items and prices
- ✅ Table Service - Validate table availability, update table status
- ✅ User Service - Validate waiter

### Business Rules
- Orders must have at least one item
- Only pending/confirmed orders can be modified
- Completed orders cannot be cancelled
- Table is marked as occupied when order is created
- Table is marked as available when order is completed/cancelled
- Tax is calculated as 10% of subtotal
- Total = Subtotal + Tax - Discount

## 🏗️ Architecture

Follows **Clean Architecture** with **Inter-Service Communication**:

```
cmd/server/          # Application entry point + gRPC clients
internal/
  ├── domain/        # Business entities (Order, OrderItem)
  ├── repository/    # Data access interfaces + in-memory impl
  ├── usecase/       # Business logic + service client interfaces
  └── delivery/grpc/ # gRPC handlers
pkg/config/          # Configuration
```

### Layers

1. **Domain Layer** (`internal/domain/`)
   - `order.go` - Order entity with 7 statuses
   - `order_item.go` - OrderItem entity
   - `errors.go` - Domain-specific errors
   - Zero external dependencies

2. **Repository Layer** (`internal/repository/`)
   - `repository.go` - Repository interface
   - `order_memory.go` - In-memory repository with table index
   - Thread-safe with RWMutex
   - Sorted results by created_at

3. **UseCase Layer** (`internal/usecase/`)
   - `order_usecase.go` - Business logic orchestration
   - 11 operations total
   - Service client interfaces for Menu, Table, User
   - Inter-service validation and communication

4. **Delivery Layer** (`internal/delivery/grpc/`)
   - `order_handler.go` - gRPC handlers
   - 9 RPC methods
   - Proto ↔ Domain conversion

5. **gRPC Clients** (`cmd/server/main.go`)
   - Menu Service client
   - Table Service client
   - User Service client
   - Validation helpers

## 🚀 Running the Service

### Prerequisites

**Required Services:**
- Menu Service (localhost:50054)
- Table Service (localhost:50053)
- User Service (localhost:50052)

### Local Development

```bash
# From project root
cd services/order-service

# Run service
go run cmd/server/main.go

# Service will start on port 50055
```

### Docker

```bash
# Build image
docker build -t order-service -f services/order-service/Dockerfile .

# Run container
docker run -p 50055:50055 \
  -e MENU_SERVICE_ADDR=menu-service:50054 \
  -e TABLE_SERVICE_ADDR=table-service:50053 \
  -e USER_SERVICE_ADDR=user-service:50052 \
  order-service
```

### Docker Compose

```bash
# From project root
docker-compose up order-service
```

## 📡 gRPC API

### CreateOrder
```bash
grpcurl -plaintext -d '{
  "table_id": "TABLE_ID",
  "waiter_id": "WAITER_ID",
  "items": [
    {"menu_item_id": "ITEM_ID_1", "quantity": 2, "notes": "No onions"},
    {"menu_item_id": "ITEM_ID_2", "quantity": 1}
  ],
  "notes": "Customer prefers window seat"
}' localhost:50055 order.OrderService/CreateOrder
```

### GetOrder
```bash
grpcurl -plaintext -d '{
  "order_id": "ORDER_ID"
}' localhost:50055 order.OrderService/GetOrder
```

### UpdateOrder
```bash
grpcurl -plaintext -d '{
  "order_id": "ORDER_ID",
  "discount": 10000,
  "notes": "Updated notes"
}' localhost:50055 order.OrderService/UpdateOrder
```

### CancelOrder
```bash
grpcurl -plaintext -d '{
  "order_id": "ORDER_ID",
  "reason": "Customer requested cancellation"
}' localhost:50055 order.OrderService/CancelOrder
```

### ListOrders
```bash
# List all orders
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:50055 order.OrderService/ListOrders

# Filter by status
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10,
  "status": "STATUS_PREPARING"
}' localhost:50055 order.OrderService/ListOrders

# Filter by table
grpcurl -plaintext -d '{
  "table_id": "TABLE_ID",
  "status": "STATUS_PENDING"
}' localhost:50055 order.OrderService/ListOrders
```

### UpdateOrderStatus
```bash
# Confirm order
grpcurl -plaintext -d '{
  "order_id": "ORDER_ID",
  "status": "STATUS_CONFIRMED"
}' localhost:50055 order.OrderService/UpdateOrderStatus

# Mark as preparing
grpcurl -plaintext -d '{
  "order_id": "ORDER_ID",
  "status": "STATUS_PREPARING"
}' localhost:50055 order.OrderService/UpdateOrderStatus

# Mark as ready
grpcurl -plaintext -d '{
  "order_id": "ORDER_ID",
  "status": "STATUS_READY"
}' localhost:50055 order.OrderService/UpdateOrderStatus

# Complete order
grpcurl -plaintext -d '{
  "order_id": "ORDER_ID",
  "status": "STATUS_COMPLETED"
}' localhost:50055 order.OrderService/UpdateOrderStatus
```

### AddOrderItem
```bash
grpcurl -plaintext -d '{
  "order_id": "ORDER_ID",
  "item": {
    "menu_item_id": "ITEM_ID",
    "quantity": 1,
    "notes": "Extra spicy"
  }
}' localhost:50055 order.OrderService/AddOrderItem
```

### RemoveOrderItem
```bash
grpcurl -plaintext -d '{
  "order_id": "ORDER_ID",
  "item_id": "ITEM_ID"
}' localhost:50055 order.OrderService/RemoveOrderItem
```

### GetOrdersByTable
```bash
grpcurl -plaintext -d '{
  "table_id": "TABLE_ID",
  "status": "STATUS_PENDING"
}' localhost:50055 order.OrderService/GetOrdersByTable
```

## 🔧 Configuration

Environment variables:

```bash
GRPC_PORT=50055                           # gRPC server port (default: 50055)
LOG_LEVEL=info                            # Logging level (default: info)
MENU_SERVICE_ADDR=localhost:50054         # Menu service address
TABLE_SERVICE_ADDR=localhost:50053        # Table service address
USER_SERVICE_ADDR=localhost:50052         # User service address
```

## 📊 Data Models

### Order
```go
type Order struct {
    OrderID   string
    TableID   string
    WaiterID  string
    Items     []*OrderItem
    Subtotal  float64
    Tax       float64      // 10% of subtotal
    Discount  float64
    Total     float64      // Subtotal + Tax - Discount
    Status    OrderStatus
    Notes     string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### OrderItem
```go
type OrderItem struct {
    ItemID       string
    MenuItemID   string
    MenuItemName string
    Quantity     int32
    UnitPrice    float64
    Subtotal     float64  // UnitPrice * Quantity
    Notes        string
    Status       OrderStatus
}
```

### OrderStatus
- `STATUS_UNKNOWN` (0)
- `STATUS_PENDING` (1) - Order created
- `STATUS_CONFIRMED` (2) - Order confirmed by staff
- `STATUS_PREPARING` (3) - Kitchen is preparing
- `STATUS_READY` (4) - Order ready for serving
- `STATUS_SERVED` (5) - Order served to customer
- `STATUS_COMPLETED` (6) - Order completed and paid
- `STATUS_CANCELLED` (7) - Order cancelled

## 🔄 Status Transitions

```
PENDING → CONFIRMED → PREPARING → READY → SERVED → COMPLETED
   ↓          ↓           ↓          ↓        ↓
   └──────────┴───────────┴──────────┴────────┴─→ CANCELLED
```

**Rules:**
- PENDING can go to CONFIRMED or CANCELLED
- CONFIRMED can go to PREPARING or CANCELLED
- PREPARING can go to READY or CANCELLED
- READY can go to SERVED or CANCELLED
- SERVED can go to COMPLETED or CANCELLED
- COMPLETED can only go to CANCELLED (refund scenario)
- CANCELLED is final

## 🧪 Testing

```bash
# Unit tests
go test ./internal/...

# Integration tests
go test ./...

# With coverage
go test -cover ./...
```

## 🎯 Use Cases

### Typical Order Flow

1. **Customer Arrives**
   ```
   Table Service: Get available table
   User Service: Get waiter
   ```

2. **Create Order**
   ```
   Menu Service: Validate menu items and get prices
   Order Service: Create order with items
   Table Service: Mark table as occupied
   ```

3. **Order Processing**
   ```
   Order Status: PENDING → CONFIRMED (waiter confirms)
   Order Status: CONFIRMED → PREPARING (kitchen starts)
   Order Status: PREPARING → READY (kitchen finished)
   Order Status: READY → SERVED (waiter serves)
   ```

4. **Complete Order**
   ```
   Payment Service: Process payment
   Order Status: SERVED → COMPLETED
   Table Service: Mark table as available
   ```

## 📝 Notes

- **Inter-Service Dependencies**: Order Service requires Menu, Table, and User services to be running
- **In-Memory Storage**: Current implementation uses in-memory repository. Ready for PostgreSQL.
- **Table Management**: Automatically updates table status (occupied when order created, available when completed/cancelled)
- **Price Locking**: Menu item prices are locked when order is created
- **Tax Calculation**: Fixed at 10% of subtotal
- **Discount Validation**: Cannot exceed subtotal
- **Quantity Limits**: Max 100 items per line item
- **Thread-Safe**: All repository operations are thread-safe

## 🔄 Future Enhancements

- [ ] PostgreSQL repository implementation
- [ ] Redis caching for active orders
- [ ] Order splitting (split bill)
- [ ] Order merging (combine tables)
- [ ] Tip/gratuity support
- [ ] Custom tax rates
- [ ] Promotional codes
- [ ] Order history archiving
- [ ] Real-time notifications (kitchen display, waiter alerts)
- [ ] Order preparation time tracking

## 📚 Dependencies

- `google.golang.org/grpc` - gRPC framework
- `google.golang.org/protobuf` - Protocol Buffers
- `github.com/google/uuid` - UUID generation
- Shared packages: logger, middleware
- External services: Menu, Table, User

## 🔗 Service Interactions

```
Order Service (50055)
├── → Menu Service (50054)
│   └── ValidateMenuItem (check availability, get price)
├── → Table Service (50053)
│   ├── ValidateTable (check availability)
│   └── UpdateTableStatus (mark occupied/available)
└── → User Service (50052)
    └── ValidateWaiter (check user exists and is active)
```

## 📄 License

Part of Restaurant Management System
