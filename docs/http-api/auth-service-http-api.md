# Auth Service HTTP API

Base URL: `http://localhost:8080`

## 1) Register
- Method: `POST`
- Path: `/auth/register`
- Auth: No

Request body:

```json
{
  "email": "staff1@restaurant.com",
  "password": "123456",
  "username": "staff1",
  "full_name": "Staff One",
  "phone": "0900000001"
}
```

Success response (`200`):

```json
{
  "user_id": "u_123",
  "message": "register successful",
  "success": true
}
```

## 2) Login
- Method: `POST`
- Path: `/auth/login`
- Auth: No

Request body:

```json
{
  "email": "staff1@restaurant.com",
  "password": "123456"
}
```

Success response (`200`):

```json
{
  "access_token": "<jwt_access_token>",
  "refresh_token": "<jwt_refresh_token>",
  "user_id": "u_123",
  "success": true,
  "message": "login successful"
}
```

## 3) Refresh token
- Method: `POST`
- Path: `/auth/refresh`
- Auth: No

Request body:

```json
{
  "refresh_token": "<jwt_refresh_token>"
}
```

Success response (`200`):

```json
{
  "access_token": "<new_jwt_access_token>",
  "success": true,
  "message": "token refreshed"
}
```

## 4) Verify token
- Method: `POST`
- Path: `/auth/verify`
- Auth: No

Request body:

```json
{
  "access_token": "<jwt_access_token>"
}
```

Success response (`200`):

```json
{
  "valid": true,
  "user_id": "u_123",
  "email": "staff1@restaurant.com",
  "roles": ["admin"],
  "expires_at": "2026-04-19T15:30:00Z"
}
```

## 5) Logout
- Method: `POST`
- Path: `/auth/logout`
- Auth: No (token is passed in body)

Request body:

```json
{
  "user_id": "u_123",
  "access_token": "<jwt_access_token>"
}
```

Success response (`200`):

```json
{
  "success": true,
  "message": "logout successful"
}
```

## 6) Change password
- Method: `POST`
- Path: `/auth/change-password`
- Auth: No (identity is passed in body)

Request body:

```json
{
  "user_id": "u_123",
  "old_password": "123456",
  "new_password": "12345678"
}
```

Success response (`200`):

```json
{
  "success": true,
  "message": "password changed"
}
```

## Common errors

```json
{
  "error": "invalid request body",
  "success": false
}
```
