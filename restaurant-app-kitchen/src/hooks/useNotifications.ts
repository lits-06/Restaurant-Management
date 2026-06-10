import { useEffect, useRef, useState } from 'react'
import { createNotificationWS, type KitchenNotification } from '../api/gateway.api'

export function useNotifications(accessToken: string | null, role: string | null) {
  const [notifications, setNotifications] = useState<KitchenNotification[]>([])
  const [connected, setConnected] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)

  useEffect(() => {
    if (!accessToken || !role) return

    const ws = createNotificationWS(accessToken, role)
    wsRef.current = ws

    ws.onopen = () => setConnected(true)
    ws.onclose = () => setConnected(false)
    ws.onerror = () => setConnected(false)
    ws.onmessage = (event) => {
      try {
        const notif: KitchenNotification = JSON.parse(event.data)
        setNotifications((prev) => [notif, ...prev.slice(0, 49)])
      } catch {
        // ignore malformed messages
      }
    }

    return () => {
      ws.close()
      wsRef.current = null
    }
  }, [accessToken, role])

  const clearNotifications = () => setNotifications([])

  return { notifications, connected, clearNotifications }
}
