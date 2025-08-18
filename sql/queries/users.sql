-- name: CreateUser :one
INSERT INTO users (name, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE name = $1 LIMIT 1;

-- name: DeleteUsers :exec
DELETE FROM users;
