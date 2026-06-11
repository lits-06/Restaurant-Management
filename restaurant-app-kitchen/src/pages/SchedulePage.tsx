import { useCallback, useEffect, useState } from 'react'
import { scheduleApi, type ShiftDto } from '../api/gateway.api'
import { useAuthStore } from '../store/authStore'

function isoMonth(d: Date) {
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`
}

function addMonths(base: Date, delta: number) {
  const d = new Date(base)
  d.setMonth(d.getMonth() + delta)
  d.setDate(1)
  return d
}

interface CreateFormState {
  date: string
  startTime: string
  endTime: string
  notes: string
}

interface Props {
  onLogout: () => void
}

export default function SchedulePage({ onLogout }: Props) {
  const { user } = useAuthStore()
  const [baseDate, setBaseDate] = useState(new Date())
  const [shifts, setShifts] = useState<ShiftDto[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [form, setForm] = useState<CreateFormState | null>(null)
  const [saving, setSaving] = useState(false)

  const month = isoMonth(baseDate)
  const monthLabel = baseDate.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })

  const load = useCallback(async () => {
    if (!user) return
    setLoading(true)
    setError('')
    try {
      const res = await scheduleApi.myShifts(user.user_id, month)
      setShifts(res.shifts ?? [])
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load shifts')
    } finally {
      setLoading(false)
    }
  }, [user, month])

  useEffect(() => { load() }, [load])

  const today = new Date().toISOString().slice(0, 10)

  const openForm = () =>
    setForm({ date: today, startTime: '08:00', endTime: '16:00', notes: '' })

  const submitCreate = async () => {
    if (!form || !user) return
    setSaving(true)
    try {
      // role = first staff role of the user
      const role = user.roles.find((r) => ['CHEF', 'WAITER', 'MANAGER', 'ADMIN'].includes(r)) ?? 'CHEF'
      await scheduleApi.create({
        user_id:    user.user_id,
        date:       form.date,
        start_time: form.startTime,
        end_time:   form.endTime,
        role,
        notes: form.notes,
      })
      setForm(null)
      load()
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Failed to register shift')
    } finally {
      setSaving(false)
    }
  }

  const deleteShift = async (id: string) => {
    if (!confirm('Delete this shift?')) return
    try {
      await scheduleApi.delete(id)
      load()
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Failed to delete shift')
    }
  }

  return (
    <div className="min-h-screen bg-gray-900 text-gray-100">
      {/* Header */}
      <div className="bg-gray-800 border-b border-gray-700 px-6 py-4 flex items-center justify-between">
        <div>
          <h1 className="text-xl font-bold text-white">My Schedule</h1>
          <p className="text-sm text-gray-400 capitalize">{monthLabel}</p>
        </div>
        <div className="flex items-center gap-3">
          <button onClick={() => setBaseDate(addMonths(baseDate, -1))} className="px-3 py-1.5 bg-gray-700 rounded-lg text-sm hover:bg-gray-600">◀</button>
          <button onClick={() => setBaseDate(addMonths(baseDate, 1))} className="px-3 py-1.5 bg-gray-700 rounded-lg text-sm hover:bg-gray-600">▶</button>
          <button
            onClick={openForm}
            className="px-4 py-1.5 bg-orange-600 text-white rounded-lg text-sm font-semibold hover:bg-orange-500"
          >
            + Register Shift
          </button>
          <button onClick={onLogout} className="px-3 py-1.5 bg-gray-700 rounded-lg text-sm hover:bg-gray-600 text-gray-300">Sign Out</button>
        </div>
      </div>

      <div className="p-6 max-w-2xl mx-auto">
        {error && <div className="bg-red-900/50 text-red-300 rounded-lg p-3 mb-4 text-sm">{error}</div>}

        {loading ? (
          <div className="text-center py-16 text-gray-500">Loading...</div>
        ) : shifts.length === 0 ? (
          <div className="text-center py-16 text-gray-500">
            <p className="text-lg mb-2">No shifts this month</p>
            <p className="text-sm">Press "+ Register Shift" to create a new shift</p>
          </div>
        ) : (
          <div className="flex flex-col gap-3">
            {shifts.map((s) => (
              <div key={s.shift_id} className="bg-gray-800 rounded-xl p-4 flex items-center justify-between border border-gray-700">
                <div>
                  <p className="font-semibold text-white">{s.date}</p>
                  <p className="text-sm text-gray-400">{s.start_time} – {s.end_time}</p>
                  {s.notes && <p className="text-xs text-gray-500 mt-1">{s.notes}</p>}
                </div>
                <div className="flex items-center gap-3">
                  <span className="text-xs font-semibold px-2 py-1 rounded-full bg-orange-900/50 text-orange-300">{s.role}</span>
                  <button
                    onClick={() => deleteShift(s.shift_id ?? '')}
                    className="text-red-400 hover:text-red-300 text-sm font-medium"
                  >
                    Delete
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Create form modal */}
      {form && (
        <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50">
          <div className="bg-gray-800 rounded-2xl shadow-2xl p-6 w-full max-w-sm border border-gray-700">
            <h2 className="text-lg font-bold text-white mb-4">Register New Shift</h2>
            <div className="flex flex-col gap-3">
              <div>
                <label className="text-xs font-semibold text-gray-400 mb-1 block">Date</label>
                <input
                  type="date"
                  value={form.date}
                  min={today}
                  onChange={(e) => setForm({ ...form, date: e.target.value })}
                  className="w-full bg-gray-700 border border-gray-600 rounded-lg px-3 py-2 text-sm text-white"
                />
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="text-xs font-semibold text-gray-400 mb-1 block">Start Time</label>
                  <input type="time" value={form.startTime} onChange={(e) => setForm({ ...form, startTime: e.target.value })}
                    className="w-full bg-gray-700 border border-gray-600 rounded-lg px-3 py-2 text-sm text-white" />
                </div>
                <div>
                  <label className="text-xs font-semibold text-gray-400 mb-1 block">End Time</label>
                  <input type="time" value={form.endTime} onChange={(e) => setForm({ ...form, endTime: e.target.value })}
                    className="w-full bg-gray-700 border border-gray-600 rounded-lg px-3 py-2 text-sm text-white" />
                </div>
              </div>
              <div>
                <label className="text-xs font-semibold text-gray-400 mb-1 block">Notes</label>
                <input type="text" value={form.notes} onChange={(e) => setForm({ ...form, notes: e.target.value })}
                  placeholder="e.g. Morning shift..." className="w-full bg-gray-700 border border-gray-600 rounded-lg px-3 py-2 text-sm text-white" />
              </div>
            </div>
            <div className="flex gap-3 mt-5">
              <button onClick={submitCreate} disabled={saving}
                className="flex-1 bg-orange-600 text-white rounded-lg py-2 text-sm font-semibold hover:bg-orange-500 disabled:opacity-60">
                {saving ? 'Saving...' : 'Register'}
              </button>
              <button onClick={() => setForm(null)} className="flex-1 bg-gray-700 text-gray-200 rounded-lg py-2 text-sm font-semibold hover:bg-gray-600">
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
