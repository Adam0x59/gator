-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES($1,$2,$3,$4,$5)
    RETURNING *
)
    SELECT
        inserted_feed_follow.*,
        feeds.name AS feed_name,
        users.name AS user_name
    FROM inserted_feed_follow
    INNER JOIN users ON users.id = inserted_feed_follow.user_id
    INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT 
    feed_follows.*, 
    feeds.name AS feed_name, 
    users.name AS user_name 
FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: DeleteFollow :exec
DELETE
FROM feed_follows
WHERE user_id = $1 AND feed_id = $2; 

-- name: MarkFeedFetched :exec
UPDATE feeds 
SET updated_at = $1, 
    last_fetched_at = $2
WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
