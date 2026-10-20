-- name: CreateUser :exec
INSERT INTO users (id, name, email, created_at, updated_at)
VALUES (sqlc.arg(id)::uuid, sqlc.arg(name), sqlc.arg(email),
        sqlc.arg(created_at), sqlc.arg(updated_at));

-- name: GetUserByID :one
SELECT id::text AS id, name, email, created_at, updated_at
FROM users
WHERE id = sqlc.arg(id)::uuid;

-- name: GetUserByEmail :one
SELECT id::text AS id, name, email, created_at, updated_at
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id::text AS id, name, email, created_at, updated_at
FROM users
ORDER BY created_at DESC, id DESC
LIMIT $1 OFFSET $2;

-- name: UpdateUser :execrows
UPDATE users
SET name = sqlc.arg(name),
    email = sqlc.arg(email),
    updated_at = sqlc.arg(updated_at)
WHERE id = sqlc.arg(id)::uuid;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE id = sqlc.arg(id)::uuid;