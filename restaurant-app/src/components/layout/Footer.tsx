export default function Footer() {
  return (
    <footer className="w-full mt-auto bg-surface-container-lowest dark:bg-surface-container-low border-t border-outline-variant dark:border-outline">
      <div className="flex flex-col md:flex-row justify-between items-center px-margin-desktop py-8 w-full max-w-container-max mx-auto gap-6">
        <div className="flex flex-col items-center md:items-start gap-2">
          <div className="text-xl font-bold text-primary">LuxeBistro</div>
          <p className="text-xs font-semibold text-on-secondary-container dark:text-secondary-fixed-dim">
            © 2024 LuxeBistro Hospitality Group
          </p>
        </div>
        <div className="flex flex-wrap justify-center gap-8">
          <a className="text-xs font-semibold text-on-secondary-container dark:text-secondary-fixed-dim hover:underline decoration-primary cursor-pointer transition-opacity" href="#">Privacy Policy</a>
          <a className="text-xs font-semibold text-on-secondary-container dark:text-secondary-fixed-dim hover:underline decoration-primary cursor-pointer transition-opacity" href="#">Terms of Service</a>
          <a className="text-xs font-semibold text-on-secondary-container dark:text-secondary-fixed-dim hover:underline decoration-primary cursor-pointer transition-opacity" href="#">Careers</a>
          <a className="text-xs font-semibold text-on-secondary-container dark:text-secondary-fixed-dim hover:underline decoration-primary cursor-pointer transition-opacity" href="#">Accessibility</a>
        </div>
        <div className="flex gap-4">
          <span className="material-symbols-outlined text-primary cursor-pointer hover:scale-110 transition-transform">share</span>
          <span className="material-symbols-outlined text-primary cursor-pointer hover:scale-110 transition-transform">location_on</span>
          <span className="material-symbols-outlined text-primary cursor-pointer hover:scale-110 transition-transform">mail</span>
        </div>
      </div>
    </footer>
  );
}