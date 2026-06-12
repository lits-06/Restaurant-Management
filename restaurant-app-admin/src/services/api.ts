import { useAdminAuthStore } from '../store/adminAuthStore';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';
const WS_BASE_URL  = import.meta.env.VITE_WS_BASE_URL  ?? 'ws://localhost:8080';

type QueryValue = string | number | boolean | undefined | null;

const buildPath = (path: string, query?: Record<string, QueryValue>) => {
  const url = new URL(path, API_BASE_URL);
  Object.entries(query ?? {}).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      url.searchParams.set(key, String(value));
    }
  });
  return url.toString();
};

// Deduplicates concurrent refresh calls into one in-flight promise
let refreshing: Promise<string | null> | null = null;

const tryRefreshToken = (): Promise<string | null> => {
  if (refreshing) return refreshing;
  const { refreshToken, user, setAuth, clearAuth } = useAdminAuthStore.getState();
  if (!refreshToken || !user) return Promise.resolve(null);

  refreshing = fetch(`${API_BASE_URL}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ refresh_token: refreshToken }),
  })
    .then(r => r.json())
    .then(data => {
      if (data.access_token) {
        setAuth(user, data.access_token, refreshToken);
        return data.access_token as string;
      }
      clearAuth();
      return null;
    })
    .catch(() => { clearAuth(); return null; })
    .finally(() => { refreshing = null; });

  return refreshing;
};

const request = async <T>(path: string, init?: RequestInit, query?: Record<string, QueryValue>): Promise<T> => {
  const token = useAdminAuthStore.getState().accessToken;

  const makeRequest = (tok: string | null) =>
    fetch(buildPath(path, query), {
      ...init,
      headers: {
        'Content-Type': 'application/json',
        ...(tok ? { Authorization: `Bearer ${tok}` } : {}),
        ...(init?.headers ?? {}),
      },
    });

  let response = await makeRequest(token);

  if (response.status === 401) {
    const newToken = await tryRefreshToken();
    if (newToken) {
      response = await makeRequest(newToken);
    }
  }

  const data = await response.json().catch(() => ({}));

  if (!response.ok || data.success === false) {
    if (response.status === 401) {
      useAdminAuthStore.getState().clearAuth();
    }
    throw new Error(data.message || data.error || `Request failed with status ${response.status}`);
  }

  return data as T;
};

// ── DTOs ──────────────────────────────────────────────────────────────────────

export interface MenuItemDto {
  item_id?: string;
  itemId?: string;
  name?: string;
  description?: string;
  price?: number;
  category?: string;
  category_id?: string;
  categoryId?: string;
  image_url?: string;
  imageUrl?: string;
}

export interface CategoryDto {
  category_id?: string;
  categoryId?: string;
  name?: string;
  description?: string;
  display_order?: number;
  displayOrder?: number;
}

export interface OrderItemDto {
  item_id?: string;
  itemId?: string;
  name?: string;
  price?: number;
  category?: string;
  image_url?: string;
  imageUrl?: string;
  quantity?: number;
  item_status?: string;
}

export interface OrderDto {
  order_id?: string;
  orderId?: string;
  table_id?: string;
  user_id?: string;
  name?: string;
  phone?: string;
  notes?: string;
  time?: string | { seconds?: number; nanos?: number };
  end_time?: string | { seconds?: number; nanos?: number };
  date?: string;
  party_size?: number;
  partySize?: number;
  status?: string;
  total?: number;
  total_price?: number;
  totalPrice?: number;
  items?: OrderItemDto[];
}

export interface TableDto {
  table_id?: string;
  table_number?: number;
  capacity?: number;
  status?: string;
}

export interface ShiftDto {
  shift_id?: string;
  user_id?: string;
  date?: string;
  start_time?: string;
  end_time?: string;
  role?: string;
  notes?: string;
  created_by?: string;
  created_at?: { seconds?: number; nanos?: number } | string;
  updated_at?: { seconds?: number; nanos?: number } | string;
}

export interface UserDto {
  user_id?: string;
  email?: string;
  username?: string;
  full_name?: string;
  phone?: string;
  status?: string;
  roles?: string[];
}

// ── API ───────────────────────────────────────────────────────────────────────

export const authApi = {
  login: (email: string, password: string) =>
    request<{ success?: boolean; access_token?: string; refresh_token?: string; user_id?: string; message?: string }>(
      '/auth/login',
      { method: 'POST', body: JSON.stringify({ email, password }) }
    ),
  logout: (refreshToken: string) =>
    request<{ success?: boolean; message?: string }>(
      '/auth/logout',
      { method: 'POST', body: JSON.stringify({ refresh_token: refreshToken }) }
    ),
  changePassword: (oldPassword: string, newPassword: string) =>
    request<{ success?: boolean; message?: string }>(
      '/auth/change-password',
      { method: 'POST', body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }) }
    ),
};

export const usersApi = {
  getOne: (id: string) =>
    request<{ user?: UserDto; success?: boolean; message?: string }>(`/users/${id}`),
  listAll: (query?: { page?: number; page_size?: number; keyword?: string }) =>
    request<{ users?: UserDto[]; total?: number }>('/users', undefined, { page_size: 200, ...query }),
  create: (payload: { email: string; password: string; username: string; full_name: string; phone: string; roles?: string[] }) =>
    request<{ user?: UserDto; success?: boolean; message?: string }>(
      '/users',
      { method: 'POST', body: JSON.stringify(payload) }
    ),
  update: (id: string, payload: { email?: string; username?: string; full_name?: string; phone?: string; status?: string }) =>
    request<{ user?: UserDto; success?: boolean; message?: string }>(
      `/users/${id}`,
      { method: 'PUT', body: JSON.stringify(payload) }
    ),
  delete: (id: string) =>
    request<{ success?: boolean }>(`/users/${id}`, { method: 'DELETE' }),
  assignRole: (id: string, roles: string[]) =>
    request<{ success?: boolean; message?: string }>(
      `/users/${id}/roles`,
      { method: 'PATCH', body: JSON.stringify({ roles }) }
    ),
  changePassword: (id: string, oldPassword: string, newPassword: string) =>
    request<{ success?: boolean; message?: string }>(
      `/users/${id}/password`,
      { method: 'PATCH', body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }) }
    ),
};

export const menuApi = {
  listItems: (query?: { page?: number; page_size?: number; category_id?: string; keyword?: string }) =>
    request<{ items?: MenuItemDto[]; total?: number }>('/menu/items', undefined, query),
  createItem: (payload: { name: string; description: string; price: number; category_id?: string; image_url: string }) =>
    request<{ item?: MenuItemDto }>('/menu/items', { method: 'POST', body: JSON.stringify(payload) }),
  updateItem: (itemId: string, payload: { name: string; description: string; price: number; category_id?: string; image_url: string }) =>
    request<{ item?: MenuItemDto }>(`/menu/items/${itemId}`, { method: 'PUT', body: JSON.stringify(payload) }),
  deleteItem: (itemId: string) =>
    request<{ success?: boolean }>(`/menu/items/${itemId}`, { method: 'DELETE' }),
  listCategories: () =>
    request<{ categories?: CategoryDto[]; total?: number }>('/menu/categories', undefined, { page_size: 100 }),
  createCategory: (payload: { name: string; description?: string; display_order?: number }) =>
    request<{ category?: CategoryDto }>('/menu/categories', { method: 'POST', body: JSON.stringify(payload) }),
  updateCategory: (id: string, payload: { name: string; description?: string; display_order?: number }) =>
    request<{ category?: CategoryDto }>(`/menu/categories/${id}`, { method: 'PUT', body: JSON.stringify(payload) }),
  deleteCategory: (id: string) =>
    request<{ success?: boolean }>(`/menu/categories/${id}`, { method: 'DELETE' }),
};

export const tablesApi = {
  list: (query?: { page?: number; page_size?: number }) =>
    request<{ tables?: TableDto[]; total?: number }>('/tables', undefined, query),
  create: (payload: { table_number: number; capacity: number }) =>
    request<{ table?: TableDto }>('/tables', { method: 'POST', body: JSON.stringify(payload) }),
  update: (id: string, payload: { table_number?: number; capacity?: number }) =>
    request<{ table?: TableDto }>(`/tables/${id}`, { method: 'PUT', body: JSON.stringify(payload) }),
  delete: (id: string) =>
    request<{ success?: boolean }>(`/tables/${id}`, { method: 'DELETE' }),
  updateStatus: (id: string, status: string) =>
    request<{ table?: TableDto }>(`/tables/${id}/status`, { method: 'PATCH', body: JSON.stringify({ status }) }),
};

export const ordersApi = {
  list: (query?: { page?: number; page_size?: number; status?: string; keyword?: string; user_id?: string; sort_order?: 'asc' | 'desc' }) =>
    request<{ orders?: OrderDto[]; total?: number }>('/orders', undefined, query),
  create: (payload: {
    name: string;
    phone: string;
    party_size: number;
    table_id?: string;
    date: string;
    time: string;
    end_time?: string;
    notes?: string;
    status?: string;
    items?: Array<{ item_id: string; quantity: number }>;
    walk_in?: boolean;
  }) =>
    request<{ success?: boolean; message?: string; order?: OrderDto }>('/orders', { method: 'POST', body: JSON.stringify(payload) }),
  update: (orderId: string, payload: {
    name: string;
    phone: string;
    notes?: string;
    date: string;
    time: string;
    end_time?: string;
    party_size: number;
    items: Array<{ item_id: string; quantity: number }>;
  }) =>
    request<{ order?: OrderDto }>(`/orders/${orderId}`, { method: 'PUT', body: JSON.stringify(payload) }),
  updateStatus: (orderId: string, status: string) =>
    request<{ order?: OrderDto }>(`/orders/${orderId}/status`, { method: 'PATCH', body: JSON.stringify({ status }) }),
  updateItemStatus: (orderId: string, itemId: string, itemStatus: string) =>
    request<{ order?: OrderDto }>(`/orders/${orderId}/items/${itemId}/status`, { method: 'PATCH', body: JSON.stringify({ item_status: itemStatus }) }),
  cancel: (orderId: string, reason = '') =>
    request<{ success?: boolean }>(`/orders/${orderId}/cancel`, { method: 'POST', body: JSON.stringify({ reason }) }),
  delete: (orderId: string) =>
    request<{ success?: boolean }>(`/orders/${orderId}`, { method: 'DELETE' }),
};

export const scheduleApi = {
  list: (query?: { month?: string; user_id?: string; role?: string }) =>
    request<{ shifts?: ShiftDto[]; total?: number }>('/schedule/shifts', undefined, query),
  create: (payload: { user_id: string; date: string; start_time: string; end_time: string; role: string; notes?: string }) =>
    request<{ shift?: ShiftDto; success?: boolean; message?: string }>('/schedule/shifts', { method: 'POST', body: JSON.stringify(payload) }),
  update: (shiftId: string, payload: { date?: string; start_time?: string; end_time?: string; notes?: string }) =>
    request<{ shift?: ShiftDto; success?: boolean; message?: string }>(`/schedule/shifts/${shiftId}`, { method: 'PUT', body: JSON.stringify(payload) }),
  delete: (shiftId: string) =>
    request<{ success?: boolean }>(`/schedule/shifts/${shiftId}`, { method: 'DELETE' }),
};

export const createNotificationWS = (token: string, role: string): WebSocket =>
  new WebSocket(`${WS_BASE_URL}/ws/notifications?token=${encodeURIComponent(token)}&role=${role}`);
