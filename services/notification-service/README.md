# Notification Service

**Port:** 50058 (gRPC)  
**Module:** `restaurant-management/services/notification-service`  
**Proto:** `proto/notification/notification.proto`  
**Database:** Không có — Redis Pub/Sub thuần túy

## Tổng quan

Notification Service là message bus cho kitchen staff. Nhận sự kiện từ order-service và đẩy real-time xuống các client (CHEF/WAITER) qua api-gateway WebSocket. Service này **không lưu** notification — là fire-and-forget, không có persistence.

## Kiến trúc

```
cmd/server/main.go
internal/
  domain/
    notification.go  ← Notification struct, type constants, role constants
  repository/
    pubsub.go        ← Redis Pub/Sub: Publish + Subscribe
  usecase/
    notification_usecase.go ← Send (assign UUID+timestamp) + Subscribe
  delivery/grpc/
    notification_handler.go ← gRPC handler (streaming Subscribe)
pkg/config/config.go
```

---

## 2 Loại Notification

| Type | Trigger | Target Role | Payload |
|------|---------|-------------|---------|
| `ORDER_CONFIRMED` | `order-service.UpdateOrderStatus → Confirmed` | `CHEF` | `order_id`, `table_id`, `customer_name`, `party_size`, `notes`, `items[]` |
| `ITEM_READY` | `order-service.UpdateOrderItemStatus → READY` | `WAITER` | `order_id`, `table_id`, `item_id`, `item_name` |

---

## gRPC API — 2 RPCs

### 1. `SendNotification`
Publish notification vào Redis channel `notifications:{target_role}`.

**Request:**
| Field | Type | Mô tả |
|-------|------|-------|
| `type` | string | `ORDER_CONFIRMED` hoặc `ITEM_READY` |
| `target_role` | string | `CHEF` hoặc `WAITER` |
| `order_id` | string | UUID |
| `table_id` | string | UUID |
| `item_id` | string | UUID (chỉ cho ITEM_READY) |
| `item_name` | string | Tên món (chỉ cho ITEM_READY) |
| `message` | string | Human-readable message |
| `customer_name` | string | (chỉ cho ORDER_CONFIRMED) |
| `party_size` | int32 | (chỉ cho ORDER_CONFIRMED) |
| `notes` | string | Yêu cầu đặc biệt (chỉ cho ORDER_CONFIRMED) |
| `items[]` | NotificationOrderItem[] | Danh sách món (chỉ cho ORDER_CONFIRMED) |

**Response:** `success`, `message`

---

### 2. `Subscribe` (Server-side Streaming)
Subscribe vào Redis channel theo role, stream notification về client.

**Request:** `role` (string: `CHEF` hoặc `WAITER`)  
**Response stream:** `Notification` messages liên tục cho đến khi client disconnect

**Notification fields:**
| Field | Type |
|-------|------|
| `id` | string (UUID, assigned bởi usecase) |
| `type` | string |
| `target_role` | string |
| `order_id` | string |
| `table_id` | string |
| `item_id` | string |
| `item_name` | string |
| `created_at` | int64 (Unix timestamp) |
| `message` | string |
| `customer_name` | string |
| `party_size` | int32 |
| `notes` | string |
| `items[]` | NotificationOrderItem[] |

---

## Redis Channels

| Channel | Subscriber |
|---------|------------|
| `notifications:CHEF` | Kitchen app (KitchenPage) |
| `notifications:WAITER` | Kitchen app (WaiterPage) |

---

## Luồng dữ liệu đầy đủ

```
order-service usecase
  └─ go func() { notifClient.SendNotification(...) }  ← fire-and-forget

notification-service gRPC handler
  └─ usecase.Send(req)
       └─ Assign UUID + created_at timestamp
       └─ repo.Publish(role, notification)
            └─ json.Marshal → redis.PUBLISH notifications:{role}

api-gateway NotificationHandler.Subscribe (per WebSocket client)
  └─ notifClient.Subscribe(ctx, role)  ← gRPC streaming call
       └─ notification-service handler
            └─ redis.Subscribe notifications:{role}
            └─ stream each message → gRPC Send
  └─ json.Marshal → websocket.WriteMessage
  └─ Drain goroutine: websocket.ReadMessage → cancel ctx → close gRPC stream
```

---

## Cấu hình (Environment Variables)

| Biến | Default | Mô tả |
|------|---------|-------|
| `SERVER_PORT` | `50058` | gRPC listen port |
| `REDIS_HOST` | `localhost` | |
| `REDIS_PORT` | `6379` | |
| `REDIS_PASSWORD` | `` | |
| `REDIS_DB` | `0` | |

---

## Vấn đề chưa giải quyết / Có thể triển khai thêm

- **Không có persistence** — nếu CHEF/WAITER offline khi notification được gửi, họ sẽ bỏ lỡ thông báo đó
- Không có delivery confirmation — không biết notification có được nhận hay không
- Không có retry mechanism nếu Redis unreachable
- Không có notification history / inbox
- Kitchen app chưa có auto-reconnect WebSocket — nếu mất kết nối phải reload trang
- `table_id` trong notification là UUID — cần thêm `table_number` vào payload để chef/waiter biết bàn số mấy
- Không có notification khi order bị cancel (CHEF có thể đang nấu món không cần thiết nữa)
- Không có threshold / rate limiting — event storm có thể flood Redis
