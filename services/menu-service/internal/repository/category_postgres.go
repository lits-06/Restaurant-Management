package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"restaurant-management/services/menu-service/internal/domain"
)

// PostgresCategoryRepository is a PostgreSQL implementation of CategoryRepository.
type PostgresCategoryRepository struct {
	db *sql.DB
}

// NewPostgresCategoryRepository creates a PostgreSQL-backed category repository.
func NewPostgresCategoryRepository(db *sql.DB) (*PostgresCategoryRepository, error) {
	repo := &PostgresCategoryRepository{db: db}
	if err := repo.ensureSchema(context.Background()); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *PostgresCategoryRepository) ensureSchema(ctx context.Context) error {
	const query = `
		CREATE TABLE IF NOT EXISTS menu_categories (
			category_id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(100) NOT NULL UNIQUE
		);
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	if err := category.Validate(); err != nil {
		return err
	}

	const query = `
		INSERT INTO menu_categories (name)
		VALUES ($1)
	`

	_, err := r.db.ExecContext(ctx, query, category.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrCategoryAlreadyExists
		}
		return err
	}

	return nil
}

func (r *PostgresCategoryRepository) GetByID(ctx context.Context, categoryID string) (*domain.Category, error) {
	const query = `
		SELECT category_id, name
		FROM menu_categories
		WHERE category_id = $1
	`

	category := &domain.Category{}
	err := r.db.QueryRowContext(ctx, query, categoryID).Scan(&category.CategoryID, &category.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}

	return category, nil
}

func (r *PostgresCategoryRepository) GetByName(ctx context.Context, name string) (*domain.Category, error) {
	const query = `
		SELECT category_id, name
		FROM menu_categories
		WHERE name = $1
	`

	category := &domain.Category{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(&category.CategoryID, &category.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrCategoryNotFound
		}
		return nil, err
	}

	return category, nil
}

func (r *PostgresCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	if err := category.Validate(); err != nil {
		return err
	}

	const query = `
		UPDATE menu_categories
		SET name = $2
		WHERE category_id = $1
	`

	result, err := r.db.ExecContext(ctx, query, category.CategoryID, category.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrCategoryAlreadyExists
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrCategoryNotFound
	}

	return nil
}

func (r *PostgresCategoryRepository) Delete(ctx context.Context, categoryID string) error {
	const query = `DELETE FROM menu_categories WHERE category_id = $1`

	result, err := r.db.ExecContext(ctx, query, categoryID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrCategoryNotFound
	}

	return nil
}

// func (r *PostgresCategoryRepository) List(ctx context.Context, page, pageSize int) ([]*domain.Category, int, error) {
// 	const countQuery = `SELECT COUNT(1) FROM menu_categories`

// 	var total int
// 	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
// 		return nil, 0, err
// 	}

// 	offset := (page - 1) * pageSize
// 	if offset < 0 {
// 		offset = 0
// 	}

// 	const listQuery = `
// 		SELECT category_id, name, description, display_order, created_at, updated_at
// 		FROM menu_categories
// 		ORDER BY display_order ASC, name ASC
// 		LIMIT $1 OFFSET $2
// 	`

// 	rows, err := r.db.QueryContext(ctx, listQuery, pageSize, offset)
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	defer rows.Close()

// 	categories := make([]*domain.Category, 0)
// 	for rows.Next() {
// 		category := &domain.Category{}
// 		if scanErr := rows.Scan(
// 			&category.CategoryID,
// 			&category.Name,
// 			&category.Description,
// 			&category.DisplayOrder,
// 			&category.CreatedAt,
// 			&category.UpdatedAt,
// 		); scanErr != nil {
// 			return nil, 0, scanErr
// 		}
// 		categories = append(categories, category)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, 0, err
// 	}

// 	return categories, total, nil
// }

func (r *PostgresCategoryRepository) ListAll(ctx context.Context) ([]*domain.Category, error) {
	const query = `
		SELECT category_id, name
		FROM menu_categories
		ORDER BY name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*domain.Category, 0)
	for rows.Next() {
		category := &domain.Category{}
		if scanErr := rows.Scan(&category.CategoryID, &category.Name); scanErr != nil {
			return nil, scanErr
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
