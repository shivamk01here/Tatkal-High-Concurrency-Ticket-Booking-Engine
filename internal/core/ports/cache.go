package ports

import (
	"context"
	"time"
)

type CacheRepository interface {
	AcquireLock(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, key string) error
}
