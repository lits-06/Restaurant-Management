package domain

import (
	"fmt"
	"strings"
	"time"
)


// OrderStatus represents the reservation/order state shown in the admin UI.
type OrderStatus string

const (
	StatusPending   OrderStatus = "Pending"
	StatusConfirmed OrderStatus = "Confirmed"
	StatusCompleted OrderStatus = "Completed"
	StatusCancelled OrderStatus = "Cancelled"
)

var allowedOrderStatuses = map[OrderStatus]struct{}{
	StatusPending:   {},
	StatusConfirmed: {},
	StatusCompleted: {},
	StatusCancelled: {},
}

// Order represents a reservation/order record used by the admin UI.
type Order struct {
	OrderID   string
	TableID   string
	UserID    string
	Name      string
	Phone     string
	Notes     string
	Time      time.Time
	EndTime   time.Time
	PartySize int32
	Status    OrderStatus
	Total     float64
	Items     []*OrderItem
}

// Validate validates the order record.
func (o *Order) Validate() error {
	if o.Name == "" {
		return ErrOrderNameRequired
	}
	if o.Phone == "" {
		return ErrOrderPhoneRequired
	}
	if o.Time.IsZero() {
		return ErrOrderTimeRequired
	}
	if !o.EndTime.IsZero() && !o.EndTime.After(o.Time) {
		return ErrOrderEndTimeInvalid
	}
	if o.PartySize <= 0 {
		return ErrOrderPartySizeInvalid
	}
	if !o.Status.IsValid() {
		return ErrOrderStatusInvalid
	}
	for _, item := range o.Items {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// IsValid reports whether the status is one of the supported UI states.
func (s OrderStatus) IsValid() bool {
	_, ok := allowedOrderStatuses[s]
	return ok
}

// String returns the string representation of the status.
func (s OrderStatus) String() string {
	return string(s)
}

// NormalizeStatus converts free-form input into a supported status value.
func NormalizeStatus(raw string) OrderStatus {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "pending":
		return StatusPending
	case "confirmed":
		return StatusConfirmed
	case "completed":
		return StatusCompleted
	case "cancelled", "canceled":
		return StatusCancelled
	default:
		return ""
	}
}

// UpdateStatus updates the order status.
func (o *Order) UpdateStatus(newStatus OrderStatus) error {
	if !newStatus.IsValid() {
		return ErrOrderStatusInvalid
	}
	if o.Status == StatusCompleted && newStatus == StatusCancelled {
		return ErrOrderInvalidStatusTransition
	}
	if o.Status == StatusCancelled && newStatus != StatusCancelled {
		return ErrOrderInvalidStatusTransition
	}
	o.Status = newStatus
	return nil
}

// Cancel marks the order as cancelled.
func (o *Order) Cancel() error {
	if o.Status == StatusCancelled {
		return ErrOrderAlreadyCancelled
	}
	if o.Status == StatusCompleted {
		return ErrOrderCannotCancelCompleted
	}
	o.Status = StatusCancelled
	return nil
}

// AddItem appends an order item.
func (o *Order) AddItem(item *OrderItem) error {
	if err := item.Validate(); err != nil {
		return err
	}
	o.Items = append(o.Items, item)
	return nil
}

// RemoveItem removes an order item by ID.
func (o *Order) RemoveItem(itemID string) error {
	for index, item := range o.Items {
		if item.ItemID == itemID {
			o.Items = append(o.Items[:index], o.Items[index+1:]...)
			return nil
		}
	}
	return ErrOrderItemNotFound
}

// NewOrder creates a new order record.
func NewOrder(tableID, userID, name, phone, notes, timeValue, endTimeValue, date string, partySize int32, status OrderStatus, items []*OrderItem) (*Order, error) {
	if status == "" {
		status = StatusPending
	}
	startTime, err := ParseReservationTime(date, timeValue)
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %w", err)
	}

	var endTime time.Time
	if strings.TrimSpace(endTimeValue) != "" {
		endTime, err = ParseReservationTime(date, endTimeValue)
		if err != nil {
			return nil, fmt.Errorf("invalid end time: %w", err)
		}
	}

	order := &Order{
		TableID:   strings.TrimSpace(tableID),
		UserID:    strings.TrimSpace(userID),
		Name:      strings.TrimSpace(name),
		Phone:     strings.TrimSpace(phone),
		Notes:     strings.TrimSpace(notes),
		Time:      startTime,
		EndTime:   endTime,
		PartySize: partySize,
		Status:    status,
		Items:     items,
	}
	order.Total = order.TotalPrice()
	if err := order.Validate(); err != nil {
		return nil, err
	}
	return order, nil
}

func ParseReservationTime(dateStr, timeStr string) (time.Time, error) {
	return time.Parse(
		"2006-01-02 15:04",
		dateStr+" "+timeStr,
	)
}

func (o *Order) TotalPrice() float64 {
	var total float64

	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}

	return total
}

// UpdateItemStatus advances the status of a single item, enforcing valid transitions.
func (o *Order) UpdateItemStatus(itemID string, next ItemStatus) error {
	if !next.IsValid() {
		return ErrOrderItemStatusInvalid
	}
	for _, item := range o.Items {
		if item.ItemID == itemID {
			if !item.ItemStatus.CanTransitionTo(next) {
				return ErrOrderItemInvalidStatusTransition
			}
			item.ItemStatus = next
			return nil
		}
	}
	return ErrOrderItemNotFound
}

// AllItemsServed reports whether every item in the order has been served.
func (o *Order) AllItemsServed() bool {
	if len(o.Items) == 0 {
		return false
	}
	for _, item := range o.Items {
		if item.ItemStatus != ItemStatusServed {
			return false
		}
	}
	return true
}

// Clone creates a shallow clone with copied item slice.
func (o *Order) Clone() *Order {
	if o == nil {
		return nil
	}
	clone := *o
	if o.Items != nil {
		clone.Items = make([]*OrderItem, len(o.Items))
		copy(clone.Items, o.Items)
	}
	return &clone
}

func (s OrderStatus) GoString() string {
	return fmt.Sprintf("OrderStatus(%q)", string(s))
}
