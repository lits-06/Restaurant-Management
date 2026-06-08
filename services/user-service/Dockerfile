# Builder stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY services/user-service/go.mod services/user-service/go.sum ./services/user-service/
COPY shared/go.mod ./shared/

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
WORKDIR /app/services/user-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/user-service ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/user-service .

# Expose port
EXPOSE 50052

# Run the binary
CMD ["./user-service"]
