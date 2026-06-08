FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY services/menu-service/go.mod services/menu-service/
COPY shared/ shared/
COPY proto/ proto/

# Download dependencies
RUN go mod download

# Copy source code
COPY services/menu-service/ services/menu-service/

# Build
WORKDIR /app/services/menu-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /menu-service ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /menu-service .

EXPOSE 50054

CMD ["./menu-service"]
