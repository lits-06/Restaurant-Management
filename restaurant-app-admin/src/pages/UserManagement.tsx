import React, { useEffect, useState, useCallback } from 'react';
import { usersApi, type UserDto } from '../services/api';

const ALL_ROLES = ['USER', 'MANAGER', 'CHEF', 'WAITER', 'ADMIN'];

const ROLE_COLORS: Record<string, string> = {
  USER:    'bg-gray-100 text-gray-600',
  MANAGER: 'bg-purple-100 text-purple-700',
  CHEF:    'bg-amber-100 text-amber-700',
  WAITER:  'bg-blue-100 text-blue-700',
  ADMIN:   'bg-green-100 text-green-700',
};

const STATUS_COLORS: Record<string, string> = {
  ACTIVE:    'text-green-600',
  INACTIVE:  'text-gray-400',
  SUSPENDED: 'text-red-500',
};

interface CreateForm {
  email: string; password: string; username: string;
  full_name: string; phone: string;
}

interface EditForm {
  email: string; username: string; full_name: string; phone: string; status: string;
}

const emptyCreate = (): CreateForm => ({
  email: '', password: '', username: '', full_name: '', phone: '',
});

const UserManagement: React.FC = () => {
  const [users, setUsers] = useState<UserDto[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');

  // Create modal
  const [showCreate, setShowCreate] = useState(false);
  const [createForm, setCreateForm] = useState<CreateForm>(emptyCreate());
  const [creating, setCreating] = useState(false);
  const [createError, setCreateError] = useState('');

  // Edit modal
  const [editUser, setEditUser] = useState<UserDto | null>(null);
  const [editForm, setEditForm] = useState<EditForm>({ email: '', username: '', full_name: '', phone: '', status: '' });
  const [saving, setSaving] = useState(false);
  const [editError, setEditError] = useState('');

  // Roles modal
  const [rolesUser, setRolesUser] = useState<UserDto | null>(null);
  const [selectedRoles, setSelectedRoles] = useState<string[]>([]);
  const [savingRoles, setSavingRoles] = useState(false);

  // Change password modal
  const [pwUser, setPwUser] = useState<UserDto | null>(null);
  const [pwForm, setPwForm] = useState({ oldPw: '', newPw: '', confirmPw: '' });
  const [savingPw, setSavingPw] = useState(false);
  const [pwError, setPwError] = useState('');

  const load = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const res = await usersApi.listAll({ page_size: 200 });
      setUsers(res.users ?? []);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load data');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { load(); }, [load]);

  const filtered = users.filter((u) => {
    if (!search) return true;
    const q = search.toLowerCase();
    return (
      u.email?.toLowerCase().includes(q) ||
      u.username?.toLowerCase().includes(q) ||
      u.full_name?.toLowerCase().includes(q)
    );
  });

  // ── Create ──────────────────────────────────────────────────────
  const submitCreate = async () => {
    if (!createForm.email || !createForm.password || !createForm.username) {
      setCreateError('Email, password and username are required');
      return;
    }
    setCreating(true);
    setCreateError('');
    try {
      await usersApi.create(createForm);
      setShowCreate(false);
      setCreateForm(emptyCreate());
      load();
    } catch (e) {
      setCreateError(e instanceof Error ? e.message : 'Failed to create user');
    } finally {
      setCreating(false);
    }
  };

  // ── Edit ─────────────────────────────────────────────────────────
  const openEdit = (u: UserDto) => {
    setEditUser(u);
    setEditForm({
      email:     u.email ?? '',
      username:  u.username ?? '',
      full_name: u.full_name ?? '',
      phone:     u.phone ?? '',
      status:    u.status ?? 'ACTIVE',
    });
    setEditError('');
  };

  const submitEdit = async () => {
    if (!editUser?.user_id) return;
    setSaving(true);
    setEditError('');
    try {
      await usersApi.update(editUser.user_id, editForm);
      setEditUser(null);
      load();
    } catch (e) {
      setEditError(e instanceof Error ? e.message : 'Failed to save');
    } finally {
      setSaving(false);
    }
  };

  // ── Delete ───────────────────────────────────────────────────────
  const deleteUser = async (id: string) => {
    if (!confirm('Delete this user? This action cannot be undone.')) return;
    try {
      await usersApi.delete(id);
      load();
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Failed to delete');
    }
  };

  // ── Roles ────────────────────────────────────────────────────────
  const openRoles = (u: UserDto) => {
    setRolesUser(u);
    setSelectedRoles(u.roles ?? []);
  };

  const toggleRole = (role: string) => {
    setSelectedRoles((prev) =>
      prev.includes(role) ? prev.filter((r) => r !== role) : [...prev, role]
    );
  };

  const submitRoles = async () => {
    if (!rolesUser?.user_id) return;
    setSavingRoles(true);
    try {
      await usersApi.assignRole(rolesUser.user_id, selectedRoles);
      setRolesUser(null);
      load();
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Failed to update roles');
    } finally {
      setSavingRoles(false);
    }
  };

  // ── Change password ───────────────────────────────────────────────
  const openPw = (u: UserDto) => {
    setPwUser(u);
    setPwForm({ oldPw: '', newPw: '', confirmPw: '' });
    setPwError('');
  };

  const submitPw = async () => {
    if (!pwUser?.user_id) return;
    if (!pwForm.newPw || pwForm.newPw.length < 6) { setPwError('New password must be at least 6 characters'); return; }
    if (pwForm.newPw !== pwForm.confirmPw) { setPwError('Passwords do not match'); return; }
    setSavingPw(true);
    setPwError('');
    try {
      await usersApi.changePassword(pwUser.user_id, pwForm.oldPw, pwForm.newPw);
      setPwUser(null);
    } catch (e) {
      setPwError(e instanceof Error ? e.message : 'Failed to change password');
    } finally {
      setSavingPw(false);
    }
  };

  return (
    <div className="flex flex-col gap-6">
      {/* Header */}
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 className="text-2xl font-bold text-[#191c1d]">User Management</h1>
          <p className="text-sm text-[#6b7280] mt-1">{users.length} users</p>
        </div>
        <div className="flex items-center gap-3">
          <input
            type="text"
            placeholder="Search by email / name..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="text-sm border border-[#d0c5af] rounded-lg px-3 py-2 bg-white text-[#4d4635] w-56"
          />
          <button
            onClick={() => { setShowCreate(true); setCreateForm(emptyCreate()); setCreateError(''); }}
            className="flex items-center gap-2 bg-[#d4af37] text-white px-4 py-2 rounded-lg text-sm font-semibold hover:bg-[#b8962d] transition-colors"
          >
            <span className="material-symbols-outlined text-base">person_add</span>
            Add User
          </button>
        </div>
      </div>

      {error && <div className="text-red-600 text-sm bg-red-50 p-3 rounded-lg">{error}</div>}

      {/* User table */}
      <div className="bg-white rounded-2xl border border-[#e5e7eb] overflow-hidden shadow-sm">
        <table className="w-full text-sm">
          <thead className="bg-[#f9fafb] border-b border-[#e5e7eb]">
            <tr>
              <th className="text-left px-4 py-3 text-xs font-semibold text-[#6b7280] uppercase tracking-wider">User</th>
              <th className="text-left px-4 py-3 text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Roles</th>
              <th className="text-left px-4 py-3 text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Status</th>
              <th className="text-right px-4 py-3 text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-[#f3f4f5]">
            {loading ? (
              <tr><td colSpan={4} className="text-center py-12 text-[#6b7280]">Loading...</td></tr>
            ) : filtered.length === 0 ? (
              <tr><td colSpan={4} className="text-center py-12 text-[#6b7280]">No users found</td></tr>
            ) : filtered.map((u) => (
              <tr key={u.user_id} className="hover:bg-[#fafafa]">
                <td className="px-4 py-3">
                  <p className="font-semibold text-[#191c1d]">{u.full_name || u.username}</p>
                  <p className="text-xs text-[#6b7280]">{u.email}</p>
                  {u.phone && <p className="text-xs text-[#9ca3af]">{u.phone}</p>}
                </td>
                <td className="px-4 py-3">
                  <div className="flex flex-wrap gap-1">
                    {(u.roles ?? []).map((r) => (
                      <span key={r} className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${ROLE_COLORS[r] ?? 'bg-gray-100 text-gray-600'}`}>
                        {r}
                      </span>
                    ))}
                  </div>
                </td>
                <td className="px-4 py-3">
                  <span className={`text-xs font-semibold ${STATUS_COLORS[u.status ?? ''] ?? 'text-gray-500'}`}>
                    {u.status ?? '—'}
                  </span>
                </td>
                <td className="px-4 py-3">
                  <div className="flex items-center justify-end gap-1">
                    <button
                      onClick={() => openEdit(u)}
                      title="Edit"
                      className="p-1.5 rounded-lg hover:bg-[#f3f4f5] text-[#6b7280] hover:text-[#191c1d] transition-colors"
                    >
                      <span className="material-symbols-outlined text-base">edit</span>
                    </button>
                    <button
                      onClick={() => openRoles(u)}
                      title="Manage roles"
                      className="p-1.5 rounded-lg hover:bg-[#f3f4f5] text-[#6b7280] hover:text-[#191c1d] transition-colors"
                    >
                      <span className="material-symbols-outlined text-base">admin_panel_settings</span>
                    </button>
                    <button
                      onClick={() => openPw(u)}
                      title="Change password"
                      className="p-1.5 rounded-lg hover:bg-[#f3f4f5] text-[#6b7280] hover:text-[#191c1d] transition-colors"
                    >
                      <span className="material-symbols-outlined text-base">key</span>
                    </button>
                    <button
                      onClick={() => deleteUser(u.user_id ?? '')}
                      title="Delete"
                      className="p-1.5 rounded-lg hover:bg-red-50 text-[#6b7280] hover:text-red-500 transition-colors"
                    >
                      <span className="material-symbols-outlined text-base">delete</span>
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Create modal */}
      {showCreate && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl shadow-xl p-6 w-full max-w-md">
            <h2 className="text-lg font-bold mb-4 text-[#191c1d]">Add User</h2>
            <div className="flex flex-col gap-3">
              {(['full_name', 'username', 'email', 'phone'] as const).map((field) => (
                <div key={field}>
                  <label className="text-xs font-semibold text-[#6b7280] mb-1 block capitalize">
                    {field === 'full_name' ? 'Full Name' : field === 'username' ? 'Username' : field === 'email' ? 'Email *' : 'Phone'}
                  </label>
                  <input
                    type={field === 'email' ? 'email' : 'text'}
                    value={createForm[field]}
                    onChange={(e) => setCreateForm({ ...createForm, [field]: e.target.value })}
                    className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm"
                  />
                </div>
              ))}
              <div>
                <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Password *</label>
                <input
                  type="password"
                  value={createForm.password}
                  onChange={(e) => setCreateForm({ ...createForm, password: e.target.value })}
                  className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm"
                />
              </div>
              {createError && <p className="text-sm text-red-600">{createError}</p>}
            </div>
            <div className="flex gap-3 mt-5">
              <button onClick={submitCreate} disabled={creating}
                className="flex-1 bg-[#d4af37] text-white rounded-lg py-2 text-sm font-semibold hover:bg-[#b8962d] disabled:opacity-60">
                {creating ? 'Creating...' : 'Create'}
              </button>
              <button onClick={() => setShowCreate(false)}
                className="flex-1 border border-[#d0c5af] rounded-lg py-2 text-sm font-semibold hover:bg-[#f3f4f5]">
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Edit modal */}
      {editUser && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl shadow-xl p-6 w-full max-w-md">
            <h2 className="text-lg font-bold mb-4 text-[#191c1d]">Edit User</h2>
            <div className="flex flex-col gap-3">
              {(['full_name', 'username', 'email', 'phone'] as const).map((field) => (
                <div key={field}>
                  <label className="text-xs font-semibold text-[#6b7280] mb-1 block">
                    {field === 'full_name' ? 'Full Name' : field === 'username' ? 'Username' : field === 'email' ? 'Email' : 'Phone'}
                  </label>
                  <input
                    type={field === 'email' ? 'email' : 'text'}
                    value={editForm[field]}
                    onChange={(e) => setEditForm({ ...editForm, [field]: e.target.value })}
                    className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm"
                  />
                </div>
              ))}
              <div>
                <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Status</label>
                <select
                  value={editForm.status}
                  onChange={(e) => setEditForm({ ...editForm, status: e.target.value })}
                  className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm"
                >
                  <option value="ACTIVE">ACTIVE</option>
                  <option value="INACTIVE">INACTIVE</option>
                  <option value="SUSPENDED">SUSPENDED</option>
                </select>
              </div>
              {editError && <p className="text-sm text-red-600">{editError}</p>}
            </div>
            <div className="flex gap-3 mt-5">
              <button onClick={submitEdit} disabled={saving}
                className="flex-1 bg-[#d4af37] text-white rounded-lg py-2 text-sm font-semibold hover:bg-[#b8962d] disabled:opacity-60">
                {saving ? 'Saving...' : 'Save'}
              </button>
              <button onClick={() => setEditUser(null)}
                className="flex-1 border border-[#d0c5af] rounded-lg py-2 text-sm font-semibold hover:bg-[#f3f4f5]">
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Roles modal */}
      {rolesUser && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl shadow-xl p-6 w-full max-w-sm">
            <h2 className="text-lg font-bold mb-1 text-[#191c1d]">Manage Roles</h2>
            <p className="text-sm text-[#6b7280] mb-4">{rolesUser.full_name || rolesUser.username}</p>
            <div className="flex flex-col gap-2">
              {ALL_ROLES.map((role) => (
                <label key={role} className="flex items-center gap-3 p-3 rounded-lg border border-[#e5e7eb] cursor-pointer hover:bg-[#fafafa]">
                  <input
                    type="checkbox"
                    checked={selectedRoles.includes(role)}
                    onChange={() => toggleRole(role)}
                    className="accent-[#d4af37]"
                  />
                  <span className={`text-xs font-bold px-2 py-1 rounded-full ${ROLE_COLORS[role]}`}>{role}</span>
                </label>
              ))}
            </div>
            <div className="flex gap-3 mt-5">
              <button onClick={submitRoles} disabled={savingRoles}
                className="flex-1 bg-[#d4af37] text-white rounded-lg py-2 text-sm font-semibold hover:bg-[#b8962d] disabled:opacity-60">
                {savingRoles ? 'Saving...' : 'Save Roles'}
              </button>
              <button onClick={() => setRolesUser(null)}
                className="flex-1 border border-[#d0c5af] rounded-lg py-2 text-sm font-semibold hover:bg-[#f3f4f5]">
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Change password modal */}
      {pwUser && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl shadow-xl p-6 w-full max-w-sm">
            <h2 className="text-lg font-bold mb-1 text-[#191c1d]">Change Password</h2>
            <p className="text-sm text-[#6b7280] mb-4">{pwUser.full_name || pwUser.username}</p>
            <div className="flex flex-col gap-3">
              <div>
                <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Current Password</label>
                <input type="password" value={pwForm.oldPw}
                  onChange={(e) => setPwForm({ ...pwForm, oldPw: e.target.value })}
                  className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm" />
              </div>
              <div>
                <label className="text-xs font-semibold text-[#6b7280] mb-1 block">New Password</label>
                <input type="password" value={pwForm.newPw}
                  onChange={(e) => setPwForm({ ...pwForm, newPw: e.target.value })}
                  className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm" />
              </div>
              <div>
                <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Confirm New Password</label>
                <input type="password" value={pwForm.confirmPw}
                  onChange={(e) => setPwForm({ ...pwForm, confirmPw: e.target.value })}
                  className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm" />
              </div>
              {pwError && <p className="text-sm text-red-600">{pwError}</p>}
            </div>
            <div className="flex gap-3 mt-5">
              <button onClick={submitPw} disabled={savingPw}
                className="flex-1 bg-[#d4af37] text-white rounded-lg py-2 text-sm font-semibold hover:bg-[#b8962d] disabled:opacity-60">
                {savingPw ? 'Saving...' : 'Change Password'}
              </button>
              <button onClick={() => setPwUser(null)}
                className="flex-1 border border-[#d0c5af] rounded-lg py-2 text-sm font-semibold hover:bg-[#f3f4f5]">
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default UserManagement;
