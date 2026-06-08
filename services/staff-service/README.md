# Staff Service

Staff management microservice for the restaurant management system.

## Features

- Staff CRUD operations
- PostgreSQL persistence
- Pagination and keyword filtering
- Staff fields based on the admin UI

## Architecture

This service follows Clean Architecture principles:

```
cmd/server/          # Application entry point
internal/
  domain/           # Business entities (Staff)
  repository/       # Data access interfaces & implementations
  usecase/          # Business logic
  delivery/grpc/    # gRPC handlers
pkg/config/         # Configuration
```

## gRPC Service Definition

```protobuf
service StaffService {
  rpc CreateStaff(CreateStaffRequest) returns (CreateStaffResponse);
  rpc GetStaff(GetStaffRequest) returns (GetStaffResponse);
  rpc UpdateStaff(UpdateStaffRequest) returns (UpdateStaffResponse);
  rpc DeleteStaff(DeleteStaffRequest) returns (DeleteStaffResponse);
  rpc ListStaff(ListStaffRequest) returns (ListStaffResponse);
}
```

## Running the Service

### Standalone

```bash
# From project root
go run services/staff-service/cmd/server/main.go
```

### With Docker

```bash
# Build
docker build -t staff-service -f services/staff-service/Dockerfile .

# Run
docker run -p 50053:50053 staff-service
```

### With Docker Compose

```bash
# From project root
docker-compose up staff-service
```

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| GRPC_PORT | 50053 | gRPC server port |
| DATABASE_HOST | localhost | PostgreSQL host |
| DATABASE_PORT | 5432 | PostgreSQL port |
| DATABASE_USER | restaurant_user | PostgreSQL user |
| DATABASE_PASSWORD | restaurant_pass | PostgreSQL password |
| DATABASE_NAME | restaurant_db | PostgreSQL database |
| DATABASE_SSLMODE | disable | PostgreSQL SSL mode |
| LOG_LEVEL | info | Log level (debug, info, warn, error) |

## Testing with grpcurl

### Create Staff

```bash
grpcurl -plaintext -d '{
  "name": "Marcus Vance",
  "role": "Floor Manager",
  "contact": "+1 (555) 012-3456",
  "until": "11:00 PM",
  "avatar": "https://example.com/avatar.png"
}' localhost:50053 staff.StaffService/CreateStaff
```

### Get Staff

```bash
grpcurl -plaintext -d '{
  "staff_id": "STAFF_ID_HERE"
}' localhost:50053 staff.StaffService/GetStaff
```

### List Staff

```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10,
  "keyword": "server"
}' localhost:50053 staff.StaffService/ListStaff
```

### Update Staff

```bash
grpcurl -plaintext -d '{
  "staff_id": "STAFF_ID_HERE",
  "name": "Marcus Vance",
  "role": "Floor Manager",
  "contact": "+1 (555) 000-0000",
  "until": "12:00 AM"
}' localhost:50053 staff.StaffService/UpdateStaff
```

### Delete Staff

```bash
grpcurl -plaintext -d '{
  "staff_id": "STAFF_ID_HERE"
}' localhost:50053 staff.StaffService/DeleteStaff
```

## Development

### Adding New Features

1. Update proto definition in `proto/staff/staff.proto`
2. Regenerate proto code: `make proto`
3. Add domain logic in `internal/domain/`
4. Update repository interface if needed
5. Implement use case in `internal/usecase/`
6. Add gRPC handler in `internal/delivery/grpc/`

### Running Tests

```bash
go test ./...
```

## Dependencies

- gRPC & Protocol Buffers
- go.uber.org/zap (logging)
- github.com/google/uuid (UUID generation)

## Related Services

- **Auth Service** - Authentication and JWT token management
- **Order Service** - Links orders to staff members
- **Report Service** - Staff activity reports

---

**Port:** 50053  
**Protocol:** gRPC  
**Status:** ✅ Implemented
