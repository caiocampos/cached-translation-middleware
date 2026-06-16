package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"cached-translation-middleware/internal/model"

	"github.com/redis/go-redis/v9"
)

// ErrCacheMiss is returned when a key is not found in the cache.
var ErrCacheMiss = errors.New("cache miss")

// TranslationCache defines the contract for caching translation results.
type TranslationCache interface {
	Get(ctx context.Context, req *model.TranslationRequest) (*model.TranslationResponse, error)
	Set(ctx context.Context, req *model.TranslationRequest, resp *model.TranslationResponse) error
}

type redisCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisCache creates a new Redis-backed TranslationCache.
func NewRedisCache(client *redis.Client, ttl time.Duration) TranslationCache {
	return &redisCache{client: client, ttl: ttl}
}

// buildKey creates a deterministic, human-readable cache key.
func buildKey(req *model.TranslationRequest) string {
	return fmt.Sprintf("translation:%s:%s:%s", req.Source, req.Target, req.Q)
}

func (c *redisCache) Get(ctx context.Context, req *model.TranslationRequest) (*model.TranslationResponse, error) {
	key := buildKey(req)

	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}

	var resp model.TranslationResponse
	if err := json.Unmarshal([]byte(val), &resp); err != nil {
		return nil, fmt.Errorf("unmarshal cached value: %w", err)
	}

	return &resp, nil
}

func (c *redisCache) Set(ctx context.Context, req *model.TranslationRequest, resp *model.TranslationResponse) error {
	key := buildKey(req)

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}

	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}
