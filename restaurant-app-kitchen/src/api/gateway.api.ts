import { useAuthStore } from '../store/authStore'

const API_BASE = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'
const WS_BASE = import.meta.env.VITE_WS_BASE_URL ?? 'ws://localhost:8080'

const req = async <T>(path: string, init?: RequestInit): Promise<T> => {
  const token = useAuthStore.getState().accessToken
  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(init?.headers ?? {}),
    },
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok || data.success === false) {
    throw new Error(data.message || data.error || `HTTP ${res.status}`)
  }
  return data as T
}

export interface OrderItemDto {
  item_id: string
  name: string
  price: number
  quantity: number
  item_status: string
}

export interface OrderDto {
  order_id: string
  table_id: string
  name: string
  phone: string
  notes: string
  time?: { seconds?: number } | string
  end_time?: { seconds?: number } | string
  party_size: number
  status: string
  total_price: number
  items: OrderItemDto[]
}

export interface KitchenNotification {
  id: string
  type: string
  target_role: string
  order_id: string
  table_id: string
  item_id: string
  item_name: string
  created_at: number
  message: string
  customer_name: string
  party_size: number
  notes: string
  items: Array<{ item_id: string; item_name: string; quantity: number }>
}

export const authApi = {
  login: (email: string, password: string) =>
    req<{ success?: boolean; access_token?: string; refresh_token?: string; user_id?: string; message?: string }>(
      '/auth/login',
      { method: 'POST', body: JSON.stringify({ email, password }) }
    ),
  logout: (refreshToken: string) =>
    req('/auth/logout', { method: 'POST', body: JSON.stringify({ refresh_token: refreshToken }) }),
}

export const usersApi = {
  getOne: (id: string) =>
    req<{ user?: { user_id: string; email: string; username: string; full_name: string; roles: string[] }; success?: boolean }>(`/users/${id}`),
}

export const ordersApi = {
  list: (query: { status?: string; page_size?: number }) => {
    const params = new URLSearchParams()
    if (query.status) params.set('status', query.status)
    params.set('page_size', String(query.page_size ?? 50))
    return req<{ orders?: OrderDto[]; total?: number }>(`/orders?${params}`)
  },
  updateItemStatus: (orderId: string, itemId: string, itemStatus: string) =>
    req<{ order?: OrderDto }>(`/orders/${orderId}/items/${itemId}/status`, {
      method: 'PATCH',
      body: JSON.stringify({ item_status: itemStatus }),
    }),
}

export const createNotificationWS = (token: string, role: string): WebSocket => {
  return new WebSocket(`${WS_BASE}/ws/notifications?token=${encodeURIComponent(token)}&role=${role}`)
}

export interface ShiftDto {
  shift_id?: string
  user_id?: string
  date?: string
  start_time?: string
  end_time?: string
  role?: string
  notes?: string
  created_by?: string
}

export const scheduleApi = {
  myShifts: (userId: string, month: string) => {
    const params = new URLSearchParams({ user_id: userId, month })
    return req<{ shifts?: ShiftDto[]; total?: number }>(`/schedule/shifts?${params}`)
  },
  create: (payload: { user_id: string; date: string; start_time: string; end_time: string; role: string; notes?: string }) =>
    req<{ shift?: ShiftDto; success?: boolean; message?: string }>('/schedule/shifts', {
      method: 'POST',
      body: JSON.stringify(payload),
    }),
  delete: (shiftId: string) =>
    req<{ success?: boolean }>(`/schedule/shifts/${shiftId}`, { method: 'DELETE' }),
}
