import React, { useState } from 'react';
import Sidebar from './components/Sidebar';
import Footer from './components/Footer';
import AnalyticsOverview from './pages/AnalyticsOverview';
import MenuManagement from './pages/MenuManagement';
import OrdersManagement from './pages/OrdersManagement';
import MonthlyScheduler from './pages/MonthlyScheduler';
import LoginPage from './pages/Login';
import { useAdminAuthStore } from './store/adminAuthStore';
import { authApi } from './services/api';

const App: React.FC = () => {
  const { user, clearAuth, refreshToken } = useAdminAuthStore();
  const [activeTab, setActiveTab] = useState<string>('Dashboard');

  const handleLogout = async () => {
    if (refreshToken) {
      try {
        await authApi.logout(refreshToken);
      } catch {
        // ignore
      }
    }
    clearAuth();
  };

  if (!user) {
    return <LoginPage onSuccess={() => {}} />;
  }

  return (
    <div className="bg-[#f8f9fa] text-[#191c1d] font-sans overflow-x-hidden min-h-screen">
      <Sidebar
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        user={user}
        onLogout={handleLogout}
      />

      <main className="ml-64 min-h-screen p-10 flex flex-col justify-between">
        <div className="flex-grow">
          <div className={activeTab === 'Dashboard' ? 'block' : 'hidden'}>
            <AnalyticsOverview />
          </div>
          <div className={activeTab === 'Menu' ? 'block' : 'hidden'}>
            <MenuManagement />
          </div>
          <div className={activeTab === 'Orders' ? 'block' : 'hidden'}>
            <OrdersManagement />
          </div>
          <div className={activeTab === 'Staff' ? 'block' : 'hidden'}>
            <MonthlyScheduler />
          </div>
        </div>
        <Footer />
      </main>
    </div>
  );
};

export default App;
