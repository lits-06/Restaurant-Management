# Table Service

**Port:** 50053 (gRPC)  
**Module:** `restaurant-management/services/table-service`  
**Proto:** `proto/table/table.proto`  
**Database:** PostgreSQL (`restaurant_db`)

## Tổng quan

Table Service quản lý **thông tin vật lý** của các bàn ăn — số bàn, sức chứa, trạng thái vật lý. Service này **không** quản lý việc bàn đang được đặt hay không — đó là trách nhiệm của order-service.

> **Nguyên tắc thiết kế:** Status của bàn (AVAILABLE/CLEANING/OUT_OF_SERVICE) chỉ phản ánh trạng thái vật lý do nhân viên cập nhật thủ công. Time-slot availability được tính từ bảng `orders` trong order-service.

## Kiến trúc

```
cmd/server/main.go
internal/
  domain/
    table.go      ← Table entity, TableStatus, validation, errors
  repository/
    table_postgres.go ← PostgreSQL implementation + ensureSchema
    repository.go     ← TableRepository interface
  usecase/
    table_usecase.go  ← Business logic
  delivery/grpc/
    table_handler.go  ← gRPC handler
pkg/config/config.go
```

---

## 3 Trạng thái Bàn

| Status | Enum | Mô tả |
|--------|------|-------|
| `AVAILABLE` | `STATUS_AVAILABLE` | Bàn sẵn sàng phục vụ |
| `CLEANING` | `STATUS_CLEANING` | Đang dọn dẹp sau khách trước |
| `OUT_OF_SERVICE` | `STATUS_OUT_OF_SERVICE` | Hỏng hoặc không dùng được |

---

## gRPC API — 7 RPCs

### 1. `CreateTable`

**Request:**
| Field | Type | Required | Validation |
|-------|------|----------|-----------|
| `table_number` | int32 | ✓ | > 0, unique |
| `capacity` | int32 | ✓ | 1–50 |

**Response:** `Table`, `success`, `message`  
**Errors:** `ALREADY_EXISTS` nếu `table_number` trùng; `INVALID_ARGUMENT` nếu capacity ngoài range

---

### 2. `GetTable`

**Request:** `table_id` (UUID)  
**Response:** `Table`, `success`, `message`

---

### 3. `UpdateTable`
Chỉ cập nhật `table_number` và `capacity`. Không thể đổi `table_id`.

**Request:** `table_id`, `table_number`, `capacity`  
**Response:** `Table`, `success`, `message`

---

### 4. `DeleteTable`

**Request:** `table_id`  
**Response:** `success`, `message`

---

### 5. `ListTables`

**Request:**
| Field | Type | Mô tả |
|-------|------|-------|
| `page` | int32 | Default: 1 |
| `page_size` | int32 | Default: 10 |
| `status` | TableStatus | Filter (optional) |

**Response:** `tables[]`, `total`, `page`, `page_size`, `success`, `message`

---

### 6. `UpdateTableStatus`
Cập nhật trạng thái vật lý của bàn. Chỉ nhận các status hợp lệ.

**Request:** `table_id`, `status` (TableStatus enum)  
**Response:** `Table`, `success`, `message`

---

### 7. `GetAvailableTables`
Được gọi bởi **order-service** trong luồng auto-assign bàn.

**Request:** `min_capacity` (int32) — party size của khách  
**SQL thực thi:**
```sql
SELECT * FROM restaurant_tables
WHERE status = 'AVAILABLE' AND capacity >= $min_capacity
ORDER BY capacity ASC, table_number ASC
```

**Response:** `tables[]`, `success`, `message`  
> Trả về tất cả bàn phù hợp, không lọc theo time-slot (order-service tự lọc).

---

## Database Schema

### Bảng `restaurant_tables`

```sql
CREATE TABLE IF NOT EXISTS restaurant_tables (
    table_id     VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    table_number INT         NOT NULL UNIQUE CHECK (table_number > 0),
    capacity     INT         NOT NULL CHECK (capacity > 0 AND capacity <= 50),
    status       VARCHAR(32) NOT NULL DEFAULT 'AVAILABLE',
    created_at   TIMESTAMP   NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP   NOT NULL DEFAULT NOW()
);
```

### Indexes

```sql
CREATE INDEX IF NOT EXISTS idx_tables_status        ON restaurant_tables(status);
CREATE INDEX IF NOT EXISTS idx_tables_table_number  ON restaurant_tables(table_number);
```

---

## Domain Errors

| Error | gRPC Code | Mô tả |
|-------|-----------|-------|
| `ErrTableNotFound` | NOT_FOUND | |
| `ErrTableNumberAlreadyExists` | ALREADY_EXISTS | `table_number` trùng |
| `ErrTableNotAvailable` | FAILED_PRECONDITION | Bàn không ở trạng thái AVAILABLE khi cần |
| `ErrInvalidTableNumber` | INVALID_ARGUMENT | ≤ 0 |
| `ErrInvalidCapacity` | INVALID_ARGUMENT | ≤ 0 |
| `ErrCapacityTooLarge` | INVALID_ARGUMENT | > 50 |
| `ErrInvalidStatus` | INVALID_ARGUMENT | Status không hợp lệ |
| `ErrInvalidTableID` | INVALID_ARGUMENT | UUID rỗng |

---

## Cấu hình (Environment Variables)

| Biến | Default |
|------|---------|
| `SERVER_PORT` | `50053` |
| `DB_HOST` | `localhost` |
| `DB_PORT` | `5432` |
| `DB_USER` | `restaurant_user` |
| `DB_PASSWORD` | `restaurant_pass` |
| `DB_NAME` | `restaurant_db` |
| `DB_SSLMODE` | `disable` |

---

## Tích hợp với các Service Khác

- **order-service** gọi `GetAvailableTables` khi `table_id` không được truyền trong `CreateOrder`
- **api-gateway** expose các HTTP endpoints `/tables/...` nhưng **không** cập nhật status tự động theo booking

---

## Vấn đề chưa giải quyết / Có thể triển khai thêm

- Status bàn phải cập nhật thủ công bởi nhân viên — chưa có auto-flow (e.g., bàn tự chuyển CLEANING khi order Completed)
- `table_id` trong notifications/kitchen app hiển thị dạng UUID cắt ngắn, chưa resolve sang `table_number`
- Không có lịch sử thay đổi status
- Capacity constraint (1–50) hard-coded trong schema — nên move vào config
- Không thể filter `GetAvailableTables` theo `table_number` range (e.g., chỉ lấy bàn VIP 10–20)
