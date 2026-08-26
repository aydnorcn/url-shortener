package cache

import (
	"context"
	"errors"
	"time"
	"url-shortener/metrics"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
	return &redisCache{
		client: client,
	}
}

func (r redisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		metrics.CacheMissesTotal.Inc()
		return "", ErrCacheMiss
	}

	if err != nil {
		return "", err
	}

	metrics.CacheHitsTotal.Inc()

	return value, nil
}

func (r redisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {

	return r.client.Set(
		ctx,
		key,
		value,
		expiration,
	).Err()
}

func (r redisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
