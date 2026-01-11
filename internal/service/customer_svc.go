package service

import (
	"bsnack/internal/domain"
	"bsnack/internal/port"
	"context"
)

type CustomerService struct {
	repoCust port.CustomerRepository
}

func NewCustomerService(rc port.CustomerRepository) *CustomerService {
	return &CustomerService{
		repoCust: rc,
	}
}

func (s *CustomerService) GetAllCustomers(ctx context.Context) ([]domain.Customer, error) {
	return s.repoCust.ListAll(ctx)
}
