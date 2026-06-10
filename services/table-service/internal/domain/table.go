package domain

import (
	"time"
)

type TableStatus string

const (
	StatusAvailable    TableStatus = "AVAILABLE"
	StatusCleaning     TableStatus = "CLEANING"
	StatusOutOfService TableStatus = "OUT_OF_SERVICE"
)

type Table struct {
	ID          string
	TableNumber int
	Capacity    int
	Status      TableStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTable(tableNumber, capacity int) (*Table, error) {
	table := &Table{
		TableNumber: tableNumber,
		Capacity:    capacity,
		Status:      StatusAvailable,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := table.Validate(); err != nil {
		return nil, err
	}
	return table, nil
}

func (t *Table) Validate() error {
	if t.TableNumber <= 0 {
		return ErrInvalidTableNumber
	}
	if t.Capacity <= 0 {
		return ErrInvalidCapacity
	}
	if t.Capacity > 50 {
		return ErrCapacityTooLarge
	}
	if !isValidStatus(t.Status) {
		return ErrInvalidStatus
	}
	return nil
}

func (t *Table) IsAvailable() bool {
	return t.Status == StatusAvailable
}

// MarkCleaning transitions AVAILABLE → CLEANING.
func (t *Table) MarkCleaning() error {
	if t.Status != StatusAvailable {
		return ErrTableNotAvailable
	}
	t.Status = StatusCleaning
	t.UpdatedAt = time.Now()
	return nil
}

// MarkAvailable transitions any state → AVAILABLE.
func (t *Table) MarkAvailable() error {
	t.Status = StatusAvailable
	t.UpdatedAt = time.Now()
	return nil
}

// MarkOutOfService transitions any state → OUT_OF_SERVICE.
func (t *Table) MarkOutOfService() error {
	t.Status = StatusOutOfService
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Table) UpdateStatus(status TableStatus) error {
	if !isValidStatus(status) {
		return ErrInvalidStatus
	}
	switch status {
	case StatusAvailable:
		return t.MarkAvailable()
	case StatusCleaning:
		return t.MarkCleaning()
	case StatusOutOfService:
		return t.MarkOutOfService()
	default:
		return ErrInvalidStatus
	}
}

func (t *Table) Update(tableNumber, capacity int) error {
	if tableNumber > 0 {
		t.TableNumber = tableNumber
	}
	if capacity > 0 {
		if capacity > 50 {
			return ErrCapacityTooLarge
		}
		t.Capacity = capacity
	}
	t.UpdatedAt = time.Now()
	return nil
}

func isValidStatus(status TableStatus) bool {
	switch status {
	case StatusAvailable, StatusCleaning, StatusOutOfService:
		return true
	default:
		return false
	}
}
