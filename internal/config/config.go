package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Host            string        `env:"HOST" default:"0.0.0.0" validate:"required"`
	Port            string        `env:"PORT" default:"8080" validate:"required"`
	Env             string        `env:"ENV" default:"development" validate:"required"`
	ScannerInterval time.Duration `env:"SCANNER_INTERVAL" default:"5m" validate:"required"`
	APIKey          string        `env:"API_KEY" validate:"required"`
	TestBool        bool          `env:"TEST_BOOL" default:"false"`

	RedisHost     string        `env:"REDIS_HOST" default:"localhost" validate:"required"`
	RedisPort     string        `env:"REDIS_PORT" default:"6379" validate:"required"`
	RedisPass     string        `env:"REDIS_PASS" default:""`
	RedisDB       int           `env:"REDIS_DB" default:"0"`
	RedisCacheTTL time.Duration `env:"REDIS_CACHE_TTL" default:"10m" validate:"required"`

	SMTPHost string `env:"SMTP_HOST" validate:"required"`
	SMTPPort string `env:"SMTP_PORT" validate:"required"`
	SMTPUser string `env:"SMTP_USER" validate:"required"`
	SMTPPass string `env:"SMTP_PASS" validate:"required"`
	SMTPFrom string `env:"SMTP_FROM"`

	PostgresDSN string `env:"POSTGRES_DSN" validate:"required"`

	GitHubToken string `env:"GITHUB_TOKEN" validate:"required"`
}

var durationType = reflect.TypeOf(time.Duration(0))

func loadFromEnv(cfg *Config) error {
	value := reflect.ValueOf(cfg).Elem()
	typ := value.Type()
	errs := make([]error, 0)

	for i := range typ.NumField() {
		field := typ.Field(i)

		envKey := field.Tag.Get("env")
		if envKey == "" {
			continue
		}

		raw, exists := os.LookupEnv(envKey)
		// Keep defaults when env var is empty.
		if !exists || raw == "" {
			continue
		}

		if err := setFieldValue(value.Field(i), raw); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", envKey, err))
		}
	}

	return errors.Join(errs...)
}

func setFieldValue(field reflect.Value, raw string) error {
	fieldType := field.Type()

	if fieldType == durationType {
		duration, err := time.ParseDuration(raw)
		if err != nil {
			return fmt.Errorf("invalid duration: %w", err)
		}

		field.SetInt(int64(duration))

		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(raw, 10, fieldType.Bits())
		if err != nil {
			return fmt.Errorf("invalid integer: %w", err)
		}

		field.SetInt(parsed)
	case reflect.Bool:
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			return fmt.Errorf("invalid boolean: %w", err)
		}

		field.SetBool(parsed)
	case reflect.Invalid, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.Array,
		reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.Struct,
		reflect.UnsafePointer:
		return fmt.Errorf("unsupported field type %s", fieldType)
	}

	return nil
}

func Load(validate *validator.Validate, envPath ...string) (*Config, error) {
	if validate == nil {
		return nil, errors.New("validator is nil")
	}

	if err := godotenv.Load(envPath...); err != nil {
		return nil, errors.Join(errors.New("failed to load environment variables from .env file"), err)
	}

	cfg := &Config{}

	if err := defaults.Set(cfg); err != nil {
		return nil, errors.Join(errors.New("failed to apply default configuration values"), err)
	}

	if err := loadFromEnv(cfg); err != nil {
		return nil, errors.Join(errors.New("failed to parse environment variables"), err)
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, errors.Join(errors.New("configuration validation failed"), err)
	}

	return cfg, nil
}

func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

func (c *Config) ServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c *Config) RedisAddress() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}
