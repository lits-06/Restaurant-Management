package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"restaurant-management/services/table-service/internal/domain"
	"restaurant-management/services/table-service/internal/repository"
)

// TableUseCase handles table business logic.
type TableUseCase struct {
	tableRepo repository.TableRepository
}

func NewTableUseCase(tableRepo repository.TableRepository) *TableUseCase {
	return &TableUseCase{tableRepo: tableRepo}
}

func (uc *TableUseCase) CreateTable(ctx context.Context, tableNumber, capacity int) (*domain.Table, error) {
	exists, err := uc.tableRepo.ExistsByTableNumber(ctx, tableNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to check table number: %w", err)
	}
	if exists {
		return nil, domain.ErrTableNumberAlreadyExists
	}

	table, err := domain.NewTable(tableNumber, capacity)
	if err != nil {
		return nil, err
	}
	table.ID = uuid.New().String()

	if err := uc.tableRepo.Create(ctx, table); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	return table, nil
}

func (uc *TableUseCase) GetTable(ctx context.Context, tableID string) (*domain.Table, error) {
	if tableID == "" {
		return nil, domain.ErrInvalidTableID
	}
	table, err := uc.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}
	return table, nil
}

func (uc *TableUseCase) UpdateTable(ctx context.Context, tableID string, tableNumber, capacity int) (*domain.Table, error) {
	table, err := uc.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}
	if err := table.Update(tableNumber, capacity); err != nil {
		return nil, err
	}
	if err := uc.tableRepo.Update(ctx, table); err != nil {
		return nil, fmt.Errorf("failed to update table: %w", err)
	}
	return table, nil
}

func (uc *TableUseCase) DeleteTable(ctx context.Context, tableID string) error {
	if tableID == "" {
		return domain.ErrInvalidTableID
	}
	if err := uc.tableRepo.Delete(ctx, tableID); err != nil {
		return fmt.Errorf("failed to delete table: %w", err)
	}
	return nil
}

func (uc *TableUseCase) ListTables(ctx context.Context, page, pageSize int, status domain.TableStatus) ([]*domain.Table, int, error) {
	tables, total, err := uc.tableRepo.List(ctx, repository.ListFilters{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tables: %w", err)
	}
	return tables, total, nil
}

func (uc *TableUseCase) UpdateTableStatus(ctx context.Context, tableID string, status domain.TableStatus) (*domain.Table, error) {
	if tableID == "" {
		return nil, domain.ErrInvalidTableID
	}
	table, err := uc.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}
	if err := table.UpdateStatus(status); err != nil {
		return nil, err
	}
	if err := uc.tableRepo.Update(ctx, table); err != nil {
		return nil, fmt.Errorf("failed to save table: %w", err)
	}
	return table, nil
}

func (uc *TableUseCase) GetAvailableTables(ctx context.Context, minCapacity int) ([]*domain.Table, error) {
	tables, err := uc.tableRepo.GetAvailableTables(ctx, minCapacity)
	if err != nil {
		return nil, fmt.Errorf("failed to get available tables: %w", err)
	}
	return tables, nil
}
