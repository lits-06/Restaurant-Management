package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"restaurant-management/services/table-service/internal/domain"
	"restaurant-management/services/table-service/internal/repository"
)

// ReservationUseCase handles reservation business logic.
type ReservationUseCase struct {
	tableRepo       repository.TableRepository
	reservationRepo repository.ReservationRepository
}

// NewReservationUseCase creates a new ReservationUseCase.
func NewReservationUseCase(tableRepo repository.TableRepository, reservationRepo repository.ReservationRepository) *ReservationUseCase {
	return &ReservationUseCase{
		tableRepo:       tableRepo,
		reservationRepo: reservationRepo,
	}
}

// CreateReservation creates a new table reservation.
func (uc *ReservationUseCase) CreateReservation(ctx context.Context, tableID, customerName, customerPhone string, startTime, endTime time.Time, notes string, items []domain.ReservationItem) (*domain.Reservation, error) {
	if tableID == "" {
		return nil, domain.ErrInvalidTableID
	}

	if _, err := uc.tableRepo.GetByID(ctx, tableID); err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}

	overlap, err := uc.reservationRepo.HasOverlappingReservation(ctx, tableID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to check reservation conflicts: %w", err)
	}
	if overlap {
		return nil, domain.ErrReservationConflict
	}

	reservation, err := domain.NewReservation(tableID, customerName, customerPhone, startTime, endTime, notes, items)
	if err != nil {
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	reservation.ID = uuid.New().String()

	if err := uc.reservationRepo.CreateReservation(ctx, reservation); err != nil {
		return nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	return reservation, nil
}

// GetReservation retrieves a reservation by ID.
func (uc *ReservationUseCase) GetReservation(ctx context.Context, reservationID string) (*domain.Reservation, error) {
	if reservationID == "" {
		return nil, domain.ErrReservationNotFound
	}

	reservation, err := uc.reservationRepo.GetReservationByID(ctx, reservationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}

	return reservation, nil
}

// ListReservations retrieves reservations with filters and pagination.
func (uc *ReservationUseCase) ListReservations(ctx context.Context, filters repository.ReservationListFilters) ([]*domain.Reservation, int, error) {
	reservations, total, err := uc.reservationRepo.ListReservations(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list reservations: %w", err)
	}

	return reservations, total, nil
}

// CancelReservation cancels an active reservation.
func (uc *ReservationUseCase) CancelReservation(ctx context.Context, reservationID string) (*domain.Reservation, error) {
	if reservationID == "" {
		return nil, domain.ErrReservationNotFound
	}

	reservation, err := uc.reservationRepo.CancelReservation(ctx, reservationID)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel reservation: %w", err)
	}

	return reservation, nil
}
