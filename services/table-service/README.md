# Table Service

Microservice quản lý thông tin vật lý của các bàn ăn trong nhà hàng.

**Port:** `50053` | **Protocol:** gRPC | **Status:** Active (enabled in docker-compose)

---

## Trách nhiệm

Table-service chỉ quản lý **thông tin tĩnh và trạng thái vật lý** của bàn:
- Bàn nào tồn tại, sức chứa bao nhiêu
- Trạng thái vật lý hiện tại: sẵn sàng / đang dọn / hỏng

**Không thuộc trách nhiệm của service này:**
- Bàn nào đang có khách trong khung giờ nào → **order-service** quản lý
- Đặt bàn trước theo thời gian → **order-service** quản lý (auto-assign)

---

## Database

Service dùng **1 bảng** trong `restaurant_db`, tự tạo khi khởi động:

```sql
CREATE TABLE restaurant_tables (
    table_id     VARCHAR(36) PRIMARY KEY,
    table_number INTEGER     NOT NULL UNIQUE,  -- số bàn (1, 2, 3, ...)
    capacity     INTEGER     NOT NULL,          -- số chỗ ngồi tối đa (1–50)
    status       VARCHAR(32) NOT NULL DEFAULT 'AVAILABLE',
    created_at   TIMESTAMP   NOT NULL,
    updated_at   TIMESTAMP   NOT NULL,
    CONSTRAINT chk_table_number CHECK (table_number > 0),
    CONSTRAINT chk_capacity     CHECK (capacity > 0 AND capacity <= 50)
);

CREATE INDEX idx_restaurant_tables_status   ON restaurant_tables(status);
CREATE INDEX idx_restaurant_tables_capacity ON restaurant_tables(capacity);
```

---

## Trạng thái bàn (TableStatus)

| Status | Ý nghĩa | Ai set |
|--------|---------|--------|
| `AVAILABLE` | Bàn sạch, sẵn sàng | Staff (sau khi dọn xong) |
| `CLEANING` | Đang dọn dẹp sau khi khách rời | Staff |
| `OUT_OF_SERVICE` | Hỏng / bảo trì | Admin |

**Quan trọng:** Status này **không liên quan** đến việc bàn có được đặt hay chưa. Một bàn `AVAILABLE` vẫn có thể đã được book trong order-service. Xem [auto-assign flow](#luồng-auto-assign) để hiểu rõ hơn.

### State machine

```
AVAILABLE    → CLEANING       (khách vừa rời, nhân viên bắt đầu dọn)
CLEANING     → AVAILABLE      (dọn xong)
AVAILABLE    → OUT_OF_SERVICE
CLEANING     → OUT_OF_SERVICE
OUT_OF_SERVICE → AVAILABLE    (sửa xong)
```

---

## gRPC API — 7 RPCs

### `CreateTable`
```
Request:  { table_number: int32, capacity: int32 }
Response: { table: Table, success: bool, message: string }
```
Kiểm tra `table_number` chưa tồn tại trước khi tạo.

### `GetTable`
```
Request:  { table_id: string }
Response: { table: Table, success: bool, message: string }
```

### `UpdateTable`
```
Request:  { table_id: string, table_number: int32, capacity: int32 }
Response: { table: Table, success: bool, message: string }
```
Partial update — chỉ cập nhật field nào được truyền vào (> 0).

### `DeleteTable`
```
Request:  { table_id: string }
Response: { success: bool, message: string }
```

### `ListTables`
```
Request:  { page: int32, page_size: int32, status: TableStatus }
Response: { tables: []Table, total: int32, page: int32, page_size: int32, ... }
```
Filter theo status (optional). Sắp xếp theo `table_number ASC`.

### `UpdateTableStatus`
```
Request:  { table_id: string, status: TableStatus }
Response: { table: Table, success: bool, message: string }
```
Áp dụng state machine — từ chối transition không hợp lệ (ví dụ: CLEANING → CLEANING).

### `GetAvailableTables`
```
Request:  { min_capacity: int32 }
Response: { tables: []Table, success: bool, message: string }
```
Trả về bàn có `status = AVAILABLE` và `capacity >= min_capacity`.
Sắp xếp theo `capacity ASC, table_number ASC` (best-fit first).
**Dùng bởi order-service** khi auto-assign bàn cho order mới.

---

## Luồng auto-assign

Khi khách tạo order mà không chỉ định bàn, order-service tự tìm bàn phù hợp:

```
POST /orders { party_size: 4, time: "19:00", end_time: "21:00" }

order-service:
  1. Gọi GetAvailableTables(min_capacity=4)
     → table-service trả về [Bàn 3 (4 chỗ), Bàn 7 (6 chỗ), Bàn 9 (8 chỗ)]

  2. Query DB nội bộ: bàn nào có order trùng [19:00, 21:00)?
     → [Bàn 3, Bàn 7] đã bị chiếm

  3. Gán Bàn 9 → tạo order với table_id = "uuid-ban-9"
```

---

## HTTP API (qua api-gateway :8080)

| Method | Path | Mô tả |
|--------|------|-------|
| `GET` | `/tables` | Danh sách bàn (`?status=AVAILABLE&page=1&page_size=20`) |
| `POST` | `/tables` | Tạo bàn mới |
| `GET` | `/tables/available` | Bàn đang AVAILABLE (`?min_capacity=4`) |
| `GET` | `/tables/{id}` | Chi tiết một bàn |
| `PUT` | `/tables/{id}` | Cập nhật thông tin bàn |
| `DELETE` | `/tables/{id}` | Xoá bàn |
| `PATCH` | `/tables/{id}/status` | Đổi trạng thái bàn |

---

## Cấu trúc thư mục

```
services/table-service/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   │   ├── table.go          # Entity + state machine
│   │   └── errors.go         # Domain errors
│   ├── repository/
│   │   ├── repository.go     # TableRepository interface
│   │   └── table_postgres.go # PostgreSQL implementation
│   ├── usecase/
│   │   └── table_usecase.go  # Business logic
│   └── delivery/grpc/
│       └── table_handler.go  # gRPC handler
└── pkg/config/config.go
```

---

## Configuration

| Env var | Default | Mô tả |
|---------|---------|-------|
| `SERVER_PORT` | `50053` | gRPC port |
| `DATABASE_HOST` | `localhost` | PostgreSQL host |
| `DATABASE_PORT` | `5432` | |
| `DATABASE_USER` | `restaurant_user` | |
| `DATABASE_PASSWORD` | `restaurant_pass` | |
| `DATABASE_NAME` | `restaurant_db` | |

---

## Build & Run

```bash
cd services/table-service
go build ./cmd/server/
go run cmd/server/main.go
```

## grpcurl — ví dụ

```bash
# Tạo bàn số 5, sức chứa 4 người
grpcurl -plaintext -d '{"table_number":5,"capacity":4}' \
  localhost:50053 table.TableService/CreateTable

# Danh sách tất cả bàn
grpcurl -plaintext -d '{"page":1,"page_size":20}' \
  localhost:50053 table.TableService/ListTables

# Bàn trống, ít nhất 4 chỗ
grpcurl -plaintext -d '{"min_capacity":4}' \
  localhost:50053 table.TableService/GetAvailableTables

# Đổi trạng thái bàn sang CLEANING
grpcurl -plaintext -d '{"table_id":"TABLE_UUID","status":"STATUS_CLEANING"}' \
  localhost:50053 table.TableService/UpdateTableStatus
```

---

## Related services

- **order-service** (port 50055) — gọi `GetAvailableTables` khi auto-assign bàn cho order mới
