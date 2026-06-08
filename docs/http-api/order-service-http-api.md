# Order Service HTTP API

Base URL: `http://localhost:8080`

## Auth rules
- Public: `GET /orders`, `GET /orders/{order_id}`, `GET /orders/table/{table_id}`
- Requires Bearer token: `POST /orders`, `PUT /orders/{order_id}`, `POST /orders/{order_id}/cancel`, `PATCH /orders/{order_id}/status`, `POST /orders/{order_id}/items`, `DELETE /orders/{order_id}/items/{item_id}`

## 1) List orders
- Method: `GET`
- Path: `/orders`
- Auth: No

Query params:
- `page` (default `1`)
- `page_size` (default `20`)
- `status` (optional): `pending`, `confirmed`, `preparing`, `ready`, `served`, `completed`, `cancelled`
- `table_id` (optional)
- `from_date` (optional, RFC3339)
- `to_date` (optional, RFC3339)

Example:
- `/orders?page=1&page_size=10&status=pending&from_date=2026-04-19T00:00:00Z&to_date=2026-04-19T23:59:59Z`

Success response (`200`):

```json
{
  "orders": [
    {
      "order_id": "ord_1001",
      "table_id": "tb_001",
      "table_number": "T01",
      "waiter_id": "u_001",
      "waiter_name": "Staff One",
      "items": [
        {
          "item_id": "oi_01",
          "menu_item_id": "mi_001",
          "menu_item_name": "Grilled Chicken Rice",
          "quantity": 2,
          "unit_price": 75000,
          "subtotal": 150000,
          "notes": "less spicy",
          "status": "STATUS_PENDING"
        }
      ],
      "subtotal": 150000,
      "tax": 15000,
      "discount": 0,
      "total": 165000,
      "status": "STATUS_PENDING",
      "notes": "customer near window",
      "created_at": "2026-04-19T10:00:00Z",
      "updated_at": "2026-04-19T10:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "success": true,
  "message": "ok"
}
```

## 2) Create order
- Method: `POST`
- Path: `/orders`
- Auth: `Authorization: Bearer <access_token>`

Request body:

```json
{
  "table_id": "tb_001",
  "waiter_id": "u_001",
  "items": [
    {
      "menu_item_id": "mi_001",
      "quantity": 2,
      "notes": "less spicy"
    }
  ],
  "notes": "customer near window"
}
```

Success response (`200`):

```json
{
  "order": {
    "order_id": "ord_1001",
    "status": "STATUS_PENDING"
  },
  "success": true,
  "message": "created"
}
```

## 3) Get order detail
- Method: `GET`
- Path: `/orders/{order_id}`
- Auth: No

Success response (`200`):

```json
{
  "order": {
    "order_id": "ord_1001",
    "table_id": "tb_001",
    "status": "STATUS_PENDING"
  },
  "success": true,
  "message": "ok"
}
```

## 4) Update order
- Method: `PUT`
- Path: `/orders/{order_id}`
- Auth: `Authorization: Bearer <access_token>`

Request body:

```json
{
  "items": [
    {
      "menu_item_id": "mi_001",
      "quantity": 3,
      "notes": "extra sauce"
    }
  ],
  "notes": "updated notes",
  "discount": 10000
}
```

Success response (`200`):

```json
{
  "order": {
    "order_id": "ord_1001",
    "discount": 10000,
    "status": "STATUS_PENDING"
  },
  "success": true,
  "message": "updated"
}
```

## 5) Cancel order
- Method: `POST`
- Path: `/orders/{order_id}/cancel`
- Auth: `Authorization: Bearer <access_token>`

Request body:

```json
{
  "reason": "customer requested cancellation"
}
```

Success response (`200`):

```json
{
  "success": true,
  "message": "cancelled"
}
```

## 6) Update order status
- Method: `PATCH`
- Path: `/orders/{order_id}/status`
- Auth: `Authorization: Bearer <access_token>`

Request body:

```json
{
  "status": "confirmed"
}
```

Status values:
- `pending`
- `confirmed`
- `preparing`
- `ready`
- `served`
- `completed`
- `cancelled`

Success response (`200`):

```json
{
  "order": {
    "order_id": "ord_1001",
    "status": "STATUS_CONFIRMED"
  },
  "success": true,
  "message": "status updated"
}
```

## 7) Add order item
- Method: `POST`
- Path: `/orders/{order_id}/items`
- Auth: `Authorization: Bearer <access_token>`

Request body:

```json
{
  "item": {
    "menu_item_id": "mi_002",
    "quantity": 1,
    "notes": "extra hot"
  }
}
```

Success response (`200`):

```json
{
  "order": {
    "order_id": "ord_1001",
    "items": []
  },
  "success": true,
  "message": "item added"
}
```

## 8) Remove order item
- Method: `DELETE`
- Path: `/orders/{order_id}/items/{item_id}`
- Auth: `Authorization: Bearer <access_token>`

Success response (`200`):

```json
{
  "order": {
    "order_id": "ord_1001",
    "items": []
  },
  "success": true,
  "message": "item removed"
}
```

## 9) Get orders by table
- Method: `GET`
- Path: `/orders/table/{table_id}`
- Auth: No

Query params:
- `status` (optional)

Success response (`200`):

```json
{
  "orders": [
    {
      "order_id": "ord_1001",
      "table_id": "tb_001",
      "status": "STATUS_PENDING"
    }
  ],
  "success": true,
  "message": "ok"
}
```

## Common errors

```json
{
  "error": "invalid from_date, expected RFC3339",
  "success": false
}
```
