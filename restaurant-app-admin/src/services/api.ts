import { useAdminAuthStore } from '../store/adminAuthStore';

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

const request = async <T>(path: string, init?: RequestInit, query?: Record<string, QueryValue>): Promise<T> => {
  const token = useAdminAuthStore.getState().accessToken;
  const response = await fetch(buildPath(path, query), {
    ...init,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(init?.headers ?? {}),
    },
  });

  const data = await response.json().catch(() => ({}));

  if (!response.ok || data.success === false) {
    throw new Error(data.message || data.error || `Request failed with status ${response.status}`);
  }

  return data as T;
};

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
}

export interface OrderDto {
  order_id?: string;
  orderId?: string;
  name?: string;
  phone?: string;
  time?: string | { seconds?: number; nanos?: number };
  date?: string;
  party_size?: number;
  partySize?: number;
  status?: string;
  total_price?: number;
  totalPrice?: number;
  items?: OrderItemDto[];
}

export interface StaffDto {
  staff_id?: string;
  staffId?: string;
  name?: string;
  role?: string;
  contact?: string;
  avatar?: string;
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
};

export const usersApi = {
  getOne: (id: string) =>
    request<{ user?: UserDto; success?: boolean; message?: string }>(`/users/${id}`),
};

export const menuApi = {
  listItems: (query?: { page?: number; page_size?: number; category_id?: string; keyword?: string }) =>
    request<{ items?: MenuItemDto[]; total?: number }>('/menu/items', undefined, query),
  createItem: (payload: { name: string; description: string; price: number; category_id?: string; category?: string; image_url: string }) =>
    request<{ item?: MenuItemDto }>('/menu/items', { method: 'POST', body: JSON.stringify(payload) }),
  updateItem: (itemId: string, payload: { name: string; description: string; price: number; category_id?: string; image_url: string }) =>
    request<{ item?: MenuItemDto }>(`/menu/items/${itemId}`, { method: 'PUT', body: JSON.stringify(payload) }),
  deleteItem: (itemId: string) => request<{ success?: boolean }>(`/menu/items/${itemId}`, { method: 'DELETE' }),
  listCategories: () => request<{ categories?: CategoryDto[]; total?: number }>('/menu/categories', undefined, { page_size: 100 }),
};

export const ordersApi = {
  list: (query?: { page?: number; page_size?: number; status?: string; keyword?: string }) =>
    request<{ orders?: OrderDto[]; total?: number }>('/orders', undefined, query),
  create: (payload: { name: string; phone: string; time: string; date: string; party_size: number; status: string; items: Array<{ item_id: string; quantity: number }> }) =>
    request<{ order?: OrderDto }>('/orders', { method: 'POST', body: JSON.stringify(payload) }),
  update: (orderId: string, payload: { name: string; phone: string; time: string; date: string; party_size: number; status: string; items: Array<{ item_id: string; quantity: number }> }) =>
    request<{ order?: OrderDto }>(`/orders/${orderId}`, { method: 'PUT', body: JSON.stringify(payload) }),
  updateStatus: (orderId: string, status: string) =>
    request<{ order?: OrderDto }>(`/orders/${orderId}/status`, { method: 'PATCH', body: JSON.stringify({ status }) }),
  cancel: (orderId: string, reason = '') =>
    request<{ success?: boolean }>(`/orders/${orderId}/cancel`, { method: 'POST', body: JSON.stringify({ reason }) }),
  delete: (orderId: string) => request<{ success?: boolean }>(`/orders/${orderId}`, { method: 'DELETE' }),
};

export const staffApi = {
  list: (query?: { page?: number; page_size?: number; keyword?: string }) =>
    request<{ staff?: StaffDto[]; total?: number }>('/staff', undefined, query),
  create: (payload: { name: string; role: string; contact: string; avatar: string }) =>
    request<{ staff?: StaffDto }>('/staff', { method: 'POST', body: JSON.stringify(payload) }),
  update: (staffId: string, payload: { name: string; role: string; contact: string; avatar: string }) =>
    request<{ staff?: StaffDto }>(`/staff/${staffId}`, { method: 'PUT', body: JSON.stringify(payload) }),
  delete: (staffId: string) => request<{ success?: boolean }>(`/staff/${staffId}`, { method: 'DELETE' }),
};
