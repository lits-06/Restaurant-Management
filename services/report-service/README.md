# Report Service

**Port:** 50059 (gRPC)  
**Module:** `restaurant-management/services/report-service`  
**Proto:** `proto/report/report.proto`  
**Database:** In-memory (PostgreSQL chưa implement)

## Tổng quan

Report Service cung cấp analytics và báo cáo cho hệ thống nhà hàng: doanh thu, đơn hàng, nhân viên, món ăn phổ biến. **Hiện tại chưa được kết nối vào docker-compose và api-gateway** — proto và service skeleton đã có nhưng chưa integrate vào luồng thực tế.

## Kiến trúc

```
cmd/server/main.go
internal/
  domain/      ← Report structs (SalesReport, OrderReport, etc.)
  repository/  ← In-memory implementation (mock data)
  usecase/     ← Report usecase
  delivery/grpc/ ← gRPC handler
pkg/config/config.go
```

---

## gRPC API — 7 RPCs

### 1. `GetSalesReport`

**Request:** `period` (DAILY/WEEKLY/MONTHLY/YEARLY/CUSTOM), `from_date`, `to_date`  
**Response:**
```
SalesReport {
  total_sales, total_tax, total_discount, net_revenue
  total_orders, average_order_value
  daily_breakdown[]  ← { date, sales, orders }
  payment_breakdown[] ← { method, amount, count }
}
```

---

### 2. `GetInventoryReport`

**Request:** `as_of_date`  
**Response:**
```
InventoryReport {
  total_inventory_value, low_stock_items_count
  items[] ← { ingredient_id, name, current_stock, minimum_stock, unit_cost, total_value, is_low_stock }
  total_waste_value
}
```

> **Note:** Inventory service chưa tồn tại — report này trả về mock data.

---

### 3. `GetOrderReport`

**Request:** `period`, `from_date`, `to_date`  
**Response:**
```
OrderReport {
  total_orders, completed_orders, cancelled_orders
  completion_rate, average_preparation_time
  hourly_breakdown[] ← { hour, order_count }
}
```

---

### 4. `GetStaffPerformanceReport`

**Request:** `period`, `from_date`, `to_date`, `staff_id` (optional, filter 1 nhân viên)  
**Response:**
```
StaffPerformance[] {
  staff_id, staff_name, orders_handled
  total_sales, average_order_value, customer_complaints
}
```

---

### 5. `GetPopularItemsReport`

**Request:** `period`, `from_date`, `to_date`, `top_n` (số món top)  
**Response:**
```
PopularItem[] {
  item_id, item_name, quantity_sold, revenue, category
}
```

---

### 6. `GetRevenueAnalytics`

**Request:** `period`, `from_date`, `to_date`  
**Response:**
```
RevenueAnalytics {
  total_revenue, growth_rate
  trend_data[] ← DailySales[]
  peak_day, peak_hour
}
```

---

### 7. `ExportReport`

**Request:** `report_type`, `format` (PDF/EXCEL/CSV/JSON), `period`, `from_date`, `to_date`  
**Response:** `file_data` (bytes), `file_name`, `content_type`, `success`, `message`

---

## Report Periods (Enum)

| Value | Mô tả |
|-------|-------|
| `PERIOD_DAILY` | Theo ngày |
| `PERIOD_WEEKLY` | Theo tuần |
| `PERIOD_MONTHLY` | Theo tháng |
| `PERIOD_YEARLY` | Theo năm |
| `PERIOD_CUSTOM` | Khoảng tùy chỉnh (dùng `from_date` + `to_date`) |

---

## Export Formats

`FORMAT_PDF`, `FORMAT_EXCEL`, `FORMAT_CSV`, `FORMAT_JSON`

---

## Cấu hình (Environment Variables)

| Biến | Default |
|------|---------|
| `SERVER_PORT` | `50059` |

---

## Vấn đề chưa giải quyết / Cần triển khai

- **Chưa kết nối vào docker-compose** — service chưa chạy trong stack
- **Chưa có HTTP route** trong api-gateway cho report endpoints
- **Database chưa implement** — hiện tại trả về mock/in-memory data
- Cần kết nối tới PostgreSQL và query từ `orders`, `order_items`, `menu_items`, `staff` tables
- Không có authentication/authorization — chỉ ADMIN mới nên có quyền truy cập
- `GetInventoryReport` cần có inventory-service (chưa tồn tại)
- `GetStaffPerformanceReport` cần join staff-service data với order data
- Export (PDF/Excel) chưa implement — cần thư viện Go cho PDF/spreadsheet generation
- AnalyticsOverview page trong admin dashboard chưa gọi report-service (dùng mock data)
