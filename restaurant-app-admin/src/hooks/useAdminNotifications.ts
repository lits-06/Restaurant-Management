import { useEffect, useRef, useState } from 'react';
import { createNotificationWS } from '../services/api';

export interface AdminNotification {
  id: string;
  type: string;
  target_role: string;
  order_id: string;
  table_id: string;
  customer_name: string;
  party_size: number;
  notes: string;
  message: string;
  created_at: number;
}

export function useAdminNotifications(accessToken: string | null, roles: string[]) {
  const [notifications, setNotifications] = useState<AdminNotification[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [connected, setConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);

  const wsRole = roles.includes('ADMIN') ? 'ADMIN' : roles.includes('MANAGER') ? 'MANAGER' : null;

  useEffect(() => {
    if (!accessToken || !wsRole) return;

    const ws = createNotificationWS(accessToken, wsRole);
    wsRef.current = ws;

    ws.onopen = () => setConnected(true);
    ws.onclose = () => setConnected(false);
    ws.onerror = () => setConnected(false);
    ws.onmessage = (event) => {
      try {
        const notif: AdminNotification = JSON.parse(event.data);
        setNotifications((prev) => [notif, ...prev.slice(0, 49)]);
        setUnreadCount((c) => c + 1);
      } catch {
        // ignore malformed messages
      }
    };

    return () => {
      ws.close();
      wsRef.current = null;
    };
  }, [accessToken, wsRole]);

  const clearNotifications = () => {
    setNotifications([]);
    setUnreadCount(0);
  };

  const markAllRead = () => setUnreadCount(0);

  return { notifications, unreadCount, connected, clearNotifications, markAllRead };
}
