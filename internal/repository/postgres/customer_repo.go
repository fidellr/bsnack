package postgres

import (
	"bsnack/internal/domain"
	"bsnack/internal/port"
	"context"
	"database/sql"
)

type CustomerRepo struct {
	db *sql.DB
}

func NewCustomerRepo(db *sql.DB) port.CustomerRepository {
	return &CustomerRepo{db: db}
}

func (r *CustomerRepo) GetByName(ctx context.Context, name string) (*domain.Customer, error) {
	c := &domain.Customer{}
	query := `SELECT id, name, points, created_at FROM customers WHERE name = $1`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&c.ID, &c.Name, &c.Points, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CustomerRepo) Create(ctx context.Context, c *domain.Customer) error {
	query := `INSERT INTO customers (name, points) VALUES ($1, $2) RETURNING id, created_at`
	return r.db.QueryRowContext(ctx, query, c.Name, c.Points).Scan(&c.ID, &c.CreatedAt)
}

func (r *CustomerRepo) UpdatePoints(ctx context.Context, id int64, points int) error {
	query := `UPDATE customers SET points = points + $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, points, id)
	return err
}

func (r *CustomerRepo) ListAll(ctx context.Context) ([]domain.Customer, error) {
	query := `SELECT id, name, points, created_at, updated_at FROM customers ORDER BY points DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []domain.Customer
	for rows.Next() {
		var c domain.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Points, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}
