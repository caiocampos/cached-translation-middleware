package service

import (
	"context"
	"errors"
	"fmt"

	"cached-translation-middleware/internal/cache"
	"cached-translation-middleware/internal/model"

	"go.uber.org/zap"
)

// MiddlewareService orchestrates cache lookup and upstream translation.
type MiddlewareService interface {
	Translate(ctx context.Context, req *model.TranslationRequest) (*model.TranslationResponse, model.CacheSource, error)
}

type middlewareService struct {
	cache       cache.TranslationCache
	translation TranslationService
	logger      *zap.Logger
}

// NewMiddlewareService creates the cache-aside middleware orchestrator.
func NewMiddlewareService(c cache.TranslationCache, t TranslationService, logger *zap.Logger) MiddlewareService {
	return &middlewareService{
		cache:       c,
		translation: t,
		logger:      logger,
	}
}

func (s *middlewareService) Translate(ctx context.Context, req *model.TranslationRequest) (*model.TranslationResponse, model.CacheSource, error) {
	// 1. Try cache first
	cached, err := s.cache.Get(ctx, req)
	if err == nil {
		s.logger.Info("cache hit",
			zap.String("source", req.Source),
			zap.String("target", req.Target),
			zap.String("q", req.Q),
		)
		return cached, model.CacheSourceHit, nil
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		// Non-fatal: log the error and fall through to upstream
		s.logger.Warn("cache get error, falling through to upstream", zap.Error(err))
	}

	// 2. Call upstream translation API
	s.logger.Info("cache miss, calling upstream",
		zap.String("source", req.Source),
		zap.String("target", req.Target),
		zap.String("q", req.Q),
	)

	resp, err := s.translation.Translate(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("upstream translation failed: %w", err)
	}

	// 3. Persist result in cache (best-effort)
	if setErr := s.cache.Set(ctx, req, resp); setErr != nil {
		s.logger.Warn("failed to store translation in cache", zap.Error(setErr))
	}

	return resp, model.CacheSourceMiss, nil
}
