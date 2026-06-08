# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
COPY services/payment-service/go.mod services/payment-service/go.sum* ./services/payment-service/
COPY shared/go.mod shared/go.sum* ./shared/

# Download dependencies
RUN cd services/payment-service && go mod download

# Copy source code
COPY proto/ ./proto/
COPY shared/ ./shared/
COPY services/payment-service/ ./services/payment-service/

# Build
RUN cd services/payment-service && CGO_ENABLED=0 GOOS=linux go build -o /payment-service ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /payment-service .

EXPOSE 50056

CMD ["./payment-service"]
