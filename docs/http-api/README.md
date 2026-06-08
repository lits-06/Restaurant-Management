# HTTP API Docs for Front-end

This folder contains HTTP API documentation grouped by service.

## Base URL
- Local: `http://localhost:8080`

## Service docs
- [Auth Service](./auth-service-http-api.md)
- [Menu Service](./menu-service-http-api.md)
- [Table Service](./table-service-http-api.md)
- [Order Service](./order-service-http-api.md)
- [User Service](./user-service-http-api.md)
- [Payment Service](./payment-service-http-api.md)
- [Inventory Service](./inventory-service-http-api.md)
- [Notification Service](./notification-service-http-api.md)
- [Report Service](./report-service-http-api.md)

## Current HTTP availability
- Available now (via API Gateway): Auth, Menu, Table, Order.
- Not exposed via HTTP yet: User, Payment, Inventory, Notification, Report.

## Common error response
Most API Gateway handlers return this structure on error:

```json
{
  "error": "message",
  "success": false
}
```

## Health check
- `GET /health`
- Response:

```json
{
  "status": "ok"
}
```
