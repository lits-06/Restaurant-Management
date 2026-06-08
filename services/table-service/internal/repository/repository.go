package repository

import (
	"context"
	"time"

	"restaurant-management/services/table-service/internal/domain"
)

// TableRepository defines the interface for table data access
type TableRepository interface {
	// Create creates a new table
	Create(ctx context.Context, table *domain.Table) error

	// GetByID retrieves a table by ID
	GetByID(ctx context.Context, id string) (*domain.Table, error)

	// GetByTableNumber retrieves a table by table number
	GetByTableNumber(ctx context.Context, tableNumber string) (*domain.Table, error)

	// Update updates a table
	Update(ctx context.Context, table *domain.Table) error

	// Delete deletes a table by ID
	Delete(ctx context.Context, id string) error

	// List retrieves tables with filters and pagination
	List(ctx context.Context, filters ListFilters) ([]*domain.Table, int, error)

	// GetAvailableTables retrieves available tables with optional filters
	GetAvailableTables(ctx context.Context, minCapacity int, location string) ([]*domain.Table, error)

	// ExistsByTableNumber checks if a table with the given number exists
	ExistsByTableNumber(ctx context.Context, tableNumber string) (bool, error)
}

// ReservationRepository defines the interface for reservation data access.
type ReservationRepository interface {
	// CreateReservation creates a new reservation.
	CreateReservation(ctx context.Context, reservation *domain.Reservation) error

	// GetReservationByID retrieves a reservation by ID.
	GetReservationByID(ctx context.Context, id string) (*domain.Reservation, error)

	// ListReservations retrieves reservations with filters and pagination.
	ListReservations(ctx context.Context, filters ReservationListFilters) ([]*domain.Reservation, int, error)

	// CancelReservation cancels a reservation by ID.
	CancelReservation(ctx context.Context, id string) (*domain.Reservation, error)

	// HasOverlappingReservation checks for overlapping reservations by time range.
	HasOverlappingReservation(ctx context.Context, tableID string, startTime, endTime time.Time) (bool, error)
}

// ListFilters represents filters for listing tables
type ListFilters struct {
	Page     int
	PageSize int
	Status   domain.TableStatus
	Location string
}

// ReservationListFilters represents filters for listing reservations.
type ReservationListFilters struct {
	Page     int
	PageSize int
	TableID  string
	Status   domain.ReservationStatus
	FromTime time.Time
	ToTime   time.Time
}
