# restaurant-app-kitchen (Kitchen Display)

**Port:** 5175 (Vite dev server)  
**Tech:** React 19 + Vite + TypeScript + TailwindCSS + Zustand + Native WebSocket

## Tổng quan

Giao diện hiển thị bếp dành cho **CHEF** và **WAITER**. Kết nối real-time qua WebSocket để nhận thông báo từ hệ thống. Dark theme, thiết kế cho tablet/màn hình bếp.

## Cấu trúc

```
src/
  store/
    authStore.ts           ← Zustand persist (key: 'kitchen-auth')
  api/
    gateway.api.ts         ← ordersApi, authApi, usersApi, scheduleApi, createNotificationWS
  hooks/
    useNotifications.ts    ← WebSocket hook, max 50 notifications
  pages/
    LoginPage.tsx          ← Dark-themed login, validate CHEF/WAITER/ADMIN/MANAGER
    KitchenPage.tsx        ← CHEF view
    WaiterPage.tsx         ← WAITER view
    SchedulePage.tsx       ← My schedule view (self-service shift registration)
  App.tsx                  ← Auth gate + role-based routing + floating switcher
  vite-env.d.ts
```

---

## Auth Store (`src/store/authStore.ts`)

**Persist key:** `kitchen-auth`

| Field | Type |
|-------|------|
| `user` | `KitchenUser \| null` (`{ user_id, email, username, full_name, roles[] }`) |
| `accessToken` | `string \| null` |
| `refreshToken` | `string \| null` |
| `setAuth(user, access, refresh)` | fn |
| `clearAuth()` | fn |

**`hasKitchenAccess(roles)`:** Trả về `true` nếu có CHEF/WAITER/ADMIN/MANAGER  
**`getDefaultRole(roles)`:** Trả về role mặc định cho routing

---

## Role Routing (`App.tsx`)

`ActiveView = 'CHEF' | 'WAITER' | 'SCHEDULE'`

| Role | Trang mặc định | Floating switcher |
|------|---------------|-------------------|
| `CHEF` | `KitchenPage` | 🍳 Bếp + 📅 Lịch |
| `WAITER` | `WaiterPage` | 🛎 Phục vụ + 📅 Lịch |
| `ADMIN` | `KitchenPage` | 🍳 Bếp + 🛎 Phục vụ + 📅 Lịch |
| `MANAGER` | `KitchenPage` | 🍳 Bếp + 🛎 Phục vụ + 📅 Lịch |

---

## Các Trang & Chức năng

### LoginPage
- Dark theme (background `gray-900`)
- Form: email + password
- Gọi `POST /auth/login` → nhận token
- Gọi `GET /users/{user_id}` → lấy roles
- Validate: phải có CHEF/WAITER/ADMIN/MANAGER → redirect sang trang phù hợp
- Nếu chỉ có role `USER` → "Không có quyền truy cập"

---

### KitchenPage (CHEF View)
**Mục đích:** Hiển thị các order đã confirmed, cho CHEF mark từng món COOKING → READY.

**Load data:** `GET /orders?status=Confirmed&page_size=50`

**Hiển thị:**
- Header: tên CHEF, indicator WebSocket (Live/Offline), nút logout
- Grid các order đang active (có ít nhất 1 món chưa SERVED)
- Mỗi order card: order ID (8 ký tự), table_id (8 ký tự UUID), tên khách, party size, notes
- Danh sách món trong order với status badge:
  - `PENDING` → badge xám "Chờ" + button "Bắt đầu nấu"
  - `COOKING` → badge cam "Đang nấu" + button "Đã xong ✓"
  - `READY` → badge xanh lá "Xong" (không có button — waiter xử lý)
  - `SERVED` → badge xám gạch "Đã mang" (strikethrough)

**Actions:**
- "Bắt đầu nấu" → `PATCH /orders/{id}/items/{itemId}/status` với `item_status: "COOKING"`
- "Đã xong ✓" → `PATCH /orders/{id}/items/{itemId}/status` với `item_status: "READY"`

**WebSocket:** Subscribe với `role=CHEF`
- Nhận `ORDER_CONFIRMED` → tự động refresh danh sách order

---

### WaiterPage (WAITER View)
**Mục đích:** Feed tất cả các món đã READY cần mang ra bàn, với notification sidebar.

**Load data:** `GET /orders?status=Confirmed&page_size=50`

**Layout:** 2 cột
- **Cột trái (main):** Danh sách tất cả items có `item_status === 'READY'` từ tất cả order
  - Mỗi item: tên món, bàn (table_id UUID), tên khách, party size, notes của order
  - Button "Đã mang ra ✓" → `PATCH /orders/{id}/items/{itemId}/status` với `item_status: "SERVED"`
  - Badge đếm số món đang sẵn sàng (animate-pulse khi > 0)
- **Cột phải (sidebar):** Lịch sử notifications (WebSocket)
  - Badge "N thông báo mới" + button "Xóa tất cả"
  - Mỗi notification: type, message, tên khách, party size, notes

**WebSocket:** Subscribe với `role=WAITER`
- Nhận `ITEM_READY` → tự động refresh orders để cập nhật feed

---

## WebSocket Hook (`src/hooks/useNotifications.ts`)

```typescript
const { notifications, connected, clearNotifications } = useNotifications(accessToken, role)
```

| Output | Type | Mô tả |
|--------|------|-------|
| `notifications` | `KitchenNotification[]` | Max 50, mới nhất ở đầu |
| `connected` | boolean | `true` nếu WebSocket ở state OPEN |
| `clearNotifications()` | fn | Xóa toàn bộ notification state |

**WebSocket URL:** `ws://localhost:8080/ws/notifications?token={accessToken}&role={role}`

**Lifecycle:**
- Mount: tạo WebSocket
- `onopen` → `connected = true`
- `onmessage` → parse JSON → prepend vào notifications (giới hạn 50)
- `onclose/onerror` → `connected = false`
- Unmount: `ws.close()`

---

## API Calls (`src/api/gateway.api.ts`)

| Function | HTTP | Endpoint |
|----------|------|---------|
| `authApi.login(email, pwd)` | POST | `/auth/login` |
| `usersApi.getOne(id)` | GET | `/users/{id}` |
| `ordersApi.list(query)` | GET | `/orders` |
| `ordersApi.updateItemStatus(orderId, itemId, status)` | PATCH | `/orders/{id}/items/{itemId}/status` |
| `scheduleApi.myShifts(userId, month)` | GET | `/schedule/shifts?user_id=...&month=...` |
| `scheduleApi.create(payload)` | POST | `/schedule/shifts` |
| `scheduleApi.delete(shiftId)` | DELETE | `/schedule/shifts/{id}` |
| `createNotificationWS(token, role)` | WS | `ws://localhost:8080/ws/notifications?token=...&role=...` |

---

## Role-based Authorization (api-gateway)

| Action | Allowed Roles |
|--------|--------------|
| Mark item → `COOKING` | CHEF, ADMIN, MANAGER |
| Mark item → `READY` | CHEF, ADMIN, MANAGER |
| Mark item → `SERVED` | WAITER, ADMIN, MANAGER |

---

---

### SchedulePage (`src/pages/SchedulePage.tsx`)
**Mục đích:** Xem và tự đăng ký ca làm việc của bản thân.

**Load data:** `GET /schedule/shifts?user_id={myId}&month={YYYY-MM}`

**Hiển thị:**
- Header: tháng hiện tại + nút tháng trước/sau + nút "Đăng ký ca"
- Danh sách ca dạng list: ngày, giờ bắt đầu–kết thúc, role, notes
- Nút "Xóa" cho mỗi ca → `DELETE /schedule/shifts/{id}`

**Đăng ký ca mới (modal):**
- Date picker (min = hôm nay)
- Start time / End time (HH:MM)
- Notes (optional)
- `POST /schedule/shifts` với `user_id` = mình, `role` = role của user hiện tại
- Ca có hiệu lực ngay (không cần duyệt)

---

## Vấn đề chưa giải quyết / Có thể triển khai thêm

- **Không có auto-reconnect WebSocket** — mất kết nối phải reload trang
- **`table_id` hiển thị dạng UUID cắt ngắn** — cần lookup `table_number` từ table-service để hiển thị "Bàn 5" thay vì "a1b2c3d4"
- Không có polling fallback khi WebSocket unavailable
- Không có sound notification khi ORDER_CONFIRMED hoặc ITEM_READY
- KitchenPage load tất cả Confirmed orders — không filter theo thời gian (có thể hiện orders cũ từ hôm trước)
- Không có trang xem lịch sử orders đã completed
- Không có multi-language support (chỉ tiếng Việt)
- Chưa có trang settings (đổi mật khẩu)
- ADMIN/MANAGER có thể thực hiện cả CHEF và WAITER actions nhưng chỉ thấy 1 view tại 1 thời điểm
- Không có offline mode — hoàn toàn dependent vào API và WebSocket
- SchedulePage không validate trùng ca (cùng ngày, cùng giờ) — backend cũng không block
