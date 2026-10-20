package postgres

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Migrator struct {
	m      *migrate.Migrate
	logger *slog.Logger
}

func NewEmbeddedMigrator(url string, logger *slog.Logger) (*Migrator, error) {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("create migration source: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", source, url)
	if err != nil {
		source.Close()
		return nil, fmt.Errorf("create embedded migrator: %w", err)
	}
	return &Migrator{m: m, logger: logger}, nil
}

func NewMigrator(migrationPath string, url string, logger *slog.Logger) (*Migrator, error) {
	m, err := migrate.New(migrationPath, url)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	return &Migrator{m: m, logger: logger}, nil
}

func (m *Migrator) Up() error {
	err := m.m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		m.logger.Info(
			"Database migrations already up to date",
			"db.system.name", "postgresql",
		)
		return nil
	}
	if err != nil {
		return fmt.Errorf("apply database migrations: %w", err)
	}
	m.logger.Info(
		"Database migrations applied successfully",
		"db.system.name", "postgresql",
	)

	return nil
}

func (m *Migrator) Close() error {
	sourceErr, databaseErr := m.m.Close()
	if err := errors.Join(sourceErr, databaseErr); err != nil {
		return fmt.Errorf("close migrator: %w", err)
	}
	return nil
}
