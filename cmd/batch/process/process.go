package process

import (
	"cached-translation-middleware/config"
	"cached-translation-middleware/internal/model"
	"cached-translation-middleware/internal/service"
	"cached-translation-middleware/internal/util"
	"context"
	"sync"

	"go.uber.org/zap"
)

func Process(logger *zap.Logger, configuration *config.Config, githubService service.GithubService, middlewareService service.MiddlewareService) {
	repos := getRepos(logger, configuration.Github, githubService)

	checkAndUpdateTranslation(logger, configuration.Translation, middlewareService, repos)
}

func getRepos(logger *zap.Logger, githubConfig config.GithubConfig, githubService service.GithubService) model.ListUserReposResponse {
	reposResponsesChan := make(chan model.ListUserReposResponse, len(githubConfig.OrgsLogin)+1)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), githubConfig.Timeout)
		defer cancel()

		reposResponse, err := githubService.GetRepos(ctx, model.UserTypeUser, githubConfig.UserLogin)
		if err != nil {
			logger.Warn("failed to get user repos", zap.String("username", githubConfig.UserLogin), zap.Error(err))
		}
		if reposResponse != nil {
			reposResponsesChan <- *reposResponse
		} else {
			reposResponsesChan <- model.ListUserReposResponse{}
		}

	}()

	for _, orgLogin := range githubConfig.OrgsLogin {
		wg.Add(1)
		go func(login string) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), githubConfig.Timeout)
			defer cancel()

			reposResponse, err := githubService.GetRepos(ctx, model.UserTypeOrg, login)
			if err != nil {
				logger.Warn("failed to get org repos", zap.String("org", login), zap.Error(err))
			}
			if reposResponse != nil {
				reposResponsesChan <- *reposResponse
			} else {
				reposResponsesChan <- model.ListUserReposResponse{}
			}
		}(orgLogin)
	}
	wg.Wait()
	close(reposResponsesChan)

	var repos model.ListUserReposResponse
	for response := range reposResponsesChan {
		repos = append(repos, response...)
	}
	return repos
}

func checkAndUpdateTranslation(logger *zap.Logger, translationConfig config.TranslationConfig, middlewareService service.MiddlewareService, repos model.ListUserReposResponse) {
	batchSize := 4

	reposLen := len(repos)
	for _, targetLanguage := range model.TargetLanguages {
		for i := 0; i < reposLen; i += batchSize {
			end := i + batchSize
			if end > reposLen {
				end = reposLen
			}
			batch := repos[i:end]

			var wg sync.WaitGroup
			wg.Add(len(batch))

			for _, repository := range batch {
				go func(repo model.RepoItem) {
					defer wg.Done()

					ctx, cancel := context.WithTimeout(context.Background(), translationConfig.Timeout)
					defer cancel()

					if repository.Description == nil {
						return
					}
					description := util.GetTextWithoutLinks(*repository.Description)
					req := model.TranslationRequest{Q: description, Source: string(model.SourceLanguage), Target: string(targetLanguage)}

					middlewareService.CheckAndUpdateTranslation(ctx, &req)
				}(repository)
			}
			wg.Wait()
		}
	}
}
