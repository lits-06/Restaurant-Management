package usecase

import (
	"context"
	"fmt"
	"restaurant-management/services/table-service/internal/domain"
	"restaurant-management/services/table-service/internal/repository"

	"github.com/google/uuid"
)

// TableUseCase handles table business logic
type TableUseCase struct {
	tableRepo repository.TableRepository
}

// NewTableUseCase creates a new TableUseCase
func NewTableUseCase(tableRepo repository.TableRepository) *TableUseCase {
	return &TableUseCase{
		tableRepo: tableRepo,
	}
}

// CreateTable creates a new table
func (uc *TableUseCase) CreateTable(ctx context.Context, tableNumber string, capacity int, location string) (*domain.Table, error) {
	// Check if table number already exists
	exists, err := uc.tableRepo.ExistsByTableNumber(ctx, tableNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to check table number existence: %w", err)
	}
	if exists {
		return nil, domain.ErrTableNumberAlreadyExists
	}

	// Create table entity
	table, err := domain.NewTable(tableNumber, capacity, location)
	if err != nil {
		return nil, fmt.Errorf("failed to create table entity: %w", err)
	}

	// Generate ID
	table.ID = uuid.New().String()

	// Save to repository
	if err := uc.tableRepo.Create(ctx, table); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return table, nil
}

// GetTable retrieves a table by ID
func (uc *TableUseCase) GetTable(ctx context.Context, tableID string) (*domain.Table, error) {
	if tableID == "" {
		return nil, domain.ErrInvalidTableNumber
	}

	table, err := uc.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}

	return table, nil
}

// GetTableByNumber retrieves a table by table number
func (uc *TableUseCase) GetTableByNumber(ctx context.Context, tableNumber string) (*domain.Table, error) {
	if tableNumber == "" {
		return nil, domain.ErrInvalidTableNumber
	}

	table, err := uc.tableRepo.GetByTableNumber(ctx, tableNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get table by number: %w", err)
	}

	return table, nil
}

// UpdateTable updates table information
func (uc *TableUseCase) UpdateTable(ctx context.Context, tableID, tableNumber string, capacity int, location string) (*domain.Table, error) {
	// Get existing table
	table, err := uc.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}

	// Update fields
	if err := table.Update(tableNumber, capacity, location); err != nil {
		return nil, fmt.Errorf("failed to update table fields: %w", err)
	}

	// Save changes
	if err := uc.tableRepo.Update(ctx, table); err != nil {
		return nil, fmt.Errorf("failed to update table: %w", err)
	}

	return table, nil
}

// DeleteTable deletes a table
func (uc *TableUseCase) DeleteTable(ctx context.Context, tableID string) error {
	if tableID == "" {
		return domain.ErrInvalidTableNumber
	}

	// Check if table exists
	table, err := uc.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}

	// Don't allow deletion of occupied tables
	if table.IsOccupied() {
		return fmt.Errorf("cannot delete occupied table")
	}

	// Delete table
	if err := uc.tableRepo.Delete(ctx, tableID); err != nil {
		return fmt.Errorf("failed to delete table: %w", err)
	}

	return nil
}

// ListTables retrieves tables with filters and pagination
func (uc *TableUseCase) ListTables(ctx context.Context, page, pageSize int, status domain.TableStatus, location string) ([]*domain.Table, int, error) {
	filters := repository.ListFilters{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		Location: location,
	}

	tables, total, err := uc.tableRepo.List(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tables: %w", err)
	}

	return tables, total, nil
}

// UpdateTableStatus updates the status of a table
func (uc *TableUseCase) UpdateTableStatus(ctx context.Context, tableID string, status domain.TableStatus, orderID string) (*domain.Table, error) {
	if tableID == "" {
		return nil, domain.ErrInvalidTableNumber
	}

	// Get table
	table, err := uc.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}

	// Update status
	if err := table.UpdateStatus(status, orderID); err != nil {
		return nil, fmt.Errorf("failed to update table status: %w", err)
	}

	// Save changes
	if err := uc.tableRepo.Update(ctx, table); err != nil {
		return nil, fmt.Errorf("failed to save table: %w", err)
	}

	return table, nil
}

// GetAvailableTables retrieves all available tables
func (uc *TableUseCase) GetAvailableTables(ctx context.Context, minCapacity int, location string) ([]*domain.Table, error) {
	tables, err := uc.tableRepo.GetAvailableTables(ctx, minCapacity, location)
	if err != nil {
		return nil, fmt.Errorf("failed to get available tables: %w", err)
	}

	return tables, nil
}

// MarkTableAvailable marks a table as available
func (uc *TableUseCase) MarkTableAvailable(ctx context.Context, tableID string) error {
	table, err := uc.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}

	if err := table.MarkAvailable(); err != nil {
		return fmt.Errorf("failed to mark table as available: %w", err)
	}

	if err := uc.tableRepo.Update(ctx, table); err != nil {
		return fmt.Errorf("failed to save table: %w", err)
	}

	return nil
}

// MarkTableOccupied marks a table as occupied
func (uc *TableUseCase) MarkTableOccupied(ctx context.Context, tableID, orderID string) error {
	table, err := uc.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}

	if err := table.MarkOccupied(orderID); err != nil {
		return fmt.Errorf("failed to mark table as occupied: %w", err)
	}

	if err := uc.tableRepo.Update(ctx, table); err != nil {
		return fmt.Errorf("failed to save table: %w", err)
	}

	return nil
}

// MarkTableReserved marks a table as reserved
func (uc *TableUseCase) MarkTableReserved(ctx context.Context, tableID string) error {
	table, err := uc.tableRepo.GetByID(ctx, tableID)
	if err != nil {
		return fmt.Errorf("failed to get table: %w", err)
	}

	if err := table.MarkReserved(); err != nil {
		return fmt.Errorf("failed to mark table as reserved: %w", err)
	}

	if err := uc.tableRepo.Update(ctx, table); err != nil {
		return fmt.Errorf("failed to save table: %w", err)
	}

	return nil
}

// FindTableForParty finds a suitable table for a party of given size
func (uc *TableUseCase) FindTableForParty(ctx context.Context, numberOfGuests int, location string) (*domain.Table, error) {
	if numberOfGuests <= 0 {
		return nil, fmt.Errorf("number of guests must be positive")
	}

	// Get available tables sorted by capacity
	tables, err := uc.tableRepo.GetAvailableTables(ctx, numberOfGuests, location)
	if err != nil {
		return nil, fmt.Errorf("failed to get available tables: %w", err)
	}

	if len(tables) == 0 {
		return nil, fmt.Errorf("no suitable table available")
	}

	// Return the first table (smallest that fits)
	return tables[0], nil
}
