# Menu Service

**Port:** 50054 (gRPC)  
**Module:** `restaurant-management/services/menu-service`  
**Proto:** `proto/menu/menu.proto`  
**Database:** PostgreSQL (`restaurant_db`)

## Tổng quan

Menu Service quản lý thực đơn nhà hàng: danh mục món ăn (Category) và các món ăn cụ thể (MenuItem). Các service khác (order-service) gọi vào đây để validate item và lấy giá.

## Kiến trúc

```
cmd/server/main.go
internal/
  domain/
    menu_item.go    ← MenuItem entity, validation
    category.go     ← Category entity, validation
    errors.go       ← Typed domain errors
  repository/
    repository.go         ← Interfaces
    menu_postgres.go      ← PostgreSQL: menu_items CRUD + ensureSchema
    category_postgres.go  ← PostgreSQL: menu_categories CRUD + ensureSchema
  usecase/
    menu_usecase.go       ← Business logic
  delivery/grpc/
    menu_handler.go       ← gRPC handler
pkg/config/config.go
```

---

## gRPC API — 10 RPCs

### MenuItem — 5 RPCs

#### 1. `CreateMenuItem`

**Request:**
| Field | Type | Required | Mô tả |
|-------|------|----------|-------|
| `name` | string | ✓ | Unique trong hệ thống |
| `description` | string | — | Mô tả món |
| `price` | double | ✓ | > 0 |
| `category` | string | ✓ | Category ID (UUID) |
| `image_url` | string | — | URL ảnh món |

**Response:** `MenuItem`, `success`, `message`

---

#### 2. `GetMenuItem`

**Request:** `item_id` (UUID)  
**Response:** `MenuItem`, `success`, `message`

---

#### 3. `UpdateMenuItem`

**Request:** `item_id`, `name`, `description`, `price`, `category_id`, `image_url`  
**Response:** `MenuItem`, `success`, `message`

---

#### 4. `DeleteMenuItem`
Cascade xóa references trong `order_items` (FK constraint).

**Request:** `item_id`  
**Response:** `success`, `message`

---

#### 5. `ListMenuItems`

**Request:**
| Field | Type | Mô tả |
|-------|------|-------|
| `page` | int32 | Default: 1 |
| `page_size` | int32 | Default: 10 |
| `category_id` | string | Filter theo category |
| `keyword` | string | Tìm trong name (ILIKE) |

**Response:** `items[]`, `total`, `page`, `page_size`, `success`, `message`

---

### Category — 5 RPCs

#### 6. `CreateCategory`

**Request:** `name` (required, unique), `description`, `display_order`  
**Response:** `Category`, `success`, `message`

> **Lưu ý:** `description` và `display_order` có trong proto nhưng **không được lưu vào DB** (schema chỉ có `category_id`, `name`).

#### 7. `GetCategory`

**Request:** `category_id`  
**Response:** `Category`, `success`, `message`

#### 8. `UpdateCategory`

**Request:** `category_id`, `name`, `description`, `display_order`  
**Response:** `Category`, `success`, `message`

#### 9. `DeleteCategory`
Cascade xóa tất cả `menu_items` trong category (FK ON DELETE CASCADE).

**Request:** `category_id`  
**Response:** `success`, `message`

#### 10. `ListCategories`

**Request:** `page`, `page_size`  
**Response:** `categories[]`, `total`, `success`, `message`

---

## Database Schema

### Bảng `menu_categories`

```sql
CREATE TABLE IF NOT EXISTS menu_categories (
    category_id VARCHAR(36)  PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL UNIQUE
);
```

### Bảng `menu_items`

```sql
CREATE TABLE IF NOT EXISTS menu_items (
    item_id     VARCHAR(36)       PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100)      NOT NULL UNIQUE,
    description TEXT              NOT NULL DEFAULT '',
    price       DOUBLE PRECISION  NOT NULL,
    category_id VARCHAR(36)       NOT NULL REFERENCES menu_categories(category_id) ON DELETE CASCADE,
    image_url   TEXT              NOT NULL DEFAULT ''
);
```

### Indexes

```sql
CREATE INDEX IF NOT EXISTS idx_menu_items_category_id ON menu_items(category_id);
```

---

## Tích hợp với các Service Khác

- **order-service** gọi `GetMenuItem(item_id)` trong `AddOrderItem` để:
  - Validate item tồn tại
  - Lấy price hiện tại (snapshot vào order_items)
  - Lấy category name để hiển thị trong kitchen
- `order_items.item_id` là FK tham chiếu `menu_items.item_id` — nếu xóa menu item, các order_items tương ứng cũng bị cascade delete

---

## Cấu hình (Environment Variables)

| Biến | Default |
|------|---------|
| `SERVER_PORT` | `50054` |
| `DB_HOST` | `localhost` |
| `DB_PORT` | `5432` |
| `DB_USER` | `restaurant_user` |
| `DB_PASSWORD` | `restaurant_pass` |
| `DB_NAME` | `restaurant_db` |
| `DB_SSLMODE` | `disable` |

---

## Vấn đề chưa giải quyết / Có thể triển khai thêm

- `description` và `display_order` của Category có trong proto nhưng không lưu vào DB — cần migration thêm cột
- Giá món (`price`) không có currency field — assume VND
- Xóa MenuItem cascade xóa `order_items` — mất lịch sử order đã đặt món đó
- Không có `is_available` flag để tạm ẩn món mà không xóa
- Không có hình ảnh upload — chỉ lưu URL text (cần CDN/storage service)
- `name` phải unique trong toàn bộ system — không thể có 2 món khác category trùng tên
- Không có versioning giá — nếu giá thay đổi, order cũ vẫn hiển thị giá mới (vì price được snapshot vào `order_items` khi order được tạo, nhưng chỉ khi thêm mới items, không phải lúc tạo order ban đầu)
