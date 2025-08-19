-- name: CreatePost :exec
INSERT INTO posts (title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (url) DO NOTHING;

-- name: GetPostsForUser :many
WITH user_feeds AS (
  SELECT feed_id
  FROM feed_follows
  WHERE user_id = $1
)
SELECT *
FROM posts
INNER JOIN user_feeds ON posts.feed_id = user_feeds.feed_id
ORDER BY published_at DESC
LIMIT $2;
