package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countUsers = `-- name: CountUsers :one
SELECT COUNT(*) FROM users
`

func (q *Queries) CountUsers(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countUsers)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users (
    id,
    role,
    username,
    first_name,
    second_name,
    email,
    password_hash,
    is_blocked
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
`

type CreateUserParams struct {
	ID           pgtype.UUID `json:"id"`
	Role         Role        `json:"role"`
	Username     string      `json:"username"`
	FirstName    string      `json:"first_name"`
	SecondName   string      `json:"second_name"`
	Email        string      `json:"email"`
	PasswordHash string      `json:"password_hash"`
	IsBlocked    bool        `json:"is_blocked"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser,
		arg.ID,
		arg.Role,
		arg.Username,
		arg.FirstName,
		arg.SecondName,
		arg.Email,
		arg.PasswordHash,
		arg.IsBlocked,
	)
	return err
}

const deleteUser = `-- name: DeleteUser :execrows
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id pgtype.UUID) (int64, error) {
	result, err := q.db.Exec(ctx, deleteUser, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, role, username, first_name, second_name, email, password_hash, is_blocked
FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Role,
		&i.Username,
		&i.FirstName,
		&i.SecondName,
		&i.Email,
		&i.PasswordHash,
		&i.IsBlocked,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, role, username, first_name, second_name, email, password_hash, is_blocked
FROM users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Role,
		&i.Username,
		&i.FirstName,
		&i.SecondName,
		&i.Email,
		&i.PasswordHash,
		&i.IsBlocked,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, role, username, first_name, second_name, email, password_hash, is_blocked
FROM users
ORDER BY username ASC, id ASC
LIMIT $2 OFFSET $1
`

type ListUsersParams struct {
	PageOffset int32 `json:"page_offset"`
	PageLimit  int32 `json:"page_limit"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers, arg.PageOffset, arg.PageLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Role,
			&i.Username,
			&i.FirstName,
			&i.SecondName,
			&i.Email,
			&i.PasswordHash,
			&i.IsBlocked,
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

const updateUser = `-- name: UpdateUser :execrows
UPDATE users
SET role = $1,
    username = $2,
    first_name = $3,
    second_name = $4,
    email = $5,
    password_hash = $6,
    is_blocked = $7
WHERE id = $8
`

type UpdateUserParams struct {
	Role         Role        `json:"role"`
	Username     string      `json:"username"`
	FirstName    string      `json:"first_name"`
	SecondName   string      `json:"second_name"`
	Email        string      `json:"email"`
	PasswordHash string      `json:"password_hash"`
	IsBlocked    bool        `json:"is_blocked"`
	ID           pgtype.UUID `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (int64, error) {
	result, err := q.db.Exec(ctx, updateUser,
		arg.Role,
		arg.Username,
		arg.FirstName,
		arg.SecondName,
		arg.Email,
		arg.PasswordHash,
		arg.IsBlocked,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}
