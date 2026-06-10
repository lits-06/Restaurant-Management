# Staff Service

**Port:** 50052 (gRPC)  
**Module:** `restaurant-management/services/staff-service`  
**Proto:** `proto/staff/staff.proto`  
**Database:** PostgreSQL (`restaurant_db`)

## Tổng quan

Staff Service quản lý hồ sơ nhân viên trong nhà hàng — thông tin hiển thị (tên, vai trò, liên lạc, avatar). Service này **độc lập** với user-service: staff record là profile HR, user record là identity auth. Hai hệ thống chưa được liên kết.

## Kiến trúc

```
cmd/server/main.go
internal/
  domain/
    staff.go      ← Staff entity, validation
    errors.go     ← Typed domain errors
  repository/
    staff_postgres.go ← PostgreSQL implementation + ensureSchema
    repository.go     ← StaffRepository interface
  usecase/
    staff_usecase.go  ← Business logic
  delivery/grpc/
    staff_handler.go  ← gRPC handler
pkg/config/config.go
```

---

## gRPC API — 5 RPCs

### 1. `CreateStaff`

**Request:**
| Field | Type | Required | Mô tả |
|-------|------|----------|-------|
| `name` | string | ✓ | Tên nhân viên |
| `role` | string | ✓ | Tên vai trò (free text, không validate enum) |
| `contact` | string | — | SĐT hoặc email liên lạc |
| `avatar` | string | — | URL ảnh đại diện |

**Response:** `Staff`, `success`, `message`

---

### 2. `GetStaff`

**Request:** `staff_id` (UUID)  
**Response:** `Staff`, `success`, `message`

---

### 3. `UpdateStaff`
Cập nhật toàn bộ các field (tên, role, liên lạc, avatar).

**Request:** `staff_id`, `name`, `role`, `contact`, `avatar`  
**Response:** `Staff`, `success`, `message`

---

### 4. `DeleteStaff`

**Request:** `staff_id`  
**Response:** `success`, `message`

---

### 5. `ListStaff`

**Request:**
| Field | Type | Mô tả |
|-------|------|-------|
| `page` | int32 | Default: 1 |
| `page_size` | int32 | Default: 10 |
| `keyword` | string | Tìm trong name (ILIKE) |

**Response:** `staff[]`, `total`, `page`, `page_size`, `success`, `message`

---

## Database Schema

### Bảng `staff`

```sql
CREATE TABLE IF NOT EXISTS staff (
    staff_id   VARCHAR(36)  PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    role       VARCHAR(100) NOT NULL DEFAULT '',
    contact    VARCHAR(255) NOT NULL DEFAULT '',
    avatar     TEXT         NOT NULL DEFAULT '',
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);
```

### Indexes

```sql
CREATE INDEX IF NOT EXISTS idx_staff_name ON staff(name);
```

---

## Cấu hình (Environment Variables)

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

## Vấn đề chưa giải quyết / Có thể triển khai thêm

- `role` field là free text — không validate enum, không liên kết với `UserRole` trong user-service
- Staff record và User record là 2 bảng riêng biệt, chưa có foreign key hay link
- Không có schedule/shift management (WeeklyScheduler UI đã có nhưng chưa có backend)
- Avatar là URL text — chưa có file upload / CDN integration
- Không có soft delete
- Không có pagination theo role
