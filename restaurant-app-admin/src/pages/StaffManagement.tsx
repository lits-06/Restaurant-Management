import React, { useEffect, useMemo, useState } from 'react';
import { staffApi, type StaffDto } from '../services/api';

interface Employee {
    id: string;
    name: string;
    role: string;
    contact: string;
    until?: string;
    avatar: string;
}

interface Shift {
    id: number;
    employeeId: string;
    type: 'Lunch' | 'Dinner';
}

interface StaffManagementProps {
    onNavigateToScheduler: () => void;
}

const getStaffId = (staff: StaffDto) => staff.staff_id ?? staff.staffId ?? '';
const mapStaff = (staff: StaffDto): Employee => ({
    id: getStaffId(staff),
    name: staff.name ?? 'Unnamed staff',
    role: staff.role ?? 'Staff',
    contact: staff.contact ?? '',
    avatar: staff.avatar || 'https://www.gstatic.com/labs-code/stitch/stitch-placeholder-300x300.svg',
});

const StaffManagement: React.FC<StaffManagementProps> = ({ onNavigateToScheduler }) => {
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [error, setError] = useState<string>('');

    const [shifts, setShifts] = useState<Shift[]>([]);

    // --- 2. STATE QUẢN LÝ UI/MODAL ---
    const [searchQuery, setSearchQuery] = useState<string>('');

    // State Modal Assign Shift
    const [isAssignModalOpen, setIsAssignModalOpen] = useState<boolean>(false);
    const [selectedStaffId, setSelectedStaffId] = useState<string>('');
    const [shiftType, setShiftType] = useState<'Lunch' | 'Dinner'>('Lunch');
    const [startTime, setStartTime] = useState<string>('');
    const [endTime, setEndTime] = useState<string>('');

    // State Modal Edit Employee
    const [isEditModalOpen, setIsEditModalOpen] = useState<boolean>(false);
    const [editingEmployee, setEditingEmployee] = useState<Employee | null>(null);

    const loadStaff = async () => {
        setIsLoading(true);
        setError('');

        try {
            const response = await staffApi.list({ page: 1, page_size: 100, keyword: searchQuery });
            const loaded = (response.staff ?? []).map(mapStaff).filter((staff) => staff.id);
            setEmployees(loaded);
            setSelectedStaffId((current) => current || loaded[0]?.id || '');
            setShifts((current) => current.length > 0 ? current : loaded.slice(0, 6).map((staff, index) => ({
                id: Date.now() + index,
                employeeId: staff.id,
                type: index < 3 ? 'Lunch' : 'Dinner',
            })));
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Không thể tải danh sách nhân viên từ máy chủ.');
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        loadStaff();
    }, []);

    // --- 3. LOGIC TÌM KIẾM & PHÂN LOẠI CA LÀM ---
    // Lọc danh sách sơ lược nhân viên ở mục Directory
    const filteredDirectory = useMemo(() => {
        return employees.filter(emp =>
            emp.name.toLowerCase().includes(searchQuery.toLowerCase())
        );
    }, [employees, searchQuery]);

    // Lấy danh sách nhân viên trực thuộc ca trưa (Lunch Shift)
    const lunchShiftStaff = useMemo(() => {
        return shifts
            .filter(s => s.type === 'Lunch')
            .map(s => employees.find(emp => emp.id === s.employeeId))
            .filter((emp): emp is Employee => !!emp);
    }, [shifts, employees]);

    // Lấy danh sách nhân viên trực thuộc ca tối (Dinner Shift)
    const dinnerShiftStaff = useMemo(() => {
        return shifts
            .filter(s => s.type === 'Dinner')
            .map(s => employees.find(emp => emp.id === s.employeeId))
            .filter((emp): emp is Employee => !!emp);
    }, [shifts, employees]);

    // Danh sách 4 nhân viên đang hiển thị trên On Duty Banner đầu trang
    const activeDutyStaff = useMemo(() => employees.slice(0, 4), [employees]);

    // --- 4. XỬ LÝ SỰ KIỆN (ACTIONS) ---
    const handleAssignShiftSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        const newShift: Shift = {
            id: Date.now(),
            employeeId: selectedStaffId,
            type: shiftType,
        };
        setShifts([...shifts, newShift]);
        setIsAssignModalOpen(false);
    };

    const handleEditClick = (employee: Employee) => {
        setEditingEmployee({ ...employee });
        setIsEditModalOpen(true);
    };

    const handleEditSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!editingEmployee) return;

        staffApi.update(editingEmployee.id, {
            name: editingEmployee.name,
            role: editingEmployee.role,
            contact: editingEmployee.contact,
            avatar: editingEmployee.avatar,
        })
            .then((response) => {
                const updated = response.staff ? mapStaff(response.staff) : editingEmployee;
                setEmployees(employees.map(emp => emp.id === editingEmployee.id ? updated : emp));
                setIsEditModalOpen(false);
                setEditingEmployee(null);
            })
            .catch((err) => setError(err instanceof Error ? err.message : 'Không thể cập nhật nhân viên.'));
    };

    const handleDeleteEmployee = (id: string) => {
        if (window.confirm('Are you sure you want to delete this employee?')) {
            staffApi.delete(id)
                .then(() => {
                    setEmployees(employees.filter(emp => emp.id !== id));
                    setShifts(shifts.filter(s => s.employeeId !== id));
                })
                .catch((err) => setError(err instanceof Error ? err.message : 'Không thể xóa nhân viên.'));
        }
    };

    return (
        <div className="flex flex-col animate-fadeIn w-full">
            {/* Upper Content Header */}
            <header className="h-24 px-margin-desktop flex items-center justify-between bg-surface-container-lowest sticky top-0 z-40 border-b border-outline-variant/10">
                <div>
                    <h2 className="font-headline-lg text-headline-lg text-on-surface">Staff Schedule</h2>
                    <p className="font-body-md text-body-md text-on-surface-variant">Wednesday, October 23rd, 2024</p>
                </div>
                <div>
                    <button
                        onClick={onNavigateToScheduler}
                        className="px-6 py-2.5 bg-primary text-on-primary font-label-sm text-label-sm rounded-lg flex items-center gap-2 hover:opacity-90 active:scale-[0.98] transition-all shadow-md"
                    >
                        <span className="material-symbols-outlined text-[18px]">calendar_add_on</span>
                        Assign Shift
                    </button>
                </div>
            </header>

            {/* Main Content Canvas */}
            <div className="px-margin-desktop py-8 max-w-container-max w-full mx-auto space-y-12">
                {error && <p className="text-sm text-error">{error}</p>}
                {isLoading && <p className="text-sm text-on-surface-variant">Đang tải danh sách từ staff-service...</p>}

                {/* Section: Currently On Duty */}
                <section>
                    <div className="flex items-center gap-4 mb-6">
                        <span className="flex h-3 w-3 rounded-full bg-primary relative">
                            <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
                        </span>
                        <h3 className="font-headline-md text-headline-md">Currently On Duty</h3>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                        {activeDutyStaff.map((staff) => (
                            <div
                                key={staff.id}
                                className="bg-surface-container-lowest border border-outline-variant rounded-xl p-5 flex items-center gap-4 staff-card-shadow border-t-4 border-t-primary hover:-translate-y-1 hover:active-glow transition-all duration-200 cursor-pointer"
                            >
                                <div className="w-12 h-12 rounded-full overflow-hidden bg-surface-container flex-shrink-0">
                                    <img alt={staff.name} className="w-full h-full object-cover" src={staff.avatar} />
                                </div>
                                <div>
                                    <p className="font-label-sm text-label-sm font-bold text-on-surface">{staff.name}</p>
                                    <p className="text-[11px] text-primary font-bold uppercase tracking-tighter">{staff.role}</p>
                                    {staff.until && <p className="text-[11px] text-on-surface-variant">Until {staff.until}</p>}
                                </div>
                            </div>
                        ))}
                    </div>
                </section>

                {/* Section: Today's Staff Shift Organization */}
                <section className="grid grid-cols-1 lg:grid-cols-2 gap-12">
                    {/* Lunch Shift Table */}
                    <div>
                        <div className="flex items-center justify-between mb-6 pb-2 border-b border-outline-variant">
                            <div className="flex items-center gap-3">
                                <span className="material-symbols-outlined text-primary">light_mode</span>
                                <h3 className="font-headline-md text-headline-md">Lunch Shift</h3>
                            </div>
                            <span className="bg-secondary-container text-on-secondary-container text-[11px] px-3 py-1 rounded-full font-bold">11:00 AM — 4:00 PM</span>
                        </div>

                        <table className="w-full border-collapse">
                            <thead>
                                <tr className="text-left border-b border-outline-variant">
                                    <th className="py-3 font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Employee</th>
                                    <th className="py-3 font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Contact</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-outline-variant/30">
                                {lunchShiftStaff.map((staff) => (
                                    <tr key={staff.id} className="group hover:bg-surface-container-low transition-colors">
                                        <td className="py-4">
                                            <div className="flex flex-col">
                                                <span className="font-body-md text-body-md font-semibold text-on-surface">{staff.name}</span>
                                                <span className="text-[12px] text-on-surface-variant">{staff.role}</span>
                                            </div>
                                        </td>
                                        <td className="py-4 font-body-md text-body-md text-on-surface-variant">{staff.contact}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>

                    {/* Dinner Shift Table */}
                    <div>
                        <div className="flex items-center justify-between mb-6 pb-2 border-b border-outline-variant">
                            <div className="flex items-center gap-3">
                                <span className="material-symbols-outlined text-primary">dark_mode</span>
                                <h3 className="font-headline-md text-headline-md">Dinner Shift</h3>
                            </div>
                            <span className="bg-secondary-container text-on-secondary-container text-[11px] px-3 py-1 rounded-full font-bold">5:00 PM — 12:00 AM</span>
                        </div>

                        <table className="w-full border-collapse">
                            <thead>
                                <tr className="text-left border-b border-outline-variant">
                                    <th className="py-3 font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Employee</th>
                                    <th className="py-3 font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider">Contact</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-outline-variant/30">
                                {dinnerShiftStaff.map((staff) => (
                                    <tr key={staff.id} className="group hover:bg-surface-container-low transition-colors">
                                        <td className="py-4">
                                            <div className="flex flex-col">
                                                <span className="font-body-md text-body-md font-semibold text-on-surface">{staff.name}</span>
                                                <span className="text-[12px] text-on-surface-variant">{staff.role}</span>
                                            </div>
                                        </td>
                                        <td className="py-4 font-body-md text-body-md text-on-surface-variant">{staff.contact}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </section>

                {/* Section: Employee Directory */}
                <section className="bg-surface-container-low rounded-2xl p-8 border border-outline-variant">
                    <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
                        <h3 className="font-headline-md text-headline-md">Employee Directory</h3>
                        <div className="relative">
                            <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[20px]">search</span>
                            <input
                                type="text"
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                className="pl-10 pr-4 py-2 bg-surface-container-lowest border border-outline-variant rounded-full text-label-sm font-label-sm focus:ring-1 focus:ring-primary focus:outline-none w-full sm:w-64"
                                placeholder="Search by name..."
                            />
                        </div>
                    </div>

                    <div className="overflow-x-auto">
                        <table className="w-full border-collapse">
                            <thead>
                                <tr className="text-left border-b border-outline-variant">
                                    <th className="py-3 font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider px-4">Name</th>
                                    <th className="py-3 font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider px-4">Contact</th>
                                    <th className="py-3 font-label-sm text-label-sm text-on-surface-variant uppercase tracking-wider text-right px-4">Actions</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-outline-variant/30">
                                {filteredDirectory.length === 0 ? (
                                    <tr>
                                        <td colSpan={3} className="text-center py-6 text-sm text-on-surface-variant italic">
                                            No employees found match your search criteria.
                                        </td>
                                    </tr>
                                ) : (
                                    filteredDirectory.map((emp) => (
                                        <tr key={emp.id} className="group hover:bg-surface-container-lowest transition-colors">
                                            <td className="py-4 px-4 font-body-md text-body-md font-semibold text-on-surface">{emp.name}</td>
                                            <td className="py-4 px-4 font-body-md text-body-md text-on-surface-variant">{emp.contact}</td>
                                            <td className="py-4 px-4 text-right">
                                                <div className="flex justify-end gap-3">
                                                    <button
                                                        onClick={() => handleEditClick(emp)}
                                                        className="material-symbols-outlined text-on-surface-variant hover:text-primary transition-colors cursor-pointer text-[20px]"
                                                    >
                                                        edit
                                                    </button>
                                                    <button
                                                        onClick={() => handleDeleteEmployee(emp.id)}
                                                        className="material-symbols-outlined text-on-surface-variant hover:text-error transition-colors cursor-pointer text-[20px]"
                                                    >
                                                        delete
                                                    </button>
                                                </div>
                                            </td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    </div>
                </section>
            </div>

            {/* --- MODAL 1: ASSIGN SHIFT (REGISTER SHIFT) --- */}
            {isAssignModalOpen && (
                <div className="fixed inset-0 bg-inverse-surface/40 backdrop-blur-sm z-[100] flex items-center justify-center p-4 transition-all animate-fadeIn">
                    <div className="bg-surface-container-lowest w-full max-w-md p-8 rounded-2xl shadow-2xl transition-all scale-100">
                        <div className="flex justify-between items-center mb-6">
                            <h3 className="font-headline-md text-headline-md">Register Shift</h3>
                            <button onClick={() => setIsAssignModalOpen(false)} className="material-symbols-outlined text-on-surface-variant hover:text-on-surface">close</button>
                        </div>

                        <form onSubmit={handleAssignShiftSubmit} className="space-y-5">
                            <div>
                                <label className="block font-label-sm text-label-sm text-on-surface-variant mb-2 uppercase tracking-widest">Select Staff Member</label>
                                <select
                                    value={selectedStaffId}
                                    onChange={(e) => setSelectedStaffId(e.target.value)}
                                    className="w-full bg-surface-container-low border border-outline-variant rounded-lg px-4 py-3 font-body-md text-body-md focus:ring-1 focus:ring-primary focus:border-primary outline-none"
                                >
                                    {employees.map(emp => (
                                        <option key={emp.id} value={emp.id}>{emp.name} ({emp.role})</option>
                                    ))}
                                </select>
                            </div>

                            <div>
                                <label className="block font-label-sm text-label-sm text-on-surface-variant mb-2 uppercase tracking-widest">Shift Type</label>
                                <div className="grid grid-cols-2 gap-4">
                                    <label className={`flex items-center gap-3 p-4 border rounded-lg cursor-pointer transition-colors ${shiftType === 'Lunch' ? 'border-primary bg-primary-container/10' : 'border-outline-variant hover:bg-primary-container/5'}`}>
                                        <input
                                            type="radio"
                                            name="shiftType"
                                            checked={shiftType === 'Lunch'}
                                            onChange={() => setShiftType('Lunch')}
                                            className="text-primary focus:ring-primary"
                                        />
                                        <span className="font-body-md text-body-md">Lunch</span>
                                    </label>
                                    <label className={`flex items-center gap-3 p-4 border rounded-lg cursor-pointer transition-colors ${shiftType === 'Dinner' ? 'border-primary bg-primary-container/10' : 'border-outline-variant hover:bg-primary-container/5'}`}>
                                        <input
                                            type="radio"
                                            name="shiftType"
                                            checked={shiftType === 'Dinner'}
                                            onChange={() => setShiftType('Dinner')}
                                            className="text-primary focus:ring-primary"
                                        />
                                        <span className="font-body-md text-body-md">Dinner</span>
                                    </label>
                                </div>
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block font-label-sm text-label-sm text-on-surface-variant mb-2 uppercase tracking-widest">Start Time</label>
                                    <input
                                        type="time"
                                        value={startTime}
                                        onChange={(e) => setStartTime(e.target.value)}
                                        className="w-full bg-surface-container-low border border-outline-variant rounded-lg px-4 py-3 font-body-md text-body-md focus:ring-1 focus:ring-primary outline-none"
                                    />
                                </div>
                                <div>
                                    <label className="block font-label-sm text-label-sm text-on-surface-variant mb-2 uppercase tracking-widest">End Time</label>
                                    <input
                                        type="time"
                                        value={endTime}
                                        onChange={(e) => setEndTime(e.target.value)}
                                        className="w-full bg-surface-container-low border border-outline-variant rounded-lg px-4 py-3 font-body-md text-body-md focus:ring-1 focus:ring-primary outline-none"
                                    />
                                </div>
                            </div>

                            <button type="submit" className="w-full bg-primary text-on-primary font-label-sm text-label-sm py-4 rounded-lg shadow-lg hover:opacity-90 active:scale-[0.98] transition-all mt-4">
                                Confirm Registration
                            </button>
                        </form>
                    </div>
                </div>
            )}

            {/* --- MODAL 2: EDIT EMPLOYEE INFORMATION --- */}
            {isEditModalOpen && editingEmployee && (
                <div className="fixed inset-0 bg-inverse-surface/40 backdrop-blur-sm z-[100] flex items-center justify-center p-4 transition-all animate-fadeIn">
                    <div className="bg-surface-container-lowest w-full max-w-md p-8 rounded-2xl shadow-2xl transition-all scale-100">
                        <div className="flex justify-between items-center mb-6">
                            <h3 className="font-headline-md text-headline-md font-playfair">Edit Employee Information</h3>
                            <button onClick={() => setIsEditModalOpen(false)} className="material-symbols-outlined text-on-surface-variant hover:text-on-surface">close</button>
                        </div>

                        <form onSubmit={handleEditSubmit} className="space-y-5">
                            <div>
                                <label className="block font-label-sm text-label-sm text-on-surface-variant mb-2 uppercase tracking-widest">Full Name</label>
                                <input
                                    type="text"
                                    value={editingEmployee.name}
                                    onChange={(e) => setEditingEmployee({ ...editingEmployee, name: e.target.value })}
                                    className="w-full bg-surface-container-low border border-outline-variant rounded-lg px-4 py-3 font-body-md text-body-md focus:ring-1 focus:ring-primary focus:border-primary outline-none"
                                    placeholder="Enter employee name"
                                    required
                                />
                            </div>
                            <div>
                                <label className="block font-label-sm text-label-sm text-on-surface-variant mb-2 uppercase tracking-widest">Phone Number</label>
                                <input
                                    type="tel"
                                    value={editingEmployee.contact}
                                    onChange={(e) => setEditingEmployee({ ...editingEmployee, contact: e.target.value })}
                                    className="w-full bg-surface-container-low border border-outline-variant rounded-lg px-4 py-3 font-body-md text-body-md focus:ring-1 focus:ring-primary focus:border-primary outline-none"
                                    placeholder="+1 (555) 000-0000"
                                    required
                                />
                            </div>

                            <div className="flex gap-4 mt-6">
                                <button
                                    type="button"
                                    onClick={() => setIsEditModalOpen(false)}
                                    className="flex-1 py-3 px-4 border border-outline-variant text-on-surface font-label-sm text-label-sm rounded-lg hover:bg-surface-container-low transition-all"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    className="flex-1 bg-primary text-on-primary font-label-sm text-label-sm py-3 px-4 rounded-lg shadow-lg hover:opacity-90 active:scale-[0.98] transition-all"
                                >
                                    Save Changes
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default StaffManagement;
