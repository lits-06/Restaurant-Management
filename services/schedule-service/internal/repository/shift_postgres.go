package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"restaurant-management/services/schedule-service/internal/domain"
)

type PostgresShiftRepository struct {
	db *sql.DB
}

func NewPostgresShiftRepository(db *sql.DB) (*PostgresShiftRepository, error) {
	r := &PostgresShiftRepository{db: db}
	if err := r.ensureSchema(context.Background()); err != nil {
		return nil, fmt.Errorf("ensureSchema: %w", err)
	}
	return r, nil
}

func (r *PostgresShiftRepository) ensureSchema(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS shifts (
			shift_id   VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id    VARCHAR(36) NOT NULL,
			date       DATE        NOT NULL,
			start_time VARCHAR(5)  NOT NULL,
			end_time   VARCHAR(5)  NOT NULL,
			role       VARCHAR(32) NOT NULL,
			notes      TEXT        NOT NULL DEFAULT '',
			created_by VARCHAR(36) NOT NULL DEFAULT '',
			created_at TIMESTAMP   NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP   NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_shifts_user_id ON shifts(user_id);
		CREATE INDEX IF NOT EXISTS idx_shifts_date    ON shifts(date);
	`)
	return err
}

func (r *PostgresShiftRepository) Create(ctx context.Context, s *domain.Shift) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO shifts (user_id, date, start_time, end_time, role, notes, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING shift_id`,
		s.UserID, s.Date, s.StartTime, s.EndTime, s.Role, s.Notes, s.CreatedBy, s.CreatedAt, s.UpdatedAt,
	).Scan(&s.ShiftID)
}

func (r *PostgresShiftRepository) GetByID(ctx context.Context, id string) (*domain.Shift, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT shift_id, user_id, date::text, start_time, end_time, role, notes, created_by, created_at, updated_at
		FROM shifts WHERE shift_id = $1`, id)
	s, err := scanShift(row)
	if err == sql.ErrNoRows {
		return nil, domain.ErrShiftNotFound
	}
	return s, err
}

func (r *PostgresShiftRepository) Update(ctx context.Context, s *domain.Shift) error {
	s.UpdatedAt = time.Now()
	res, err := r.db.ExecContext(ctx, `
		UPDATE shifts SET date = $1, start_time = $2, end_time = $3, notes = $4, updated_at = $5
		WHERE shift_id = $6`,
		s.Date, s.StartTime, s.EndTime, s.Notes, s.UpdatedAt, s.ShiftID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrShiftNotFound
	}
	return nil
}

func (r *PostgresShiftRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM shifts WHERE shift_id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrShiftNotFound
	}
	return nil
}

func (r *PostgresShiftRepository) List(ctx context.Context, f ListFilters) ([]*domain.Shift, int, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 {
		f.PageSize = 50
	}

	args := []any{}
	conds := []string{}
	i := 1

	if f.Month != "" && len(f.Month) == 7 {
		// "YYYY-MM" → date range
		conds = append(conds, fmt.Sprintf("date >= $%d::date AND date < ($%d::date + INTERVAL '1 month')", i, i))
		args = append(args, f.Month+"-01")
		i++
	}
	if f.UserID != "" {
		conds = append(conds, fmt.Sprintf("user_id = $%d", i))
		args = append(args, f.UserID)
		i++
	}
	if f.Role != "" {
		conds = append(conds, fmt.Sprintf("role = $%d", i))
		args = append(args, f.Role)
		i++
	}

	where := ""
	if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
	}

	var total int
	countArgs := make([]any, len(args))
	copy(countArgs, args)
	if err := r.db.QueryRowContext(ctx,
		fmt.Sprintf("SELECT COUNT(*) FROM shifts %s", where), countArgs...,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.PageSize
	args = append(args, f.PageSize, offset)
	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf(`
			SELECT shift_id, user_id, date::text, start_time, end_time, role, notes, created_by, created_at, updated_at
			FROM shifts %s
			ORDER BY date ASC, start_time ASC
			LIMIT $%d OFFSET $%d`, where, i, i+1),
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var shifts []*domain.Shift
	for rows.Next() {
		s, err := scanShiftRow(rows)
		if err != nil {
			return nil, 0, err
		}
		shifts = append(shifts, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return shifts, total, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanShift(s scanner) (*domain.Shift, error) {
	var sh domain.Shift
	if err := s.Scan(
		&sh.ShiftID, &sh.UserID, &sh.Date, &sh.StartTime, &sh.EndTime,
		&sh.Role, &sh.Notes, &sh.CreatedBy, &sh.CreatedAt, &sh.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &sh, nil
}

func scanShiftRow(rows *sql.Rows) (*domain.Shift, error) {
	var sh domain.Shift
	if err := rows.Scan(
		&sh.ShiftID, &sh.UserID, &sh.Date, &sh.StartTime, &sh.EndTime,
		&sh.Role, &sh.Notes, &sh.CreatedBy, &sh.CreatedAt, &sh.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &sh, nil
}
