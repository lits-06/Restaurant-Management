# Table Service HTTP API

Base URL: `http://localhost:8080`

## Auth rules
- Public: `GET /tables`, `GET /tables/{table_id}`, `GET /tables/available`
- Requires Bearer token: `POST /tables`, `PUT /tables/{table_id}`, `DELETE /tables/{table_id}`, `PATCH /tables/{table_id}/status`
- Requires Bearer token: `POST /tables/{table_id}/reservations`, `GET /tables/{table_id}/reservations`, `GET /reservations/{reservation_id}`, `POST /reservations/{reservation_id}/cancel`

## 1) List tables
- Method: `GET`
- Path: `/tables`
- Auth: No

Query params:
- `page` (default `1`)
- `page_size` (default `10`)
- `status` (optional): `available`, `occupied`, `reserved`, `cleaning`, `out_of_service`
- `location` (optional)

Example:
- `/tables?page=1&page_size=10&status=available&location=Main`

Success response (`200`):

```json
{
  "tables": [
    {
      "table_id": "tb_001",
      "table_number": "T01",
      "capacity": 4,
      "status": "STATUS_AVAILABLE",
      "location": "Main Hall",
      "current_order_id": ""
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "success": true,
  "message": "ok"
}
```

## 2) Create table
- Method: `POST`
- Path: `/tables`
- Auth: `Authorization: Bearer <access_token>`

Request body:

```json
{
  "table_number": "T01",
  "capacity": 4,
  "location": "Main Hall"
}
```

Success response (`200`):

```json
{
  "table": {
    "table_id": "tb_001",
    "table_number": "T01",
    "capacity": 4,
    "status": "STATUS_AVAILABLE",
    "location": "Main Hall"
  },
  "success": true,
  "message": "created"
}
```

## 3) Get table detail
- Method: `GET`
- Path: `/tables/{table_id}`
- Auth: No

Success response (`200`):

```json
{
  "table": {
    "table_id": "tb_001",
    "table_number": "T01",
    "capacity": 4,
    "status": "STATUS_AVAILABLE",
    "location": "Main Hall"
  },
  "success": true,
  "message": "ok"
}
```

## 4) Update table
- Method: `PUT`
- Path: `/tables/{table_id}`
- Auth: `Authorization: Bearer <access_token>`

Request body:

```json
{
  "table_number": "T01-A",
  "capacity": 6,
  "location": "VIP Room"
}
```

Success response (`200`):

```json
{
  "table": {
    "table_id": "tb_001",
    "table_number": "T01-A",
    "capacity": 6,
    "status": "STATUS_AVAILABLE",
    "location": "VIP Room"
  },
  "success": true,
  "message": "updated"
}
```

## 5) Delete table
- Method: `DELETE`
- Path: `/tables/{table_id}`
- Auth: `Authorization: Bearer <access_token>`

Success response (`200`):

```json
{
  "success": true,
  "message": "deleted"
}
```

## 6) Update table status
- Method: `PATCH`
- Path: `/tables/{table_id}/status`
- Auth: `Authorization: Bearer <access_token>`

Request body:

```json
{
  "status": "occupied",
  "order_id": "ord_1001"
}
```

`status` values:
- `available`
- `occupied`
- `reserved`
- `cleaning`
- `out_of_service`

Success response (`200`):

```json
{
  "table": {
    "table_id": "tb_001",
    "status": "STATUS_OCCUPIED",
    "current_order_id": "ord_1001"
  },
  "success": true,
  "message": "status updated"
}
```

## 7) Get available tables
- Method: `GET`
- Path: `/tables/available`
- Auth: No

Query params:
- `min_capacity` (optional, int)
- `location` (optional)

Example:
- `/tables/available?min_capacity=4&location=VIP`

Success response (`200`):

```json
{
  "tables": [
    {
      "table_id": "tb_002",
      "table_number": "T02",
      "capacity": 6,
      "status": "STATUS_AVAILABLE",
      "location": "VIP Room"
    }
  ],
  "success": true,
  "message": "ok"
}
```

## 8) Create reservation
- Method: `POST`
- Path: `/tables/{table_id}/reservations`
- Auth: `Authorization: Bearer <access_token>`

Request body:

```json
{
  "customer_name": "Nguyen Van A",
  "customer_phone": "0901234567",
  "notes": "Window seat",
  "start_time": "2026-05-27T18:00:00Z",
  "end_time": "2026-05-27T20:00:00Z",
  "items": [
    {
      "menu_item_id": "menu_001",
      "quantity": 2,
      "note": "No spicy"
    }
  ]
}
```

Success response (`200`):

```json
{
  "reservation": {
    "reservation_id": "res_001",
    "table_id": "tb_001",
    "customer_name": "Nguyen Van A",
    "customer_phone": "0901234567",
    "notes": "Window seat",
    "status": "RESERVATION_STATUS_RESERVED",
    "start_time": "2026-05-27T18:00:00Z",
    "end_time": "2026-05-27T20:00:00Z",
    "items": [
      {
        "menu_item_id": "menu_001",
        "quantity": 2,
        "note": "No spicy"
      }
    ]
  },
  "success": true,
  "message": "created"
}
```

## 9) List reservations by table
- Method: `GET`
- Path: `/tables/{table_id}/reservations`
- Auth: `Authorization: Bearer <access_token>`

Query params:
- `status` (optional): `reserved`, `cancelled`, `completed`
- `from` (optional, RFC3339)
- `to` (optional, RFC3339)
- `page` (default `1`)
- `page_size` (default `10`)

Example:
- `/tables/tb_001/reservations?status=reserved&from=2026-05-27T00:00:00Z&to=2026-05-28T00:00:00Z`

Success response (`200`):

```json
{
  "reservations": [
    {
      "reservation_id": "res_001",
      "table_id": "tb_001",
      "customer_name": "Nguyen Van A",
      "status": "RESERVATION_STATUS_RESERVED",
      "start_time": "2026-05-27T18:00:00Z",
      "end_time": "2026-05-27T20:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10,
  "success": true,
  "message": "ok"
}
```

## 10) Get reservation detail
- Method: `GET`
- Path: `/reservations/{reservation_id}`
- Auth: `Authorization: Bearer <access_token>`

Success response (`200`):

```json
{
  "reservation": {
    "reservation_id": "res_001",
    "table_id": "tb_001",
    "customer_name": "Nguyen Van A",
    "status": "RESERVATION_STATUS_RESERVED",
    "start_time": "2026-05-27T18:00:00Z",
    "end_time": "2026-05-27T20:00:00Z"
  },
  "success": true,
  "message": "ok"
}
```

## 11) Cancel reservation
- Method: `POST`
- Path: `/reservations/{reservation_id}/cancel`
- Auth: `Authorization: Bearer <access_token>`

Success response (`200`):

```json
{
  "reservation": {
    "reservation_id": "res_001",
    "status": "RESERVATION_STATUS_CANCELLED"
  },
  "success": true,
  "message": "cancelled"
}
```

## Common errors

```json
{
  "error": "invalid table status",
  "success": false
}
```
