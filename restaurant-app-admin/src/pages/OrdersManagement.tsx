import React, { useState, useMemo, useEffect } from 'react';
import { menuApi, ordersApi, tablesApi, type MenuItemDto, type OrderDto, type OrderItemDto, type TableDto } from '../services/api';

// ── types ──────────────────────────────────────────────────────────────────────

interface OrderItem {
  id: string;
  name: string;
  price: number;
  quantity: number;
  itemStatus: string;
}

interface Reservation {
  id: string;
  tableId: string;
  userId: string;
  name: string;
  phone: string;
  notes: string;
  date: string;
  time: string;
  endTime: string;
  partySize: number;
  status: 'Confirmed' | 'Pending' | 'Cancelled' | 'Completed';
  items: OrderItem[];
  total: number;
}

// ── helpers ────────────────────────────────────────────────────────────────────

const STATUS_OPTIONS = ['Confirmed', 'Pending', 'Completed', 'Cancelled'] as const;

const getOrderId     = (o: OrderDto)     => o.order_id ?? o.orderId ?? '';
const getOrderItemId = (i: OrderItemDto) => i.item_id  ?? i.itemId  ?? '';

const parseTs = (v: OrderDto['time']): Date | null => {
  if (!v) return null;
  if (typeof v === 'string') return new Date(v);
  if (typeof v === 'object' && v.seconds) return new Date(v.seconds * 1000);
  return null;
};

const fmtTime     = (d: Date | null) => d ? d.toTimeString().slice(0, 5) : '';
const fmtDate     = (d: Date | null) => d ? d.toISOString().slice(0, 10) : new Date().toISOString().slice(0, 10);
const fmtVnd      = (n: number)      => `${Math.round(n).toLocaleString('vi-VN')}đ`;
const shortId     = (id: string)     => id.slice(0, 8);

const STATUS_STYLE: Record<string, string> = {
  Confirmed: 'bg-green-100 text-green-800',
  Completed: 'bg-green-100 text-green-800',
  Cancelled: 'bg-red-100 text-red-800',
  Pending:   'bg-amber-100 text-amber-800',
};

const ITEM_STATUS_CHIP: Record<string, string> = {
  PENDING: 'bg-gray-100 text-gray-600',
  COOKING: 'bg-orange-100 text-orange-700',
  READY:   'bg-green-100 text-green-700',
  SERVED:  'bg-blue-100 text-blue-700',
};
const ITEM_STATUS_LABEL: Record<string, string> = {
  PENDING: 'Chờ', COOKING: 'Đang nấu', READY: 'Xong', SERVED: 'Đã mang',
};

const mapOrderToReservation = (o: OrderDto): Reservation => {
  const timeDate = parseTs(o.time);
  const endDate  = parseTs(o.end_time);
  const apiTotal = o.total ?? o.total_price ?? o.totalPrice ?? 0;
  const itemsSum = (o.items ?? []).reduce((s, i) => s + (i.price ?? 0) * (i.quantity ?? 1), 0);
  return {
    id:        getOrderId(o),
    tableId:   o.table_id ?? '',
    userId:    o.user_id ?? '',
    name:      o.name ?? 'Khách vãng lai',
    phone:     o.phone ?? '',
    notes:     o.notes ?? '',
    date:      fmtDate(timeDate),
    time:      fmtTime(timeDate),
    endTime:   fmtTime(endDate),
    partySize: o.party_size ?? o.partySize ?? 1,
    status:    (o.status as Reservation['status']) || 'Pending',
    items:     (o.items ?? []).map(i => ({
      id:         getOrderItemId(i),
      name:       i.name ?? 'Món ăn',
      price:      i.price ?? 0,
      quantity:   i.quantity ?? 1,
      itemStatus: (i.item_status ?? 'PENDING').toUpperCase(),
    })),
    total: apiTotal > 0 ? apiTotal : itemsSum,
  };
};

// ── component ──────────────────────────────────────────────────────────────────

const OrdersManagement: React.FC = () => {
  // ── data state ────────────────────────────────────────────────────────────────
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [menuItems, setMenuItems]       = useState<MenuItemDto[]>([]);
  const [tables, setTables]             = useState<TableDto[]>([]);
  const [loading, setLoading]           = useState(true);
  const [error, setError]               = useState('');

  // ── filter/pagination ─────────────────────────────────────────────────────────
  const [searchQuery, setSearchQuery]   = useState('');
  const [statusFilter, setStatusFilter] = useState('All Statuses');
  const [currentPage, setCurrentPage]   = useState(1);
  const PER_PAGE = 5;

  // ── drawers / modals ──────────────────────────────────────────────────────────
  const [drawerRes, setDrawerRes]       = useState<Reservation | null>(null);
  const [editModal, setEditModal]       = useState<Reservation | null>(null);
  const [orderModal, setOrderModal]     = useState<Reservation | null>(null);
  const [addSearch, setAddSearch]       = useState('');
  const [statusBusy, setStatusBusy]     = useState<string | null>(null);

  // ── load ───────────────────────────────────────────────────────────────────────
  const loadOrders = async () => {
    setLoading(true);
    setError('');
    try {
      const res = await ordersApi.list({ page: 1, page_size: 100 });
      setReservations((res.orders ?? []).map(mapOrderToReservation).filter(r => r.id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không thể tải đơn đặt bàn.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadOrders();
    menuApi.listItems({ page: 1, page_size: 200 }).then(r => setMenuItems(r.items ?? [])).catch(() => {});
    tablesApi.list({ page_size: 100 }).then(r => setTables(r.tables ?? [])).catch(() => {});
  }, []);

  // ── filtering / pagination ────────────────────────────────────────────────────
  const filtered = useMemo(() =>
    reservations.filter(r =>
      r.name.toLowerCase().includes(searchQuery.toLowerCase()) &&
      (statusFilter === 'All Statuses' || r.status === statusFilter)
    ),
    [reservations, searchQuery, statusFilter]
  );

  const totalPages = Math.max(1, Math.ceil(filtered.length / PER_PAGE));
  const safePage   = Math.min(currentPage, totalPages);
  const pageSlice  = filtered.slice((safePage - 1) * PER_PAGE, safePage * PER_PAGE);

  useEffect(() => { if (currentPage > totalPages) setCurrentPage(totalPages); }, [currentPage, totalPages]);

  // ── table name lookup ─────────────────────────────────────────────────────────
  const tableLabel = (tableId: string) => {
    if (!tableId) return '—';
    const t = tables.find(t => t.table_id === tableId);
    return t?.table_number != null ? `Bàn ${t.table_number}` : shortId(tableId);
  };

  // ── status update ─────────────────────────────────────────────────────────────
  const handleStatusUpdate = async (res: Reservation, newStatus: string) => {
    setStatusBusy(res.id);
    try {
      await ordersApi.updateStatus(res.id, newStatus);
      const updated = { ...res, status: newStatus as Reservation['status'] };
      setReservations(prev => prev.map(r => r.id === res.id ? updated : r));
      if (drawerRes?.id === res.id) setDrawerRes(updated);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không thể cập nhật trạng thái.');
    } finally {
      setStatusBusy(null);
    }
  };

  // ── edit booking ──────────────────────────────────────────────────────────────
  const saveBooking = async () => {
    if (!editModal) return;
    try {
      const res = await ordersApi.update(editModal.id, {
        name:       editModal.name,
        phone:      editModal.phone,
        notes:      editModal.notes,
        date:       editModal.date,
        time:       editModal.time,
        end_time:   editModal.endTime || undefined,
        party_size: editModal.partySize,
        items:      editModal.items.map(i => ({ item_id: i.id, quantity: i.quantity })),
      });
      const updated = res.order ? mapOrderToReservation(res.order) : editModal;
      setReservations(prev => prev.map(r => r.id === editModal.id ? updated : r));
      setEditModal(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không thể lưu.');
    }
  };

  // ── edit order items ──────────────────────────────────────────────────────────
  const handleQty = (itemId: string, delta: number) => {
    if (!orderModal) return;
    setOrderModal({
      ...orderModal,
      items: orderModal.items
        .map(i => i.id === itemId ? { ...i, quantity: Math.max(0, i.quantity + delta) } : i)
        .filter(i => i.quantity > 0),
    });
  };

  const handleAddItem = (menuItem: MenuItemDto) => {
    if (!orderModal) return;
    const id = menuItem.item_id ?? menuItem.itemId ?? '';
    const existing = orderModal.items.find(i => i.id === id);
    if (existing) {
      setOrderModal({ ...orderModal, items: orderModal.items.map(i => i.id === id ? { ...i, quantity: i.quantity + 1 } : i) });
    } else {
      setOrderModal({ ...orderModal, items: [...orderModal.items, { id, name: menuItem.name ?? '', price: menuItem.price ?? 0, quantity: 1, itemStatus: 'PENDING' }] });
    }
    setAddSearch('');
  };

  const saveOrderItems = async () => {
    if (!orderModal) return;
    try {
      const res = await ordersApi.update(orderModal.id, {
        name:       orderModal.name,
        phone:      orderModal.phone,
        notes:      orderModal.notes,
        date:       orderModal.date,
        time:       orderModal.time,
        end_time:   orderModal.endTime || undefined,
        party_size: orderModal.partySize,
        items:      orderModal.items.map(i => ({ item_id: i.id, quantity: i.quantity })),
      });
      const updated = res.order ? mapOrderToReservation(res.order) : orderModal;
      setReservations(prev => prev.map(r => r.id === orderModal.id ? updated : r));
      setOrderModal(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Không thể lưu.');
    }
  };

  const menuSuggestions = useMemo(() => {
    if (!addSearch.trim()) return [];
    return menuItems.filter(i => (i.name ?? '').toLowerCase().includes(addSearch.toLowerCase())).slice(0, 8);
  }, [menuItems, addSearch]);

  // ── render ────────────────────────────────────────────────────────────────────
  return (
    <div className="flex flex-col animate-fadeIn">
      {/* Header */}
      <header className="mb-10 flex md:items-end justify-between gap-6">
        <div>
          <nav className="flex gap-2 text-[10px] text-on-surface-variant uppercase tracking-widest mb-2">
            <span>Admin</span><span>/</span><span className="text-primary font-bold">Đặt bàn</span>
          </nav>
          <h2 className="font-serif text-5xl font-bold text-on-surface">Quản lý Đặt bàn</h2>
          <p className="text-on-surface-variant text-sm mt-2">Xem, xác nhận và quản lý đơn đặt bàn của khách.</p>
        </div>
        <button onClick={loadOrders} className="flex items-center gap-2 px-4 py-2 bg-white border border-outline-variant/30 rounded-lg text-xs font-semibold hover:bg-[#f3f4f5] transition-all whitespace-nowrap">
          <span className="material-symbols-outlined text-base">refresh</span>
          Làm mới
        </button>
      </header>

      {/* Filter bar */}
      <section className="bg-white p-6 rounded-xl shadow-sm mb-6 flex flex-wrap items-center gap-4 border border-outline-variant/10">
        <div className="flex-grow relative">
          <span className="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-on-surface-variant">search</span>
          <input
            type="text"
            value={searchQuery}
            onChange={e => { setSearchQuery(e.target.value); setCurrentPage(1); }}
            className="w-full pl-12 pr-4 py-3 bg-[#f3f4f5] border-none rounded-lg focus:ring-2 focus:ring-[#d4af37] text-sm"
            placeholder="Tìm theo tên khách..."
          />
        </div>
        <select
          value={statusFilter}
          onChange={e => { setStatusFilter(e.target.value); setCurrentPage(1); }}
          className="bg-[#f3f4f5] border-none rounded-lg py-3 px-4 text-xs font-semibold focus:ring-2 focus:ring-[#d4af37]"
        >
          <option>All Statuses</option>
          {STATUS_OPTIONS.map(s => <option key={s}>{s}</option>)}
        </select>
      </section>

      {error  && <p className="mb-4 text-sm text-red-600">{error}</p>}
      {loading && <p className="mb-4 text-sm text-on-surface-variant">Đang tải...</p>}

      {/* Table */}
      <div className="bg-white rounded-xl shadow-sm overflow-hidden border border-outline-variant/20">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-[#f3f4f5] border-b border-outline-variant/30">
              <th className="px-5 py-4 text-xs font-semibold text-on-surface-variant uppercase">Khách hàng</th>
              <th className="px-5 py-4 text-xs font-semibold text-on-surface-variant uppercase">Thời gian</th>
              <th className="px-5 py-4 text-xs font-semibold text-on-surface-variant uppercase">Bàn / Khách</th>
              <th className="px-5 py-4 text-xs font-semibold text-on-surface-variant uppercase">Trạng thái</th>
              <th className="px-5 py-4 text-xs font-semibold text-on-surface-variant uppercase text-right">Thao tác</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-outline-variant/10">
            {pageSlice.length === 0 ? (
              <tr><td colSpan={5} className="text-center py-12 text-sm text-on-surface-variant italic">Không có đơn nào phù hợp.</td></tr>
            ) : pageSlice.map(res => (
              <tr key={res.id} className="hover:bg-[#f3f4f5]/50 transition-colors">
                <td className="px-5 py-4">
                  <p className="text-sm font-bold text-on-surface">{res.name}</p>
                  <p className="text-[12px] text-on-surface-variant">{res.phone}</p>
                  {res.notes && (
                    <p className="text-[11px] text-[#735c00] mt-0.5 italic truncate max-w-[200px]" title={res.notes}>
                      📝 {res.notes}
                    </p>
                  )}
                </td>
                <td className="px-5 py-4">
                  <p className="text-sm text-on-surface font-medium">
                    {res.time}{res.endTime ? ` – ${res.endTime}` : ''}
                  </p>
                  <p className="text-[12px] text-on-surface-variant italic">{res.date}</p>
                </td>
                <td className="px-5 py-4">
                  <p className="text-xs font-semibold text-on-surface">{tableLabel(res.tableId)}</p>
                  <div className="flex items-center gap-1 text-on-surface-variant mt-0.5">
                    <span className="material-symbols-outlined text-sm">groups</span>
                    <span className="text-xs">{res.partySize} khách</span>
                  </div>
                </td>
                <td className="px-5 py-4">
                  <span className={`text-[10px] px-3 py-1 rounded-full font-bold uppercase ${STATUS_STYLE[res.status] ?? 'bg-gray-100 text-gray-700'}`}>
                    {res.status}
                  </span>
                </td>
                <td className="px-5 py-4 text-right whitespace-nowrap">
                  <button onClick={() => setEditModal({ ...res })} className="p-2 text-on-surface-variant hover:text-[#735c00]" title="Sửa thông tin đặt chỗ">
                    <span className="material-symbols-outlined text-xl">edit_note</span>
                  </button>
                  <button onClick={() => { setOrderModal({ ...res }); setAddSearch(''); }} className="p-2 text-on-surface-variant hover:text-[#735c00]" title="Sửa món ăn">
                    <span className="material-symbols-outlined text-xl">restaurant</span>
                  </button>
                  <button onClick={() => setDrawerRes(res)} className="p-2 text-on-surface-variant hover:text-[#735c00]" title="Xem chi tiết">
                    <span className="material-symbols-outlined text-xl">more_vert</span>
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {/* Pagination */}
        <div className="px-6 py-4 flex items-center justify-between bg-[#f3f4f5] border-t border-outline-variant/20">
          <p className="text-xs font-semibold text-on-surface-variant">
            {filtered.length === 0 ? '0' : `${(safePage - 1) * PER_PAGE + 1}–${Math.min(safePage * PER_PAGE, filtered.length)}`} / {filtered.length} đơn
          </p>
          <div className="flex gap-1">
            <button onClick={() => setCurrentPage(p => Math.max(p - 1, 1))} disabled={safePage === 1} className="p-2 rounded hover:bg-[#e1e3e4] disabled:opacity-30">
              <span className="material-symbols-outlined text-sm">chevron_left</span>
            </button>
            {Array.from({ length: totalPages }, (_, i) => i + 1).map(p => (
              <button key={p} onClick={() => setCurrentPage(p)} className={`w-8 h-8 rounded text-xs font-semibold ${safePage === p ? 'bg-[#735c00] text-white' : 'hover:bg-[#e1e3e4]'}`}>{p}</button>
            ))}
            <button onClick={() => setCurrentPage(p => Math.min(p + 1, totalPages))} disabled={safePage === totalPages} className="p-2 rounded hover:bg-[#e1e3e4] disabled:opacity-30">
              <span className="material-symbols-outlined text-sm">chevron_right</span>
            </button>
          </div>
        </div>
      </div>

      {/* ────────────────────────────────────────────────────────────────────────
          DRAWER: Chi tiết đơn
      ──────────────────────────────────────────────────────────────────────── */}
      {drawerRes && (
        <div className="fixed inset-0 z-50 flex justify-end">
          <div className="absolute inset-0 bg-black/40 backdrop-blur-sm" onClick={() => setDrawerRes(null)} />
          <div className="relative h-full w-full max-w-md bg-white shadow-2xl flex flex-col">
            {/* Drawer header */}
            <div className="p-6 border-b flex items-start justify-between bg-[#f3f4f5]">
              <div>
                <h3 className="font-serif text-xl font-bold">{drawerRes.name}</h3>
                <p className="text-on-surface-variant text-xs mt-0.5">{drawerRes.phone} • {drawerRes.partySize} khách</p>
                <span className={`mt-2 inline-block text-[10px] px-3 py-1 rounded-full font-bold uppercase ${STATUS_STYLE[drawerRes.status] ?? ''}`}>
                  {drawerRes.status}
                </span>
              </div>
              <button className="p-2 hover:bg-[#e1e3e4] rounded-full mt-1" onClick={() => setDrawerRes(null)}>
                <span className="material-symbols-outlined">close</span>
              </button>
            </div>

            <div className="flex-grow overflow-y-auto p-6 space-y-5">
              {/* Booking info grid */}
              <div className="grid grid-cols-2 gap-3 text-xs">
                <div className="bg-[#f8f9fa] rounded-lg p-3">
                  <p className="text-on-surface-variant font-semibold uppercase mb-1">Ngày</p>
                  <p className="font-bold text-on-surface">{drawerRes.date}</p>
                </div>
                <div className="bg-[#f8f9fa] rounded-lg p-3">
                  <p className="text-on-surface-variant font-semibold uppercase mb-1">Giờ</p>
                  <p className="font-bold text-on-surface">{drawerRes.time}{drawerRes.endTime ? ` – ${drawerRes.endTime}` : ''}</p>
                </div>
                <div className="bg-[#f8f9fa] rounded-lg p-3">
                  <p className="text-on-surface-variant font-semibold uppercase mb-1">Bàn</p>
                  <p className="font-bold text-on-surface">{tableLabel(drawerRes.tableId)}</p>
                  {drawerRes.tableId && <p className="font-mono text-[10px] text-on-surface-variant mt-0.5">{drawerRes.tableId}</p>}
                </div>
                <div className="bg-[#f8f9fa] rounded-lg p-3">
                  <p className="text-on-surface-variant font-semibold uppercase mb-1">Mã đơn</p>
                  <p className="font-mono text-[11px] text-on-surface break-all">{drawerRes.id}</p>
                </div>
              </div>

              {/* Notes */}
              {drawerRes.notes && (
                <div className="p-3 bg-amber-50 rounded-lg border border-amber-200">
                  <p className="text-xs font-semibold text-amber-800 uppercase mb-1">📝 Ghi chú khách hàng</p>
                  <p className="text-sm text-amber-900">{drawerRes.notes}</p>
                </div>
              )}

              {/* Status action buttons */}
              <div>
                <p className="text-xs font-semibold text-on-surface-variant uppercase mb-2">Cập nhật trạng thái</p>
                <div className="flex gap-2 flex-wrap">
                  {drawerRes.status === 'Pending' && (
                    <button
                      disabled={!!statusBusy}
                      onClick={() => handleStatusUpdate(drawerRes, 'Confirmed')}
                      className="px-4 py-2 bg-green-600 text-white text-xs font-semibold rounded-lg hover:bg-green-700 disabled:opacity-50 transition-colors"
                    >
                      {statusBusy === drawerRes.id ? '...' : '✓ Xác nhận đơn'}
                    </button>
                  )}
                  {drawerRes.status === 'Confirmed' && (
                    <button
                      disabled={!!statusBusy}
                      onClick={() => handleStatusUpdate(drawerRes, 'Completed')}
                      className="px-4 py-2 bg-[#735c00] text-white text-xs font-semibold rounded-lg hover:bg-[#5d4a00] disabled:opacity-50 transition-colors"
                    >
                      {statusBusy === drawerRes.id ? '...' : '✓ Hoàn thành'}
                    </button>
                  )}
                  {(drawerRes.status === 'Pending' || drawerRes.status === 'Confirmed') && (
                    <button
                      disabled={!!statusBusy}
                      onClick={() => handleStatusUpdate(drawerRes, 'Cancelled')}
                      className="px-4 py-2 bg-red-50 text-red-700 border border-red-200 text-xs font-semibold rounded-lg hover:bg-red-100 disabled:opacity-50 transition-colors"
                    >
                      Hủy đơn
                    </button>
                  )}
                </div>
              </div>

              {/* Items list */}
              <div>
                <p className="text-xs font-semibold text-on-surface-variant uppercase mb-3">Danh sách món ({drawerRes.items.length})</p>
                {drawerRes.items.length === 0 ? (
                  <p className="text-xs text-on-surface-variant italic">Chưa có món nào.</p>
                ) : (
                  <div className="space-y-3">
                    {drawerRes.items.map(item => (
                      <div key={item.id} className="flex items-center justify-between pb-2 border-b border-gray-100 last:border-0">
                        <div className="flex-grow">
                          <p className="text-sm font-bold text-on-surface">{item.name}</p>
                          <div className="flex items-center gap-2 mt-0.5">
                            <p className="text-on-surface-variant text-xs">
                              {item.quantity} × {fmtVnd(item.price)}
                            </p>
                            <span className={`text-[10px] px-2 py-0.5 rounded-full font-semibold ${ITEM_STATUS_CHIP[item.itemStatus] ?? 'bg-gray-100 text-gray-600'}`}>
                              {ITEM_STATUS_LABEL[item.itemStatus] ?? item.itemStatus}
                            </span>
                          </div>
                        </div>
                        <p className="text-sm font-bold ml-4">{fmtVnd(item.quantity * item.price)}</p>
                      </div>
                    ))}
                  </div>
                )}
              </div>

              {/* Total */}
              <div className="bg-[#f3f4f5] p-4 rounded-lg">
                <div className="flex justify-between items-center">
                  <span className="font-bold text-on-surface text-sm">Tổng cộng</span>
                  <span className="text-xl font-serif font-bold text-[#735c00]">{fmtVnd(drawerRes.total)}</span>
                </div>
              </div>
            </div>

            <div className="p-6 border-t bg-white">
              <button className="w-full bg-[#735c00] text-white text-xs font-semibold py-3 rounded-lg hover:bg-[#5d4a00] transition-colors" onClick={() => setDrawerRes(null)}>
                Đóng
              </button>
            </div>
          </div>
        </div>
      )}

      {/* ────────────────────────────────────────────────────────────────────────
          MODAL: Sửa thông tin đặt chỗ
      ──────────────────────────────────────────────────────────────────────── */}
      {editModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
          <div className="bg-white rounded-xl shadow-2xl w-full max-w-md p-8">
            <div className="flex items-center justify-between mb-6">
              <h3 className="font-serif text-2xl font-bold text-on-surface">Sửa Thông Tin</h3>
              <button className="p-1 hover:bg-[#f3f4f5] rounded-full" onClick={() => setEditModal(null)}>
                <span className="material-symbols-outlined">close</span>
              </button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">Tên khách hàng</label>
                <input className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={editModal.name} onChange={e => setEditModal({ ...editModal, name: e.target.value })} />
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">Số điện thoại</label>
                <input className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={editModal.phone} onChange={e => setEditModal({ ...editModal, phone: e.target.value })} />
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">Ghi chú đặc biệt</label>
                <textarea className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm resize-none" rows={2} placeholder="Dị ứng, yêu cầu đặc biệt..." value={editModal.notes} onChange={e => setEditModal({ ...editModal, notes: e.target.value })} />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Ngày</label>
                  <input type="date" className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={editModal.date} onChange={e => setEditModal({ ...editModal, date: e.target.value })} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Số khách</label>
                  <input type="number" min={1} className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={editModal.partySize} onChange={e => setEditModal({ ...editModal, partySize: parseInt(e.target.value) || 1 })} />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Giờ bắt đầu</label>
                  <input type="time" className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={editModal.time} onChange={e => setEditModal({ ...editModal, time: e.target.value })} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Giờ kết thúc</label>
                  <input type="time" className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={editModal.endTime} onChange={e => setEditModal({ ...editModal, endTime: e.target.value })} />
                </div>
              </div>
              <div className="flex gap-4 mt-2">
                <button className="flex-1 border border-outline-variant text-xs py-3 rounded-lg" onClick={() => setEditModal(null)}>Hủy</button>
                <button className="flex-1 bg-[#735c00] text-white text-xs py-3 rounded-lg hover:bg-[#5d4a00]" onClick={saveBooking}>Lưu thay đổi</button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* ────────────────────────────────────────────────────────────────────────
          MODAL: Sửa món ăn trong đơn
      ──────────────────────────────────────────────────────────────────────── */}
      {orderModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
          <div className="bg-white rounded-xl shadow-2xl w-full max-w-xl max-h-[90vh] overflow-y-auto p-8">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h3 className="font-serif text-2xl font-bold text-on-surface">Sửa Món Ăn</h3>
                <p className="text-xs text-on-surface-variant">{orderModal.name} — {orderModal.date}</p>
              </div>
              <button className="p-1 hover:bg-[#f3f4f5] rounded-full" onClick={() => setOrderModal(null)}>
                <span className="material-symbols-outlined">close</span>
              </button>
            </div>

            {/* Current items */}
            <div className="space-y-2 mb-5 max-h-52 overflow-y-auto pr-1">
              {orderModal.items.length === 0 ? (
                <p className="text-center text-xs text-on-surface-variant py-4 italic">Chưa có món nào.</p>
              ) : orderModal.items.map(item => (
                <div key={item.id} className="flex items-center justify-between p-3 border border-gray-200 rounded-lg">
                  <div className="flex-grow">
                    <p className="text-sm font-bold text-on-surface">{item.name}</p>
                    <p className="text-on-surface-variant text-xs">{fmtVnd(item.price)}</p>
                  </div>
                  <div className="flex items-center border border-outline-variant rounded-lg overflow-hidden bg-white mx-3">
                    <button type="button" className="px-3 py-1 hover:bg-gray-100" onClick={() => handleQty(item.id, -1)}>−</button>
                    <span className="px-4 py-1 font-mono text-sm border-x border-outline-variant">{item.quantity}</span>
                    <button type="button" className="px-3 py-1 hover:bg-gray-100" onClick={() => handleQty(item.id, 1)}>+</button>
                  </div>
                  <span className="text-sm font-semibold text-[#735c00] w-24 text-right">
                    {fmtVnd(item.quantity * item.price)}
                  </span>
                </div>
              ))}
            </div>

            {/* Add item from menu */}
            <div className="p-4 bg-[#ffe088]/10 rounded-xl border border-[#ffe088]/30">
              <label className="block text-xs font-bold text-[#574500] mb-2 uppercase">Thêm món từ thực đơn</label>
              <div className="relative">
                <input
                  className="w-full px-3 py-2 bg-white border border-outline-variant rounded-lg text-xs"
                  placeholder="Gõ tên món để tìm kiếm..."
                  type="text"
                  value={addSearch}
                  onChange={e => setAddSearch(e.target.value)}
                />
                {menuSuggestions.length > 0 && (
                  <div className="absolute top-full left-0 right-0 mt-1 bg-white border border-outline-variant rounded-lg shadow-lg z-10 max-h-44 overflow-y-auto">
                    {menuSuggestions.map(item => (
                      <button
                        key={item.item_id ?? item.itemId}
                        type="button"
                        className="w-full text-left px-3 py-2.5 hover:bg-[#f3f4f5] flex justify-between items-center text-xs border-b border-gray-50 last:border-0"
                        onClick={() => handleAddItem(item)}
                      >
                        <span className="font-semibold">{item.name}</span>
                        <span className="text-[#735c00] font-semibold ml-4">{fmtVnd(item.price ?? 0)}</span>
                      </button>
                    ))}
                  </div>
                )}
              </div>
            </div>

            {/* Subtotal preview */}
            <div className="mt-4 flex justify-between items-center text-sm font-bold text-on-surface">
              <span>Tạm tính</span>
              <span className="text-[#735c00]">
                {fmtVnd(orderModal.items.reduce((s, i) => s + i.price * i.quantity, 0))}
              </span>
            </div>

            <div className="flex gap-4 mt-5">
              <button className="flex-1 border border-outline-variant text-xs py-3 rounded-lg" onClick={() => setOrderModal(null)}>Hủy</button>
              <button className="flex-1 bg-[#735c00] text-white text-xs py-3 rounded-lg hover:bg-[#5d4a00]" onClick={saveOrderItems}>Lưu thay đổi</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default OrdersManagement;
