# Restaurant Management Microservices System Plan

## Problem Statement
Tạo một hệ thống quản lý nhà hàng theo kiến trúc microservices sử dụng Golang, Clean Architecture và gRPC để giao tiếp giữa các services.

## Architecture Overview

### Services
1. **Auth Service**: JWT + Refresh Token authentication
2. **User Service**: Quản lý user (staff, manager, admin)
3. **Table Service**: Quản lý bàn ăn
4. **Menu Service**: Quản lý thực đơn và món ăn
5. **Order Service**: Quản lý đơn hàng
6. **Payment Service**: Xử lý thanh toán
7. **Inventory Service**: Quản lý kho nguyên liệu
8. **Notification Service**: Gửi thông báo
9. **Report Service**: Báo cáo và thống kê

### Technology Stack
- **Language**: Golang
- **Communication**: gRPC
- **Architecture**: Clean Architecture (Entity, Use Case, Controller, Framework layers)
- **Database**: Interface-based repository pattern (DB implementation agnostic)
- **Infrastructure**: Docker, docker-compose, API Gateway

### Clean Architecture Structure (per service)
```
service-name/
├── cmd/                    # Application entry points
│   └── server/
├── internal/
│   ├── domain/            # Entities (Business models)
│   ├── usecase/           # Business logic / Use cases
│   ├── delivery/          # Controllers (gRPC handlers)
│   │   └── grpc/
│   └── repository/        # Data access layer (interfaces + implementations)
├── pkg/                   # Shared packages
│   ├── logger/
│   ├── config/
│   └── middleware/
└── proto/                 # Protocol buffer definitions
```

## Implementation Todos

### Phase 1: Project Setup & Shared Infrastructure
- setup-project-structure: Tạo cấu trúc thư mục cho toàn bộ project
- setup-proto-definitions: Tạo proto files cho tất cả services
- setup-shared-packages: Tạo shared packages (logger, config, middleware, jwt util)
- setup-infrastructure: Tạo Docker, docker-compose, Makefile
- setup-api-gateway: Tạo API Gateway để routing HTTP -> gRPC

### Phase 2: Auth Service
- auth-domain-layer: Tạo domain entities cho Auth (User credentials, tokens)
- auth-repository-layer: Tạo repository interfaces và mock implementation
- auth-usecase-layer: Implement business logic (login, register, refresh token, verify token)
- auth-grpc-delivery: Implement gRPC handlers cho Auth Service
- auth-server: Tạo main server file và configuration

### Phase 3: User Service
- user-domain-layer: Tạo domain entities (User, Role, Profile)
- user-repository-layer: Tạo repository interfaces
- user-usecase-layer: Implement CRUD operations, role management
- user-grpc-delivery: Implement gRPC handlers
- user-server: Tạo main server file

### Phase 4: Table Service
- table-domain-layer: Tạo domain entities (Table, TableStatus)
- table-repository-layer: Tạo repository interfaces
- table-usecase-layer: Implement table management (create, update, status tracking)
- table-grpc-delivery: Implement gRPC handlers
- table-server: Tạo main server file

### Phase 5: Menu Service
- menu-domain-layer: Tạo domain entities (MenuItem, Category, Price)
- menu-repository-layer: Tạo repository interfaces
- menu-usecase-layer: Implement menu management (CRUD, category management)
- menu-grpc-delivery: Implement gRPC handlers
- menu-server: Tạo main server file

### Phase 6: Order Service
- order-domain-layer: Tạo domain entities (Order, OrderItem, OrderStatus)
- order-repository-layer: Tạo repository interfaces
- order-usecase-layer: Implement order workflow (create, update, status transitions)
- order-grpc-client: Integrate với Menu, Table, User services
- order-grpc-delivery: Implement gRPC handlers
- order-server: Tạo main server file

### Phase 7: Payment Service
- payment-domain-layer: Tạo domain entities (Payment, Transaction, PaymentMethod)
- payment-repository-layer: Tạo repository interfaces
- payment-usecase-layer: Implement payment processing logic
- payment-grpc-client: Integrate với Order service
- payment-grpc-delivery: Implement gRPC handlers
- payment-server: Tạo main server file

### Phase 8: Inventory Service
- inventory-domain-layer: Tạo domain entities (Ingredient, Stock, StockMovement)
- inventory-repository-layer: Tạo repository interfaces
- inventory-usecase-layer: Implement inventory management (CRUD, stock tracking)
- inventory-grpc-client: Integrate với Order service
- inventory-grpc-delivery: Implement gRPC handlers
- inventory-server: Tạo main server file

### Phase 9: Notification Service
- notification-domain-layer: Tạo domain entities (Notification, NotificationType)
- notification-repository-layer: Tạo repository interfaces
- notification-usecase-layer: Implement notification logic (email, SMS, push)
- notification-grpc-delivery: Implement gRPC handlers
- notification-server: Tạo main server file

### Phase 10: Report Service
- report-domain-layer: Tạo domain entities (Report, ReportType, Analytics)
- report-repository-layer: Tạo repository interfaces
- report-usecase-layer: Implement reporting logic (sales, inventory, performance)
- report-grpc-client: Integrate với Order, Payment, Inventory services
- report-grpc-delivery: Implement gRPC handlers
- report-server: Tạo main server file

### Phase 11: Integration & Testing
- integration-docker: Finalize Docker configurations cho tất cả services
- integration-compose: Configure docker-compose với tất cả services
- integration-gateway: Complete API Gateway implementation
- create-makefile: Tạo Makefile với commands để build, run, test
- create-readme: Tạo documentation cho project

## Notes

### Clean Architecture Principles
- **Independence**: Business logic không phụ thuộc vào framework, UI, database
- **Testability**: Use cases có thể test riêng biệt
- **Dependency Rule**: Dependencies chỉ point inward (Framework -> Controller -> Use Case -> Entity)

### gRPC Communication Pattern
- Services giao tiếp với nhau qua gRPC
- API Gateway nhận HTTP requests và forward đến các gRPC services
- Use protobuf cho data serialization

### Repository Pattern
- Repository interfaces defined trong domain/usecase layer
- Implementations trong repository layer
- Dễ dàng swap database implementations (PostgreSQL, MongoDB, Redis, etc.)

### Security Considerations
- Auth Service generates JWT tokens
- API Gateway validates JWT trước khi forward requests
- Services có thể verify JWT nếu cần (optional)
- Sensitive data cần được encrypt

### Service Communication Example
Order Service -> Menu Service: Get menu item details
Order Service -> Table Service: Check table availability
Order Service -> Inventory Service: Update stock after order
Order Service -> Payment Service: Process payment
Payment Service -> Notification Service: Send payment confirmation
