package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"traci/backend/internal/config"
	"traci/backend/internal/infrastructure"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("application stopped", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg := config.NewConfig()
	if err := cfg.Load(); err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	app, err := infrastructure.NewApp(ctx, *cfg, logger)
	if err != nil {
		return err
	}
	defer app.Close()
	return app.Run(ctx)
}
