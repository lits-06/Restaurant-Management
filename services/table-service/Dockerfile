# Builder stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
COPY services/table-service/go.mod services/table-service/go.sum ./services/table-service/
COPY shared/go.mod ./shared/

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
WORKDIR /app/services/table-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/table-service ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/table-service .

# Expose port
EXPOSE 50053

# Run the binary
CMD ["./table-service"]
