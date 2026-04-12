package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
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

	SMTPHost string `env:"SMTP_HOST" validate:"required"`
	SMTPPort string `env:"SMTP_PORT" validate:"required"`
	SMTPUser string `env:"SMTP_USER" validate:"required"`
	SMTPPass string `env:"SMTP_PASSW" validate:"required"`
	SMTPFrom string `env:"SMTP_FROM"`

	PostgresDSN string `env:"POSTGRES_DSN" validate:"required"`

	GitHubToken string `env:"GITHUB_TOKEN" validate:"required"`
}

func NewConfig(validator *validator.Validate, envPath ...string) error {
	once.Do(func() {
		slog.Info("Initializing configuration...")

		if err := godotenv.Load(envPath...); err != nil {
			errInit = errors.Join(errors.New("failed to load environment variables from .env file"), err)
			return
		}

		scannerInterval, err := time.ParseDuration(os.Getenv("SCANNER_INTERVAL"))
		if err != nil {
			errInit = errors.Join(errors.New("failed to parse SCANNER_INTERVAL"), err)
			return
		}

		cfg := &Config{
			Host:            os.Getenv("HOST"),
			Port:            os.Getenv("PORT"),
			Env:             os.Getenv("ENV"),
			ScannerInterval: scannerInterval,

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

func IsDevelopment() bool {
	return instance.Env == "development"
}

func IsProduction() bool {
	return instance.Env == "production"
}

func GetServerAddress() string {
	return fmt.Sprintf("%s:%s", instance.Host, instance.Port)
}
