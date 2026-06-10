module restaurant-management/services/auth-service

go 1.26.2

require (
	github.com/redis/go-redis/v9 v9.18.0
	go.uber.org/zap v1.27.1
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
	restaurant-management v0.0.0
)

replace restaurant-management => ../..

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
)
