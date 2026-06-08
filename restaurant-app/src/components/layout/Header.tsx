export default function Header() {
  return (
    <header className="w-full top-0 sticky z-50 bg-surface border-b border-outline-variant shadow-sm">
      <nav className="flex justify-between items-center px-margin-desktop h-16 w-full max-w-container-max mx-auto">
        <div className="font-headline-md text-headline-md font-bold text-primary cursor-pointer">
          LuxeBistro
        </div>

        <div className="hidden md:flex gap-8 items-center">
          <a
            href="#"
            className="font-body-md text-body-md text-primary border-b-2 border-primary pb-1"
          >
            Home
          </a>

          <a
            href="#"
            className="font-body-md text-body-md text-on-surface-variant hover:text-primary"
          >
            Menu
          </a>

          <a
            href="#"
            className="font-body-md text-body-md text-on-surface-variant hover:text-primary"
          >
            Reservation
          </a>

          <a
            href="#"
            className="font-body-md text-body-md text-on-surface-variant hover:text-primary"
          >
            Contact
          </a>
        </div>

        <button className="bg-primary text-on-primary px-6 py-2 rounded-lg font-label-sm text-label-sm hover:opacity-90">
          Book Now
        </button>
      </nav>
    </header>
  );
}