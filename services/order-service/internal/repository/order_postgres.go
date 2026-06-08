package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
			name VARCHAR(120) NOT NULL,
			phone VARCHAR(120) NOT NULL DEFAULT '',
			time TIMESTAMP NOT NULL,
			party_size INTEGER NOT NULL,
			status VARCHAR(32) NOT NULL,
			total DOUBLE PRECISION NOT NULL DEFAULT 0
		);

		CREATE TABLE IF NOT EXISTS order_items (
			item_id VARCHAR(36) PRIMARY KEY REFERENCES menu_items(item_id),
			order_id VARCHAR(36) NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
			name VARCHAR(200) NOT NULL,
			price DOUBLE PRECISION NOT NULL,
			quantity INTEGER NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
		CREATE INDEX IF NOT EXISTS idx_orders_name ON orders(name);
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
			name, phone, time, party_size, status, total
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING order_id
	`

	if err := tx.QueryRowContext(ctx, insertOrderQuery,
		order.Name,
		order.Phone,
		order.Time,
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
		SELECT order_id, name, phone, time, party_size, status, total
		FROM orders
		WHERE order_id = $1
	`

	order := &domain.Order{}
	if err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&order.OrderID,
		&order.Name,
		&order.Phone,
		&order.Time,
		&order.PartySize,
		&order.Status,
		&order.Total,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
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
		SET name = $2, phone = $3, time = $4, party_size = $5, status = $6, total = $7
		WHERE order_id = $1
	`

	result, err := tx.ExecContext(ctx, updateOrderQuery,
		order.OrderID,
		order.Name,
		order.Phone,
		order.Time,
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

func (r *PostgresOrderRepository) List(ctx context.Context, page, pageSize int, status domain.OrderStatus, keyword string) ([]*domain.Order, int, error) {
	clauses := make([]string, 0, 2)
	args := make([]any, 0, 3)

	if status != "" {
		args = append(args, string(status))
		clauses = append(clauses, fmt.Sprintf("status = $%d", len(args)))
	}
	if keyword != "" {
		args = append(args, "%"+keyword+"%")
		clauses = append(clauses, fmt.Sprintf("(name ILIKE $%d OR phone ILIKE $%d)", len(args), len(args)))
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
		SELECT order_id, name, phone, time, party_size, status, total
		FROM orders
	`
	if whereClause != "" {
		listQuery += " WHERE " + whereClause
	}
	listQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)

	queryArgs := append(args, pageSize, offset)
	rows, err := r.db.QueryContext(ctx, listQuery, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	orders := make([]*domain.Order, 0)
	for rows.Next() {
		order := &domain.Order{}
		if err := rows.Scan(
			&order.OrderID,
			&order.Name,
			&order.Phone,
			&order.Time,
			&order.PartySize,
			&order.Status,
			&order.Total,
		); err != nil {
			return nil, 0, err
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

func (r *PostgresOrderRepository) insertOrderItems(ctx context.Context, tx *sql.Tx, orderID string, items []*domain.OrderItem) error {
	const query = `
		INSERT INTO order_items (item_id, order_id, name, price, quantity)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, item := range items {
		if _, err := tx.ExecContext(ctx, query, item.ItemID, orderID, item.Name, item.Price, item.Quantity); err != nil {
			return err
		}
	}

	return nil
}

func (r *PostgresOrderRepository) getOrderItems(ctx context.Context, q sqlQueryer, orderID string) ([]*domain.OrderItem, error) {
	const query = `
		SELECT item_id, name, price, quantity
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
		if err := rows.Scan(&item.ItemID, &item.Name, &item.Price, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
