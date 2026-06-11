import { useAuthStore } from '../store/authStore'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';

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
  const { refreshToken, user, setAuth, clearAuth } = useAuthStore.getState();
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
  const token = useAuthStore.getState().accessToken;

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
      useAuthStore.getState().clearAuth();
    }
    throw new Error(data.message || data.error || `Request failed with status ${response.status}`);
  }

  return data as T;
};

// --- DTOs ---

export interface MenuItemDto {
  item_id?: string;
  itemId?: string;
  name?: string;
  description?: string;
  price?: number;
  category?: string;
  category_id?: string;
  image_url?: string;
  imageUrl?: string;
}

export interface CategoryDto {
  category_id?: string;
  categoryId?: string;
  name?: string;
}

export interface OrderItemDto {
  item_id?: string;
  name?: string;
  price?: number;
  category?: string;
  image_url?: string;
  quantity?: number;
  item_status?: string;
}

export interface OrderDto {
  order_id?: string;
  table_id?: string;
  name?: string;
  phone?: string;
  notes?: string;
  party_size?: number;
  status?: string;
  total_price?: number;
  // proto Timestamp encodes as RFC3339 string or {seconds, nanos}
  time?: string | { seconds?: string | number; nanos?: number };
  end_time?: string | { seconds?: string | number; nanos?: number };
  items?: OrderItemDto[];
}

export interface TableDto {
  table_id?: string;
  table_number?: number;
  capacity?: number;
  status?: string;
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

// --- Menu ---

export const menuApi = {
  listItems: (query?: { page?: number; page_size?: number; keyword?: string; category_id?: string }) =>
    request<{ items?: MenuItemDto[]; total?: number }>('/menu/items', undefined, query),

  listCategories: () =>
    request<{ categories?: CategoryDto[]; total?: number }>('/menu/categories', undefined, { page_size: 100 }),
};

// --- Auth ---

export const authApi = {
  login: (email: string, password: string) =>
    request<{ success?: boolean; access_token?: string; refresh_token?: string; user_id?: string; message?: string }>(
      '/auth/login',
      { method: 'POST', body: JSON.stringify({ email, password }) }
    ),

  register: (data: { email: string; password: string; username: string; full_name: string; phone: string }) =>
    request<{ success?: boolean; user_id?: string; message?: string }>(
      '/auth/register',
      { method: 'POST', body: JSON.stringify(data) }
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

// --- Users ---

export const usersApi = {
  getOne: (id: string) =>
    request<{ user?: UserDto; success?: boolean; message?: string }>(`/users/${id}`),
};

// --- Orders ---

export const ordersApi = {
  create: (payload: {
    name: string;
    phone: string;
    date: string;
    time: string;
    end_time?: string;
    party_size: number;
    notes?: string;
    table_id?: string;
    items: Array<{ item_id: string; quantity: number }>;
  }) =>
    request<{ order?: OrderDto; success?: boolean; message?: string }>(
      '/orders',
      { method: 'POST', body: JSON.stringify(payload) }
    ),

  getOne: (id: string) =>
    request<{ order?: OrderDto; success?: boolean; message?: string }>(`/orders/${id}`),

  list: (query?: { page?: number; page_size?: number; status?: string; keyword?: string; user_id?: string }) =>
    request<{ orders?: OrderDto[]; total?: number }>('/orders', undefined, query),

  cancel: (id: string) =>
    request<{ success?: boolean; message?: string }>(`/orders/${id}/cancel`, { method: 'POST' }),

  addItem: (orderId: string, item: { item_id: string; name: string; price: number; quantity: number }) =>
    request<{ order?: OrderDto; success?: boolean; message?: string }>(
      `/orders/${orderId}/items`,
      { method: 'POST', body: JSON.stringify({ item }) }
    ),
};

// --- Tables ---

export const tableApi = {
  getOne: (id: string) =>
    request<{ table?: TableDto; success?: boolean; message?: string }>(`/tables/${id}`),

  list: (query?: { page?: number; page_size?: number }) =>
    request<{ tables?: TableDto[]; total?: number }>('/tables', undefined, query),
};
