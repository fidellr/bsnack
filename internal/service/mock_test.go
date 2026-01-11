package service_test

import (
	"bsnack/internal/domain"
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockProductRepo mocks port.ProductRepository
type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) Create(ctx context.Context, p *domain.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *MockProductRepo) GetByDate(ctx context.Context, date string) ([]domain.Product, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]domain.Product), args.Error(1)
}
func (m *MockProductRepo) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}
func (m *MockProductRepo) UpdateStock(ctx context.Context, id int64, delta int) error {
	args := m.Called(ctx, id, delta)
	return args.Error(0)
}

// MockCustomerRepo mocks port.CustomerRepository
type MockCustomerRepo struct {
	mock.Mock
}

func (m *MockCustomerRepo) GetByName(ctx context.Context, name string) (*domain.Customer, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}
func (m *MockCustomerRepo) Create(ctx context.Context, c *domain.Customer) error {
	args := m.Called(ctx, c)
	c.ID = 1 // Simulate DB assigning ID
	return args.Error(0)
}
func (m *MockCustomerRepo) UpdatePoints(ctx context.Context, id int64, points int) error {
	args := m.Called(ctx, id, points)
	return args.Error(0)
}
func (m *MockCustomerRepo) ListAll(ctx context.Context) ([]domain.Customer, error) {
	return nil, nil
}

// MockTransactionRepo mocks port.TransactionRepository
type MockTransactionRepo struct {
	mock.Mock
}

func (m *MockTransactionRepo) Create(ctx context.Context, t *domain.Transaction) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}
func (m *MockTransactionRepo) GetReport(ctx context.Context, start, end string) (*domain.SalesReport, error) {
	args := m.Called(ctx, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SalesReport), args.Error(1)
}

// MockCacheRepo mocks port.CacheRepository
type MockCacheRepo struct {
	mock.Mock
}

func (m *MockCacheRepo) GetReport(ctx context.Context, start, end string) (*domain.SalesReport, error) {
	args := m.Called(ctx, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SalesReport), args.Error(1)
}
func (m *MockCacheRepo) SetReport(ctx context.Context, start, end string, report *domain.SalesReport, ttl time.Duration) error {
	args := m.Called(ctx, start, end, report, ttl)
	return args.Error(0)
}
func (m *MockCacheRepo) InvalidateProducts(ctx context.Context, date string) error {
	return nil
}
