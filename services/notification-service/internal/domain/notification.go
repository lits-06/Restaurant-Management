package domain

const (
	TypeOrderConfirmed = "ORDER_CONFIRMED"
	TypeItemReady      = "ITEM_READY"

	RoleChef   = "CHEF"
	RoleWaiter = "WAITER"
)

type OrderItem struct {
	ItemID   string `json:"item_id"`
	ItemName string `json:"item_name"`
	Quantity int32  `json:"quantity"`
}

type Notification struct {
	ID           string      `json:"id"`
	Type         string      `json:"type"`
	TargetRole   string      `json:"target_role"`
	OrderID      string      `json:"order_id"`
	TableID      string      `json:"table_id"`
	ItemID       string      `json:"item_id"`
	ItemName     string      `json:"item_name"`
	CreatedAt    int64       `json:"created_at"`
	Message      string      `json:"message"`
	CustomerName string      `json:"customer_name"`
	PartySize    int32       `json:"party_size"`
	Notes        string      `json:"notes"`
	Items        []OrderItem `json:"items,omitempty"`
}
