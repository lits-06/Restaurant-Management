package domain

import (
	"errors"
	"strings"
	"time"
)

// TableStatus represents the status of a table
type TableStatus string

const (
	StatusAvailable     TableStatus = "AVAILABLE"
	StatusOccupied      TableStatus = "OCCUPIED"
	StatusReserved      TableStatus = "RESERVED"
	StatusCleaning      TableStatus = "CLEANING"
	StatusOutOfService  TableStatus = "OUT_OF_SERVICE"
)

// Table represents a restaurant table
type Table struct {
	ID             string
	TableNumber    string
	Capacity       int
	Status         TableStatus
	Location       string      // e.g., "Main Hall", "VIP Room", "Outdoor"
	CurrentOrderID string      // ID of the current order (if occupied)
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewTable creates a new Table with validation
func NewTable(tableNumber string, capacity int, location string) (*Table, error) {
	table := &Table{
		TableNumber: strings.TrimSpace(tableNumber),
		Capacity:    capacity,
		Location:    strings.TrimSpace(location),
		Status:      StatusAvailable,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := table.Validate(); err != nil {
		return nil, err
	}

	return table, nil
}

// Validate validates the table fields
func (t *Table) Validate() error {
	if t.TableNumber == "" {
		return ErrInvalidTableNumber
	}

	if t.Capacity <= 0 {
		return ErrInvalidCapacity
	}

	if t.Capacity > 50 {
		return ErrCapacityTooLarge
	}

	if t.Location == "" {
		return ErrInvalidLocation
	}

	if !isValidStatus(t.Status) {
		return ErrInvalidStatus
	}

	return nil
}

// IsAvailable checks if the table is available for seating
func (t *Table) IsAvailable() bool {
	return t.Status == StatusAvailable
}

// IsOccupied checks if the table is currently occupied
func (t *Table) IsOccupied() bool {
	return t.Status == StatusOccupied
}

// IsReserved checks if the table is reserved
func (t *Table) IsReserved() bool {
	return t.Status == StatusReserved
}

// CanBeOccupied checks if the table can be occupied (available or reserved)
func (t *Table) CanBeOccupied() bool {
	return t.Status == StatusAvailable || t.Status == StatusReserved
}

// MarkAvailable marks the table as available
func (t *Table) MarkAvailable() error {
	if t.Status == StatusOutOfService {
		return ErrTableOutOfService
	}

	t.Status = StatusAvailable
	t.CurrentOrderID = ""
	t.UpdatedAt = time.Now()
	return nil
}

// MarkOccupied marks the table as occupied
func (t *Table) MarkOccupied(orderID string) error {
	if !t.CanBeOccupied() {
		return ErrTableNotAvailable
	}

	if orderID == "" {
		return ErrInvalidOrderID
	}

	t.Status = StatusOccupied
	t.CurrentOrderID = orderID
	t.UpdatedAt = time.Now()
	return nil
}

// MarkReserved marks the table as reserved
func (t *Table) MarkReserved() error {
	if t.Status != StatusAvailable {
		return ErrTableNotAvailable
	}

	t.Status = StatusReserved
	t.UpdatedAt = time.Now()
	return nil
}

// MarkCleaning marks the table as being cleaned
func (t *Table) MarkCleaning() error {
	if t.Status != StatusOccupied {
		return errors.New("can only mark occupied tables for cleaning")
	}

	t.Status = StatusCleaning
	t.CurrentOrderID = ""
	t.UpdatedAt = time.Now()
	return nil
}

// MarkOutOfService marks the table as out of service
func (t *Table) MarkOutOfService() error {
	if t.Status == StatusOccupied {
		return errors.New("cannot mark occupied table as out of service")
	}

	t.Status = StatusOutOfService
	t.CurrentOrderID = ""
	t.UpdatedAt = time.Now()
	return nil
}

// UpdateStatus updates the table status with validation
func (t *Table) UpdateStatus(status TableStatus, orderID string) error {
	if !isValidStatus(status) {
		return ErrInvalidStatus
	}

	switch status {
	case StatusAvailable:
		return t.MarkAvailable()
	case StatusOccupied:
		return t.MarkOccupied(orderID)
	case StatusReserved:
		return t.MarkReserved()
	case StatusCleaning:
		return t.MarkCleaning()
	case StatusOutOfService:
		return t.MarkOutOfService()
	default:
		return ErrInvalidStatus
	}
}

// Update updates table information
func (t *Table) Update(tableNumber string, capacity int, location string) error {
	if tableNumber != "" {
		tableNumber = strings.TrimSpace(tableNumber)
		if tableNumber == "" {
			return ErrInvalidTableNumber
		}
		t.TableNumber = tableNumber
	}

	if capacity > 0 {
		if capacity > 50 {
			return ErrCapacityTooLarge
		}
		t.Capacity = capacity
	}

	if location != "" {
		location = strings.TrimSpace(location)
		if location == "" {
			return ErrInvalidLocation
		}
		t.Location = location
	}

	t.UpdatedAt = time.Now()
	return nil
}

// HasSufficientCapacity checks if table has enough capacity for given number of guests
func (t *Table) HasSufficientCapacity(numberOfGuests int) bool {
	return t.Capacity >= numberOfGuests
}

// Helper functions

func isValidStatus(status TableStatus) bool {
	switch status {
	case StatusAvailable, StatusOccupied, StatusReserved, StatusCleaning, StatusOutOfService:
		return true
	default:
		return false
	}
}
