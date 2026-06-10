import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export interface KitchenUser {
  user_id: string
  email: string
  username: string
  full_name: string
  roles: string[]
}

interface AuthState {
  user: KitchenUser | null
  accessToken: string | null
  refreshToken: string | null
  setAuth: (user: KitchenUser, accessToken: string, refreshToken: string) => void
  clearAuth: () => void
}

export const KITCHEN_ROLES = ['CHEF', 'WAITER', 'ADMIN', 'MANAGER']

export const hasKitchenAccess = (roles: string[]) =>
  roles.some((r) => KITCHEN_ROLES.includes(r))

export const getDefaultRole = (roles: string[]): 'CHEF' | 'WAITER' | 'ADMIN' | null => {
  if (roles.includes('CHEF')) return 'CHEF'
  if (roles.includes('WAITER')) return 'WAITER'
  if (roles.includes('ADMIN') || roles.includes('MANAGER')) return 'ADMIN'
  return null
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      setAuth: (user, accessToken, refreshToken) => set({ user, accessToken, refreshToken }),
      clearAuth: () => set({ user: null, accessToken: null, refreshToken: null }),
    }),
    { name: 'kitchen-auth' }
  )
)
