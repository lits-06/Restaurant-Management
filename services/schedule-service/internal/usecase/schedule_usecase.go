package usecase

import (
	"context"
	"fmt"

	"restaurant-management/services/schedule-service/internal/domain"
	"restaurant-management/services/schedule-service/internal/repository"
)

type ScheduleUseCase struct {
	repo repository.ShiftRepository
}

func NewScheduleUseCase(repo repository.ShiftRepository) *ScheduleUseCase {
	return &ScheduleUseCase{repo: repo}
}

func (uc *ScheduleUseCase) CreateShift(ctx context.Context, userID, date, startTime, endTime, role, notes, createdBy string) (*domain.Shift, error) {
	shift, err := domain.NewShift(userID, date, startTime, endTime, role, notes, createdBy)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.Create(ctx, shift); err != nil {
		return nil, fmt.Errorf("failed to create shift: %w", err)
	}
	return shift, nil
}

func (uc *ScheduleUseCase) GetShift(ctx context.Context, id string) (*domain.Shift, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *ScheduleUseCase) UpdateShift(ctx context.Context, id, date, startTime, endTime, notes string) (*domain.Shift, error) {
	shift, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if date != "" {
		shift.Date = date
	}
	if startTime != "" {
		shift.StartTime = startTime
	}
	if endTime != "" {
		shift.EndTime = endTime
	}
	shift.Notes = notes

	if err := shift.Validate(); err != nil {
		return nil, err
	}
	if err := uc.repo.Update(ctx, shift); err != nil {
		return nil, fmt.Errorf("failed to update shift: %w", err)
	}
	return shift, nil
}

func (uc *ScheduleUseCase) DeleteShift(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *ScheduleUseCase) ListShifts(ctx context.Context, month, userID, role string, page, pageSize int) ([]*domain.Shift, int, error) {
	return uc.repo.List(ctx, repository.ListFilters{
		Month:    month,
		UserID:   userID,
		Role:     role,
		Page:     page,
		PageSize: pageSize,
	})
}
