package initialization

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

type Cache struct {
	rdb *redis.Client
}

func InitCache(cfg config.CacheConfig) (*Cache, error) {
	rAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	options := &redis.Options{
		Addr:     rAddr,
		Password: cfg.Password,
		DB:       0,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	}
	if cfg.UseSSL {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	rdb := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Cache{
		rdb,
	}, nil
}

func (c *Cache) Client() *redis.Client {
	return c.rdb
}

func (c *Cache) Close() {
	_ = c.rdb.Close()
}
