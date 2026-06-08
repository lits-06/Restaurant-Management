# User Service

User management microservice for the restaurant management system.

## Features

- User CRUD operations
- Role management (Admin, Manager, Waiter, Chef, Cashier)
- User status tracking (Active, Inactive, Suspended)
- Email and username uniqueness validation
- Password hashing with bcrypt
- Pagination and filtering

## Architecture

This service follows Clean Architecture principles:

```
cmd/server/          # Application entry point
internal/
  domain/           # Business entities (User)
  repository/       # Data access interfaces & implementations
  usecase/          # Business logic
  delivery/grpc/    # gRPC handlers
pkg/config/         # Configuration
```

## gRPC Service Definition

```protobuf
service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  rpc AssignRole(AssignRoleRequest) returns (AssignRoleResponse);
  rpc GetUserRoles(GetUserRolesRequest) returns (GetUserRolesResponse);
}
```

## Running the Service

### Standalone

```bash
# From project root
go run services/user-service/cmd/server/main.go
```

### With Docker

```bash
# Build
docker build -t user-service -f services/user-service/Dockerfile .

# Run
docker run -p 50052:50052 user-service
```

### With Docker Compose

```bash
# From project root
docker-compose up user-service
```

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| SERVER_PORT | 50052 | gRPC server port |
| LOG_LEVEL | info | Log level (debug, info, warn, error) |

## Testing with grpcurl

### Create User

```bash
grpcurl -plaintext -d '{
  "email": "john@example.com",
  "username": "johndoe",
  "full_name": "John Doe",
  "phone": "+1234567890",
  "password": "password123",
  "roles": ["ROLE_WAITER"]
}' localhost:50052 user.UserService/CreateUser
```

### Get User

```bash
grpcurl -plaintext -d '{
  "user_id": "USER_ID_HERE"
}' localhost:50052 user.UserService/GetUser
```

### List Users

```bash
grpcurl -plaintext -d '{
  "page": 1,
  "page_size": 10
}' localhost:50052 user.UserService/ListUsers
```

### Update User

```bash
grpcurl -plaintext -d '{
  "user_id": "USER_ID_HERE",
  "email": "newemail@example.com",
  "full_name": "John Updated"
}' localhost:50052 user.UserService/UpdateUser
```

### Assign Roles

```bash
grpcurl -plaintext -d '{
  "user_id": "USER_ID_HERE",
  "roles": ["ROLE_ADMIN", "ROLE_MANAGER"]
}' localhost:50052 user.UserService/AssignRole
```

### Get User Roles

```bash
grpcurl -plaintext -d '{
  "user_id": "USER_ID_HERE"
}' localhost:50052 user.UserService/GetUserRoles
```

### Delete User

```bash
grpcurl -plaintext -d '{
  "user_id": "USER_ID_HERE"
}' localhost:50052 user.UserService/DeleteUser
```

## User Roles

- `ROLE_ADMIN` - Full system access
- `ROLE_MANAGER` - Restaurant management
- `ROLE_WAITER` - Order taking and serving
- `ROLE_CHEF` - Kitchen operations
- `ROLE_CASHIER` - Payment processing

## User Status

- `STATUS_ACTIVE` - User can login and use the system
- `STATUS_INACTIVE` - User account is deactivated
- `STATUS_SUSPENDED` - User account is temporarily suspended

## Development

### Adding New Features

1. Update proto definition in `proto/user/user.proto`
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
- golang.org/x/crypto/bcrypt (password hashing)
- github.com/google/uuid (UUID generation)

## Related Services

- **Auth Service** - Authentication and JWT token management
- **Order Service** - Links orders to users (waiters, chefs)
- **Report Service** - User activity reports

---

**Port:** 50052  
**Protocol:** gRPC  
**Status:** ✅ Implemented
