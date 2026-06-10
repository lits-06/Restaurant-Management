package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"restaurant-management/services/user-service/internal/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) (*PostgresUserRepository, error) {
	r := &PostgresUserRepository{db: db}
	if err := r.ensureSchema(context.Background()); err != nil {
		return nil, fmt.Errorf("ensureSchema: %w", err)
	}
	return r, nil
}

func (r *PostgresUserRepository) ensureSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			user_id    VARCHAR(36)  PRIMARY KEY,
			email      VARCHAR(255) NOT NULL UNIQUE,
			username   VARCHAR(100) NOT NULL UNIQUE,
			full_name  VARCHAR(255) NOT NULL,
			phone      VARCHAR(50)  NOT NULL DEFAULT '',
			password   VARCHAR(255) NOT NULL,
			status     VARCHAR(32)  NOT NULL DEFAULT 'ACTIVE',
			created_at TIMESTAMP    NOT NULL,
			updated_at TIMESTAMP    NOT NULL
		);
		ALTER TABLE users DROP COLUMN IF EXISTS roles;
		CREATE TABLE IF NOT EXISTS user_roles (
			user_id VARCHAR(36) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
			role    VARCHAR(32) NOT NULL,
			PRIMARY KEY (user_id, role)
		);
		CREATE INDEX IF NOT EXISTS idx_users_status    ON users(status);
		CREATE INDEX IF NOT EXISTS idx_users_username  ON users(username);
		CREATE INDEX IF NOT EXISTS idx_user_roles_role ON user_roles(role);
	`)
	return err
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (user_id, email, username, full_name, phone, password, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		user.ID, user.Email, user.Username, user.FullName, user.Phone,
		user.Password, string(user.Status), user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, role := range user.Roles {
		if _, err = tx.ExecContext(ctx,
			`INSERT INTO user_roles (user_id, role) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
			user.ID, string(role),
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT user_id, email, username, full_name, phone, password, status, created_at, updated_at
		FROM users WHERE user_id = $1`, id)
	u, err := scanUser(row)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	u.Roles, err = r.getRoles(ctx, u.ID)
	return u, err
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT user_id, email, username, full_name, phone, password, status, created_at, updated_at
		FROM users WHERE email = $1`, email)
	u, err := scanUser(row)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	u.Roles, err = r.getRoles(ctx, u.ID)
	return u, err
}

func (r *PostgresUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT user_id, email, username, full_name, phone, password, status, created_at, updated_at
		FROM users WHERE username = $1`, username)
	u, err := scanUser(row)
	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	u.Roles, err = r.getRoles(ctx, u.ID)
	return u, err
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE users
		SET email = $1, username = $2, full_name = $3, phone = $4,
		    password = $5, status = $6, updated_at = $7
		WHERE user_id = $8`,
		user.Email, user.Username, user.FullName, user.Phone,
		user.Password, string(user.Status), user.UpdatedAt, user.ID,
	)
	if err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM user_roles WHERE user_id = $1`, user.ID); err != nil {
		return err
	}
	for _, role := range user.Roles {
		if _, err = tx.ExecContext(ctx,
			`INSERT INTO user_roles (user_id, role) VALUES ($1, $2)`,
			user.ID, string(role),
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE user_id = $1`, id)
	return err
}

func (r *PostgresUserRepository) List(ctx context.Context, f ListFilters) ([]*domain.User, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 20
	}

	args := []any{}
	conds := []string{}
	i := 1

	if f.Status != "" {
		conds = append(conds, fmt.Sprintf("u.status = $%d", i))
		args = append(args, string(f.Status))
		i++
	}
	if f.Role != "" {
		conds = append(conds, fmt.Sprintf(
			"EXISTS (SELECT 1 FROM user_roles r2 WHERE r2.user_id = u.user_id AND r2.role = $%d)", i,
		))
		args = append(args, string(f.Role))
		i++
	}
	if f.Keyword != "" {
		kw := "%" + f.Keyword + "%"
		conds = append(conds, fmt.Sprintf(
			"(u.email ILIKE $%d OR u.username ILIKE $%d OR u.full_name ILIKE $%d)", i, i+1, i+2,
		))
		args = append(args, kw, kw, kw)
		i += 3
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	var total int
	countArgs := make([]any, len(args))
	copy(countArgs, args)
	if err := r.db.QueryRowContext(ctx,
		fmt.Sprintf("SELECT COUNT(*) FROM users u %s", where), countArgs...,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.PageSize
	args = append(args, f.PageSize, offset)
	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf(`
			SELECT u.user_id, u.email, u.username, u.full_name, u.phone, u.password,
			       u.status, u.created_at, u.updated_at
			FROM users u %s
			ORDER BY u.created_at DESC LIMIT $%d OFFSET $%d`, where, i, i+1),
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u, err := scanUserRow(rows)
		if err != nil {
			return nil, 0, err
		}
		u.Roles, err = r.getRoles(ctx, u.ID)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *PostgresUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email,
	).Scan(&exists)
	return exists, err
}

func (r *PostgresUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, username,
	).Scan(&exists)
	return exists, err
}

// ── helpers ───────────────────────────────────────────────────

func (r *PostgresUserRepository) getRoles(ctx context.Context, userID string) ([]domain.UserRole, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT role FROM user_roles WHERE user_id = $1 ORDER BY role`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.UserRole
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, domain.UserRole(role))
	}
	return roles, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(s scanner) (*domain.User, error) {
	var u domain.User
	var statusStr string
	if err := s.Scan(
		&u.ID, &u.Email, &u.Username, &u.FullName, &u.Phone,
		&u.Password, &statusStr, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	u.Status = domain.UserStatus(statusStr)
	return &u, nil
}

func scanUserRow(rows *sql.Rows) (*domain.User, error) {
	var u domain.User
	var statusStr string
	if err := rows.Scan(
		&u.ID, &u.Email, &u.Username, &u.FullName, &u.Phone,
		&u.Password, &statusStr, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	u.Status = domain.UserStatus(statusStr)
	return &u, nil
}
