# Restaurant Management System

Hệ thống quản lý nhà hàng full-stack với kiến trúc microservices. Backend Go (Clean Architecture + gRPC), 3 React frontends, giao tiếp qua HTTP/REST API Gateway.

---

## Tổng quan Kiến trúc

```
┌─────────────────────────────────────────────────────────────┐
│                      3 Frontend Apps                         │
│  restaurant-app    restaurant-app-admin   restaurant-app-   │
│     :5173              :5174              kitchen :5175      │
│   (Khách hàng)       (Admin/Staff)      (Bếp/Phục vụ)      │
└─────────────────────────────────────────────────────────────┘
                              │ HTTP/REST + WebSocket
                              ▼
                    ┌─────────────────┐
                    │   API Gateway   │ :8080
                    │  (HTTP → gRPC)  │
                    └────────┬────────┘
                             │ gRPC
          ┌──────────────────┼──────────────────┐
          │          │       │        │          │
    auth :50051  menu :50054 │  staff :50052   user :50056
                     │  order :50055   table :50053
                     │  notification :50058
                     │  report :50059
                             │
              ┌──────────────┴──────────────┐
              │                             │
         PostgreSQL :5432              Redis :6379
          (restaurant_db)          (tokens + pub/sub)
```

---

## Cấu trúc Repository

```
Restaurant_Management/
├── restaurant-app/            # Customer SPA (port 5173)
├── restaurant-app-admin/      # Admin dashboard (port 5174)
├── restaurant-app-kitchen/    # Kitchen display (port 5175)
├── api-gateway/               # HTTP → gRPC layer (port 8080)
├── services/
│   ├── auth-service/          # Authentication, JWT, Redis (port 50051)
│   ├── schedule-service/      # Shift scheduling for kitchen staff (port 50052)
│   ├── table-service/         # Table physical registry (port 50053)
│   ├── menu-service/          # Menu items + categories (port 50054)
│   ├── order-service/         # Orders, auto-assign, item status (port 50055)
│   ├── user-service/          # Users, roles, credentials (port 50056)
│   ├── notification-service/  # Redis Pub/Sub, kitchen push (port 50058)
│   └── report-service/        # Analytics (port 50059, not yet wired)
├── proto/                     # .proto source files (source of truth)
├── shared/                    # Shared Go packages (logger, db, jwt, errors)
├── scripts/                   # Proto generation scripts
├── docker-compose.yml
└── go.mod / go.sum            # Root Go module
```

---

## Port Map

| Component | Protocol | Port |
|-----------|----------|------|
| restaurant-app | HTTP | 5173 |
| restaurant-app-admin | HTTP | 5174 |
| restaurant-app-kitchen | HTTP | 5175 |
| api-gateway | HTTP/REST + WebSocket | 8080 |
| auth-service | gRPC | 50051 |
| schedule-service | gRPC | 50052 |
| table-service | gRPC | 50053 |
| menu-service | gRPC | 50054 |
| order-service | gRPC | 50055 |
| user-service | gRPC | 50056 |
| notification-service | gRPC | 50058 |
| report-service | gRPC | 50059 |
| PostgreSQL | TCP | 5432 |
| Redis | TCP | 6379 |

---

## 5 Roles

| Role | Mô tả | App |
|------|-------|-----|
| `USER` | Khách hàng — đặt bàn, xem order của bản thân | restaurant-app |
| `MANAGER` | Tạo order walk-in, cập nhật trạng thái | restaurant-app-admin |
| `CHEF` | Xem món, mark COOKING/READY | restaurant-app-kitchen |
| `WAITER` | Nhận thông báo, mark SERVED | restaurant-app-kitchen |
| `ADMIN` | Toàn quyền — tất cả apps | Tất cả |

---

## Luồng Chính

### Khách đặt bàn (Customer Flow)
```
1. Vào ReservationPage → chọn ngày/giờ/số khách/món ăn
2. POST /orders → order-service
3. Auto-assign table (nếu không chỉ định bàn cụ thể)
4. Order được tạo với status: Pending
5. Khách nhận confirmation (hiện tại chưa có email)
```

### Luồng bếp (Kitchen Flow)
```
1. MANAGER/ADMIN: PATCH /orders/{id}/status → Confirmed
2. → notification-service publish ORDER_CONFIRMED → Redis → WebSocket → KitchenPage
3. CHEF: mark item COOKING, sau đó READY
4. → notification-service publish ITEM_READY → Redis → WebSocket → WaiterPage
5. WAITER: mark item SERVED (fetch đến bàn)
6. MANAGER: PATCH /orders/{id}/status → Completed (sau thanh toán)
```

---

## Chạy hệ thống

### Docker Compose (Full Stack)

```bash
# Build và chạy tất cả services
docker-compose up --build

# Chạy không build lại
docker-compose up

# Dừng
docker-compose down
```

**Dev credentials PostgreSQL:**
```
POSTGRES_USER=restaurant_user
POSTGRES_PASSWORD=restaurant_pass
POSTGRES_DB=restaurant_db
```

### Chạy từng service riêng (Development)

```bash
# Backend service
cd services/order-service
go run cmd/server/main.go

# Frontend
cd restaurant-app
npm install && npm run dev
```

### Build Go binaries (cho Docker)

```bash
# Mỗi service build trong directory của nó
cd services/auth-service && go build -o ../../server ./cmd/server/
cd services/order-service && go build -o ../../server ./cmd/server/
# ... tương tự cho các service khác
cd api-gateway && go build -o server ./cmd/server/
```

### Regenerate Proto

```bash
export PATH=$PATH:~/go/bin
bash scripts/generate-proto.sh
```

---

## Service Documentation

| Service | README |
|---------|--------|
| API Gateway | [api-gateway/README.md](api-gateway/README.md) |
| Auth Service | [services/auth-service/README.md](services/auth-service/README.md) |
| User Service | [services/user-service/README.md](services/user-service/README.md) |
| Order Service | [services/order-service/README.md](services/order-service/README.md) |
| Table Service | [services/table-service/README.md](services/table-service/README.md) |
| Menu Service | [services/menu-service/README.md](services/menu-service/README.md) |
| Schedule Service | [services/schedule-service/README.md](services/schedule-service/README.md) |
| Notification Service | [services/notification-service/README.md](services/notification-service/README.md) |
| Report Service | [services/report-service/README.md](services/report-service/README.md) |
| Customer App | [restaurant-app/README.md](restaurant-app/README.md) |
| Admin Dashboard | [restaurant-app-admin/README.md](restaurant-app-admin/README.md) |
| Kitchen Display | [restaurant-app-kitchen/README.md](restaurant-app-kitchen/README.md) |

---

## Tech Stack

### Backend
- **Go 1.25+** — mỗi service có `go.mod` riêng
- **gRPC** (`google.golang.org/grpc v1.80`) + **Protocol Buffers v3**
- **PostgreSQL 15** — single shared `restaurant_db`
- **Redis 7** — refresh tokens (auth-service) + Pub/Sub (notification-service)
- **go.uber.org/zap** — structured logging
- **github.com/golang-jwt/jwt/v5** — JWT HS256

### Frontend
- **React 19** + **Vite** + **TypeScript**
- **TailwindCSS** — styling
- **Zustand** — state management + persist
- **Native `fetch` / WebSocket API** — HTTP + real-time

---

## Vấn đề chưa giải quyết / Roadmap

### Quan trọng
- [ ] Role-based middleware cho menu/table/user routes trong api-gateway
- [ ] `table_id` UUID → `table_number` lookup trong kitchen app
- [ ] Token refresh flow trong tất cả frontends
- [ ] OrdersManagement (admin) cần hiển thị `notes`, `end_time`, `item_status`
- [ ] WebSocket auto-reconnect trong kitchen app

### Business Logic
- [ ] Operating hours validation ở backend (hiện chỉ có frontend)
- [ ] Notification khi order bị cancel (CHEF đang nấu)
- [ ] Concurrency control cho auto-assign table (race condition)
- [ ] Email/SMS confirmation sau khi đặt bàn
- [ ] Xung đột ca làm việc — schedule-service chưa validate trùng shift cùng user cùng ngày giờ

### Infrastructure
- [ ] Report service kết nối PostgreSQL và integrate vào API Gateway
- [ ] Rate limiting tại API Gateway
- [ ] Circuit breaker cho inter-service gRPC calls
- [ ] Notification persistence (hiện là fire-and-forget)
