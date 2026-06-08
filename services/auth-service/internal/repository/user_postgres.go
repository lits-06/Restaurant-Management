package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"restaurant-management/services/auth-service/internal/domain"
)

// PostgresUserRepository is a PostgreSQL implementation of UserRepository.
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a PostgreSQL-backed user repository.
func NewPostgresUserRepository(db *sql.DB) (*PostgresUserRepository, error) {
	repo := &PostgresUserRepository{db: db}
	if err := repo.ensureSchema(context.Background()); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *PostgresUserRepository) ensureSchema(ctx context.Context) error {
	const query = `
		CREATE TABLE IF NOT EXISTS auth_users (
			user_id VARCHAR(36) PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	const query = `
		INSERT INTO auth_users (user_id, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.UserID,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	const query = `
		SELECT user_id, email, password, created_at, updated_at
		FROM auth_users
		WHERE user_id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.UserID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
		SELECT user_id, email, password, created_at, updated_at
		FROM auth_users
		WHERE email = $1
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	const query = `
		UPDATE auth_users
		SET email = $2, password = $3, updated_at = $4
		WHERE user_id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		user.UserID,
		user.Email,
		user.Password,
		user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrUserAlreadyExists
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, userID string) error {
	const query = `DELETE FROM auth_users WHERE user_id = $1`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *PostgresUserRepository) Exists(ctx context.Context, email string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM auth_users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
