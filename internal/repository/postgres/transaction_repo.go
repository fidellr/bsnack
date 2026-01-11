package postgres

import (
	"bsnack/internal/domain"
	"bsnack/internal/port"
	"context"
	"database/sql"
)

type TransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) port.TransactionRepository {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) Create(ctx context.Context, t *domain.Transaction) error {
	query := `
		INSERT INTO transactions (customer_id, product_id, quantity, total_price, transaction_date) 
		VALUES ($1, $2, $3, $4, $5) RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		t.CustomerID, t.ProductID, t.Quantity, t.TotalPrice, t.TransactionDate,
	).Scan(&t.ID)
}

func (r *TransactionRepo) GetReport(ctx context.Context, start, end string) (*domain.SalesReport, error) {
	report := &domain.SalesReport{
		StartDate:    start,
		EndDate:      end,
		Transactions: []domain.Transaction{},
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	queryAgg := `
        SELECT 
            COUNT(DISTINCT customer_id), 
            COALESCE(SUM(quantity), 0), 
            COALESCE(SUM(total_price), 0)
        FROM transactions 
        WHERE transaction_date::date >= $1::date 
          AND transaction_date::date <= $2::date`

	err = tx.QueryRowContext(ctx, queryAgg, start, end).Scan(
		&report.TotalCustomers, &report.TotalProducts, &report.TotalIncome,
	)
	if err != nil {
		return nil, err
	}

	queryBest := `
        SELECT p.name || ' - ' || p.flavor
        FROM transactions t
        JOIN products p ON t.product_id = p.id
        WHERE t.transaction_date::date >= $1::date 
          AND t.transaction_date::date <= $2::date
        GROUP BY p.name, p.flavor
        ORDER BY SUM(t.quantity) DESC LIMIT 1`

	var bestSeller sql.NullString
	err = tx.QueryRowContext(ctx, queryBest, start, end).Scan(&bestSeller)
	if err == nil && bestSeller.Valid {
		report.BestSeller = bestSeller.String
	} else {
		report.BestSeller = "No sales yet"
	}

	// compare the month/year of customer creation with the month/year of the transaction.
	queryList := `
        SELECT 
            t.id, 
            t.customer_id,
            t.product_id,
            c.name, 
            p.name, 
            p.size, 
            p.flavor, 
            t.quantity, 
            t.total_price, 
            t.transaction_date,
            (EXTRACT(MONTH FROM c.created_at) = EXTRACT(MONTH FROM t.transaction_date) AND 
             EXTRACT(YEAR FROM c.created_at) = EXTRACT(YEAR FROM t.transaction_date)) as is_new
        FROM transactions t
        JOIN customers c ON t.customer_id = c.id
        JOIN products p ON t.product_id = p.id
        WHERE t.transaction_date::date >= $1::date 
          AND t.transaction_date::date <= $2::date
        ORDER BY t.transaction_date DESC`

	rows, err := tx.QueryContext(ctx, queryList, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var trx domain.Transaction
		if err := rows.Scan(
			&trx.ID,
			&trx.CustomerID,
			&trx.ProductID,
			&trx.CustomerName,
			&trx.ProductName,
			&trx.ProductSize,
			&trx.ProductFlavor,
			&trx.Quantity,
			&trx.TotalPrice,
			&trx.TransactionDate,
			&trx.IsNewCustomer,
		); err != nil {
			return nil, err
		}

		report.Transactions = append(report.Transactions, trx)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return report, nil
}
