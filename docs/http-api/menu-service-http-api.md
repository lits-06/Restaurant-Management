# Menu Service HTTP API

Base URL: `http://localhost:8080`

## Auth rules
- Public: `GET /menu/items`, `GET /menu/categories`
- Requires Bearer token: `POST /menu/items`, `POST /menu/categories`

## 1) List menu items
- Method: `GET`
- Path: `/menu/items`
- Auth: No

Query params:
- `page` (default `1`)
- `page_size` (default `20`)
- `category_id` (optional)
- `status` (optional): `available`, `out_of_stock`, `discontinued`

Example:
- `/menu/items?page=1&page_size=20&status=available`

Success response (`200`):

```json
{
  "items": [
    {
      "item_id": "mi_001",
      "name": "Grilled Chicken Rice",
      "description": "Rice with grilled chicken",
      "price": 75000,
      "category_id": "cat_main",
      "category_name": "Main Dishes",
      "image_url": "https://example.com/chicken-rice.jpg",
      "status": "STATUS_AVAILABLE",
      "preparation_time": 20,
      "ingredients": ["rice", "chicken", "sauce"],
      "created_at": "2026-04-19T10:00:00Z",
      "updated_at": "2026-04-19T10:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 20,
  "success": true,
  "message": "ok"
}
```

## 2) Create menu item
- Method: `POST`
- Path: `/menu/items`
- Auth: `Authorization: Bearer <access_token>`

Request body:

```json
{
  "name": "Grilled Chicken Rice",
  "description": "Rice with grilled chicken",
  "price": 75000,
  "category_id": "cat_main",
  "image_url": "https://example.com/chicken-rice.jpg",
  "preparation_time": 20,
  "ingredients": ["rice", "chicken", "sauce"]
}
```

Success response (`200`):

```json
{
  "item": {
    "item_id": "mi_001",
    "name": "Grilled Chicken Rice",
    "status": "STATUS_AVAILABLE"
  },
  "success": true,
  "message": "created"
}
```

## 3) Get all categories
- Method: `GET`
- Path: `/menu/categories`
- Auth: No

Success response (`200`):

```json
{
  "categories": [
    {
      "category_id": "cat_main",
      "name": "Main Dishes",
      "description": "Main course items",
      "display_order": 1,
      "created_at": "2026-04-19T10:00:00Z",
      "updated_at": "2026-04-19T10:00:00Z"
    }
  ],
  "total": 1,
  "success": true,
  "message": "ok"
}
```

## 4) Create category
- Method: `POST`
- Path: `/menu/categories`
- Auth: `Authorization: Bearer <access_token>`

Request body:

```json
{
  "name": "Main Dishes",
  "description": "Main course items",
  "display_order": 1
}
```

Success response (`200`):

```json
{
  "category": {
    "category_id": "cat_main",
    "name": "Main Dishes",
    "description": "Main course items",
    "display_order": 1
  },
  "success": true,
  "message": "created"
}
```

## Common errors

```json
{
  "error": "missing Authorization header",
  "success": false
}
```
