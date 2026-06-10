# restaurant-app (Customer Frontend)

**Port:** 5173 (Vite dev server)  
**Tech:** React 19 + Vite + TypeScript + TailwindCSS + Zustand

## Tổng quan

Giao diện dành cho **khách hàng** — xem menu, đặt bàn, xem và quản lý đơn đặt bàn của bản thân. Có thể dùng không cần đăng nhập (walk-in) hoặc đăng nhập để liên kết đơn với tài khoản.

## Cấu trúc

```
src/
  api/
    gateway.api.ts     ← Tất cả API calls, auto-inject Authorization header
  store/
    authStore.ts       ← Zustand persist (key: 'luxe-customer-auth')
  components/
    booking/           ← Reservation components (date/time picker, menu selector)
    common/            ← Button, Container, TopNavBar, SectionTitle
    home/              ← Hero, FeaturedDishes, Philosophy, CTA, Testimonials
    layout/            ← Header, Footer
    ui/                ← UI primitives
  pages/
    HomePage.tsx       ← Landing page
    MenuPage.tsx       ← Hiển thị menu đầy đủ
    ReservationPage.tsx ← Đặt bàn + chọn món
    ContactPage.tsx    ← Thông tin liên hệ
    LoginPage.tsx      ← Login + Register tabs
    MyOrdersPage.tsx   ← Danh sách đơn của user
  hooks/
    useAutoSlider.ts
    useHeaderScroll.ts
  lib/
    constants.ts
    types.ts
  assets/images/       ← Food images (jpg)
```

---

## Auth Store (`src/store/authStore.ts`)

**Persist key:** `luxe-customer-auth`

| Field | Type | Mô tả |
|-------|------|-------|
| `user` | `CustomerUser \| null` | `{ user_id, email, username, full_name, phone, roles[] }` |
| `accessToken` | `string \| null` | JWT access token |
| `refreshToken` | `string \| null` | JWT refresh token |
| `setAuth(user, access, refresh)` | fn | Login/register |
| `clearAuth()` | fn | Logout |

---

## Routing

Dùng `useState` để điều hướng (`currentPage`), không dùng react-router.

| Page | Path state | Auth required |
|------|-----------|---------------|
| HomePage | `home` | Không |
| MenuPage | `menu` | Không |
| ReservationPage | `reservation` | Không (walk-in OK) |
| ContactPage | `contact` | Không |
| LoginPage | `login` | Không |
| MyOrdersPage | `my-orders` | **Bắt buộc** |

**loginRedirect:** Khi truy cập trang cần auth, lưu trang đích vào state, redirect sang LoginPage. Sau login thành công → redirect về trang đích.

---

## Các Trang & Chức năng

### HomePage
- Hero section với CTA "Book a Table" và "Explore Menu"
- FeaturedDishes: hiển thị 4–6 món nổi bật (hard-coded hoặc từ API)
- Philosophy, Testimonials, Footer

### MenuPage
- Tải danh sách món từ `GET /menu/items`
- Filter theo category
- Tìm kiếm theo tên
- Hiển thị ảnh, tên, giá, mô tả từng món

### ReservationPage
- **Bước 1:** Chọn ngày, giờ đến, giờ về (default: +2h), số khách, tên, SĐT, ghi chú
  - Time picker bị giới hạn 10:00–22:00 (frontend-only)
  - `end_time` mặc định = `time + 2h`, user có thể điều chỉnh
- **Bước 2:** Chọn món từ menu thực tế (gọi `GET /menu/items`)
  - Có fallback về `mockDishes` nếu API lỗi
  - Tính tổng tiền trước order
- **Submit:** `POST /orders` với `user_id` từ JWT (nếu đã đăng nhập)
  - `table_id` không truyền → auto-assign bởi order-service
- Nếu chưa đăng nhập → vẫn cho phép đặt (walk-in order)

### LoginPage
- **Tab Login:** email + password → `POST /auth/login`
  - Sau login: lấy user info từ `GET /users/{user_id}` → lưu vào store
- **Tab Register:** email + password + username + full_name + phone → `POST /auth/register`
  - Sau register thành công: pre-fill email vào form login
- Redirect về `loginRedirect` page sau login thành công

### MyOrdersPage
- Load orders của user: `GET /orders?user_id={id}`
- Hiển thị danh sách đơn: ngày, giờ, bàn, party size, status, tổng tiền, danh sách món
- **Thêm món:** chọn từ menu → `POST /orders/{id}/items`
- **Hủy đơn (2 bước):** confirm dialog → `POST /orders/{id}/cancel`
  - Chỉ cancel được đơn `Pending` hoặc `Confirmed`

---

## API Calls (`src/api/gateway.api.ts`)

| Function | HTTP | Endpoint | Auth |
|----------|------|---------|------|
| `menuApi.listItems(query)` | GET | `/menu/items` | Không |
| `authApi.login(email, pwd)` | POST | `/auth/login` | Không |
| `authApi.register(data)` | POST | `/auth/register` | Không |
| `authApi.logout(refreshToken)` | POST | `/auth/logout` | Bearer |
| `usersApi.getOne(id)` | GET | `/users/{id}` | Bearer |
| `ordersApi.create(payload)` | POST | `/orders` | Optional |
| `ordersApi.getOne(id)` | GET | `/orders/{id}` | Optional |
| `ordersApi.list(query)` | GET | `/orders` | Optional |
| `ordersApi.cancel(id)` | POST | `/orders/{id}/cancel` | Optional |
| `ordersApi.addItem(orderId, item)` | POST | `/orders/{id}/items` | Optional |
| `tableApi.getOne(id)` | GET | `/tables/{id}` | Không |

**Base URL:** `VITE_API_BASE_URL` env (default: `http://localhost:8080`)

---

## TopNavBar

Hiển thị khác nhau theo auth state:

| Auth State | Hiển thị |
|-----------|---------|
| Chưa đăng nhập | Logo + Nav links + "Đăng nhập" + "Book Now" |
| Đã đăng nhập | Logo + Nav links + "My Orders" + username + "Logout" |

---

## Role trong Frontend

App này chỉ dùng role `USER`. Staff roles (ADMIN/MANAGER/CHEF/WAITER) redirect sang admin hoặc kitchen app nếu login.

---

## Vấn đề chưa giải quyết / Có thể triển khai thêm

- Không có token refresh flow — khi access token hết hạn, user phải login lại
- `ReservationPage` dùng mock dishes làm fallback thay vì proper error state
- Operating hours (10:00–22:00) chỉ validate ở frontend
- Không hiển thị `table_number` sau khi đặt bàn thành công (chỉ có `table_id` UUID)
- `end_time` mặc định `+2h` không cộng qua midnight chính xác (edge case)
- MyOrdersPage không hiển thị `item_status` từng món (chưa update cho field mới)
- Không có real-time update trạng thái order (phải reload trang)
- Chưa có trang profile/account settings
- Chưa có trang xem lịch sử đặt bàn theo calendar view
- Không có email confirmation sau khi đặt bàn
