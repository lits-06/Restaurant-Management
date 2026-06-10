# Order Service

Microservice quản lý đơn đặt bàn (reservation + order) cho hệ thống nhà hàng.

**Port:** `50055` | **Protocol:** gRPC | **DB:** PostgreSQL (`restaurant_db`)

---

## Trách nhiệm

- Tạo và quản lý đơn đặt bàn (tên, SĐT, thời gian, số người, ghi chú, món ăn)
- **Tự động gán bàn** khi khách không chỉ định bàn (gọi table-service + kiểm tra xung đột giờ)
- Validate món ăn và lấy giá qua menu-service gRPC
- Theo dõi trạng thái order: Pending → Confirmed → Completed (hoặc Cancelled)

---

## Order Entity

| Field | Type | Ghi chú |
|-------|------|---------|
| `order_id` | string | UUID, DB-generated |
| `table_id` | string | Auto-assigned nếu không truyền vào |
| `name` | string | Tên khách hàng (bắt buộc) |
| `phone` | string | SĐT khách hàng (bắt buộc) |
| `notes` | string | Ghi chú đặc biệt (dị ứng, yêu cầu ghế…) — luôn overwrite khi update |
| `time` | timestamp | Giờ bắt đầu |
| `end_time` | timestamp | Giờ kết thúc |
| `party_size` | int32 | Số người (bắt buộc, > 0) |
| `status` | string | Pending / Confirmed / Completed / Cancelled |
| `total` | float64 | Tổng tiền (sum of items) |
| `items` | OrderItem[] | Danh sách món kèm số lượng |

### OrderItem

| Field | Ghi chú |
|-------|---------|
| `item_id` | UUID của order item |
| `name` | Tên món (lấy từ menu lúc tạo) |
| `price` | Giá tại thời điểm tạo order |
| `category` | Category của món |
| `image_url` | URL ảnh |
| `quantity` | Số lượng |

---

## Status workflow

```
Pending → Confirmed → Completed
   ↓           ↓
  Cancelled  Cancelled
```

- `Pending`: mới tạo
- `Confirmed`: admin xác nhận
- `Completed`: order hoàn tất
- `Cancelled`: từ Pending hoặc Confirmed (không thể cancel Completed)

---

## 9 RPCs

### `CreateOrder`

```json
Request:
{
  "name": "Nguyen Van A",
  "phone": "0901234567",
  "date": "2026-06-15",
  "time": "19:00",
  "end_time": "21:00",
  "party_size": 4,
  "notes": "dị ứng hải sản, cần ghế cao cho trẻ em",
  "table_id": "",          // optional — bỏ trống để auto-assign
  "items": [{ "item_id": "UUID", "quantity": 2 }]
}
```

**Nếu `table_id` bỏ trống** → auto-assign:
1. Gọi `table-service.GetAvailableTables(min_capacity=party_size)` → bàn AVAILABLE ≥ party_size, sorted capacity ASC
2. Query `orders` DB: bàn nào trùng khung giờ `[time, end_time)` với status != Cancelled
3. Gán bàn nhỏ nhất không bị xung đột (best-fit)
4. Trả lỗi `ErrNoTableAvailable` nếu không còn bàn phù hợp

**Nếu `TABLE_SERVICE_ADDR` không được set** → auto-assign bị tắt, `table_id` bắt buộc.

### `GetOrder`
```
Request:  { order_id: string }
Response: { order: Order, success: bool, message: string }
```

### `UpdateOrder`
```
Request:  { order_id, name, phone, date, time, end_time, party_size, table_id, notes, items }
Response: { order: Order, success: bool, message: string }
```
Chỉ update Pending/Confirmed orders. `notes` luôn được overwrite (truyền rỗng = xóa notes).

### `DeleteOrder`
```
Request:  { order_id: string }
Response: { success: bool, message: string }
```

### `CancelOrder`
```
Request:  { order_id: string, reason: string }
Response: { success: bool, message: string }
```

### `ListOrders`
```
Request:  { page: int32, page_size: int32, status: string, keyword: string }
Response: { orders: []Order, total, page, page_size, success, message }
```
Filter theo status và/hoặc keyword (tên, SĐT).

### `UpdateOrderStatus`
```
Request:  { order_id: string, status: string }
Response: { order: Order, success: bool, message: string }
```

### `AddOrderItem`
```
Request:  { order_id: string, item: { item_id, quantity } }
Response: { order: Order, success: bool, message: string }
```
Validate item qua menu-service. Chỉ áp dụng cho Pending/Confirmed orders.

### `RemoveOrderItem`
```
Request:  { order_id: string, item_id: string }
Response: { order: Order, success: bool, message: string }
```

---

## Database

```sql
CREATE TABLE orders (
    order_id   VARCHAR(36) PRIMARY KEY,
    table_id   VARCHAR(36) NOT NULL DEFAULT '',
    name       VARCHAR(255) NOT NULL,
    phone      VARCHAR(50) NOT NULL,
    notes      TEXT NOT NULL DEFAULT '',
    time       TIMESTAMP NOT NULL,
    end_time   TIMESTAMP,
    party_size INTEGER NOT NULL,
    status     VARCHAR(32) NOT NULL DEFAULT 'Pending',
    total      NUMERIC(10,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE order_items (
    item_id    VARCHAR(36) PRIMARY KEY,
    order_id   VARCHAR(36) NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    price      NUMERIC(10,2) NOT NULL,
    category   VARCHAR(100) NOT NULL DEFAULT '',
    image_url  TEXT NOT NULL DEFAULT '',
    quantity   INTEGER NOT NULL
);
```

Schema auto-created on startup (`ensureSchema`).

---

## Inter-service communication

```
order-service (50055)
├── → menu-service (50054)
│   └── GetMenuItem — validate item, lấy name/price/category/image
└── → table-service (50053)  [optional — requires TABLE_SERVICE_ADDR]
    └── GetAvailableTables(min_capacity) — danh sách bàn AVAILABLE để auto-assign
```

**Adapter pattern:** `tableAdapter` và `menuAdapter` trong `cmd/server/main.go` implements các interface trong `usecase` layer.

---

## Cấu trúc thư mục

```
services/order-service/
├── cmd/server/main.go           # Entry point — wires adapters + gRPC server
├── internal/
│   ├── domain/
│   │   ├── order.go             # Order entity, status workflow, validation
│   │   ├── order_item.go        # OrderItem entity
│   │   └── errors.go            # Domain errors (ErrNoTableAvailable, ErrTableRequired…)
│   ├── repository/
│   │   ├── repository.go        # OrderRepository interface + GetOccupiedTableIDs
│   │   └── order_postgres.go    # PostgreSQL implementation
│   ├── usecase/
│   │   └── order_usecase.go     # Business logic, autoAssignTable, TableServiceClient interface
│   └── delivery/grpc/
│       └── order_handler.go     # gRPC handlers, proto ↔ domain mapping
└── pkg/config/config.go
```

---

## Configuration

| Env var | Default | Mô tả |
|---------|---------|-------|
| `GRPC_PORT` | `50055` | gRPC server port |
| `DATABASE_HOST` | `localhost` | PostgreSQL host |
| `DATABASE_PORT` | `5432` | |
| `DATABASE_USER` | `restaurant_user` | |
| `DATABASE_PASSWORD` | `restaurant_pass` | |
| `DATABASE_NAME` | `restaurant_db` | |
| `DATABASE_SSLMODE` | `disable` | |
| `MENU_SERVICE_ADDR` | `localhost:50054` | Menu service gRPC address |
| `TABLE_SERVICE_ADDR` | `""` | Table service address — **rỗng = auto-assign bị tắt** |
| `LOG_LEVEL` | `info` | |

---

## Build & Run

```bash
cd services/order-service
go build ./cmd/server/
go run cmd/server/main.go
```

Cross-compile cho Alpine container (docker-compose):
```bash
GOOS=linux GOARCH=amd64 go build -o server ./cmd/server/
```

---

## grpcurl — ví dụ

```bash
# Tạo order — auto-assign bàn (không truyền table_id)
grpcurl -plaintext -d '{
  "name": "Nguyen Van A",
  "phone": "0901234567",
  "date": "2026-06-15",
  "time": "19:00",
  "end_time": "21:00",
  "party_size": 4,
  "notes": "dị ứng hải sản"
}' localhost:50055 order.OrderService/CreateOrder

# Lấy order
grpcurl -plaintext -d '{"order_id":"ORDER_UUID"}' localhost:50055 order.OrderService/GetOrder

# Danh sách orders đang Pending
grpcurl -plaintext -d '{"page":1,"page_size":20,"status":"Pending"}' localhost:50055 order.OrderService/ListOrders

# Xác nhận order
grpcurl -plaintext -d '{"order_id":"ORDER_UUID","status":"Confirmed"}' localhost:50055 order.OrderService/UpdateOrderStatus

# Huỷ order
grpcurl -plaintext -d '{"order_id":"ORDER_UUID","reason":"Khach huy"}' localhost:50055 order.OrderService/CancelOrder
```

---

## Related services

- **table-service** (port 50053) — cung cấp `GetAvailableTables` cho auto-assign
- **menu-service** (port 50054) — validate và lấy thông tin món ăn
