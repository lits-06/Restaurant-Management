import React, { useCallback, useEffect, useState } from 'react';
import { scheduleApi, usersApi, type ShiftDto, type UserDto } from '../services/api';

// ── helpers ──────────────────────────────────────────────────────────────────

const roleChip: Record<string, { bg: string; text: string; emoji: string }> = {
  CHEF:    { bg: 'bg-amber-100', text: 'text-amber-800', emoji: '🍳' },
  WAITER:  { bg: 'bg-blue-100',  text: 'text-blue-800',  emoji: '🛎' },
  MANAGER: { bg: 'bg-purple-100',text: 'text-purple-800',emoji: '👔' },
  ADMIN:   { bg: 'bg-green-100', text: 'text-green-800', emoji: '⚙️' },
};

function isoMonth(year: number, month: number) {
  return `${year}-${String(month + 1).padStart(2, '0')}`;
}

function daysInMonth(year: number, month: number) {
  return new Date(year, month + 1, 0).getDate();
}

// Returns Mon-indexed day (0=Mon … 6=Sun) for the 1st of the month
function firstWeekday(year: number, month: number) {
  const d = new Date(year, month, 1).getDay(); // 0=Sun
  return (d + 6) % 7; // shift so Mon=0
}

function calendarGrid(year: number, month: number) {
  const offset = firstWeekday(year, month);
  const total = daysInMonth(year, month);
  const cells: (number | null)[] = Array(offset).fill(null);
  for (let d = 1; d <= total; d++) cells.push(d);
  while (cells.length % 7 !== 0) cells.push(null);
  return cells;
}

function dateStr(year: number, month: number, day: number) {
  return `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
}

function displayName(user: UserDto) {
  return user.full_name || user.username || user.email || user.user_id || '?';
}

// ── types ─────────────────────────────────────────────────────────────────────

interface CreateModalState {
  date: string;
  userID: string;
  startTime: string;
  endTime: string;
  role: string;
  notes: string;
}

interface DetailState {
  shift: ShiftDto;
  userName: string;
}

interface EditModalState {
  shiftId: string;
  date: string;
  startTime: string;
  endTime: string;
  notes: string;
}

// ── component ─────────────────────────────────────────────────────────────────

const STAFF_ROLES = ['CHEF', 'WAITER', 'MANAGER', 'ADMIN'];

const MonthlyScheduler: React.FC = () => {
  const now = new Date();
  const [year, setYear]   = useState(now.getFullYear());
  const [month, setMonth] = useState(now.getMonth()); // 0-indexed

  const [shifts, setShifts]   = useState<ShiftDto[]>([]);
  const [users, setUsers]     = useState<UserDto[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError]     = useState('');

  // Filters
  const [roleFilter, setRoleFilter] = useState('');

  // Create modal
  const [createModal, setCreateModal] = useState<CreateModalState | null>(null);
  const [creating, setCreating] = useState(false);

  // Detail popover
  const [detail, setDetail] = useState<DetailState | null>(null);

  // Edit modal
  const [editModal, setEditModal] = useState<EditModalState | null>(null);
  const [editing, setEditing] = useState(false);

  const userMap = new Map(users.map((u) => [u.user_id ?? '', u]));

  // ── data loading ─────────────────────────────────────────────
  const loadData = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const [shiftsRes, usersRes] = await Promise.all([
        scheduleApi.list({ month: isoMonth(year, month) }),
        usersApi.listAll(),
      ]);
      setShifts(shiftsRes.shifts ?? []);
      setUsers((usersRes.users ?? []).filter((u) => u.roles?.some((r) => STAFF_ROLES.includes(r))));
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load data');
    } finally {
      setLoading(false);
    }
  }, [year, month]);

  useEffect(() => { loadData(); }, [loadData]);

  // ── navigation ───────────────────────────────────────────────
  const prevMonth = () => {
    if (month === 0) { setYear(y => y - 1); setMonth(11); }
    else setMonth(m => m - 1);
  };
  const nextMonth = () => {
    if (month === 11) { setYear(y => y + 1); setMonth(0); }
    else setMonth(m => m + 1);
  };

  // ── create ───────────────────────────────────────────────────
  const openCreateModal = (day: number) => {
    setCreateModal({
      date: dateStr(year, month, day),
      userID: users[0]?.user_id ?? '',
      startTime: '08:00',
      endTime: '16:00',
      role: 'CHEF',
      notes: '',
    });
  };

  const submitCreate = async () => {
    if (!createModal) return;
    setCreating(true);
    try {
      await scheduleApi.create({
        user_id:    createModal.userID,
        date:       createModal.date,
        start_time: createModal.startTime,
        end_time:   createModal.endTime,
        role:       createModal.role,
        notes:      createModal.notes,
      });
      setCreateModal(null);
      loadData();
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Failed to create shift');
    } finally {
      setCreating(false);
    }
  };

  // ── delete ───────────────────────────────────────────────────
  const deleteShift = async (shiftId: string) => {
    if (!confirm('Delete this shift?')) return;
    try {
      await scheduleApi.delete(shiftId);
      setDetail(null);
      loadData();
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Failed to delete shift');
    }
  };

  // ── edit ─────────────────────────────────────────────────────
  const openEditModal = (s: ShiftDto) => {
    setDetail(null);
    setEditModal({
      shiftId: s.shift_id ?? '',
      date: s.date ?? '',
      startTime: s.start_time ?? '',
      endTime: s.end_time ?? '',
      notes: s.notes ?? '',
    });
  };

  const submitEdit = async () => {
    if (!editModal) return;
    setEditing(true);
    try {
      await scheduleApi.update(editModal.shiftId, {
        date:       editModal.date,
        start_time: editModal.startTime,
        end_time:   editModal.endTime,
        notes:      editModal.notes,
      });
      setEditModal(null);
      loadData();
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Failed to save');
    } finally {
      setEditing(false);
    }
  };

  // ── render ───────────────────────────────────────────────────
  const cells = calendarGrid(year, month);
  const monthLabel = new Date(year, month, 1).toLocaleDateString('en-US', { month: 'long', year: 'numeric' });

  const shiftsOnDay = (day: number) => {
    const d = dateStr(year, month, day);
    return shifts.filter(
      (s) => s.date === d && (!roleFilter || s.role === roleFilter)
    );
  };

  return (
    <div className="flex flex-col gap-6">
      {/* Header */}
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 className="text-2xl font-bold text-[#191c1d]">Staff Schedule</h1>
          <p className="text-sm text-[#6b7280] mt-1 capitalize">{monthLabel}</p>
        </div>
        <div className="flex items-center gap-3">
          <select
            value={roleFilter}
            onChange={(e) => setRoleFilter(e.target.value)}
            className="text-sm border border-[#d0c5af] rounded-lg px-3 py-2 bg-white text-[#4d4635]"
          >
            <option value="">All Roles</option>
            {STAFF_ROLES.map((r) => <option key={r} value={r}>{r}</option>)}
          </select>
          <button onClick={prevMonth} className="px-3 py-2 border border-[#d0c5af] rounded-lg text-sm hover:bg-[#e1e3e4]">◀ Previous</button>
          <button onClick={nextMonth} className="px-3 py-2 border border-[#d0c5af] rounded-lg text-sm hover:bg-[#e1e3e4]">Next ▶</button>
        </div>
      </div>

      {error && <div className="text-red-600 text-sm bg-red-50 p-3 rounded-lg">{error}</div>}

      {/* Calendar grid */}
      <div className="bg-white rounded-2xl border border-[#e5e7eb] overflow-hidden shadow-sm">
        {/* Weekday headers */}
        <div className="grid grid-cols-7 border-b border-[#e5e7eb]">
          {['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'].map((d) => (
            <div key={d} className="text-center text-xs font-semibold text-[#6b7280] py-3 bg-[#f9fafb]">{d}</div>
          ))}
        </div>

        {/* Cells */}
        {loading ? (
          <div className="text-center py-16 text-[#6b7280]">Loading...</div>
        ) : (
          <div className="grid grid-cols-7">
            {cells.map((day, i) => (
              <div
                key={i}
                className={`min-h-[110px] border-b border-r border-[#e5e7eb] p-2 ${day ? 'bg-white hover:bg-[#fafafa]' : 'bg-[#f9fafb]'}`}
              >
                {day && (
                  <>
                    <div className="flex items-center justify-between mb-1">
                      <span className={`text-xs font-semibold ${dateStr(year, month, day) === new Date().toISOString().slice(0, 10) ? 'text-[#d4af37]' : 'text-[#374151]'}`}>
                        {day}
                      </span>
                      <button
                        onClick={() => openCreateModal(day)}
                        className="w-5 h-5 rounded-full bg-[#e5e7eb] hover:bg-[#d4af37] hover:text-white text-[#6b7280] text-xs flex items-center justify-center transition-colors"
                        title="Add shift"
                      >
                        +
                      </button>
                    </div>
                    <div className="flex flex-col gap-0.5 overflow-hidden">
                      {shiftsOnDay(day).map((s) => {
                        const chip = roleChip[s.role ?? ''] ?? { bg: 'bg-gray-100', text: 'text-gray-700', emoji: '' };
                        const u = userMap.get(s.user_id ?? '');
                        const name = u ? (u.full_name?.split(' ').pop() ?? u.username ?? '?') : (s.user_id?.slice(0, 6) ?? '?');
                        return (
                          <button
                            key={s.shift_id}
                            onClick={() => setDetail({ shift: s, userName: u ? displayName(u) : (s.user_id ?? '?') })}
                            className={`text-[10px] font-medium px-1.5 py-0.5 rounded truncate text-left ${chip.bg} ${chip.text}`}
                            title={`${displayName(u ?? {})} ${s.start_time}–${s.end_time}`}
                          >
                            {chip.emoji} {name}
                          </button>
                        );
                      })}
                    </div>
                  </>
                )}
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Create modal */}
      {createModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl shadow-xl p-6 w-full max-w-md">
            <h2 className="text-lg font-bold mb-4 text-[#191c1d]">Add Shift — {createModal.date}</h2>
            <div className="flex flex-col gap-3">
              <div>
                <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Staff Member</label>
                <select
                  value={createModal.userID}
                  onChange={(e) => setCreateModal({ ...createModal, userID: e.target.value })}
                  className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm"
                >
                  {users.map((u) => (
                    <option key={u.user_id} value={u.user_id}>{displayName(u)} ({u.roles?.join(', ')})</option>
                  ))}
                </select>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Start Time</label>
                  <input type="time" value={createModal.startTime} onChange={(e) => setCreateModal({ ...createModal, startTime: e.target.value })}
                    className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm" />
                </div>
                <div>
                  <label className="text-xs font-semibold text-[#6b7280] mb-1 block">End Time</label>
                  <input type="time" value={createModal.endTime} onChange={(e) => setCreateModal({ ...createModal, endTime: e.target.value })}
                    className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm" />
                </div>
              </div>
              <div>
                <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Role</label>
                <select value={createModal.role} onChange={(e) => setCreateModal({ ...createModal, role: e.target.value })}
                  className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm">
                  {STAFF_ROLES.map((r) => <option key={r} value={r}>{r}</option>)}
                </select>
              </div>
              <div>
                <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Notes</label>
                <input type="text" value={createModal.notes} onChange={(e) => setCreateModal({ ...createModal, notes: e.target.value })}
                  placeholder="e.g. Morning shift, substitute..." className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm" />
              </div>
            </div>
            <div className="flex gap-3 mt-5">
              <button onClick={submitCreate} disabled={creating}
                className="flex-1 bg-[#d4af37] text-white rounded-lg py-2 text-sm font-semibold hover:bg-[#b8962d] disabled:opacity-60">
                {creating ? 'Saving...' : 'Save Shift'}
              </button>
              <button onClick={() => setCreateModal(null)}
                className="flex-1 border border-[#d0c5af] rounded-lg py-2 text-sm font-semibold hover:bg-[#f3f4f5]">
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Detail popover */}
      {detail && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50" onClick={() => setDetail(null)}>
          <div className="bg-white rounded-2xl shadow-xl p-6 w-full max-w-sm" onClick={(e) => e.stopPropagation()}>
            <div className="flex items-start justify-between mb-3">
              <div>
                <p className="font-bold text-[#191c1d]">{detail.userName}</p>
                <p className="text-sm text-[#6b7280]">{detail.shift.date} · {detail.shift.start_time}–{detail.shift.end_time}</p>
              </div>
              <span className={`text-xs font-semibold px-2 py-1 rounded-full ${roleChip[detail.shift.role ?? '']?.bg ?? 'bg-gray-100'} ${roleChip[detail.shift.role ?? '']?.text ?? 'text-gray-700'}`}>
                {detail.shift.role}
              </span>
            </div>
            {detail.shift.notes && <p className="text-sm text-[#374151] mb-4">{detail.shift.notes}</p>}
            <div className="flex gap-2">
              <button
                onClick={() => openEditModal(detail.shift)}
                className="flex-1 bg-[#d4af37] text-white border border-[#d4af37] rounded-lg py-2 text-sm font-semibold hover:bg-[#b8962d]"
              >
                Edit Shift
              </button>
              <button
                onClick={() => deleteShift(detail.shift.shift_id ?? '')}
                className="flex-1 bg-red-50 text-red-600 border border-red-200 rounded-lg py-2 text-sm font-semibold hover:bg-red-100"
              >
                Delete Shift
              </button>
            </div>
          </div>
        </div>
      )}
      {/* Edit modal */}
      {editModal && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl shadow-xl p-6 w-full max-w-md">
            <h2 className="text-lg font-bold mb-4 text-[#191c1d]">Edit Shift</h2>
            <div className="flex flex-col gap-3">
              <div>
                <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Date</label>
                <input
                  type="date"
                  value={editModal.date}
                  onChange={(e) => setEditModal({ ...editModal, date: e.target.value })}
                  className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm"
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Start Time</label>
                  <input
                    type="time"
                    value={editModal.startTime}
                    onChange={(e) => setEditModal({ ...editModal, startTime: e.target.value })}
                    className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm"
                  />
                </div>
                <div>
                  <label className="text-xs font-semibold text-[#6b7280] mb-1 block">End Time</label>
                  <input
                    type="time"
                    value={editModal.endTime}
                    onChange={(e) => setEditModal({ ...editModal, endTime: e.target.value })}
                    className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm"
                  />
                </div>
              </div>
              <div>
                <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Notes</label>
                <input
                  type="text"
                  value={editModal.notes}
                  onChange={(e) => setEditModal({ ...editModal, notes: e.target.value })}
                  placeholder="e.g. Morning shift..."
                  className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm"
                />
              </div>
            </div>
            <div className="flex gap-3 mt-5">
              <button
                onClick={submitEdit}
                disabled={editing}
                className="flex-1 bg-[#d4af37] text-white rounded-lg py-2 text-sm font-semibold hover:bg-[#b8962d] disabled:opacity-60"
              >
                {editing ? 'Saving...' : 'Save Changes'}
              </button>
              <button
                onClick={() => setEditModal(null)}
                className="flex-1 border border-[#d0c5af] rounded-lg py-2 text-sm font-semibold hover:bg-[#f3f4f5]"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default MonthlyScheduler;
