# Auth Service

**Port:** 50051 (gRPC)  
**Module:** `restaurant-management/services/auth-service`  
**Proto:** `proto/auth/auth.proto`

## Tổng quan

Auth Service xử lý toàn bộ luồng xác thực (authentication). Service này **không có database riêng** — mọi thông tin người dùng đều được ủy thác (delegate) sang user-service qua gRPC. Chỉ dùng **Redis** để lưu refresh token.

## Kiến trúc

```
cmd/server/main.go          ← Entry point, wiring
internal/
  domain/                   ← JWT config, token structs
  usecase/auth_usecase.go   ← Business logic (login, register, refresh...)
  delivery/grpc/            ← gRPC handler
  grpcclient/user_client.go ← gRPC client gọi user-service
pkg/config/config.go        ← Env vars loader
```

## gRPC API

### 1. `Register`
Tạo tài khoản mới. Delegate sang `user-service.CreateUser`.

**Request:**
| Field | Type | Required | Mô tả |
|-------|------|----------|-------|
| `email` | string | ✓ | Email duy nhất |
| `password` | string | ✓ | Plain text, bcrypt bởi user-service |
| `username` | string | ✓ | Tối thiểu 3 ký tự, duy nhất |
| `full_name` | string | ✓ | Họ và tên |
| `phone` | string | — | Số điện thoại |

**Response:** `user_id`, `message`, `success`

---

### 2. `Login`
Xác thực email/password, trả về JWT tokens.

**Flow:**
1. Gọi `user-service.VerifyCredentials(email, password)`
2. Nếu thành công → `GenerateAccessToken(user_id, email, roles)`
3. Lưu refresh token vào Redis: key `refresh_token:<token>`, value `{user_id, email}`
4. Trả về `access_token` + `refresh_token`

**Request:** `email`, `password`  
**Response:** `access_token`, `refresh_token`, `user_id`, `success`, `message`

---

### 3. `RefreshToken`
Cấp lại access token từ refresh token còn hợp lệ.

**Flow:**
1. Validate JWT của refresh token
2. Kiểm tra refresh token tồn tại trong Redis
3. Gọi `user-service.GetUser(user_id)` để lấy roles mới nhất
4. Tạo access token mới

**Request:** `refresh_token`  
**Response:** `access_token`, `success`, `message`

---

### 4. `VerifyToken`
Kiểm tra access token còn hợp lệ hay không. Dùng bởi api-gateway trên mỗi request cần auth.

**Request:** `access_token`  
**Response:** `valid`, `user_id`, `email`, `expires_at`, `roles[]`

---

### 5. `Logout`
Thu hồi refresh token khỏi Redis.

**Request:** `user_id`, `access_token`, `refresh_token`  
> **Quan trọng:** `refresh_token` bắt buộc — đây là field thực sự bị xóa khỏi Redis. `access_token` không được lưu nên không thể thu hồi.

**Response:** `success`, `message`

---

### 6. `ChangePassword`
Thay đổi mật khẩu. Delegate sang `user-service.ChangePassword`.

**Request:** `user_id`, `old_password`, `new_password`  
**Response:** `success`, `message`

---

## JWT

| Thuộc tính | Giá trị |
|-----------|---------|
| Algorithm | HS256 |
| Secret | `JWT_SECRET` env var (log warning nếu dùng default) |
| Access token TTL | `JWT_ACCESS_MINUTES` (default: 15 phút) |
| Refresh token TTL | `JWT_REFRESH_HOURS` (default: 168 giờ = 7 ngày) |
| Claims | `user_id`, `email`, `roles []string`, `exp` |

---

## Database

Không có PostgreSQL. Chỉ dùng Redis.

### Redis Keys
| Key pattern | Type | Value | TTL |
|------------|------|-------|-----|
| `refresh_token:<token_string>` | Hash | `{user_id, email}` | JWT_REFRESH_HOURS |

---

## Cấu hình (Environment Variables)

| Biến | Default | Mô tả |
|------|---------|-------|
| `SERVER_PORT` | `50051` | gRPC listen port |
| `USER_SERVICE_ADDR` | `localhost:50056` | Địa chỉ user-service |
| `REDIS_HOST` | `localhost` | Redis host |
| `REDIS_PORT` | `6379` | Redis port |
| `REDIS_PASSWORD` | `` | Redis password |
| `REDIS_DB` | `0` | Redis database number |
| `JWT_SECRET` | `default-secret` | **Thay trong production** |
| `JWT_ACCESS_MINUTES` | `15` | Access token TTL (phút) |
| `JWT_REFRESH_HOURS` | `168` | Refresh token TTL (giờ) |

---

## Roles được xử lý

Auth service không kiểm soát roles — chỉ đọc từ user-service và đưa vào JWT claims.

5 roles: `USER`, `MANAGER`, `CHEF`, `WAITER`, `ADMIN`

---

## Vấn đề chưa giải quyết / Có thể triển khai thêm

- Access token không thể thu hồi (blacklist) — logout chỉ vô hiệu refresh token
- Không có rate limiting trên Login (brute-force protection)
- Không có multi-device session management (chỉ 1 refresh token/user)
- Không có email verification sau Register
- Không có forgot password / reset password flow
- Redis key TTL không đồng bộ với JWT exp nếu `JWT_REFRESH_HOURS` thay đổi sau khi token đã tạo
