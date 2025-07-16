-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
    $5,
    $6
)
RETURNING id, created_at, updated_at, name, url, user_id;

-- name: Feeds :many
SELECT feeds.name, feeds.url, users.name AS uname FROM feeds
INNER JOIN users ON feeds.user_id = users.id;