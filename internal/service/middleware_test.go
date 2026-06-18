package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"cached-translation-middleware/internal/cache"
	"cached-translation-middleware/internal/model"
	"cached-translation-middleware/internal/service"

	"go.uber.org/zap"
)

// ── Mocks ─────────────────────────────────────────────────────────────────────

type mockCache struct {
	getResp    *model.TranslationResponse
	getErr     error
	getTTLResp *time.Duration
	getTTLErr  error
	setErr     error
	setCalled  bool
}

func (m *mockCache) Get(_ context.Context, _ *model.TranslationRequest) (*model.TranslationResponse, error) {
	return m.getResp, m.getErr
}

func (m *mockCache) GetTTL(_ context.Context, _ *model.TranslationRequest) (*time.Duration, error) {
	return m.getTTLResp, m.getTTLErr
}

func (m *mockCache) Set(_ context.Context, _ *model.TranslationRequest, _ *model.TranslationResponse) error {
	m.setCalled = true
	return m.setErr
}

func (m *mockCache) SetWithTTL(_ context.Context, _ *model.TranslationRequest, _ *model.TranslationResponse, _ time.Duration) error {
	m.setCalled = true
	return m.setErr
}

type mockTranslation struct {
	resp *model.TranslationResponse
	err  error
}

func (m *mockTranslation) Translate(_ context.Context, _ *model.TranslationRequest) (*model.TranslationResponse, error) {
	return m.resp, m.err
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func TestMiddlewareService_CacheHit(t *testing.T) {
	cached := &model.TranslationResponse{TranslatedText: "Hola"}
	mc := &mockCache{getResp: cached, getErr: nil}
	mt := &mockTranslation{}

	svc := service.NewMiddlewareService(mc, mt, zap.NewNop())
	req := &model.TranslationRequest{Q: "Hello", Source: "en", Target: "es"}

	resp, src, err := svc.Translate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src != model.CacheSourceHit {
		t.Errorf("expected cache hit, got %s", src)
	}
	if resp.TranslatedText != "Hola" {
		t.Errorf("expected 'Hola', got %s", resp.TranslatedText)
	}
}

func TestMiddlewareService_CacheMissFallsThrough(t *testing.T) {
	mc := &mockCache{getErr: cache.ErrCacheMiss}
	mt := &mockTranslation{resp: &model.TranslationResponse{TranslatedText: "Hola"}}

	svc := service.NewMiddlewareService(mc, mt, zap.NewNop())
	req := &model.TranslationRequest{Q: "Hello", Source: "en", Target: "es"}

	resp, src, err := svc.Translate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src != model.CacheSourceMiss {
		t.Errorf("expected upstream, got %s", src)
	}
	if resp.TranslatedText != "Hola" {
		t.Errorf("expected 'Hola', got %s", resp.TranslatedText)
	}
	if !mc.setCalled {
		t.Error("expected cache Set to be called")
	}
}

func TestMiddlewareService_UpstreamError(t *testing.T) {
	mc := &mockCache{getErr: cache.ErrCacheMiss}
	mt := &mockTranslation{err: errors.New("timeout")}

	svc := service.NewMiddlewareService(mc, mt, zap.NewNop())
	req := &model.TranslationRequest{Q: "Hello", Source: "en", Target: "es"}

	_, _, err := svc.Translate(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
