import { useState } from 'react';
import TopNavBar from './components/common/TopNavBar';
import HomePage from './pages/HomePage';       // Cụm trang chủ mới tạo ở Bước 1
import ContactPage from './pages/ContactPage';   // Trang liên hệ đã tạo ở bài trước
import ReservationPage from './pages/ReservationPage'; // Trang đặt bàn mới tạo
import Footer from './components/layout/Footer';
import MenuPage from './pages/MenuPage';

// type Page = 'home' | 'contact' | 'menu' | 'about';
type Page = 'home' | 'contact' | 'reservation' | 'menu';

function App() {
  // Quản lý trang hiện tại ('home' hoặc 'contact')
  const [currentPage, setCurrentPage] = useState<'home' | 'contact' | 'reservation' | 'menu'>('home');

  const renderPage = () => {
    const pageMap: Record<Page, React.ReactNode> = {
      home: <HomePage />,
      contact: <ContactPage />,
      reservation: <ReservationPage />,
      menu: <MenuPage />,
      // about: <About />,
    };
    return pageMap[currentPage] || <HomePage />;
  };

  return (
    <div className="min-h-screen bg-background text-on-surface flex flex-col">
      {/* Thanh điều hướng nhận State để điều khiển menu */}
      <TopNavBar currentPage={currentPage} setCurrentPage={setCurrentPage} />

      {/* Render có điều kiện dựa trên State */}
      {renderPage()}

      {/* Chân trang dùng chung */}
      <Footer />
    </div>
  );
}

export default App;