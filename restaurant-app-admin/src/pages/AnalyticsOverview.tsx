import React, { useEffect, useMemo, useState } from 'react';
import Header from '../components/HeaderDashboard';
import KPIGrid, { type KPIItem } from '../components/KPIGrid';
import PerformanceTable from '../components/PerformanceTable';
import { menuApi, ordersApi, type MenuItemDto, type OrderDto } from '../services/api';

// ── helpers ────────────────────────────────────────────────────────────────────

const fmtVnd = (n: number) => `${Math.round(n).toLocaleString('vi-VN')}đ`;

const getMenuItemId = (item: MenuItemDto) => item.item_id ?? item.itemId ?? '';

const getOrderDate = (order: OrderDto): Date | null => {
  const v = order.time;
  if (!v) return null;
  if (typeof v === 'string') return new Date(v);
  if (typeof v === 'object' && v.seconds) return new Date(v.seconds * 1000);
  return null;
};

const getOrderTotal = (order: OrderDto) => {
  const explicit = order.total ?? order.total_price ?? order.totalPrice;
  if (explicit && explicit > 0) return explicit;
  return (order.items ?? []).reduce((s, i) => s + (i.price ?? 0) * (i.quantity ?? 0), 0);
};

const trendLabel = (curr: number, prev: number): { label: string; isUp: boolean | null } => {
  if (prev === 0) return { label: 'No data for previous month', isUp: null };
  const pct = ((curr - prev) / prev) * 100;
  const sign = pct >= 0 ? '+' : '';
  return { label: `${sign}${pct.toFixed(1)}% vs previous month`, isUp: pct >= 0 };
};

// ── component ──────────────────────────────────────────────────────────────────

const AnalyticsOverview: React.FC = () => {
  const now = new Date();
  const [year, setYear]   = useState(now.getFullYear());
  const [month, setMonth] = useState(now.getMonth()); // 0-indexed

  const [allOrders, setAllOrders]   = useState<OrderDto[]>([]);
  const [menuItems, setMenuItems]   = useState<MenuItemDto[]>([]);
  const [loading, setLoading]       = useState(true);
  const [error, setError]           = useState('');

  useEffect(() => {
    let active = true;
    setLoading(true);
    Promise.all([
      ordersApi.list({ page: 1, page_size: 500 }),
      menuApi.listItems({ page: 1, page_size: 100 }),
    ])
      .then(([ordersRes, menuRes]) => {
        if (!active) return;
        setAllOrders(ordersRes.orders ?? []);
        setMenuItems(menuRes.items ?? []);
      })
      .catch(err => { if (active) setError(err instanceof Error ? err.message : 'Failed to load data.'); })
      .finally(() => { if (active) setLoading(false); });
    return () => { active = false; };
  }, []);

  // ── filter by selected month ────────────────────────────────────────────────

  const prevMonth = month === 0 ? 11 : month - 1;
  const prevYear  = month === 0 ? year - 1 : year;

  const monthOrders = useMemo(() =>
    allOrders.filter(o => {
      const d = getOrderDate(o);
      return d && d.getFullYear() === year && d.getMonth() === month;
    }), [allOrders, year, month]);

  const prevMonthOrders = useMemo(() =>
    allOrders.filter(o => {
      const d = getOrderDate(o);
      return d && d.getFullYear() === prevYear && d.getMonth() === prevMonth;
    }), [allOrders, prevYear, prevMonth]);

  // ── KPI calculations ────────────────────────────────────────────────────────

  const revenue     = useMemo(() => monthOrders.reduce((s, o) => s + getOrderTotal(o), 0), [monthOrders]);
  const prevRevenue = useMemo(() => prevMonthOrders.reduce((s, o) => s + getOrderTotal(o), 0), [prevMonthOrders]);

  const covers     = useMemo(() => monthOrders.reduce((s, o) => s + (o.party_size ?? o.partySize ?? 0), 0), [monthOrders]);
  const prevCovers = useMemo(() => prevMonthOrders.reduce((s, o) => s + (o.party_size ?? o.partySize ?? 0), 0), [prevMonthOrders]);

  const avgOrder     = monthOrders.length > 0 ? revenue / monthOrders.length : 0;
  const prevAvgOrder = prevMonthOrders.length > 0 ? prevRevenue / prevMonthOrders.length : 0;

  const completedCount = monthOrders.filter(o => (o.status ?? '').toLowerCase() === 'completed').length;
  const cancelledCount = monthOrders.filter(o => (o.status ?? '').toLowerCase() === 'cancelled').length;

  const revTrend    = trendLabel(revenue, prevRevenue);
  const orderTrend  = trendLabel(monthOrders.length, prevMonthOrders.length);
  const coversTrend = trendLabel(covers, prevCovers);
  const avgTrend    = trendLabel(avgOrder, prevAvgOrder);

  const kpis: KPIItem[] = [
    {
      title: 'Monthly Revenue',
      value: fmtVnd(revenue),
      trend: revTrend.label,
      isUp: revTrend.isUp,
      isPrimary: true,
    },
    {
      title: 'Total Reservations',
      value: monthOrders.length.toLocaleString('vi-VN'),
      trend: orderTrend.label,
      isUp: orderTrend.isUp,
      isPrimary: false,
    },
    {
      title: 'Avg. Value / Order',
      value: fmtVnd(avgOrder),
      trend: avgTrend.label,
      isUp: avgTrend.isUp,
      isPrimary: false,
    },
    {
      title: 'Total Covers',
      value: covers.toLocaleString('vi-VN'),
      trend: coversTrend.label,
      isUp: coversTrend.isUp,
      isPrimary: false,
    },
  ];

  // ── order status breakdown ───────────────────────────────────────────────────

  const pendingCount   = monthOrders.filter(o => (o.status ?? '').toLowerCase() === 'pending').length;
  const confirmedCount = monthOrders.filter(o => (o.status ?? '').toLowerCase() === 'confirmed').length;

  // ── top dishes in selected month ─────────────────────────────────────────────

  const performanceRows = useMemo(() => {
    const menuById = new Map(menuItems.map(i => [getMenuItemId(i), i]));
    const totals   = new Map<string, { name: string; unitsSold: number; revenue: number }>();

    monthOrders.flatMap(o => o.items ?? []).forEach(item => {
      const itemId   = item.item_id ?? item.itemId ?? '';
      const menuItem = menuById.get(itemId);
      const name     = item.name ?? menuItem?.name ?? (itemId || 'Unknown');
      const price    = item.price ?? menuItem?.price ?? 0;
      const qty      = item.quantity ?? 0;
      const key      = itemId || name;
      const current  = totals.get(key) ?? { name, unitsSold: 0, revenue: 0 };
      current.unitsSold += qty;
      current.revenue   += price * qty;
      totals.set(key, current);
    });

    return Array.from(totals.values())
      .sort((a, b) => b.revenue - a.revenue)
      .slice(0, 5)
      .map(i => ({ name: i.name, unitsSold: i.unitsSold, revenue: fmtVnd(i.revenue) }));
  }, [menuItems, monthOrders]);

  // ── render ──────────────────────────────────────────────────────────────────

  return (
    <div className="flex flex-col gap-10">
      <Header year={year} month={month} onChange={(y, m) => { setYear(y); setMonth(m); }} />

      {error && <p className="text-sm text-red-600">{error}</p>}
      {loading && <p className="text-sm text-on-surface-variant">Loading data...</p>}

      <KPIGrid kpis={kpis} />

      {/* Order status breakdown */}
      {monthOrders.length > 0 && (
        <div className="grid grid-cols-4 gap-4">
          {[
            { label: 'Pending',   count: pendingCount,   color: 'bg-amber-50 border-amber-200 text-amber-800' },
            { label: 'Confirmed', count: confirmedCount, color: 'bg-blue-50 border-blue-200 text-blue-800' },
            { label: 'Completed', count: completedCount, color: 'bg-green-50 border-green-200 text-green-800' },
            { label: 'Cancelled', count: cancelledCount, color: 'bg-red-50 border-red-200 text-red-800' },
          ].map(({ label, count, color }) => (
            <div key={label} className={`p-4 rounded-xl border ${color} flex justify-between items-center`}>
              <span className="text-xs font-semibold">{label}</span>
              <span className="font-serif text-2xl font-bold">{count}</span>
            </div>
          ))}
        </div>
      )}

      {/* Top dishes */}
      <div className="grid grid-cols-12 gap-6">
        <PerformanceTable dishes={performanceRows} />
      </div>
    </div>
  );
};

export default AnalyticsOverview;
