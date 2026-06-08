# Table Service

Table management microservice for the restaurant management system.

## Features

- Table CRUD operations
- Table status management (Available, Occupied, Reserved, Cleaning, Out of Service)
- Location-based table organization
- Capacity tracking
- Available table search with filters
- Table number uniqueness validation

## Architecture

This service follows Clean Architecture principles:

```
cmd/server/          # Application entry point
internal/
  domain/           # Business entities (Table)
  repository/       # Data access interfaces & implementations
  usecase/          # Business logic
  delivery/grpc/    # gRPC handlers
pkg/config/         # Configuration
```

## gRPC Service Definition

```protobuf
service TableService {
  rpc CreateTable(CreateTableRequest) returns (CreateTableResponse);
  rpc GetTable(GetTableRequest) returns (GetTableResponse);
  rpc UpdateTable(UpdateTableRequest) returns (UpdateTableResponse);
  rpc DeleteTable(DeleteTableRequest) returns (DeleteTableResponse);
  rpc ListTables(ListTablesRequest) returns (ListTablesResponse);
  rpc UpdateTableStatus(UpdateTableStatusRequest) returns (UpdateTableStatusResponse);
  rpc GetAvailableTables(GetAvailableTablesRequest) returns (GetAvailableTablesResponse);
}
```

## Running the Service

### Standalone

```bash
# From project root
go run services/table-service/cmd/server/main.go
```

### With Docker

```bash
# Build
docker build -t table-service -f services/table-service/Dockerfile .

# Run
docker run -p 50053:50053 table-service
```

### With Docker Compose

```bash
# From project root
docker-compose up table-service
```

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| SERVER_PORT | 50053 | gRPC server port |
| LOG_LEVEL | info | Log level (debug, info, warn, error) |

## Testing with grpcurl

### Create Table

```bash
grpcurl -plaintext -d '{
  "table_number": "T01",
  "capacity": 4,
  "location": "Main Hall"
}' localhost:50053 table.TableService/CreateTable
```

### Get Table

```bash
grpcurl -plaintext -d '{
  "table_id": "TABLE_ID_HERE"
}' localhost:50053 table.TableService/GetTable
```

### List Tables

```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10,
  "location": "Main Hall"
}' localhost:50053 table.TableService/ListTables
```

### Update Table

```bash
grpcurl -plaintext -d '{
  "table_id": "TABLE_ID_HERE",
  "table_number": "T01-Updated",
  "capacity": 6
}' localhost:50053 table.TableService/UpdateTable
```

### Update Table Status

```bash
# Mark as occupied
grpcurl -plaintext -d '{
  "table_id": "TABLE_ID_HERE",
  "status": "STATUS_OCCUPIED",
  "order_id": "ORDER_123"
}' localhost:50053 table.TableService/UpdateTableStatus

# Mark as available
grpcurl -plaintext -d '{
  "table_id": "TABLE_ID_HERE",
  "status": "STATUS_AVAILABLE"
}' localhost:50053 table.TableService/UpdateTableStatus

# Mark as reserved
grpcurl -plaintext -d '{
  "table_id": "TABLE_ID_HERE",
  "status": "STATUS_RESERVED"
}' localhost:50053 table.TableService/UpdateTableStatus
```

### Get Available Tables

```bash
# All available tables
grpcurl -plaintext -d '{}' localhost:50053 table.TableService/GetAvailableTables

# Available tables with minimum capacity
grpcurl -plaintext -d '{
  "min_capacity": 4,
  "location": "Main Hall"
}' localhost:50053 table.TableService/GetAvailableTables
```

### Delete Table

```bash
grpcurl -plaintext -d '{
  "table_id": "TABLE_ID_HERE"
}' localhost:50053 table.TableService/DeleteTable
```

## Table Statuses

- `STATUS_AVAILABLE` - Table is ready for customers
- `STATUS_OCCUPIED` - Table currently has customers
- `STATUS_RESERVED` - Table is reserved for upcoming customers
- `STATUS_CLEANING` - Table is being cleaned after customers leave
- `STATUS_OUT_OF_SERVICE` - Table is temporarily unavailable

## Status Transitions

```
AVAILABLE → OCCUPIED (when customers are seated)
AVAILABLE → RESERVED (when table is reserved)
RESERVED → OCCUPIED (when reserved customers arrive)
OCCUPIED → CLEANING (when customers leave)
CLEANING → AVAILABLE (when cleaning is done)
ANY → OUT_OF_SERVICE (except OCCUPIED)
OUT_OF_SERVICE → AVAILABLE (when back in service)
```

## Business Rules

1. **Table Number Uniqueness**: Each table must have a unique table number
2. **Capacity Limits**: Tables can hold 1-50 people
3. **Deletion Restriction**: Cannot delete occupied tables
4. **Status Changes**: Some status transitions are restricted for data integrity
5. **Location Organization**: Tables can be organized by location (e.g., Main Hall, VIP Room, Outdoor)

## Development

### Adding New Features

1. Update proto definition in `proto/table/table.proto`
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

- **Order Service** - Links orders to tables
- **User Service** - Waiter assignment to tables
- **Menu Service** - Orders placed at tables

## Use Cases

1. **Restaurant Floor Management**: Track table availability in real-time
2. **Reservation System**: Reserve tables for upcoming customers
3. **Capacity Planning**: Find suitable tables for party sizes
4. **Location-based Organization**: Manage different dining areas

---

**Port:** 50053  
**Protocol:** gRPC  
**Status:** ✅ Implemented
