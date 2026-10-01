-- name: CreateUser :exec
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
    sqlc.arg(id),
    sqlc.arg(role),
    sqlc.arg(username),
    sqlc.arg(first_name),
    sqlc.arg(second_name),
    sqlc.arg(email),
    sqlc.arg(password_hash),
    sqlc.arg(is_blocked)
);

-- name: GetUserByID :one
SELECT id, role, username, first_name, second_name, email, password_hash, is_blocked
FROM users
WHERE id = sqlc.arg(id);

-- name: GetUserByEmail :one
SELECT id, role, username, first_name, second_name, email, password_hash, is_blocked
FROM users
WHERE email = sqlc.arg(email);

-- name: ListUsers :many
SELECT id, role, username, first_name, second_name, email, password_hash, is_blocked
FROM users
ORDER BY username ASC, id ASC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: UpdateUser :execrows
UPDATE users
SET role = sqlc.arg(role),
    username = sqlc.arg(username),
    first_name = sqlc.arg(first_name),
    second_name = sqlc.arg(second_name),
    email = sqlc.arg(email),
    password_hash = sqlc.arg(password_hash),
    is_blocked = sqlc.arg(is_blocked)
WHERE id = sqlc.arg(id);

-- name: DeleteUser :execrows
DELETE FROM users
WHERE id = sqlc.arg(id);
