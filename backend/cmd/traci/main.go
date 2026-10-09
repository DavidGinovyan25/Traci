package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"traci/backend/db/postgres"
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
	command := "serve"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if len(os.Args) > 2 || (command != "serve" && command != "create-admin" && command != "migrate") {
		return fmt.Errorf("usage: traci [serve|create-admin|migrate]")
	}
	cfg := config.NewConfig()
	if err := cfg.Load(); err != nil {
		return err
	}
	if command == "migrate" {
		m, err := postgres.NewEmbeddedMigrator(cfg.PostgresURL)
		if err != nil {
			return err
		}
		defer m.Close()
		return m.Up()
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if command == "create-admin" {
		input, err := readAdminInput(os.Stdin, os.Stdout)
		if err != nil {
			return err
		}
		user, err := infrastructure.CreateAdmin(ctx, *cfg, input)
		if err != nil {
			return err
		}
		logger.Info("administrator created", "id", user.ID, "email", user.Email)
		return nil
	}
	app, err := infrastructure.NewApp(ctx, *cfg, logger)
	if err != nil {
		return err
	}
	defer app.Close()
	return app.Run(ctx)
}
