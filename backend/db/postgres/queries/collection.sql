-- name: CreateCollectionItem :exec
INSERT INTO user_movies (
    id, user_id, movie_id, personal_rating, review, status, created_at, updated_at
) VALUES (
    sqlc.arg(id), sqlc.arg(user_id), sqlc.arg(movie_id),
    sqlc.narg(personal_rating), sqlc.narg(review), sqlc.arg(status),
    sqlc.arg(created_at), sqlc.arg(updated_at)
);

-- name: GetCollectionItem :one
SELECT id, user_id, movie_id, personal_rating, review, status, created_at, updated_at
FROM user_movies
WHERE user_id = sqlc.arg(user_id) AND movie_id = sqlc.arg(movie_id);

-- name: ListCollectionItems :many
SELECT id, user_id, movie_id, personal_rating, review, status, created_at, updated_at
FROM user_movies
WHERE user_id = sqlc.arg(user_id)
  AND (sqlc.narg(status)::movie_status IS NULL OR status = sqlc.narg(status)::movie_status)
ORDER BY created_at DESC, id ASC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountCollectionItems :one
SELECT COUNT(*)
FROM user_movies
WHERE user_id = sqlc.arg(user_id)
  AND (sqlc.narg(status)::movie_status IS NULL OR status = sqlc.narg(status)::movie_status);

-- name: UpdateCollectionItem :execrows
UPDATE user_movies
SET personal_rating = sqlc.narg(personal_rating),
    review = sqlc.narg(review),
    status = sqlc.arg(status),
    updated_at = now()
WHERE user_id = sqlc.arg(user_id) AND movie_id = sqlc.arg(movie_id);

-- name: DeleteCollectionItem :execrows
DELETE FROM user_movies
WHERE user_id = sqlc.arg(user_id) AND movie_id = sqlc.arg(movie_id);
