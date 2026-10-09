-- name: CreateMovie :exec
INSERT INTO movies (
    id, name, duration_min, genre, release_year, description, created_at, updated_at
) VALUES (
    sqlc.arg(id), sqlc.arg(name), sqlc.arg(duration_min), sqlc.arg(genre),
    sqlc.arg(release_year), sqlc.narg(description), sqlc.arg(created_at), sqlc.arg(updated_at)
);

-- name: GetMovieByID :one
SELECT id, name, duration_min, genre, release_year, description, created_at, updated_at
FROM movies
WHERE id = sqlc.arg(id);

-- name: ListMovies :many
SELECT id, name, duration_min, genre, release_year, description, created_at, updated_at
FROM movies
WHERE (sqlc.narg(genre)::movie_genre IS NULL OR genre = sqlc.narg(genre)::movie_genre)
  AND name ILIKE '%' || sqlc.arg(search)::text || '%'
ORDER BY created_at DESC, id ASC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountMovies :one
SELECT COUNT(*)
FROM movies
WHERE (sqlc.narg(genre)::movie_genre IS NULL OR genre = sqlc.narg(genre)::movie_genre)
  AND name ILIKE '%' || sqlc.arg(search)::text || '%';

-- name: UpdateMovie :execrows
UPDATE movies
SET name = sqlc.arg(name),
    duration_min = sqlc.arg(duration_min),
    genre = sqlc.arg(genre),
    release_year = sqlc.arg(release_year),
    description = sqlc.narg(description),
    updated_at = now()
WHERE id = sqlc.arg(id);

-- name: DeleteMovie :execrows
DELETE FROM movies
WHERE id = sqlc.arg(id);
