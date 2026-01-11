package service_test

import (
	"bsnack/internal/domain"
	"bsnack/internal/service"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPurchase_Success_NewCustomer(t *testing.T) {
	mockProd := new(MockProductRepo)
	mockCust := new(MockCustomerRepo)
	mockTrans := new(MockTransactionRepo)
	mockCache := new(MockCacheRepo)

	svc := service.NewTransactionService(mockProd, mockCust, mockTrans, mockCache)
	ctx := context.TODO()

	req := service.PurchaseRequest{
		CustomerName: "Budi",
		ProductID:    1,
		Quantity:     2,
	}

	product := &domain.Product{ID: 1, Price: 10000, Quantity: 10} // 10k price

	mockProd.On("GetByID", ctx, int64(1)).Return(product, nil)

	// customer not found -> register new customer
	mockCust.On("GetByName", ctx, "Budi").Return(nil, errors.New("not found"))
	mockCust.On("Create", ctx, mock.AnythingOfType("*domain.Customer")).Return(nil)

	// deduct Stock (qty 2)
	mockProd.On("UpdateStock", ctx, int64(1), -2).Return(nil)

	// calculation: (10,000 * 2) / 1000 = 20 points
	mockCust.On("UpdatePoints", ctx, int64(1), 20).Return(nil)

	mockTrans.On("Create", ctx, mock.AnythingOfType("*domain.Transaction")).Return(nil)

	err := svc.Purchase(ctx, req)

	assert.NoError(t, err)
	mockProd.AssertExpectations(t)
	mockCust.AssertExpectations(t)
	mockTrans.AssertExpectations(t)
}

func TestPurchase_InsufficientStock(t *testing.T) {
	mockProd := new(MockProductRepo)
	svc := service.NewTransactionService(mockProd, nil, nil, nil)

	product := &domain.Product{ID: 1, Quantity: 1}
	mockProd.On("GetByID", context.TODO(), int64(1)).Return(product, nil)

	req := service.PurchaseRequest{ProductID: 1, Quantity: 2}
	err := svc.Purchase(context.TODO(), req)

	assert.Error(t, err)
	assert.Equal(t, "insufficient stock", err.Error())
}

func TestRedeem_Success(t *testing.T) {
	mockProd := new(MockProductRepo)
	mockCust := new(MockCustomerRepo)
	svc := service.NewTransactionService(mockProd, mockCust, nil, nil)
	ctx := context.TODO()

	product := &domain.Product{ID: 1, Size: domain.SizeSmall, Quantity: 10}
	customer := &domain.Customer{ID: 5, Name: "Fery", Points: 250} // Has 250

	mockProd.On("GetByID", ctx, int64(1)).Return(product, nil)
	mockCust.On("GetByName", ctx, "Fery").Return(customer, nil)

	mockCust.On("UpdatePoints", ctx, int64(5), -200).Return(nil)
	mockProd.On("UpdateStock", ctx, int64(1), -1).Return(nil)

	err := svc.Redeem(ctx, "Fery", 1)
	assert.NoError(t, err)
}

func TestRedeem_InsufficientPoints(t *testing.T) {
	mockProd := new(MockProductRepo)
	mockCust := new(MockCustomerRepo)
	svc := service.NewTransactionService(mockProd, mockCust, nil, nil)
	ctx := context.TODO()

	product := &domain.Product{ID: 1, Size: domain.SizeSmall}
	customer := &domain.Customer{ID: 5, Points: 50}

	mockProd.On("GetByID", ctx, int64(1)).Return(product, nil)
	mockCust.On("GetByName", ctx, "Fery").Return(customer, nil)

	err := svc.Redeem(ctx, "Fery", 1)

	assert.Error(t, err)
	assert.Equal(t, "insufficient points", err.Error())
}

func TestGetReport_CacheHit(t *testing.T) {
	mockCache := new(MockCacheRepo)
	mockTrans := new(MockTransactionRepo)
	svc := service.NewTransactionService(nil, nil, mockTrans, mockCache)
	ctx := context.TODO()

	cachedReport := &domain.SalesReport{TotalIncome: 50000}

	mockCache.On("GetReport", ctx, "2025-01-01", "2025-01-31").Return(cachedReport, nil)

	res, err := svc.GetReport(ctx, "2025-01-01", "2025-01-31")

	assert.NoError(t, err)
	assert.Equal(t, 50000.0, res.TotalIncome)
	mockTrans.AssertNotCalled(t, "GetReport")
}
