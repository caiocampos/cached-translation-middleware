package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"cached-translation-middleware/internal/model"
)

// GithubService defines the contract for calling the upstream Github API.
type GithubService interface {
	GetRepos(ctx context.Context, userType model.UserType, login string) (*model.ListUserReposResponse, error)
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

func (s *githubService) GetRepos(ctx context.Context, userType model.UserType, login string) (*model.ListUserReposResponse, error) {
	url := strings.Join([]string{s.apiURL, string(userType), login, "repos"}, "/")

	log.Println(url)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	var repos model.ListUserReposResponse
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("decode upstream response: %w", err)
	}

	return &repos, nil
}
