-- name: CreateFeedFollow :one
WITH user_lookup AS (
  SELECT id, name 
  FROM users 
  WHERE users.name = $1
),
inserted_feed_follow AS (
  INSERT INTO feed_follows (user_id, feed_id, created_at, updated_at)
  SELECT user_lookup.id, $2, $3, $4
  FROM user_lookup
  RETURNING *
)
SELECT
  inserted_feed_follow.*,
  feeds.name AS feed_name,
  user_lookup.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id
INNER JOIN user_lookup ON true;


-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = $1
AND feed_follows.feed_id = $2;


-- name: GetFeedFolllowsForUser :many
WITH user_lookup AS (
  SELECT id
  FROM users 
  WHERE users.name = $1
)
SELECT 
  feeds.name as feed_name,
  feed_follows.feed_id as feed_id,
  feed_follows.user_id as user_id
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN user_lookup ON feed_follows.user_id = user_lookup.id
WHERE user_lookup.id = feed_follows.user_id;
