package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	"restaurant-management/services/menu-service/internal/domain"
)

// PostgresMenuItemRepository is a PostgreSQL implementation of MenuItemRepository.
type PostgresMenuItemRepository struct {
	db *sql.DB
}

// NewPostgresMenuItemRepository creates a PostgreSQL-backed menu item repository.
func NewPostgresMenuItemRepository(db *sql.DB) (*PostgresMenuItemRepository, error) {
	repo := &PostgresMenuItemRepository{db: db}
	if err := repo.ensureSchema(context.Background()); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *PostgresMenuItemRepository) ensureSchema(ctx context.Context) error {
	const query = `
		CREATE TABLE IF NOT EXISTS menu_items (
			item_id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(100) NOT NULL UNIQUE,
			description TEXT NOT NULL DEFAULT '',
			price DOUBLE PRECISION NOT NULL,
			category_id VARCHAR(36) NOT NULL REFERENCES menu_categories(category_id) ON DELETE CASCADE,
			image_url TEXT NOT NULL DEFAULT ''
		);

		CREATE INDEX IF NOT EXISTS idx_menu_items_category_id ON menu_items(category_id);
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresMenuItemRepository) Create(ctx context.Context, item *domain.MenuItem) error {
	if err := item.Validate(); err != nil {
		return err
	}

	const query = `
		INSERT INTO menu_items (name, description, price, category_id, image_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING item_id
	`

	err := r.db.QueryRowContext(ctx, query,
		item.Name,
		item.Description,
		item.Price,
		item.CategoryID,
		item.ImageURL,
	).Scan(&item.ItemID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.ErrMenuItemAlreadyExists
			}
			if pgErr.Code == "23503" {
				return domain.ErrCategoryNotFound
			}
		}
		return err
	}

	return nil
}

func (r *PostgresMenuItemRepository) GetByID(ctx context.Context, itemID string) (*domain.MenuItem, error) {
	const query = `
		SELECT item_id, name, description, price, category_id, image_url
		FROM menu_items
		WHERE item_id = $1
	`

	item, err := r.querySingle(ctx, query, itemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrMenuItemNotFound
		}
		return nil, err
	}

	return item, nil
}

func (r *PostgresMenuItemRepository) GetByName(ctx context.Context, name string) (*domain.MenuItem, error) {
	const query = `
		SELECT item_id, name, description, price, category_id, image_url
		FROM menu_items
		WHERE name = $1
	`

	item, err := r.querySingle(ctx, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrMenuItemNotFound
		}
		return nil, err
	}

	return item, nil
}

func (r *PostgresMenuItemRepository) Update(ctx context.Context, item *domain.MenuItem) error {
	if err := item.Validate(); err != nil {
		return err
	}

	const query = `
		UPDATE menu_items
		SET
			name = $2,
			description = $3,
			price = $4,
			category_id = $5,
			image_url = $6
		WHERE item_id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		item.ItemID,
		item.Name,
		item.Description,
		item.Price,
		item.CategoryID,
		item.ImageURL,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.ErrMenuItemAlreadyExists
			}
			if pgErr.Code == "23503" {
				return domain.ErrCategoryNotFound
			}
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrMenuItemNotFound
	}

	return nil
}

func (r *PostgresMenuItemRepository) Delete(ctx context.Context, itemID string) error {
	const query = `DELETE FROM menu_items WHERE item_id = $1`

	result, err := r.db.ExecContext(ctx, query, itemID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrMenuItemNotFound
	}

	return nil
}

func (r *PostgresMenuItemRepository) List(ctx context.Context, page, pageSize int, categoryID string, keyword string) ([]*domain.MenuItem, int, error) {
	baseQuery := `
		SELECT item_id, name, description, price, category_id, image_url
		FROM menu_items
	`
	countQuery := `SELECT COUNT(1) FROM menu_items`

	conditions := make([]string, 0, 2)
	args := make([]interface{}, 0, 4)

	if categoryID != "" {
		args = append(args, categoryID)
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", len(args)))
	}

	if keyword != "" {
		args = append(args, "%"+keyword+"%")
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", len(args)))
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

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	args = append(args, pageSize, offset)
	baseQuery += fmt.Sprintf(" ORDER BY name ASC LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]*domain.MenuItem, 0)
	for rows.Next() {
		item, scanErr := scanMenuItem(rows)
		if scanErr != nil {
			return nil, 0, scanErr
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *PostgresMenuItemRepository) ListByCategory(ctx context.Context, categoryID string) ([]*domain.MenuItem, error) {
	const query = `
		SELECT item_id, name, description, price, category_id, image_url
		FROM menu_items
		WHERE category_id = $1
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.MenuItem, 0)
	for rows.Next() {
		item, scanErr := scanMenuItem(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *PostgresMenuItemRepository) querySingle(ctx context.Context, query string, arg interface{}) (*domain.MenuItem, error) {
	row := r.db.QueryRowContext(ctx, query, arg)
	return scanMenuItem(row)
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanMenuItem(s scanner) (*domain.MenuItem, error) {
	item := &domain.MenuItem{}
	err := s.Scan(
		&item.ItemID,
		&item.Name,
		&item.Description,
		&item.Price,
		&item.CategoryID,
		&item.ImageURL,
	)
	if err != nil {
		return nil, err
	}

	return item, nil
}
