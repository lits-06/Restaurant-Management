FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY shared/ ./shared/
COPY proto/ ./proto/
COPY services/auth-service/ ./services/auth-service/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./services/auth-service/cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 50051

CMD ["./main"]
