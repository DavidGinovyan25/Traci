package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, url string, logger *slog.Logger) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		logger.ErrorContext(
			ctx,
			"PostgreSQL pool creation failed",
			"db.system.name", "postgresql",
			"error", err,
		)
		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		logger.ErrorContext(
			ctx,
			"PostgreSQL connection failed",
			"db.system.name", "postgresql",
			"error", err,
		)
		return nil, fmt.Errorf("connect to PostgreSQL: %w", err)
	}
	logger.InfoContext(
		ctx,
		"PostgreSQL connection established",
		"db.system.name", "postgresql",
	)
	return &Postgres{pool: pool}, nil
}

func (p *Postgres) Close() {
	p.pool.Close()
}

func (p *Postgres) Pool() *pgxpool.Pool {
	return p.pool
}
