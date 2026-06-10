# API Gateway

**Port:** 8080 (HTTP/REST)  
**Module:** `restaurant-management` (root module)

## Tổng quan

API Gateway là điểm vào duy nhất (single entry point) của toàn bộ hệ thống. Nhận HTTP/REST requests từ 3 frontend apps, dịch sang gRPC calls đến các backend services, và trả về JSON response. Không có business logic — chỉ translation + authentication middleware.

## Kiến trúc

```
cmd/server/main.go          ← Entry point: wiring, mux setup, CORS, graceful shutdown
internal/
  grpcclient/               ← gRPC clients cho từng service
    auth_client.go
    menu_client.go
    order_client.go
    staff_client.go
    table_client.go
    user_client.go
    notification_client.go
  handler/                  ← HTTP handlers
    auth_handler.go
    menu_handler.go
    order_handler.go
    staff_handler.go
    table_handler.go
    user_handler.go
    notification_handler.go ← WebSocket upgrade + gRPC streaming bridge
```

---

## CORS

Allowed origins: `http://localhost:5173`, `http://localhost:5174`, `http://localhost:5175`

```
Access-Control-Allow-Headers:  Content-Type, Authorization
Access-Control-Allow-Methods:  GET, POST, PUT, PATCH, DELETE, OPTIONS
Access-Control-Allow-Credentials: true
```

WebSocket endpoint `/ws/notifications` bypass CORS middleware — có origin check riêng trong `wsUpgrader.CheckOrigin`.

---

## Authentication

Hầu hết endpoints không require auth ở gateway level — auth được check ở logic handler.

Pattern chung:
1. Extract `Authorization: Bearer <token>` header
2. Gọi `auth-service.VerifyToken(token)` → nhận `user_id`, `email`, `roles[]`
3. Dùng `roles[]` để authorize action

**Helper functions trong `order_handler.go`:**
- `verifyCaller(r)` — optional: extract token nếu có, trả về `callerInfo{UserID, Roles}` hoặc nil
- `checkOrderAccess(r, orderUserID)` — combined auth+authz: nếu order có user_id, caller phải là owner hoặc staff
- `checkUserIDAccess(r, targetUserID)` — cho filter `user_id` trong ListOrders
- `canMarkItemStatus(roles, targetStatus)` — kiểm tra role cho phép mark item status không

---

## HTTP → gRPC Endpoint Map

### Auth (`/auth/*`)

| Method | Path | Handler | gRPC |
|--------|------|---------|------|
| POST | `/auth/register` | `Register` | `auth.Register` |
| POST | `/auth/login` | `Login` | `auth.Login` |
| POST | `/auth/refresh` | `RefreshToken` | `auth.RefreshToken` |
| POST | `/auth/verify` | `VerifyToken` | `auth.VerifyToken` |
| POST | `/auth/logout` | `Logout` | `auth.Logout` |
| POST | `/auth/change-password` | `ChangePassword` | `auth.ChangePassword` |

**Login request body:**
```json
{ "email": "...", "password": "..." }
```
**Login response:**
```json
{ "access_token": "...", "refresh_token": "...", "user_id": "...", "success": true }
```

**Logout request body:** `{ "refresh_token": "..." }` — bắt buộc có refresh_token

---

### Menu (`/menu/*`)

| Method | Path | Handler | Auth | gRPC |
|--------|------|---------|------|------|
| GET | `/menu/items` | `ListMenuItems` | Không | `menu.ListMenuItems` |
| POST | `/menu/items` | `CreateMenuItem` | Không* | `menu.CreateMenuItem` |
| GET | `/menu/items/{id}` | `GetMenuItem` | Không | `menu.GetMenuItem` |
| PUT | `/menu/items/{id}` | `UpdateMenuItem` | Không* | `menu.UpdateMenuItem` |
| DELETE | `/menu/items/{id}` | `DeleteMenuItem` | Không* | `menu.DeleteMenuItem` |
| GET | `/menu/categories` | `ListCategories` | Không | `menu.ListCategories` |
| POST | `/menu/categories` | `CreateCategory` | Không* | `menu.CreateCategory` |
| GET | `/menu/categories/{id}` | `GetCategory` | Không | `menu.GetCategory` |
| PUT | `/menu/categories/{id}` | `UpdateCategory` | Không* | `menu.UpdateCategory` |
| DELETE | `/menu/categories/{id}` | `DeleteCategory` | Không* | `menu.DeleteCategory` |

> *Chưa có role-based middleware — tất cả request đều pass qua.

**Query params `GET /menu/items`:** `page`, `page_size`, `category_id`, `keyword`

---

### Orders (`/orders/*`)

| Method | Path | Auth | Role | gRPC |
|--------|------|------|------|------|
| GET | `/orders` | Optional | Xem tất cả: không; filter `user_id`: phải là owner hoặc staff | `order.ListOrders` |
| POST | `/orders` | Optional | Không bắt buộc (walk-in) | `order.CreateOrder` |
| GET | `/orders/{id}` | Conditional | Nếu order có user_id: phải là owner hoặc staff | `order.GetOrder` |
| PUT | `/orders/{id}` | Conditional | Như trên | `order.UpdateOrder` |
| DELETE | `/orders/{id}` | Conditional | Như trên | `order.DeleteOrder` |
| POST | `/orders/{id}/cancel` | Conditional | Như trên | `order.CancelOrder` |
| PATCH | `/orders/{id}/status` | Required | ADMIN/MANAGER/CHEF/WAITER | `order.UpdateOrderStatus` |
| POST | `/orders/{id}/items` | Conditional | Như trên | `order.AddOrderItem` |
| DELETE | `/orders/{id}/items/{itemId}` | Conditional | Như trên | `order.RemoveOrderItem` |
| PATCH | `/orders/{id}/items/{itemId}/status` | Required | COOKING/READY: CHEF+ADMIN+MANAGER; SERVED: WAITER+ADMIN+MANAGER | `order.UpdateOrderItemStatus` |

**POST `/orders` request body:**
```json
{
  "name": "Nguyen Van A",
  "phone": "0901234567",
  "date": "2026-06-15",
  "time": "19:00",
  "end_time": "21:00",
  "party_size": 4,
  "notes": "dị ứng hải sản",
  "items": [{ "item_id": "...", "quantity": 2 }]
}
```
`table_id` optional — bỏ trống để auto-assign.  
`user_id` không truyền từ client — api-gateway extract từ JWT nếu có token.

**Query params `GET /orders`:** `page`, `page_size`, `status`, `keyword`, `user_id`

**PATCH `/orders/{id}/items/{itemId}/status` body:** `{ "item_status": "COOKING" }`

---

### Tables (`/tables/*`)

| Method | Path | Auth | gRPC |
|--------|------|------|------|
| GET | `/tables` | Không | `table.ListTables` |
| POST | `/tables` | Không* | `table.CreateTable` |
| GET | `/tables/available` | Không | `table.GetAvailableTables` |
| GET | `/tables/{id}` | Không | `table.GetTable` |
| PUT | `/tables/{id}` | Không* | `table.UpdateTable` |
| DELETE | `/tables/{id}` | Không* | `table.DeleteTable` |
| PATCH | `/tables/{id}/status` | Không* | `table.UpdateTableStatus` |

**Query params `GET /tables`:** `page`, `page_size`, `status`  
**Query params `GET /tables/available`:** `min_capacity`

**POST `/tables` body:** `{ "table_number": 5, "capacity": 4 }`  
**PATCH `/tables/{id}/status` body:** `{ "status": "CLEANING" }`

---

### Schedule (`/schedule/*`)

| Method | Path | Auth | Role | gRPC |
|--------|------|------|------|------|
| GET | `/schedule/shifts` | Required | Any staff | `schedule.ListShifts` |
| POST | `/schedule/shifts` | Required | Self: any staff; Others: ADMIN/MANAGER | `schedule.CreateShift` |
| GET | `/schedule/shifts/{id}` | Required | Owner hoặc ADMIN/MANAGER | `schedule.GetShift` |
| PUT | `/schedule/shifts/{id}` | Required | Owner hoặc ADMIN/MANAGER | `schedule.UpdateShift` |
| DELETE | `/schedule/shifts/{id}` | Required | Owner hoặc ADMIN/MANAGER | `schedule.DeleteShift` |

**Query params `GET /schedule/shifts`:** `month` (YYYY-MM), `user_id`, `role`, `page`, `page_size`

**POST `/schedule/shifts` body:**
```json
{
  "user_id": "...",       // optional — defaults to caller; if different from caller → requires ADMIN/MANAGER
  "date": "2026-06-20",
  "start_time": "08:00",
  "end_time": "16:00",
  "role": "CHEF",
  "notes": "Ca sáng"
}
```

---

### Users (`/users/*`)

| Method | Path | Auth | gRPC |
|--------|------|------|------|
| GET | `/users/by-email?email=` | Không* | `user.GetUserByEmail` |
| GET | `/users` | Không* | `user.ListUsers` |
| POST | `/users` | Không* | `user.CreateUser` |
| GET | `/users/{id}` | Không* | `user.GetUser` |
| PUT | `/users/{id}` | Không* | `user.UpdateUser` |
| DELETE | `/users/{id}` | Không* | `user.DeleteUser` |
| GET | `/users/{id}/roles` | Không* | `user.GetUserRoles` |
| PATCH | `/users/{id}/roles` | Không* | `user.AssignRole` |
| PATCH | `/users/{id}/password` | Không* | `user.ChangePassword` |

---

### WebSocket (`/ws/*`)

| Path | Auth | Mô tả |
|------|------|-------|
| `GET /ws/notifications?token=<jwt>&role=<CHEF\|WAITER>` | JWT query param | Real-time notifications |

**Luồng WebSocket:**
1. Upgrade HTTP → WebSocket
2. Validate JWT từ `?token=` param
3. Mở `notification-service.Subscribe(role)` gRPC streaming
4. Forward mỗi `Notification` message → JSON → WebSocket
5. Drain goroutine: `conn.ReadMessage()` → nếu lỗi (client close) → cancel ctx → đóng gRPC stream

---

### Health

| Method | Path | Mô tả |
|--------|------|-------|
| GET | `/health` | `{"status":"ok"}` |

---

## Cấu hình (Environment Variables)

| Biến | Default |
|------|---------|
| `SERVER_PORT` | `8080` |
| `AUTH_SERVICE_HOST` | `localhost` |
| `AUTH_SERVICE_PORT` | `50051` |
| `MENU_SERVICE_HOST` | `localhost` |
| `MENU_SERVICE_PORT` | `50054` |
| `ORDER_SERVICE_HOST` | `localhost` |
| `ORDER_SERVICE_PORT` | `50055` |
| `SCHEDULE_SERVICE_HOST` | `localhost` |
| `SCHEDULE_SERVICE_PORT` | `50052` |
| `TABLE_SERVICE_HOST` | `localhost` |
| `TABLE_SERVICE_PORT` | `50053` |
| `USER_SERVICE_HOST` | `localhost` |
| `USER_SERVICE_PORT` | `50056` |
| `NOTIFICATION_SERVICE_HOST` | `localhost` |
| `NOTIFICATION_SERVICE_PORT` | `50058` |
| `ENVIRONMENT` | `development` |

---

## Build

```bash
# Từ repo root (api-gateway dùng root go.mod)
cd api-gateway
go build -o server ./cmd/server/
```

---

## Vấn đề chưa giải quyết / Có thể triển khai thêm

- Không có role-based middleware cho menu/staff/table/user routes — hiện tại mọi request đều pass
- Không có rate limiting
- Không có request logging (chỉ log lúc startup)
- Không có request timeout per-endpoint (chỉ có `ReadHeaderTimeout: 10s`)
- WebSocket không có auto-reconnect (client-side)
- Không có circuit breaker cho gRPC calls — nếu 1 service down, request bị block cho đến timeout
- REPORT service chưa có HTTP routes
- Không có API versioning (`/v1/`, `/v2/`)
- gRPC connections không được pooled — tạo 1 connection per service khi startup
