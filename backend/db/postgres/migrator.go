package postgres

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrations embed.FS

func NewEmbeddedMigrator(url string) (*Migrator, error) {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithSourceInstance("iofs", source, url)
	if err != nil {
		source.Close()
		return nil, err
	}
	return &Migrator{m: m}, nil
}

type Migrator struct {
	m *migrate.Migrate
}

func NewMigrator(migrationPath string, url string) (*Migrator, error) {
	m, err := migrate.New(migrationPath, url)
	if err != nil {
		return nil, err
	}
	return &Migrator{m: m}, nil
}

func (m *Migrator) Up() error {
	if err := m.m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("ошибка в поднятии миграций %w", err)
	}
	return nil
}

func (m *Migrator) Close() error {
	sourceErr, databaseErr := m.m.Close()

	if err := errors.Join(sourceErr, databaseErr); err != nil {
		return fmt.Errorf("ошибка закрытия мигратора: %w", err)
	}
	return nil
}
