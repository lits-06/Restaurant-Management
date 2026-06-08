# Table Service HTTP API - Postman Test Guide

Tai lieu nay cung cap cac API HTTP mau de test Table Service bang Postman.

Luu y:
- Table Service goc la gRPC.
- Cac API HTTP duoi day di qua API Gateway.

## 1) Base config

- Base URL: `http://localhost:8080`
- Content-Type: `application/json`
- Authorization (cho API tao/sua/xoa): `Bearer <ACCESS_TOKEN>`

Nen tao Postman Environment voi cac bien:
- `base_url` = `http://localhost:8080`
- `access_token` = token lay tu API login
- `table_id` = ID ban sau khi tao

## 2) Lay access token truoc khi test

### Request
- Method: `POST`
- URL: `{{base_url}}/auth/login`
- Body (raw JSON):

```json
{
  "email": "admin@restaurant.com",
  "password": "123456"
}
```

### Ghi chu
- Lay `access_token` tu response va gan vao bien `{{access_token}}`.

## 3) Create table

### Request
- Method: `POST`
- URL: `{{base_url}}/tables`
- Headers:
  - `Authorization: Bearer {{access_token}}`
- Body (raw JSON):

```json
{
  "table_number": "T01",
  "capacity": 4,
  "location": "Main Hall"
}
```

### Ket qua mong doi
- HTTP `200`
- Response co `table.table_id`.
- Gan `table_id` vao bien `{{table_id}}`.

## 4) List tables

### Request
- Method: `GET`
- URL: `{{base_url}}/tables?page=1&page_size=10&status=available&location=Main`

### Ghi chu
- `status` ho tro: `available`, `occupied`, `reserved`, `cleaning`, `out_of_service`.

## 5) Get table by ID

### Request
- Method: `GET`
- URL: `{{base_url}}/tables/{{table_id}}`

## 6) Update table info

### Request
- Method: `PUT`
- URL: `{{base_url}}/tables/{{table_id}}`
- Headers:
  - `Authorization: Bearer {{access_token}}`
- Body (raw JSON):

```json
{
  "table_number": "T01-A",
  "capacity": 6,
  "location": "VIP Room"
}
```

## 7) Update table status

### 7.1 Mark occupied

### Request
- Method: `PATCH`
- URL: `{{base_url}}/tables/{{table_id}}/status`
- Headers:
  - `Authorization: Bearer {{access_token}}`
- Body (raw JSON):

```json
{
  "status": "occupied",
  "order_id": "ORDER_1001"
}
```

### 7.2 Mark cleaning

### Request
- Method: `PATCH`
- URL: `{{base_url}}/tables/{{table_id}}/status`
- Headers:
  - `Authorization: Bearer {{access_token}}`
- Body (raw JSON):

```json
{
  "status": "cleaning"
}
```

### 7.3 Mark available

### Request
- Method: `PATCH`
- URL: `{{base_url}}/tables/{{table_id}}/status`
- Headers:
  - `Authorization: Bearer {{access_token}}`
- Body (raw JSON):

```json
{
  "status": "available"
}
```

## 8) Get available tables

### Request
- Method: `GET`
- URL: `{{base_url}}/tables/available?min_capacity=4&location=VIP`

## 9) Delete table

### Request
- Method: `DELETE`
- URL: `{{base_url}}/tables/{{table_id}}`
- Headers:
  - `Authorization: Bearer {{access_token}}`

## 10) Error cases nen test

### 10.1 Create table thieu du lieu

- Body:

```json
{
  "table_number": "",
  "capacity": 0,
  "location": ""
}
```

- Mong doi: loi `InvalidArgument` tu table-service qua gateway.

### 10.2 Update status occupied nhung khong co order_id

- Body:

```json
{
  "status": "occupied"
}
```

- Mong doi: loi validate do `order_id` bat buoc khi occupied.

### 10.3 Delete table dang occupied

- Mong doi: bi chan boi business rule, khong xoa duoc.

## 11) Bo status hop le

- `available`
- `occupied`
- `reserved`
- `cleaning`
- `out_of_service`

Ban cung co the gui theo enum string:
- `STATUS_AVAILABLE`
- `STATUS_OCCUPIED`
- `STATUS_RESERVED`
- `STATUS_CLEANING`
- `STATUS_OUT_OF_SERVICE`
