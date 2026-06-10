# restaurant-app-admin (Admin Dashboard)

**Port:** 5174 (Vite dev server)  
**Tech:** React 19 + Vite + TypeScript + TailwindCSS + Zustand

## Tổng quan

Dashboard quản lý dành cho **nhân viên nhà hàng** (ADMIN, MANAGER, CHEF, WAITER). Cung cấp giao diện quản lý menu, đơn đặt bàn, nhân viên, và phân tích doanh thu.

## Cấu trúc

```
src/
  store/
    adminAuthStore.ts      ← Zustand persist (key: 'luxe-admin-auth')
  services/
    api.ts                 ← fetch wrapper, auto-inject Bearer token
    auth.ts                ← Login helper
  components/
    Footer.tsx
    HeaderDashboard.tsx
    KPIGrid.tsx
    PerformanceTable.tsx
    Sidebar.tsx            ← Navigation + user info + logout
  pages/
    Login.tsx              ← Login form, validate staff role
    AnalyticsOverview.tsx  ← Dashboard KPIs + charts
    MenuManagement.tsx     ← CRUD menu items
    OrdersManagement.tsx   ← Danh sách và quản lý đơn hàng
    MonthlyScheduler.tsx   ← Lịch làm việc theo tháng (calendar grid)
  App.tsx                  ← Auth gate + router
```

---

## Auth Store (`src/store/adminAuthStore.ts`)

**Persist key:** `luxe-admin-auth`

| Field | Type |
|-------|------|
| `user` | `AdminUser \| null` (`{ user_id, email, username, full_name, roles[] }`) |
| `accessToken` | `string \| null` |
| `refreshToken` | `string \| null` |
| `setAuth(user, access, refresh)` | fn |
| `clearAuth()` | fn |

**`hasAdminAccess(roles)`:** Trả về `true` nếu roles chứa ít nhất 1 trong `[ADMIN, MANAGER, CHEF, WAITER]`

---

## Auth Gate

`App.tsx` kiểm tra `useAdminAuthStore().user`:
- Nếu `null` → render `<LoginPage />`
- Nếu có → render dashboard với Sidebar + active page

---

## Routing (Sidebar navigation)

| Page | Roles được dùng |
|------|----------------|
| AnalyticsOverview | Tất cả staff |
| MenuManagement | ADMIN, MANAGER |
| OrdersManagement | Tất cả staff |
| StaffManagement | ADMIN, MANAGER |
| WeeklyScheduler | ADMIN, MANAGER |

> Frontend hiện tại **chưa block routes theo role** — mọi staff đều thấy tất cả menu items. Cần implement role-based sidebar visibility.

---

## Các Trang & Chức năng

### Login (`src/pages/Login.tsx`)
- Form: email + password
- Gọi `POST /auth/login` → nhận token
- Gọi `GET /users/{user_id}` → lấy user info + roles
- **Validate role:** Gọi `hasAdminAccess(roles)` — nếu không có staff role → "Không có quyền truy cập"
- Lưu vào `adminAuthStore` → App.tsx re-render dashboard

### AnalyticsOverview (`src/pages/AnalyticsOverview.tsx`)
- KPI cards: tổng đơn, doanh thu, tỷ lệ hoàn thành (mock data)
- Chart doanh thu theo thời gian (mock data)
- **Chưa kết nối report-service**

### MenuManagement (`src/pages/MenuManagement.tsx`)
- Load menu items: `GET /menu/items` (tất cả, page_size lớn)
- Filter theo category tab
- Tìm kiếm theo tên
- **Thêm món:** Modal form → `POST /menu/items` (`name`, `price`, `description`, `image`, `category`)
- **Sửa món:** Modal form → `PUT /menu/items/{id}`
- **Xóa món:** `DELETE /menu/items/{id}`
- Hiển thị ảnh món (`image_url`), giá, tên, category
- VIP badge nếu giá ≥ 1.000.000đ

### OrdersManagement (`src/pages/OrdersManagement.tsx`)
- Load orders: `GET /orders?page_size=100&status={filter}&keyword={search}`
- Filter theo status (All/Confirmed/Pending/Completed/Cancelled)
- Tìm kiếm theo tên khách
- Pagination (5 đơn/trang)
- **Xem chi tiết:** Slide-over drawer — tên, SĐT, giờ, party size, danh sách món
- **Cập nhật thông tin đặt chỗ:** Modal edit → `PUT /orders/{id}` (tên, SĐT, ngày, giờ, party size, status)
- **Sửa món ăn trong order:** Modal → thêm/xóa items → `PUT /orders/{id}` với items mới
- Hiển thị: ID, tên khách, SĐT, ngày giờ, party size, status, tổng tiền
- **Chưa hiển thị:** `notes`, `end_time`, `user_id`, `item_status` (các field mới chưa được update)

### MonthlyScheduler (`src/pages/MonthlyScheduler.tsx`)
- Thay thế `StaffManagement.tsx` và `WeeklyScheduler.tsx` (đã xóa)
- **Calendar grid** tháng: hàng = tuần (T2–CN), mỗi ô = 1 ngày
- Load danh sách nhân viên: `GET /users?page_size=200` → filter staff roles client-side
- Load shifts: `GET /schedule/shifts?month=YYYY-MM`
- **Chip màu** theo role trong mỗi ô ngày: amber=CHEF, blue=WAITER, purple=MANAGER
- **Tạo ca:** Click "+" trong ô ngày → modal chọn nhân viên + giờ + role + notes → `POST /schedule/shifts`
- **Xóa ca:** Click chip → popover chi tiết → `DELETE /schedule/shifts/{id}`
- **Filter** theo role qua dropdown ở header
- **Navigation** tháng trước/sau

---

## Sidebar (`src/components/Sidebar.tsx`)
- Hiển thị tên user, email từ `adminAuthStore`
- Navigation links với active state
- Logout: gọi `POST /auth/logout` → `clearAuth()` → re-render Login

---

## API Service (`src/services/api.ts`)
- Native `fetch` wrapper với `baseURL: http://localhost:8080` (không dùng axios)
- Tự động inject `Authorization: Bearer {accessToken}` từ `adminAuthStore`
- Exports: `menuApi`, `ordersApi`, `scheduleApi`, `usersApi`, `authApi`
- `scheduleApi.list/create/update/delete` → `/schedule/shifts`
- `usersApi.listAll()` → `GET /users?page_size=200` (dùng trong MonthlyScheduler để load staff picker)

---

## Roles & Quyền Truy Cập

| Role | Mô tả quyền trong Admin App |
|------|----------------------------|
| `ADMIN` | Full access — tất cả trang |
| `MANAGER` | Full access — tất cả trang |
| `CHEF` | Orders, Analytics |
| `WAITER` | Orders, Analytics |

> Hiện tại route không được guard theo role — chỉ check có staff role hay không để vào dashboard.

---

## Vấn đề chưa giải quyết / Có thể triển khai thêm

- **OrdersManagement chưa hiển thị `notes`, `end_time`, `user_id`, `item_status`** — cần update sau khi thêm các field này
- **Role-based route guard** chưa implement — CHEF/WAITER có thể truy cập MenuManagement
- AnalyticsOverview dùng mock data — chưa kết nối report-service
- Không có token refresh — khi access token hết hạn phải login lại
- Không có confirmation dialog khi xóa menu (có thể xóa nhầm)
- Không có pagination cho menu (load hết 1 lần)
- OrdersManagement không có `UpdateOrderStatus` cho từng order (chỉ có PUT toàn bộ order)
- Không có `UpdateOrderItemStatus` trong admin dashboard — CHEF/WAITER phải dùng kitchen app
- Không có real-time update khi có order mới (phải reload thủ công)
- MonthlyScheduler không hiển thị tên nhân viên nếu `/users` API lỗi (fallback về UUID)
- MonthlyScheduler không validate trùng ca (cùng nhân viên, cùng giờ)
