package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

type AuthBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func auth(cfg *Config) (string, error) {
	bodyAuth := &AuthBody{
		Email:    cfg.Email,
		Password: cfg.Password,
	}
	body, err := json.Marshal(bodyAuth)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(cfg.AuthURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication returned HTTP %d", resp.StatusCode)
	}

	var login LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&login); err != nil {
		return "", fmt.Errorf("decode auth response (HTTP %d, Content-Type %q): %w", resp.StatusCode, resp.Header.Get("Content-Type"), err)
	}
	if login.AccessToken == "" {
		return "", fmt.Errorf("authentication response contains no access token")
	}
	return login.AccessToken, nil
}

func main() {
	logger := newLogger()
	cfg := NewConfig()
	if err := cfg.Load(); err != nil {
		logger.Error("Configuration loading failed", "error", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(cfg.Workers)

	ctx := context.Background()

	started := time.Now()
	logger.InfoContext(ctx, "Probe run started", "probe.workers", cfg.Workers)

	logger.InfoContext(ctx, "Authentication started", "http.request.method", http.MethodPost)
	token, err := auth(cfg)
	if err != nil {
		logger.ErrorContext(ctx, "Authentication failed",
			"error.type", "authentication_error",
			"error.message", err.Error(),
		)
		return
	}
	logger.InfoContext(ctx, "Authentication completed")

	for range cfg.Workers {
		go func() {
			defer wg.Done()
			requestStarted := time.Now()
			requestLogger := logger.With("http.request.method", http.MethodGet)
			ctxT, cancel := context.WithTimeout(ctx, cfg.Timeout)
			defer cancel()
			r, err := http.NewRequestWithContext(ctxT, http.MethodGet, cfg.TargetURL, nil)
			r.Header.Set("Authorization", "Bearer "+token)
			if err != nil {
				requestLogger.ErrorContext(ctxT, "Failed to create probe request",
					"error.type", "request_creation_error",
				)
				return

			}
			requestLogger = requestLogger.With(
				"server.address", r.URL.Hostname(),
				"url.path", r.URL.Path,
			)
			resp, err := http.DefaultClient.Do(r)
			if err != nil {
				if errors.Is(err, http.ErrHandlerTimeout) {
					requestLogger.WarnContext(ctxT, "Probe request timed out",
						"error.type", "timeout",
						"duration_ms", float64(time.Since(requestStarted).Microseconds())/1000,
					)
				} else {
					requestLogger.WarnContext(ctxT, "Probe request failed",
						"error.type", "request_error",
						"duration_ms", float64(time.Since(requestStarted).Microseconds())/1000,
					)
				}
				return
			}
			requestLogger.InfoContext(ctxT, "Probe response received",
				"http.response.status_code", resp.StatusCode,
				"duration_ms", float64(time.Since(requestStarted).Microseconds())/1000,
			)
			defer resp.Body.Close()
		}()
	}
	wg.Wait()
	logger.InfoContext(ctx, "Probe run completed",
		"probe.workers", cfg.Workers,
		"duration_ms", float64(time.Since(started).Microseconds())/1000,
	)
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey && len(groups) == 0 {
				a.Key = "severity_text"
			}
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.String(
					"timestamp",
					a.Value.Time().UTC().Format("2006-01-02T15:04:05.000Z"),
				)
			}
			if a.Key == slog.MessageKey && len(groups) == 0 {
				a.Key = "message"
			}
			return a
		},
	})).With("service.name", "traci-prober")
}
