import React from 'react';
import { AdminUser } from '../store/adminAuthStore';

interface SidebarProps {
  activeTab: string;
  setActiveTab: (tab: string) => void;
  user: AdminUser;
  onLogout: () => void;
}

const Sidebar: React.FC<SidebarProps> = ({ activeTab, setActiveTab, user, onLogout }) => {
  const menuItems = [
    { icon: 'dashboard',       label: 'Dashboard' },
    { icon: 'restaurant_menu', label: 'Menu' },
    { icon: 'receipt_long',    label: 'Orders' },
    { icon: 'calendar_month',  label: 'Staff' },
    { icon: 'table_restaurant',label: 'Tables' },
    { icon: 'group',           label: 'Users' },
  ];

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
      <div className="px-6 py-8">
        <h1 className="font-serif text-2xl text-[#735c00] font-semibold">LuxeBistro Admin</h1>
        <p className="text-xs font-semibold tracking-wider text-[#4d4635] opacity-70 mt-1 uppercase">
          Service Mode: Dinner
        </p>
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
