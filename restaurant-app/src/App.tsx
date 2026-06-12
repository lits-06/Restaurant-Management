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

const SESSION_KEY = 'luxe-customer-page';

function App() {
  const { user, clearAuth, refreshToken } = useAuthStore();
  const [currentPage, setCurrentPage] = useState<Page>(() => {
    const saved = sessionStorage.getItem(SESSION_KEY) as Page | null;
    if (saved && saved !== 'login') {
      if (!PROTECTED.includes(saved) || !!useAuthStore.getState().user) return saved;
    }
    return 'home';
  });
  const [loginRedirect, setLoginRedirect] = useState<Page>('home');

  const navigateTo = (page: Page) => {
    if (PROTECTED.includes(page) && !user) {
      setLoginRedirect(page);
      sessionStorage.setItem(SESSION_KEY, 'login');
      setCurrentPage('login');
      return;
    }
    sessionStorage.setItem(SESSION_KEY, page);
    setCurrentPage(page);
  };

  const handleLoginSuccess = () => {
    const target = loginRedirect;
    sessionStorage.setItem(SESSION_KEY, target);
    setCurrentPage(target);
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
    sessionStorage.setItem(SESSION_KEY, 'home');
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

  const isLoginPage = currentPage === 'login';

  return (
    <div className="min-h-screen bg-background text-on-surface flex flex-col">
      {!isLoginPage && (
        <TopNavBar
          currentPage={currentPage}
          navigateTo={navigateTo}
          onLogout={handleLogout}
        />
      )}
      {renderPage()}
      {!isLoginPage && <Footer />}
    </div>
  );
}

export default App;
