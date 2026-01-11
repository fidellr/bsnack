package postgres

import (
	"bsnack/internal/domain"
	"bsnack/internal/port"
	"context"
	"database/sql"
	"fmt"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) port.ProductRepository {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(ctx context.Context, p *domain.Product) error {
	query := `
		INSERT INTO products (name, type, flavor, size, price, quantity, manufacturing_date) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRowContext(ctx, query,
		p.Name, p.Type, p.Flavor, p.Size, p.Price, p.Quantity, p.ManufacturingDate,
	).Scan(&p.ID)
}

func (r *ProductRepo) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	p := &domain.Product{}
	query := `SELECT id, name, type, flavor, size, price, quantity, manufacturing_date FROM products WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Type, &p.Flavor, &p.Size, &p.Price, &p.Quantity, &p.ManufacturingDate,
	)
	return p, err
}

func (r *ProductRepo) GetByDate(ctx context.Context, date string) ([]domain.Product, error) {
	query := `
		SELECT id, name, type, flavor, size, price, quantity, manufacturing_date 
		FROM products WHERE manufacturing_date = $1`

	rows, err := r.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Flavor, &p.Size, &p.Price, &p.Quantity, &p.ManufacturingDate); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
func (r *ProductRepo) UpdateStock(ctx context.Context, id int64, delta int) error {
	// fmt.Printf("DEBUG: Updating Stock for ID %d with delta %d\n", id, delta)

	query := `UPDATE products SET quantity = quantity + $1 WHERE id = $2`

	res, err := r.db.ExecContext(ctx, query, delta, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("product with id %d not found during stock update", id)
	}

	return nil
}
