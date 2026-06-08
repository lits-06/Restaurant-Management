# Protobuf Generation

This project uses Protocol Buffers and gRPC for service communication.

## Prerequisites

### Install Protocol Buffer Compiler

Check installation:

```bash
protoc --version
```

### Install Go Plugins

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Make sure Go binaries are available in PATH:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Verify:

```bash
which protoc-gen-go
which protoc-gen-go-grpc
```

---

## Directory Structure

```text
proto/
├── auth/
│   └── auth.proto
├── user/
│   └── user.proto
└── product/
    └── product.proto
```

---

## Generate Protobuf Files

Grant execution permission:

```bash
chmod +x scripts/gen_proto.sh
```

Run:

```bash
./scripts/gen_proto.sh
```

Or:

```bash
bash scripts/gen_proto.sh
```

---

## Generated Files

For each `.proto` file, the following files will be generated:

```text
*.pb.go
*_grpc.pb.go
```

Example:

```text
proto/auth/auth.pb.go
proto/auth/auth_grpc.pb.go
```

---

## Regenerate After Changes

Whenever a `.proto` file is modified, rerun:

```bash
./scripts/gen_proto.sh
```
