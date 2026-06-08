# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api-gateway ./api-gateway/cmd/server

# Runtime stage
FROM alpine:3.21

WORKDIR /app
COPY --from=builder /bin/api-gateway /app/api-gateway

EXPOSE 8080

CMD ["/app/api-gateway"]
