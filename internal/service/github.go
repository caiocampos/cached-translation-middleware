package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cached-translation-middleware/internal/model"
)

type UserType string

const (
	UserTypeUser UserType = "users"
	UserTypeOrg  UserType = "orgs"
)

// GithubService defines the contract for calling the upstream Github API.
type GithubService interface {
	GetRepos(ctx context.Context, userType UserType, login string) (*model.ListUserReposResponse, error)
}

type githubService struct {
	apiURL string
	client *http.Client
}

// NewGithubService creates a new HTTP client for the upstream Github API.
func NewGithubService(apiURL string, client *http.Client) GithubService {
	return &githubService{
		apiURL: apiURL,
		client: client,
	}
}

func (s *githubService) GetRepos(ctx context.Context, userType UserType, login string) (*model.ListUserReposResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.apiURL, nil)
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

	var Github model.ListUserReposResponse
	if err := json.NewDecoder(resp.Body).Decode(&Github); err != nil {
		return nil, fmt.Errorf("decode upstream response: %w", err)
	}

	return &Github, nil
}
