# User Service

**Port:** 50056 (gRPC)  
**Module:** `restaurant-management/services/user-service`  
**Proto:** `proto/user/user.proto`  
**Database:** PostgreSQL (`restaurant_db`)

## Tổng quan

User Service quản lý danh tính và phân quyền cho toàn bộ người dùng hệ thống. Đây là nguồn dữ liệu chính (source of truth) về user — auth-service không có DB riêng mà delegate sang đây.

## Kiến trúc

```
cmd/server/main.go
internal/
  domain/
    user.go          ← User entity, UserRole, UserStatus, validation
    errors.go        ← Typed domain errors
  repository/
    repository.go    ← UserRepository interface
    user_postgres.go ← PostgreSQL implementation + ensureSchema
  usecase/
    user_usecase.go  ← Business logic
  delivery/grpc/
    user_handler.go  ← gRPC handler (maps proto ↔ domain)
pkg/config/config.go ← Env vars loader
```

## 5 Roles

| Role | Enum | Mô tả |
|------|------|-------|
| `USER` | `ROLE_USER` | Khách hàng — đặt bàn, xem order của bản thân |
| `MANAGER` | `ROLE_MANAGER` | Tạo order walk-in, cập nhật trạng thái order |
| `CHEF` | `ROLE_CHEF` | Xem món trong order, mark COOKING/READY |
| `WAITER` | `ROLE_WAITER` | Nhận thông báo từ bếp, mark SERVED |
| `ADMIN` | `ROLE_ADMIN` | Toàn quyền — dashboard, menu, users, orders, tables |

User có thể có **nhiều roles đồng thời** (e.g., ADMIN+CHEF).  
Default khi tạo user: `[USER]` nếu không truyền roles.

## 3 Trạng thái User

| Status | Mô tả |
|--------|-------|
| `ACTIVE` | Tài khoản hoạt động bình thường |
| `INACTIVE` | Tài khoản bị vô hiệu hóa |
| `SUSPENDED` | Tài khoản bị đình chỉ (block login ở VerifyCredentials) |

---

## gRPC API — 10 RPCs

### 1. `CreateUser`

**Request:**
| Field | Type | Required | Mô tả |
|-------|------|----------|-------|
| `email` | string | ✓ | Unique |
| `username` | string | ✓ | Unique, tối thiểu 3 ký tự |
| `full_name` | string | ✓ | |
| `phone` | string | — | |
| `password` | string | ✓ | Plain text, được bcrypt trong repository |
| `roles` | UserRole[] | — | Default `[ROLE_USER]` |

**Response:** `User`, `success`, `message`

---

### 2. `GetUser`

**Request:** `user_id` (UUID)  
**Response:** `User`, `success`, `message`

---

### 3. `GetUserByEmail`

**Request:** `email`  
**Response:** `User`, `success`, `message` (reuses `GetUserResponse`)

---

### 4. `UpdateUser`
Cập nhật profile. **Không cập nhật password hoặc roles** — dùng `ChangePassword` / `AssignRole`.

**Request:** `user_id`, `email`, `username`, `full_name`, `phone`, `status`  
**Response:** `User`, `success`, `message`

---

### 5. `DeleteUser`
Xóa vĩnh viễn. Cascade xóa `user_roles`.

**Request:** `user_id`  
**Response:** `success`, `message`

---

### 6. `ListUsers`

**Request:**
| Field | Type | Mô tả |
|-------|------|-------|
| `page` | int32 | Default: 1 |
| `page_size` | int32 | Default: 10 |
| `status` | UserStatus | Filter |
| `role` | UserRole | Filter |
| `keyword` | string | ILIKE trên email, username, full_name |

**Response:** `users[]`, `total`, `page`, `page_size`, `success`, `message`

---

### 7. `AssignRole`
**Thay thế toàn bộ roles** của user (không append). Dùng transaction.

**Request:** `user_id`, `roles[]`  
**Response:** `success`, `message`

---

### 8. `GetUserRoles`

**Request:** `user_id`  
**Response:** `roles[]`, `success`, `message`

---

### 9. `ChangePassword`
Kiểm tra `old_password` trước khi cập nhật.

**Request:** `user_id`, `old_password`, `new_password`  
**Response:** `success`, `message`

---

### 10. `VerifyCredentials`
Xác thực email + password. Được gọi bởi **auth-service** trong luồng Login.

**Request:** `email`, `password` (plain text)  
**Flow:** Tìm user theo email → bcrypt.Compare(password, hash) → kiểm tra status không SUSPENDED  
**Response:** `success`, `message`, `user_id`, `email`, `roles[]`

---

## Database Schema

### Bảng `users`

```sql
CREATE TABLE IF NOT EXISTS users (
    user_id    VARCHAR(36)  PRIMARY KEY,
    email      VARCHAR(255) NOT NULL UNIQUE,
    username   VARCHAR(100) NOT NULL UNIQUE,
    full_name  VARCHAR(255) NOT NULL,
    phone      VARCHAR(50)  NOT NULL DEFAULT '',
    password   VARCHAR(255) NOT NULL,   -- bcrypt hash
    status     VARCHAR(32)  NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NOT NULL
);
```

### Bảng `user_roles`

```sql
CREATE TABLE IF NOT EXISTS user_roles (
    user_id VARCHAR(36) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    role    VARCHAR(32) NOT NULL,
    PRIMARY KEY (user_id, role)
);
```

### Indexes

```sql
CREATE INDEX IF NOT EXISTS idx_users_status    ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_username  ON users(username);
CREATE INDEX IF NOT EXISTS idx_user_roles_role ON user_roles(role);
```

> **Migration note:** `ensureSchema` chạy `ALTER TABLE users DROP COLUMN IF EXISTS roles` để xóa cột roles cũ (thiết kế trước lưu comma-separated string).

---

## Domain Errors

| Error | gRPC Code | Mô tả |
|-------|-----------|-------|
| `ErrUserNotFound` | NOT_FOUND | User ID không tồn tại |
| `ErrEmailAlreadyExists` | ALREADY_EXISTS | Email đã dùng |
| `ErrUsernameAlreadyExists` | ALREADY_EXISTS | Username đã dùng |
| `ErrInvalidCredentials` | UNAUTHENTICATED | Sai email/password |
| `ErrAccountSuspended` | PERMISSION_DENIED | Tài khoản bị đình chỉ |

---

## Cấu hình (Environment Variables)

| Biến | Default | Mô tả |
|------|---------|-------|
| `SERVER_PORT` | `50056` | gRPC listen port |
| `DB_HOST` | `localhost` | |
| `DB_PORT` | `5432` | |
| `DB_USER` | `restaurant_user` | |
| `DB_PASSWORD` | `restaurant_pass` | |
| `DB_NAME` | `restaurant_db` | |
| `DB_SSLMODE` | `disable` | |

---

## Vấn đề chưa giải quyết / Có thể triển khai thêm

- `AssignRole` thay thế toàn bộ — muốn thêm/xóa 1 role cần call `GetUserRoles` trước
- Không có audit log khi role thay đổi
- Không validate phone format
- Xóa user là hard delete — orders với `user_id` vẫn còn (orphaned reference)
- Không có soft delete / deactivation-only mode
- Password strength không được kiểm tra
- Offset pagination có thể chậm với data lớn — nên dùng cursor pagination
