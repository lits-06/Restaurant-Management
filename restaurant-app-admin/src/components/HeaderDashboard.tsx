import React from 'react';

interface HeaderProps {
  selectedMonth: string;
  setSelectedMonth: (month: string) => void;
}

const Header: React.FC<HeaderProps> = ({ selectedMonth, setSelectedMonth }) => {
  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

  return (
    <header className="flex justify-between items-end mb-12">
      <div>
        <nav className="flex items-center gap-2 text-[#4d4635] text-xs font-semibold mb-2">
          <span>Admin</span>
          <span className="material-symbols-outlined text-[14px]">chevron_right</span>
          <span className="text-[#735c00]">Detailed Reports</span>
        </nav>
        <h2 className="font-serif text-5xl font-bold text-[#191c1d]">Analytics Overview</h2>
      </div>

      <div className="flex items-center gap-4">
        {/* Date Range Picker Dropdown */}
        <div className="relative group">
          <button className="flex items-center bg-white border border-[#d0c5af]/50 rounded-lg px-4 py-2.5 shadow-sm hover:border-[#735c00] transition-all gap-3 min-w-[280px] text-left">
            <span className="material-symbols-outlined text-[#735c00] text-[20px]">calendar_month</span>
            <div className="flex flex-col">
              <span className="text-[10px] uppercase tracking-wider text-[#4d4635] opacity-70">Select Month</span>
              <span className="text-base text-[#191c1d] font-semibold">{selectedMonth}</span>
            </div>
            <span className="material-symbols-outlined text-[#4d4635] text-[18px] ml-auto group-hover:text-[#735c00] transition-colors">
              expand_more
            </span>
          </button>

          {/* Dropdown Menu */}
          <div className="absolute top-full right-0 mt-3 bg-white border border-[#d0c5af]/20 rounded-xl shadow-2xl z-[60] overflow-hidden opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-300 origin-top-right transform scale-95 group-hover:scale-100 w-fit">
            <div className="p-6 bg-white w-80">
              <div className="flex justify-between items-center mb-6">
                <button className="material-symbols-outlined text-[#4d4635] hover:text-[#735c00] transition-colors">chevron_left</button>
                <span className="font-serif text-xl font-semibold text-[#191c1d]">2023</span>
                <button className="material-symbols-outlined text-[#4d4635] hover:text-[#735c00] transition-colors">chevron_right</button>
              </div>
              <div className="grid grid-cols-3 gap-3">
                {months.map((m) => {
                  const isSelected = selectedMonth.startsWith(m) || (m === 'Oct' && selectedMonth.includes('October'));
                  return (
                    <button
                      key={m}
                      onClick={() => setSelectedMonth(`${m === 'Oct' ? 'October' : m} 2023`)}
                      className={`py-3 rounded-lg text-xs font-semibold transition-all ${
                        isSelected
                          ? 'bg-[#d4af37] text-[#554300] font-bold shadow-sm ring-1 ring-[#735c00]/20'
                          : 'hover:bg-[#edeeef]'
                      }`}
                    >
                      {m}
                    </button>
                  );
                })}
              </div>
              <div className="mt-6 pt-4 border-t border-[#d0c5af]/10 flex justify-end">
                <button className="px-6 py-2 text-xs font-semibold bg-[#735c00] text-white rounded-lg shadow-md hover:opacity-90 active:scale-95 transition-all">
                  Confirm Selection
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Export Button */}
        <button className="bg-[#735c00] text-white text-xs font-semibold px-6 py-3 rounded-lg flex items-center gap-2 hover:opacity-90 active:scale-95 transition-all shadow-md">
          <span className="material-symbols-outlined">download</span>
          Export Report
        </button>
      </div>
    </header>
  );
};

export default Header;