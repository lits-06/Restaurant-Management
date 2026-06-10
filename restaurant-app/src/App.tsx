import { useState } from 'react';
import TopNavBar from './components/common/TopNavBar';
import HomePage from './pages/HomePage';
import ContactPage from './pages/ContactPage';
import ReservationPage from './pages/ReservationPage';
import Footer from './components/layout/Footer';
import MenuPage from './pages/MenuPage';
import LoginPage from './pages/LoginPage';
import MyOrdersPage from './pages/MyOrdersPage';
import { useAuthStore } from './store/authStore';
import { authApi } from './api/gateway.api';

type Page = 'home' | 'contact' | 'reservation' | 'menu' | 'login' | 'my-orders';

const PROTECTED: Page[] = ['reservation', 'my-orders'];

function App() {
  const { user, clearAuth, refreshToken } = useAuthStore();
  const [currentPage, setCurrentPage] = useState<Page>('home');
  // Stores where to redirect after successful login
  const [loginRedirect, setLoginRedirect] = useState<Page>('home');

  const navigateTo = (page: Page) => {
    if (PROTECTED.includes(page) && !user) {
      setLoginRedirect(page);
      setCurrentPage('login');
      return;
    }
    setCurrentPage(page);
  };

  const handleLoginSuccess = () => {
    setCurrentPage(loginRedirect);
    setLoginRedirect('home');
  };

  const handleLogout = async () => {
    if (refreshToken) {
      try {
        await authApi.logout(refreshToken);
      } catch {
        // ignore network errors on logout
      }
    }
    clearAuth();
    setCurrentPage('home');
  };

  const renderPage = () => {
    switch (currentPage) {
      case 'home':
        return <HomePage />;
      case 'menu':
        return <MenuPage />;
      case 'reservation':
        return <ReservationPage onNeedLogin={() => { setLoginRedirect('reservation'); setCurrentPage('login'); }} />;
      case 'contact':
        return <ContactPage />;
      case 'login':
        return <LoginPage onSuccess={handleLoginSuccess} />;
      case 'my-orders':
        return <MyOrdersPage />;
      default:
        return <HomePage />;
    }
  };

  return (
    <div className="min-h-screen bg-background text-on-surface flex flex-col">
      <TopNavBar
        currentPage={currentPage}
        navigateTo={navigateTo}
        onLogout={handleLogout}
      />
      {renderPage()}
      <Footer />
    </div>
  );
}

export default App;
