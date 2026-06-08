package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"restaurant-management/services/table-service/internal/domain"
)

// PostgresReservationRepository is a PostgreSQL implementation of ReservationRepository.
type PostgresReservationRepository struct {
	db *sql.DB
}

// NewPostgresReservationRepository creates a PostgreSQL-backed reservation repository.
func NewPostgresReservationRepository(db *sql.DB) (*PostgresReservationRepository, error) {
	repo := &PostgresReservationRepository{db: db}
	if err := repo.ensureSchema(context.Background()); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *PostgresReservationRepository) ensureSchema(ctx context.Context) error {
	const query = `
		CREATE TABLE IF NOT EXISTS table_reservations (
			reservation_id VARCHAR(36) PRIMARY KEY,
			table_id VARCHAR(36) NOT NULL REFERENCES restaurant_tables(table_id) ON DELETE CASCADE,
			customer_name VARCHAR(100) NOT NULL,
			customer_phone VARCHAR(32) NOT NULL,
			notes TEXT NOT NULL DEFAULT '',
			status VARCHAR(32) NOT NULL,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			CONSTRAINT chk_table_reservations_time CHECK (end_time > start_time)
		);

		CREATE INDEX IF NOT EXISTS idx_table_reservations_table_id ON table_reservations(table_id);
		CREATE INDEX IF NOT EXISTS idx_table_reservations_status ON table_reservations(status);
		CREATE INDEX IF NOT EXISTS idx_table_reservations_start_time ON table_reservations(start_time);
		CREATE INDEX IF NOT EXISTS idx_table_reservations_end_time ON table_reservations(end_time);

		CREATE TABLE IF NOT EXISTS table_reservation_items (
			reservation_id VARCHAR(36) NOT NULL REFERENCES table_reservations(reservation_id) ON DELETE CASCADE,
			menu_item_id VARCHAR(36) NOT NULL,
			quantity INTEGER NOT NULL,
			note TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (reservation_id, menu_item_id)
		);
		CREATE INDEX IF NOT EXISTS idx_table_reservation_items_reservation_id ON table_reservation_items(reservation_id);
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresReservationRepository) CreateReservation(ctx context.Context, reservation *domain.Reservation) error {
	if err := reservation.Validate(); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	const insertReservation = `
		INSERT INTO table_reservations (
			reservation_id,
			table_id,
			customer_name,
			customer_phone,
			notes,
			status,
			start_time,
			end_time,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err = tx.ExecContext(ctx, insertReservation,
		reservation.ID,
		reservation.TableID,
		reservation.CustomerName,
		reservation.CustomerPhone,
		reservation.Notes,
		string(reservation.Status),
		reservation.StartTime,
		reservation.EndTime,
		reservation.CreatedAt,
		reservation.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return domain.ErrTableNotFound
		}
		return err
	}

	if len(reservation.Items) > 0 {
		const insertItem = `
			INSERT INTO table_reservation_items (
				reservation_id,
				menu_item_id,
				quantity,
				note
			)
			VALUES ($1, $2, $3, $4)
		`
		for _, item := range reservation.Items {
			_, err = tx.ExecContext(ctx, insertItem,
				reservation.ID,
				item.MenuItemID,
				item.Quantity,
				item.Note,
			)
			if err != nil {
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *PostgresReservationRepository) GetReservationByID(ctx context.Context, id string) (*domain.Reservation, error) {
	const query = `
		SELECT reservation_id, table_id, customer_name, customer_phone, notes, status, start_time, end_time, created_at, updated_at
		FROM table_reservations
		WHERE reservation_id = $1
	`

	reservation := &domain.Reservation{}
	var status string
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&reservation.ID,
		&reservation.TableID,
		&reservation.CustomerName,
		&reservation.CustomerPhone,
		&reservation.Notes,
		&status,
		&reservation.StartTime,
		&reservation.EndTime,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrReservationNotFound
		}
		return nil, err
	}

	reservation.Status = domain.ReservationStatus(status)

	items, err := r.getReservationItems(ctx, reservation.ID)
	if err != nil {
		return nil, err
	}
	reservation.Items = items
	return reservation, nil
}

func (r *PostgresReservationRepository) ListReservations(ctx context.Context, filters ReservationListFilters) ([]*domain.Reservation, int, error) {
	baseQuery := `
		SELECT reservation_id, table_id, customer_name, customer_phone, notes, status, start_time, end_time, created_at, updated_at
		FROM table_reservations
	`
	countQuery := `SELECT COUNT(1) FROM table_reservations`

	conditions := make([]string, 0, 4)
	args := make([]any, 0, 6)

	if strings.TrimSpace(filters.TableID) != "" {
		args = append(args, strings.TrimSpace(filters.TableID))
		conditions = append(conditions, fmt.Sprintf("table_id = $%d", len(args)))
	}
	if filters.Status != "" {
		args = append(args, string(filters.Status))
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)))
	}
	if !filters.FromTime.IsZero() {
		args = append(args, filters.FromTime)
		conditions = append(conditions, fmt.Sprintf("start_time >= $%d", len(args)))
	}
	if !filters.ToTime.IsZero() {
		args = append(args, filters.ToTime)
		conditions = append(conditions, fmt.Sprintf("end_time <= $%d", len(args)))
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
	baseQuery += fmt.Sprintf(" ORDER BY start_time ASC LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	reservations := make([]*domain.Reservation, 0)
	for rows.Next() {
		reservation, scanErr := scanReservation(rows)
		if scanErr != nil {
			return nil, 0, scanErr
		}
		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	for _, reservation := range reservations {
		items, itemsErr := r.getReservationItems(ctx, reservation.ID)
		if itemsErr != nil {
			return nil, 0, itemsErr
		}
		reservation.Items = items
	}

	return reservations, total, nil
}

func (r *PostgresReservationRepository) CancelReservation(ctx context.Context, id string) (*domain.Reservation, error) {
	reservation, err := r.GetReservationByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := reservation.Cancel(); err != nil {
		return nil, err
	}

	const query = `
		UPDATE table_reservations
		SET status = $2, updated_at = $3
		WHERE reservation_id = $1
	`

	_, err = r.db.ExecContext(ctx, query, reservation.ID, string(reservation.Status), reservation.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (r *PostgresReservationRepository) HasOverlappingReservation(ctx context.Context, tableID string, startTime, endTime time.Time) (bool, error) {
	const query = `
		SELECT EXISTS(
			SELECT 1
			FROM table_reservations
			WHERE table_id = $1
			  AND status = $2
			  AND NOT (end_time <= $3 OR start_time >= $4)
		)
	`

	var exists bool
	if err := r.db.QueryRowContext(ctx, query, tableID, string(domain.ReservationStatusReserved), startTime, endTime).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PostgresReservationRepository) getReservationItems(ctx context.Context, reservationID string) ([]domain.ReservationItem, error) {
	const query = `
		SELECT menu_item_id, quantity, note
		FROM table_reservation_items
		WHERE reservation_id = $1
		ORDER BY menu_item_id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, reservationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.ReservationItem, 0)
	for rows.Next() {
		var item domain.ReservationItem
		if scanErr := rows.Scan(&item.MenuItemID, &item.Quantity, &item.Note); scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func scanReservation(scanner interface {
	Scan(dest ...any) error
}) (*domain.Reservation, error) {
	reservation := &domain.Reservation{}
	var status string

	if err := scanner.Scan(
		&reservation.ID,
		&reservation.TableID,
		&reservation.CustomerName,
		&reservation.CustomerPhone,
		&reservation.Notes,
		&status,
		&reservation.StartTime,
		&reservation.EndTime,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
	); err != nil {
		return nil, err
	}

	reservation.Status = domain.ReservationStatus(status)
	return reservation, nil
}
