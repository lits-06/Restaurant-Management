# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
COPY services/notification-service/go.mod services/notification-service/go.sum* ./services/notification-service/
COPY shared/go.mod shared/go.sum* ./shared/

# Download dependencies
RUN cd services/notification-service && go mod download

# Copy source code
COPY proto/ ./proto/
COPY shared/ ./shared/
COPY services/notification-service/ ./services/notification-service/

# Build
RUN cd services/notification-service && CGO_ENABLED=0 GOOS=linux go build -o /notification-service ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /notification-service .

EXPOSE 50058

CMD ["./notification-service"]
