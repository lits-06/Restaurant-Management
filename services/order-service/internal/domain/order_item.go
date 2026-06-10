package domain

// ItemStatus represents the kitchen lifecycle of a single order item.
type ItemStatus string

const (
	ItemStatusPending  ItemStatus = "PENDING"
	ItemStatusCooking  ItemStatus = "COOKING"
	ItemStatusReady    ItemStatus = "READY"
	ItemStatusServed   ItemStatus = "SERVED"
)

var allowedItemStatuses = map[ItemStatus]struct{}{
	ItemStatusPending: {},
	ItemStatusCooking: {},
	ItemStatusReady:   {},
	ItemStatusServed:  {},
}

// validItemTransitions defines which next statuses are allowed from each state.
var validItemTransitions = map[ItemStatus]map[ItemStatus]struct{}{
	ItemStatusPending: {ItemStatusCooking: {}},
	ItemStatusCooking: {ItemStatusReady: {}},
	ItemStatusReady:   {ItemStatusServed: {}},
	ItemStatusServed:  {},
}

func (s ItemStatus) IsValid() bool {
	_, ok := allowedItemStatuses[s]
	return ok
}

// CanTransitionTo reports whether transitioning from s to next is allowed.
func (s ItemStatus) CanTransitionTo(next ItemStatus) bool {
	allowed, ok := validItemTransitions[s]
	if !ok {
		return false
	}
	_, ok = allowed[next]
	return ok
}

// OrderItem represents a menu item selected for an order.
type OrderItem struct {
	ItemID     string
	Name       string
	Price      float64
	Quantity   int32
	ItemStatus ItemStatus
}

// Validate validates the order item.
func (i *OrderItem) Validate() error {
	if i.Name == "" {
		return ErrOrderItemNameRequired
	}
	if i.Price < 0 {
		return ErrOrderItemPriceInvalid
	}
	if i.Quantity <= 0 {
		return ErrOrderItemQuantityInvalid
	}
	return nil
}

// NewOrderItem creates a new order item with status PENDING.
func NewOrderItem(itemID string, name string, price float64, quantity int32) (*OrderItem, error) {
	item := &OrderItem{
		ItemID:     itemID,
		Name:       name,
		Price:      price,
		Quantity:   quantity,
		ItemStatus: ItemStatusPending,
	}
	if err := item.Validate(); err != nil {
		return nil, err
	}
	return item, nil
}
