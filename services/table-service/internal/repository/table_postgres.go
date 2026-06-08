package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	"restaurant-management/services/table-service/internal/domain"
)

// PostgresTableRepository is a PostgreSQL implementation of TableRepository.
type PostgresTableRepository struct {
	db *sql.DB
}

// NewPostgresTableRepository creates a PostgreSQL-backed table repository.
func NewPostgresTableRepository(db *sql.DB) (*PostgresTableRepository, error) {
	repo := &PostgresTableRepository{db: db}
	if err := repo.ensureSchema(context.Background()); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *PostgresTableRepository) ensureSchema(ctx context.Context) error {
	const query = `
		CREATE TABLE IF NOT EXISTS restaurant_tables (
			table_id VARCHAR(36) PRIMARY KEY,
			table_number VARCHAR(50) NOT NULL UNIQUE,
			capacity INTEGER NOT NULL,
			status VARCHAR(32) NOT NULL,
			location VARCHAR(100) NOT NULL,
			current_order_id VARCHAR(36) NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			CONSTRAINT chk_restaurant_tables_capacity CHECK (capacity > 0 AND capacity <= 50)
		);

		CREATE INDEX IF NOT EXISTS idx_restaurant_tables_status ON restaurant_tables(status);
		CREATE INDEX IF NOT EXISTS idx_restaurant_tables_location ON restaurant_tables(location);
		CREATE INDEX IF NOT EXISTS idx_restaurant_tables_capacity ON restaurant_tables(capacity);
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresTableRepository) Create(ctx context.Context, table *domain.Table) error {
	if err := table.Validate(); err != nil {
		return err
	}

	const query = `
		INSERT INTO restaurant_tables (
			table_id,
			table_number,
			capacity,
			status,
			location,
			current_order_id,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		table.ID,
		table.TableNumber,
		table.Capacity,
		string(table.Status),
		table.Location,
		table.CurrentOrderID,
		table.CreatedAt,
		table.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrTableNumberAlreadyExists
		}
		return err
	}

	return nil
}

func (r *PostgresTableRepository) GetByID(ctx context.Context, id string) (*domain.Table, error) {
	const query = `
		SELECT table_id, table_number, capacity, status, location, current_order_id, created_at, updated_at
		FROM restaurant_tables
		WHERE table_id = $1
	`

	table := &domain.Table{}
	var status string
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&table.ID,
		&table.TableNumber,
		&table.Capacity,
		&status,
		&table.Location,
		&table.CurrentOrderID,
		&table.CreatedAt,
		&table.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTableNotFound
		}
		return nil, err
	}

	table.Status = domain.TableStatus(status)
	return table, nil
}

func (r *PostgresTableRepository) GetByTableNumber(ctx context.Context, tableNumber string) (*domain.Table, error) {
	const query = `
		SELECT table_id, table_number, capacity, status, location, current_order_id, created_at, updated_at
		FROM restaurant_tables
		WHERE table_number = $1
	`

	table := &domain.Table{}
	var status string
	if err := r.db.QueryRowContext(ctx, query, tableNumber).Scan(
		&table.ID,
		&table.TableNumber,
		&table.Capacity,
		&status,
		&table.Location,
		&table.CurrentOrderID,
		&table.CreatedAt,
		&table.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTableNotFound
		}
		return nil, err
	}

	table.Status = domain.TableStatus(status)
	return table, nil
}

func (r *PostgresTableRepository) Update(ctx context.Context, table *domain.Table) error {
	if err := table.Validate(); err != nil {
		return err
	}

	const query = `
		UPDATE restaurant_tables
		SET
			table_number = $2,
			capacity = $3,
			status = $4,
			location = $5,
			current_order_id = $6,
			updated_at = $7
		WHERE table_id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		table.ID,
		table.TableNumber,
		table.Capacity,
		string(table.Status),
		table.Location,
		table.CurrentOrderID,
		table.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrTableNumberAlreadyExists
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrTableNotFound
	}

	return nil
}

func (r *PostgresTableRepository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM restaurant_tables WHERE table_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrTableNotFound
	}

	return nil
}

func (r *PostgresTableRepository) List(ctx context.Context, filters ListFilters) ([]*domain.Table, int, error) {
	baseQuery := `
		SELECT table_id, table_number, capacity, status, location, current_order_id, created_at, updated_at
		FROM restaurant_tables
	`
	countQuery := `SELECT COUNT(1) FROM restaurant_tables`

	conditions := make([]string, 0, 2)
	args := make([]any, 0, 4)

	if filters.Status != "" {
		args = append(args, string(filters.Status))
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)))
	}
	if strings.TrimSpace(filters.Location) != "" {
		args = append(args, "%"+strings.TrimSpace(filters.Location)+"%")
		conditions = append(conditions, fmt.Sprintf("location ILIKE $%d", len(args)))
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	page := filters.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filters.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)
	baseQuery += fmt.Sprintf(" ORDER BY table_number ASC LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tables := make([]*domain.Table, 0)
	for rows.Next() {
		table, scanErr := scanTable(rows)
		if scanErr != nil {
			return nil, 0, scanErr
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return tables, total, nil
}

func (r *PostgresTableRepository) GetAvailableTables(ctx context.Context, minCapacity int, location string) ([]*domain.Table, error) {
	query := `
		SELECT table_id, table_number, capacity, status, location, current_order_id, created_at, updated_at
		FROM restaurant_tables
		WHERE status = $1
	`

	args := make([]any, 0, 3)
	args = append(args, string(domain.StatusAvailable))

	if minCapacity > 0 {
		args = append(args, minCapacity)
		query += fmt.Sprintf(" AND capacity >= $%d", len(args))
	}
	if strings.TrimSpace(location) != "" {
		args = append(args, "%"+strings.TrimSpace(location)+"%")
		query += fmt.Sprintf(" AND location ILIKE $%d", len(args))
	}

	query += " ORDER BY capacity ASC, table_number ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := make([]*domain.Table, 0)
	for rows.Next() {
		table, scanErr := scanTable(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (r *PostgresTableRepository) ExistsByTableNumber(ctx context.Context, tableNumber string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM restaurant_tables WHERE table_number = $1)`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, tableNumber).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func scanTable(scanner interface {
	Scan(dest ...any) error
}) (*domain.Table, error) {
	table := &domain.Table{}
	var status string

	if err := scanner.Scan(
		&table.ID,
		&table.TableNumber,
		&table.Capacity,
		&status,
		&table.Location,
		&table.CurrentOrderID,
		&table.CreatedAt,
		&table.UpdatedAt,
	); err != nil {
		return nil, err
	}

	table.Status = domain.TableStatus(status)
	return table, nil
}
