package infrastructure

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"traci/backend/internal/config"
	"traci/backend/internal/testutil"
)

func appTestConfig() config.Config {
	return config.Config{JWTSecret: strings.Repeat("a", 32), TokenTTL: time.Hour, ConnectTimeout: time.Second, ShutdownTimeout: time.Second, HTTPAddr: "127.0.0.1:0"}
}

func TestAppStartupErrors(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	t.Run("invalid credentials", func(t *testing.T) {
		cfg := appTestConfig()
		cfg.JWTSecret = "short"
		app, err := NewApp(context.Background(), cfg, logger)
		require.Error(t, err, "app=%v err=%v", app, err)
		require.Nil(t, app, "app=%v err=%v", app, err)
	})
	t.Run("invalid database URL", func(t *testing.T) {
		cfg := appTestConfig()
		cfg.PostgresURL = "://invalid"
		app, err := NewApp(context.Background(), cfg, logger)
		require.Error(t, err, "app=%v err=%v", app, err)
		require.Nil(t, app, "app=%v err=%v", app, err)
	})
	t.Run("canceled connection", func(t *testing.T) {
		cfg := appTestConfig()
		cfg.PostgresURL = "postgres://test:test@127.0.0.1:1/test?sslmode=disable"
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		app, err := NewApp(ctx, cfg, logger)
		require.Nil(t, app, "app=%v err=%v", app, err)
		require.ErrorIs(t, err, context.Canceled, "app=%v err=%v", app, err)
	})
}

func TestAppPostgresLifecycle(t *testing.T) {
	cfg := appTestConfig()
	cfg.PostgresURL = testutil.Postgres(t)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	cfg.HTTPAddr = listener.Addr().String()
	require.NoError(t, listener.Close())
	app, err := NewApp(context.Background(), cfg, slog.New(slog.NewTextHandler(io.Discard, nil)))
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- app.Run(ctx) }()
	stopped := false
	t.Cleanup(func() {
		cancel()
		if !stopped {
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				assert.Fail(t, "app did not stop")
			}
		}
		app.Close()
	})
	client := &http.Client{Timeout: time.Second}
	url := "http://" + cfg.HTTPAddr
	deadline := time.Now().Add(5 * time.Second)
	for {
		response, err := client.Get(url + "/healthz")
		if err == nil {
			response.Body.Close()
			assert.EqualValues(t, 200, response.StatusCode, "health HTTP %d", response.StatusCode)
			break
		}
		select {
		case err := <-done:
			stopped = true
			require.FailNow(t, "test failed", "app stopped during startup: %v", err)
		default:
		}
		require.False(t, time.Now().After(deadline), "app not ready: %v", err)
		time.Sleep(10 * time.Millisecond)
	}
	response, err := client.Post(url+"/auth/register", "application/json", strings.NewReader(`{"username":"lifecycle","first_name":"Test","second_name":"User","email":"lifecycle@example.com","password":"password123"}`))
	require.NoError(t, err)
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(response.Body)
		require.FailNow(t, "test failed", "registration HTTP %d: %s", response.StatusCode, body)
	}
	var user struct {
		ID string `json:"id"`
	}
	require.NoError(t, json.NewDecoder(response.Body).Decode(&user))
	var email string
	require.NoError(t, app.database.Pool().QueryRow(context.Background(), "SELECT email FROM users WHERE id = $1", user.ID).Scan(&email))
	assert.Equal(t, "lifecycle@example.com", email, "unexpected persisted user: %s", email)
	cancel()
	select {
	case err := <-done:
		stopped = true
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		require.FailNow(t, "app shutdown timed out")
	}
	app.Close()
	{
		count := app.database.Pool().Stat().TotalConns()
		assert.EqualValues(t, 0, count, "database connections remain: %d", count)
	}
}
