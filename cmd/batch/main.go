package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cached-translation-middleware/cmd/batch/process"
	"cached-translation-middleware/config"
	"cached-translation-middleware/internal/cache"
	"cached-translation-middleware/internal/service"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// ── Logger ──────────────────────────────────────────────────────────────
	logger, _ := zap.NewProduction()
	defer logger.Sync() //nolint:errcheck

	// ── Config ──────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// ── Redis ───────────────────────────────────────────────────────────────
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Fatal("failed to connect to Redis", zap.Error(err))
	}
	logger.Info("connected to Redis", zap.String("addr", cfg.Redis.Addr))

	// ── Dependencies ────────────────────────────────────────────────────────
	translationCache := cache.NewRedisCache(rdb, cfg.Redis.CacheTTL)

	httpTranslationClient := &http.Client{Timeout: cfg.Translation.Timeout}
	translationSvc := service.NewTranslationService(cfg.Translation.APIURL, httpTranslationClient)

	httpGithubClient := &http.Client{Timeout: cfg.Github.Timeout}
	githubSvc := service.NewGithubService(cfg.Github.APIURL, httpGithubClient)

	middlewareSvc := service.NewMiddlewareService(translationCache, translationSvc, logger)

	logger.Info("starting process")
	process.Process(logger, cfg, githubSvc, middlewareSvc)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down process...")

	if err := rdb.Close(); err != nil {
		logger.Warn("error closing Redis connection", zap.Error(err))
	}

	logger.Info("server exited gracefully")
}
