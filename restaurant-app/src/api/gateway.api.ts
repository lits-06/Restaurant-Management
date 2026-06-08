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
  const response = await fetch(buildPath(path, query), {
    ...init,
    headers: {
      'Content-Type': 'application/json',
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
  image_url?: string;
  imageUrl?: string;
}

export interface OrderDto {
  order_id?: string;
  orderId?: string;
}

export const menuApi = {
  listItems: (query?: { page?: number; page_size?: number; keyword?: string; category_id?: string }) =>
    request<{ items?: MenuItemDto[]; total?: number }>('/menu/items', undefined, query),
};

export const ordersApi = {
  create: (payload: {
    name: string;
    phone: string;
    time: string;
    date: string;
    party_size: number;
    status: string;
    items: Array<{ item_id: string; quantity: number }>;
  }) => request<{ order?: OrderDto; message?: string }>('/orders', { method: 'POST', body: JSON.stringify(payload) }),
};
