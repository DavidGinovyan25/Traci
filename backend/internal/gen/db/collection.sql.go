package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countCollectionItems = `-- name: CountCollectionItems :one
SELECT COUNT(*)
FROM user_movies
WHERE user_id = $1
  AND ($2::movie_status IS NULL OR status = $2::movie_status)
`

type CountCollectionItemsParams struct {
	UserID pgtype.UUID     `json:"user_id"`
	Status NullMovieStatus `json:"status"`
}

func (q *Queries) CountCollectionItems(ctx context.Context, arg CountCollectionItemsParams) (int64, error) {
	row := q.db.QueryRow(ctx, countCollectionItems, arg.UserID, arg.Status)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createCollectionItem = `-- name: CreateCollectionItem :exec
INSERT INTO user_movies (
    id, user_id, movie_id, personal_rating, review, status, created_at, updated_at
) VALUES (
    $1, $2, $3,
    $4, $5, $6,
    $7, $8
)
`

type CreateCollectionItemParams struct {
	ID             pgtype.UUID        `json:"id"`
	UserID         pgtype.UUID        `json:"user_id"`
	MovieID        pgtype.UUID        `json:"movie_id"`
	PersonalRating pgtype.Int4        `json:"personal_rating"`
	Review         pgtype.Text        `json:"review"`
	Status         MovieStatus        `json:"status"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
	UpdatedAt      pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) CreateCollectionItem(ctx context.Context, arg CreateCollectionItemParams) error {
	_, err := q.db.Exec(ctx, createCollectionItem,
		arg.ID,
		arg.UserID,
		arg.MovieID,
		arg.PersonalRating,
		arg.Review,
		arg.Status,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const deleteCollectionItem = `-- name: DeleteCollectionItem :execrows
DELETE FROM user_movies
WHERE user_id = $1 AND movie_id = $2
`

type DeleteCollectionItemParams struct {
	UserID  pgtype.UUID `json:"user_id"`
	MovieID pgtype.UUID `json:"movie_id"`
}

func (q *Queries) DeleteCollectionItem(ctx context.Context, arg DeleteCollectionItemParams) (int64, error) {
	result, err := q.db.Exec(ctx, deleteCollectionItem, arg.UserID, arg.MovieID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getCollectionItem = `-- name: GetCollectionItem :one
SELECT id, user_id, movie_id, personal_rating, review, status, created_at, updated_at
FROM user_movies
WHERE user_id = $1 AND movie_id = $2
`

type GetCollectionItemParams struct {
	UserID  pgtype.UUID `json:"user_id"`
	MovieID pgtype.UUID `json:"movie_id"`
}

type GetCollectionItemRow struct {
	ID             pgtype.UUID        `json:"id"`
	UserID         pgtype.UUID        `json:"user_id"`
	MovieID        pgtype.UUID        `json:"movie_id"`
	PersonalRating pgtype.Int4        `json:"personal_rating"`
	Review         pgtype.Text        `json:"review"`
	Status         MovieStatus        `json:"status"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
	UpdatedAt      pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) GetCollectionItem(ctx context.Context, arg GetCollectionItemParams) (GetCollectionItemRow, error) {
	row := q.db.QueryRow(ctx, getCollectionItem, arg.UserID, arg.MovieID)
	var i GetCollectionItemRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.MovieID,
		&i.PersonalRating,
		&i.Review,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listCollectionItems = `-- name: ListCollectionItems :many
SELECT id, user_id, movie_id, personal_rating, review, status, created_at, updated_at
FROM user_movies
WHERE user_id = $1
  AND ($2::movie_status IS NULL OR status = $2::movie_status)
ORDER BY created_at DESC, id ASC
LIMIT $4 OFFSET $3
`

type ListCollectionItemsParams struct {
	UserID     pgtype.UUID     `json:"user_id"`
	Status     NullMovieStatus `json:"status"`
	PageOffset int32           `json:"page_offset"`
	PageLimit  int32           `json:"page_limit"`
}

type ListCollectionItemsRow struct {
	ID             pgtype.UUID        `json:"id"`
	UserID         pgtype.UUID        `json:"user_id"`
	MovieID        pgtype.UUID        `json:"movie_id"`
	PersonalRating pgtype.Int4        `json:"personal_rating"`
	Review         pgtype.Text        `json:"review"`
	Status         MovieStatus        `json:"status"`
	CreatedAt      pgtype.Timestamptz `json:"created_at"`
	UpdatedAt      pgtype.Timestamptz `json:"updated_at"`
}

func (q *Queries) ListCollectionItems(ctx context.Context, arg ListCollectionItemsParams) ([]ListCollectionItemsRow, error) {
	rows, err := q.db.Query(ctx, listCollectionItems,
		arg.UserID,
		arg.Status,
		arg.PageOffset,
		arg.PageLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListCollectionItemsRow
	for rows.Next() {
		var i ListCollectionItemsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.MovieID,
			&i.PersonalRating,
			&i.Review,
			&i.Status,
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

const updateCollectionItem = `-- name: UpdateCollectionItem :execrows
UPDATE user_movies
SET personal_rating = $1,
    review = $2,
    status = $3,
    updated_at = now()
WHERE user_id = $4 AND movie_id = $5
`

type UpdateCollectionItemParams struct {
	PersonalRating pgtype.Int4 `json:"personal_rating"`
	Review         pgtype.Text `json:"review"`
	Status         MovieStatus `json:"status"`
	UserID         pgtype.UUID `json:"user_id"`
	MovieID        pgtype.UUID `json:"movie_id"`
}

func (q *Queries) UpdateCollectionItem(ctx context.Context, arg UpdateCollectionItemParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateCollectionItem,
		arg.PersonalRating,
		arg.Review,
		arg.Status,
		arg.UserID,
		arg.MovieID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
