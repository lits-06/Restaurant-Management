interface TopNavBarProps {
  currentPage: 'home' | 'contact' | 'reservation' | 'menu';
  setCurrentPage: (page: 'home' | 'contact' | 'reservation' | 'menu') => void;
}

export default function TopNavBar({ currentPage, setCurrentPage }: TopNavBarProps) {
  return (
    <header className="w-full top-0 sticky bg-surface z-50 shadow-sm border-b border-outline-variant">
      <nav className="flex justify-between items-center px-margin-desktop h-16 w-full max-w-container-max mx-auto">
        
        {/* Click vào Logo quay về Home */}
        <div 
          onClick={() => setCurrentPage('home')} 
          className="font-headline-md text-headline-md font-bold text-primary cursor-pointer"
        >
          LuxeBistro
        </div>

        {/* Menu điều hướng */}
        <div className="hidden md:flex items-center gap-8 font-body-md text-body-md">
          <button
            onClick={() => setCurrentPage('home')}
            className={`transition-colors cursor-pointer active:opacity-80 pb-1 ${
              currentPage === 'home' 
                ? 'text-primary border-b-2 border-primary font-semibold' 
                : 'text-on-surface-variant hover:text-primary'
            }`}
          >
            Home
          </button>
          
          <button
            onClick={() => setCurrentPage('menu')}
            className={`transition-colors cursor-pointer active:opacity-80 pb-1 ${
              currentPage === 'menu' 
                ? 'text-primary border-b-2 border-primary font-semibold' 
                : 'text-on-surface-variant hover:text-primary'
            }`}
          >
            Menu
          </button>

          <button
            onClick={() => setCurrentPage('reservation')}
            className={`transition-colors cursor-pointer active:opacity-80 pb-1 ${
              currentPage === 'reservation' 
                ? 'text-primary border-b-2 border-primary font-semibold' 
                : 'text-on-surface-variant hover:text-primary'
            }`}
          >
            Reservation
          </button>

          <button
            onClick={() => setCurrentPage('contact')}
            className={`transition-colors cursor-pointer active:opacity-80 pb-1 ${
              currentPage === 'contact' 
                ? 'text-primary border-b-2 border-primary font-semibold' 
                : 'text-on-surface-variant hover:text-primary'
            }`}
          >
            Contact
          </button>
        </div>

        <button className="bg-primary-container text-on-primary-container px-6 py-2 rounded-lg font-label-sm hover:opacity-90 transition-all cursor-pointer active:scale-95">
          Book Now
        </button>
      </nav>
    </header>
  );
}