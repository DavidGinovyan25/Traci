package postgres

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

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
