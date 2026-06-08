import React, { useState } from 'react';
import Sidebar from './components/Sidebar';
import Footer from './components/Footer';

// Import trang AnalyticsOverview nguyên vẹn của bạn
import AnalyticsOverview from './pages/AnalyticsOverview';
import MenuManagement from './pages/MenuManagement';
import OrdersManagement from './pages/OrdersManagement';
import StaffManagement from './pages/StaffManagement';
import WeeklyScheduler from './pages/WeeklyScheduler';

const App: React.FC = () => {
  // Sidebar chỉ quản lý việc chuyển đổi tab ở đây
  const [activeTab, setActiveTab] = useState<string>('Dashboard');

  const [staffView, setStaffView] = useState<'management' | 'scheduler'>('management');

  const handleSetActiveTab = (tab: string) => {
    if (tab !== 'Staff') setStaffView('management');
    setActiveTab(tab);
  };

  return (
    <div className="bg-[#f8f9fa] text-[#191c1d] font-sans overflow-x-hidden min-h-screen">
      {/* Sidebar quản lý điều hướng */}
      <Sidebar activeTab={activeTab} setActiveTab={setActiveTab} />

      {/* Khu vực hiển thị nội dung các trang */}
      <main className="ml-64 min-h-screen p-10 flex flex-col justify-between">
        <div className="flex-grow">
          
          {/* TRANG ANALYTICS OVERVIEW GIỮ NGUYÊN VẸN, TỰ QUẢN LÝ STATE */}
          <div className={activeTab === 'Dashboard' ? 'block' : 'hidden'}>
            <AnalyticsOverview />
          </div>
          <div className={activeTab === 'Menu' ? 'block' : 'hidden'}>
            <MenuManagement />
          </div>
          <div className={activeTab === 'Orders' ? 'block' : 'hidden'}>
            <OrdersManagement />
          </div>
          {/* Staff section: 2 sub-views, dùng visibility để giữ state */}
          {activeTab === 'Staff' && (
            <>
              <div className={staffView === 'management' ? 'block' : 'hidden'}>
                <StaffManagement
                  onNavigateToScheduler={() => setStaffView('scheduler')}
                />
              </div>
              <div className={staffView === 'scheduler' ? 'block' : 'hidden'}>
                <WeeklyScheduler
                  onBack={() => setStaffView('management')}
                />
              </div>
            </>
          )}

        </div>

        {/* Footer chung nằm dưới cùng */}
        <Footer />
      </main>
    </div>
  );
};

export default App;