module restaurant-management/api-gateway

go 1.26.2

require (
	go.uber.org/zap v1.27.1
	google.golang.org/grpc v1.80.0
	restaurant-management v0.0.0-00010101000000-000000000000
)

replace restaurant-management => ..

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
