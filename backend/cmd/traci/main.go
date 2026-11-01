package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"traci/backend/db/postgres"
	"traci/backend/internal/application"
	"traci/backend/internal/config"
	"traci/backend/internal/infrastructure"
	httptransport "traci/backend/internal/transport/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.NewConfig()
	if err := cfg.Load(); err != nil {
		fmt.Fprintln(os.Stderr, "configuration loading failed:", err)
		os.Exit(1)
	}
	logger := newLogger(cfg)
	if err := run(logger, cfg); err != nil {
		var serverErr *httptransport.ServerError
		if errors.As(err, &serverErr) {
			logger.Error(
				"Application stopped with error",
				"error", err,
				"server.address", serverErr.Host,
				"server.port", serverErr.Port,
			)
		} else {
			logger.Error(
				"Application stopped with error",
				"error", err,
			)
		}
		os.Exit(1)
	}
}

func getMetrics() (*httptransport.HTTPMetrics, *application.BusinessMetrics) {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	httpMetrics := httptransport.NewHTTPMetrics(reg)
	businessMetrics := application.NewBusinessMetrics(reg)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	server := &http.Server{
		Addr: ":2112",
		Handler: mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Metrics server stopped", "error", err)
		}
	}()
	return httpMetrics, businessMetrics
}

func newLogger(cfg *config.Config) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
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
	})).With(
		"service.name", cfg.ServiceName,
		"service.version", "v1",
		"k8s.deployment.name", cfg.DeploymentName,
		"k8s.pod.name", cfg.PodName,
		"k8s.pod.uid", cfg.PodUID,
		"k8s.namespace.name", cfg.NamespaceName,
		"k8s.node.name", cfg.NodeName,
	)
}

func run(logger *slog.Logger, cfg *config.Config) (runErr error) {
	command := "serve"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if len(os.Args) > 2 || (command != "serve" && command != "create-admin" && command != "migrate") {
		return fmt.Errorf("usage: traci [serve|create-admin|migrate]")
	}
	if command == "migrate" {
		logger.Info("Database migration started", "db.system.name", "postgresql")
		m, err := postgres.NewEmbeddedMigrator(cfg.PostgresURL, logger)
		if err != nil {
			return fmt.Errorf("initialize database migrator: %w", err)
		}
		defer func() {
			if err := m.Close(); err != nil {
				runErr = errors.Join(runErr, err)
			}
		}()
		return m.Up()
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if command == "create-admin" {
		input, err := readAdminInput(os.Stdin, os.Stdout)
		if err != nil {
			return err
		}
		if _, err := infrastructure.CreateAdmin(ctx, *cfg, input, logger); err != nil {
			return fmt.Errorf("create administrator: %w", err)
		}
		return nil
	}
	httpMetrics, businessMetrics := getMetrics()
	logger.Info("Application starting")
	app, err := infrastructure.NewApp(ctx, *cfg, logger, httpMetrics, businessMetrics)
	if err != nil {
		return fmt.Errorf("initialize application: %w", err)
	}
	defer app.Close()
	if err := app.Run(ctx); err != nil {
		return err
	}
	logger.Info("Application stopped")
	return nil
}
