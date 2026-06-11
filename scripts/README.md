# Scripts

Build, proto generation, and Docker orchestration helpers for the Restaurant Management System.

---

## Makefile

Run from the `scripts/` directory.

### Full stack

```bash
make start       # docker compose up -d (all 8 app services)
make stop        # docker compose stop
make restart     # rebuild all binaries then stop + start
make logs        # tail logs for all app services
make db          # open psql shell in the postgres container
```

### Build all binaries

```bash
make rebuild     # compile all 8 services (static Linux binaries)
```

Each service binary is written to `<service-dir>/server` and volume-mounted into Docker.

### Build + restart a single service

```bash
make restart-gateway
make restart-auth
make restart-menu
make restart-schedule
make restart-table
make restart-order
make restart-notification
make restart-user
```

### Frontend dev servers

```bash
make fe-customer   # port 5173
make fe-admin      # port 5174
make fe-kitchen    # port 5175
```

### Proto generation

```bash
make proto         # runs generate-proto.sh with ~/go/bin in PATH
```

---

## Proto Generation

### Prerequisites

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Regenerate all protos

```bash
# from repo root
bash scripts/generate-proto.sh

# or via Makefile (from scripts/)
make proto
```

### Regenerate a single proto

```bash
# from repo root
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/<service>/<file>.proto
```

### Active proto files

```
proto/
├── auth/         auth.proto
├── menu/         menu.proto
├── notification/ notification.proto
├── order/        order.proto
├── report/       report.proto
├── schedule/     schedule.proto
├── table/        table.proto
└── user/         user.proto
```

Each `.proto` generates two files: `*.pb.go` and `*_grpc.pb.go`. Do not edit generated files directly.

---

## Windows

Use `generate-proto.bat` instead of `generate-proto.sh` on Windows (same logic, batch syntax).
