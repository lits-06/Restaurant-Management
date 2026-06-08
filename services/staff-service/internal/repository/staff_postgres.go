package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"restaurant-management/services/staff-service/internal/domain"

	"github.com/jackc/pgx/v5/pgconn"
)

// PostgresStaffRepository is a PostgreSQL implementation of StaffRepository.
type PostgresStaffRepository struct {
	db *sql.DB
}

// NewPostgresStaffRepository creates a PostgreSQL-backed staff repository.
func NewPostgresStaffRepository(db *sql.DB) (*PostgresStaffRepository, error) {
	repo := &PostgresStaffRepository{db: db}
	if err := repo.ensureSchema(context.Background()); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *PostgresStaffRepository) ensureSchema(ctx context.Context) error {
	const query = `
		CREATE EXTENSION IF NOT EXISTS pgcrypto;

		CREATE TABLE IF NOT EXISTS staff_members (
			staff_id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(120) NOT NULL,
			role VARCHAR(100) NOT NULL,
			contact VARCHAR(100) NOT NULL,
			avatar TEXT NOT NULL DEFAULT ''
		);

		CREATE INDEX IF NOT EXISTS idx_staff_members_name ON staff_members(name);
		CREATE INDEX IF NOT EXISTS idx_staff_members_role ON staff_members(role);
		CREATE INDEX IF NOT EXISTS idx_staff_members_contact ON staff_members(contact);
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

// Create inserts a new staff member.
func (r *PostgresStaffRepository) Create(ctx context.Context, staff *domain.Staff) error {
	if err := staff.Validate(); err != nil {
		return err
	}

	const query = `
		INSERT INTO staff_members (
			name,
			role,
			contact,
			avatar
		)
		VALUES ($1, $2, $3, $4)
		RETURNING staff_id
	`

	err := r.db.QueryRowContext(ctx, query,
		staff.Name,
		staff.Role,
		staff.Contact,
		staff.Avatar,
	).Scan(&staff.StaffID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrStaffNotFound
		}
		return err
	}

	return nil
}

// GetByID retrieves a staff member by ID.
func (r *PostgresStaffRepository) GetByID(ctx context.Context, staffID string) (*domain.Staff, error) {
	const query = `
		SELECT staff_id, name, role, contact, avatar,
		FROM staff_members
		WHERE staff_id = $1
	`

	staff := &domain.Staff{}
	if err := r.db.QueryRowContext(ctx, query, staffID).Scan(
		&staff.StaffID,
		&staff.Name,
		&staff.Role,
		&staff.Contact,
		&staff.Avatar,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrStaffNotFound
		}
		return nil, err
	}

	return staff, nil
}

// Update updates an existing staff member.
func (r *PostgresStaffRepository) Update(ctx context.Context, staff *domain.Staff) error {
	if err := staff.Validate(); err != nil {
		return err
	}

	const query = `
		UPDATE staff_members
		SET
			name = $2,
			role = $3,
			contact = $4,
			avatar = $5
		WHERE staff_id = $1
		RETURNING staff_id
	`

	if err := r.db.QueryRowContext(ctx, query,
		staff.StaffID,
		staff.Name,
		staff.Role,
		staff.Contact,
		staff.Avatar,
	).Scan(&staff.StaffID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrStaffNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrStaffNotFound
		}
		return err
	}

	return nil
}

// Delete deletes a staff member by ID.
func (r *PostgresStaffRepository) Delete(ctx context.Context, staffID string) error {
	const query = `DELETE FROM staff_members WHERE staff_id = $1`

	result, err := r.db.ExecContext(ctx, query, staffID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrStaffNotFound
	}

	return nil
}

// List retrieves staff members with pagination and keyword search.
func (r *PostgresStaffRepository) List(ctx context.Context, page, pageSize int, keyword string) ([]*domain.Staff, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	baseQuery := `
		SELECT staff_id, name, role, contact, avatar
		FROM staff_members
	`
	countQuery := `SELECT COUNT(1) FROM staff_members`

	conditions := make([]string, 0, 1)
	args := make([]interface{}, 0, 3)

	if keyword != "" {
		args = append(args, "%"+keyword+"%")
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR role ILIKE $%d OR contact ILIKE $%d)", len(args), len(args), len(args)))
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

	staffMembers := make([]*domain.Staff, 0)
	for rows.Next() {
		staff, scanErr := scanStaff(rows)
		if scanErr != nil {
			return nil, 0, scanErr
		}
		staffMembers = append(staffMembers, staff)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return staffMembers, total, nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanStaff(s scanner) (*domain.Staff, error) {
	staff := &domain.Staff{}
	if err := s.Scan(
		&staff.StaffID,
		&staff.Name,
		&staff.Role,
		&staff.Contact,
		&staff.Avatar,
	); err != nil {
		return nil, err
	}

	return staff, nil
}
