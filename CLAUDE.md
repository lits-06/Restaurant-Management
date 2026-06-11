# Restaurant Management System — CLAUDE.md

## Project Overview

Full-stack restaurant management system built with microservices architecture. Backend in Go (Clean Architecture + gRPC), three React frontends (customer app + admin dashboard + kitchen display), communicating via an HTTP/REST API Gateway.

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

> **Note:** `payment-service` and `staff-service` proto/files have been removed. `user-service` has been recreated with a new design (port 50056, 5 roles, PostgreSQL). `staff-service` was replaced by `schedule-service` (2026-06-10); the `services/staff-service/` and `proto/staff/` directories were deleted (2026-06-11).

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
- **Native WebSocket API** — real-time notifications (no socket.io)
- **lucide-react** — icons
- No react-router — state-based routing (`useState` for view switching)

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
docker-compose up          # starts postgres, redis, auth, menu, schedule, table, order, notification, user, api-gateway
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
AUTH_SERVICE_HOST=localhost     AUTH_SERVICE_PORT=50051
SCHEDULE_SERVICE_HOST=localhost SCHEDULE_SERVICE_PORT=50052
TABLE_SERVICE_HOST=localhost    TABLE_SERVICE_PORT=50053
MENU_SERVICE_HOST=localhost     MENU_SERVICE_PORT=50054
ORDER_SERVICE_HOST=localhost    ORDER_SERVICE_PORT=50055
USER_SERVICE_HOST=localhost     USER_SERVICE_PORT=50056
NOTIFICATION_SERVICE_HOST=localhost NOTIFICATION_SERVICE_PORT=50058
ENVIRONMENT=development
```

### Implemented Endpoints

Auth legend: `public` = no auth, `optional` = token used if present, `staff` = any staff role, `admin` = ADMIN/MANAGER only, `chef` = CHEF/ADMIN/MANAGER, `waiter` = WAITER/ADMIN/MANAGER

```
# ── Auth ──────────────────────────────────────────────────────── public
POST   /auth/register
POST   /auth/login
POST   /auth/refresh
POST   /auth/verify
POST   /auth/logout
POST   /auth/change-password

# ── Menu ──────────────────────────────────────────────────────── public (no role checks)
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

# ── Orders ────────────────────────────────────────────────────── mixed
GET    /orders                                   # public; filter by user_id → owner or staff only
POST   /orders                                   # optional auth; user_id auto-set from token if present
GET    /orders/{id}                              # optional; owner or staff
PUT    /orders/{id}                              # optional; owner or staff
DELETE /orders/{id}                              # optional; owner or staff
POST   /orders/{id}/cancel                       # optional; owner or staff
PATCH  /orders/{id}/status                       # staff (ADMIN/MANAGER/CHEF/WAITER)
POST   /orders/{id}/items                        # optional; owner or staff
DELETE /orders/{id}/items/{itemId}               # optional; owner or staff
PATCH  /orders/{id}/items/{itemId}/status        # chef (COOKING/READY) or waiter (SERVED)

# ── Schedule ──────────────────────────────────────────────────── staff required
GET    /schedule/shifts                          # staff; any role
POST   /schedule/shifts                          # staff; self-register or admin assigns others
GET    /schedule/shifts/{id}                     # staff; owner or admin
PUT    /schedule/shifts/{id}                     # staff; owner or admin
DELETE /schedule/shifts/{id}                     # staff; owner or admin

# ── Tables ────────────────────────────────────────────────────── public (no role checks)
GET    /tables
POST   /tables
GET    /tables/available
GET    /tables/{id}
PUT    /tables/{id}
DELETE /tables/{id}
PATCH  /tables/{id}/status

# ── Users ─────────────────────────────────────────────────────── public (no role checks)
GET    /users/by-email?email=...
GET    /users
POST   /users
GET    /users/{id}
PUT    /users/{id}
DELETE /users/{id}
GET    /users/{id}/roles
PATCH  /users/{id}/roles
PATCH  /users/{id}/password

# ── Misc ──────────────────────────────────────────────────────────────
GET    /health                                   # public

GET    /ws/notifications?token=<jwt>&role=<CHEF|WAITER>   # WebSocket; token required; role-validated
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

Pub/Sub message bus for kitchen staff. **No DB** — pure Redis Pub/Sub.

#### 2 RPCs:
- `SendNotification(SendNotificationRequest)` — publish notification to Redis channel `notifications:{target_role}`
- `Subscribe(SubscribeRequest) returns (stream Notification)` — subscribe Redis channel, stream notifications to client (used by WebSocket gateway)

#### Notification types:
| Type | Trigger | Target | Payload |
|------|---------|--------|---------|
| `ORDER_CONFIRMED` | `UpdateOrderStatus → Confirmed` | `CHEF` | order_id, table_id, customer_name, party_size, notes, items[] |
| `ITEM_READY` | `UpdateOrderItemStatus → READY` | `WAITER` | order_id, table_id, item_id, item_name |

#### Notification entity fields:
`id`, `type`, `target_role`, `order_id`, `table_id`, `item_id`, `item_name`, `created_at` (Unix), `message`, `customer_name`, `party_size`, `notes`, `items[]`

**Redis channels:** `notifications:CHEF`, `notifications:WAITER`

**Notification flow** (fire-and-forget, does not block order operations):
```
order-service usecase → notifClient.SendNotification (background goroutine)
  → notification-service gRPC → Redis PUBLISH notifications:{role}
  → api-gateway WebSocket → browser
```

**Config env vars:** `SERVER_PORT=50058`, `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_DB`

---

### Report Service (`services/report-service`, port 50059)
**Proto:** `proto/report/report.proto`
**Status:** gRPC handler implemented — **NOT wired into docker-compose or api-gateway yet.**

#### 7 RPCs:
| RPC | Request key fields | Response |
|-----|--------------------|----------|
| `GetSalesReport` | period, from_date, to_date | SalesReport (revenue, order_count, avg_order_value) |
| `GetInventoryReport` | as_of_date | InventoryReport (items with stock levels) |
| `GetOrderReport` | period, from_date, to_date | OrderReport (counts by status, hourly breakdown) |
| `GetStaffPerformanceReport` | period, from_date, to_date, staff_id | StaffPerformance[] (orders handled, avg time) |
| `GetPopularItemsReport` | period, from_date, to_date, top_n | PopularItem[] (name, order_count, revenue) |
| `GetRevenueAnalytics` | period, from_date, to_date | RevenueAnalytics (daily/monthly breakdown) |
| `ExportReport` | report_type, format, period, from_date, to_date | file_data (bytes), file_name, content_type |

**Known issues in current implementation:**
- `pkg/config/config.go` has `UserServiceAddr` defaulting to `localhost:50052` (wrong — should be `localhost:50056`)
- No HTTP routes registered in api-gateway
- Not in docker-compose

**Config env vars:** `SERVER_PORT=50059`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`, `ORDER_SERVICE_ADDR`, `MENU_SERVICE_ADDR`, `USER_SERVICE_ADDR`

---

### User Service (`services/user-service`, port 50056)
**Proto:** `proto/user/user.proto`
**Status:** Active — enabled in docker-compose, routes wired in api-gateway.

Manages user identity and access control for the entire system.

#### 5 Roles:
| Role | Description |
|------|-------------|
| `USER` | Customer — make reservations, view own orders (must be logged in) |
| `MANAGER` | Create walk-in orders, add items, update order status after payment |
| `CHEF` | View order items, mark items as cooked |
| `WAITER` | Receive notifications from chef, serve items to customers |
| `ADMIN` | Full access — dashboard, manage menu/users/orders/tables |

#### 10 RPCs:
CreateUser, GetUser, GetUserByEmail, UpdateUser, DeleteUser, ListUsers, AssignRole, GetUserRoles, ChangePassword, VerifyCredentials

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

Default role on user creation: `USER`

**DB table:** `users`

**Env:** `SERVER_PORT=50056`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`

DB schema: `users` table (no roles column) + `user_roles` junction table. All Create/Update use `sql.Tx` transactions.

---

### Schedule Service (`services/schedule-service`, port 50052)
**Proto:** `proto/schedule/schedule.proto`
**Status:** Active — replaced staff-service (2026-06-10).

Manages work shifts for kitchen staff. **No status field** — a shift's existence means it is confirmed. Delete = cancel shift.

- 5 RPCs: CreateShift, GetShift, UpdateShift, DeleteShift, ListShifts
- `Shift`: ShiftID, UserID, Date (YYYY-MM-DD), StartTime (HH:MM), EndTime, Role, Notes, CreatedBy, CreatedAt, UpdatedAt
- `user_id` links to user-service (UUID) — NOT enforced at DB level
- `ListShifts` filters: month (YYYY-MM), user_id, role, page, page_size
- Validation: user_id required, date format, StartTime < EndTime, Role in {CHEF, WAITER, MANAGER, ADMIN}

---

### Table Service (`services/table-service`, port 50053)
**Proto:** `proto/table/table.proto`
**Status:** Active — enabled in docker-compose, routes wired in api-gateway.

Manages **static information and physical status** of tables. Does not manage time-slot booking.

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
- Status reflects physical state only — unrelated to booking availability
- Time-slot availability is derived from the `orders` table in order-service

**DB table:** `restaurant_tables` (1 bảng duy nhất)

**`GetAvailableTables`** — called by order-service during auto-assign:
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
- `proto/schedule/schedule.proto`
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
    LoginPage.tsx    # Dark theme, checks CHEF/WAITER/ADMIN/MANAGER role
    KitchenPage.tsx  # CHEF view — confirmed orders, mark COOKING/READY per item
    WaiterPage.tsx   # WAITER view — READY items feed + notification sidebar, mark SERVED
    SchedulePage.tsx # My Schedule — view/register own shifts (self-service)
  App.tsx      # Auth gate + role routing; floating switcher 🍳/🛎/📅 for all roles
```

**Role routing** (`ActiveView = 'CHEF' | 'WAITER' | 'SCHEDULE'`):
- CHEF → KitchenPage default + floating switcher (🍳 Kitchen + 📅 Schedule)
- WAITER → WaiterPage default + floating switcher (🛎 Service + 📅 Schedule)
- ADMIN/MANAGER → KitchenPage default + floating switcher (🍳 Kitchen + 🛎 Service + 📅 Schedule)

**WebSocket:** `ws://localhost:8080/ws/notifications?token=<jwt>&role=<CHEF|WAITER>`
- KitchenPage (CHEF): receives `ORDER_CONFIRMED` → auto-refresh order list
- WaiterPage (WAITER): receives `ITEM_READY` → auto-refresh + show notification sidebar

**Note:** Kitchen app notifications still display `table_id` as truncated UUID (known gap) — customer app `MyOrdersPage` already fixed this (uses `tableApi.list()` to resolve table_number).

---

### restaurant-app (Customer)
**State management:** `src/store/authStore.ts` — Zustand + persist (`luxe-customer-auth`), holds `user`, `accessToken`, `refreshToken`.
**Routing:** `App.tsx` uses `currentPage` useState (not react-router). Protected pages: `reservation`, `my-orders` — redirect to `login` with `loginRedirect` state to return after success.
```
src/
  api/            # gateway.api.ts — token refresh interceptor + all API calls
  store/          # authStore.ts
  components/
    booking/      # Reservation components
    common/       # Button, Container, TopNavBar (auth-aware), SectionTitle
    home/         # Hero, FeaturedDishes, Philosophy, CTA, Testimonials
    layout/       # Header, Footer
    ui/           # UI primitives
  pages/
    HomePage, ContactPage
    MenuPage.tsx     # loads from API: GET /menu/categories + GET /menu/items; category filter tabs
    ReservationPage.tsx  # POST /orders flow; API menu for pre-order
    LoginPage.tsx    # login + register tabs; pre-fills email on register success
    MyOrdersPage.tsx # loads via ordersApi.list({user_id}); resolves table_id→table_number; add items, 2-step cancel
  hooks/          # useAutoSlider, useHeaderScroll
  lib/            # constants.ts, types.ts
  assets/images/  # food images (jpg)
```

**`gateway.api.ts` key exports:**
- `menuApi.listItems`, `menuApi.listCategories`
- `tableApi.getOne`, `tableApi.list`
- `ordersApi.create/getOne/list/cancel/addItem`
- `authApi.login/register/logout/changePassword`
- Token refresh: 401 → `POST /auth/refresh` → retry; `clearAuth()` on second 401

**TopNavBar** shows: logo + nav links + "My Orders" (auth only) + user name (`full_name || username`) + "Sign Out" (auth) OR "Sign In" + "Book Now" (anonymous).

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
    MonthlyScheduler.tsx   # monthly shift calendar grid; create + edit + delete shifts
    TableManagement.tsx    # CRUD bàn: create/edit/delete + status update (AVAILABLE/CLEANING/OUT_OF_SERVICE)
    UserManagement.tsx     # CRUD user + role assignment + change password
  services/       # api.ts — token refresh interceptor + all endpoints
```

**`services/api.ts` key exports:**
- `ordersApi`: `list`, `update` (notes+end_time), `updateStatus`, `updateItemStatus`, `cancel`, `delete`
- `tablesApi`: `list`, `create`, `update`, `delete`, `updateStatus`
- `menuApi`: `listItems`, `createItem`, `updateItem`, `deleteItem`, `listCategories`, `createCategory`, `updateCategory`, `deleteCategory`
- `usersApi`: `getOne`, `listAll`, `create`, `update`, `delete`, `assignRole`, `changePassword`
- `scheduleApi`: `list`, `create`, `update`, `delete`
- `authApi`: `login`, `logout`, `changePassword`
- Token refresh: 401 → `POST /auth/refresh` → retry; `clearAuth()` on second 401
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
9. **Redis Pub/Sub** — notification-service uses Redis channels (`notifications:CHEF`, `notifications:WAITER`) as message bus; order-service publishes fire-and-forget
10. **WebSocket + gRPC streaming** — api-gateway bridges browser WebSocket → notification-service `Subscribe` streaming RPC; context cancel on WebSocket close automatically closes gRPC stream

---

## Current State (as of 2026-06-11)

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

### UI Localization — all 3 apps (2026-06-11)

All frontend text translated to **English**. No Vietnamese strings remain in any component.

**restaurant-app (customer):**
- All pages translated: TopNavBar, LoginPage, MenuPage, ReservationPage, MyOrdersPage
- Currency changed from USD (`$X.XX`) to VND (`X.000 ₫` via `toLocaleString('vi-VN')`, no decimals)
- `TopNavBar`: "Sign In" / "Sign Out" / "Book Now"
- `LoginPage`: login + register tabs; password confirm is client-side only (never sent to server)
- **`username` field note:** Collected at registration, required by user-service (min 3 chars). In customer app only used as fallback display in TopNavBar: `{user.full_name || user.username}`. Since `full_name` is also required, username is effectively never shown. Candidate for auto-generation from email prefix.

**restaurant-app-admin (admin):**
- All 9 pages + components translated: Login, Sidebar, HeaderDashboard, AnalyticsOverview, MenuManagement, TableManagement, UserManagement, MonthlyScheduler, OrdersManagement
- Month labels use `'en-US'` locale (January…December). Weekday headers: Mon/Tue/Wed/Thu/Fri/Sat/Sun

**restaurant-app-kitchen (kitchen display):**
- All 5 files translated: App.tsx, LoginPage, KitchenPage, WaiterPage, SchedulePage
- Item action buttons: "Start Cooking" / "Done ✓" (CHEF), "Served ✓" (WAITER)
- Notification sidebar: "Notifications" / "Clear" / "No notifications"
- Time locale in WaiterPage notification timestamps: `'en-US'`

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
- **Kitchen app**: `SchedulePage.tsx` — self-service shift registration/deletion. `scheduleApi` added to `gateway.api.ts`. App.tsx adds 📅 Schedule tab to floating switcher for all roles.
- **`docker-compose.yml`**: staff-service block → schedule-service. api-gateway env `SCHEDULE_SERVICE_*`.

### Admin dashboard update (2026-06-10)
- **`HeaderDashboard.tsx`**: Rewritten — accepts `year/month/onChange` props; functional year navigation; click month = apply immediately; future months disabled.
- **`AnalyticsOverview.tsx`**: Real month filtering — loads 500 orders once, filters client-side by selected month. 4 KPIs (revenue, orders, avg/order, covers) with trend % vs previous month. Order status breakdown row (Pending/Confirmed/Completed/Cancelled). Top 5 dishes ranked by revenue in selected month. Currency changed to VNĐ.
- **`OrdersManagement.tsx`**: Full update — table row shows `end_time`, table name (`Table X` resolved from `tablesApi`), notes (italic gold). Drawer shows notes banner, full `table_id`, item_status badges (Pending/Cooking/Done/Served), status action buttons (Confirm/Complete/Cancel Order) via `PATCH /orders/{id}/status`. Edit booking modal adds `end_time` and `notes` fields. Edit items modal replaces hardcoded text input with real menu search dropdown.
- **`MenuManagement.tsx`**: Fixed duplicate `mapMenuItem` bug. Category input changed to `<select>` dropdown from `menuApi.listCategories`. Payload now sends correct `category_id` (UUID), not category name.
- **`services/api.ts`**: Added `item_status` to `OrderItemDto`; added `table_id`, `user_id`, `notes`, `end_time`, `total` to `OrderDto`; added `tablesApi`; updated `ordersApi.update` to accept `notes`/`end_time` (removed `status` — use `updateStatus` instead); added `ordersApi.updateItemStatus`.

### Frontend API alignment + new admin pages (2026-06-11)

**Bug fixes:**
- **`MenuPage.tsx` (customer)**: Replaced hardcoded static `MENU_DATA` array with real API calls — `GET /menu/categories` (dynamic tabs) + `GET /menu/items?page_size=100` (filtered client-side by `category_id`).
- **`MyOrdersPage.tsx` (customer)**: Added `tableApi.list()` call on mount to build `Map<table_id, table_number>`; now shows "Table X" instead of UUID slice.
- **Token refresh interceptor (all 3 apps)**: `gateway.api.ts` / `services/api.ts` now detect HTTP 401, call `POST /auth/refresh` (deduplicating concurrent calls), retry original request. On refresh failure, `clearAuth()` is called. Pattern: module-level `let refreshing: Promise<string|null>|null` prevents parallel refresh storms.
- **`order-service` `UpdateOrder` preserve `item_status`**: `resolveItems` previously reset all items to PENDING. Fixed in `usecase/order_usecase.go` — after resolving new items, a `existingStatus` map is built from the current order; matching `item_id`s restore their COOKING/READY/SERVED status. New items default to PENDING.
- **`MonthlyScheduler.tsx` (admin)**: Added Edit shift — detail popover gains "Edit Shift" button → edit modal with date/start_time/end_time/notes → `PUT /schedule/shifts/{id}`.

**New admin pages (all backed by existing api-gateway endpoints):**
- **`TableManagement.tsx`**: Full CRUD for tables — grid card layout sorted by `table_number`. Create (`POST /tables`), edit (`PUT /tables/{id}`), delete (`DELETE /tables/{id}`), status update (`PATCH /tables/{id}/status`) with AVAILABLE/CLEANING/OUT_OF_SERVICE radio picker.
- **`UserManagement.tsx`**: Full user management — searchable table. Create (`POST /users`), edit (`PUT /users/{id}`), delete (`DELETE /users/{id}`), role assignment (`PATCH /users/{id}/roles` — checkbox for all 5 roles), change password (`PATCH /users/{id}/password` — requires old + new password).
- **`App.tsx` + `Sidebar.tsx` (admin)**: Added "Tables" and "Users" tabs to sidebar nav.

**New API functions added (no backend change needed):**
- `restaurant-app/gateway.api.ts`: `menuApi.listCategories()`, `tableApi.list()`, `authApi.changePassword()`
- `restaurant-app-admin/services/api.ts`: `authApi.changePassword`, `usersApi.create/update/delete/assignRole/changePassword`, `menuApi.createCategory/updateCategory/deleteCategory`, `tablesApi.create/update/delete/updateStatus`, `scheduleApi.update` (was already defined, now used)
- `restaurant-app-kitchen/gateway.api.ts`: token refresh interceptor added (no new endpoints)

### API audit + DeleteOrder fix (2026-06-11)
- **`services/order-service/internal/usecase/order_usecase.go`**: Added `DeleteOrder` usecase method — it was missing entirely (only the repository had `Delete`).
- **`services/order-service/internal/delivery/grpc/order_handler.go`**: Uncommented `DeleteOrder` gRPC handler — was commented out, causing HTTP `DELETE /orders/{id}` to return gRPC "Unimplemented" error.
- **Full API audit** (see below for complete RPC/HTTP inventory per service).

### Cleanup + infra fixes (2026-06-11)
- **`services/staff-service/`** deleted — fully replaced by `schedule-service`; no active references remained.
- **`proto/staff/`** deleted — generated `.pb.go` files removed alongside source `.proto`.
- **`docker-compose.yml`** fixes:
  - Added `AUTH_SERVICE_PORT=50051` to api-gateway env (was missing; worked by default, now explicit).
  - Added `menu-service` to `order-service` `depends_on` (order calls menu-service for item validation).
  - Added `schedule-service` to `api-gateway` `depends_on` (gateway connects to schedule-service).
- **`services/table-service/pkg/config/config.go`** rewritten — removed viper dependency that silently returned empty config when no YAML file present (every Docker deploy). Replaced with direct `os.Getenv` pattern consistent with all other services. `go mod tidy` removed orphaned viper dependency.
- **`scripts/Makefile`** updated — removed `staff-service`, added all 8 active services with individual `build-*` and `restart-*` targets.
- **`scripts/README.md`** updated — rewrote to reflect actual proto layout and Makefile-based workflow.

### Known gaps / TODO

**Backend:**
- **report-service** not wired into docker-compose or api-gateway (7 RPCs implemented, 0 HTTP routes)
- **`report-service/pkg/config/config.go`** has `UserServiceAddr` defaulting to wrong port (50052 instead of 50056)
- **No auth middleware** on menu/table/user routes in api-gateway — only order and schedule endpoints have auth
- **schedule-service** does not validate shift conflicts (same user, same day, overlapping time)
- **Operating hours** (10:00–22:00) only enforced at frontend — no backend validation

**Frontend:**
- **`table_id` in kitchen app** (KitchenPage/WaiterPage notification cards) still shows truncated UUID — needs `table_number` lookup from table-service (customer app MyOrdersPage already fixed this pattern)
- **Kitchen app WebSocket** has no auto-reconnect on disconnect; no polling fallback
- **AnalyticsOverview** loads max 500 orders — monthly stats incomplete if >500 orders exist
- **Category management UI** not in admin dashboard (API client has `createCategory/updateCategory/deleteCategory` but no dedicated page yet)

**Validation gaps (audited 2026-06-11):**
- **Email:** Frontend relies on `type="email"` browser validation. Backend `isValidEmail()` only checks `strings.Contains(email, "@") && strings.Contains(email, ".")` — no proper RFC check
- **Phone:** No validation anywhere — frontend uses `type="tel"` (no format check), backend only does `strings.TrimSpace()`
- **Password:** Frontend enforces `length >= 8`. Backend only checks `!= ""` — a 1-char password bypasses if sent directly to API
- **Username:** Required by user-service (min 3 chars), but in customer app only used as fallback display name if `full_name` is empty (since `full_name` is also required, username is never shown). Candidate for auto-generation from email prefix to reduce registration friction
