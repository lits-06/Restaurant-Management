import React, { useState } from 'react';

const MONTH_SHORT = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
const MONTH_FULL  = ['January', 'February', 'March', 'April', 'May', 'June',
                     'July', 'August', 'September', 'October', 'November', 'December'];

interface HeaderProps {
  year: number;
  month: number; // 0-indexed
  onChange: (year: number, month: number) => void;
}

const Header: React.FC<HeaderProps> = ({ year, month, onChange }) => {
  const now = new Date();
  const [pickerYear, setPickerYear] = useState(year);

  return (
    <header className="flex justify-between items-end mb-12">
      <div>
        <nav className="flex items-center gap-2 text-[#4d4635] text-xs font-semibold mb-2">
          <span>Admin</span>
          <span className="material-symbols-outlined text-[14px]">chevron_right</span>
          <span className="text-[#735c00]">Analytics Report</span>
        </nav>
        <h2 className="font-serif text-5xl font-bold text-[#191c1d]">Analytics Overview</h2>
      </div>

      <div className="flex items-center gap-4">
        {/* Month picker */}
        <div className="relative group">
          <button className="flex items-center bg-white border border-[#d0c5af]/50 rounded-lg px-4 py-2.5 shadow-sm hover:border-[#735c00] transition-all gap-3 min-w-[220px] text-left">
            <span className="material-symbols-outlined text-[#735c00] text-[20px]">calendar_month</span>
            <div className="flex flex-col">
              <span className="text-[10px] uppercase tracking-wider text-[#4d4635] opacity-70">Reporting Period</span>
              <span className="text-base text-[#191c1d] font-semibold">{MONTH_FULL[month]} {year}</span>
            </div>
            <span className="material-symbols-outlined text-[#4d4635] text-[18px] ml-auto group-hover:text-[#735c00] transition-colors">
              expand_more
            </span>
          </button>

          {/* Dropdown */}
          <div className="absolute top-full right-0 mt-3 bg-white border border-[#d0c5af]/20 rounded-xl shadow-2xl z-[60] overflow-hidden opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 origin-top-right scale-95 group-hover:scale-100 w-72">
            <div className="p-5">
              {/* Year nav */}
              <div className="flex justify-between items-center mb-4">
                <button
                  className="w-8 h-8 flex items-center justify-center rounded-full hover:bg-[#edeeef] text-[#4d4635] hover:text-[#735c00] transition-colors"
                  onClick={() => setPickerYear(y => y - 1)}
                >
                  <span className="material-symbols-outlined text-[20px]">chevron_left</span>
                </button>
                <span className="font-serif text-lg font-semibold text-[#191c1d]">{pickerYear}</span>
                <button
                  className="w-8 h-8 flex items-center justify-center rounded-full hover:bg-[#edeeef] text-[#4d4635] hover:text-[#735c00] transition-colors disabled:opacity-30 disabled:cursor-not-allowed"
                  onClick={() => setPickerYear(y => Math.min(y + 1, now.getFullYear()))}
                  disabled={pickerYear >= now.getFullYear()}
                >
                  <span className="material-symbols-outlined text-[20px]">chevron_right</span>
                </button>
              </div>

              {/* Month grid */}
              <div className="grid grid-cols-3 gap-2">
                {MONTH_SHORT.map((m, idx) => {
                  const isFuture = pickerYear === now.getFullYear() && idx > now.getMonth();
                  const isSelected = pickerYear === year && idx === month;
                  return (
                    <button
                      key={m}
                      disabled={isFuture}
                      onClick={() => onChange(pickerYear, idx)}
                      className={`py-2.5 rounded-lg text-xs font-semibold transition-all ${
                        isSelected
                          ? 'bg-[#d4af37] text-[#554300] shadow-sm'
                          : isFuture
                          ? 'opacity-25 cursor-not-allowed'
                          : 'hover:bg-[#edeeef] text-[#191c1d]'
                      }`}
                    >
                      {m}
                    </button>
                  );
                })}
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
