package service

import (
	"bsnack/internal/domain"
	"bsnack/internal/port"
	"bsnack/pkg/logger"
	"context"
	"errors"
	"math"
	"time"
)

type TransactionService struct {
	repoProd  port.ProductRepository
	repoCust  port.CustomerRepository
	repoTrans port.TransactionRepository
	repoCache port.CacheRepository
}

func NewTransactionService(
	rp port.ProductRepository,
	rc port.CustomerRepository,
	rt port.TransactionRepository,
	cache port.CacheRepository,
) *TransactionService {
	return &TransactionService{
		repoProd:  rp,
		repoCust:  rc,
		repoTrans: rt,
		repoCache: cache,
	}
}

type PurchaseRequest struct {
	CustomerName    string `json:"customer_name"`
	ProductID       int64  `json:"product_id"`
	Quantity        int    `json:"quantity"`
	TransactionDate string `json:"transaction_date"`
}

func (s *TransactionService) Purchase(ctx context.Context, req PurchaseRequest) error {
	if req.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	product, err := s.repoProd.GetByID(ctx, req.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	if product.Quantity < req.Quantity {
		return errors.New("insufficient stock")
	}

	customer, err := s.repoCust.GetByName(ctx, req.CustomerName)
	if err != nil || customer == nil {
		customer = &domain.Customer{Name: req.CustomerName, Points: 0}
		if err := s.repoCust.Create(ctx, customer); err != nil {
			return err
		}
	}

	totalPrice := product.Price * float64(req.Quantity)

	// rule: 1 Point per Rp 1,000
	pointsEarned := int(math.Floor(totalPrice / 1000))

	if err := s.repoProd.UpdateStock(ctx, product.ID, -req.Quantity); err != nil {
		return err
	}
	if err := s.repoCust.UpdatePoints(ctx, customer.ID, pointsEarned); err != nil {
		return err
	}
	var txDate time.Time
	if req.TransactionDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.TransactionDate)
		if err != nil {
			return errors.New("invalid transaction_date format (use YYYY-MM-DD)")
		}
		txDate = parsedDate
	} else {
		txDate = time.Now()
	}

	tx := &domain.Transaction{
		CustomerID:      customer.ID,
		ProductID:       product.ID,
		Quantity:        req.Quantity,
		TotalPrice:      totalPrice,
		TransactionDate: txDate,
	}
	return s.repoTrans.Create(ctx, tx)
}

// Redeem handles point exchange for products
func (s *TransactionService) Redeem(ctx context.Context, customerName string, productID int64) error {
	product, err := s.repoProd.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	var cost int
	switch product.Size {
	case domain.SizeSmall:
		cost = 200
	case domain.SizeMedium:
		cost = 300
	case domain.SizeLarge:
		cost = 500
	default:
		return errors.New("invalid product size")
	}

	customer, err := s.repoCust.GetByName(ctx, customerName)
	if err != nil {
		return errors.New("customer not found")
	}

	if customer.Points < cost {
		logger.Warn("redemption failed: insufficient points",
			"customer", customerName,
			"points", customer.Points,
			"required", cost)
		return errors.New("insufficient points")
	}

	if err := s.repoCust.UpdatePoints(ctx, customer.ID, -cost); err != nil {
		return err
	}
	if err := s.repoProd.UpdateStock(ctx, product.ID, -1); err != nil {
		return err
	}

	return nil
}

// GetReport uses Cache-Aside pattern
func (s *TransactionService) GetReport(ctx context.Context, start, end string) (*domain.SalesReport, error) {
	if end == "" {
		now := time.Now()
		end = now.Format("2006-01-02")
	}

	cached, err := s.repoCache.GetReport(ctx, start, end)
	if err == nil && cached != nil {
		return cached, nil
	}

	report, err := s.repoTrans.GetReport(ctx, start, end)
	if err != nil {
		return nil, err
	}

	_ = s.repoCache.SetReport(ctx, start, end, report, 5*time.Minute)

	return report, nil
}
