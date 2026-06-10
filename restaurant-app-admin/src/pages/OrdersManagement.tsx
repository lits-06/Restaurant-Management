import React, { useState, useMemo, useEffect } from 'react';
import { ordersApi, type OrderDto, type OrderItemDto } from '../services/api';

interface OrderItem {
    id: string;
    name: string;
    price: number;
    quantity: number;
}

interface Reservation {
    id: string;
    name: string;
    phone: string;
    time: string;
    date: string;
    partySize: number;
    status: 'Confirmed' | 'Pending' | 'Cancelled' | 'Completed';
    items: OrderItem[];
}

const statusOptions = ['Confirmed', 'Pending', 'Completed', 'Cancelled'] as const;

const getOrderId = (order: OrderDto) => order.order_id ?? order.orderId ?? '';
const getOrderItemId = (item: OrderItemDto) => item.item_id ?? item.itemId ?? '';
const getTimestampDate = (value: OrderDto['time']) => {
    if (typeof value === 'string') return new Date(value);
    if (value?.seconds) return new Date(value.seconds * 1000);
    return null;
};
const mapOrderToReservation = (order: OrderDto): Reservation => {
    const date = getTimestampDate(order.time);

    return {
        id: getOrderId(order),
        name: order.name ?? 'Unknown guest',
        phone: order.phone ?? '',
        time: date ? date.toTimeString().slice(0, 5) : '19:00',
        date: date ? date.toISOString().slice(0, 10) : new Date().toISOString().slice(0, 10),
        partySize: order.party_size ?? order.partySize ?? 1,
        status: (order.status as Reservation['status']) || 'Pending',
        items: (order.items ?? []).map((item) => ({
            id: getOrderItemId(item),
            name: item.name ?? 'Menu item',
            price: item.price ?? 0,
            quantity: item.quantity ?? 1,
        })),
    };
};

const toOrderPayload = (reservation: Reservation) => ({
    name: reservation.name,
    phone: reservation.phone,
    time: reservation.time,
    date: reservation.date,
    party_size: reservation.partySize,
    status: reservation.status,
    items: reservation.items.map((item) => ({
        item_id: item.id,
        quantity: item.quantity,
    })),
});

const OrdersManagement: React.FC = () => {
    // 1. State Tìm kiếm & Bộ lọc
    const [searchQuery, setSearchQuery] = useState<string>('');
    const [statusFilter, setStatusFilter] = useState<string>('All Statuses');
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [error, setError] = useState<string>('');

    // 2. State Phân trang (Pagination)
    const [currentPage, setCurrentPage] = useState<number>(1);
    const itemsPerPage = 5; // Số lượng hàng hiển thị trên mỗi trang

    // 2. State Slide-over Drawer (Xem chi tiết đơn hàng)
    const [isDrawerOpen, setIsDrawerOpen] = useState<boolean>(false);
    const [selectedReservation, setSelectedReservation] = useState<Reservation | null>(null);

    // 3. State Modal Chỉnh sửa Đặt chỗ (Hành chính)
    const [isEditModalOpen, setIsEditModalOpen] = useState<boolean>(false);
    const [editingReservation, setEditingReservation] = useState<Reservation | null>(null);

    // 4. State Modal Chỉnh sửa Chi tiết Đơn hàng & Món ăn
    const [isEditOrderModalOpen, setIsEditOrderModalOpen] = useState<boolean>(false);
    const [editingOrder, setEditingOrder] = useState<Reservation | null>(null);

    const [reservations, setReservations] = useState<Reservation[]>([]);

    const loadOrders = async () => {
        setIsLoading(true);
        setError('');

        try {
            const response = await ordersApi.list({
                page: 1,
                page_size: 100,
                status: statusFilter === 'All Statuses' ? '' : statusFilter,
                keyword: searchQuery,
            });
            console.log('Raw orders response:', response);
            setReservations((response.orders ?? []).map(mapOrderToReservation).filter((order) => order.id));
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Không thể tải đơn đặt bàn từ máy chủ.');
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        loadOrders();
    }, []);

    // --- LOGIC LỌC & TÌM KIẾM DỮ LIỆU ---
    const filteredReservations = useMemo(() => {
        // Reset về trang 1 nếu người dùng thay đổi bộ lọc hoặc từ khóa tìm kiếm
        return reservations.filter((res) => {
            const matchesSearch = res.name.toLowerCase().includes(searchQuery.toLowerCase());
            const matchesStatus = statusFilter === 'All Statuses' || res.status === statusFilter;
            return matchesSearch && matchesStatus;
        });
    }, [reservations, searchQuery, statusFilter]);

    // --- LOGIC PHÂN TRANG (PAGINATION) ---
    const totalItems = filteredReservations.length;
    const totalPages = Math.ceil(totalItems / itemsPerPage) || 1;

    // Điều chỉnh trang hiện tại nếu vượt quá tổng số trang sau khi lọc
    const safeCurrentPage = Math.min(currentPage, totalPages);

    const startIndex = (safeCurrentPage - 1) * itemsPerPage;
    const endIndex = Math.min(startIndex + itemsPerPage, totalItems);

    // Cắt mảng dữ liệu chỉ hiển thị cho trang hiện tại
    const paginatedReservations = useMemo(() => {
        return filteredReservations.slice(startIndex, endIndex);
    }, [filteredReservations, startIndex, endIndex]);

    useEffect(() => {
        if (currentPage > totalPages) {
            setCurrentPage(totalPages);
        }
    }, [currentPage, totalPages]);

    // --- Các hàm tính toán tài chính ---
    const calculateSubtotal = (items: OrderItem[]) => items.reduce((sum, item) => sum + item.price * item.quantity, 0);
    const calculateTax = (subtotal: number) => subtotal * 0.08;
    const calculateServiceFee = (subtotal: number) => subtotal * 0.15;
    const calculateTotal = (subtotal: number) => subtotal + calculateTax(subtotal) + calculateServiceFee(subtotal);

    // --- Xử lý sự kiện tăng/giảm/xóa món ăn trực tiếp trên State mẫu ---
    const handleUpdateQty = (itemId: string, delta: number) => {
        if (!editingOrder) return;
        const updatedItems = editingOrder.items
            .map((item) => (item.id === itemId ? { ...item, quantity: Math.max(0, item.quantity + delta) } : item))
            .filter((item) => item.quantity > 0);
        setEditingOrder({ ...editingOrder, items: updatedItems });
    };

    const handleRemoveItem = (itemId: string) => {
        if (!editingOrder) return;
        setEditingOrder({ ...editingOrder, items: editingOrder.items.filter((item) => item.id !== itemId) });
    };

    const saveOrderChanges = () => {
        if (!editingOrder) return;
        ordersApi.update(editingOrder.id, toOrderPayload(editingOrder))
            .then((response) => {
                const updated = response.order ? mapOrderToReservation(response.order) : editingOrder;
                setReservations(reservations.map((res) => (res.id === editingOrder.id ? updated : res)));
                setIsEditOrderModalOpen(false);
            })
            .catch((err) => setError(err instanceof Error ? err.message : 'Không thể lưu chi tiết đơn hàng.'));
    };

    const saveReservationChanges = () => {
        if (!editingReservation) return;
        ordersApi.update(editingReservation.id, toOrderPayload(editingReservation))
            .then((response) => {
                const updated = response.order ? mapOrderToReservation(response.order) : editingReservation;
                setReservations(reservations.map((res) => (res.id === editingReservation.id ? updated : res)));
                setIsEditModalOpen(false);
            })
            .catch((err) => setError(err instanceof Error ? err.message : 'Không thể lưu thông tin đặt chỗ.'));
    };

    return (
        <div className="flex flex-col animate-fadeIn">
            {/* Header */}
            <header className="mb-10 flex flex-col md:flex-row md:items-end justify-between gap-6">
                <div>
                    <nav className="flex gap-2 text-[10px] text-on-surface-variant uppercase tracking-widest mb-2">
                        <span>Admin</span>
                        <span>/</span>
                        <span className="text-primary font-bold">Reservations</span>
                    </nav>
                    <h2 className="font-serif text-5xl font-bold text-on-surface">Upcoming Reservations</h2>
                    <p className="text-on-surface-variant text-sm mt-2 max-w-xl">
                        Manage guest arrivals, modify table assignments, and track seating status for the current service period.
                    </p>
                </div>
            </header>

            {/* Bộ lọc & Tìm kiếm */}
            <section className="bg-white p-6 rounded-xl shadow-sm mb-6 flex flex-wrap items-center gap-4 border border-outline-variant/10">
                <div className="flex-grow relative">
                    <span className="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-on-surface-variant">search</span>
                    <input
                        type="text"
                        value={searchQuery}
                        onChange={(e) => {
                            setSearchQuery(e.target.value);
                            setCurrentPage(1); // Reset về trang 1 khi gõ tìm kiếm
                        }}
                        className="w-full pl-12 pr-4 py-3 bg-[#f3f4f5] border-none rounded-lg focus:ring-2 focus:ring-[#d4af37] text-sm text-on-surface"
                        placeholder="Search by guest name..."
                    />
                </div>
                <div className="flex items-center gap-3">
                    <select
                        value={statusFilter}
                        onChange={(e) => {
                            setStatusFilter(e.target.value);
                            setCurrentPage(1); // Reset về trang 1 khi thay đổi bộ lọc
                        }}
                        className="bg-[#f3f4f5] border-none rounded-lg py-3 px-4 text-xs font-semibold focus:ring-2 focus:ring-[#d4af37] pr-10"
                    >
                        <option>All Statuses</option>
                        {statusOptions.map((status) => (
                            <option key={status}>{status}</option>
                        ))}
                    </select>
                </div>
            </section>

            {error && <p className="mb-4 text-sm text-red-600">{error}</p>}
            {isLoading && <p className="mb-4 text-sm text-on-surface-variant">Đang tải danh sách từ order-service...</p>}

            {/* Bảng dữ liệu Đặt bàn */}
            <div className="bg-white rounded-xl shadow-sm overflow-hidden border border-outline-variant/20">
                <table className="w-full text-left border-collapse">
                    <thead>
                        <tr className="bg-[#f3f4f5] border-b border-outline-variant/30">
                            <th className="px-6 py-4 text-xs font-semibold text-on-surface-variant uppercase">Guest Details</th>
                            <th className="px-6 py-4 text-xs font-semibold text-on-surface-variant uppercase">Time & Date</th>
                            <th className="px-6 py-4 text-xs font-semibold text-on-surface-variant uppercase">Party</th>
                            <th className="px-6 py-4 text-xs font-semibold text-on-surface-variant uppercase">Status</th>
                            <th className="px-6 py-4 text-xs font-semibold text-on-surface-variant uppercase text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-outline-variant/10">
                        {paginatedReservations.length === 0 ? (
                            <tr>
                                <td colSpan={5} className="text-center py-10 text-sm text-on-surface-variant italic">
                                    Không tìm thấy lịch đặt bàn nào phù hợp.
                                </td>
                            </tr>
                        ) : (
                            paginatedReservations.map((res) => (
                                <tr key={res.id} className="hover:bg-[#f3f4f5]/50 transition-colors">
                                    <td className="px-6 py-5">
                                        <div className="flex items-center gap-3">
                                            <div>
                                                <p className="text-sm font-bold text-on-surface">{res.name}</p>
                                                <p className="text-[12px] text-on-surface-variant">{res.phone}</p>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-6 py-5">
                                        <p className="text-sm text-on-surface">{res.time}</p>
                                        <p className="text-[12px] text-on-surface-variant italic">{res.date}</p>
                                    </td>
                                    <td className="px-6 py-5">
                                        <div className="flex items-center gap-1 text-on-surface">
                                            <span className="material-symbols-outlined text-base">groups</span>
                                            <span className="text-sm">{res.partySize}</span>
                                        </div>
                                    </td>
                                    <td className="px-6 py-5">
                                        <span className={`text-[10px] px-3 py-1 rounded-full font-bold uppercase ${res.status === 'Confirmed' || res.status === 'Completed' ? 'bg-green-100 text-green-800' : res.status === 'Cancelled' ? 'bg-red-100 text-red-800' : 'bg-amber-100 text-amber-800'}`}>
                                            {res.status}
                                        </span>
                                    </td>
                                    <td className="px-6 py-5 text-right whitespace-nowrap">
                                        <button onClick={() => { setEditingReservation(res); setIsEditModalOpen(true); }} className="p-2 text-on-surface-variant hover:text-[#735c00]" title="Sửa thông tin đặt chỗ">
                                            <span className="material-symbols-outlined text-xl">edit_note</span>
                                        </button>
                                        <button onClick={() => { setEditingOrder(res); setIsEditOrderModalOpen(true); }} className="p-2 text-on-surface-variant hover:text-[#735c00]" title="Sửa chi tiết món ăn">
                                            <span className="material-symbols-outlined text-xl">restaurant</span>
                                        </button>
                                        <button onClick={() => { setSelectedReservation(res); setIsDrawerOpen(true); }} className="p-2 text-on-surface-variant hover:text-[#735c00]" title="Xem chi tiết đơn">
                                            <span className="material-symbols-outlined text-xl">more_vert</span>
                                        </button>
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>

            {/* Pagination */}
            <div className="px-6 py-4 flex items-center justify-between bg-[#f3f4f5] border-t border-outline-variant/20">
                <p className="text-xs font-semibold text-on-surface-variant">
                    Showing {totalItems === 0 ? 0 : startIndex + 1} to {endIndex} of {totalItems} Reservations
                </p>

                <div className="flex gap-1">
                    <button
                        onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
                        disabled={safeCurrentPage === 1}
                        className="p-2 rounded hover:bg-[#e1e3e4] transition-colors disabled:opacity-30 disabled:hover:bg-transparent"
                    >
                        <span className="material-symbols-outlined text-sm">
                            chevron_left
                        </span>
                    </button>

                    {Array.from({ length: totalPages }, (_, index) => {
                        const pageNumber = index + 1;

                        return (
                            <button
                                key={pageNumber}
                                onClick={() => setCurrentPage(pageNumber)}
                                className={`w-8 h-8 rounded text-xs font-semibold transition-colors ${safeCurrentPage === pageNumber
                                        ? 'bg-[#735c00] text-white'
                                        : 'hover:bg-[#e1e3e4] text-on-surface'
                                    }`}
                            >
                                {pageNumber}
                            </button>
                        );
                    })}

                    <button
                        onClick={() =>
                            setCurrentPage((prev) => Math.min(prev + 1, totalPages))
                        }
                        disabled={safeCurrentPage === totalPages}
                        className="p-2 rounded hover:bg-[#e1e3e4] transition-colors disabled:opacity-30 disabled:hover:bg-transparent"
                    >
                        <span className="material-symbols-outlined text-sm">
                            chevron_right
                        </span>
                    </button>
                </div>
            </div>

            {/* --- 1. SLIDE-OVER DRAWER: XEM TỔNG HỢP ĐƠN HÀNG --- */}
            {isDrawerOpen && selectedReservation && (
                <div className="fixed inset-0 z-50 flex justify-end">
                    <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" onClick={() => setIsDrawerOpen(false)} />
                    <div className="relative h-full w-full max-w-md bg-white shadow-2xl flex flex-col animate-slideLeft">
                        <div className="p-6 border-b flex items-center justify-between bg-[#f3f4f5]">
                            <div className="flex items-center gap-3">
                                <div>
                                    <div className="flex items-center gap-2">
                                        <h3 className="font-serif text-xl font-bold">{selectedReservation.name}</h3>
                                    </div>
                                    <p className="text-on-surface-variant text-xs">• {selectedReservation.partySize} Guests</p>
                                </div>
                            </div>
                            <button className="p-2 hover:bg-[#e1e3e4] rounded-full" onClick={() => setIsDrawerOpen(false)}>
                                <span className="material-symbols-outlined">close</span>
                            </button>
                        </div>

                        <div className="flex-grow overflow-y-auto p-6 space-y-6">
                            <div>
                                <label className="block text-xs font-semibold text-on-surface-variant uppercase mb-2">Trạng thái đặt món</label>
                                <div className="p-3 bg-[#f8f9fa] rounded-lg text-sm font-bold text-primary">{selectedReservation.status}</div>
                            </div>

                            <div>
                                <h4 className="text-xs font-semibold text-on-surface-variant uppercase mb-4">Danh sách món đã gọi</h4>
                                <div className="space-y-4">
                                    {selectedReservation.items.map((item) => (
                                        <div key={item.id} className="flex items-center justify-between border-b border-gray-100 pb-2">
                                            <div>
                                                <p className="text-sm font-bold text-on-surface">{item.name}</p>
                                                <p className="text-on-surface-variant text-xs">SL: {item.quantity} × ${item.price.toFixed(2)}</p>
                                            </div>
                                            <p className="text-sm font-bold">${(item.quantity * item.price).toFixed(2)}</p>
                                        </div>
                                    ))}
                                </div>
                            </div>

                            <div className="bg-[#f3f4f5] p-4 rounded-lg space-y-2">
                                <div className="flex justify-between text-xs text-on-surface-variant"><span>Tạm tính</span><span>${calculateSubtotal(selectedReservation.items).toFixed(2)}</span></div>
                                <div className="flex justify-between text-xs text-on-surface-variant"><span>Thuế (8%)</span><span>${calculateTax(calculateSubtotal(selectedReservation.items)).toFixed(2)}</span></div>
                                <div className="flex justify-between text-xs text-on-surface-variant"><span>Phí phục vụ (15%)</span><span>${calculateServiceFee(calculateSubtotal(selectedReservation.items)).toFixed(2)}</span></div>
                                <hr className="border-outline-variant/30 my-2" />
                                <div className="flex justify-between items-center">
                                    <span className="font-bold text-on-surface text-sm">Tổng cộng</span>
                                    <span className="text-xl font-serif font-bold text-[#735c00]">${calculateTotal(calculateSubtotal(selectedReservation.items)).toFixed(2)}</span>
                                </div>
                            </div>
                        </div>

                        <div className="p-6 border-t bg-white">
                            <button className="w-full bg-[#735c00] text-white text-xs font-semibold py-4 rounded-lg" onClick={() => setIsDrawerOpen(false)}>Đóng hóa đơn</button>
                        </div>
                    </div>
                </div>
            )}

            {/* --- 2. MODAL: SỬA HÀNH CHÍNH ĐẶT CHỖ --- */}
            {isEditModalOpen && editingReservation && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
                    <div className="bg-white rounded-xl shadow-2xl w-full max-w-md p-8 animate-scaleUp">
                        <div className="flex items-center justify-between mb-6">
                            <h3 className="font-serif text-2xl font-bold text-on-surface">Sửa Thông Tin Đặt Chỗ</h3>
                            <button className="p-1 hover:bg-[#f3f4f5] rounded-full" onClick={() => setIsEditModalOpen(false)}>
                                <span className="material-symbols-outlined">close</span>
                            </button>
                        </div>
                        <div className="space-y-4">
                            <div>
                                <label className="block text-xs font-semibold text-on-surface-variant mb-1">Tên khách hàng</label>
                                <input className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" type="text" value={editingReservation.name} onChange={(e) => setEditingReservation({ ...editingReservation, name: e.target.value })} />
                            </div>
                            <div>
                                <label className="block text-xs font-semibold text-on-surface-variant mb-1">Số điện thoại / Ghi chú</label>
                                <input className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" type="text" value={editingReservation.phone} onChange={(e) => setEditingReservation({ ...editingReservation, phone: e.target.value })} />
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-xs font-semibold text-on-surface-variant mb-1">Số lượng khách</label>
                                    <input className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" type="number" value={editingReservation.partySize} onChange={(e) => setEditingReservation({ ...editingReservation, partySize: parseInt(e.target.value) || 1 })} />
                                </div>
                                <div>
                                    <label className="block text-xs font-semibold text-on-surface-variant mb-1">Giờ đặt bàn</label>
                                    <input className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" type="time" value={editingReservation.time} onChange={(e) => setEditingReservation({ ...editingReservation, time: e.target.value })} />
                                </div>
                            </div>
                            <div className="flex gap-4 mt-6">
                                <button className="flex-1 border border-outline-variant text-xs py-3 rounded-lg" onClick={() => setIsEditModalOpen(false)}>Hủy</button>
                                <button className="flex-1 bg-[#735c00] text-white text-xs py-3 rounded-lg" onClick={saveReservationChanges}>Lưu thay đổi</button>
                            </div>
                        </div>
                    </div>
                </div>
            )}

            {/* --- 3. MODAL: SỬA CHI TIẾT ĐƠN HÀNG & MÓN ĂN --- */}
            {isEditOrderModalOpen && editingOrder && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
                    <div className="bg-white rounded-xl shadow-2xl w-full max-w-xl max-h-[90vh] overflow-y-auto p-8 animate-scaleUp">
                        <div className="flex items-center justify-between mb-6">
                            <h3 className="font-serif text-2xl font-bold text-on-surface">Edit Order Details</h3>
                            <button className="p-1 hover:bg-[#f3f4f5] rounded-full" onClick={() => setIsEditOrderModalOpen(false)}>
                                <span className="material-symbols-outlined">close</span>
                            </button>
                        </div>

                        <div className="space-y-6">
                            {/* Thẻ thông tin nhanh */}
                            <div className="grid grid-cols-2 gap-4 p-4 bg-[#f3f4f5] rounded-xl text-sm">
                                <div className="col-span-2">
                                    <label className="block text-[11px] font-bold text-on-surface-variant uppercase">Guest Name</label>
                                    <p className="font-bold text-on-surface">{editingOrder.name}</p>
                                </div>
                                <div>
                                    <label className="block text-[11px] font-bold text-on-surface-variant uppercase">Reservation</label>
                                    <p className="text-on-surface">{editingOrder.date} • {editingOrder.time}</p>
                                </div>
                                <div>
                                    <label className="block text-[11px] font-bold text-on-surface-variant uppercase">Trạng thái bàn</label>
                                    <select
                                        className="w-full mt-1 bg-white border border-outline-variant/30 rounded-md p-1.5 text-xs"
                                        value={editingOrder.status}
                                        onChange={(e) => setEditingOrder({ ...editingOrder, status: e.target.value as Reservation['status'] })}
                                    >
                                        {statusOptions.map((status) => (
                                            <option key={status} value={status}>{status}</option>
                                        ))}
                                    </select>
                                </div>
                            </div>

                            {/* Khu vực quản lý danh sách món */}
                            <div className="space-y-4">
                                <h4 className="text-xs font-semibold text-on-surface-variant uppercase tracking-widest">Order Items</h4>
                                <div className="space-y-3 max-h-48 overflow-y-auto pr-1">
                                    {editingOrder.items.length === 0 ? (
                                        <p className="text-center text-xs text-on-surface-variant py-4 italic">Chưa có món ăn nào trong thực đơn pre-order.</p>
                                    ) : (
                                        editingOrder.items.map((item) => (
                                            <div key={item.id} className="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                                                <div className="flex-grow">
                                                    <p className="text-sm font-bold text-on-surface">{item.name}</p>
                                                    <p className="text-on-surface-variant text-xs">${item.price.toFixed(2)} ea</p>
                                                </div>
                                                <div className="flex items-center gap-4">
                                                    <div className="flex items-center border border-outline-variant rounded-lg overflow-hidden bg-white">
                                                        <button type="button" className="px-3 py-1 hover:bg-gray-100" onClick={() => handleUpdateQty(item.id, -1)}>-</button>
                                                        <span className="px-4 py-1 font-mono text-sm border-x border-outline-variant">{item.quantity}</span>
                                                        <button type="button" className="px-3 py-1 hover:bg-gray-100" onClick={() => handleUpdateQty(item.id, 1)}>+</button>
                                                    </div>
                                                    <button type="button" className="text-red-600 hover:opacity-70" onClick={() => handleRemoveItem(item.id)}>
                                                        <span className="material-symbols-outlined">delete</span>
                                                    </button>
                                                </div>
                                            </div>
                                        ))
                                    )}
                                </div>
                            </div>

                            {/* Thanh thêm món nhanh giả lập */}
                            <div className="p-4 bg-[#ffe088]/10 rounded-xl border border-[#ffe088]/30">
                                <label className="block text-xs font-bold text-[#574500] mb-2 uppercase">Add New Item</label>
                                <div className="flex gap-2">
                                    <input className="flex-grow pl-3 pr-4 py-2 bg-white border border-outline-variant rounded-lg text-xs" placeholder="Nhập tên món ăn mới gọi thêm tại bàn..." type="text" id="itemNameInput" />
                                    <button
                                        type="button"
                                        className="bg-[#735c00] text-white px-4 py-2 rounded-lg text-xs font-semibold"
                                        onClick={() => {
                                            const input = document.getElementById('itemNameInput') as HTMLInputElement;
                                            if (input && input.value.trim()) {
                                                const newItem: OrderItem = { id: input.value.trim(), name: input.value.trim(), price: 25.0, quantity: 1 };
                                                setEditingOrder({ ...editingOrder, items: [...editingOrder.items, newItem] });
                                                input.value = '';
                                            }
                                        }}
                                    >
                                        Add
                                    </button>
                                </div>
                            </div>
                        </div>

                        <div className="flex gap-4 mt-8">
                            <button className="flex-1 border border-outline-variant text-xs py-3 rounded-lg" onClick={() => setIsEditOrderModalOpen(false)}>Cancel</button>
                            <button className="flex-1 bg-[#735c00] text-white text-xs py-3 rounded-lg" onClick={saveOrderChanges}>Save Changes</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default OrdersManagement;
