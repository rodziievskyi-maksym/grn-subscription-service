package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var (
	instance *Config
	once     sync.Once
	errInit  error
)

type Config struct {
	Host            string        `env:"HOST" default:"0.0.0.0" validate:"required"`
	Port            string        `env:"PORT" default:"8080" validate:"required"`
	Env             string        `env:"ENV" default:"development" validate:"required"`
	ScannerInterval time.Duration `env:"SCANNER_INTERVAL" default:"5m" validate:"required"`
	APIKey          string        `env:"API_KEY" validate:"required"`

	RedisHost     string        `env:"REDIS_HOST" default:"localhost" validate:"required"`
	RedisPort     string        `env:"REDIS_PORT" default:"6379" validate:"required"`
	RedisPass     string        `env:"REDIS_PASS" default:""`
	RedisDB       int           `env:"REDIS_DB" default:"0"`
	RedisCacheTTL time.Duration `env:"REDIS_CACHE_TTL" default:"10m" validate:"required"`

	SMTPHost string `env:"SMTP_HOST" validate:"required"`
	SMTPPort string `env:"SMTP_PORT" validate:"required"`
	SMTPUser string `env:"SMTP_USER" validate:"required"`
	SMTPPass string `env:"SMTP_PASSW" validate:"required"`
	SMTPFrom string `env:"SMTP_FROM"`

	PostgresDSN string `env:"POSTGRES_DSN" validate:"required"`

	GitHubToken string `env:"GITHUB_TOKEN" validate:"required"`
}

func parseTimeDuration(value string) time.Duration {
	duration, err := time.ParseDuration(value)
	if err != nil {
		panic(err)
	}

	return duration
}

func NewConfig(validator *validator.Validate, envPath ...string) error {
	once.Do(func() {
		slog.Info("Initializing configuration...")

		if err := godotenv.Load(envPath...); err != nil {
			errInit = errors.Join(errors.New("failed to load environment variables from .env file"), err)
			return
		}

		redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

		cfg := &Config{
			Host:            os.Getenv("HOST"),
			Port:            os.Getenv("PORT"),
			Env:             os.Getenv("ENV"),
			ScannerInterval: parseTimeDuration(os.Getenv("SCANNER_INTERVAL")),
			APIKey:          os.Getenv("API_KEY"),

			RedisHost:     os.Getenv("REDIS_HOST"),
			RedisPort:     os.Getenv("REDIS_PORT"),
			RedisPass:     os.Getenv("REDIS_PASS"),
			RedisDB:       redisDB,
			RedisCacheTTL: parseTimeDuration(os.Getenv("REDIS_CACHE_TTL")),

			SMTPHost: os.Getenv("SMTP_HOST"),
			SMTPPort: os.Getenv("SMTP_PORT"),
			SMTPUser: os.Getenv("SMTP_USER"),
			SMTPPass: os.Getenv("SMTP_PASS"),
			SMTPFrom: os.Getenv("SMTP_FROM"),

			PostgresDSN: os.Getenv("POSTGRES_DSN"),

			GitHubToken: os.Getenv("GITHUB_TOKEN"),
		}

		if err := defaults.Set(cfg); err != nil {
			errInit = errors.Join(errors.New("failed to apply default configuration values"), err)
			return
		}

		if err := validator.Struct(cfg); err != nil {
			errInit = errors.Join(errors.New("configuration validation failed"), err)
			return
		}

		instance = cfg
		errInit = nil

		slog.Info("Configuration initialized")
	})

	if errInit != nil {
		return errInit
	}

	return nil
}

func Cfg() *Config {
	return instance
}

func IsProduction() bool {
	return instance.Env == "production"
}

func GetServerAddress() string {
	return fmt.Sprintf("%s:%s", instance.Host, instance.Port)
}

func GetRedisAddress() string {
	return fmt.Sprintf("%s:%s", instance.RedisHost, instance.RedisPort)
}
