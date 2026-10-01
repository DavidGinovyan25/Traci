package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countMovies = `-- name: CountMovies :one
SELECT COUNT(*)
FROM movies
WHERE ($1::movie_genre IS NULL OR genre = $1::movie_genre)
  AND name ILIKE '%' || $2::text || '%'
`

type CountMoviesParams struct {
	Genre  NullMovieGenre `json:"genre"`
	Search string         `json:"search"`
}

func (q *Queries) CountMovies(ctx context.Context, arg CountMoviesParams) (int64, error) {
	row := q.db.QueryRow(ctx, countMovies, arg.Genre, arg.Search)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createMovie = `-- name: CreateMovie :exec
INSERT INTO movies (
    id, name, duration_min, genre, release_year, description, created_at, updated_at
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7, $8
)
`

type CreateMovieParams struct {
	ID          pgtype.UUID        `json:"id"`
	Name        string             `json:"name"`
	DurationMin int32              `json:"duration_min"`
	Genre       MovieGenre         `json:"genre"`
	ReleaseYear int32              `json:"release_year"`
	Description pgtype.Text        `json:"description"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	UpdatedAt   pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) CreateMovie(ctx context.Context, arg CreateMovieParams) error {
	_, err := q.db.Exec(ctx, createMovie,
		arg.ID,
		arg.Name,
		arg.DurationMin,
		arg.Genre,
		arg.ReleaseYear,
		arg.Description,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const deleteMovie = `-- name: DeleteMovie :execrows
DELETE FROM movies
WHERE id = $1
`

func (q *Queries) DeleteMovie(ctx context.Context, id pgtype.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteMovie, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getMovieByID = `-- name: GetMovieByID :one
SELECT id, name, duration_min, genre, release_year, description, created_at, updated_at
FROM movies
WHERE id = $1
`

func (q *Queries) GetMovieByID(ctx context.Context, id pgtype.UUID) (Movie, error) {
	row := q.db.QueryRow(ctx, getMovieByID, id)
	var i Movie
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.DurationMin,
		&i.Genre,
		&i.ReleaseYear,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listMovies = `-- name: ListMovies :many
SELECT id, name, duration_min, genre, release_year, description, created_at, updated_at
FROM movies
WHERE ($1::movie_genre IS NULL OR genre = $1::movie_genre)
  AND name ILIKE '%' || $2::text || '%'
ORDER BY created_at DESC, id ASC
LIMIT $4 OFFSET $3
`

type ListMoviesParams struct {
	Genre      NullMovieGenre `json:"genre"`
	Search     string         `json:"search"`
	PageOffset int32          `json:"page_offset"`
	PageLimit  int32          `json:"page_limit"`
}

func (q *Queries) ListMovies(ctx context.Context, arg ListMoviesParams) ([]Movie, error) {
	rows, err := q.db.Query(ctx, listMovies,
		arg.Genre,
		arg.Search,
		arg.PageOffset,
		arg.PageLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Movie
	for rows.Next() {
		var i Movie
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.DurationMin,
			&i.Genre,
			&i.ReleaseYear,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateMovie = `-- name: UpdateMovie :execrows
UPDATE movies
SET name = $1,
    duration_min = $2,
    genre = $3,
    release_year = $4,
    description = $5,
    updated_at = now()
WHERE id = $6
`

type UpdateMovieParams struct {
	Name        string      `json:"name"`
	DurationMin int32       `json:"duration_min"`
	Genre       MovieGenre  `json:"genre"`
	ReleaseYear int32       `json:"release_year"`
	Description pgtype.Text `json:"description"`
	ID          pgtype.UUID `json:"id"`
}

func (q *Queries) UpdateMovie(ctx context.Context, arg UpdateMovieParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateMovie,
		arg.Name,
		arg.DurationMin,
		arg.Genre,
		arg.ReleaseYear,
		arg.Description,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
