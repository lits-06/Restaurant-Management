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

| Tab label | Component | Ghi chú |
|-----------|-----------|---------|
| Dashboard | `AnalyticsOverview.tsx` | Tất cả staff |
| Menu | `MenuManagement.tsx` | ADMIN, MANAGER (chưa block route) |
| Orders | `OrdersManagement.tsx` | Tất cả staff |
| Staff | `MonthlyScheduler.tsx` | ADMIN, MANAGER (chưa block route) |

> Frontend hiện tại **chưa block routes theo role** — mọi staff đều thấy tất cả tab.

---

## Các Trang & Chức năng

### Login (`src/pages/Login.tsx`)
- Form: email + password
- Gọi `POST /auth/login` → nhận token
- Gọi `GET /users/{user_id}` → lấy user info + roles
- **Validate role:** Gọi `hasAdminAccess(roles)` — nếu không có staff role → "Không có quyền truy cập"
- Lưu vào `adminAuthStore` → App.tsx re-render dashboard

### AnalyticsOverview (`src/pages/AnalyticsOverview.tsx`)
- **Month picker**: chọn năm/tháng → filter dữ liệu theo tháng đó (click tháng = apply ngay)
- Load 500 orders + 100 menu items từ API một lần khi mount; filter client-side theo tháng được chọn
- **4 KPI card** (dữ liệu thật, so sánh trend % với tháng trước):
  - Doanh thu tháng — tổng `order.total` hoặc sum items
  - Tổng đơn đặt bàn
  - Giá trị TB/đơn
  - Tổng khách phục vụ (sum `party_size`)
- **Breakdown trạng thái** đơn trong tháng: Chờ xác nhận / Đã xác nhận / Hoàn thành / Đã hủy
- **Bảng Top 5 món** theo doanh thu trong tháng đang xem
- Đơn vị tiền: VNĐ (`toLocaleString('vi-VN')đ`)

### MenuManagement (`src/pages/MenuManagement.tsx`)
- Load menu items: `GET /menu/items?page_size=100`
- Load categories: `GET /menu/categories` → dùng làm tabs filter và dropdown trong modal
- Filter theo category tab + tìm kiếm theo tên
- **Thêm món:** Modal form → `POST /menu/items` (`name`, `price`, `description`, `image_url`, `category_id`)
- **Sửa món:** Modal form → `PUT /menu/items/{id}`
- **Xóa món:** `DELETE /menu/items/{id}` + confirm dialog
- Hiển thị ảnh, giá (VNĐ), tên, category name, VIP badge nếu ≥ 1.000.000đ
- **Category input là `<select>` dropdown** từ API (không phải text input) — payload gửi đúng `category_id` UUID

### OrdersManagement (`src/pages/OrdersManagement.tsx`)
- Load orders: `GET /orders?page_size=100`
- Load tables: `GET /tables?page_size=100` → resolve `table_id` → `"Bàn X"` (fallback UUID 8 ký tự)
- Load menu items: `GET /menu/items?page_size=200` → dùng cho search picker trong modal sửa món
- Filter theo status + tìm kiếm theo tên khách; pagination 5 đơn/trang
- Nút **Làm mới** để reload orders

**Bảng hiển thị:**
- Tên khách + SĐT + ghi chú ngắn (italic vàng nếu `notes` không trống)
- Giờ bắt đầu–kết thúc + ngày (`19:00 – 21:00`)
- Tên bàn (`Bàn 5`) + số khách
- Status badge

**3 action buttons mỗi hàng:**
1. `edit_note` → **Modal sửa thông tin**: tên, SĐT, ghi chú đặc biệt, ngày, giờ bắt đầu/kết thúc, số khách → `PUT /orders/{id}`
2. `restaurant` → **Modal sửa món**: tăng/giảm/xóa items + search box tìm từ menu thật → `PUT /orders/{id}`
3. `more_vert` → **Drawer chi tiết**: thông tin đặt bàn, ghi chú (banner vàng), `item_status` badge từng món (Chờ/Đang nấu/Xong/Đã mang), **nút cập nhật trạng thái**, tổng tiền

**Cập nhật trạng thái đơn (trong drawer):**
- Pending → nút **Xác nhận đơn** → `PATCH /orders/{id}/status { status: "Confirmed" }`
- Confirmed → nút **Hoàn thành** → `PATCH /orders/{id}/status { status: "Completed" }`
- Pending/Confirmed → nút **Hủy đơn** → `PATCH /orders/{id}/status { status: "Cancelled" }`

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
- Exports: `menuApi`, `ordersApi`, `tablesApi`, `scheduleApi`, `usersApi`, `authApi`
- `ordersApi.update(id, payload)` — payload gồm `name, phone, notes?, date, time, end_time?, party_size, items[]` (**không có `status`** — dùng `updateStatus`)
- `ordersApi.updateStatus(id, status)` → `PATCH /orders/{id}/status`
- `ordersApi.updateItemStatus(orderId, itemId, itemStatus)` → `PATCH /orders/{id}/items/{itemId}/status`
- `tablesApi.list({ page_size? })` → `GET /tables` (dùng để resolve `table_id` → `table_number`)
- `scheduleApi.list/create/update/delete` → `/schedule/shifts`
- `usersApi.listAll()` → `GET /users?page_size=200`

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

- **Role-based route guard** chưa implement — CHEF/WAITER có thể truy cập MenuManagement, Staff tab
- Không có token refresh — khi access token hết hạn phải login lại
- Không có pagination cho menu (load hết 100 items một lần)
- Không có real-time update khi có order mới (phải bấm Làm mới thủ công)
- MonthlyScheduler không hiển thị tên nhân viên nếu `/users` API lỗi (fallback về UUID)
- MonthlyScheduler không validate trùng ca (cùng nhân viên, cùng giờ)
- AnalyticsOverview load max 500 orders — nếu >500 orders thì thống kê theo tháng có thể thiếu
- OrdersManagement chưa có `UpdateOrderItemStatus` trong admin (CHEF/WAITER dùng kitchen app)
- Không có `user_id` hiển thị trong OrdersManagement (chỉ hiển thị tên khách)
