-- name: CreateUser :exec
INSERT INTO users (id, name, email, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);

-- name: GetUserByID :one
SELECT id, name, email, created_at, updated_at
FROM users
WHERE id = ?;

-- name: GetUserByEmail :one
SELECT id, name, email, created_at, updated_at
FROM users
WHERE email = ?;

-- name: ListUsers :many
SELECT id, name, email, created_at, updated_at
FROM users
ORDER BY created_at DESC, id DESC
LIMIT ? OFFSET ?;

-- name: UpdateUser :execrows
UPDATE users
SET name = ?, email = ?, updated_at = ?
WHERE id = ?;

-- name: DeleteUser :execrows
DELETE FROM users WHERE id = ?;