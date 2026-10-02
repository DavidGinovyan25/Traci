package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"traci/backend/db/postgres"
	"traci/backend/internal/testutil"
)

func TestPostgresPoolAndMigrations(t *testing.T) {
	url := testutil.Postgres(t)
	ctx := context.Background()
	database, err := postgres.NewPostgres(ctx, url)
	require.NoError(t, err)
	defer database.Close()
	var version int64
	var dirty bool
	require.NoError(t, database.Pool().QueryRow(ctx, "SELECT version, dirty FROM schema_migrations").Scan(&version, &dirty))
	assert.NotZero(t, version, "invalid migration state: version=%d dirty=%v", version, dirty)
	require.False(t, dirty, "invalid migration state: version=%d dirty=%v", version, dirty)
	for _, table := range []string{"users", "movies", "user_movies"} {
		var exists bool
		require.NoError(t, database.Pool().QueryRow(ctx, "SELECT to_regclass($1) IS NOT NULL", table).Scan(&exists))
		require.True(t, exists, "missing table %s", table)
	}
	var extension bool
	require.NoError(t, database.Pool().QueryRow(ctx, "SELECT EXISTS (SELECT FROM pg_extension WHERE extname = 'pg_trgm')").Scan(&extension))
	require.True(t, extension, "search extension missing")
}

func TestPostgresRejectsInvalidURL(t *testing.T) {
	{
		database, err := postgres.NewPostgres(context.Background(), "://invalid")
		require.Error(t, err, "got database=%v err=%v", database, err)
		require.Nil(t, database, "got database=%v err=%v", database, err)
	}
}
