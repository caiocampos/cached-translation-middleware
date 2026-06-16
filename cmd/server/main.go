package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cached-translation-middleware/config"
	"cached-translation-middleware/internal/cache"
	"cached-translation-middleware/internal/handler"
	"cached-translation-middleware/internal/middleware"
	"cached-translation-middleware/internal/service"

	"github.com/gin-gonic/gin"
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

	httpClient := &http.Client{Timeout: cfg.Translation.Timeout}
	translationSvc := service.NewTranslationService(cfg.Translation.APIURL, httpClient)

	middlewareSvc := service.NewMiddlewareService(translationCache, translationSvc, logger)

	translationHandler := handler.NewTranslationHandler(middlewareSvc, logger)

	// ── Router ──────────────────────────────────────────────────────────────
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.RequestLogger(logger))

	router.GET("/health", translationHandler.HealthCheck)
	router.POST("/translate", translationHandler.Translate)

	// ── Server ──────────────────────────────────────────────────────────────
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start in a goroutine so we can listen for shutdown signals
	go func() {
		logger.Info("server starting", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	// ── Graceful Shutdown ───────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("forced shutdown", zap.Error(err))
	}

	if err := rdb.Close(); err != nil {
		logger.Warn("error closing Redis connection", zap.Error(err))
	}

	logger.Info("server exited gracefully")
}
