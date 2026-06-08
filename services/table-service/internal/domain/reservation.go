package domain

import (
	"strings"
	"time"
)

// ReservationStatus represents reservation lifecycle state.
type ReservationStatus string

const (
	ReservationStatusReserved  ReservationStatus = "RESERVED"
	ReservationStatusCancelled ReservationStatus = "CANCELLED"
	ReservationStatusCompleted ReservationStatus = "COMPLETED"
)

// ReservationItem represents a menu item requested for a reservation.
type ReservationItem struct {
	MenuItemID string
	Quantity   int
	Note       string
}

// Reservation represents a table reservation.
type Reservation struct {
	ID            string
	TableID       string
	CustomerName  string
	CustomerPhone string
	Notes         string
	Status        ReservationStatus
	StartTime     time.Time
	EndTime       time.Time
	Items         []ReservationItem
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewReservation creates a new Reservation with validation.
func NewReservation(tableID, customerName, customerPhone string, startTime, endTime time.Time, notes string, items []ReservationItem) (*Reservation, error) {
	reservation := &Reservation{
		TableID:       strings.TrimSpace(tableID),
		CustomerName:  strings.TrimSpace(customerName),
		CustomerPhone: strings.TrimSpace(customerPhone),
		Notes:         strings.TrimSpace(notes),
		Status:        ReservationStatusReserved,
		StartTime:     startTime,
		EndTime:       endTime,
		Items:         items,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := reservation.Validate(); err != nil {
		return nil, err
	}

	return reservation, nil
}

// Validate validates reservation fields.
func (r *Reservation) Validate() error {
	if r.TableID == "" {
		return ErrInvalidTableID
	}
	if r.CustomerName == "" {
		return ErrInvalidReservationCustomer
	}
	if r.EndTime.IsZero() || r.StartTime.IsZero() {
		return ErrInvalidReservationTime
	}
	if !r.EndTime.After(r.StartTime) {
		return ErrInvalidReservationTime
	}
	if !isValidReservationStatus(r.Status) {
		return ErrInvalidReservationStatus
	}

	for _, item := range r.Items {
		if strings.TrimSpace(item.MenuItemID) == "" {
			return ErrInvalidReservationItem
		}
		if item.Quantity <= 0 {
			return ErrInvalidReservationItem
		}
	}

	return nil
}

// Cancel marks reservation as cancelled.
func (r *Reservation) Cancel() error {
	if r.Status == ReservationStatusCancelled {
		return ErrReservationAlreadyCancelled
	}
	if r.Status == ReservationStatusCompleted {
		return ErrReservationAlreadyCompleted
	}

	r.Status = ReservationStatusCancelled
	r.UpdatedAt = time.Now()
	return nil
}

// Complete marks reservation as completed.
func (r *Reservation) Complete() error {
	if r.Status == ReservationStatusCancelled {
		return ErrReservationAlreadyCancelled
	}
	if r.Status == ReservationStatusCompleted {
		return ErrReservationAlreadyCompleted
	}

	r.Status = ReservationStatusCompleted
	r.UpdatedAt = time.Now()
	return nil
}

func isValidReservationStatus(status ReservationStatus) bool {
	switch status {
	case ReservationStatusReserved, ReservationStatusCancelled, ReservationStatusCompleted:
		return true
	default:
		return false
	}
}
