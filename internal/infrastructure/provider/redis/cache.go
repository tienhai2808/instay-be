package redis

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/redis/go-redis/v9"
)

type cacheProviderImpl struct {
	rdb *redis.Client
}

var (
	incrementAndGetScript = redis.NewScript(`
		local count = redis.call("INCR", KEYS[1])
		if count == 1 and tonumber(ARGV[1]) > 0 then
			redis.call("PEXPIRE", KEYS[1], ARGV[1])
		end
		return count
	`)
	renewLeaseScript = redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("PEXPIRE", KEYS[1], ARGV[2])
		end
		return 0
	`)
	releaseLeaseScript = redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		end
		return 0
	`)
)

func NewCacheProvider(rdb *redis.Client) port.CacheProvider {
	return &cacheProviderImpl{rdb}
}

func (p *cacheProviderImpl) SetString(ctx context.Context, key, str string, ttl time.Duration) error {
	return p.rdb.Set(ctx, key, str, ttl).Err()
}

func (p *cacheProviderImpl) GetString(ctx context.Context, key string) (string, error) {
	str, err := p.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return str, nil
}

func (p *cacheProviderImpl) SetObject(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	return p.rdb.Set(ctx, key, data, ttl).Err()
}

func (p *cacheProviderImpl) Exists(ctx context.Context, key string) (bool, error) {
	count, err := p.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (p *cacheProviderImpl) GetObject(ctx context.Context, key string) ([]byte, error) {
	data, err := p.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return []byte(data), nil
}

func (p *cacheProviderImpl) Del(ctx context.Context, key string) error {
	return p.rdb.Del(ctx, key).Err()
}

func (p *cacheProviderImpl) GetInt(ctx context.Context, key string) (int, error) {
	num, err := p.rdb.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return num, nil
}

func (p *cacheProviderImpl) Increment(ctx context.Context, key string) error {
	return p.rdb.Incr(ctx, key).Err()
}

func (p *cacheProviderImpl) IncrementAndGet(ctx context.Context, key string, ttlIfNew time.Duration) (int, error) {
	result, err := incrementAndGetScript.Run(ctx, p.rdb, []string{key}, ttlIfNew.Milliseconds()).Int()
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (p *cacheProviderImpl) AcquireLease(ctx context.Context, key, owner string, ttl time.Duration) (bool, error) {
	return p.rdb.SetNX(ctx, key, owner, ttl).Result()
}

func (p *cacheProviderImpl) RenewLease(ctx context.Context, key, owner string, ttl time.Duration) (bool, error) {
	result, err := renewLeaseScript.Run(ctx, p.rdb, []string{key}, owner, ttl.Milliseconds()).Int()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func (p *cacheProviderImpl) ReleaseLease(ctx context.Context, key, owner string) (bool, error) {
	result, err := releaseLeaseScript.Run(ctx, p.rdb, []string{key}, owner).Int()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func (p *cacheProviderImpl) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := p.rdb.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return ttl, nil
}

func (p *cacheProviderImpl) Publish(ctx context.Context, channel string, data []byte) error {
	return p.rdb.Publish(ctx, channel, data).Err()
}

func (p *cacheProviderImpl) Subscribe(ctx context.Context, channel string, handler func(ctx context.Context, data []byte)) error {
	pubsub := p.rdb.Subscribe(ctx, channel)

	if _, err := pubsub.Receive(ctx); err != nil {
		return err
	}

	go func() {
		defer func() {
			_ = pubsub.Close()
		}()

		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				handler(ctx, []byte(msg.Payload))
			}
		}
	}()

	return nil
}
