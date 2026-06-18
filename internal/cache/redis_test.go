package cache_test

import (
	"context"
	"testing"
	"time"

	"cached-translation-middleware/internal/cache"
	"cached-translation-middleware/internal/model"

	"github.com/redis/go-redis/v9"
)

func newTestRedis(t *testing.T) *redis.Client {
	t.Helper()
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 15})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skipf("Redis not available, skipping: %v", err)
	}
	t.Cleanup(func() { rdb.FlushDB(context.Background()) })
	return rdb
}

func TestRedisCache_SetAndGet(t *testing.T) {
	rdb := newTestRedis(t)
	c := cache.NewRedisCache(rdb, time.Minute)
	ctx := context.Background()

	req := &model.TranslationRequest{Q: "Hello", Source: "en", Target: "es"}
	resp := &model.TranslationResponse{TranslatedText: "Hola"}

	if err := c.Set(ctx, req, resp); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	got, err := c.Get(ctx, req)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}

	if got.TranslatedText != resp.TranslatedText {
		t.Errorf("expected %q, got %q", resp.TranslatedText, got.TranslatedText)
	}
}

func TestRedisCache_SetAndGet_long(t *testing.T) {
	rdb := newTestRedis(t)
	c := cache.NewRedisCache(rdb, time.Minute)
	ctx := context.Background()

	req := &model.TranslationRequest{Q: "Learning a new language opens up a world of incredible opportunities and allows you to connect with people from completely different cultures", Source: "en", Target: "es"}
	resp := &model.TranslationResponse{TranslatedText: "Aprender un nuevo idioma abre un mundo de oportunidades increíbles y te permite conectar con personas de culturas completamente diferentes"}

	if err := c.Set(ctx, req, resp); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	got, err := c.Get(ctx, req)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}

	if got.TranslatedText != resp.TranslatedText {
		t.Errorf("expected %q, got %q", resp.TranslatedText, got.TranslatedText)
	}
}

func TestRedisCache_Get_MissReturnsError(t *testing.T) {
	rdb := newTestRedis(t)
	c := cache.NewRedisCache(rdb, time.Minute)
	ctx := context.Background()

	req := &model.TranslationRequest{Q: "nonexistent", Source: "en", Target: "fr"}

	_, err := c.Get(ctx, req)
	if err != cache.ErrCacheMiss {
		t.Errorf("expected ErrCacheMiss, got %v", err)
	}
}

func TestRedisCache_GetTTL(t *testing.T) {
	rdb := newTestRedis(t)
	c := cache.NewRedisCache(rdb, time.Minute)
	ctx := context.Background()

	req := &model.TranslationRequest{Q: "Hello", Source: "en", Target: "es"}

	got, err := c.GetTTL(ctx, req)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}

	if *got != time.Minute {
		t.Errorf("expected %q, got %q", time.Minute, got)
	}
}

func TestRedisCache_GetTTL_MissReturnsError(t *testing.T) {
	rdb := newTestRedis(t)
	c := cache.NewRedisCache(rdb, time.Minute)
	ctx := context.Background()

	req := &model.TranslationRequest{Q: "nonexistent", Source: "en", Target: "fr"}

	_, err := c.GetTTL(ctx, req)
	if err != cache.ErrCacheMiss {
		t.Errorf("expected ErrCacheMiss, got %v", err)
	}
}
