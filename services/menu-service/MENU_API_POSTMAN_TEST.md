# Menu API Postman Test Script

This document describes HTTP test scripts for API Gateway endpoints that proxy to Menu Service.

## Base URL

- `http://localhost:8080`

## Auth Rules

- `GET /menu/items`: public endpoint, does not require login.
- `GET /menu/categories`: public endpoint, does not require login.
- `POST /menu/items`: requires `Authorization: Bearer <access_token>`.
- `POST /menu/categories`: requires `Authorization: Bearer <access_token>`.

## 1) Login to get access token

Method: `POST`
URL: `/auth/login`

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

Save token in Postman variable:
- Name: `access_token`
- Value: response field `access_token`

## 2) Create category (requires Authorization)

Method: `POST`
URL: `/menu/categories`

Headers:
- `Content-Type: application/json`
- `Authorization: Bearer {{access_token}}`

Body:

```json
{
  "name": "Main Dishes",
  "description": "Main course items",
  "display_order": 1
}
```

Expected response:
- `success = true`
- `category.category_id` is returned.

Save category id in variable:
- Name: `category_id`
- Value: response field `category.category_id`

## 3) Create menu item (requires Authorization)

Method: `POST`
URL: `/menu/items`

Headers:
- `Content-Type: application/json`
- `Authorization: Bearer {{access_token}}`

Body:

```json
{
  "name": "Grilled Chicken Rice",
  "description": "Rice with grilled chicken",
  "price": 75000,
  "category_id": "{{category_id}}",
  "image_url": "https://example.com/chicken-rice.jpg",
  "preparation_time": 20,
  "ingredients": ["rice", "chicken", "sauce"]
}
```

Expected response:
- `success = true`
- `item` object is returned.

## 4) List menu items (public, no Authorization)

Method: `GET`
URL: `/menu/items?page=1&page_size=20`

Optional query params:
- `category_id=<category_id>`
- `status=available|out_of_stock|discontinued`

Expected response:
- `success = true`
- `items` array contains created item.

## 5) Get all categories (public, no Authorization)

Method: `GET`
URL: `/menu/categories`

Expected response:
- `success = true`
- `categories` array is returned.

## 6) Negative tests

### 6.1 Create menu item without token

Method: `POST`
URL: `/menu/items`

Headers:
- `Content-Type: application/json`

Expected response:
- HTTP `401`
- JSON error message about missing or invalid Authorization header.

### 6.2 Create category with invalid token

Method: `POST`
URL: `/menu/categories`

Headers:
- `Content-Type: application/json`
- `Authorization: Bearer invalid-token`

Expected response:
- HTTP `401`
- JSON error message about invalid or expired token.

## Postman Pre-request Script (optional)

Use this if you want to auto-inject token into all protected requests:

```javascript
const token = pm.environment.get("access_token") || pm.collectionVariables.get("access_token");
if (token) {
  pm.request.headers.upsert({ key: "Authorization", value: `Bearer ${token}` });
}
```
