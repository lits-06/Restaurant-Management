import React, { useState, useMemo } from 'react';

interface Employee {
  id: number;
  name: string;
  role: string;
  avatar: string;
}

interface ShiftAssignment {
  employeeId: number;
  dayIndex: number; // 0: Mon, 1: Tue, 2: Wed, ..., 6: Sun
  shiftType: 'Lunch' | 'Dinner';
}

interface WeeklySchedulerProps {
  onBack: () => void;
}

const WeeklyScheduler: React.FC<WeeklySchedulerProps> = ({ onBack }) => {
  // --- 1. DANH SÁCH NHÂN VIÊN ---
  const [employees] = useState<Employee[]>([
    { id: 1, name: 'Marcus Vance', role: 'Floor Manager', avatar: 'https://lh3.googleusercontent.com/aida-public/AB6AXuAwt2OoNL8JlNjKbO-PatCMVjdZPnCkxWbPudcq2DKWsv1M9LEkH7yaGR9v_5_ksX05XLxQi0YRCpm3jeqPQeaFLUKdmF2hFLUZQar9xccTYa7R8Bu3q62UhcfrGNejJfwlV0eukhNTukoXQ2HLDrSAw0BEXSCR2aJohBlcgk8kTfWhmZSQPltAFwmU6qIMb2b5UlC54AxF5heowbGGOdlsRPtOVnSQ4wBHwqd6xzaHsRBM9ww7K8pdGO0VNQTFWRibwZjjh40JhNM' },
    { id: 2, name: 'Elena Rodriguez', role: 'Executive Chef', avatar: 'https://lh3.googleusercontent.com/aida-public/AB6AXuCVCurC_obnRjWkEXqHam0yMqr3fXVO5ycEQ9FALxgySnji17b-e-lMRionHfOxS3jtcQFpE_r8tnzzcfcCq7I6deNbh8Zo43JTlR4C5HUw_1KBgTTqTYiZ_d5EWHcpY-RpCkc8CeDCJM7ilfAo_68Jv_k2_FoMH7wKKTTEgnB3X6fS12Tg2gFLVvPTkI77KxYUNE-VQcLhU4hFWH5eYUFoznD0RteuBiunv1IpZfl87Je--rc8YS1QB2nEN5oyqr38EgZhSXnWCmg' },
    { id: 3, name: 'Julian Schmidt', role: 'Sommelier', avatar: 'https://lh3.googleusercontent.com/aida-public/AB6AXuCxfdcU44u5KlWLvRwmvqZytikLbpmJLvAdm6xVttE9aXNYyhcWsxnI47DV2RqfJ2v3QEaBkAPKGNBrwCy-zhuSD_OTijK9oe5d2Eild9NDTKE7yzv9t9zs7IBGalAZAGBE8rZEtSJcG5hch2Mfuqsz9GUyxIjilTX4lV-3iYu1x3BtwTUiZeZjEZoQNZjVsajPVfHelclXPK3hIAOG9vOoICSF32i21x-YOkxeS5_2kmCdc-V-6gb0ER1HBzRUsMpRhWd1B7PI3Ls' },
    { id: 4, name: 'Sarah Lin', role: 'Lead Server', avatar: 'https://lh3.googleusercontent.com/aida-public/AB6AXuC8S1jzbxuTxDiDCEnXfFfMQ022ly-aAu_HQLb3OAKX0LUbkJFwExp0fUCJeaNRrvds5-NuluGo752CYe-Hskx03Ui_NCTTcdBZ-ZU0h72ZGSlfmu5YxvKWvyN_aYKLDMd3ylXZ8zO0_GvyNG-FAUmXO0XCAI9yGkQ5Ir70Q5cKHFKEu7n6d5f5sNyh3TYroY-AIA4pYOmFDqvxuJSee81x_f627np7_dyJhz972P8rCl_0pVzanCXnqk3nOavFs1x66RmsrW1Yz7M' },
  ]);

  // --- 2. STATES QUẢN LÝ ĐIỀU HƯỚNG & LỊCH TRÌNH ---
  const [selectedEmployeeId, setSelectedEmployeeId] = useState<number>(1);
  const [searchQuery, setSearchQuery] = useState<string>('');
  
  // Lưu trữ danh sách các ca đã gán (Mock dữ liệu mẫu ban đầu)
  const [assignments, setAssignments] = useState<ShiftAssignment[]>([
    { employeeId: 1, dayIndex: 0, shiftType: 'Lunch' },  // Marcus trực Trưa Thứ 2
    { employeeId: 1, dayIndex: 1, shiftType: 'Lunch' },  // Marcus trực Trưa Thứ 3
    { employeeId: 2, dayIndex: 3, shiftType: 'Dinner' }, // Elena trực Tối Thứ 5
  ]);

  // Cấu trúc danh sách các ngày trong tuần
  const days = [
    { label: 'Mon', date: '21', isToday: false },
    { label: 'Tue', date: '22', isToday: false },
    { label: 'Wed', date: '23', isToday: true }, // Ngày hiện tại được highlight
    { label: 'Thu', date: '24', isToday: false },
    { label: 'Fri', date: '25', isToday: false },
    { label: 'Sat', date: '26', isToday: false },
    { label: 'Sun', date: '27', isToday: false, isSunday: true },
  ];

  // --- 3. BỘ LỌC TÌM KIẾM NHÂN VIÊN ---
  const filteredEmployees = useMemo(() => {
    return employees.filter(emp =>
      emp.name.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [employees, searchQuery]);

  // --- 4. LOGIC XỬ LÝ CLICK Ô LỊCH (GRID CELL TOGGLE) ---
  const checkIsAssigned = (employeeId: number, dayIndex: number, shiftType: 'Lunch' | 'Dinner') => {
    return assignments.some(
      item => item.employeeId === employeeId && item.dayIndex === dayIndex && item.shiftType === shiftType
    );
  };

  const handleCellClick = (dayIndex: number, shiftType: 'Lunch' | 'Dinner') => {
    const isAlreadyAssigned = checkIsAssigned(selectedEmployeeId, dayIndex, shiftType);

    if (isAlreadyAssigned) {
      // Nếu đã gán ca trước đó -> Click để xóa gán lịch
      setAssignments(prev =>
        prev.filter(
          item => !(item.employeeId === selectedEmployeeId && item.dayIndex === dayIndex && item.shiftType === shiftType)
        )
      );
    } else {
      // Nếu chưa gán ca -> Thêm mới vào lịch trình nhân viên hiện tại
      setAssignments(prev => [...prev, { employeeId: selectedEmployeeId, dayIndex, shiftType }]);
    }
  };

  return (
    <div className="bg-surface text-on-surface min-h-screen flex w-full animate-fadeIn">
      <main className="flex-1 flex flex-col min-h-screen h-screen overflow-hidden">
        
        {/* Header thành phần */}
        <header className="h-24 px-margin-desktop flex items-center justify-between bg-surface-container-lowest sticky top-0 z-40 border-b border-outline-variant">
          <div className="flex items-center gap-4">
            <button
              onClick={onBack}
              className="p-2 hover:bg-surface-container-low rounded-full transition-colors flex items-center justify-center"
            >
              <span className="material-symbols-outlined">arrow_back</span>
            </button>
            <div>
              <h2 className="font-headline-lg text-headline-lg text-on-surface">Weekly Staff Scheduler</h2>
              <p className="font-body-md text-body-md text-on-surface-variant">Manage and assign shifts for the upcoming week</p>
            </div>
          </div>
          
          <div className="flex items-center gap-4">
            <div className="flex items-center bg-surface-container-low px-4 py-2 rounded-lg border border-outline-variant">
              <span className="material-symbols-outlined text-on-surface-variant mr-2">calendar_today</span>
              <span className="font-body-md text-body-md font-semibold">Oct 21 — Oct 27, 2024</span>
              <button className="ml-4 material-symbols-outlined text-on-surface-variant hover:text-primary transition-colors">chevron_left</button>
              <button className="ml-2 material-symbols-outlined text-on-surface-variant hover:text-primary transition-colors">chevron_right</button>
            </div>
          </div>
        </header>

        {/* Scheduler Content Body */}
        <div className="flex flex-1 overflow-hidden">
          
          {/* Left Sidebar: Employee Selector */}
          <aside className="w-80 border-r border-outline-variant bg-surface-container-lowest flex flex-col flex-shrink-0">
            <div className="p-6 border-b border-outline-variant">
              <h3 className="font-label-sm text-label-sm text-on-surface-variant uppercase tracking-widest mb-4">Select Staff Member</h3>
              <div className="relative">
                <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[20px]">search</span>
                <input 
                  type="text"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 pr-4 py-2.5 w-full bg-surface-container-low border border-outline-variant rounded-lg text-label-sm font-label-sm focus:ring-1 focus:ring-primary focus:outline-none" 
                  placeholder="Search employees..." 
                />
              </div>
            </div>
            
            <div className="flex-1 overflow-y-auto p-4 space-y-2">
              {filteredEmployees.map((emp) => {
                const isSelected = emp.id === selectedEmployeeId;
                return (
                  <div 
                    key={emp.id}
                    onClick={() => setSelectedEmployeeId(emp.id)}
                    className={`p-3 border rounded-xl flex items-center gap-3 cursor-pointer transition-all ${
                      isSelected 
                        ? 'bg-primary-container/20 border-primary' 
                        : 'hover:bg-surface-container-low border-transparent'
                    }`}
                  >
                    <div className="w-10 h-10 rounded-full bg-surface-container overflow-hidden shrink-0">
                      <img alt={emp.name} className="w-full h-full object-cover" src={emp.avatar} />
                    </div>
                    <div className="flex-1 overflow-hidden">
                      <p className="font-label-sm text-label-sm font-bold text-on-surface truncate">{emp.name}</p>
                      <p className={`text-[11px] font-bold uppercase truncate ${isSelected ? 'text-primary' : 'text-on-surface-variant'}`}>
                        {emp.role}
                      </p>
                    </div>
                    {isSelected && (
                      <span className="material-symbols-outlined text-primary text-[18px]">check_circle</span>
                    )}
                  </div>
                );
              })}
            </div>
          </aside>

          {/* Main Scheduler Grid */}
          <section className="flex-1 p-8 overflow-y-auto">
            <div className="bg-surface-container-lowest rounded-2xl border border-outline-variant shadow-sm overflow-hidden">
              
              {/* Days Headers */}
              <div className="grid grid-cols-8 border-b border-outline-variant bg-surface-container-low">
                <div className="p-4 border-r border-outline-variant flex flex-col justify-center bg-surface-container-lowest">
                  <span className="font-label-sm text-label-sm text-on-surface-variant uppercase tracking-widest text-center">Shift</span>
                </div>
                {days.map((day, idx) => (
                  <div 
                    key={idx} 
                    className={`p-4 border-r border-outline-variant text-center last:border-r-0 ${
                      day.isToday ? 'bg-primary-container/10' : ''
                    }`}
                  >
                    <p className={`font-label-sm text-label-sm uppercase mb-1 ${
                      day.isToday ? 'text-primary' : day.isSunday ? 'text-error' : 'text-on-surface-variant'
                    }`}>
                      {day.label}
                    </p>
                    <p className={`font-headline-md text-headline-md ${day.isToday ? 'text-primary' : 'text-on-surface'}`}>
                      {day.date}
                    </p>
                  </div>
                ))}
              </div>

              {/* Grid Rows */}
              
              {/* 1. Lunch Shift Row */}
              <div className="grid grid-cols-8 border-b border-outline-variant">
                <div className="p-6 border-r border-outline-variant flex flex-col justify-center bg-surface-container-low/50">
                  <div className="flex items-center gap-2 mb-1">
                    <span className="material-symbols-outlined text-primary text-[20px]">light_mode</span>
                    <span className="font-body-md text-body-md font-bold">Lunch</span>
                  </div>
                  <span className="text-[9px] text-on-surface-variant uppercase font-semibold tracking-wider">11:00 — 16:00</span>
                </div>
                
                {days.map((_, idx) => {
                  const isActive = checkIsAssigned(selectedEmployeeId, idx, 'Lunch');
                  return (
                    <div 
                      key={idx}
                      onClick={() => handleCellClick(idx, 'Lunch')}
                      className={`grid-cell border-r border-outline-variant last:border-r-0 hover:bg-surface-container-low cursor-pointer transition-all min-h-[160px] flex flex-col items-center justify-center p-2 group relative ${
                        isActive ? 'grid-cell-active' : ''
                      }`}
                    >
                      {!isActive && (
                        <span className="material-symbols-outlined text-outline-variant opacity-0 group-hover:opacity-100 transition-opacity">
                          add_circle
                        </span>
                      )}
                    </div>
                  );
                })}
              </div>

              {/* 2. Dinner Shift Row */}
              <div className="grid grid-cols-8">
                <div className="p-6 border-r border-outline-variant flex flex-col justify-center bg-surface-container-low/50">
                  <div className="flex items-center gap-2 mb-1">
                    <span className="material-symbols-outlined text-primary text-[20px]">dark_mode</span>
                    <span className="font-body-md text-body-md font-bold">Dinner</span>
                  </div>
                  <span className="text-[9px] text-on-surface-variant uppercase font-semibold tracking-wider">17:00 — 00:00</span>
                </div>
                
                {days.map((_, idx) => {
                  const isActive = checkIsAssigned(selectedEmployeeId, idx, 'Dinner');
                  return (
                    <div 
                      key={idx}
                      onClick={() => handleCellClick(idx, 'Dinner')}
                      className={`grid-cell border-r border-outline-variant last:border-r-0 hover:bg-surface-container-low cursor-pointer transition-all min-h-[160px] flex flex-col items-center justify-center p-2 group relative ${
                        isActive ? 'grid-cell-active' : ''
                      }`}
                    >
                      {!isActive && (
                        <span className="material-symbols-outlined text-outline-variant opacity-0 group-hover:opacity-100 transition-opacity">
                          add_circle
                        </span>
                      )}
                    </div>
                  );
                })}
              </div>

            </div>
          </section>
        </div>
      </main>
    </div>
  );
};

export default WeeklyScheduler;