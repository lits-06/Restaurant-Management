import React, { useEffect, useMemo, useState } from 'react';
import Header from '../components/HeaderDashboard';
import KPIGrid, { type KPIItem } from '../components/KPIGrid';
import PerformanceTable from '../components/PerformanceTable';
import { menuApi, ordersApi, staffApi, type MenuItemDto, type OrderDto } from '../services/api';

const currency = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' });

const getMenuItemId = (item: MenuItemDto) => item.item_id ?? item.itemId ?? '';
const getOrderItems = (orders: OrderDto[]) => orders.flatMap((order) => order.items ?? []);
const getOrderTotal = (order: OrderDto) => {
  const explicitTotal = order.total_price ?? order.totalPrice;
  if (explicitTotal) return explicitTotal;

  return (order.items ?? []).reduce((sum, item) => sum + (item.price ?? 0) * (item.quantity ?? 0), 0);
};

const AnalyticsOverview: React.FC = () => {
  // Trạng thái chọn tháng được cô lập hoàn toàn bên trong trang này
  const [selectedMonth, setSelectedMonth] = useState<string>('October 2023');
  const [orders, setOrders] = useState<OrderDto[]>([]);
  const [menuItems, setMenuItems] = useState<MenuItemDto[]>([]);
  const [staffTotal, setStaffTotal] = useState<number>(0);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    let isMounted = true;

    Promise.all([
      ordersApi.list({ page: 1, page_size: 100 }),
      menuApi.listItems({ page: 1, page_size: 100 }),
      staffApi.list({ page: 1, page_size: 100 }),
    ])
      .then(([ordersResponse, menuResponse, staffResponse]) => {
        if (!isMounted) return;
        setOrders(ordersResponse.orders ?? []);
        setMenuItems(menuResponse.items ?? []);
        setStaffTotal(staffResponse.total ?? staffResponse.staff?.length ?? 0);
      })
      .catch((err) => {
        if (isMounted) {
          setError(err instanceof Error ? err.message : 'Không thể tải dữ liệu analytics từ gateway.');
        }
      });

    return () => {
      isMounted = false;
    };
  }, []);

  const totalRevenue = useMemo(() => orders.reduce((sum, order) => sum + getOrderTotal(order), 0), [orders]);
  const totalOrders = orders.length;
  const averageOrderValue = totalOrders > 0 ? totalRevenue / totalOrders : 0;

  const kpis: KPIItem[] = [
    {
      title: 'Total Revenue',
      value: currency.format(totalRevenue),
      trend: 'From order-service',
      isUp: totalRevenue > 0,
      isPrimary: true,
    },
    {
      title: 'Total Orders',
      value: totalOrders.toLocaleString('en-US'),
      trend: 'Live order records',
      isUp: totalOrders > 0,
      isPrimary: false,
    },
    {
      title: 'Avg. Order Value',
      value: currency.format(averageOrderValue),
      trend: 'Calculated in frontend',
      isUp: null,
      isPrimary: false,
    },
    {
      title: 'Staff Members',
      value: staffTotal.toLocaleString('en-US'),
      trend: 'From staff-service',
      isUp: staffTotal > 0,
      isPrimary: false,
    },
  ];

  const performanceRows = useMemo(() => {
    const menuById = new Map(menuItems.map((item) => [getMenuItemId(item), item]));
    const totals = new Map<string, { name: string; unitsSold: number; revenue: number }>();

    getOrderItems(orders).forEach((item) => {
      const itemId = item.item_id ?? item.itemId ?? '';
      const menuItem = menuById.get(itemId);
      const name = item.name ?? menuItem?.name ?? (itemId || 'Unknown item');
      const price = item.price ?? menuItem?.price ?? 0;
      const quantity = item.quantity ?? 0;
      const current = totals.get(itemId || name) ?? { name, unitsSold: 0, revenue: 0 };

      current.unitsSold += quantity;
      current.revenue += price * quantity;
      totals.set(itemId || name, current);
    });

    return Array.from(totals.values())
      .sort((a, b) => b.revenue - a.revenue)
      .slice(0, 5)
      .map((item) => ({
        name: item.name,
        unitsSold: item.unitsSold,
        revenue: currency.format(item.revenue),
      }));
  }, [menuItems, orders]);

  return (
    <div className="flex flex-col gap-12">
      {/* Khối tiêu đề & Lịch chọn tháng của trang */}
      <Header 
        selectedMonth={selectedMonth} 
        setSelectedMonth={setSelectedMonth} 
      />

      {error && <p className="text-sm text-red-600">{error}</p>}

      {/* Khối các chỉ số doanh thu */}
      <KPIGrid kpis={kpis} />

      {/* Bảng danh sách món ăn */}
      <div className="grid grid-cols-12 gap-6">
        <PerformanceTable dishes={performanceRows} />
      </div>
    </div>
  );
};

export default AnalyticsOverview;
