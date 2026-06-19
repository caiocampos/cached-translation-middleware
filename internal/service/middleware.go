package service

import (
	"context"
	"errors"
	"fmt"

	"cached-translation-middleware/internal/cache"
	"cached-translation-middleware/internal/model"
	"cached-translation-middleware/internal/util"

	"go.uber.org/zap"
)

type MiddlewareService interface {
	Translate(ctx context.Context, req *model.TranslationRequest) (*model.TranslationResponse, model.CacheSource, error)
	CheckAndUpdateTranslation(ctx context.Context, req *model.TranslationRequest) error
}

type middlewareService struct {
	cache       cache.TranslationCache
	translation TranslationService
	logger      *zap.Logger
}

func NewMiddlewareService(c cache.TranslationCache, t TranslationService, logger *zap.Logger) MiddlewareService {
	return &middlewareService{
		cache:       c,
		translation: t,
		logger:      logger,
	}
}

func (s *middlewareService) Translate(ctx context.Context, req *model.TranslationRequest) (*model.TranslationResponse, model.CacheSource, error) {
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
		s.logger.Warn("cache get error, falling through to upstream", zap.Error(err))
	}

	s.logger.Info("cache miss, calling upstream",
		zap.String("source", req.Source),
		zap.String("target", req.Target),
		zap.String("q", req.Q),
	)

	resp, err := s.translation.Translate(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("upstream translation failed: %w", err)
	}

	if setErr := s.cache.Set(ctx, req, resp); setErr != nil {
		s.logger.Warn("failed to store translation in cache", zap.Error(setErr))
	}

	return resp, model.CacheSourceMiss, nil
}

func (s *middlewareService) CheckAndUpdateTranslation(ctx context.Context, req *model.TranslationRequest) error {
	ttl, err := s.cache.GetTTL(ctx, req)
	if err == nil && *ttl > util.OneDay {
		s.logger.Info("valid cache, no need for update",
			zap.String("source", req.Source),
			zap.String("target", req.Target),
			zap.String("q", req.Q),
		)
		return nil
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		s.logger.Warn("cache get error, falling through to upstream", zap.Error(err))
	}

	s.logger.Info("cache miss, calling upstream",
		zap.String("source", req.Source),
		zap.String("target", req.Target),
		zap.String("q", req.Q),
	)

	resp, err := s.translation.Translate(ctx, req)
	if err != nil {
		return fmt.Errorf("upstream translation failed: %w", err)
	}

	if setErr := s.cache.Set(ctx, req, resp); setErr != nil {
		s.logger.Warn("failed to store translation in cache", zap.Error(setErr))
	}

	return nil
}
