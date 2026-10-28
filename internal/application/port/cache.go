package port

import (
	"context"
	"time"
)

type CacheProvider interface {
	SetObject(ctx context.Context, key string, data []byte, ttl time.Duration) error

	GetObject(ctx context.Context, key string) ([]byte, error)

	Del(ctx context.Context, key string) error

	SetString(ctx context.Context, key, str string, ttl time.Duration) error

	GetString(ctx context.Context, key string) (string, error)

	GetInt(ctx context.Context, key string) (int, error)

	Increment(ctx context.Context, key string) error

	IncrementAndGet(ctx context.Context, key string, ttlIfNew time.Duration) (int, error)

	AcquireLease(ctx context.Context, key, owner string, ttl time.Duration) (bool, error)

	RenewLease(ctx context.Context, key, owner string, ttl time.Duration) (bool, error)

	ReleaseLease(ctx context.Context, key, owner string) (bool, error)

	Exists(ctx context.Context, key string) (bool, error)

	TTL(ctx context.Context, key string) (time.Duration, error)

	Publish(ctx context.Context, channel string, data []byte) error

	Subscribe(ctx context.Context, channel string, handler func(ctx context.Context, data []byte)) error
}
