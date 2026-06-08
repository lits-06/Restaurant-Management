FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY services/order-service/go.mod services/order-service/
COPY shared/ shared/
COPY proto/ proto/

# Download dependencies
RUN go mod download

# Copy source code
COPY services/order-service/ services/order-service/

# Build
WORKDIR /app/services/order-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /order-service ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /order-service .

EXPOSE 50055

CMD ["./order-service"]
