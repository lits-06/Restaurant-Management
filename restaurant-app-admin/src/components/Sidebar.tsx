import React from 'react';
import { AdminUser } from '../store/adminAuthStore';

interface SidebarProps {
  activeTab: string;
  setActiveTab: (tab: string) => void;
  user: AdminUser;
  onLogout: () => void;
  unreadCount?: number;
  onBellClick?: () => void;
  connected?: boolean;
  bellRef?: React.RefObject<HTMLButtonElement>;
}

const Sidebar: React.FC<SidebarProps> = ({ activeTab, setActiveTab, user, onLogout, unreadCount = 0, onBellClick, connected, bellRef }) => {
  const roles = user.roles;
  const isAdmin = roles.includes('ADMIN');
  const isAdminOrManager = isAdmin || roles.includes('MANAGER');
  const isStaff = isAdminOrManager || roles.includes('CHEF') || roles.includes('WAITER');

  const allMenuItems = [
    { icon: 'dashboard',        label: 'Dashboard',  visible: isAdminOrManager },
    { icon: 'restaurant_menu',  label: 'Menu',       visible: isAdminOrManager },
    { icon: 'receipt_long',     label: 'Orders',     visible: isStaff },
    { icon: 'calendar_month',   label: 'Staff',      visible: isStaff },
    { icon: 'table_restaurant', label: 'Tables',     visible: isAdminOrManager },
    { icon: 'group',            label: 'Users',      visible: isAdmin },
  ];
  const menuItems = allMenuItems.filter(item => item.visible);

  const initials = (user.full_name || user.username || 'A')
    .split(' ')
    .map((w) => w[0])
    .join('')
    .slice(0, 2)
    .toUpperCase();

  const displayRole = user.roles.includes('ADMIN')
    ? 'Admin'
    : user.roles.includes('MANAGER')
    ? 'Manager'
    : user.roles.includes('CHEF')
    ? 'Chef'
    : user.roles.includes('WAITER')
    ? 'Waiter'
    : user.roles[0] ?? 'Staff';

  return (
    <aside className="h-screen w-64 fixed left-0 top-0 bg-[#edeeef] flex flex-col py-2 shadow-md z-50">
      <div className="px-6 py-8 flex items-start justify-between">
        <div>
          <h1 className="font-serif text-2xl text-[#735c00] font-semibold">LuxeBistro Admin</h1>
          <p className="text-xs font-semibold tracking-wider text-[#4d4635] opacity-70 mt-1 uppercase">
            Service Mode: Dinner
          </p>
        </div>
        {onBellClick && (
          <button
            ref={bellRef}
            onClick={onBellClick}
            title={connected ? 'Notifications (connected)' : 'Notifications (disconnected)'}
            className="relative mt-1 p-1.5 rounded-lg hover:bg-[#e1e3e4] transition-all text-[#4d4635]"
          >
            <span className="material-symbols-outlined text-xl">notifications</span>
            {unreadCount > 0 && (
              <span className="absolute -top-0.5 -right-0.5 min-w-[16px] h-4 px-0.5 bg-red-500 text-white text-[9px] font-bold rounded-full flex items-center justify-center leading-none">
                {unreadCount > 99 ? '99+' : unreadCount}
              </span>
            )}
            {!connected && (
              <span className="absolute bottom-0 right-0 w-2 h-2 bg-gray-400 rounded-full border border-[#edeeef]" />
            )}
          </button>
        )}
      </div>

      <nav className="flex-grow">
        <ul className="space-y-1">
          {menuItems.map((item, index) => (
            <li
              key={index}
              onClick={() => setActiveTab(item.label)}
              className={`mx-2 my-1 px-4 py-3 rounded-lg transition-all cursor-pointer flex items-center gap-3 ${
                activeTab === item.label
                  ? 'bg-[#d4af37] text-[#554300] shadow-sm'
                  : 'text-[#4d4635] hover:bg-[#e1e3e4]'
              }`}
            >
              <span className="material-symbols-outlined">{item.icon}</span>
              <span className="text-xs font-semibold tracking-wider">{item.label}</span>
            </li>
          ))}
        </ul>
      </nav>

      <div className="mt-auto px-4 py-4 border-t border-[#d0c5af] space-y-2">
        <div className="flex items-center gap-3 p-2 bg-[#f3f4f5] rounded-lg">
          <div className="w-10 h-10 rounded-full bg-[#735c00] flex items-center justify-center text-white font-bold text-sm flex-shrink-0">
            {initials}
          </div>
          <div className="min-w-0">
            <p className="text-xs font-semibold text-[#191c1d] truncate">{user.full_name || user.username}</p>
            <p className="text-[10px] text-[#4d4635] uppercase tracking-wider">{displayRole}</p>
          </div>
        </div>
        <button
          onClick={onLogout}
          className="w-full flex items-center gap-2 px-3 py-2 rounded-lg text-[#4d4635] hover:bg-[#e1e3e4] transition-all text-xs font-semibold"
        >
          <span className="material-symbols-outlined text-base">logout</span>
          Sign Out
        </button>
      </div>
    </aside>
  );
};

export default Sidebar;
