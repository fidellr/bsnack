package port

import (
	"bsnack/internal/domain"
	"context"
)

// ProductRepository defines interactions with product data
type ProductRepository interface {
	Create(ctx context.Context, p *domain.Product) error
	GetByDate(ctx context.Context, date string) ([]domain.Product, error)
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	UpdateStock(ctx context.Context, id int64, delta int) error
}

// CustomerRepository defines interactions with customer data
type CustomerRepository interface {
	GetByName(ctx context.Context, name string) (*domain.Customer, error)
	Create(ctx context.Context, c *domain.Customer) error
	UpdatePoints(ctx context.Context, id int64, points int) error
	ListAll(ctx context.Context) ([]domain.Customer, error)
}

// TransactionRepository defines interactions with sales data
type TransactionRepository interface {
	Create(ctx context.Context, t *domain.Transaction) error
	// GetReport aggregates data for the specific date range
	GetReport(ctx context.Context, startDate, endDate string) (*domain.SalesReport, error)
}
