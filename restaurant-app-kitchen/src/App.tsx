import { useState, useEffect } from 'react'
import { authApi } from './api/gateway.api'
import { getDefaultRole, useAuthStore } from './store/authStore'
import LoginPage from './pages/LoginPage'
import KitchenPage from './pages/KitchenPage'
import WaiterPage from './pages/WaiterPage'
import SchedulePage from './pages/SchedulePage'

type ActiveView = 'CHEF' | 'WAITER' | 'SCHEDULE'

export default function App() {
  const { user, refreshToken, clearAuth } = useAuthStore()
  const [view, setView] = useState<ActiveView | null>(() =>
    sessionStorage.getItem('kitchen-view') as ActiveView | null
  )

  const activeView = view ?? (getDefaultRole(user?.roles ?? []) === 'WAITER' ? 'WAITER' : 'CHEF')

  useEffect(() => {
    sessionStorage.setItem('kitchen-view', activeView)
  }, [activeView])

  const handleLoginSuccess = () => {
    const s = useAuthStore.getState()
    const role = s.user ? getDefaultRole(s.user.roles) : null
    const newView: ActiveView = role === 'WAITER' ? 'WAITER' : 'CHEF'
    sessionStorage.setItem('kitchen-view', newView)
    setView(newView)
  }

  const handleLogout = async () => {
    try {
      if (refreshToken) await authApi.logout(refreshToken)
    } catch {}
    sessionStorage.removeItem('kitchen-view')
    clearAuth()
    setView(null)
  }

  if (!user) {
    return <LoginPage onSuccess={handleLoginSuccess} />
  }

  // ADMIN/MANAGER can switch views
  const canSwitchView =
    user.roles.includes('ADMIN') || user.roles.includes('MANAGER')

  return (
    <div className="relative">
      {/* Floating tab switcher — always visible for ADMIN/MANAGER; CHEF/WAITER can access Schedule */}
      <div className="fixed bottom-4 right-4 z-50 flex gap-2 bg-gray-800 rounded-full p-1 shadow-xl border border-gray-700">
        {(canSwitchView || user.roles.includes('CHEF')) && (
          <button
            onClick={() => setView('CHEF')}
            className={`px-3 py-1.5 rounded-full text-sm font-medium transition-colors ${
              activeView === 'CHEF' ? 'bg-orange-600 text-white' : 'text-gray-400 hover:text-white'
            }`}
          >
            🍳 Kitchen
          </button>
        )}
        {(canSwitchView || user.roles.includes('WAITER')) && activeView !== 'WAITER' && (
          <button
            onClick={() => setView('WAITER')}
            className={`px-3 py-1.5 rounded-full text-sm font-medium transition-colors ${
              activeView === 'WAITER' ? 'bg-green-600 text-white' : 'text-gray-400 hover:text-white'
            }`}
          >
            🛎 Service
          </button>
        )}
        <button
          onClick={() => setView('SCHEDULE')}
          className={`px-3 py-1.5 rounded-full text-sm font-medium transition-colors ${
            activeView === 'SCHEDULE' ? 'bg-blue-600 text-white' : 'text-gray-400 hover:text-white'
          }`}
        >
          📅 Schedule
        </button>
      </div>

      {activeView === 'CHEF' ? (
        <KitchenPage onLogout={handleLogout} />
      ) : activeView === 'WAITER' ? (
        <WaiterPage onLogout={handleLogout} />
      ) : (
        <SchedulePage onLogout={handleLogout} />
      )}
    </div>
  )
}
