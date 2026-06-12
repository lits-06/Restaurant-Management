import React, { useState, useEffect, useRef } from 'react';
import Sidebar from './components/Sidebar';
import Footer from './components/Footer';
import AnalyticsOverview from './pages/AnalyticsOverview';
import MenuManagement from './pages/MenuManagement';
import OrdersManagement from './pages/OrdersManagement';
import MonthlyScheduler from './pages/MonthlyScheduler';
import TableManagement from './pages/TableManagement';
import UserManagement from './pages/UserManagement';
import LoginPage from './pages/Login';
import { useAdminAuthStore } from './store/adminAuthStore';
import { authApi } from './services/api';
import { useAdminNotifications, type AdminNotification } from './hooks/useAdminNotifications';

const getDefaultTab = (roles: string[]) => {
  if (roles.includes('ADMIN') || roles.includes('MANAGER')) return 'Dashboard';
  return 'Orders';
};

const isTabVisibleForRoles = (tab: string, roles: string[]) => {
  const isAdmin = roles.includes('ADMIN');
  const isAdminOrManager = isAdmin || roles.includes('MANAGER');
  const isStaff = isAdminOrManager || roles.includes('CHEF') || roles.includes('WAITER');
  const visible: Record<string, boolean> = {
    Dashboard: isAdminOrManager,
    Menu: isAdminOrManager,
    Tables: isAdminOrManager,
    Orders: isStaff,
    Staff: isStaff,
    Users: isAdmin,
  };
  return visible[tab] ?? false;
};

const NOTIF_TYPE_LABEL: Record<string, string> = {
  ORDER_CREATED: 'New Order',
  ORDER_STATUS_CHANGED: 'Status Changed',
};

const NOTIF_TYPE_COLOR: Record<string, string> = {
  ORDER_CREATED: 'border-amber-400',
  ORDER_STATUS_CHANGED: 'border-blue-400',
};

const App: React.FC = () => {
  const { user, clearAuth, refreshToken, accessToken } = useAdminAuthStore();

  const [activeTab, setActiveTab] = useState<string>(() => {
    const saved = sessionStorage.getItem('luxe-admin-tab');
    if (saved && isTabVisibleForRoles(saved, user?.roles ?? [])) return saved;
    return getDefaultTab(user?.roles ?? []);
  });

  const [showNotifPanel, setShowNotifPanel] = useState(false);
  const [ordersRefreshSignal, setOrdersRefreshSignal] = useState(0);
  const panelRef = useRef<HTMLDivElement>(null);
  const bellRef = useRef<HTMLButtonElement>(null);

  const { notifications, unreadCount, connected, clearNotifications, markAllRead } =
    useAdminNotifications(user ? accessToken : null, user?.roles ?? []);

  // Auto-refresh Orders tab when relevant notifications arrive
  const prevNotifCount = useRef(0);
  useEffect(() => {
    if (notifications.length > prevNotifCount.current) {
      setOrdersRefreshSignal((s) => s + 1);
    }
    prevNotifCount.current = notifications.length;
  }, [notifications.length]);

  // Close panel when clicking outside (but not on the bell button itself)
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      const target = e.target as Node;
      const inPanel = panelRef.current?.contains(target);
      const inBell = bellRef.current?.contains(target);
      if (!inPanel && !inBell) setShowNotifPanel(false);
    };
    if (showNotifPanel) document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, [showNotifPanel]);

  useEffect(() => {
    sessionStorage.setItem('luxe-admin-tab', activeTab);
  }, [activeTab]);

  const handleLogout = async () => {
    if (refreshToken) {
      try {
        await authApi.logout(refreshToken);
      } catch {
        // ignore
      }
    }
    sessionStorage.removeItem('luxe-admin-tab');
    clearAuth();
  };

  const handleBellClick = () => {
    setShowNotifPanel((v) => !v);
    markAllRead();
  };

  if (!user) {
    return <LoginPage onSuccess={() => {
      const u = useAdminAuthStore.getState().user;
      if (u) {
        const tab = getDefaultTab(u.roles);
        sessionStorage.setItem('luxe-admin-tab', tab);
        setActiveTab(tab);
      }
    }} />;
  }

  return (
    <div className="bg-[#f8f9fa] text-[#191c1d] font-sans overflow-x-hidden min-h-screen">
      <Sidebar
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        user={user}
        onLogout={handleLogout}
        unreadCount={unreadCount}
        onBellClick={user.roles.some(r => r === 'ADMIN' || r === 'MANAGER') ? handleBellClick : undefined}
        connected={connected}
        bellRef={bellRef}
      />

      {/* Notification panel */}
      {showNotifPanel && (
        <div
          ref={panelRef}
          className="fixed top-0 left-64 h-screen w-80 bg-white shadow-2xl z-40 flex flex-col border-r border-gray-200"
        >
          <div className="flex items-center justify-between px-4 py-4 border-b border-gray-100">
            <span className="font-semibold text-sm text-[#191c1d]">Notifications</span>
            <div className="flex items-center gap-2">
              {notifications.length > 0 && (
                <button
                  onClick={clearNotifications}
                  className="text-xs text-gray-400 hover:text-gray-600"
                >
                  Clear all
                </button>
              )}
              <button onClick={() => setShowNotifPanel(false)} className="text-gray-400 hover:text-gray-600">
                <span className="material-symbols-outlined text-base">close</span>
              </button>
            </div>
          </div>
          <div className="flex-1 overflow-y-auto">
            {notifications.length === 0 ? (
              <p className="text-center text-xs text-gray-400 mt-10">No notifications</p>
            ) : (
              notifications.map((n, i) => (
                <NotifItem key={n.id || i} notif={n} onOrderClick={() => {
                  setActiveTab('Orders');
                  setShowNotifPanel(false);
                }} />
              ))
            )}
          </div>
        </div>
      )}

      <main className="ml-64 min-h-screen p-10 flex flex-col justify-between">
        <div className="flex-grow">
          <div className={activeTab === 'Dashboard' ? 'block' : 'hidden'}>
            <AnalyticsOverview />
          </div>
          <div className={activeTab === 'Menu' ? 'block' : 'hidden'}>
            <MenuManagement />
          </div>
          <div className={activeTab === 'Orders' ? 'block' : 'hidden'}>
            <OrdersManagement refreshSignal={ordersRefreshSignal} />
          </div>
          <div className={activeTab === 'Staff' ? 'block' : 'hidden'}>
            <MonthlyScheduler />
          </div>
          <div className={activeTab === 'Tables' ? 'block' : 'hidden'}>
            <TableManagement />
          </div>
          <div className={activeTab === 'Users' ? 'block' : 'hidden'}>
            <UserManagement />
          </div>
        </div>
        <Footer />
      </main>
    </div>
  );
};

interface NotifItemProps {
  notif: AdminNotification;
  onOrderClick: () => void;
}

const NotifItem: React.FC<NotifItemProps> = ({ notif, onOrderClick }) => {
  const color = NOTIF_TYPE_COLOR[notif.type] ?? 'border-gray-300';
  const label = NOTIF_TYPE_LABEL[notif.type] ?? notif.type;
  const ts = notif.created_at
    ? new Date(notif.created_at * 1000).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })
    : '';
  return (
    <div
      className={`border-l-4 ${color} px-4 py-3 border-b border-gray-50 cursor-pointer hover:bg-gray-50`}
      onClick={onOrderClick}
    >
      <div className="flex items-center justify-between mb-0.5">
        <span className="text-xs font-semibold text-gray-700">{label}</span>
        {ts && <span className="text-[10px] text-gray-400">{ts}</span>}
      </div>
      <p className="text-xs text-gray-600 leading-snug">{notif.message}</p>
      {notif.order_id && (
        <p className="text-[10px] text-gray-400 mt-0.5">Order #{notif.order_id.slice(0, 8)}</p>
      )}
    </div>
  );
};

export default App;
