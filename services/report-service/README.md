# Report Service

**Port:** 50059
**Language:** Go 1.21+
**Architecture:** Clean Architecture
**Database:** In-memory (PostgreSQL support planned)

## Overview

Report Service provides comprehensive analytics and reporting across the restaurant management platform. It integrates with multiple services to aggregate data and generate insights.

## Features

- **Sales Reports** - Revenue analytics by period (daily/weekly/monthly/yearly)
- **Inventory Reports** - Stock snapshots and low-stock alerts
- **Order Analytics** - Order completion rates and preparation times
- **Staff Performance** - Metrics for individual staff members or all staff
- **Popular Items** - Top-selling menu items by revenue
- **Revenue Analytics** - Trend analysis and peak hour identification
- **Export Functionality** - PDF, Excel, CSV, JSON formats

## Architecture

### Layers

```
Delivery Layer (gRPC)
    ↓
Use Case Layer (Business Logic)
    ↓
Repository Layer (Data Access)
    ↓
Domain Layer (Entities)
```

### Key Components

#### Domain (`internal/domain/`)
- `report.go` - Report entities (SalesReport, OrderReport, etc.)
- `errors.go` - Domain-specific errors

#### Repository (`internal/repository/`)
- `repository.go` - Interface definitions
- `report_memory.go` - In-memory implementation

#### UseCase (`internal/usecase/`)
- `report_usecase.go` - Business logic (7 operations)
- Service client interfaces for inter-service calls

#### Delivery (`internal/delivery/grpc/`)
- `report_handler.go` - gRPC handler implementations

#### Server (`cmd/server/`)
- `main.go` - Entry point, DI setup, server configuration

## RPC Methods

### 1. GetSalesReport
Generates sales analytics for a date range.

```protobuf
rpc GetSalesReport(GetSalesReportRequest) returns (GetSalesReportResponse);
```

**Response includes:**
- Total sales, tax, discount, net revenue
- Order count and average order value
- Daily breakdown
- Payment method breakdown

### 2. GetInventoryReport
Provides inventory snapshot as of a specific date.

```protobuf
rpc GetInventoryReport(GetInventoryReportRequest) returns (GetInventoryReportResponse);
```

**Response includes:**
- Total inventory value
- Low stock item count
- Individual ingredient stock levels
- Total waste value

### 3. GetOrderReport
Generates order analytics for a date range.

```protobuf
rpc GetOrderReport(GetOrderReportRequest) returns (GetOrderReportResponse);
```

**Response includes:**
- Total, completed, cancelled orders
- Completion rate
- Average preparation time
- Hourly order breakdown

### 4. GetStaffPerformanceReport
Generates staff performance metrics (all or specific staff).

```protobuf
rpc GetStaffPerformanceReport(GetStaffPerformanceReportRequest) returns (GetStaffPerformanceReportResponse);
```

**Response includes:**
- Orders handled
- Total sales and average order value
- Customer complaints count

### 5. GetPopularItemsReport
Identifies top-selling menu items by quantity or revenue.

```protobuf
rpc GetPopularItemsReport(GetPopularItemsReportRequest) returns (GetPopularItemsReportResponse);
```

**Response includes:**
- Item name and ID
- Quantity sold and revenue
- Item category

### 6. GetRevenueAnalytics
Analyzes revenue trends and peak periods.

```protobuf
rpc GetRevenueAnalytics(GetRevenueAnalyticsRequest) returns (GetRevenueAnalyticsResponse);
```

**Response includes:**
- Total revenue and growth rate
- Revenue trend data
- Peak day and hour

### 7. ExportReport
Exports report in specified format (PDF, Excel, CSV, JSON).

```protobuf
rpc ExportReport(ExportReportRequest) returns (ExportReportResponse);
```

## Usage

### Build

```bash
go build -o report-service ./cmd/server
```

### Run Locally

```bash
go run ./cmd/server/main.go
```

### Environment Variables

```env
# gRPC Configuration
REPORT_GRPC_PORT=50059
REPORT_GRPC_HOST=0.0.0.0

# Service Addresses
ORDER_SERVICE_ADDR=localhost:50055
PAYMENT_SERVICE_ADDR=localhost:50056
INVENTORY_SERVICE_ADDR=localhost:50057
USER_SERVICE_ADDR=localhost:50052
MENU_SERVICE_ADDR=localhost:50054

# Logging
LOG_LEVEL=info
ENVIRONMENT=development
```

### Test via gRPC

```bash
# Get sales report
grpcurl -plaintext -d '{
  "period": 2,
  "from_date": {"seconds": 1712188800},
  "to_date": {"seconds": 1712275200}
}' localhost:50059 report.ReportService/GetSalesReport

# Get inventory report
grpcurl -plaintext -d '{}' localhost:50059 report.ReportService/GetInventoryReport

# Get staff performance
grpcurl -plaintext -d '{
  "period": 2,
  "from_date": {"seconds": 1712188800},
  "to_date": {"seconds": 1712275200}
}' localhost:50059 report.ReportService/GetStaffPerformanceReport
```

## Dependencies

### External Services
- **Order Service** - Fetch order data for analytics
- **Payment Service** - Fetch payment and revenue data
- **Inventory Service** - Fetch stock information
- **User Service** - Fetch staff information
- **Menu Service** - Fetch menu item details

### Libraries
- `go.uber.org/zap` - Logging
- `google.golang.org/grpc` - gRPC framework
- `google.golang.org/protobuf` - Protocol buffers

## Data Flow

```
Client Request
    ↓
gRPC Handler (delivery)
    ↓
UseCase (business logic)
    ↓
[Fetch data from other services]
    ↓
Repository (cache/persist)
    ↓
Response
```

## Future Enhancements

- PostgreSQL repository implementation for persistent storage
- Real data aggregation from Order/Payment/Inventory services
- File export functionality (PDF generation, Excel output)
- Advanced filtering and custom date ranges
- Caching layer for frequently accessed reports
- Scheduled report generation and email delivery
- Real-time dashboard updates via WebSocket

## Development Notes

### Adding a New Report Type

1. Add entity to `domain/report.go`
2. Add repository methods to `repository/repository.go`
3. Implement methods in `report_memory.go`
4. Add usecase method to `report_usecase.go`
5. Add handler to `delivery/grpc/report_handler.go`
6. Update proto definitions and regenerate

### Testing Strategy

1. **Domain:** Validate entity methods and errors
2. **Repository:** Test data storage and retrieval
3. **UseCase:** Mock service clients and verify logic
4. **Delivery:** Test request/response mapping

---

**Last Updated:** 2026-04-09
**Status:** Phase 10 Implementation
**Service Number:** 8/9 (Report Service)
