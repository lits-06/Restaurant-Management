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
			table_id     VARCHAR(36)  PRIMARY KEY DEFAULT gen_random_uuid(),
			table_number INTEGER      NOT NULL UNIQUE,
			capacity     INTEGER      NOT NULL,
			status       VARCHAR(32)  NOT NULL DEFAULT 'AVAILABLE',
			created_at   TIMESTAMP    NOT NULL,
			updated_at   TIMESTAMP    NOT NULL,
			CONSTRAINT chk_table_number  CHECK (table_number > 0),
			CONSTRAINT chk_capacity      CHECK (capacity > 0 AND capacity <= 50)
		);

		-- drop legacy column if migrating from older schema
		ALTER TABLE restaurant_tables DROP COLUMN IF EXISTS location;
		ALTER TABLE restaurant_tables DROP COLUMN IF EXISTS current_order_id;

		CREATE INDEX IF NOT EXISTS idx_restaurant_tables_status   ON restaurant_tables(status);
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
		INSERT INTO restaurant_tables (table_id, table_number, capacity, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		table.ID, table.TableNumber, table.Capacity, string(table.Status),
		table.CreatedAt, table.UpdatedAt,
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
		SELECT table_id, table_number, capacity, status, created_at, updated_at
		FROM restaurant_tables WHERE table_id = $1
	`
	return r.scanOne(r.db.QueryRowContext(ctx, query, id))
}

func (r *PostgresTableRepository) GetByTableNumber(ctx context.Context, tableNumber int) (*domain.Table, error) {
	const query = `
		SELECT table_id, table_number, capacity, status, created_at, updated_at
		FROM restaurant_tables WHERE table_number = $1
	`
	return r.scanOne(r.db.QueryRowContext(ctx, query, tableNumber))
}

func (r *PostgresTableRepository) Update(ctx context.Context, table *domain.Table) error {
	if err := table.Validate(); err != nil {
		return err
	}

	const query = `
		UPDATE restaurant_tables
		SET table_number = $2, capacity = $3, status = $4, updated_at = $5
		WHERE table_id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		table.ID, table.TableNumber, table.Capacity, string(table.Status), table.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrTableNumberAlreadyExists
		}
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrTableNotFound
	}
	return nil
}

func (r *PostgresTableRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM restaurant_tables WHERE table_id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrTableNotFound
	}
	return nil
}

func (r *PostgresTableRepository) List(ctx context.Context, filters ListFilters) ([]*domain.Table, int, error) {
	conditions := make([]string, 0, 1)
	args := make([]any, 0, 3)

	if filters.Status != "" {
		args = append(args, string(filters.Status))
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM restaurant_tables`+whereClause, args...).Scan(&total); err != nil {
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

	query := `SELECT table_id, table_number, capacity, status, created_at, updated_at FROM restaurant_tables` +
		whereClause +
		fmt.Sprintf(" ORDER BY table_number ASC LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	tables := make([]*domain.Table, 0)
	for rows.Next() {
		t, err := scanTableRow(rows)
		if err != nil {
			return nil, 0, err
		}
		tables = append(tables, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return tables, total, nil
}

func (r *PostgresTableRepository) GetAvailableTables(ctx context.Context, minCapacity int) ([]*domain.Table, error) {
	query := `
		SELECT table_id, table_number, capacity, status, created_at, updated_at
		FROM restaurant_tables
		WHERE status = $1
	`
	args := []any{string(domain.StatusAvailable)}

	if minCapacity > 0 {
		args = append(args, minCapacity)
		query += fmt.Sprintf(" AND capacity >= $%d", len(args))
	}
	query += " ORDER BY capacity ASC, table_number ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := make([]*domain.Table, 0)
	for rows.Next() {
		t, err := scanTableRow(rows)
		if err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

func (r *PostgresTableRepository) ExistsByTableNumber(ctx context.Context, tableNumber int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM restaurant_tables WHERE table_number = $1)`, tableNumber,
	).Scan(&exists)
	return exists, err
}

// scanOne wraps QueryRowContext result into a Table, mapping ErrNoRows → ErrTableNotFound.
func (r *PostgresTableRepository) scanOne(row *sql.Row) (*domain.Table, error) {
	t := &domain.Table{}
	var status string
	if err := row.Scan(&t.ID, &t.TableNumber, &t.Capacity, &status, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTableNotFound
		}
		return nil, err
	}
	t.Status = domain.TableStatus(status)
	return t, nil
}

func scanTableRow(scanner interface {
	Scan(dest ...any) error
}) (*domain.Table, error) {
	t := &domain.Table{}
	var status string
	if err := scanner.Scan(&t.ID, &t.TableNumber, &t.Capacity, &status, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}
	t.Status = domain.TableStatus(status)
	return t, nil
}
