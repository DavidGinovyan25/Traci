package httptransport

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"traci/backend/internal/config"
)

func runTestServer(t *testing.T, handler http.Handler, timeout time.Duration) (context.CancelFunc, <-chan error, *Server, string) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	server := NewServer(config.Config{HTTPAddr: "127.0.0.1:0", ShutdownTimeout: timeout}, handler, slog.New(slog.NewTextHandler(io.Discard, nil)))
	ready := make(chan string, 1)
	server.server.BaseContext = func(listener net.Listener) context.Context {
		ready <- listener.Addr().String()
		return context.Background()
	}
	done := make(chan error, 1)
	go func() { done <- server.Run(ctx) }()
	t.Cleanup(func() { cancel(); server.server.Close() })
	select {
	case addr := <-ready:
		return cancel, done, server, "http://" + addr
	case err := <-done:
		require.FailNow(t, "test failed", "server failed to start: %v", err)
	case <-time.After(5 * time.Second):
		require.FailNow(t, "server startup timed out")
	}
	return nil, nil, nil, ""
}

func awaitServer(t *testing.T, done <-chan error) error {
	t.Helper()
	select {
	case err := <-done:
		return err
	case <-time.After(5 * time.Second):
		require.FailNow(t, "server did not stop")
		return nil
	}
}

func TestServerShutdownWaitsForActiveRequest(t *testing.T) {
	entered, release := make(chan struct{}), make(chan struct{})
	defer close(release)
	cancel, done, server, url := runTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { close(entered); <-release; w.Write([]byte("completed")) }), time.Second)
	response := make(chan error, 1)
	go func() {
		client := &http.Client{Timeout: 5 * time.Second}
		res, err := client.Get(url)
		if err == nil {
			defer res.Body.Close()
			var body []byte
			body, err = io.ReadAll(res.Body)
			if err == nil && string(body) != "completed" {
				err = errors.New("incomplete response")
			}
		}
		response <- err
	}()
	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		require.FailNow(t, "request not received")
	}
	shuttingDown := make(chan struct{})
	server.server.RegisterOnShutdown(func() { close(shuttingDown) })
	cancel()
	select {
	case <-shuttingDown:
	case <-time.After(5 * time.Second):
		require.FailNow(t, "shutdown not started")
	}
	select {
	case err := <-done:
		require.FailNow(t, "test failed", "shutdown returned before request completed: %v", err)
	default:
	}
	release <- struct{}{}
	select {
	case err := <-response:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		require.FailNow(t, "request did not complete")
	}
	require.NoError(t, awaitServer(t, done))
}

func TestServerShutdownTimeoutClosesConnections(t *testing.T) {
	entered, finished := make(chan struct{}), make(chan struct{})
	cancel, done, _, url := runTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { close(entered); <-r.Context().Done(); close(finished) }), 50*time.Millisecond)
	response := make(chan error, 1)
	go func() {
		res, err := (&http.Client{Timeout: 5 * time.Second}).Get(url)
		if res != nil {
			res.Body.Close()
		}
		response <- err
	}()
	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		require.FailNow(t, "request not received")
	}
	cancel()
	{
		err := awaitServer(t, done)
		require.ErrorIs(t, err, context.DeadlineExceeded, "expected shutdown timeout, got %v", err)
	}
	select {
	case <-finished:
	case <-time.After(5 * time.Second):
		require.FailNow(t, "active connection was not canceled")
	}
	select {
	case err := <-response:
		require.Error(t, err, "expected connection closure")
	case <-time.After(5 * time.Second):
		require.FailNow(t, "client did not finish")
	}
}

func TestServerBindFailure(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	defer listener.Close()
	server := NewServer(config.Config{HTTPAddr: listener.Addr().String(), ShutdownTimeout: time.Second}, http.NotFoundHandler(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	{
		err := server.Run(ctx)
		require.Error(t, err, "occupied port did not cause startup failure")
	}
}

func TestServerAlreadyCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	server := NewServer(config.Config{HTTPAddr: "127.0.0.1:0", ShutdownTimeout: time.Second}, http.NotFoundHandler(), slog.New(slog.NewTextHandler(io.Discard, nil)))
	require.NoError(t, server.Run(ctx))
}
