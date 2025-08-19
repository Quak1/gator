-- name: CreateFeed :one
INSERT INTO feeds (name, url, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeeds :many
SELECT *, users.name as username FROM feeds
INNER JOIN users
ON feeds.user_id = users.id;

-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE feeds.url = $1;
