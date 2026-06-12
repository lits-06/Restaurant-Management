package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"restaurant-management/services/order-service/internal/domain"
)

type sqlQueryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// PostgresOrderRepository is a PostgreSQL implementation of OrderRepository.
type PostgresOrderRepository struct {
	db *sql.DB
}

// NewPostgresOrderRepository creates a PostgreSQL-backed order repository.
func NewPostgresOrderRepository(db *sql.DB) (*PostgresOrderRepository, error) {
	repo := &PostgresOrderRepository{db: db}
	if err := repo.ensureSchema(context.Background()); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *PostgresOrderRepository) ensureSchema(ctx context.Context) error {
	const query = `
		CREATE TABLE IF NOT EXISTS orders (
			order_id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
			table_id VARCHAR(36) NOT NULL DEFAULT '',
			user_id VARCHAR(36) NOT NULL DEFAULT '',
			name VARCHAR(120) NOT NULL,
			phone VARCHAR(120) NOT NULL DEFAULT '',
			notes TEXT NOT NULL DEFAULT '',
			time TIMESTAMP NOT NULL,
			end_time TIMESTAMP,
			party_size INTEGER NOT NULL,
			status VARCHAR(32) NOT NULL,
			total DOUBLE PRECISION NOT NULL DEFAULT 0
		);

		ALTER TABLE orders ADD COLUMN IF NOT EXISTS table_id VARCHAR(36) NOT NULL DEFAULT '';
		ALTER TABLE orders ADD COLUMN IF NOT EXISTS notes TEXT NOT NULL DEFAULT '';
		ALTER TABLE orders ADD COLUMN IF NOT EXISTS end_time TIMESTAMP;
		ALTER TABLE orders ADD COLUMN IF NOT EXISTS user_id VARCHAR(36) NOT NULL DEFAULT '';

		CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);

		CREATE TABLE IF NOT EXISTS order_items (
			item_id     VARCHAR(36) NOT NULL REFERENCES menu_items(item_id),
			order_id    VARCHAR(36) NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
			name        VARCHAR(200) NOT NULL,
			price       DOUBLE PRECISION NOT NULL,
			quantity    INTEGER NOT NULL,
			item_status VARCHAR(16) NOT NULL DEFAULT 'PENDING',
			PRIMARY KEY (order_id, item_id)
		);

		ALTER TABLE order_items ADD COLUMN IF NOT EXISTS item_status VARCHAR(16) NOT NULL DEFAULT 'PENDING';

		-- Migrate: fix order_items PK from single item_id to composite (order_id, item_id)
		DO $$
		DECLARE col_count INT;
		BEGIN
			SELECT COUNT(*) INTO col_count
			FROM information_schema.key_column_usage
			WHERE table_name = 'order_items' AND constraint_name = 'order_items_pkey';
			IF col_count = 1 THEN
				ALTER TABLE order_items DROP CONSTRAINT order_items_pkey;
				ALTER TABLE order_items ADD PRIMARY KEY (order_id, item_id);
			END IF;
		END $$;

		CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
		CREATE INDEX IF NOT EXISTS idx_orders_name ON orders(name);
		CREATE INDEX IF NOT EXISTS idx_orders_table_id ON orders(table_id);
		CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgresOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	if order.Status == "" {
		order.Status = domain.StatusPending
	}

	if err := order.Validate(); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const insertOrderQuery = `
		INSERT INTO orders (
			table_id, user_id, name, phone, notes, time, end_time, party_size, status, total
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING order_id
	`

	var endTime *time.Time
	if !order.EndTime.IsZero() {
		endTime = &order.EndTime
	}

	if err := tx.QueryRowContext(ctx, insertOrderQuery,
		order.TableID,
		order.UserID,
		order.Name,
		order.Phone,
		order.Notes,
		order.Time,
		endTime,
		order.PartySize,
		string(order.Status),
		order.Total,
	).Scan(&order.OrderID); err != nil {
		return err
	}

	if err := r.insertOrderItems(ctx, tx, order.OrderID, order.Items); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresOrderRepository) GetByID(ctx context.Context, orderID string) (*domain.Order, error) {
	const query = `
		SELECT order_id, table_id, user_id, name, phone, notes, time, end_time, party_size, status, total
		FROM orders
		WHERE order_id = $1
	`

	order := &domain.Order{}
	var endTime sql.NullTime
	if err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&order.OrderID,
		&order.TableID,
		&order.UserID,
		&order.Name,
		&order.Phone,
		&order.Notes,
		&order.Time,
		&endTime,
		&order.PartySize,
		&order.Status,
		&order.Total,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}
	if endTime.Valid {
		order.EndTime = endTime.Time
	}

	items, err := r.getOrderItems(ctx, r.db, order.OrderID)
	if err != nil {
		return nil, err
	}
	order.Items = items
	return order, nil
}

func (r *PostgresOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	if order.Status == "" {
		order.Status = domain.StatusPending
	}
	if err := order.Validate(); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const updateOrderQuery = `
		UPDATE orders
		SET table_id = $2, name = $3, phone = $4, notes = $5, time = $6, end_time = $7,
		    party_size = $8, status = $9, total = $10
		WHERE order_id = $1
	`
	// user_id is intentionally not updated — it is set at creation and is immutable

	var endTime *time.Time
	if !order.EndTime.IsZero() {
		endTime = &order.EndTime
	}

	result, err := tx.ExecContext(ctx, updateOrderQuery,
		order.OrderID,
		order.TableID,
		order.Name,
		order.Phone,
		order.Notes,
		order.Time,
		endTime,
		order.PartySize,
		string(order.Status),
		order.Total,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrOrderNotFound
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM order_items WHERE order_id = $1`, order.OrderID); err != nil {
		return err
	}
	if err := r.insertOrderItems(ctx, tx, order.OrderID, order.Items); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresOrderRepository) Delete(ctx context.Context, orderID string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM orders WHERE order_id = $1`, orderID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
}

func (r *PostgresOrderRepository) List(ctx context.Context, page, pageSize int, status domain.OrderStatus, keyword, userID, sortOrder string) ([]*domain.Order, int, error) {
	clauses := make([]string, 0, 3)
	args := make([]any, 0, 4)

	if status != "" {
		args = append(args, string(status))
		clauses = append(clauses, fmt.Sprintf("status = $%d", len(args)))
	}
	if keyword != "" {
		args = append(args, "%"+keyword+"%")
		clauses = append(clauses, fmt.Sprintf("(name ILIKE $%d OR phone ILIKE $%d)", len(args), len(args)))
	}
	if userID != "" {
		args = append(args, userID)
		clauses = append(clauses, fmt.Sprintf("user_id = $%d", len(args)))
	}

	whereClause := strings.Join(clauses, " AND ")
	countQuery := `SELECT COUNT(1) FROM orders`
	if whereClause != "" {
		countQuery += " WHERE " + whereClause
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	listQuery := `
		SELECT order_id, table_id, user_id, name, phone, notes, time, end_time, party_size, status, total
		FROM orders
	`
	if whereClause != "" {
		listQuery += " WHERE " + whereClause
	}
	dir := "DESC"
	if strings.EqualFold(strings.TrimSpace(sortOrder), "asc") {
		dir = "ASC"
	}
	listQuery += fmt.Sprintf(" ORDER BY time %s LIMIT $%d OFFSET $%d", dir, len(args)+1, len(args)+2)

	queryArgs := append(args, pageSize, offset)
	rows, err := r.db.QueryContext(ctx, listQuery, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	orders := make([]*domain.Order, 0)
	for rows.Next() {
		order := &domain.Order{}
		var endTime sql.NullTime
		if err := rows.Scan(
			&order.OrderID,
			&order.TableID,
			&order.UserID,
			&order.Name,
			&order.Phone,
			&order.Notes,
			&order.Time,
			&endTime,
			&order.PartySize,
			&order.Status,
			&order.Total,
		); err != nil {
			return nil, 0, err
		}
		if endTime.Valid {
			order.EndTime = endTime.Time
		}

		items, err := r.getOrderItems(ctx, r.db, order.OrderID)
		if err != nil {
			return nil, 0, err
		}
		order.Items = items
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *PostgresOrderRepository) GetOccupiedTableIDs(ctx context.Context, startTime, endTime time.Time) ([]string, error) {
	// Finds table IDs with a non-cancelled order whose window overlaps [startTime, endTime).
	// Overlap condition: order.time < endTime AND (order.end_time IS NULL OR order.end_time > startTime)
	const query = `
		SELECT DISTINCT table_id
		FROM orders
		WHERE status != $1
		  AND table_id != ''
		  AND time < $2
		  AND (end_time IS NULL OR end_time > $3)
	`
	rows, err := r.db.QueryContext(ctx, query, string(domain.StatusCancelled), endTime, startTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *PostgresOrderRepository) insertOrderItems(ctx context.Context, tx *sql.Tx, orderID string, items []*domain.OrderItem) error {
	const query = `
		INSERT INTO order_items (item_id, order_id, name, price, quantity, item_status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	for _, item := range items {
		status := item.ItemStatus
		if status == "" {
			status = domain.ItemStatusPending
		}
		if _, err := tx.ExecContext(ctx, query, item.ItemID, orderID, item.Name, item.Price, item.Quantity, string(status)); err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresOrderRepository) getOrderItems(ctx context.Context, q sqlQueryer, orderID string) ([]*domain.OrderItem, error) {
	const query = `
		SELECT item_id, name, price, quantity, item_status
		FROM order_items
		WHERE order_id = $1
		ORDER BY name ASC
	`

	rows, err := q.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*domain.OrderItem, 0)
	for rows.Next() {
		item := &domain.OrderItem{}
		var itemStatus string
		if err := rows.Scan(&item.ItemID, &item.Name, &item.Price, &item.Quantity, &itemStatus); err != nil {
			return nil, err
		}
		item.ItemStatus = domain.ItemStatus(itemStatus)
		if item.ItemStatus == "" {
			item.ItemStatus = domain.ItemStatusPending
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// UpdateItemStatus updates a single order item's status in-place without touching other fields.
func (r *PostgresOrderRepository) UpdateItemStatus(ctx context.Context, orderID, itemID string, status domain.ItemStatus) error {
	const query = `
		UPDATE order_items
		SET item_status = $3
		WHERE order_id = $1 AND item_id = $2
	`
	result, err := r.db.ExecContext(ctx, query, orderID, itemID, string(status))
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return domain.ErrOrderItemNotFound
	}
	return nil
}
