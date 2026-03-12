package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shivamk01here/tatkal-engine/internal/core/ports"
)

type cacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) ports.CacheRepository {
	return &cacheRepository{client: client}
}

func (c *cacheRepository) AcquireLock(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	success, err := c.client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

func (c *cacheRepository) ReleaseLock(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
