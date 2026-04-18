package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rodziievskyi-maksym/grn-subscription-service/internal/config"
	"github.com/stretchr/testify/require"
)

func createTempEnvFile(t *testing.T) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), ".env")
	err := os.WriteFile(path, []byte(""), 0o600)
	require.NoError(t, err)

	return path
}

func setRequiredEnv(t *testing.T) {
	t.Helper()

	t.Setenv("API_KEY", "secret-api-key")
	t.Setenv("SMTP_HOST", "smtp.gmail.com")
	t.Setenv("SMTP_PORT", "587")
	t.Setenv("SMTP_USER", "noreply@test.com")
	t.Setenv("SMTP_PASS", "password")
	t.Setenv("POSTGRES_DSN", "postgres://test:test@localhost:5432/test?sslmode=disable")
	t.Setenv("GITHUB_TOKEN", "ghp_test")
}

func TestLoad_AppliesDefaultsAndOverrides(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("REDIS_DB", "3")
	t.Setenv("REDIS_CACHE_TTL", "30m")

	cfg, err := config.Load(validator.New(), createTempEnvFile(t))
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, "127.0.0.1", cfg.Host)
	require.Equal(t, "8080", cfg.Port)
	require.Equal(t, "development", cfg.Env)
	require.Equal(t, 3, cfg.RedisDB)
	require.Equal(t, 30*time.Minute, cfg.RedisCacheTTL)
	require.Equal(t, 5*time.Minute, cfg.ScannerInterval)
	require.Equal(t, "password", cfg.SMTPPass)
}

func TestLoad_ReturnsErrorOnInvalidDuration(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("SCANNER_INTERVAL", "not-a-duration")

	_, err := config.Load(validator.New(), createTempEnvFile(t))
	require.Error(t, err)
	require.ErrorContains(t, err, "failed to parse environment variables")
	require.ErrorContains(t, err, "SCANNER_INTERVAL")
}

func TestLoad_ReturnsErrorOnInvalidInteger(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("REDIS_DB", "abc")

	_, err := config.Load(validator.New(), createTempEnvFile(t))
	require.Error(t, err)
	require.ErrorContains(t, err, "failed to parse environment variables")
	require.ErrorContains(t, err, "REDIS_DB")
}
