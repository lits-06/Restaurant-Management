# Restaurant Management System — CLAUDE.md

## Project Overview

Full-stack restaurant management system built with microservices architecture. Backend in Go (Clean Architecture + gRPC), two React frontends (customer app + admin dashboard), communicating via an HTTP/REST API Gateway.

---

## Repository Structure

```
Restaurant_Management/
├── restaurant-app/          # Customer-facing React + Vite SPA (port 5173)
├── restaurant-app-admin/    # Admin dashboard React + Vite SPA (port 5174)
├── api-gateway/             # HTTP → gRPC translation layer (port 8080)
├── services/                # Go microservices (gRPC)
│   ├── auth-service/        # port 50051
│   ├── staff-service/       # port 50052
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

### Frontend (both apps)
- **React 19** + **Vite 8** + **TypeScript**
- **TailwindCSS** — styling
- **Zustand** — state management
- **@tanstack/react-query** — server state / data fetching
- **react-router-dom v7** — routing
- **react-hook-form + zod** — form validation
- **socket.io-client** — real-time updates
- **axios** — HTTP client
- **@radix-ui/react-*** — accessible UI primitives
- **lucide-react** — icons

### Backend (all services)
- **Go 1.25+** — each service has its **own `go.mod`** (module path: `restaurant-management/services/<name>-service`)
- **gRPC** (`google.golang.org/grpc v1.80`)
- **Protocol Buffers v3**
- **PostgreSQL 15** — primary database (single shared `restaurant_db`)
- **Redis 7** — token/session cache (auth-service only)
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
| api-gateway          | HTTP/REST | 8080  |
| auth-service         | gRPC      | 50051 |
| staff-service        | gRPC      | 50052 |
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
docker-compose up          # starts postgres, redis, auth, menu, staff, table, order, user, api-gateway
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

**CORS allowed origins:** `http://localhost:5173`, `http://localhost:5174`

**Environment Variables:**
```
SERVER_PORT=8080
AUTH_SERVICE_HOST=localhost   AUTH_SERVICE_PORT=50051
MENU_SERVICE_HOST=localhost   MENU_SERVICE_PORT=50054
ORDER_SERVICE_HOST=localhost  ORDER_SERVICE_PORT=50055
STAFF_SERVICE_HOST=localhost  STAFF_SERVICE_PORT=50052
TABLE_SERVICE_HOST=localhost  TABLE_SERVICE_PORT=50053
USER_SERVICE_HOST=localhost   USER_SERVICE_PORT=50056
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

- 9 RPCs: CreateOrder, GetOrder, UpdateOrder, DeleteOrder, CancelOrder, ListOrders, UpdateOrderStatus, AddOrderItem, RemoveOrderItem
- **4 statuses:** `Pending → Confirmed → Completed` (or `Cancelled` from Pending/Confirmed)
- Validates menu items via menu-service gRPC (price lookup)
- **Auto-assigns table** when `table_id` is omitted in CreateOrder (calls table-service `GetAvailableTables` + checks own DB for time conflicts)

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
| `items` | OrderItem[] | Menu items + quantity |

**DB tables:** `orders`, `order_items` (FK → `menu_items`)

**Auto-assign logic** (`usecase/order_usecase.go → autoAssignTable`):
1. Call `table-service.GetAvailableTables(min_capacity=party_size)` — tables sorted by capacity ASC
2. Query own `orders` table for table IDs with conflicting time window (status != Cancelled)
3. Pick first candidate not in the occupied set (best-fit)
4. Returns `ErrNoTableAvailable` if all candidates are booked

**Env:** `TABLE_SERVICE_ADDR=table-service:50053` — if empty, auto-assign is disabled and `table_id` is required.

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
- Shared helpers: `verifyCaller(r)` (optional token extract), `checkOrderAccess(r, orderUserID)` (combined auth+authz), `checkUserIDAccess(r, targetUserID)` (for list filter)

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

### Staff Service (`services/staff-service`, port 50052)
**Proto:** `proto/staff/staff.proto`

- 5 RPCs: CreateStaff, GetStaff, UpdateStaff, DeleteStaff, ListStaff
- `Staff`: StaffID, Name, Role, Contact, Avatar
- Pagination + keyword search

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
  components/     # Footer, HeaderDashboard, KPIGrid, PerformanceTable, Sidebar (shows user info + logout)
  pages/
    Login.tsx            # wired to authApi; validates staff role before granting access
    AnalyticsOverview.tsx
    MenuManagement.tsx
    OrdersManagement.tsx
    StaffManagement.tsx
    WeeklyScheduler.tsx
  services/       # api.ts (injects auth token in all requests), auth.ts
```

---

## Infrastructure

### Docker Compose
**Dev credentials:**
```
POSTGRES_USER=restaurant_user
POSTGRES_PASSWORD=restaurant_pass
POSTGRES_DB=restaurant_db
```

**Active services:** postgres, redis, auth-service, menu-service, staff-service, table-service, order-service, user-service, api-gateway

Binaries are compiled locally and volume-mounted as `./server` inside each Alpine container.

---

## Key Architectural Patterns

1. **Clean Architecture** — domain → repository → usecase → delivery; dependencies point inward only
2. **Per-service go.mod** — each service is its own Go module linked to root via `replace restaurant-management => ../..`
3. **Microservices via gRPC** — inter-service calls use generated gRPC clients; api-gateway clients live in `api-gateway/internal/grpcclient/`
4. **API Gateway** — single HTTP entry point; handles CORS, translates JSON ↔ protobuf
5. **Repository Pattern** — PostgreSQL implementations; `ensureSchema` auto-creates tables + `ALTER TABLE ADD COLUMN IF NOT EXISTS` for safe migrations
6. **JWT + Redis** — stateless access tokens, revocable refresh tokens stored in Redis
7. **Proto-first design** — `.proto` files define all contracts; regenerate with `scripts/generate-proto.sh`
8. **Structured logging** — `go.uber.org/zap` via `shared/pkg/logger`

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

### Known gaps / TODO
- Operating hours (10:00–22:00) enforced only in frontend — no backend validation in order-service
- Admin dashboard (`OrdersManagement.tsx`) may need update to reflect new order fields (`notes`, `end_time`, `user_id`, auto-assigned `table_id`)
- No role-based middleware for menu/staff/table/user routes in api-gateway — only order endpoints have auth
