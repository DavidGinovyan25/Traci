package infrastructure

import (
	"context"
	"log/slog"

	"traci/backend/db/postgres"
	"traci/backend/internal/application"
	"traci/backend/internal/config"
	"traci/backend/internal/domain"
	"traci/backend/internal/gen/db"
	"traci/backend/internal/infrastructure/auth"
	sqlrepository "traci/backend/internal/repository/sql"
	httptransport "traci/backend/internal/transport/http"
)

type App struct {
	server   *httptransport.Server
	database *postgres.Postgres
}

func NewApp(ctx context.Context, cfg config.Config, logger *slog.Logger) (*App, error) {
	credentials, err := auth.NewCredentials(cfg.JWTSecret, cfg.TokenTTL)
	if err != nil {
		return nil, err
	}
	connectCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()
	database, err := postgres.NewPostgres(connectCtx, cfg.PostgresURL)
	if err != nil {
		return nil, err
	}
	queries := db.New(database.Pool())
	service := application.NewService(
		sqlrepository.NewUserRepository(queries), sqlrepository.NewMovieRepository(queries),
		sqlrepository.NewCollectionRepository(queries), credentials,
	)
	handler := httptransport.NewHandler(service, logger, cfg.TokenTTL)
	router, err := httptransport.NewRouter(handler, logger)
	if err != nil {
		database.Close()
		return nil, err
	}
	return &App{server: httptransport.NewServer(cfg, router, logger), database: database}, nil
}

func (a *App) Run(ctx context.Context) error { return a.server.Run(ctx) }
func (a *App) Close()                        { a.database.Close() }

func CreateAdmin(ctx context.Context, cfg config.Config, input application.Register) (*domain.User, error) {
	credentials, err := auth.NewCredentials(cfg.JWTSecret, cfg.TokenTTL)
	if err != nil {
		return nil, err
	}
	connectCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()
	database, err := postgres.NewPostgres(connectCtx, cfg.PostgresURL)
	if err != nil {
		return nil, err
	}
	defer database.Close()
	service := application.NewService(sqlrepository.NewUserRepository(db.New(database.Pool())), nil, nil, credentials)
	return service.CreateAdmin(ctx, input)
}
