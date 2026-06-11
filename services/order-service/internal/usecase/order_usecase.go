package usecase

import (
	"context"
	"fmt"

	"restaurant-management/proto/menu"
	notifpb "restaurant-management/proto/notification"
	"restaurant-management/services/order-service/internal/domain"
	"restaurant-management/services/order-service/internal/repository"
)

// MenuServiceClient validates menu items via menu-service gRPC.
type MenuServiceClient interface {
	ValidateMenuItem(ctx context.Context, itemID string) (*menu.MenuItem, error)
}

// TableInfo is the minimal table data needed for auto-assignment.
type TableInfo struct {
	TableID  string
	Capacity int32
}

// TableServiceClient fetches physically available tables from table-service.
// Called during CreateOrder when no table_id is provided.
// If nil, the caller must supply table_id explicitly.
type TableServiceClient interface {
	GetAvailableTables(ctx context.Context, minCapacity int32) ([]TableInfo, error)
}

// NotificationServiceClient sends push notifications to kitchen/waiter staff.
// If nil, notifications are silently skipped.
type NotificationServiceClient interface {
	SendNotification(ctx context.Context, req *notifpb.SendNotificationRequest) (*notifpb.SendNotificationResponse, error)
}

// OrderUseCase handles order business logic.
type OrderUseCase struct {
	orderRepo   repository.OrderRepository
	menuClient  MenuServiceClient
	tableClient  TableServiceClient          // optional; nil = auto-assign disabled
	notifClient NotificationServiceClient    // optional; nil = notifications disabled
}

// NewOrderUseCase creates a new OrderUseCase.
// tableClient and notifClient may be nil.
func NewOrderUseCase(
	orderRepo repository.OrderRepository,
	menuClient MenuServiceClient,
	tableClient TableServiceClient,
	notifClient NotificationServiceClient,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo:   orderRepo,
		menuClient:  menuClient,
		tableClient: tableClient,
		notifClient: notifClient,
	}
}

// OrderItemRequest represents a request to add an order item.
type OrderItemRequest struct {
	MenuItemID string
	Quantity   int32
}

// CreateOrder creates a new order.
// If tableID is empty and tableClient is configured, a suitable table is auto-assigned
// based on party size and time-slot availability.
func (uc *OrderUseCase) CreateOrder(ctx context.Context, tableID, userID, name, phone, notes, timeStr, endTimeStr, date string, partySize int32, status string, itemRequests []OrderItemRequest) (*domain.Order, error) {
	items, err := uc.resolveItems(ctx, itemRequests)
	if err != nil {
		return nil, err
	}

	order, err := domain.NewOrder(tableID, userID, name, phone, notes, timeStr, endTimeStr, date, partySize, domain.OrderStatus(status), items)
	if err != nil {
		return nil, fmt.Errorf("failed to build order: %w", err)
	}

	if order.TableID == "" {
		assigned, err := uc.autoAssignTable(ctx, partySize, order)
		if err != nil {
			return nil, err
		}
		order.TableID = assigned
	}

	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}
	return order, nil
}

// autoAssignTable finds the smallest available table for the order's time window.
func (uc *OrderUseCase) autoAssignTable(ctx context.Context, partySize int32, order *domain.Order) (string, error) {
	if uc.tableClient == nil {
		return "", domain.ErrTableRequired
	}

	// 1. Physically available tables with sufficient capacity, sorted by capacity ASC (best-fit).
	candidates, err := uc.tableClient.GetAvailableTables(ctx, partySize)
	if err != nil {
		return "", fmt.Errorf("failed to fetch available tables: %w", err)
	}
	if len(candidates) == 0 {
		return "", domain.ErrNoTableAvailable
	}

	// 2. Table IDs already booked in the requested time window.
	occupiedIDs, err := uc.orderRepo.GetOccupiedTableIDs(ctx, order.Time, order.EndTime)
	if err != nil {
		return "", fmt.Errorf("failed to check occupied tables: %w", err)
	}
	occupied := make(map[string]struct{}, len(occupiedIDs))
	for _, id := range occupiedIDs {
		occupied[id] = struct{}{}
	}

	// 3. Pick the first candidate not occupied — smallest capacity that fits.
	for _, t := range candidates {
		if _, busy := occupied[t.TableID]; !busy {
			return t.TableID, nil
		}
	}

	return "", domain.ErrNoTableAvailable
}

// resolveItems validates each menu item and builds domain OrderItems.
func (uc *OrderUseCase) resolveItems(ctx context.Context, reqs []OrderItemRequest) ([]*domain.OrderItem, error) {
	items := make([]*domain.OrderItem, 0, len(reqs))
	for _, req := range reqs {
		item, err := uc.menuClient.ValidateMenuItem(ctx, req.MenuItemID)
		if err != nil {
			return nil, fmt.Errorf("menu item validation failed for %s: %w", req.MenuItemID, err)
		}
		oItem, err := domain.NewOrderItem(item.ItemId, item.Name, item.Price, req.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
		items = append(items, oItem)
	}
	return items, nil
}

// GetOrder retrieves an order by ID.
func (uc *OrderUseCase) GetOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return order, nil
}

// UpdateOrder updates an existing order.
// notes is always overwritten (pass the current value to keep it, empty string to clear it).
func (uc *OrderUseCase) UpdateOrder(ctx context.Context, orderID, tableID, name, phone, notes, timeStr, endTimeStr, date string, partySize int32, status string, itemRequests []OrderItemRequest) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	if tableID != "" {
		order.TableID = tableID
	}
	if name != "" {
		order.Name = name
	}
	if phone != "" {
		order.Phone = phone
	}
	// notes is always updated so callers can clear it intentionally
	order.Notes = notes

	if timeStr != "" && date != "" {
		order.Time, _ = domain.ParseReservationTime(date, timeStr)
		if endTimeStr != "" {
			order.EndTime, _ = domain.ParseReservationTime(date, endTimeStr)
		}
	}
	if partySize > 0 {
		order.PartySize = partySize
	}
	if status != "" {
		order.Status = domain.OrderStatus(status)
	}

	if itemRequests != nil {
		newItems, err := uc.resolveItems(ctx, itemRequests)
		if err != nil {
			return nil, err
		}
		// Preserve item_status for items already in the order so COOKING/READY/SERVED
		// status is not reset to PENDING by an admin edit.
		existingStatus := make(map[string]domain.ItemStatus)
		for _, existing := range order.Items {
			existingStatus[existing.ItemID] = existing.ItemStatus
		}
		for _, item := range newItems {
			if status, ok := existingStatus[item.ItemID]; ok {
				item.ItemStatus = status
			}
		}
		order.Items = newItems
	}

	order.Total = order.TotalPrice()
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}
	return order, nil
}

// CancelOrder cancels an order.
func (uc *OrderUseCase) CancelOrder(ctx context.Context, orderID string) error {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if err := order.Cancel(); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to save cancelled order: %w", err)
	}
	return nil
}

// ListOrders retrieves orders with pagination and filters.
func (uc *OrderUseCase) ListOrders(ctx context.Context, page, pageSize int, status domain.OrderStatus, keyword, userID string) ([]*domain.Order, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	orders, total, err := uc.orderRepo.List(ctx, page, pageSize, status, keyword, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list orders: %w", err)
	}
	return orders, total, nil
}

// UpdateOrderStatus updates the status of an order.
// When transitioning to Confirmed, notifies CHEF with full order details.
func (uc *OrderUseCase) UpdateOrderStatus(ctx context.Context, orderID string, newStatus domain.OrderStatus) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if err := order.UpdateStatus(newStatus); err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}
	if newStatus == domain.StatusConfirmed && uc.notifClient != nil {
		uc.notifyChef(ctx, order)
	}
	return order, nil
}

// AddOrderItem adds an item to an order.
func (uc *OrderUseCase) AddOrderItem(ctx context.Context, orderID string, req OrderItemRequest) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	item, err := uc.menuClient.ValidateMenuItem(ctx, req.MenuItemID)
	if err != nil {
		return nil, fmt.Errorf("menu item validation failed: %w", err)
	}
	oItem, err := domain.NewOrderItem(item.ItemId, item.Name, item.Price, req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to create order item: %w", err)
	}
	if err := order.AddItem(oItem); err != nil {
		return nil, fmt.Errorf("failed to add item: %w", err)
	}
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}
	return order, nil
}

// UpdateOrderItemStatus advances the item_status of a single item within an order.
// Validates the transition (PENDING→COOKING→READY→SERVED), then persists only
// the item row (no full order update needed). Returns the full refreshed order.
// When transitioning to READY, notifies WAITER.
func (uc *OrderUseCase) UpdateOrderItemStatus(ctx context.Context, orderID, itemID string, next domain.ItemStatus) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	// Validate transition via domain logic before writing.
	if err := order.UpdateItemStatus(itemID, next); err != nil {
		return nil, err
	}
	if err := uc.orderRepo.UpdateItemStatus(ctx, orderID, itemID, next); err != nil {
		return nil, fmt.Errorf("failed to update item status: %w", err)
	}
	if next == domain.ItemStatusReady && uc.notifClient != nil {
		// Find item name for notification.
		for _, item := range order.Items {
			if item.ItemID == itemID {
				uc.notifyWaiter(ctx, order, item)
				break
			}
		}
	}
	// Return fresh order so caller has up-to-date item list.
	return uc.orderRepo.GetByID(ctx, orderID)
}

// DeleteOrder permanently deletes an order and its items.
func (uc *OrderUseCase) DeleteOrder(ctx context.Context, orderID string) error {
	if _, err := uc.orderRepo.GetByID(ctx, orderID); err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if err := uc.orderRepo.Delete(ctx, orderID); err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	return nil
}

// RemoveOrderItem removes an item from an order.
func (uc *OrderUseCase) RemoveOrderItem(ctx context.Context, orderID, itemID string) (*domain.Order, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	if err := order.RemoveItem(itemID); err != nil {
		return nil, fmt.Errorf("failed to remove item: %w", err)
	}
	if len(order.Items) == 0 {
		return nil, fmt.Errorf("cannot remove last item from order")
	}
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}
	return order, nil
}

// notifyChef sends ORDER_CONFIRMED notification to CHEF channel.
// Runs in a background goroutine — failures are silently ignored so they don't block order updates.
func (uc *OrderUseCase) notifyChef(ctx context.Context, order *domain.Order) {
	pbItems := make([]*notifpb.NotificationOrderItem, 0, len(order.Items))
	for _, it := range order.Items {
		pbItems = append(pbItems, &notifpb.NotificationOrderItem{
			ItemId:   it.ItemID,
			ItemName: it.Name,
			Quantity: it.Quantity,
		})
	}
	req := &notifpb.SendNotificationRequest{
		Type:         "ORDER_CONFIRMED",
		TargetRole:   "CHEF",
		OrderId:      order.OrderID,
		TableId:      order.TableID,
		CustomerName: order.Name,
		PartySize:    order.PartySize,
		Notes:        order.Notes,
		Message:      fmt.Sprintf("Order mới - %s - %d người", order.Name, order.PartySize),
		Items:        pbItems,
	}
	go func() {
		bgCtx := context.Background()
		uc.notifClient.SendNotification(bgCtx, req) //nolint:errcheck
	}()
}

// notifyWaiter sends ITEM_READY notification to WAITER channel.
func (uc *OrderUseCase) notifyWaiter(ctx context.Context, order *domain.Order, item *domain.OrderItem) {
	req := &notifpb.SendNotificationRequest{
		Type:       "ITEM_READY",
		TargetRole: "WAITER",
		OrderId:    order.OrderID,
		TableId:    order.TableID,
		ItemId:     item.ItemID,
		ItemName:   item.Name,
		Message:    fmt.Sprintf("%s sẵn sàng - Bàn %s", item.Name, order.TableID),
	}
	go func() {
		bgCtx := context.Background()
		uc.notifClient.SendNotification(bgCtx, req) //nolint:errcheck
	}()
}
