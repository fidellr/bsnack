package redis

import (
	"bsnack/internal/domain"
	"bsnack/internal/port"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) port.CacheRepository {
	return &RedisRepo{client: client}
}

// reportKey generates a unique key based on the date range
func (r *RedisRepo) reportKey(start, end string) string {
	return fmt.Sprintf("report:%s:%s", start, end)
}

func (r *RedisRepo) GetReport(ctx context.Context, start, end string) (*domain.SalesReport, error) {
	val, err := r.client.Get(ctx, r.reportKey(start, end)).Result()
	if err == redis.Nil {
		return nil, nil // cache miss
	}
	if err != nil {
		return nil, err
	}

	var report domain.SalesReport
	if err := json.Unmarshal([]byte(val), &report); err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *RedisRepo) SetReport(ctx context.Context, start, end string, report *domain.SalesReport, ttl time.Duration) error {
	data, err := json.Marshal(report)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.reportKey(start, end), data, ttl).Err()
}

func (r *RedisRepo) InvalidateProducts(ctx context.Context, date string) error {
	// delete specific key if caching product lists
	return r.client.Del(ctx, fmt.Sprintf("products:%s", date)).Err()
}
