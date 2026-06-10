# Order Service

**Port:** 50055 (gRPC)  
**Module:** `restaurant-management/services/order-service`  
**Proto:** `proto/order/order.proto`  
**Database:** PostgreSQL (`restaurant_db`)

## Tổng quan

Order Service quản lý toàn bộ vòng đời của đơn đặt bàn: từ tạo đặt chỗ, xử lý món ăn trong bếp, đến hoàn thành. Đây là service phức tạp nhất — tích hợp với table-service (auto-assign bàn), menu-service (lookup giá), và notification-service (push thông báo bếp).

## Kiến trúc

```
cmd/server/main.go           ← Entry point + notifAdapter
internal/
  domain/
    order.go                 ← Order entity, state machines, validation
    order_item.go            ← OrderItem + ItemStatus state machine
    errors.go                ← Typed domain errors
    order_repository.go      ← OrderRepository interface
  repository/
    order_postgres.go        ← PostgreSQL: CRUD + GetOccupiedTableIDs + UpdateItemStatus
  usecase/
    order_usecase.go         ← Business logic, auto-assign, notifications
  delivery/grpc/
    order_handler.go         ← gRPC handler
pkg/config/config.go
```

---

## Order Status (4 trạng thái)

```
Pending → Confirmed → Completed
    ↘           ↘
     Cancelled   Cancelled
```

| Status | Mô tả | Có thể chuyển sang |
|--------|-------|-------------------|
| `Pending` | Vừa tạo, chờ xác nhận | Confirmed, Cancelled |
| `Confirmed` | Đã xác nhận, bếp đang xử lý | Completed, Cancelled |
| `Completed` | Đã thanh toán, kết thúc | (không thể chuyển) |
| `Cancelled` | Đã hủy | (không thể chuyển) |

> Completed không thể chuyển sang Cancelled. Cancelled không thể chuyển sang bất kỳ trạng thái nào.

---

## Item Status (4 trạng thái)

```
PENDING → COOKING → READY → SERVED
```

| Status | Người thực hiện | Trigger thông báo |
|--------|-----------------|-------------------|
| `PENDING` | Default khi tạo | — |
| `COOKING` | CHEF (hoặc ADMIN/MANAGER) | — |
| `READY` | CHEF (hoặc ADMIN/MANAGER) | Notification → WAITER (`ITEM_READY`) |
| `SERVED` | WAITER (hoặc ADMIN/MANAGER) | — |

---

## gRPC API — 10 RPCs

### 1. `CreateOrder`

**Request:**
| Field | Type | Required | Mô tả |
|-------|------|----------|-------|
| `name` | string | ✓ | Tên khách hàng |
| `phone` | string | ✓ | SĐT khách |
| `date` | string | ✓ | Format: `YYYY-MM-DD` |
| `time` | string | ✓ | Format: `HH:MM` |
| `end_time` | string | — | Format: `HH:MM`, phải sau `time` |
| `party_size` | int32 | ✓ | > 0 |
| `notes` | string | — | Yêu cầu đặc biệt (dị ứng, v.v.) |
| `table_id` | string | — | UUID. **Bỏ trống = auto-assign** |
| `user_id` | string | — | UUID. Trích từ JWT bởi api-gateway |
| `items` | OrderItemRequest[] | — | `item_id`, `quantity` |

**Auto-assign table** (khi `table_id` trống):
1. Gọi `table-service.GetAvailableTables(min_capacity=party_size)` → danh sách bàn AVAILABLE có đủ chỗ, sắp xếp capacity ASC
2. Query `orders` tìm các `table_id` có order conflict trong khung giờ (status != Cancelled)
3. Chọn bàn đầu tiên không trong tập bận → best-fit

**Response:** `Order`, `success`, `message`

---

### 2. `GetOrder`

**Request:** `order_id`  
**Response:** `Order` (bao gồm `items[]`), `success`, `message`

---

### 3. `UpdateOrder`
Cập nhật thông tin đặt bàn. **`user_id` không được cập nhật** (immutable sau khi tạo).

**Request:** `order_id`, `name`, `phone`, `date`, `time`, `end_time`, `party_size`, `notes`, `table_id`, `items[]`  
**Response:** `Order`, `success`, `message`

> Update xóa và tạo lại toàn bộ `order_items` trong transaction.

---

### 4. `DeleteOrder`

**Request:** `order_id`  
**Response:** `success`, `message`

---

### 5. `CancelOrder`
Chuyển status → `Cancelled`. Không thể cancel order đã `Completed`.

**Request:** `order_id`, `reason`  
**Response:** `success`, `message`

---

### 6. `ListOrders`

**Request:**
| Field | Type | Mô tả |
|-------|------|-------|
| `page` | int32 | Default: 1 |
| `page_size` | int32 | Default: 10 |
| `status` | string | Filter: Pending/Confirmed/Completed/Cancelled |
| `keyword` | string | ILIKE trên `name`, `phone` |
| `user_id` | string | Filter orders của 1 user cụ thể |

**Response:** `orders[]`, `total`, `page`, `page_size`, `success`, `message`

---

### 7. `UpdateOrderStatus`
Chuyển trạng thái order. Chỉ dành cho staff (ADMIN/MANAGER/CHEF/WAITER).

**Request:** `order_id`, `status`  
**Side effect:** Khi status = `Confirmed` → fire-and-forget `notifyChef` (notification-service)  
**Response:** `Order`, `success`, `message`

---

### 8. `AddOrderItem`
Thêm 1 món vào order. Gọi `menu-service.GetMenuItem` để validate và lấy price.

**Request:** `order_id`, `item` (item_id, name, price, quantity)  
**Response:** `Order`, `success`, `message`

---

### 9. `RemoveOrderItem`

**Request:** `order_id`, `item_id`  
**Response:** `Order`, `success`, `message`

---

### 10. `UpdateOrderItemStatus`
Cập nhật trạng thái của 1 món. Kiểm tra state machine trước khi lưu.

**Request:** `order_id`, `item_id`, `item_status`  
**Side effect:** Khi `item_status = READY` → fire-and-forget `notifyWaiter`  
**Response:** `Order`, `success`, `message`

---

## Database Schema

### Bảng `orders`

```sql
CREATE TABLE IF NOT EXISTS orders (
    order_id   VARCHAR(36)        PRIMARY KEY DEFAULT gen_random_uuid(),
    table_id   VARCHAR(36)        NOT NULL DEFAULT '',
    user_id    VARCHAR(36)        NOT NULL DEFAULT '',
    name       VARCHAR(120)       NOT NULL,
    phone      VARCHAR(120)       NOT NULL DEFAULT '',
    notes      TEXT               NOT NULL DEFAULT '',
    time       TIMESTAMP          NOT NULL,
    end_time   TIMESTAMP,                             -- nullable
    party_size INTEGER            NOT NULL,
    status     VARCHAR(32)        NOT NULL,           -- Pending/Confirmed/Completed/Cancelled
    total      DOUBLE PRECISION   NOT NULL DEFAULT 0
);
```

### Bảng `order_items`

```sql
CREATE TABLE IF NOT EXISTS order_items (
    item_id     VARCHAR(36)      PRIMARY KEY REFERENCES menu_items(item_id),
    order_id    VARCHAR(36)      NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    name        VARCHAR(200)     NOT NULL,
    price       DOUBLE PRECISION NOT NULL,
    quantity    INTEGER          NOT NULL,
    item_status VARCHAR(16)      NOT NULL DEFAULT 'PENDING'
);
```

### Indexes

```sql
CREATE INDEX IF NOT EXISTS idx_orders_user_id     ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status      ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_name        ON orders(name);
CREATE INDEX IF NOT EXISTS idx_orders_table_id    ON orders(table_id);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
```

---

## Inter-service Dependencies

| Service | Khi nào gọi | Nếu không có |
|---------|-------------|--------------|
| `menu-service` | AddOrderItem (validate + lấy price) | Error |
| `table-service` | CreateOrder khi `table_id` trống | `table_id` bắt buộc |
| `notification-service` | UpdateOrderStatus→Confirmed, UpdateOrderItemStatus→READY | Silently skip |

---

## Domain Errors

| Error | Mô tả |
|-------|-------|
| `ErrOrderNotFound` | Order ID không tồn tại |
| `ErrOrderNameRequired` | Thiếu tên khách |
| `ErrOrderPhoneRequired` | Thiếu SĐT |
| `ErrOrderTimeRequired` | Thiếu giờ đặt |
| `ErrOrderEndTimeInvalid` | `end_time` phải sau `time` |
| `ErrOrderPartySizeInvalid` | party_size ≤ 0 |
| `ErrOrderStatusInvalid` | Status không hợp lệ |
| `ErrOrderInvalidStatusTransition` | Vi phạm state machine |
| `ErrOrderAlreadyCancelled` | Order đã bị cancel |
| `ErrOrderCannotCancelCompleted` | Không thể cancel order completed |
| `ErrOrderItemNotFound` | Item không có trong order |
| `ErrOrderItemStatusInvalid` | Item status không hợp lệ |
| `ErrOrderItemInvalidStatusTransition` | Vi phạm item state machine |
| `ErrNoTableAvailable` | Không có bàn phù hợp cho khung giờ |
| `ErrTableRequired` | Không có table-service, `table_id` bắt buộc |

---

## Cấu hình (Environment Variables)

| Biến | Default | Mô tả |
|------|---------|-------|
| `SERVER_PORT` | `50055` | |
| `DB_HOST/PORT/USER/PASSWORD/NAME/SSLMODE` | — | PostgreSQL |
| `MENU_SERVICE_ADDR` | `localhost:50054` | Bắt buộc |
| `TABLE_SERVICE_ADDR` | `` | Trống = auto-assign disabled |
| `NOTIFICATION_SERVICE_ADDR` | `` | Trống = thông báo silently skipped |

---

## Vấn đề chưa giải quyết / Có thể triển khai thêm

- Operating hours (10:00–22:00) chỉ validate ở frontend, không có ở backend
- Không validate `end_time` không vượt quá giờ đóng cửa
- `item_id` trong `order_items` là PRIMARY KEY → 1 order chỉ có 1 dòng/món (không thể order cùng 1 món 2 dòng khác nhau)
- Khi `AllItemsServed()` = true, Manager phải manually `UpdateOrderStatus → Completed` — không có auto-complete
- Không có notification khi order bị cancel
- Không có history/audit trail của status changes
- Overlap detection trong auto-assign tính theo (`table_id`, `time`, `end_time`) — nếu `end_time` NULL thì overlap condition khác
- Không có concurrency control — race condition nếu 2 user tạo order cùng lúc và cùng grab 1 bàn
