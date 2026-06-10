# schedule-service

**Port:** 50052 (gRPC)  
**Module:** `restaurant-management/services/schedule-service`  
**Database:** PostgreSQL — `shifts` table

Quản lý ca làm việc cho nhân viên bếp. Thay thế `staff-service` (đã xóa).

Triết lý thiết kế: **shift tồn tại = đã xác nhận**. Không có status. Muốn hủy ca → xóa (`DeleteShift`).

---

## 5 RPCs

| RPC | Input | Output |
|-----|-------|--------|
| `CreateShift` | `user_id, date, start_time, end_time, role, notes, created_by` | `Shift` |
| `GetShift` | `shift_id` | `Shift` |
| `UpdateShift` | `shift_id, date?, start_time?, end_time?, notes?` | `Shift` |
| `DeleteShift` | `shift_id` | `{success, message}` |
| `ListShifts` | `month?, user_id?, role?, page, page_size` | `Shift[]` |

---

## Shift Entity

| Field | Type | Mô tả |
|-------|------|-------|
| `shift_id` | string (UUID) | DB-generated via `gen_random_uuid()` |
| `user_id` | string | FK logic tới user-service (không enforce ở DB) |
| `date` | string | Dạng `YYYY-MM-DD` |
| `start_time` | string | Dạng `HH:MM` |
| `end_time` | string | Dạng `HH:MM` — phải > `start_time` |
| `role` | string | `CHEF` / `WAITER` / `MANAGER` / `ADMIN` |
| `notes` | string | Ghi chú tùy chọn |
| `created_by` | string | `user_id` của người tạo ca |
| `created_at` | timestamp | |
| `updated_at` | timestamp | |

---

## Database Schema

```sql
CREATE TABLE IF NOT EXISTS shifts (
  shift_id   VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    VARCHAR(36) NOT NULL,
  date       DATE        NOT NULL,
  start_time VARCHAR(5)  NOT NULL,
  end_time   VARCHAR(5)  NOT NULL,
  role       VARCHAR(32) NOT NULL,
  notes      TEXT        NOT NULL DEFAULT '',
  created_by VARCHAR(36) NOT NULL DEFAULT '',
  created_at TIMESTAMP   NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP   NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_shifts_user_id ON shifts(user_id);
CREATE INDEX IF NOT EXISTS idx_shifts_date    ON shifts(date);
```

Schema tự tạo qua `ensureSchema` khi service khởi động.

---

## ListShifts Filter

Filter theo tháng: `month = "YYYY-MM"` → query `date >= YYYY-MM-01 AND date < YYYY-(MM+1)-01`

---

## Validation

- `user_id` required
- `date` phải đúng định dạng `YYYY-MM-DD` (10 ký tự)
- `start_time` và `end_time` phải đúng định dạng `HH:MM`
- `start_time < end_time`
- `role` phải thuộc `{CHEF, WAITER, MANAGER, ADMIN}`

---

## Folder Structure

```
services/schedule-service/
  cmd/server/main.go                      ← Entry point, wiring
  internal/
    domain/shift.go, errors.go            ← Entity, validation, errors
    repository/repository.go              ← Interface
    repository/shift_postgres.go          ← PostgreSQL implementation
    usecase/schedule_usecase.go           ← Business logic
    delivery/grpc/schedule_handler.go     ← gRPC handler
  pkg/config/config.go                    ← Env loader
  go.mod                                  ← module: restaurant-management/services/schedule-service
  Dockerfile
```

---

## Environment Variables

| Biến | Default |
|------|---------|
| `SERVER_PORT` | `50052` |
| `DB_HOST` | `localhost` |
| `DB_PORT` | `5432` |
| `DB_USER` | `restaurant_user` |
| `DB_PASSWORD` | `restaurant_pass` |
| `DB_NAME` | `restaurant_db` |
| `DB_SSLMODE` | `disable` |

---

## Build

```bash
cd services/schedule-service
go build -o server ./cmd/server/
```

---

## Vấn đề chưa giải quyết / Có thể triển khai thêm

- Không validate xung đột ca (cùng user, cùng ngày, cùng giờ) — có thể tạo nhiều shift chồng chéo
- `user_id` không được validate với user-service — shift có thể tham chiếu user không tồn tại
- Không có notification khi admin tạo ca cho nhân viên
- Chưa có recurring shifts (lịch lặp tuần/tháng)
