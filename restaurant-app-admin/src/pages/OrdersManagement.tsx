import React, { useState, useMemo, useEffect } from 'react';
import { menuApi, ordersApi, tablesApi, type MenuItemDto, type CategoryDto, type OrderDto, type OrderItemDto, type TableDto } from '../services/api';

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

type ViewMode = 'upcoming' | 'today' | 'week' | 'month' | 'all';

const VIEW_TABS: { key: ViewMode; label: string }[] = [
  { key: 'upcoming', label: 'Upcoming' },
  { key: 'today',    label: 'Today' },
  { key: 'week',     label: 'This Week' },
  { key: 'month',    label: 'This Month' },
  { key: 'all',      label: 'All' },
];

interface WalkInState {
  name: string;
  phone: string;
  partySize: number;
  tableId: string;
  date: string;
  time: string;
  endTime: string;
  notes: string;
  status: string;
  items: OrderItem[];
  search: string;
}

const nowHHMM = () => {
  const d = new Date();
  return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
};

const addHoursToTime = (hhmm: string, h: number): string => {
  const [hh, mm] = hhmm.split(':').map(Number);
  const total = hh * 60 + mm + h * 60;
  return `${String(Math.floor(total / 60) % 24).padStart(2, '0')}:${String(total % 60).padStart(2, '0')}`;
};

const getMonWeekBounds = (d: Date): { start: string; end: string } => {
  const day = d.getDay();
  const diff = day === 0 ? -6 : 1 - day;
  const mon = new Date(d); mon.setDate(d.getDate() + diff);
  const sun = new Date(mon); sun.setDate(mon.getDate() + 6);
  return { start: fmtDate(mon), end: fmtDate(sun) };
};

const defaultWalkIn = (): WalkInState => {
  const t = nowHHMM();
  return {
    name: '', phone: '', partySize: 2, tableId: '',
    date: fmtDate(new Date()), time: t, endTime: addHoursToTime(t, 2),
    notes: '', status: 'Confirmed', items: [], search: '',
  };
};

const getOrderId     = (o: OrderDto)     => o.order_id ?? o.orderId ?? '';
const getOrderItemId = (i: OrderItemDto) => i.item_id  ?? i.itemId  ?? '';

const parseTs = (v: OrderDto['time']): Date | null => {
  if (!v) return null;
  if (typeof v === 'string') return new Date(v);
  if (typeof v === 'object' && v.seconds) return new Date(v.seconds * 1000);
  return null;
};

const fmtTime = (d: Date | null) => d ? d.toTimeString().slice(0, 5) : '';
const fmtDate = (d: Date | null) => {
  const target = d ?? new Date();
  const y = target.getFullYear();
  const m = String(target.getMonth() + 1).padStart(2, '0');
  const day = String(target.getDate()).padStart(2, '0');
  return `${y}-${m}-${day}`;
};
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
  PENDING: 'Pending', COOKING: 'Cooking', READY: 'Done', SERVED: 'Served',
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
    name:      o.name ?? 'Walk-in',
    phone:     o.phone ?? '',
    notes:     o.notes ?? '',
    date:      fmtDate(timeDate),
    time:      fmtTime(timeDate),
    endTime:   fmtTime(endDate),
    partySize: o.party_size ?? o.partySize ?? 1,
    status:    (o.status as Reservation['status']) || 'Pending',
    items:     (o.items ?? []).map(i => ({
      id:         getOrderItemId(i),
      name:       i.name ?? 'Item',
      price:      i.price ?? 0,
      quantity:   i.quantity ?? 1,
      itemStatus: (i.item_status ?? 'PENDING').toUpperCase(),
    })),
    total: apiTotal > 0 ? apiTotal : itemsSum,
  };
};

// ── MenuBrowser ────────────────────────────────────────────────────────────────

interface MenuBrowserProps {
  menuItems: MenuItemDto[];
  categories: CategoryDto[];
  currentItems: OrderItem[];
  search: string;
  catFilter: string;
  onSearchChange: (s: string) => void;
  onCatChange: (id: string) => void;
  onAdd: (item: MenuItemDto) => void;
  onQtyChange: (itemId: string, delta: number) => void;
}

const MenuBrowser: React.FC<MenuBrowserProps> = ({
  menuItems, categories, currentItems, search, catFilter,
  onSearchChange, onCatChange, onAdd, onQtyChange,
}) => {
  const itemId = (m: MenuItemDto) => m.item_id ?? m.itemId ?? '';

  const visible = useMemo(() => menuItems.filter(m => {
    if (catFilter && (m.category ?? '') !== catFilter) return false;
    if (search && !(m.name ?? '').toLowerCase().includes(search.toLowerCase())) return false;
    return true;
  }), [menuItems, catFilter, search]);

  const qtyOf = (m: MenuItemDto) => currentItems.find(i => i.id === itemId(m))?.quantity ?? 0;

  return (
    <div>
      <input
        type="text"
        className="w-full px-3 py-2 bg-white border border-outline-variant/30 rounded-lg text-xs mb-2 focus:ring-2 focus:ring-[#d4af37] outline-none"
        placeholder="Search menu items..."
        value={search}
        onChange={e => onSearchChange(e.target.value)}
      />

      {/* Category pills */}
      <div className="flex gap-1 flex-wrap mb-3">
        <button
          type="button"
          onClick={() => onCatChange('')}
          className={`px-2.5 py-1 rounded-full text-[10px] font-semibold transition-all ${!catFilter ? 'bg-[#735c00] text-white' : 'bg-[#f3f4f5] text-[#4d4635] hover:bg-[#e1e3e4]'}`}
        >All</button>
        {categories.map(cat => {
          const name = cat.name ?? '';
          return (
            <button
              key={name} type="button"
              onClick={() => onCatChange(catFilter === name ? '' : name)}
              className={`px-2.5 py-1 rounded-full text-[10px] font-semibold transition-all ${catFilter === name ? 'bg-[#735c00] text-white' : 'bg-[#f3f4f5] text-[#4d4635] hover:bg-[#e1e3e4]'}`}
            >{cat.name}</button>
          );
        })}
      </div>

      {/* Item grid — scrollable */}
      <div className="grid grid-cols-2 gap-2 max-h-56 overflow-y-auto pr-1">
        {visible.length === 0 ? (
          <p className="col-span-2 text-center text-xs text-gray-400 py-6 italic">No items found.</p>
        ) : visible.map(item => {
          const id  = itemId(item);
          const qty = qtyOf(item);
          return (
            <div key={id}
              className="flex flex-col justify-between p-2.5 border border-gray-200 rounded-lg bg-white hover:border-[#d4af37] transition-colors"
            >
              <div className="mb-2">
                <p className="text-xs font-semibold text-on-surface leading-snug line-clamp-2">{item.name}</p>
                <p className="text-[11px] text-[#735c00] font-semibold mt-0.5">{fmtVnd(item.price ?? 0)}</p>
              </div>
              {qty === 0 ? (
                <button type="button"
                  className="w-full py-1 bg-[#f8f9fa] hover:bg-[#ffe088]/40 border border-outline-variant/30 rounded text-xs font-bold text-[#735c00] transition-colors"
                  onClick={() => onAdd(item)}
                >+ Add</button>
              ) : (
                <div className="flex items-center border border-[#d4af37] rounded overflow-hidden">
                  <button type="button" className="px-2.5 py-1 hover:bg-gray-100 text-sm font-bold leading-none" onClick={() => onQtyChange(id, -1)}>−</button>
                  <span className="flex-1 text-center text-xs font-mono font-bold text-[#735c00]">{qty}</span>
                  <button type="button" className="px-2.5 py-1 hover:bg-gray-100 text-sm font-bold leading-none" onClick={() => onQtyChange(id, 1)}>+</button>
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
};

// ── component ──────────────────────────────────────────────────────────────────

const OrdersManagement: React.FC<{ refreshSignal?: number }> = ({ refreshSignal }) => {
  // ── data state ────────────────────────────────────────────────────────────────
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [menuItems, setMenuItems]       = useState<MenuItemDto[]>([]);
  const [categories, setCategories]     = useState<CategoryDto[]>([]);
  const [tables, setTables]             = useState<TableDto[]>([]);
  const [loading, setLoading]           = useState(true);
  const [error, setError]               = useState('');

  // ── filter/pagination ─────────────────────────────────────────────────────────
  const [viewMode, setViewMode]         = useState<ViewMode>('upcoming');
  const [searchQuery, setSearchQuery]   = useState('');
  const [statusFilter, setStatusFilter] = useState('All Statuses');
  const [currentPage, setCurrentPage]   = useState(1);
  const PER_PAGE = 10;

  // ── drawers / modals ──────────────────────────────────────────────────────────
  const [drawerRes, setDrawerRes]       = useState<Reservation | null>(null);
  const [editModal, setEditModal]       = useState<Reservation | null>(null);
  const [orderModal, setOrderModal]     = useState<Reservation | null>(null);
  const [addSearch, setAddSearch]       = useState('');
  const [addCatFilter, setAddCatFilter] = useState('');
  const [walkInCatFilter, setWalkInCatFilter] = useState('');
  const [statusBusy, setStatusBusy]     = useState<string | null>(null);
  const [walkIn, setWalkIn]             = useState<WalkInState | null>(null);
  const [walkInBusy, setWalkInBusy]     = useState(false);
  const [walkInError, setWalkInError]   = useState('');

  // ── load ───────────────────────────────────────────────────────────────────────
  const loadOrders = async () => {
    setLoading(true);
    setError('');
    try {
      const res = await ordersApi.list({ page: 1, page_size: 500, sort_order: 'asc' });
      setReservations((res.orders ?? []).map(mapOrderToReservation).filter(r => r.id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load orders.');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadOrders();
    menuApi.listItems({ page: 1, page_size: 200 }).then(r => setMenuItems(r.items ?? [])).catch(() => {});
    menuApi.listCategories().then(r => setCategories(r.categories ?? [])).catch(() => {});
    tablesApi.list({ page_size: 100 }).then(r => setTables(r.tables ?? [])).catch(() => {});
  }, []);

  useEffect(() => {
    if (refreshSignal) loadOrders();
  }, [refreshSignal]);

  // ── filtering / pagination ────────────────────────────────────────────────────
  const filtered = useMemo(() => {
    const now = new Date();
    const todayStr = fmtDate(now);
    const { start: weekStart, end: weekEnd } = getMonWeekBounds(now);
    const monthPrefix = todayStr.slice(0, 7); // YYYY-MM

    const matches = reservations.filter(r => {
      if (statusFilter !== 'All Statuses' && r.status !== statusFilter) return false;
      if (searchQuery && !r.name.toLowerCase().includes(searchQuery.toLowerCase())) return false;
      if (viewMode === 'upcoming') {
        const dt = new Date(`${r.date}T${r.time || '00:00'}:00`);
        return r.date === todayStr && dt > now;
      }
      if (viewMode === 'today')  return r.date === todayStr;
      if (viewMode === 'week')   return r.date >= weekStart && r.date <= weekEnd;
      if (viewMode === 'month')  return r.date.startsWith(monthPrefix);
      return true; // 'all'
    });

    const cmp = (a: Reservation, b: Reservation) => {
      const da = new Date(`${a.date}T${a.time || '00:00'}:00`).getTime();
      const db = new Date(`${b.date}T${b.time || '00:00'}:00`).getTime();
      return da - db;
    };
    if (viewMode === 'all') return matches.sort((a, b) => cmp(b, a)); // DESC
    return matches.sort(cmp); // ASC for time-oriented views
  }, [reservations, searchQuery, statusFilter, viewMode]);

  const totalPages = Math.max(1, Math.ceil(filtered.length / PER_PAGE));
  const safePage   = Math.min(currentPage, totalPages);
  const pageSlice  = filtered.slice((safePage - 1) * PER_PAGE, safePage * PER_PAGE);

  useEffect(() => { if (currentPage > totalPages) setCurrentPage(totalPages); }, [currentPage, totalPages]);

  // ── table name lookup ─────────────────────────────────────────────────────────
  const tableLabel = (tableId: string) => {
    if (!tableId) return '—';
    const t = tables.find(t => t.table_id === tableId);
    return t?.table_number != null ? `Table ${t.table_number}` : shortId(tableId);
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
      setError(err instanceof Error ? err.message : 'Failed to update status.');
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
      setError(err instanceof Error ? err.message : 'Failed to save.');
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
      setError(err instanceof Error ? err.message : 'Failed to save.');
    }
  };


  const createWalkIn = async () => {
    if (!walkIn) return;
    setWalkInError('');
    setWalkInBusy(true);
    try {
      const res = await ordersApi.create({
        name:       walkIn.name,
        phone:      walkIn.phone,
        party_size: walkIn.partySize,
        table_id:   walkIn.tableId || undefined,
        date:       walkIn.date,
        time:       walkIn.time,
        end_time:   walkIn.endTime || undefined,
        notes:      walkIn.notes,
        status:     walkIn.status,
        items:      walkIn.items.map(i => ({ item_id: i.id, quantity: i.quantity })),
        walk_in:    true,
      });
      if (res.success === false) { setWalkInError(res.message ?? 'Failed to create order.'); return; }
      setWalkIn(null);
      loadOrders();
    } catch (err) {
      setWalkInError(err instanceof Error ? err.message : 'Failed to create order.');
    } finally {
      setWalkInBusy(false);
    }
  };

  const walkInAddItem = (menuItem: MenuItemDto) => {
    if (!walkIn) return;
    const id = menuItem.item_id ?? menuItem.itemId ?? '';
    const existing = walkIn.items.find(i => i.id === id);
    const items = existing
      ? walkIn.items.map(i => i.id === id ? { ...i, quantity: i.quantity + 1 } : i)
      : [...walkIn.items, { id, name: menuItem.name ?? '', price: menuItem.price ?? 0, quantity: 1, itemStatus: 'PENDING' }];
    setWalkIn({ ...walkIn, items, search: '' });
  };

  const walkInQty = (id: string, delta: number) => {
    if (!walkIn) return;
    setWalkIn({ ...walkIn, items: walkIn.items.map(i => i.id === id ? { ...i, quantity: Math.max(0, i.quantity + delta) } : i).filter(i => i.quantity > 0) });
  };

  // ── render ────────────────────────────────────────────────────────────────────
  return (
    <div className="flex flex-col animate-fadeIn">
      {/* Header */}
      <header className="mb-6 flex md:items-end justify-between gap-6">
        <div>
          <nav className="flex gap-2 text-[10px] text-on-surface-variant uppercase tracking-widest mb-2">
            <span>Admin</span><span>/</span><span className="text-primary font-bold">Orders</span>
          </nav>
          <h2 className="font-serif text-5xl font-bold text-on-surface">Orders Management</h2>
          <p className="text-on-surface-variant text-sm mt-2">View, confirm and manage customer reservations.</p>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => { setWalkIn(defaultWalkIn()); setWalkInError(''); setWalkInCatFilter(''); }}
            className="flex items-center gap-2 px-4 py-2 bg-[#735c00] text-white rounded-lg text-xs font-semibold hover:bg-[#5d4a00] transition-all whitespace-nowrap shadow-sm"
          >
            <span className="material-symbols-outlined text-base">add</span>
            Walk-in Order
          </button>
          <button onClick={loadOrders} className="flex items-center gap-2 px-4 py-2 bg-white border border-outline-variant/30 rounded-lg text-xs font-semibold hover:bg-[#f3f4f5] transition-all whitespace-nowrap">
            <span className="material-symbols-outlined text-base">refresh</span>
            Refresh
          </button>
        </div>
      </header>

      {/* View mode tabs */}
      <div className="flex gap-1 p-1 bg-white rounded-xl shadow-sm mb-4 border border-outline-variant/10 w-fit">
        {VIEW_TABS.map(tab => (
          <button
            key={tab.key}
            onClick={() => { setViewMode(tab.key); setCurrentPage(1); }}
            className={`px-4 py-2 rounded-lg text-xs font-semibold transition-all ${
              viewMode === tab.key
                ? 'bg-[#735c00] text-white shadow-sm'
                : 'text-[#4d4635] hover:bg-[#f3f4f5]'
            }`}
          >
            {tab.label}
            <span className={`ml-1.5 text-[10px] px-1.5 py-0.5 rounded-full ${viewMode === tab.key ? 'bg-white/20 text-white' : 'bg-[#e1e3e4] text-[#4d4635]'}`}>
              {tab.key === 'all' ? reservations.length : (() => {
                const now = new Date();
                const todayStr = fmtDate(now);
                const { start: ws, end: we } = getMonWeekBounds(now);
                const mp = todayStr.slice(0, 7);
                return reservations.filter(r => {
                  if (tab.key === 'upcoming') { const dt = new Date(`${r.date}T${r.time || '00:00'}:00`); return r.date === todayStr && dt > now; }
                  if (tab.key === 'today')  return r.date === todayStr;
                  if (tab.key === 'week')   return r.date >= ws && r.date <= we;
                  if (tab.key === 'month')  return r.date.startsWith(mp);
                  return true;
                }).length;
              })()}
            </span>
          </button>
        ))}
      </div>

      {/* Filter bar */}
      <section className="bg-white p-4 rounded-xl shadow-sm mb-6 flex flex-wrap items-center gap-4 border border-outline-variant/10">
        <div className="flex-grow relative">
          <span className="material-symbols-outlined absolute left-4 top-1/2 -translate-y-1/2 text-on-surface-variant">search</span>
          <input
            type="text"
            value={searchQuery}
            onChange={e => { setSearchQuery(e.target.value); setCurrentPage(1); }}
            className="w-full pl-12 pr-4 py-3 bg-[#f3f4f5] border-none rounded-lg focus:ring-2 focus:ring-[#d4af37] text-sm"
            placeholder="Search by guest name..."
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
      {loading && <p className="mb-4 text-sm text-on-surface-variant">Loading...</p>}

      {/* Table */}
      <div className="bg-white rounded-xl shadow-sm overflow-hidden border border-outline-variant/20">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-[#f3f4f5] border-b border-outline-variant/30">
              <th className="px-5 py-4 text-xs font-semibold text-on-surface-variant uppercase">Guest</th>
              <th className="px-5 py-4 text-xs font-semibold text-on-surface-variant uppercase">Time</th>
              <th className="px-5 py-4 text-xs font-semibold text-on-surface-variant uppercase">Table / Guests</th>
              <th className="px-5 py-4 text-xs font-semibold text-on-surface-variant uppercase">Status</th>
              <th className="px-5 py-4 text-xs font-semibold text-on-surface-variant uppercase text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-outline-variant/10">
            {pageSlice.length === 0 ? (
              <tr><td colSpan={5} className="text-center py-12 text-sm text-on-surface-variant italic">No orders found.</td></tr>
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
                    <span className="text-xs">{res.partySize} guests</span>
                  </div>
                </td>
                <td className="px-5 py-4">
                  <span className={`text-[10px] px-3 py-1 rounded-full font-bold uppercase ${STATUS_STYLE[res.status] ?? 'bg-gray-100 text-gray-700'}`}>
                    {res.status}
                  </span>
                </td>
                <td className="px-5 py-4 text-right whitespace-nowrap">
                  <button onClick={() => setEditModal({ ...res })} className="p-2 text-on-surface-variant hover:text-[#735c00]" title="Edit booking">
                    <span className="material-symbols-outlined text-xl">edit_note</span>
                  </button>
                  <button onClick={() => { setOrderModal({ ...res }); setAddSearch(''); setAddCatFilter(''); }} className="p-2 text-on-surface-variant hover:text-[#735c00]" title="Edit items">
                    <span className="material-symbols-outlined text-xl">restaurant</span>
                  </button>
                  <button onClick={() => setDrawerRes(res)} className="p-2 text-on-surface-variant hover:text-[#735c00]" title="View details">
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
            {filtered.length === 0 ? '0' : `${(safePage - 1) * PER_PAGE + 1}–${Math.min(safePage * PER_PAGE, filtered.length)}`} / {filtered.length} orders
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
                <p className="text-on-surface-variant text-xs mt-0.5">{drawerRes.phone} • {drawerRes.partySize} guests</p>
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
                  <p className="text-on-surface-variant font-semibold uppercase mb-1">Date</p>
                  <p className="font-bold text-on-surface">{drawerRes.date}</p>
                </div>
                <div className="bg-[#f8f9fa] rounded-lg p-3">
                  <p className="text-on-surface-variant font-semibold uppercase mb-1">Time</p>
                  <p className="font-bold text-on-surface">{drawerRes.time}{drawerRes.endTime ? ` – ${drawerRes.endTime}` : ''}</p>
                </div>
                <div className="bg-[#f8f9fa] rounded-lg p-3">
                  <p className="text-on-surface-variant font-semibold uppercase mb-1">Table</p>
                  <p className="font-bold text-on-surface">{tableLabel(drawerRes.tableId)}</p>
                  {drawerRes.tableId && <p className="font-mono text-[10px] text-on-surface-variant mt-0.5">{drawerRes.tableId}</p>}
                </div>
                <div className="bg-[#f8f9fa] rounded-lg p-3">
                  <p className="text-on-surface-variant font-semibold uppercase mb-1">Order ID</p>
                  <p className="font-mono text-[11px] text-on-surface break-all">{drawerRes.id}</p>
                </div>
              </div>

              {/* Notes */}
              {drawerRes.notes && (
                <div className="p-3 bg-amber-50 rounded-lg border border-amber-200">
                  <p className="text-xs font-semibold text-amber-800 uppercase mb-1">📝 Customer Notes</p>
                  <p className="text-sm text-amber-900">{drawerRes.notes}</p>
                </div>
              )}

              {/* Status action buttons */}
              <div>
                <p className="text-xs font-semibold text-on-surface-variant uppercase mb-2">Update Status</p>
                <div className="flex gap-2 flex-wrap">
                  {drawerRes.status === 'Pending' && (
                    <button
                      disabled={!!statusBusy}
                      onClick={() => handleStatusUpdate(drawerRes, 'Confirmed')}
                      className="px-4 py-2 bg-green-600 text-white text-xs font-semibold rounded-lg hover:bg-green-700 disabled:opacity-50 transition-colors"
                    >
                      {statusBusy === drawerRes.id ? '...' : '✓ Confirm'}
                    </button>
                  )}
                  {drawerRes.status === 'Confirmed' && (
                    <button
                      disabled={!!statusBusy}
                      onClick={() => handleStatusUpdate(drawerRes, 'Completed')}
                      className="px-4 py-2 bg-[#735c00] text-white text-xs font-semibold rounded-lg hover:bg-[#5d4a00] disabled:opacity-50 transition-colors"
                    >
                      {statusBusy === drawerRes.id ? '...' : '✓ Complete'}
                    </button>
                  )}
                  {(drawerRes.status === 'Pending' || drawerRes.status === 'Confirmed') && (
                    <button
                      disabled={!!statusBusy}
                      onClick={() => handleStatusUpdate(drawerRes, 'Cancelled')}
                      className="px-4 py-2 bg-red-50 text-red-700 border border-red-200 text-xs font-semibold rounded-lg hover:bg-red-100 disabled:opacity-50 transition-colors"
                    >
                      Cancel Order
                    </button>
                  )}
                </div>
              </div>

              {/* Items list */}
              <div>
                <p className="text-xs font-semibold text-on-surface-variant uppercase mb-3">Items ({drawerRes.items.length})</p>
                {drawerRes.items.length === 0 ? (
                  <p className="text-xs text-on-surface-variant italic">No items yet.</p>
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
                  <span className="font-bold text-on-surface text-sm">Total</span>
                  <span className="text-xl font-serif font-bold text-[#735c00]">{fmtVnd(drawerRes.total)}</span>
                </div>
              </div>
            </div>

            <div className="p-6 border-t bg-white">
              <button className="w-full bg-[#735c00] text-white text-xs font-semibold py-3 rounded-lg hover:bg-[#5d4a00] transition-colors" onClick={() => setDrawerRes(null)}>
                Close
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
              <h3 className="font-serif text-2xl font-bold text-on-surface">Edit Booking</h3>
              <button className="p-1 hover:bg-[#f3f4f5] rounded-full" onClick={() => setEditModal(null)}>
                <span className="material-symbols-outlined">close</span>
              </button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">Guest Name</label>
                <input className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={editModal.name} onChange={e => setEditModal({ ...editModal, name: e.target.value })} />
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">Phone Number</label>
                <input className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={editModal.phone} onChange={e => setEditModal({ ...editModal, phone: e.target.value })} />
              </div>
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">Special Notes</label>
                <textarea className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm resize-none" rows={2} placeholder="Allergies, special requests..." value={editModal.notes} onChange={e => setEditModal({ ...editModal, notes: e.target.value })} />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Date</label>
                  <input type="date" className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={editModal.date} onChange={e => setEditModal({ ...editModal, date: e.target.value })} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Guests</label>
                  <input type="number" min={1} className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={editModal.partySize} onChange={e => setEditModal({ ...editModal, partySize: parseInt(e.target.value) || 1 })} />
                </div>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Start Time</label>
                  <input type="time" className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={editModal.time} onChange={e => setEditModal({ ...editModal, time: e.target.value })} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">End Time</label>
                  <input type="time" className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={editModal.endTime} onChange={e => setEditModal({ ...editModal, endTime: e.target.value })} />
                </div>
              </div>
              <div className="flex gap-4 mt-2">
                <button className="flex-1 border border-outline-variant text-xs py-3 rounded-lg" onClick={() => setEditModal(null)}>Cancel</button>
                <button className="flex-1 bg-[#735c00] text-white text-xs py-3 rounded-lg hover:bg-[#5d4a00]" onClick={saveBooking}>Save Changes</button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* ────────────────────────────────────────────────────────────────────────
          MODAL: Sửa món ăn trong đơn
      ──────────────────────────────────────────────────────────────────────── */}
      {/* ────────────────────────────────────────────────────────────────────────
          MODAL: Walk-in Order
      ──────────────────────────────────────────────────────────────────────── */}
      {walkIn && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
          <div className="bg-white rounded-xl shadow-2xl w-full max-w-lg max-h-[92vh] overflow-y-auto p-8">
            <div className="flex items-center justify-between mb-6">
              <div>
                <h3 className="font-serif text-2xl font-bold text-on-surface">Walk-in Order</h3>
                <p className="text-xs text-on-surface-variant mt-0.5">Customer is present — no advance booking required</p>
              </div>
              <button className="p-1 hover:bg-[#f3f4f5] rounded-full" onClick={() => setWalkIn(null)}>
                <span className="material-symbols-outlined">close</span>
              </button>
            </div>

            <div className="space-y-4">
              {/* Name + Phone */}
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Guest Name *</label>
                  <input className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" placeholder="Nguyen Van A" value={walkIn.name} onChange={e => setWalkIn({ ...walkIn, name: e.target.value })} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Phone *</label>
                  <input className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" placeholder="09xxxxxxxx" value={walkIn.phone} onChange={e => setWalkIn({ ...walkIn, phone: e.target.value })} />
                </div>
              </div>

              {/* Party size + Table + Status */}
              <div className="grid grid-cols-3 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Guests *</label>
                  <input type="number" min={1} className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={walkIn.partySize} onChange={e => setWalkIn({ ...walkIn, partySize: parseInt(e.target.value) || 1 })} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Table</label>
                  <select className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={walkIn.tableId} onChange={e => setWalkIn({ ...walkIn, tableId: e.target.value })}>
                    <option value="">Auto-assign</option>
                    {tables.slice().sort((a, b) => (a.table_number ?? 0) - (b.table_number ?? 0)).map(t => (
                      <option key={t.table_id} value={t.table_id}>Table {t.table_number} (cap {t.capacity})</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Status</label>
                  <select className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={walkIn.status} onChange={e => setWalkIn({ ...walkIn, status: e.target.value })}>
                    {STATUS_OPTIONS.map(s => <option key={s}>{s}</option>)}
                  </select>
                </div>
              </div>

              {/* Date + Time */}
              <div className="grid grid-cols-3 gap-4">
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Date</label>
                  <input type="date" className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={walkIn.date} onChange={e => setWalkIn({ ...walkIn, date: e.target.value })} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">Start Time</label>
                  <input type="time" className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={walkIn.time} onChange={e => setWalkIn({ ...walkIn, time: e.target.value })} />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-on-surface-variant mb-1">End Time</label>
                  <input type="time" className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm" value={walkIn.endTime} onChange={e => setWalkIn({ ...walkIn, endTime: e.target.value })} />
                </div>
              </div>

              {/* Notes */}
              <div>
                <label className="block text-xs font-semibold text-on-surface-variant mb-1">Notes</label>
                <textarea className="w-full bg-[#f8f9fa] border border-outline-variant/30 rounded-lg p-3 text-sm resize-none" rows={2} placeholder="Allergies, special requests..." value={walkIn.notes} onChange={e => setWalkIn({ ...walkIn, notes: e.target.value })} />
              </div>

              {/* Items */}
              <div className="p-4 bg-[#ffe088]/10 rounded-xl border border-[#ffe088]/30">
                <div className="flex items-center justify-between mb-3">
                  <p className="text-xs font-bold text-[#574500] uppercase">Add Items (optional)</p>
                  {walkIn.items.length > 0 && (
                    <span className="text-xs font-bold text-[#735c00]">
                      Subtotal: {fmtVnd(walkIn.items.reduce((s, i) => s + i.price * i.quantity, 0))}
                    </span>
                  )}
                </div>
                <MenuBrowser
                  menuItems={menuItems}
                  categories={categories}
                  currentItems={walkIn.items}
                  search={walkIn.search}
                  catFilter={walkInCatFilter}
                  onSearchChange={s => setWalkIn({ ...walkIn, search: s })}
                  onCatChange={setWalkInCatFilter}
                  onAdd={walkInAddItem}
                  onQtyChange={walkInQty}
                />
              </div>

              {walkInError && <p className="text-xs text-red-600 bg-red-50 p-3 rounded-lg">{walkInError}</p>}

              <div className="flex gap-4 mt-2">
                <button className="flex-1 border border-outline-variant text-xs py-3 rounded-lg" onClick={() => setWalkIn(null)}>Cancel</button>
                <button
                  disabled={walkInBusy || !walkIn.name || !walkIn.phone}
                  className="flex-1 bg-[#735c00] text-white text-xs py-3 rounded-lg hover:bg-[#5d4a00] disabled:opacity-50"
                  onClick={createWalkIn}
                >
                  {walkInBusy ? 'Creating...' : 'Create Walk-in Order'}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {orderModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
          <div className="bg-white rounded-xl shadow-2xl w-full max-w-xl max-h-[90vh] overflow-y-auto p-8">
            <div className="flex items-center justify-between mb-4">
              <div>
                <h3 className="font-serif text-2xl font-bold text-on-surface">Edit Items</h3>
                <p className="text-xs text-on-surface-variant">{orderModal.name} — {orderModal.date}</p>
              </div>
              <button className="p-1 hover:bg-[#f3f4f5] rounded-full" onClick={() => setOrderModal(null)}>
                <span className="material-symbols-outlined">close</span>
              </button>
            </div>

            {/* Current items */}
            <div className="space-y-2 mb-5 max-h-52 overflow-y-auto pr-1">
              {orderModal.items.length === 0 ? (
                <p className="text-center text-xs text-on-surface-variant py-4 italic">No items yet.</p>
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
              <p className="text-xs font-bold text-[#574500] mb-3 uppercase">Add from Menu</p>
              <MenuBrowser
                menuItems={menuItems}
                categories={categories}
                currentItems={orderModal.items}
                search={addSearch}
                catFilter={addCatFilter}
                onSearchChange={setAddSearch}
                onCatChange={setAddCatFilter}
                onAdd={handleAddItem}
                onQtyChange={handleQty}
              />
            </div>

            {/* Subtotal preview */}
            <div className="mt-4 flex justify-between items-center text-sm font-bold text-on-surface">
              <span>Subtotal</span>
              <span className="text-[#735c00]">
                {fmtVnd(orderModal.items.reduce((s, i) => s + i.price * i.quantity, 0))}
              </span>
            </div>

            <div className="flex gap-4 mt-5">
              <button className="flex-1 border border-outline-variant text-xs py-3 rounded-lg" onClick={() => setOrderModal(null)}>Cancel</button>
              <button className="flex-1 bg-[#735c00] text-white text-xs py-3 rounded-lg hover:bg-[#5d4a00]" onClick={saveOrderItems}>Save Changes</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default OrdersManagement;
