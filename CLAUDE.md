# Restaurant Management System — CLAUDE.md

## Project Overview

Full-stack restaurant management system built with microservices architecture. Backend in Go (Clean Architecture + gRPC), two React frontends (customer app + admin dashboard), communicating via an HTTP/REST API Gateway.

---

## Repository Structure

```
Restaurant_Management/
├── restaurant-app/          # Customer-facing React + Vite SPA (port 5173)
├── restaurant-app-admin/    # Admin dashboard React + Vite SPA (port 5174)
├── restaurant-app-kitchen/  # Kitchen display React + Vite SPA (port 5175) — CHEF + WAITER
├── api-gateway/             # HTTP → gRPC translation layer (port 8080)
├── services/                # Go microservices (gRPC)
│   ├── auth-service/        # port 50051
│   ├── schedule-service/    # port 50052
│   ├── table-service/       # port 50053
│   ├── menu-service/        # port 50054
│   ├── order-service/       # port 50055
│   ├── user-service/        # port 50056
│   ├── notification-service/# port 50058
│   └── report-service/      # port 50059
├── proto/                   # .proto definitions (source of truth for all gRPC contracts)
├── shared/                  # Shared Go packages (logger, db, jwt, middleware, errors)
├── scripts/                 # Proto generation scripts
├── docker-compose.yml       # Full stack orchestration
└── go.mod / go.sum          # Root Go module (restaurant-management)
```

> **Note:** `payment-service` proto/files have been removed. `user-service` has been recreated with a new design (port 50056, 5 roles, PostgreSQL).

---

## Tech Stack

### Frontend — restaurant-app & restaurant-app-admin
- **React 19** + **Vite 8** + **TypeScript**
- **TailwindCSS** — styling
- **Zustand** — state management
- **@tanstack/react-query** — server state / data fetching
- **react-router-dom v7** — routing
- **react-hook-form + zod** — form validation
- **axios** — HTTP client
- **@radix-ui/react-*** — accessible UI primitives
- **lucide-react** — icons

### Frontend — restaurant-app-kitchen
- **React 19** + **Vite 8** + **TypeScript**
- **TailwindCSS** — styling
- **Zustand** — state management (persist key: `kitchen-auth`)
- **Native WebSocket API** — real-time notifications (không dùng socket.io)
- **lucide-react** — icons
- No react-router — state-based routing (`useState` cho view switching)

### Backend (all services)
- **Go 1.25+** — each service has its **own `go.mod`** (module path: `restaurant-management/services/<name>-service`)
- **gRPC** (`google.golang.org/grpc v1.80`)
- **Protocol Buffers v3**
- **PostgreSQL 15** — primary database (single shared `restaurant_db`)
- **Redis 7** — token/session cache (auth-service) + Pub/Sub message bus (notification-service)
- **go.uber.org/zap** — structured logging
- **github.com/golang-jwt/jwt/v5** — JWT auth
- **github.com/spf13/viper** — config management
- **golang.org/x/crypto** — bcrypt password hashing

---

## Port Map

| Component            | Protocol  | Port  |
|----------------------|-----------|-------|
| customer frontend    | HTTP      | 5173  |
| admin frontend       | HTTP      | 5174  |
| kitchen display      | HTTP      | 5175  |
| api-gateway          | HTTP/REST | 8080  |
| auth-service         | gRPC      | 50051 |
| schedule-service     | gRPC      | 50052 |
| table-service        | gRPC      | 50053 |
| menu-service         | gRPC      | 50054 |
| order-service        | gRPC      | 50055 |
| user-service         | gRPC      | 50056 |
| notification-service | gRPC      | 50058 |
| report-service       | gRPC      | 50059 |
| PostgreSQL           | TCP       | 5432  |
| Redis                | TCP       | 6379  |

---

## Development Commands

### Frontend
```bash
cd restaurant-app          # or restaurant-app-admin
npm install
npm run dev                # Vite dev server
npm run build
npm run lint
```

### Backend (individual service)
Each service has its own `go.mod` — **always build from within the service directory**:
```bash
cd services/order-service
go build ./cmd/server/
go run cmd/server/main.go
```

### Full Stack (Docker)
```bash
docker-compose up          # starts postgres, redis, auth, menu, staff, table, order, notification, user, api-gateway
docker-compose up --build  # rebuild images
docker-compose down
```

### Proto Generation
```bash
# Requires protoc-gen-go and protoc-gen-go-grpc to be in PATH
export PATH=$PATH:~/go/bin

# Regenerate a single proto
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/<service>/<file>.proto

# Or regenerate all
bash scripts/generate-proto.sh
```

---

## API Gateway

**File:** `api-gateway/cmd/server/main.go`

Translates HTTP REST → gRPC. Uses Go's standard `net/http` mux (no framework).

**CORS allowed origins (HTTP + WebSocket):** `http://localhost:5173`, `http://localhost:5174`, `http://localhost:5175`

**Environment Variables:**
```
SERVER_PORT=8080
AUTH_SERVICE_HOST=localhost   AUTH_SERVICE_PORT=50051
MENU_SERVICE_HOST=localhost   MENU_SERVICE_PORT=50054
ORDER_SERVICE_HOST=localhost  ORDER_SERVICE_PORT=50055
STAFF_SERVICE_HOST=localhost  STAFF_SERVICE_PORT=50052
TABLE_SERVICE_HOST=localhost  TABLE_SERVICE_PORT=50053
USER_SERVICE_HOST=localhost        USER_SERVICE_PORT=50056
NOTIFICATION_SERVICE_HOST=localhost NOTIFICATION_SERVICE_PORT=50058
ENVIRONMENT=development
```

### Implemented Endpoints

```
POST   /auth/register
POST   /auth/login
POST   /auth/refresh
POST   /auth/verify
POST   /auth/logout
POST   /auth/change-password

GET    /menu/items
POST   /menu/items
GET    /menu/items/{id}
PUT    /menu/items/{id}
DELETE /menu/items/{id}

GET    /menu/categories
POST   /menu/categories
GET    /menu/categories/{id}
PUT    /menu/categories/{id}
DELETE /menu/categories/{id}

GET    /orders
POST   /orders
GET    /orders/{id}
PUT    /orders/{id}
DELETE /orders/{id}
POST   /orders/{id}/cancel
PATCH  /orders/{id}/status
POST   /orders/{id}/items
DELETE /orders/{id}/items/{itemId}
PATCH  /orders/{id}/items/{itemId}/status

GET    /staff
POST   /staff
GET    /staff/{id}
PUT    /staff/{id}
DELETE /staff/{id}

GET    /tables
POST   /tables
GET    /tables/available
GET    /tables/{id}
PUT    /tables/{id}
DELETE /tables/{id}
PATCH  /tables/{id}/status

GET    /users/by-email?email=...
GET    /users
POST   /users
GET    /users/{id}
PUT    /users/{id}
DELETE /users/{id}
GET    /users/{id}/roles
PATCH  /users/{id}/roles
PATCH  /users/{id}/password

GET    /health

GET    /ws/notifications?token=<jwt>&role=<CHEF|WAITER>   (WebSocket upgrade)
```

---

## Backend Service Architecture

All services follow **Clean Architecture**:

```
services/<name>-service/
├── cmd/server/main.go       # Entry point, wires all layers
├── internal/
│   ├── domain/              # Entities, interfaces, errors (no external deps)
│   ├── repository/          # DB access (PostgreSQL implementation)
│   ├── usecase/             # Business logic
│   └── delivery/grpc/       # gRPC handler (calls usecase)
└── pkg/config/              # Config struct + env loader
```

---

### Auth Service (`services/auth-service`, port 50051)
**Proto:** `proto/auth/auth.proto`

- 6 RPCs: Register, Login, RefreshToken, VerifyToken, Logout, ChangePassword
- **No PostgreSQL** — delegates all user operations to user-service via gRPC
- **Redis only** — stores refresh tokens as `refresh_token:<token>` hash keys
- JWT (HS256): short-lived access token + long-lived refresh token
  - Claims: `user_id`, `email`, `roles []string`
  - Secret from `JWT_SECRET` env var (logs warning if default is used)
- **Login flow**: call user-service `VerifyCredentials(email, password)` → get user_id + email + roles → `GenerateAccessToken(user_id, email, roles)` → store refresh token in Redis
- **Logout**: deletes the refresh token from Redis using `refresh_token` field in request (not access token)
- **RefreshToken**: validate JWT → check Redis → call user-service `GetUser(user_id)` for fresh email+roles → new access token
- **ChangePassword**: delegates directly to user-service `ChangePassword` RPC
- **`LogoutRequest`** has `refresh_token` field (required for proper token revocation)
- **`VerifyTokenResponse`** includes `repeated string roles`

**Config env vars:** `SERVER_PORT`, `USER_SERVICE_ADDR`, `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`, `JWT_SECRET`, `JWT_ACCESS_MINUTES`, `JWT_REFRESH_HOURS`

---

### Menu Service (`services/menu-service`, port 50054)
**Proto:** `proto/menu/menu.proto`

- 10 RPCs: CRUD for `MenuItem` + CRUD for `Category`
- `MenuItem`: ID, name, price, description, category_id, image_url
- `Category`: ID, name
- Filtering by category/keyword, pagination

---

### Order Service (`services/order-service`, port 50055)
**Proto:** `proto/order/order.proto`

- **10 RPCs:** CreateOrder, GetOrder, UpdateOrder, DeleteOrder, CancelOrder, ListOrders, UpdateOrderStatus, AddOrderItem, RemoveOrderItem, **UpdateOrderItemStatus**
- **4 order statuses:** `Pending → Confirmed → Completed` (or `Cancelled` from Pending/Confirmed)
- **4 item statuses:** `PENDING → COOKING → READY → SERVED` (enforced state machine, role-gated)
- Validates menu items via menu-service gRPC (price lookup)
- **Auto-assigns table** when `table_id` is omitted in CreateOrder (calls table-service `GetAvailableTables` + checks own DB for time conflicts)
- **Notifies kitchen staff** via notification-service when order → Confirmed (CHEF) or item → READY (WAITER)

**`Order` entity fields:**
| Field | Type | Notes |
|-------|------|-------|
| `order_id` | string | UUID, DB-generated |
| `table_id` | string | Assigned table — auto-filled if omitted |
| `user_id` | string | Auth link to user-service (empty for walk-in orders, immutable after creation) |
| `name` | string | Customer name — snapshot at booking time, not updated if user profile changes |
| `phone` | string | Customer phone — snapshot at booking time |
| `notes` | string | Special requests (allergies, preferences, etc.) — always overwritten on update |
| `time` | timestamp | Start time of the booking |
| `end_time` | timestamp | End time |
| `party_size` | int32 | Number of guests |
| `status` | string | Pending/Confirmed/Completed/Cancelled |
| `total` | float64 | Sum of item prices |
| `items` | OrderItem[] | Menu items + quantity + `item_status` |

**`OrderItem` fields:** `item_id`, `name`, `price`, `category`, `image_url`, `quantity`, `item_status` (PENDING/COOKING/READY/SERVED)

**DB tables:** `orders`, `order_items` (FK → `menu_items`, includes `item_status VARCHAR(16) NOT NULL DEFAULT 'PENDING'`)

**Auto-assign logic** (`usecase/order_usecase.go → autoAssignTable`):
1. Call `table-service.GetAvailableTables(min_capacity=party_size)` — tables sorted by capacity ASC
2. Query own `orders` table for table IDs with conflicting time window (status != Cancelled)
3. Pick first candidate not in the occupied set (best-fit)
4. Returns `ErrNoTableAvailable` if all candidates are booked

**Kitchen flow:**
```
Staff confirm order → UpdateOrderStatus(Confirmed) → 🔔 CHEF notified (ORDER_CONFIRMED)
CHEF → UpdateOrderItemStatus(COOKING) → UpdateOrderItemStatus(READY) → 🔔 WAITER notified (ITEM_READY)
WAITER → UpdateOrderItemStatus(SERVED)
Manager → UpdateOrderStatus(Completed) when customer pays (manual)
```

**Env:**
- `TABLE_SERVICE_ADDR=table-service:50053` — empty = auto-assign disabled, `table_id` required
- `NOTIFICATION_SERVICE_ADDR=notification-service:50058` — empty = notifications silently skipped
- `MENU_SERVICE_ADDR=menu-service:50054`

**CreateOrder request** (HTTP JSON):
```json
{
  "name": "Nguyen Van A",
  "phone": "0901234567",
  "notes": "dị ứng hải sản, cần ghế cao cho trẻ em",
  "date": "2026-06-15",
  "time": "19:00",
  "end_time": "21:00",
  "party_size": 4,
  "items": [{ "item_id": "...", "quantity": 2 }]
}
```
`table_id` optional — system auto-assigns. Admin can override via UpdateOrder.
`user_id` — extracted from JWT by api-gateway if token is present; omitted by anonymous/walk-in callers.

**Authorization model (api-gateway `order_handler.go`):**
- `CreateOrder`: optional auth — extracts `user_id` from token if present, no auth required
- `GetOrder`, `UpdateOrder`, `DeleteOrder`, `CancelOrder`, `AddOrderItem`, `RemoveOrderItem`: pre-fetches order; if `order.user_id != ""`, caller must be authenticated as owner OR have staff role
- `ListOrders`: if `user_id` query param present, caller must be that user or staff; no filter = open
- `UpdateOrderStatus`: staff-only (ADMIN/MANAGER/CHEF/WAITER)
- `UpdateOrderItemStatus`: role-gated by target status — `COOKING`/`READY` → CHEF+ADMIN+MANAGER; `SERVED` → WAITER+ADMIN+MANAGER
- Shared helpers: `verifyCaller(r)` (optional token extract), `checkOrderAccess(r, orderUserID)` (combined auth+authz), `checkUserIDAccess(r, targetUserID)` (for list filter), `canMarkItemStatus(roles, targetStatus)`

---

### Notification Service (`services/notification-service`, port 50058)
**Proto:** `proto/notification/notification.proto`
**Status:** Active — enabled in docker-compose.

Pub/Sub message bus cho kitchen staff. **Không có DB** — dùng Redis Pub/Sub thuần túy.

#### 2 RPCs:
- `SendNotification(SendNotificationRequest)` — publish notification vào Redis channel `notifications:{target_role}`
- `Subscribe(SubscribeRequest) returns (stream Notification)` — subscribe Redis channel, stream notifications đến client (dùng cho WebSocket gateway)

#### Notification types:
| Type | Trigger | Target | Payload |
|------|---------|--------|---------|
| `ORDER_CONFIRMED` | `UpdateOrderStatus → Confirmed` | `CHEF` | order_id, table_id, customer_name, party_size, notes, items[] |
| `ITEM_READY` | `UpdateOrderItemStatus → READY` | `WAITER` | order_id, table_id, item_id, item_name |

#### Notification entity fields:
`id`, `type`, `target_role`, `order_id`, `table_id`, `item_id`, `item_name`, `created_at` (Unix), `message`, `customer_name`, `party_size`, `notes`, `items[]`

**Redis channels:** `notifications:CHEF`, `notifications:WAITER`

**Notification flow** (fire-and-forget, không block order operations):
```
order-service usecase → notifClient.SendNotification (background goroutine)
  → notification-service gRPC → Redis PUBLISH notifications:{role}
  → (future) api-gateway WebSocket → browser
```

**Config env vars:** `SERVER_PORT=50058`, `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`

---

### User Service (`services/user-service`, port 50056)
**Proto:** `proto/user/user.proto`
**Status:** Active — enabled in docker-compose, routes wired in api-gateway.

Quản lý **danh tính và phân quyền** cho toàn bộ người dùng hệ thống.

#### 5 Roles:
| Role | Mô tả |
|------|-------|
| `USER` | Khách hàng — đặt bàn, xem order của bản thân (phải đăng nhập) |
| `MANAGER` | Tạo order walk-in, thêm món, cập nhật trạng thái order sau thanh toán |
| `CHEF` | Xem món trong order, đánh dấu món đã nấu xong |
| `WAITER` | Nhận thông báo từ chef, mang món ra phục vụ |
| `ADMIN` | Toàn quyền — dashboard, quản lý menu/user/order/bàn |

#### 9 RPCs:
CreateUser, GetUser, GetUserByEmail, UpdateUser, DeleteUser, ListUsers, AssignRole, GetUserRoles, ChangePassword

#### User entity:
| Field | Type | Notes |
|-------|------|-------|
| `user_id` | string | UUID |
| `email` | string | Unique |
| `username` | string | Unique, min 3 chars |
| `full_name` | string | |
| `phone` | string | |
| `roles` | []UserRole | Comma-separated in DB (e.g. `"USER,ADMIN"`) |
| `status` | UserStatus | ACTIVE / INACTIVE / SUSPENDED |

Default role khi tạo user: `USER`

**DB table:** `users`

**Env:** `SERVER_PORT=50056`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`

DB schema: `users` table (no roles column) + `user_roles` junction table. All Create/Update use `sql.Tx` transactions.

---

### Schedule Service (`services/schedule-service`, port 50052)
**Proto:** `proto/schedule/schedule.proto`
**Status:** Active — replaced staff-service (2026-06-10).

Quản lý ca làm việc cho nhân viên bếp. **Không có status field** — shift tồn tại = đã xác nhận. Xóa = hủy ca.

- 5 RPCs: CreateShift, GetShift, UpdateShift, DeleteShift, ListShifts
- `Shift`: ShiftID, UserID, Date (YYYY-MM-DD), StartTime (HH:MM), EndTime, Role, Notes, CreatedBy, CreatedAt, UpdatedAt
- `user_id` links to user-service (UUID) — NOT enforced at DB level
- `ListShifts` filters: month (YYYY-MM), user_id, role, page, page_size
- Validation: user_id required, date format, StartTime < EndTime, Role in {CHEF, WAITER, MANAGER, ADMIN}

---

### Table Service (`services/table-service`, port 50053)
**Proto:** `proto/table/table.proto`
**Status:** Active — enabled in docker-compose, routes wired in api-gateway.

Quản lý **thông tin tĩnh và trạng thái vật lý** của bàn. Không quản lý time-slot booking.

#### 7 RPCs:
CreateTable, GetTable, UpdateTable, DeleteTable, ListTables, UpdateTableStatus, GetAvailableTables

#### Table entity:
| Field | Type | Notes |
|-------|------|-------|
| `table_id` | string | UUID |
| `table_number` | int32 | Unique positive integer (1, 2, 3…) |
| `capacity` | int32 | 1–50 |
| `status` | TableStatus | AVAILABLE / CLEANING / OUT_OF_SERVICE |

#### Table statuses: `AVAILABLE`, `CLEANING`, `OUT_OF_SERVICE`
- Status chỉ phản ánh trạng thái vật lý — không liên quan đến booking
- Time-slot availability được tính từ bảng `orders` trong order-service

**DB table:** `restaurant_tables` (1 bảng duy nhất)

**`GetAvailableTables`** — được gọi bởi order-service khi auto-assign:
```sql
SELECT * FROM restaurant_tables
WHERE status = 'AVAILABLE' AND capacity >= $min_capacity
ORDER BY capacity ASC, table_number ASC
```

---

## Shared Packages (`shared/pkg/`)

| Package      | Purpose |
|--------------|---------|
| `config/`    | Env var loading with defaults |
| `database/`  | PostgreSQL connection helper |
| `logger/`    | Zap logger wrapper (`logger.Log`) |
| `jwt/`       | JWT creation, validation, refresh |
| `middleware/`| gRPC logging/recovery interceptors, auth middleware |
| `errors/`    | Domain error types + gRPC error code mapping |
| `utils/`     | Common helpers |

---

## Protocol Buffers

Active proto definitions (source of truth — do not edit generated `*.pb.go` files directly):
- `proto/auth/auth.proto`
- `proto/menu/menu.proto`
- `proto/order/order.proto`
- `proto/staff/staff.proto`
- `proto/table/table.proto`
- `proto/user/user.proto`
- `proto/notification/notification.proto`
- `proto/report/report.proto`

---

## Frontend Apps

### restaurant-app-kitchen (Kitchen Display)
**Port 5175** — Dành cho CHEF và WAITER. Auth gate giống admin app.

```
src/
  store/       # authStore.ts (key: 'kitchen-auth') — KitchenUser, hasKitchenAccess, getDefaultRole
  api/         # gateway.api.ts — ordersApi, authApi, usersApi, scheduleApi, createNotificationWS
  hooks/       # useNotifications.ts — WebSocket hook, giữ tối đa 50 notif
  pages/
    LoginPage.tsx    # Dark theme, kiểm tra CHEF/WAITER/ADMIN/MANAGER role
    KitchenPage.tsx  # CHEF view — confirmed orders, mark COOKING/READY per item
    WaiterPage.tsx   # WAITER view — READY items feed + notification sidebar, mark SERVED
    SchedulePage.tsx # My schedule — view/register own shifts (self-service)
  App.tsx      # Auth gate + role routing; floating switcher 🍳/🛎/📅 cho tất cả roles
```

**Role routing** (`ActiveView = 'CHEF' | 'WAITER' | 'SCHEDULE'`):
- CHEF → KitchenPage mặc định + floating switcher (🍳 Bếp + 📅 Lịch)
- WAITER → WaiterPage mặc định + floating switcher (🛎 Phục vụ + 📅 Lịch)
- ADMIN/MANAGER → KitchenPage mặc định + floating switcher (🍳 Bếp + 🛎 Phục vụ + 📅 Lịch)

**WebSocket:** `ws://localhost:8080/ws/notifications?token=<jwt>&role=<CHEF|WAITER>`
- KitchenPage (CHEF): nhận `ORDER_CONFIRMED` → auto-refresh danh sách order
- WaiterPage (WAITER): nhận `ITEM_READY` → auto-refresh + hiển thị notification sidebar

**Known gap:** `table_id` hiển thị dạng UUID (8 ký tự đầu) — chưa resolve `table_number` từ table-service.

---

### restaurant-app (Customer)
**State management:** `src/store/authStore.ts` — Zustand + persist (`luxe-customer-auth`), holds `user`, `accessToken`, `refreshToken`.
**Routing:** `App.tsx` uses `currentPage` useState (not react-router). Protected pages: `reservation`, `my-orders` — redirect to `login` with `loginRedirect` state to return after success.
```
src/
  api/            # gateway.api.ts — all API calls, auto-injects Authorization header
  store/          # authStore.ts
  components/
    booking/      # Reservation components
    common/       # Button, Container, TopNavBar (auth-aware), SectionTitle
    home/         # Hero, FeaturedDishes, Philosophy, CTA, Testimonials
    layout/       # Header, Footer
    ui/           # UI primitives
  pages/
    HomePage, MenuPage, ReservationPage, ContactPage
    LoginPage.tsx    # login + register tabs; pre-fills email on register success
    MyOrdersPage.tsx # loads via ordersApi.list({user_id}); add items, 2-step cancel
  hooks/          # useAutoSlider, useHeaderScroll
  lib/            # constants.ts, types.ts
  assets/images/  # food images (jpg)
```

**TopNavBar** shows: logo + nav links + "My Orders" (auth only) + user name + logout (auth) OR "Đăng nhập" + "Book Now" (anonymous).

### restaurant-app-admin (Admin Dashboard)
**State management:** `src/store/adminAuthStore.ts` — Zustand + persist (`luxe-admin-auth`). `hasAdminAccess(roles)` checks for ADMIN/MANAGER/CHEF/WAITER.
**Auth gate:** `App.tsx` renders `<LoginPage>` if no user in store.
```
src/
  store/          # adminAuthStore.ts
  components/     # Footer, HeaderDashboard (year/month picker), KPIGrid, PerformanceTable, Sidebar
  pages/
    Login.tsx            # wired to authApi; validates staff role before granting access
    AnalyticsOverview.tsx  # real month filtering, trend vs prev month, status breakdown
    MenuManagement.tsx     # CRUD with category dropdown from API
    OrdersManagement.tsx   # full order mgmt: status update, notes, end_time, table name, item_status
    MonthlyScheduler.tsx   # monthly shift calendar grid
  services/       # api.ts (injects auth token in all requests), auth.ts
```

**`services/api.ts` key exports:**
- `ordersApi`: `list`, `update` (notes+end_time), `updateStatus`, `updateItemStatus`, `cancel`, `delete`
- `tablesApi.list` — load all tables for UUID → table_number resolution
- `menuApi`, `scheduleApi`, `usersApi`, `authApi`
- `TableDto`: `table_id, table_number, capacity, status`

---

## Infrastructure

### Docker Compose
**Dev credentials:**
```
POSTGRES_USER=restaurant_user
POSTGRES_PASSWORD=restaurant_pass
POSTGRES_DB=restaurant_db
```

**Active services:** postgres, redis, auth-service, menu-service, schedule-service, table-service, order-service, notification-service, user-service, api-gateway

Binaries are compiled locally and volume-mounted as `./server` inside each Alpine container.

---

## Key Architectural Patterns

1. **Clean Architecture** — domain → repository → usecase → delivery; dependencies point inward only
2. **Per-service go.mod** — each service is its own Go module linked to root via `replace restaurant-management => ../..`
3. **Microservices via gRPC** — inter-service calls use generated gRPC clients; api-gateway clients live in `api-gateway/internal/grpcclient/`
4. **API Gateway** — single HTTP entry point; handles CORS, translates JSON ↔ protobuf; WebSocket upgrade for `/ws/notifications`
5. **Repository Pattern** — PostgreSQL implementations; `ensureSchema` auto-creates tables + `ALTER TABLE ADD COLUMN IF NOT EXISTS` for safe migrations
6. **JWT + Redis** — stateless access tokens, revocable refresh tokens stored in Redis
7. **Proto-first design** — `.proto` files define all contracts; regenerate with `scripts/generate-proto.sh`
8. **Structured logging** — `go.uber.org/zap` via `shared/pkg/logger`
9. **Redis Pub/Sub** — notification-service dùng Redis channels (`notifications:CHEF`, `notifications:WAITER`) làm message bus; order-service publish fire-and-forget
10. **WebSocket + gRPC streaming** — api-gateway bridges browser WebSocket → notification-service `Subscribe` streaming RPC; context cancel trên WebSocket close tự đóng gRPC stream

---

## Current State (as of 2026-06-10)

### Backend refactor (2026-06-09)
- **table-service**: Stripped down to physical registry only. Removed all Reservation RPCs, `location` field, `OCCUPIED`/`RESERVED` statuses. Now has 7 RPCs, 1 DB table, 3 statuses (AVAILABLE/CLEANING/OUT_OF_SERVICE). Re-enabled in docker-compose.
- **order-service**: Added auto-assign table logic — `table_id` is now optional in `CreateOrder`. Added `notes` field (customer special requests). Added `GetOccupiedTableIDs` query for time-conflict detection. Added `TABLE_SERVICE_ADDR` config.
- **proto/table**: Completely rewritten — `table_number` changed to `int32`, all Reservation messages/RPCs removed.
- **proto/order**: Added `notes`, `end_time` fields.
- **api-gateway**: Table handler rewritten (removed reservation routes). Order handler wired with notes support.
- **docker-compose**: table-service enabled; order-service gets `TABLE_SERVICE_ADDR=table-service:50053`.

### Frontend update (2026-06-09)
- **`restaurant-app/src/api/gateway.api.ts`**: Removed `ReservationDto`, `tableApi.createReservation` (called `/reservations` which no longer exists). Updated `ordersApi.create` to accept `end_time`, `notes`, optional `table_id`. Updated `TableDto` (`table_number: number`, removed `location`). Replaced `tableApi.getAvailable` with `tableApi.getOne`.
- **`restaurant-app/src/pages/ReservationPage.tsx`**: Replaced old reservation flow with `POST /orders`. Removed explicit table selection (auto-assign). Added `selectedEndTime` (default start+2h, user-adjustable). Time pickers bounded to operating hours 10:00–22:00 (frontend-only constraint). Simplified billing to pre-order subtotal only.

### In progress / uncommitted changes
- `restaurant-app-admin/src/pages/MenuManagement.tsx` — modified
- `restaurant-app-admin/src/pages/OrdersManagement.tsx` — modified
- Food item images added to both apps (`src/assets/images/`, `public/images/`) — untracked

### User Service + auth-service integration (2026-06-09)
- **`proto/user/user.proto`**: 10 RPCs — added `VerifyCredentials(email, password)` returning `{success, message, user_id, email, roles[]}`. Roles use `user_roles` junction table.
- **`services/user-service`**: Two-table schema — `users` (no roles column) + `user_roles` junction table. `ensureSchema` includes `ALTER TABLE users DROP COLUMN IF EXISTS roles` for safe migration. All Create/Update use `sql.Tx` transactions. Added `VerifyCredentials` usecase + handler. Added `ErrInvalidCredentials`, `ErrAccountSuspended` domain errors.
- **`services/auth-service`**: Complete rewrite — removed PostgreSQL, removed `auth_users` table. Now Redis-only. Delegates all user ops to user-service via `internal/grpcclient/user_client.go`. JWT claims now include `roles []string`. Logout properly revokes refresh token.
- **`shared/pkg/jwt`**: `GenerateAccessToken` signature updated to accept `roles []string`. `Claims` struct has `Roles []string`.
- **`api-gateway`**: User gRPC client + HTTP handler added. 10 routes wired under `/users/`. `logoutRequest` updated with `refresh_token` field.
- **`docker-compose.yml`**: user-service added on port 50056. auth-service: removed DATABASE_* vars, added `JWT_SECRET`, `JWT_ACCESS_MINUTES`, `JWT_REFRESH_HOURS`, `USER_SERVICE_ADDR`. api-gateway gets `USER_SERVICE_HOST/PORT` env vars.

### Auth + authorization on order endpoints (2026-06-10)
- **`proto/order/order.proto`**: Added `user_id` to `Order` (field 12), `CreateOrderRequest` (field 11), `ListOrdersRequest` (field 5).
- **`services/order-service`**: Full user_id propagation — domain, repository (DB column + index), usecase, gRPC handler. `user_id` is immutable after creation (UPDATE explicitly excludes it). `ListOrders` filters by `user_id` when non-empty.
- **`api-gateway/internal/handler/order_handler.go`**: Added `verifyCaller` (optional token extract), `checkOrderAccess` (combined auth+authz for owned orders), `checkUserIDAccess` (for list filter). All order endpoints now enforce ownership/role checks. `UpdateOrderStatus` is staff-only.
- **Customer app auth flow**: `store/authStore.ts` (Zustand persist), `LoginPage.tsx` (login+register), `MyOrdersPage.tsx` (loads via `ordersApi.list({user_id})`). `App.tsx` routes with `loginRedirect` state — redirects back to intended page after login. `TopNavBar.tsx` auth-aware.
- **Admin app auth gate**: `store/adminAuthStore.ts`, `Login.tsx` wired to API with role validation, `App.tsx` renders login if no user, `Sidebar.tsx` shows real user info + logout, `services/api.ts` injects Bearer token.

### Item status + notification-service (2026-06-10)
- **`proto/order/order.proto`**: Added `item_status` field (string) to `OrderItem`. Added `UpdateOrderItemStatus` RPC (10th RPC).
- **`services/order-service`**: `ItemStatus` type with state machine (`PENDING→COOKING→READY→SERVED`). `UpdateItemStatus` method on `Order` domain entity. `UpdateItemStatus` SQL (targeted row update). `UpdateOrderItemStatus` usecase + gRPC handler. `notifyChef` / `notifyWaiter` background goroutines triggered on status change.
- **`proto/notification/notification.proto`**: New proto — `SendNotification` + `Subscribe` (streaming) RPCs.
- **`services/notification-service`**: New microservice (port 50058) — Redis Pub/Sub backend, no DB. `PubSubRepository`, `NotificationUseCase`, gRPC handler.
- **`api-gateway`**: `canMarkItemStatus(roles, targetStatus)` helper — CHEF marks COOKING/READY, WAITER marks SERVED. New route `PATCH /orders/{id}/items/{itemId}/status`.
- **`docker-compose.yml`**: notification-service added. order-service gets `NOTIFICATION_SERVICE_ADDR`.

### WebSocket + Kitchen app (2026-06-10)
- **`api-gateway`**: `gorilla/websocket` dependency added. `GET /ws/notifications?token=<jwt>&role=<CHEF|WAITER>` endpoint. `NotificationHandler` — upgrade WebSocket, verify JWT, open `notification-service.Subscribe` gRPC streaming, forward JSON notifications. Origin check: 5173/5174/5175. gRPC stream cancelled when WebSocket closes via drain goroutine.
- **`api-gateway/internal/grpcclient/notification_client.go`**: `NotificationClient` wrapping gRPC streaming `Subscribe` RPC.
- **`restaurant-app-kitchen/`**: New React app (port 5175) — dark-themed kitchen display. `KitchenPage` (CHEF): confirmed orders, mark COOKING/READY. `WaiterPage` (WAITER): READY items feed + notification sidebar, mark SERVED. Role-auto-routing; ADMIN/MANAGER floating view switcher.
- **`docker-compose.yml`**: api-gateway gets `NOTIFICATION_SERVICE_HOST/PORT` env vars.

### Schedule service — replaced staff-service (2026-06-10)
- **`proto/schedule/schedule.proto`**: New proto — 5 RPCs (CreateShift/GetShift/UpdateShift/DeleteShift/ListShifts). `Shift` message has no status field.
- **`services/schedule-service/`**: New microservice (port 50052). Clean Architecture. `shifts` DB table. `ListShifts` filters by `month/user_id/role`. Validation: date format, EndTime > StartTime, role in allowed set.
- **`api-gateway`**: `schedule_client.go` + `schedule_handler.go` (5 routes under `/schedule/shifts`). Auth required on all endpoints. POST: self-register any staff; assign others → ADMIN/MANAGER only. GET/PUT/DELETE: owner or ADMIN/MANAGER.
- **Admin app**: `MonthlyScheduler.tsx` — monthly calendar grid (replaced `StaffManagement.tsx` + `WeeklyScheduler.tsx`). Color-coded chips by role. Create/delete shifts. Loads staff names from `/users?page_size=200`.
- **Kitchen app**: `SchedulePage.tsx` — self-service shift registration/deletion. `scheduleApi` added to `gateway.api.ts`. App.tsx adds 📅 Lịch tab to floating switcher for all roles.
- **`docker-compose.yml`**: staff-service block → schedule-service. api-gateway env `SCHEDULE_SERVICE_*`.

### Admin dashboard update (2026-06-10)
- **`HeaderDashboard.tsx`**: Rewritten — accepts `year/month/onChange` props; functional year navigation; click month = apply immediately; future months disabled.
- **`AnalyticsOverview.tsx`**: Real month filtering — loads 500 orders once, filters client-side by selected month. 4 KPIs (revenue, orders, avg/order, covers) with trend % vs previous month. Order status breakdown row (Pending/Confirmed/Completed/Cancelled). Top 5 dishes ranked by revenue in selected month. Currency changed to VNĐ.
- **`OrdersManagement.tsx`**: Full update — table row shows `end_time`, table name (`Bàn X` resolved from `tablesApi`), notes (italic gold). Drawer shows notes banner, full `table_id`, item_status badges (Chờ/Đang nấu/Xong/Đã mang), status action buttons (Xác nhận/Hoàn thành/Hủy) via `PATCH /orders/{id}/status`. Edit booking modal adds `end_time` and `notes` fields. Edit items modal replaces hardcoded text input with real menu search dropdown.
- **`MenuManagement.tsx`**: Fixed duplicate `mapMenuItem` bug. Category input changed to `<select>` dropdown from `menuApi.listCategories`. Payload now sends correct `category_id` (UUID), not category name.
- **`services/api.ts`**: Added `item_status` to `OrderItemDto`; added `table_id`, `user_id`, `notes`, `end_time`, `total` to `OrderDto`; added `tablesApi`; updated `ordersApi.update` to accept `notes`/`end_time` (removed `status` — use `updateStatus` instead); added `ordersApi.updateItemStatus`.

### Known gaps / TODO
- Operating hours (10:00–22:00) enforced only in frontend — no backend validation in order-service
- No role-based middleware for menu/schedule/table/user routes in api-gateway — only order endpoints have auth
- `table_id` trong notification/kitchen app hiển thị dạng UUID cắt ngắn — chưa resolve `table_number` từ table-service
- Kitchen app không có auto-reconnect khi WebSocket mất kết nối
- schedule-service không validate trùng ca (cùng user, cùng ngày, cùng giờ)
- AnalyticsOverview load max 500 orders — nếu hệ thống có >500 orders thì thống kê theo tháng sẽ thiếu
