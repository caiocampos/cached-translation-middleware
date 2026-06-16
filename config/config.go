package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig
	Redis      RedisConfig
	Translation TranslationConfig
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

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("TRANSLATION_API_URL", "http://localhost:5000/translate")
	viper.SetDefault("TRANSLATION_API_TIMEOUT", "10s")
	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("REDIS_CACHE_TTL", "24h")

	// Ignore error if .env file not found (env vars may be set directly)
	_ = viper.ReadInConfig()

	translationTimeout, err := time.ParseDuration(viper.GetString("TRANSLATION_API_TIMEOUT"))
	if err != nil {
		translationTimeout = 10 * time.Second
	}

	cacheTTL, err := time.ParseDuration(viper.GetString("REDIS_CACHE_TTL"))
	if err != nil {
		cacheTTL = 24 * time.Hour
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
	}, nil
}
