import React, { useEffect, useState, useCallback } from 'react';
import { tablesApi, type TableDto } from '../services/api';

const STATUS_LABELS: Record<string, string> = {
  AVAILABLE:      'Available',
  CLEANING:       'Cleaning',
  OUT_OF_SERVICE: 'Out of Service',
};

const STATUS_COLORS: Record<string, string> = {
  AVAILABLE:      'bg-green-100 text-green-700',
  CLEANING:       'bg-amber-100 text-amber-700',
  OUT_OF_SERVICE: 'bg-red-100 text-red-600',
};

const ALL_STATUSES = ['AVAILABLE', 'CLEANING', 'OUT_OF_SERVICE'];

interface FormState {
  table_number: string;
  capacity: string;
}

const emptyForm = (): FormState => ({ table_number: '', capacity: '' });

const TableManagement: React.FC = () => {
  const [tables, setTables] = useState<TableDto[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  // Create modal
  const [showCreate, setShowCreate] = useState(false);
  const [createForm, setCreateForm] = useState<FormState>(emptyForm());
  const [creating, setCreating] = useState(false);
  const [createError, setCreateError] = useState('');

  // Edit modal
  const [editTable, setEditTable] = useState<TableDto | null>(null);
  const [editForm, setEditForm] = useState<FormState>(emptyForm());
  const [saving, setSaving] = useState(false);
  const [editError, setEditError] = useState('');

  // Status modal
  const [statusTable, setStatusTable] = useState<TableDto | null>(null);
  const [newStatus, setNewStatus] = useState('');
  const [updatingStatus, setUpdatingStatus] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const res = await tablesApi.list({ page_size: 200 });
      const sorted = (res.tables ?? []).slice().sort(
        (a, b) => (a.table_number ?? 0) - (b.table_number ?? 0)
      );
      setTables(sorted);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load data');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { load(); }, [load]);

  // ── Create ──────────────────────────────────────────────────────
  const submitCreate = async () => {
    const num = parseInt(createForm.table_number);
    const cap = parseInt(createForm.capacity);
    if (!num || num < 1) { setCreateError('Invalid table number'); return; }
    if (!cap || cap < 1 || cap > 50) { setCreateError('Capacity must be 1–50'); return; }
    setCreating(true);
    setCreateError('');
    try {
      await tablesApi.create({ table_number: num, capacity: cap });
      setShowCreate(false);
      setCreateForm(emptyForm());
      load();
    } catch (e) {
      setCreateError(e instanceof Error ? e.message : 'Failed to create table');
    } finally {
      setCreating(false);
    }
  };

  // ── Edit ─────────────────────────────────────────────────────────
  const openEdit = (t: TableDto) => {
    setEditTable(t);
    setEditForm({
      table_number: String(t.table_number ?? ''),
      capacity:     String(t.capacity ?? ''),
    });
    setEditError('');
  };

  const submitEdit = async () => {
    if (!editTable?.table_id) return;
    const num = parseInt(editForm.table_number);
    const cap = parseInt(editForm.capacity);
    if (!num || num < 1) { setEditError('Invalid table number'); return; }
    if (!cap || cap < 1 || cap > 50) { setEditError('Capacity must be 1–50'); return; }
    setSaving(true);
    setEditError('');
    try {
      await tablesApi.update(editTable.table_id, { table_number: num, capacity: cap });
      setEditTable(null);
      load();
    } catch (e) {
      setEditError(e instanceof Error ? e.message : 'Failed to save');
    } finally {
      setSaving(false);
    }
  };

  // ── Delete ───────────────────────────────────────────────────────
  const deleteTable = async (id: string) => {
    if (!confirm('Delete this table? This action cannot be undone.')) return;
    try {
      await tablesApi.delete(id);
      load();
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Failed to delete');
    }
  };

  // ── Status ───────────────────────────────────────────────────────
  const openStatusModal = (t: TableDto) => {
    setStatusTable(t);
    setNewStatus(t.status ?? 'AVAILABLE');
  };

  const submitStatus = async () => {
    if (!statusTable?.table_id) return;
    setUpdatingStatus(true);
    try {
      await tablesApi.updateStatus(statusTable.table_id, newStatus);
      setStatusTable(null);
      load();
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Failed to update status');
    } finally {
      setUpdatingStatus(false);
    }
  };

  return (
    <div className="flex flex-col gap-6">
      {/* Header */}
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div>
          <h1 className="text-2xl font-bold text-[#191c1d]">Table Management</h1>
          <p className="text-sm text-[#6b7280] mt-1">{tables.length} tables · Physical status</p>
        </div>
        <button
          onClick={() => { setShowCreate(true); setCreateForm(emptyForm()); setCreateError(''); }}
          className="flex items-center gap-2 bg-[#d4af37] text-white px-4 py-2 rounded-lg text-sm font-semibold hover:bg-[#b8962d] transition-colors"
        >
          <span className="material-symbols-outlined text-base">add</span>
          Add Table
        </button>
      </div>

      {error && <div className="text-red-600 text-sm bg-red-50 p-3 rounded-lg">{error}</div>}

      {/* Table grid */}
      {loading ? (
        <div className="text-center py-16 text-[#6b7280]">Loading...</div>
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
          {tables.map((t) => (
            <div
              key={t.table_id}
              className="bg-white rounded-xl border border-[#e5e7eb] p-4 shadow-sm flex flex-col gap-3"
            >
              <div className="flex items-start justify-between">
                <div>
                  <p className="font-bold text-xl text-[#191c1d]">Table {t.table_number}</p>
                  <p className="text-xs text-[#6b7280]">{t.capacity} seats</p>
                </div>
                <span className={`text-[10px] font-bold px-2 py-1 rounded-full ${STATUS_COLORS[t.status ?? ''] ?? 'bg-gray-100 text-gray-600'}`}>
                  {STATUS_LABELS[t.status ?? ''] ?? t.status}
                </span>
              </div>
              <div className="flex gap-1 mt-auto">
                <button
                  onClick={() => openEdit(t)}
                  className="flex-1 text-xs py-1.5 rounded-lg border border-[#d0c5af] text-[#4d4635] hover:bg-[#f3f4f5] transition-colors"
                >
                  Edit
                </button>
                <button
                  onClick={() => openStatusModal(t)}
                  className="flex-1 text-xs py-1.5 rounded-lg border border-[#d4af37] text-[#735c00] hover:bg-[#d4af37]/10 transition-colors"
                >
                  Status
                </button>
                <button
                  onClick={() => deleteTable(t.table_id ?? '')}
                  className="text-xs px-2 py-1.5 rounded-lg border border-red-200 text-red-500 hover:bg-red-50 transition-colors"
                >
                  <span className="material-symbols-outlined text-sm leading-none">delete</span>
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {tables.length === 0 && !loading && (
        <div className="text-center py-16 text-[#6b7280]">No tables yet. Add your first table.</div>
      )}

      {/* Create modal */}
      {showCreate && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl shadow-xl p-6 w-full max-w-sm">
            <h2 className="text-lg font-bold mb-4 text-[#191c1d]">Add New Table</h2>
            <div className="flex flex-col gap-3">
              <div>
                <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Table Number</label>
                <input
                  type="number" min="1"
                  value={createForm.table_number}
                  onChange={(e) => setCreateForm({ ...createForm, table_number: e.target.value })}
                  placeholder="e.g. 5"
                  className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm"
                />
              </div>
              <div>
                <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Capacity (1–50)</label>
                <input
                  type="number" min="1" max="50"
                  value={createForm.capacity}
                  onChange={(e) => setCreateForm({ ...createForm, capacity: e.target.value })}
                  placeholder="e.g. 4"
                  className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm"
                />
              </div>
              {createError && <p className="text-sm text-red-600">{createError}</p>}
            </div>
            <div className="flex gap-3 mt-5">
              <button
                onClick={submitCreate}
                disabled={creating}
                className="flex-1 bg-[#d4af37] text-white rounded-lg py-2 text-sm font-semibold hover:bg-[#b8962d] disabled:opacity-60"
              >
                {creating ? 'Creating...' : 'Create'}
              </button>
              <button
                onClick={() => setShowCreate(false)}
                className="flex-1 border border-[#d0c5af] rounded-lg py-2 text-sm font-semibold hover:bg-[#f3f4f5]"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Edit modal */}
      {editTable && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl shadow-xl p-6 w-full max-w-sm">
            <h2 className="text-lg font-bold mb-4 text-[#191c1d]">Edit Table {editTable.table_number}</h2>
            <div className="flex flex-col gap-3">
              <div>
                <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Table Number</label>
                <input
                  type="number" min="1"
                  value={editForm.table_number}
                  onChange={(e) => setEditForm({ ...editForm, table_number: e.target.value })}
                  className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm"
                />
              </div>
              <div>
                <label className="text-xs font-semibold text-[#6b7280] mb-1 block">Capacity (1–50)</label>
                <input
                  type="number" min="1" max="50"
                  value={editForm.capacity}
                  onChange={(e) => setEditForm({ ...editForm, capacity: e.target.value })}
                  className="w-full border border-[#d0c5af] rounded-lg px-3 py-2 text-sm"
                />
              </div>
              {editError && <p className="text-sm text-red-600">{editError}</p>}
            </div>
            <div className="flex gap-3 mt-5">
              <button
                onClick={submitEdit}
                disabled={saving}
                className="flex-1 bg-[#d4af37] text-white rounded-lg py-2 text-sm font-semibold hover:bg-[#b8962d] disabled:opacity-60"
              >
                {saving ? 'Saving...' : 'Save'}
              </button>
              <button
                onClick={() => setEditTable(null)}
                className="flex-1 border border-[#d0c5af] rounded-lg py-2 text-sm font-semibold hover:bg-[#f3f4f5]"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Status modal */}
      {statusTable && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
          <div className="bg-white rounded-2xl shadow-xl p-6 w-full max-w-sm">
            <h2 className="text-lg font-bold mb-1 text-[#191c1d]">Update Status</h2>
            <p className="text-sm text-[#6b7280] mb-4">Table {statusTable.table_number}</p>
            <div className="flex flex-col gap-2">
              {ALL_STATUSES.map((s) => (
                <label key={s} className="flex items-center gap-3 p-3 rounded-lg border border-[#e5e7eb] cursor-pointer hover:bg-[#fafafa]">
                  <input
                    type="radio"
                    name="status"
                    value={s}
                    checked={newStatus === s}
                    onChange={() => setNewStatus(s)}
                    className="accent-[#d4af37]"
                  />
                  <span className={`text-xs font-bold px-2 py-1 rounded-full ${STATUS_COLORS[s]}`}>
                    {STATUS_LABELS[s]}
                  </span>
                </label>
              ))}
            </div>
            <div className="flex gap-3 mt-5">
              <button
                onClick={submitStatus}
                disabled={updatingStatus}
                className="flex-1 bg-[#d4af37] text-white rounded-lg py-2 text-sm font-semibold hover:bg-[#b8962d] disabled:opacity-60"
              >
                {updatingStatus ? 'Saving...' : 'Update'}
              </button>
              <button
                onClick={() => setStatusTable(null)}
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

export default TableManagement;
