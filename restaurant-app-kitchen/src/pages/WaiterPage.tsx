import { useCallback, useEffect, useState } from 'react'
import { ordersApi, type KitchenNotification, type OrderDto } from '../api/gateway.api'
import { useNotifications } from '../hooks/useNotifications'
import { useAuthStore } from '../store/authStore'

interface Props {
  onLogout: () => void
}

export default function WaiterPage({ onLogout }: Props) {
  const { user, accessToken } = useAuthStore()
  const [orders, setOrders] = useState<OrderDto[]>([])
  const [loading, setLoading] = useState(true)
  const { notifications, connected, clearNotifications } = useNotifications(accessToken, 'WAITER')

  const fetchOrders = useCallback(async () => {
    try {
      const res = await ordersApi.list({ status: 'Confirmed', page_size: 50 })
      setOrders(res.orders ?? [])
    } catch {
      // silently ignore
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { fetchOrders() }, [fetchOrders])

  // When ITEM_READY arrives, refresh orders to get updated statuses
  useEffect(() => {
    if (notifications.length === 0) return
    if (notifications[0].type === 'ITEM_READY') fetchOrders()
  }, [notifications, fetchOrders])

  const markServed = async (orderId: string, itemId: string) => {
    try {
      const res = await ordersApi.updateItemStatus(orderId, itemId, 'SERVED')
      if (res.order) {
        setOrders((prev) => prev.map((o) => (o.order_id === orderId ? res.order! : o)))
      }
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Update failed')
    }
  }

  const readyItems = orders.flatMap((o) =>
    (o.items ?? [])
      .filter((i) => i.item_status === 'READY')
      .map((i) => ({ ...i, order: o }))
  )

  return (
    <div className="min-h-screen bg-gray-900 text-white flex flex-col">
      {/* Header */}
      <header className="bg-gray-800 border-b border-gray-700 px-6 py-3 flex items-center justify-between shrink-0">
        <div className="flex items-center gap-3">
          <span className="text-2xl">🛎</span>
          <div>
            <h1 className="font-bold text-lg leading-tight">Service</h1>
            <p className="text-xs text-gray-400">{user?.full_name}</p>
          </div>
        </div>
        <div className="flex items-center gap-4">
          <span className={`text-xs px-2 py-1 rounded-full ${connected ? 'bg-green-800 text-green-300' : 'bg-red-900 text-red-300'}`}>
            {connected ? '● Live' : '○ Offline'}
          </span>
          {readyItems.length > 0 && (
            <span className="bg-green-700 text-white text-xs font-bold px-2 py-1 rounded-full animate-pulse">
              {readyItems.length} ready
            </span>
          )}
          <button onClick={onLogout} className="text-sm text-gray-400 hover:text-white transition-colors">
            Sign Out
          </button>
        </div>
      </header>

      <div className="flex flex-1 overflow-hidden">
        {/* Main: READY items */}
        <main className="flex-1 overflow-y-auto p-6">
          <h2 className="text-lg font-semibold text-green-400 mb-4">
            Items to Serve ({readyItems.length})
          </h2>

          {loading ? (
            <p className="text-gray-400 text-center py-12">Loading...</p>
          ) : readyItems.length === 0 ? (
            <div className="text-center py-16 text-gray-500">
              <div className="text-5xl mb-4">☕</div>
              <p className="text-lg">No items ready yet</p>
            </div>
          ) : (
            <div className="space-y-3">
              {readyItems.map(({ order, ...item }) => (
                <div
                  key={`${order.order_id}-${item.item_id}`}
                  className="bg-gray-800 border border-green-700 rounded-xl px-5 py-4 flex items-center gap-4"
                >
                  <div className="w-2 h-10 bg-green-500 rounded-full shrink-0" />
                  <div className="flex-1">
                    <p className="font-semibold text-white text-base">
                      {item.name} <span className="text-gray-400 font-normal">×{item.quantity}</span>
                    </p>
                    <p className="text-sm text-gray-400 mt-0.5">
                      {order.name} · {order.party_size} guests
                      <span className="text-orange-400 ml-2 font-mono text-xs">
                        Table: {order.table_id.slice(0, 8)}…
                      </span>
                    </p>
                  </div>
                  <button
                    onClick={() => markServed(order.order_id, item.item_id)}
                    className="bg-green-600 hover:bg-green-500 text-white font-semibold px-4 py-2 rounded-lg transition-colors shrink-0"
                  >
                    Served ✓
                  </button>
                </div>
              ))}
            </div>
          )}

          {/* All orders summary */}
          <div className="mt-8">
            <h2 className="text-sm font-semibold text-gray-400 mb-3 uppercase tracking-wide">
              All Active Orders ({orders.length})
            </h2>
            <div className="space-y-2">
              {orders.map((order) => {
                const r = order.items?.filter((i) => i.item_status === 'READY').length ?? 0
                const s = order.items?.filter((i) => i.item_status === 'SERVED').length ?? 0
                const total = order.items?.length ?? 0
                return (
                  <div key={order.order_id} className="bg-gray-800 rounded-lg px-4 py-2.5 flex items-center justify-between">
                    <div>
                      <span className="font-medium text-white text-sm">{order.name}</span>
                      <span className="text-gray-500 text-xs ml-2">{order.party_size} guests</span>
                    </div>
                    <div className="flex items-center gap-2 text-xs">
                      {r > 0 && <span className="bg-green-800 text-green-300 px-2 py-0.5 rounded">{r} ready</span>}
                      <span className="text-gray-500">{s}/{total} served</span>
                    </div>
                  </div>
                )
              })}
            </div>
          </div>
        </main>

        {/* Sidebar: notification feed */}
        <aside className="w-72 bg-gray-800 border-l border-gray-700 flex flex-col shrink-0">
          <div className="px-4 py-3 border-b border-gray-700 flex items-center justify-between">
            <h3 className="text-sm font-semibold text-gray-300">Notifications</h3>
            {notifications.length > 0 && (
              <button onClick={clearNotifications} className="text-xs text-gray-500 hover:text-gray-300">
                Clear
              </button>
            )}
          </div>
          <div className="flex-1 overflow-y-auto p-3 space-y-2">
            {notifications.length === 0 ? (
              <p className="text-gray-500 text-xs text-center py-6">No notifications</p>
            ) : (
              notifications.map((n) => <NotificationItem key={n.id} notif={n} />)
            )}
          </div>
        </aside>
      </div>
    </div>
  )
}

function NotificationItem({ notif }: { notif: KitchenNotification }) {
  const time = notif.created_at
    ? new Date(notif.created_at * 1000).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
    : ''

  return (
    <div className="bg-gray-700 rounded-lg px-3 py-2.5 border-l-4 border-green-500">
      <p className="text-sm text-white font-medium leading-snug">{notif.message || notif.item_name}</p>
      {notif.table_id && (
        <p className="text-xs text-gray-400 mt-0.5">
          Table: <span className="font-mono">{notif.table_id.slice(0, 8)}…</span>
        </p>
      )}
      <p className="text-xs text-gray-500 mt-1">{time}</p>
    </div>
  )
}
