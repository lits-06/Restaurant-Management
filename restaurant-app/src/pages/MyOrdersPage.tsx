import React, { useCallback, useEffect, useState, useMemo } from 'react';
import { ordersApi, menuApi, type OrderDto, type MenuItemDto } from '../api/gateway.api';
import { useAuthStore } from '../store/authStore';

type OrderTimeValue = string | { seconds?: string | number; nanos?: number } | undefined;

const parseOrderTime = (time: OrderTimeValue): Date | null => {
  if (!time) return null;
  if (typeof time === 'string') return new Date(time);
  if (time.seconds !== undefined) return new Date(Number(time.seconds) * 1000);
  return null;
};

const formatOrderTime = (time: OrderTimeValue): string => {
  const d = parseOrderTime(time);
  if (!d || isNaN(d.getTime())) return '—';
  return d.toLocaleString('vi-VN', { dateStyle: 'medium', timeStyle: 'short' });
};

const canModifyOrder = (order: OrderDto): boolean => {
  if (!['Pending', 'Confirmed'].includes(order.status || '')) return false;
  const t = parseOrderTime(order.time);
  if (!t) return true;
  return new Date() < t;
};

const STATUS_STYLE: Record<string, string> = {
  Pending: 'bg-amber-100 text-amber-800',
  Confirmed: 'bg-blue-100 text-blue-800',
  Completed: 'bg-green-100 text-green-800',
  Cancelled: 'bg-gray-100 text-gray-600',
};

const STATUS_LABEL: Record<string, string> = {
  Pending: 'Chờ xác nhận',
  Confirmed: 'Đã xác nhận',
  Completed: 'Hoàn thành',
  Cancelled: 'Đã hủy',
};

const getMenuItemId = (item: MenuItemDto) => item.item_id ?? item.itemId ?? '';

interface AddItemsModalProps {
  orderId: string;
  onClose: () => void;
  onDone: (updated: OrderDto) => void;
}

const AddItemsModal: React.FC<AddItemsModalProps> = ({ orderId, onClose, onDone }) => {
  const [menuItems, setMenuItems] = useState<MenuItemDto[]>([]);
  const [cart, setCart] = useState<Record<string, number>>({});
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    menuApi
      .listItems({ page: 1, page_size: 100 })
      .then((r) => setMenuItems(r.items ?? []))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const updateQty = (id: string, delta: number) => {
    setCart((prev) => {
      const q = (prev[id] || 0) + delta;
      if (q <= 0) {
        const { [id]: _, ...rest } = prev;
        return rest;
      }
      return { ...prev, [id]: q };
    });
  };

  const total = useMemo(
    () =>
      Object.entries(cart).reduce((sum, [id, qty]) => {
        const item = menuItems.find((m) => getMenuItemId(m) === id);
        return sum + (item?.price ?? 0) * qty;
      }, 0),
    [cart, menuItems]
  );

  const handleConfirm = async () => {
    if (Object.keys(cart).length === 0) return;
    setError('');
    setSubmitting(true);
    try {
      let lastOrder: OrderDto | undefined;
      for (const [id, qty] of Object.entries(cart)) {
        const item = menuItems.find((m) => getMenuItemId(m) === id);
        if (!item) continue;
        const res = await ordersApi.addItem(orderId, {
          item_id: id,
          name: item.name ?? '',
          price: item.price ?? 0,
          quantity: qty,
        });
        lastOrder = res.order;
      }
      if (lastOrder) onDone(lastOrder);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Thêm món thất bại.');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm px-4">
      <div className="bg-surface-container-lowest rounded-2xl shadow-2xl border border-outline-variant/30 w-full max-w-lg max-h-[85vh] flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-outline-variant/30">
          <h3 className="font-bold text-lg text-on-surface">Thêm món vào đơn</h3>
          <button type="button" onClick={onClose} className="text-on-surface-variant hover:text-primary transition-colors">
            <span className="material-symbols-outlined">close</span>
          </button>
        </div>

        {/* Menu list */}
        <div className="flex-1 overflow-y-auto px-6 py-4 space-y-3">
          {loading ? (
            <p className="text-center text-on-surface-variant py-8">Đang tải thực đơn...</p>
          ) : menuItems.length === 0 ? (
            <p className="text-center text-on-surface-variant py-8">Không có món nào.</p>
          ) : (
            menuItems.map((item) => {
              const id = getMenuItemId(item);
              const qty = cart[id] || 0;
              return (
                <div
                  key={id}
                  className={`flex items-center justify-between rounded-xl p-3 border transition-all ${
                    qty > 0 ? 'border-primary/40 bg-primary/5' : 'border-outline-variant/30 bg-surface'
                  }`}
                >
                  <div className="flex-1 min-w-0 mr-4">
                    <p className="font-semibold text-sm text-on-surface truncate">{item.name}</p>
                    <p className="text-xs text-primary font-bold">${item.price?.toFixed(2)}</p>
                  </div>
                  <div className="flex items-center gap-2 rounded-full border border-outline-variant/30 px-2 py-1 bg-surface-container-low">
                    <button
                      type="button"
                      onClick={() => updateQty(id, -1)}
                      className="w-7 h-7 rounded-full flex items-center justify-center hover:bg-surface-container-high text-primary"
                    >
                      <span className="material-symbols-outlined text-base">remove</span>
                    </button>
                    <span className="w-5 text-center text-sm font-bold">{qty}</span>
                    <button
                      type="button"
                      onClick={() => updateQty(id, 1)}
                      className="w-7 h-7 rounded-full flex items-center justify-center bg-primary text-on-primary hover:opacity-90"
                    >
                      <span className="material-symbols-outlined text-base">add</span>
                    </button>
                  </div>
                </div>
              );
            })
          )}
        </div>

        {/* Footer */}
        <div className="px-6 py-4 border-t border-outline-variant/30 space-y-3">
          {total > 0 && (
            <div className="flex justify-between text-sm font-bold text-primary">
              <span>Cộng thêm</span>
              <span>+${total.toFixed(2)}</span>
            </div>
          )}
          {error && <p className="text-sm text-error">{error}</p>}
          <div className="flex gap-3">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 h-11 rounded-lg border border-outline-variant text-on-surface-variant font-semibold text-sm hover:bg-surface-container-low transition-all"
            >
              Hủy
            </button>
            <button
              type="button"
              onClick={handleConfirm}
              disabled={submitting || Object.keys(cart).length === 0}
              className="flex-1 h-11 rounded-lg bg-primary text-on-primary font-semibold text-sm hover:opacity-90 active:scale-[0.98] transition-all disabled:opacity-50"
            >
              {submitting ? 'Đang gửi...' : `Thêm ${Object.keys(cart).length > 0 ? `(${Object.values(cart).reduce((a, b) => a + b, 0)} món)` : ''}`}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

const MyOrdersPage: React.FC = () => {
  const { user } = useAuthStore();
  const [orders, setOrders] = useState<OrderDto[]>([]);
  const [loading, setLoading] = useState(true);
  const [cancellingId, setCancellingId] = useState<string | null>(null);
  const [addItemsOrderId, setAddItemsOrderId] = useState<string | null>(null);
  const [cancelConfirmId, setCancelConfirmId] = useState<string | null>(null);
  const [actionError, setActionError] = useState<Record<string, string>>({});

  const loadOrders = useCallback(() => {
    if (!user) { setLoading(false); return; }
    setLoading(true);
    ordersApi
      .list({ user_id: user.user_id, page: 1, page_size: 50 })
      .then((res) => {
        const fetched = res.orders ?? [];
        fetched.sort((a, b) => {
          const ta = parseOrderTime(a.time);
          const tb = parseOrderTime(b.time);
          if (ta && tb) return tb.getTime() - ta.getTime();
          return 0;
        });
        setOrders(fetched);
      })
      .catch(() => setOrders([]))
      .finally(() => setLoading(false));
  }, [user]);

  useEffect(() => { loadOrders(); }, [loadOrders]);

  const handleCancel = async (orderId: string) => {
    setCancellingId(orderId);
    setActionError((prev) => ({ ...prev, [orderId]: '' }));
    try {
      await ordersApi.cancel(orderId);
      setOrders((prev) =>
        prev.map((o) => (o.order_id === orderId ? { ...o, status: 'Cancelled' } : o))
      );
    } catch (err) {
      setActionError((prev) => ({
        ...prev,
        [orderId]: err instanceof Error ? err.message : 'Hủy đơn thất bại.',
      }));
    } finally {
      setCancellingId(null);
      setCancelConfirmId(null);
    }
  };

  const handleAddItemsDone = (updated: OrderDto) => {
    setOrders((prev) =>
      prev.map((o) => (o.order_id === updated.order_id ? updated : o))
    );
    setAddItemsOrderId(null);
  };

  if (!user) {
    return (
      <div className="flex-1 flex items-center justify-center py-24">
        <div className="text-center space-y-3">
          <span className="material-symbols-outlined text-5xl text-on-surface-variant">lock</span>
          <p className="text-on-surface-variant">Vui lòng đăng nhập để xem đơn hàng.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-background text-on-surface min-h-[60vh] flex-1">
      <main className="max-w-container-max mx-auto px-margin-mobile md:px-margin-desktop py-12">
        {/* Header */}
        <div className="mb-8">
          <h1 className="font-headline-xl text-headline-xl text-on-surface mb-1">Đơn đặt bàn của tôi</h1>
          <p className="text-body-md text-on-surface-variant">
            Xem lại lịch sử đặt bàn, thêm món hoặc hủy đơn trước giờ ăn.
          </p>
        </div>

        {loading ? (
          <div className="flex items-center justify-center py-24">
            <span className="material-symbols-outlined animate-spin text-4xl text-primary">progress_activity</span>
          </div>
        ) : orders.length === 0 ? (
          <div className="text-center py-24 space-y-4">
            <div className="w-20 h-20 bg-primary-container/20 rounded-full flex items-center justify-center mx-auto">
              <span className="material-symbols-outlined text-4xl text-primary">receipt_long</span>
            </div>
            <h3 className="font-bold text-lg text-on-surface">Chưa có đơn nào</h3>
            <p className="text-on-surface-variant text-sm">Đặt bàn để bắt đầu trải nghiệm tại LuxeBistro.</p>
          </div>
        ) : (
          <div className="space-y-4">
            {orders.map((order) => {
              const modifiable = canModifyOrder(order);
              const isCancel = cancellingId === order.order_id;

              return (
                <div
                  key={order.order_id}
                  className="bg-surface-container-lowest rounded-xl border border-outline-variant/30 shadow-sm overflow-hidden"
                >
                  {/* Order header row */}
                  <div className="flex flex-wrap items-center justify-between gap-3 px-6 py-4 border-b border-outline-variant/20">
                    <div className="space-y-0.5">
                      <div className="flex items-center gap-2">
                        <span className="font-bold text-on-surface text-sm">
                          #{order.order_id?.slice(-8).toUpperCase()}
                        </span>
                        <span
                          className={`text-xs font-semibold px-2 py-0.5 rounded-full ${
                            STATUS_STYLE[order.status ?? ''] ?? 'bg-gray-100 text-gray-600'
                          }`}
                        >
                          {STATUS_LABEL[order.status ?? ''] ?? order.status}
                        </span>
                      </div>
                      <p className="text-xs text-on-surface-variant">
                        {order.name} · {order.phone}
                      </p>
                    </div>
                    <div className="text-right space-y-0.5">
                      <p className="text-xs font-medium text-on-surface">{formatOrderTime(order.time)}</p>
                      {order.end_time && (
                        <p className="text-xs text-on-surface-variant">đến {formatOrderTime(order.end_time)}</p>
                      )}
                    </div>
                  </div>

                  {/* Body */}
                  <div className="px-6 py-4 grid md:grid-cols-2 gap-6">
                    {/* Info */}
                    <div className="space-y-2 text-sm">
                      {order.party_size && (
                        <div className="flex items-center gap-2 text-on-surface-variant">
                          <span className="material-symbols-outlined text-base">group</span>
                          <span>{order.party_size} khách</span>
                        </div>
                      )}
                      {order.table_id && (
                        <div className="flex items-center gap-2 text-on-surface-variant">
                          <span className="material-symbols-outlined text-base">table_restaurant</span>
                          <span>Bàn {order.table_id.slice(-4).toUpperCase()}</span>
                        </div>
                      )}
                      {order.notes && (
                        <div className="flex items-start gap-2 text-on-surface-variant">
                          <span className="material-symbols-outlined text-base mt-0.5">note</span>
                          <span className="italic">{order.notes}</span>
                        </div>
                      )}
                    </div>

                    {/* Items */}
                    <div>
                      {(order.items ?? []).length > 0 ? (
                        <div className="space-y-1.5">
                          <p className="text-xs font-semibold text-primary uppercase tracking-wider mb-2">Pre-order</p>
                          {(order.items ?? []).map((item, i) => (
                            <div key={i} className="flex justify-between text-sm">
                              <span className="text-on-surface">{item.name ?? item.item_id}</span>
                              <span className="text-on-surface-variant font-medium">×{item.quantity}</span>
                            </div>
                          ))}
                          {(order.total_price ?? 0) > 0 && (
                            <div className="flex justify-between text-sm font-bold text-primary border-t border-dashed border-outline-variant pt-2 mt-2">
                              <span>Tổng</span>
                              <span>${order.total_price}</span>
                            </div>
                          )}
                        </div>
                      ) : (
                        <p className="text-sm text-on-surface-variant italic">Chưa có món pre-order</p>
                      )}
                    </div>
                  </div>

                  {/* Actions */}
                  {(modifiable || actionError[order.order_id ?? '']) && (
                    <div className="px-6 py-3 border-t border-outline-variant/20 bg-surface-container-low/50 flex flex-wrap items-center gap-3">
                      {modifiable && (
                        <>
                          <button
                            type="button"
                            onClick={() => setAddItemsOrderId(order.order_id ?? '')}
                            className="flex items-center gap-1.5 text-sm font-semibold text-primary hover:bg-primary/10 px-3 py-1.5 rounded-lg transition-all"
                          >
                            <span className="material-symbols-outlined text-base">add_circle</span>
                            Thêm món
                          </button>

                          {cancelConfirmId !== order.order_id ? (
                            <button
                              type="button"
                              onClick={() => setCancelConfirmId(order.order_id ?? '')}
                              className="flex items-center gap-1.5 text-sm font-semibold text-error hover:bg-error/10 px-3 py-1.5 rounded-lg transition-all"
                            >
                              <span className="material-symbols-outlined text-base">cancel</span>
                              Hủy đơn
                            </button>
                          ) : (
                            <div className="flex items-center gap-2">
                              <span className="text-sm text-on-surface-variant">Xác nhận hủy?</span>
                              <button
                                type="button"
                                onClick={() => handleCancel(order.order_id ?? '')}
                                disabled={isCancel}
                                className="text-sm font-bold text-white bg-error px-3 py-1.5 rounded-lg hover:opacity-90 disabled:opacity-60 transition-all"
                              >
                                {isCancel ? '...' : 'Xác nhận'}
                              </button>
                              <button
                                type="button"
                                onClick={() => setCancelConfirmId(null)}
                                className="text-sm text-on-surface-variant hover:text-on-surface px-2 py-1.5 transition-all"
                              >
                                Không
                              </button>
                            </div>
                          )}
                        </>
                      )}
                      {actionError[order.order_id ?? ''] && (
                        <p className="text-sm text-error ml-auto">{actionError[order.order_id ?? '']}</p>
                      )}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </main>

      {/* Add Items Modal */}
      {addItemsOrderId && (
        <AddItemsModal
          orderId={addItemsOrderId}
          onClose={() => setAddItemsOrderId(null)}
          onDone={handleAddItemsDone}
        />
      )}
    </div>
  );
};

export default MyOrdersPage;
