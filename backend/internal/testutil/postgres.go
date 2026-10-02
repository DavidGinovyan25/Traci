package testutil

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"

	"traci/backend/db/postgres"
)

func Postgres(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("integration test requires Docker")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	container, err := pgcontainer.Run(ctx, "postgres:17-alpine", pgcontainer.WithDatabase("traci_test"), pgcontainer.WithUsername("traci_test"), pgcontainer.WithPassword("test_password"), pgcontainer.BasicWaitStrategies())
	testcontainers.CleanupContainer(t, container)
	require.NoError(t, err)
	url, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok, "cannot locate migrations")
	migrationPath := "file://" + filepath.Join(filepath.Dir(file), "../../db/postgres/migrations")
	migrator, err := postgres.NewMigrator(migrationPath, url)
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, migrator.Close())
	}()
	for range 2 {
		require.NoError(t, migrator.Up())
	}
	return url
}
