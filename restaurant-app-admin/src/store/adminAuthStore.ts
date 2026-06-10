import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export interface AdminUser {
  user_id: string
  email: string
  username: string
  full_name: string
  roles: string[]
}

interface AdminAuthState {
  user: AdminUser | null
  accessToken: string | null
  refreshToken: string | null
  setAuth: (user: AdminUser, accessToken: string, refreshToken: string) => void
  clearAuth: () => void
}

const STAFF_ROLES = ['ADMIN', 'MANAGER', 'CHEF', 'WAITER']

export const hasAdminAccess = (roles: string[]) =>
  roles.some((r) => STAFF_ROLES.includes(r))

export const useAdminAuthStore = create<AdminAuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      setAuth: (user, accessToken, refreshToken) => set({ user, accessToken, refreshToken }),
      clearAuth: () => set({ user: null, accessToken: null, refreshToken: null }),
    }),
    { name: 'luxe-admin-auth' }
  )
)
