package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID              uuid.UUID `json:"id"`
	CustomerID      int64     `json:"customer_id"`
	CustomerName    string    `json:"customer_name"`
	ProductID       int64     `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductSize     string    `json:"product_size"`
	ProductFlavor   string    `json:"product_flavor"`
	Quantity        int       `json:"quantity"`
	TotalPrice      float64   `json:"total_price"`
	TransactionDate time.Time `json:"transaction_date"`
	IsNewCustomer   bool      `json:"is_new_customer"`
}

type SalesReport struct {
	StartDate      string        `json:"start_date"`
	EndDate        string        `json:"end_date"`
	TotalCustomers int           `json:"total_customers"`
	TotalProducts  int           `json:"total_products"`
	TotalIncome    float64       `json:"total_income"`
	BestSeller     string        `json:"best_seller"`
	Transactions   []Transaction `json:"transactions"`
}
