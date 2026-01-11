package service

import (
	"bsnack/internal/domain"
	"bsnack/internal/port"
	"context"
)

type ProductService struct {
	repo port.ProductRepository
}

func NewProductService(repo port.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) AddProduct(ctx context.Context, p *domain.Product) error {
	return s.repo.Create(ctx, p)
}

func (s *ProductService) GetProductsByDate(ctx context.Context, date string) ([]domain.Product, error) {
	return s.repo.GetByDate(ctx, date)
}
