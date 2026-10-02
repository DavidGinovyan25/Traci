package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentWithoutFile(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("JWT_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("POSTGRES_URL", "postgres://example/database")
	t.Setenv("HTTP_ADDR", ":9090")
	t.Setenv("SHUTDOWN_TIMEOUT", "3s")
	cfg := NewConfig()
	require.NoError(t, cfg.Load())
	assert.Equal(t, ":9090", cfg.HTTPAddr, "unexpected config: %+v", cfg)
	assert.Equal(t, 3*time.Second, cfg.ShutdownTimeout, "unexpected config: %+v", cfg)
}

func TestEnvironmentOverridesFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte("POSTGRES_URL=postgres://file/database\nHTTP_ADDR=:8000\n"), 0600))
	t.Setenv("JWT_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("POSTGRES_URL", "postgres://environment/database")
	t.Setenv("HTTP_ADDR", ":9000")
	cfg := NewConfig()
	require.NoError(t, cfg.Load())
	assert.Equal(t, "postgres://environment/database", cfg.PostgresURL, "environment did not override file")
	assert.Equal(t, ":9000", cfg.HTTPAddr, "environment did not override file")
}

func TestInvalidConfigurationPreservesPreviousValue(t *testing.T) {
	for _, tc := range []struct{ name, key, value string }{
		{"missing database", "POSTGRES_URL", ""},
		{"short secret", "JWT_SECRET", "short"},
		{"missing secret", "JWT_SECRET", ""},
		{"invalid token duration", "TOKEN_TTL", "invalid"},
		{"short token duration", "TOKEN_TTL", "500ms"},
		{"zero token duration", "TOKEN_TTL", "0s"},
		{"negative connect timeout", "CONNECT_TIMEOUT", "-1s"},
		{"zero connect timeout", "CONNECT_TIMEOUT", "0s"},
		{"invalid connect timeout", "CONNECT_TIMEOUT", "invalid"},
		{"zero shutdown timeout", "SHUTDOWN_TIMEOUT", "0s"},
		{"negative shutdown timeout", "SHUTDOWN_TIMEOUT", "-1s"},
		{"invalid shutdown timeout", "SHUTDOWN_TIMEOUT", "invalid"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			t.Setenv("POSTGRES_URL", "postgres://example/database")
			t.Setenv("JWT_SECRET", "0123456789abcdef0123456789abcdef")
			t.Setenv("HTTP_ADDR", ":8080")
			t.Setenv("TOKEN_TTL", "1h")
			t.Setenv("CONNECT_TIMEOUT", "5s")
			t.Setenv("SHUTDOWN_TIMEOUT", "10s")
			cfg := NewConfig()
			require.NoError(t, cfg.Load())
			previous := *cfg
			t.Setenv(tc.key, tc.value)
			{
				err := cfg.Load()
				require.Error(t, err, "invalid configuration accepted")
			}
			assert.Equal(t, previous, *cfg, "failed reload changed current configuration")
		})
	}
}

func TestUnreadableConfigFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	require.NoError(t, os.Mkdir(filepath.Join(dir, ".env"), 0700))
	{
		err := NewConfig().Load()
		require.Error(t, err, "invalid config file accepted")
	}
}
