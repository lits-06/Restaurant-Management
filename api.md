# API Documentation - Restaurant Management System (Microservices)

## Auth Service
- **Đăng ký tài khoản nhân viên**
- **Đăng nhập hệ thống**
- **Đăng xuất**
- **Refresh token (JWT)**

## User/Staff Service
- **Tạo nhân viên mới (admin)**
- **Cập nhật thông tin nhân viên**
- **Xóa nhân viên**
- **Phân quyền** (admin / thu ngân / phục vụ / bếp)

## Table Service
- **Tạo bàn (table)**
- **Cập nhật trạng thái bàn** (trống / đang dùng / đã đặt)
- **Xem danh sách bàn**

## Menu Service
- **Thêm món ăn**
- **Cập nhật món ăn**
- **Xóa món ăn**
- **Xem menu theo category**

## Order Service
- **Tạo order mới** (khách vào bàn)
- **Thêm món vào order**
- **Cập nhật số lượng món**
- **Xóa món khỏi order**
- **Gửi order xuống bếp**
- **Đóng order** (khi thanh toán xong)

## Payment Service
- **Tính tổng tiền**
- **Thanh toán** (cash / chuyển khoản / QR)
- **Xuất hóa đơn**

## Inventory Service
- **Nhập nguyên liệu**
- **Trừ nguyên liệu khi có order**

## Notification Service
- **Thông báo cho bếp khi có order mới**
- **Thông báo cho phục vụ khi món đã xong**

## Report Service
- **Xem báo cáo doanh thu**
- **Xem thống kê món bán chạy**