package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEnvironmentWithoutFile(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("POSTGRES_URL", "postgres://example/database")
	t.Setenv("HTTP_ADDR", ":9090")
	t.Setenv("SHUTDOWN_TIMEOUT", "3s")
	cfg := NewConfig()
	if err := cfg.Load(); err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPAddr != ":9090" || cfg.ShutdownTimeout != 3*time.Second {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestEnvironmentOverridesFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("POSTGRES_URL=postgres://file/database\nHTTP_ADDR=:8000\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("POSTGRES_URL", "postgres://environment/database")
	t.Setenv("HTTP_ADDR", ":9000")
	cfg := NewConfig()
	if err := cfg.Load(); err != nil {
		t.Fatal(err)
	}
	if cfg.PostgresURL != "postgres://environment/database" || cfg.HTTPAddr != ":9000" {
		t.Fatalf("environment did not override file")
	}
}
