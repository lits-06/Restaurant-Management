# Auth Service Test Script

## 1) Start services

Run from project root:

```bash
docker compose up -d postgres redis auth-service api-gateway
```

Check health:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status":"ok"}
```

## 2) Register account (HTTP JSON via API Gateway)

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "staff1@example.com",
    "password": "123456",
    "username": "staff1",
    "full_name": "Staff One",
    "phone": "0909000111"
  }'
```

## 3) Login and save tokens

```bash
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "staff1@example.com",
    "password": "123456"
  }')

echo "$LOGIN_RESPONSE"

ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token')
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.refresh_token')
USER_ID=$(echo "$LOGIN_RESPONSE" | jq -r '.user_id')

echo "ACCESS_TOKEN=$ACCESS_TOKEN"
echo "REFRESH_TOKEN=$REFRESH_TOKEN"
echo "USER_ID=$USER_ID"
```

## 4) Verify token

```bash
curl -X POST http://localhost:8080/auth/verify \
  -H "Content-Type: application/json" \
  -d "{\"access_token\":\"$ACCESS_TOKEN\"}"
```

## 5) Refresh token

```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}"
```

## 6) Change password

```bash
curl -X POST http://localhost:8080/auth/change-password \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER_ID\",\"old_password\":\"123456\",\"new_password\":\"1234567\"}"
```

## 7) Logout

```bash
curl -X POST http://localhost:8080/auth/logout \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER_ID\",\"access_token\":\"$ACCESS_TOKEN\"}"
```

## 8) Optional: gRPC test directly (without API Gateway)

Requires grpcurl:

```bash
grpcurl -plaintext localhost:50051 list auth.AuthService
```

Login via gRPC:

```bash
grpcurl -plaintext -d '{"email":"staff1@example.com","password":"123456"}' \
  localhost:50051 auth.AuthService/Login
```

## Notes

- This script expects jq installed for token extraction.
- Password policy in auth-service requires at least 6 characters.
- If register returns user already exists, use another email.
