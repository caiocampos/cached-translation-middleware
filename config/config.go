package config

import (
	"cached-translation-middleware/internal/util"
	"encoding/json"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App         AppConfig
	Redis       RedisConfig
	Translation TranslationConfig
	Github      GithubConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	CacheTTL time.Duration
}

type TranslationConfig struct {
	APIURL  string
	Timeout time.Duration
}

type GithubConfig struct {
	APIURL    string
	UserLogin string
	OrgsLogin []string
	Timeout   time.Duration
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("TRANSLATION_API_URL", "http://127.0.0.1:5000")
	viper.SetDefault("TRANSLATION_API_TIMEOUT", util.TwentySecondsString)
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_USER", "default")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_CACHE_TTL", util.SevenDaysString)
	viper.SetDefault("GITHUB_API_URL", "https://api.github.com")
	viper.SetDefault("GITHUB_API_TIMEOUT", util.TenSecondsString)
	viper.SetDefault("GITHUB_API_USER_LOGIN", "")
	viper.SetDefault("GITHUB_API_ORGS_LOGIN", "[]")

	// Ignore error if .env file not found (env vars may be set directly)
	_ = viper.ReadInConfig()

	translationTimeout, err := time.ParseDuration(viper.GetString("TRANSLATION_API_TIMEOUT"))
	if err != nil {
		translationTimeout = util.TwentySeconds
	}

	githubTimeout, err := time.ParseDuration(viper.GetString("GITHUB_API_TIMEOUT"))
	if err != nil {
		githubTimeout = util.TenSeconds
	}

	cacheTTL, err := time.ParseDuration(viper.GetString("REDIS_CACHE_TTL"))
	if err != nil {
		cacheTTL = util.SevenDays
	}

	var orgsLogin []string
	orgsLoginString := viper.GetString("GITHUB_API_ORGS_LOGIN")
	if orgsLoginString == "" {
		orgsLogin = []string{}
	} else {
		err := json.Unmarshal([]byte(orgsLoginString), &orgsLogin)
		if err != nil {
			orgsLogin = []string{}
		}
	}

	return &Config{
		App: AppConfig{
			Port: viper.GetString("APP_PORT"),
			Env:  viper.GetString("APP_ENV"),
		},
		Redis: RedisConfig{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
			CacheTTL: cacheTTL,
		},
		Translation: TranslationConfig{
			APIURL:  viper.GetString("TRANSLATION_API_URL"),
			Timeout: translationTimeout,
		},
		Github: GithubConfig{
			APIURL:    viper.GetString("GITHUB_API_URL"),
			Timeout:   githubTimeout,
			UserLogin: viper.GetString("GITHUB_API_USER_LOGIN"),
			OrgsLogin: orgsLogin,
		},
	}, nil
}
