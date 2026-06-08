package domain

// OrderItem represents a menu item selected for an order.
type OrderItem struct {
	ItemID   string
	Name     string
	Price    float64
	Quantity int32
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

// NewOrderItem creates a new order item.
func NewOrderItem(itemID string, name string, price float64, quantity int32) (*OrderItem, error) {
	item := &OrderItem{
		ItemID:   itemID,
		Name:     name,
		Price:    price,
		Quantity: quantity,
	}
	if err := item.Validate(); err != nil {
		return nil, err
	}
	return item, nil
}
