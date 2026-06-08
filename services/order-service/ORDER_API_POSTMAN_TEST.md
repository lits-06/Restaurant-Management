# Order API Postman Test Script

This document describes HTTP test scripts for API Gateway endpoints that proxy to Order Service.

## Base URL

- `http://localhost:8080`

## Auth Rules

- `GET /orders`: public endpoint.
- `GET /orders/{order_id}`: public endpoint.
- `GET /orders/table/{table_id}`: public endpoint.
- `POST /orders`: requires `Authorization: Bearer <access_token>`.
- `PUT /orders/{order_id}`: requires `Authorization: Bearer <access_token>`.
- `POST /orders/{order_id}/cancel`: requires `Authorization: Bearer <access_token>`.
- `PATCH /orders/{order_id}/status`: requires `Authorization: Bearer <access_token>`.
- `POST /orders/{order_id}/items`: requires `Authorization: Bearer <access_token>`.
- `DELETE /orders/{order_id}/items/{item_id}`: requires `Authorization: Bearer <access_token>`.

## Environment Variables

Create a Postman environment with:

- `base_url` = `http://localhost:8080`
- `access_token` = token from login response
- `table_id` = table id for order creation
- `waiter_id` = waiter user id
- `order_id` = created order id
- `order_item_id` = existing item id in an order
- `menu_item_id_1` = menu item id
- `menu_item_id_2` = another menu item id

## 1) Login to get access token

Method: `POST`
URL: `{{base_url}}/auth/login`

Headers:

- `Content-Type: application/json`

Body:

```json
{
  "email": "admin@restaurant.com",
  "password": "123456"
}
```

Expected response:

- `success = true`
- `access_token` is returned.

## 2) Create order

Method: `POST`
URL: `{{base_url}}/orders`

Headers:

- `Content-Type: application/json`
- `Authorization: Bearer {{access_token}}`

Body:

```json
{
  "table_id": "{{table_id}}",
  "waiter_id": "{{waiter_id}}",
  "items": [
    {
      "menu_item_id": "{{menu_item_id_1}}",
      "quantity": 2,
      "notes": "less spicy"
    },
    {
      "menu_item_id": "{{menu_item_id_2}}",
      "quantity": 1,
      "notes": "no ice"
    }
  ],
  "notes": "customer near window"
}
```

Expected response:

- `success = true`
- `order.order_id` is returned.

Save:

- `order_id = order.order_id`
- `order_item_id = order.items[0].item_id`

## 3) List orders

Method: `GET`
URL: `{{base_url}}/orders?page=1&page_size=10&status=pending`

Optional query params:

- `table_id=<table_id>`
- `status=pending|confirmed|preparing|ready|served|completed|cancelled`
- `from_date=2026-04-01T00:00:00Z`
- `to_date=2026-04-30T23:59:59Z`

## 4) Get order by id

Method: `GET`
URL: `{{base_url}}/orders/{{order_id}}`

## 5) Update order

Method: `PUT`
URL: `{{base_url}}/orders/{{order_id}}`

Headers:

- `Content-Type: application/json`
- `Authorization: Bearer {{access_token}}`

Body:

```json
{
  "items": [
    {
      "menu_item_id": "{{menu_item_id_1}}",
      "quantity": 3,
      "notes": "extra sauce"
    }
  ],
  "notes": "updated notes",
  "discount": 10000
}
```

## 6) Update order status

Method: `PATCH`
URL: `{{base_url}}/orders/{{order_id}}/status`

Headers:

- `Content-Type: application/json`
- `Authorization: Bearer {{access_token}}`

Body:

```json
{
  "status": "confirmed"
}
```

You can test status flow with:

- `confirmed`
- `preparing`
- `ready`
- `served`
- `completed`

## 7) Add order item

Method: `POST`
URL: `{{base_url}}/orders/{{order_id}}/items`

Headers:

- `Content-Type: application/json`
- `Authorization: Bearer {{access_token}}`

Body:

```json
{
  "item": {
    "menu_item_id": "{{menu_item_id_2}}",
    "quantity": 1,
    "notes": "extra hot"
  }
}
```

## 8) Remove order item

Method: `DELETE`
URL: `{{base_url}}/orders/{{order_id}}/items/{{order_item_id}}`

Headers:

- `Authorization: Bearer {{access_token}}`

## 9) Get orders by table

Method: `GET`
URL: `{{base_url}}/orders/table/{{table_id}}?status=pending`

## 10) Cancel order

Method: `POST`
URL: `{{base_url}}/orders/{{order_id}}/cancel`

Headers:

- `Content-Type: application/json`
- `Authorization: Bearer {{access_token}}`

Body:

```json
{
  "reason": "customer requested cancellation"
}
```

## 11) Negative tests

### 11.1 Create order without token

- Request: `POST /orders` without `Authorization`.
- Expected: HTTP `401`.

### 11.2 Invalid status

- Request: `PATCH /orders/{order_id}/status` with body `{ "status": "abc" }`.
- Expected: HTTP `400`.

### 11.3 Invalid date format

- Request: `GET /orders?from_date=04-19-2026`.
- Expected: HTTP `400` with date format error.
