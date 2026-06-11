import { useAuthStore } from '../../store/authStore';

type Page = 'home' | 'contact' | 'reservation' | 'menu' | 'login' | 'my-orders';

interface TopNavBarProps {
  currentPage: Page;
  navigateTo: (page: Page) => void;
  onLogout: () => void;
}

export default function TopNavBar({ currentPage, navigateTo, onLogout }: TopNavBarProps) {
  const { user } = useAuthStore();

  const navItem = (label: string, page: Page) => (
    <button
      key={page}
      onClick={() => navigateTo(page)}
      className={`transition-colors cursor-pointer active:opacity-80 pb-1 ${
        currentPage === page
          ? 'text-primary border-b-2 border-primary font-semibold'
          : 'text-on-surface-variant hover:text-primary'
      }`}
    >
      {label}
    </button>
  );

  return (
    <header className="w-full top-0 sticky bg-surface z-50 shadow-sm border-b border-outline-variant">
      <nav className="flex justify-between items-center px-margin-desktop h-16 w-full max-w-container-max mx-auto">

        <div
          onClick={() => navigateTo('home')}
          className="font-headline-md text-headline-md font-bold text-primary cursor-pointer"
        >
          LuxeBistro
        </div>

        <div className="hidden md:flex items-center gap-8 font-body-md text-body-md">
          {navItem('Home', 'home')}
          {navItem('Menu', 'menu')}
          {navItem('Reservation', 'reservation')}
          {navItem('Contact', 'contact')}
          {user && navItem('My Orders', 'my-orders')}
        </div>

        <div className="flex items-center gap-3">
          {user ? (
            <>
              <span className="hidden sm:block text-sm text-on-surface-variant font-medium">
                {user.full_name || user.username}
              </span>
              <button
                onClick={onLogout}
                className="flex items-center gap-1.5 text-sm font-semibold text-on-surface-variant border border-outline-variant px-4 py-2 rounded-lg hover:bg-surface-container-low transition-all"
              >
                <span className="material-symbols-outlined text-base">logout</span>
                Sign Out
              </button>
            </>
          ) : (
            <button
              onClick={() => navigateTo('login')}
              className="text-sm font-semibold text-primary border border-primary/30 px-4 py-2 rounded-lg hover:bg-primary/5 transition-all"
            >
              Sign In
            </button>
          )}
        </div>
      </nav>
    </header>
  );
}
