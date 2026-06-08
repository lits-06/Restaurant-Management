package usecase

import (
	"context"
	"fmt"

	"restaurant-management/services/staff-service/internal/domain"
	"restaurant-management/services/staff-service/internal/repository"
)

// StaffUseCase handles staff business logic.
type StaffUseCase struct {
	staffRepo repository.StaffRepository
}

// NewStaffUseCase creates a new StaffUseCase.
func NewStaffUseCase(staffRepo repository.StaffRepository) *StaffUseCase {
	return &StaffUseCase{staffRepo: staffRepo}
}

// CreateStaff creates a new staff member.
func (uc *StaffUseCase) CreateStaff(ctx context.Context, name, role, contact, avatar string) (*domain.Staff, error) {
	staff, err := domain.NewStaff(name, role, contact, avatar)
	if err != nil {
		return nil, fmt.Errorf("failed to create staff entity: %w", err)
	}

	if err := uc.staffRepo.Create(ctx, staff); err != nil {
		return nil, fmt.Errorf("failed to create staff: %w", err)
	}

	return staff, nil
}

// GetStaff retrieves a staff member by ID.
func (uc *StaffUseCase) GetStaff(ctx context.Context, staffID string) (*domain.Staff, error) {
	if staffID == "" {
		return nil, domain.ErrInvalidStaffID
	}

	staff, err := uc.staffRepo.GetByID(ctx, staffID)
	if err != nil {
		return nil, fmt.Errorf("failed to get staff: %w", err)
	}

	return staff, nil
}

// UpdateStaff updates staff information.
func (uc *StaffUseCase) UpdateStaff(ctx context.Context, staffID, name, role, contact, avatar string) (*domain.Staff, error) {
	if staffID == "" {
		return nil, domain.ErrInvalidStaffID
	}

	staff, err := uc.staffRepo.GetByID(ctx, staffID)
	if err != nil {
		return nil, fmt.Errorf("failed to get staff: %w", err)
	}

	if err := staff.Update(name, role, contact, avatar); err != nil {
		return nil, fmt.Errorf("failed to update staff fields: %w", err)
	}

	if err := uc.staffRepo.Update(ctx, staff); err != nil {
		return nil, fmt.Errorf("failed to update staff: %w", err)
	}

	return staff, nil
}

// DeleteStaff deletes a staff member.
func (uc *StaffUseCase) DeleteStaff(ctx context.Context, staffID string) error {
	if staffID == "" {
		return domain.ErrInvalidStaffID
	}

	if _, err := uc.staffRepo.GetByID(ctx, staffID); err != nil {
		return fmt.Errorf("failed to get staff: %w", err)
	}

	if err := uc.staffRepo.Delete(ctx, staffID); err != nil {
		return fmt.Errorf("failed to delete staff: %w", err)
	}

	return nil
}

// ListStaff retrieves staff members with pagination and keyword search.
func (uc *StaffUseCase) ListStaff(ctx context.Context, page, pageSize int, keyword string) ([]*domain.Staff, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	staffMembers, total, err := uc.staffRepo.List(ctx, page, pageSize, keyword)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list staff: %w", err)
	}

	return staffMembers, total, nil
}
