package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CountCollectionItems(ctx context.Context, arg CountCollectionItemsParams) (int64, error)
	CountMovies(ctx context.Context, arg CountMoviesParams) (int64, error)
	CountUsers(ctx context.Context) (int64, error)
	CreateCollectionItem(ctx context.Context, arg CreateCollectionItemParams) error
	CreateMovie(ctx context.Context, arg CreateMovieParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) error
	DeleteCollectionItem(ctx context.Context, arg DeleteCollectionItemParams) (int64, error)
	DeleteMovie(ctx context.Context, id pgtype.UUID) (int64, error)
	DeleteUser(ctx context.Context, id pgtype.UUID) (int64, error)
	GetCollectionItem(ctx context.Context, arg GetCollectionItemParams) (GetCollectionItemRow, error)
	GetMovieByID(ctx context.Context, id pgtype.UUID) (Movie, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id pgtype.UUID) (User, error)
	ListCollectionItems(ctx context.Context, arg ListCollectionItemsParams) ([]ListCollectionItemsRow, error)
	ListMovies(ctx context.Context, arg ListMoviesParams) ([]Movie, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateCollectionItem(ctx context.Context, arg UpdateCollectionItemParams) (int64, error)
	UpdateMovie(ctx context.Context, arg UpdateMovieParams) (int64, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (int64, error)
}

var _ Querier = (*Queries)(nil)
