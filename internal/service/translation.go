package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cached-translation-middleware/internal/model"
)

// TranslationService defines the contract for calling the upstream translation API.
type TranslationService interface {
	Translate(ctx context.Context, req *model.TranslationRequest) (*model.TranslationResponse, error)
}

type translationService struct {
	apiURL string
	client *http.Client
}

// NewTranslationService creates a new HTTP client for the upstream translation API.
func NewTranslationService(apiURL string, client *http.Client) TranslationService {
	return &translationService{
		apiURL: apiURL,
		client: client,
	}
}

func (s *translationService) Translate(ctx context.Context, req *model.TranslationRequest) (*model.TranslationResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("upstream request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream returned status %d", resp.StatusCode)
	}

	var translation model.TranslationResponse
	if err := json.NewDecoder(resp.Body).Decode(&translation); err != nil {
		return nil, fmt.Errorf("decode upstream response: %w", err)
	}

	return &translation, nil
}
