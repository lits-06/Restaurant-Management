import { useCallback, useEffect, useState } from 'react'
import { ordersApi, type OrderDto, type OrderItemDto } from '../api/gateway.api'
import { useNotifications } from '../hooks/useNotifications'
import { useAuthStore } from '../store/authStore'

const STATUS_STYLE: Record<string, string> = {
  PENDING:  'bg-gray-700 text-gray-300',
  COOKING:  'bg-amber-600 text-white',
  READY:    'bg-green-600 text-white',
  SERVED:   'bg-gray-600 text-gray-400 line-through',
}

const STATUS_LABEL: Record<string, string> = {
  PENDING: 'Chờ',
  COOKING: 'Đang nấu',
  READY:   'Xong',
  SERVED:  'Đã mang',
}

interface Props {
  onLogout: () => void
}

export default function KitchenPage({ onLogout }: Props) {
  const { user, accessToken } = useAuthStore()
  const [orders, setOrders] = useState<OrderDto[]>([])
  const [loading, setLoading] = useState(true)
  const { notifications, connected } = useNotifications(accessToken, 'CHEF')

  const fetchOrders = useCallback(async () => {
    try {
      const res = await ordersApi.list({ status: 'Confirmed', page_size: 50 })
      setOrders(res.orders ?? [])
    } catch {
      // silently retry on next notification
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { fetchOrders() }, [fetchOrders])

  // Refresh when a new ORDER_CONFIRMED arrives
  useEffect(() => {
    if (notifications.length === 0) return
    const latest = notifications[0]
    if (latest.type === 'ORDER_CONFIRMED') fetchOrders()
  }, [notifications, fetchOrders])

  const markItem = async (orderId: string, itemId: string, next: string) => {
    try {
      const res = await ordersApi.updateItemStatus(orderId, itemId, next)
      if (res.order) {
        setOrders((prev) => prev.map((o) => (o.order_id === orderId ? res.order! : o)))
      }
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Lỗi cập nhật')
    }
  }

  const nextAction = (item: OrderItemDto): { label: string; next: string } | null => {
    if (item.item_status === 'PENDING')  return { label: 'Bắt đầu nấu', next: 'COOKING' }
    if (item.item_status === 'COOKING')  return { label: 'Đã xong ✓', next: 'READY' }
    return null
  }

  const activeOrders = orders.filter((o) =>
    o.items?.some((i) => i.item_status !== 'SERVED')
  )

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      {/* Header */}
      <header className="bg-gray-800 border-b border-gray-700 px-6 py-3 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <span className="text-2xl">🍳</span>
          <div>
            <h1 className="font-bold text-lg leading-tight">Bếp — Kitchen</h1>
            <p className="text-xs text-gray-400">{user?.full_name}</p>
          </div>
        </div>
        <div className="flex items-center gap-4">
          <span className={`text-xs px-2 py-1 rounded-full ${connected ? 'bg-green-800 text-green-300' : 'bg-red-900 text-red-300'}`}>
            {connected ? '● Live' : '○ Offline'}
          </span>
          <span className="text-sm text-gray-400">{activeOrders.length} order đang chờ</span>
          <button onClick={onLogout} className="text-sm text-gray-400 hover:text-white transition-colors">
            Đăng xuất
          </button>
        </div>
      </header>

      {/* Notification banner */}
      {notifications[0]?.type === 'ORDER_CONFIRMED' && (
        <div className="bg-orange-600 px-6 py-2 text-sm font-medium animate-pulse">
          🔔 Order mới: {notifications[0].message}
        </div>
      )}

      <main className="p-6">
        {loading ? (
          <p className="text-gray-400 text-center py-12">Đang tải...</p>
        ) : activeOrders.length === 0 ? (
          <div className="text-center py-16 text-gray-500">
            <div className="text-5xl mb-4">✅</div>
            <p className="text-lg">Không có order nào đang chờ</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
            {activeOrders.map((order) => (
              <OrderCard key={order.order_id} order={order} onMark={markItem} />
            ))}
          </div>
        )}
      </main>
    </div>
  )
}

function OrderCard({
  order,
  onMark,
}: {
  order: OrderDto
  onMark: (orderId: string, itemId: string, next: string) => void
}) {
  const pendingCount = order.items?.filter((i) => i.item_status === 'PENDING').length ?? 0
  const cookingCount = order.items?.filter((i) => i.item_status === 'COOKING').length ?? 0
  const readyCount   = order.items?.filter((i) => i.item_status === 'READY').length ?? 0

  return (
    <div className="bg-gray-800 rounded-xl border border-gray-700 overflow-hidden">
      {/* Order header */}
      <div className="px-4 py-3 bg-gray-750 border-b border-gray-700 flex items-start justify-between">
        <div>
          <p className="font-semibold text-white">{order.name}</p>
          <p className="text-xs text-gray-400 mt-0.5">
            {order.party_size} người · Bàn: <span className="text-orange-400 font-mono text-xs">{order.table_id.slice(0, 8)}…</span>
          </p>
          {order.notes && (
            <p className="text-xs text-amber-300 mt-1 bg-amber-900/30 rounded px-2 py-0.5">
              ⚠ {order.notes}
            </p>
          )}
        </div>
        <div className="flex gap-1 text-xs flex-wrap justify-end">
          {pendingCount > 0 && <span className="bg-gray-700 text-gray-300 px-1.5 py-0.5 rounded">{pendingCount} chờ</span>}
          {cookingCount > 0 && <span className="bg-amber-700 text-white px-1.5 py-0.5 rounded">{cookingCount} nấu</span>}
          {readyCount > 0   && <span className="bg-green-700 text-white px-1.5 py-0.5 rounded">{readyCount} xong</span>}
        </div>
      </div>

      {/* Items */}
      <ul className="divide-y divide-gray-700">
        {order.items?.map((item) => {
          const action = nextAction(item)
          return (
            <li key={item.item_id} className="px-4 py-3 flex items-center gap-3">
              <span className={`text-xs px-2 py-0.5 rounded-full font-medium shrink-0 ${STATUS_STYLE[item.item_status] ?? STATUS_STYLE.PENDING}`}>
                {STATUS_LABEL[item.item_status] ?? item.item_status}
              </span>
              <span className={`flex-1 text-sm ${item.item_status === 'SERVED' ? 'text-gray-500' : 'text-white'}`}>
                {item.name} <span className="text-gray-400">×{item.quantity}</span>
              </span>
              {action && (
                <button
                  onClick={() => onMark(order.order_id, item.item_id, action.next)}
                  className="text-xs bg-orange-600 hover:bg-orange-500 text-white px-2.5 py-1 rounded-lg transition-colors shrink-0"
                >
                  {action.label}
                </button>
              )}
            </li>
          )
        })}
      </ul>
    </div>
  )
}

function nextAction(item: OrderItemDto): { label: string; next: string } | null {
  if (item.item_status === 'PENDING')  return { label: 'Bắt đầu nấu', next: 'COOKING' }
  if (item.item_status === 'COOKING')  return { label: 'Đã xong ✓', next: 'READY' }
  return null
}
