import React from 'react';

const Footer: React.FC = () => {
  return (
    <footer className="flex flex-col md:flex-row justify-between items-center py-8 w-full border-t border-[#d0c5af] mt-16">
      <p className="text-xs text-[#5d6466] opacity-60">
        © 2024 LuxeBistro Hospitality Group
      </p>
      <div className="flex gap-6 mt-4 md:mt-0">
        <a className="text-xs text-[#5d6466] hover:underline decoration-[#735c00] transition-opacity" href="#privacy">
          Privacy Policy
        </a>
        <a className="text-xs text-[#5d6466] hover:underline decoration-[#735c00] transition-opacity" href="#terms">
          Terms of Service
        </a>
        <a className="text-xs text-[#5d6466] hover:underline decoration-[#735c00] transition-opacity" href="#accessibility">
          Accessibility
        </a>
      </div>
    </footer>
  );
};

export default Footer;