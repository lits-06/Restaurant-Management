import { useState } from 'react'
import { authApi, usersApi } from '../api/gateway.api'
import { hasKitchenAccess, useAuthStore, type KitchenUser } from '../store/authStore'

interface Props {
  onSuccess: () => void
}

export default function LoginPage({ onSuccess }: Props) {
  const setAuth = useAuthStore((s) => s.setAuth)
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const loginRes = await authApi.login(email, password)
      if (!loginRes.access_token || !loginRes.user_id) {
        throw new Error(loginRes.message ?? 'Login failed')
      }
      const userRes = await usersApi.getOne(loginRes.user_id)
      const user = userRes.user
      if (!user || !hasKitchenAccess(user.roles ?? [])) {
        throw new Error('Tài khoản không có quyền truy cập Kitchen Display')
      }
      const kitchenUser: KitchenUser = {
        user_id: user.user_id ?? '',
        email: user.email ?? '',
        username: user.username ?? '',
        full_name: user.full_name ?? '',
        roles: user.roles ?? [],
      }
      setAuth(kitchenUser, loginRes.access_token, loginRes.refresh_token ?? '')
      onSuccess()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Đăng nhập thất bại')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-900 flex items-center justify-center p-4">
      <div className="bg-gray-800 rounded-2xl p-8 w-full max-w-sm shadow-xl">
        <div className="text-center mb-8">
          <div className="text-4xl mb-2">🍳</div>
          <h1 className="text-2xl font-bold text-white">Kitchen Display</h1>
          <p className="text-gray-400 text-sm mt-1">Dành cho CHEF và WAITER</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm text-gray-400 mb-1">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className="w-full bg-gray-700 text-white rounded-lg px-4 py-2.5 outline-none focus:ring-2 focus:ring-orange-500"
              placeholder="staff@restaurant.com"
            />
          </div>
          <div>
            <label className="block text-sm text-gray-400 mb-1">Mật khẩu</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="w-full bg-gray-700 text-white rounded-lg px-4 py-2.5 outline-none focus:ring-2 focus:ring-orange-500"
              placeholder="••••••••"
            />
          </div>

          {error && (
            <p className="text-red-400 text-sm bg-red-900/30 rounded-lg px-3 py-2">{error}</p>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-orange-500 hover:bg-orange-600 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg transition-colors"
          >
            {loading ? 'Đang đăng nhập...' : 'Đăng nhập'}
          </button>
        </form>
      </div>
    </div>
  )
}
