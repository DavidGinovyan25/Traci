package sql

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"traci/backend/internal/domain"
	"traci/backend/internal/gen/db"
)

type MovieRepository struct{ queries db.Querier }

var _ domain.MovieRepository = (*MovieRepository)(nil)

func NewMovieRepository(queries db.Querier) *MovieRepository {
	return &MovieRepository{queries: queries}
}

func (r *MovieRepository) Create(ctx context.Context, movie *domain.Movie) error {
	duration, err := databaseInt(movie.DurationMin)
	if err != nil {
		return fmt.Errorf("create movie: %w", err)
	}
	year, err := databaseInt(movie.ReleaseYear)
	if err != nil {
		return fmt.Errorf("create movie: %w", err)
	}
	err = r.queries.CreateMovie(ctx, db.CreateMovieParams{
		ID: databaseUUID(movie.ID), Name: movie.Name, DurationMin: duration,
		Genre: db.MovieGenre(movie.Genre), ReleaseYear: year, Description: databaseText(movie.Description),
		CreatedAt: databaseTime(movie.CreatedAt), UpdatedAt: databaseTime(movie.UpdatedAt),
	})
	return repositoryError("create movie", err)
}

func (r *MovieRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Movie, error) {
	row, err := r.queries.GetMovieByID(ctx, databaseUUID(id))
	if err != nil {
		return nil, repositoryError("get movie by id", err)
	}
	movie := movieFromRow(row)
	return &movie, nil
}

func (r *MovieRepository) List(ctx context.Context, filter domain.MovieFilter) ([]domain.Movie, int, error) {
	l, o, err := pageParams(filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list movies: %w", err)
	}
	genre := db.NullMovieGenre{}
	if filter.Genre != nil {
		genre = db.NullMovieGenre{MovieGenre: db.MovieGenre(*filter.Genre), Valid: true}
	}
	rows, err := r.queries.ListMovies(ctx, db.ListMoviesParams{
		Genre: genre, Search: filter.Search, PageLimit: l, PageOffset: o,
	})
	if err != nil {
		return nil, 0, repositoryError("list movies", err)
	}
	total, err := r.queries.CountMovies(ctx, db.CountMoviesParams{Genre: genre, Search: filter.Search})
	if err != nil {
		return nil, 0, repositoryError("count movies", err)
	}
	movies := make([]domain.Movie, len(rows))
	for i, row := range rows {
		movies[i] = movieFromRow(row)
	}
	return movies, int(total), nil
}

func (r *MovieRepository) Update(ctx context.Context, movie *domain.Movie) error {
	duration, err := databaseInt(movie.DurationMin)
	if err != nil {
		return fmt.Errorf("update movie: %w", err)
	}
	year, err := databaseInt(movie.ReleaseYear)
	if err != nil {
		return fmt.Errorf("update movie: %w", err)
	}
	rows, err := r.queries.UpdateMovie(ctx, db.UpdateMovieParams{
		ID: databaseUUID(movie.ID), Name: movie.Name, DurationMin: duration,
		Genre: db.MovieGenre(movie.Genre), ReleaseYear: year, Description: databaseText(movie.Description),
	})
	return affectedRows("update movie", rows, err)
}

func (r *MovieRepository) Delete(ctx context.Context, id uuid.UUID) error {
	rows, err := r.queries.DeleteMovie(ctx, databaseUUID(id))
	return affectedRows("delete movie", rows, err)
}

func movieFromRow(row db.Movie) domain.Movie {
	return domain.Movie{ID: uuid.UUID(row.ID.Bytes), Name: row.Name, DurationMin: int(row.DurationMin),
		Genre: domain.MovieGenre(row.Genre), ReleaseYear: int(row.ReleaseYear),
		Description: textPointer(row.Description), CreatedAt: row.CreatedAt.Time, UpdatedAt: row.UpdatedAt.Time}
}
