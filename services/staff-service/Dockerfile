# Builder stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY services/staff-service/go.mod services/staff-service/go.sum ./services/staff-service/
COPY shared/go.mod ./shared/

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
WORKDIR /app/services/staff-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/staff-service ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/staff-service .

# Expose port
EXPOSE 50052

# Run the binary
CMD ["./staff-service"]
