package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(c *redis.Client) *RedisCache {
	return &RedisCache{client: c}
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *RedisCache) SetNx(ctx context.Context, key string, data []byte, ttl time.Duration) (bool, error) {
	ok, err := r.client.SetNX(ctx, key, data, ttl).Result()
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	err := r.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}
