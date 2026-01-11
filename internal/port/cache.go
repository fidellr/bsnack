package port

import (
	"bsnack/internal/domain"
	"context"
	"time"
)

// CacheRepository defines Redis operations
type CacheRepository interface {
	GetReport(ctx context.Context, start, end string) (*domain.SalesReport, error)
	SetReport(ctx context.Context, start, end string, report *domain.SalesReport, ttl time.Duration) error

	InvalidateProducts(ctx context.Context, date string) error
}
