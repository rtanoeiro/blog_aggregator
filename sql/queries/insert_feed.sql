-- name: InsertFeed :one
INSERT INTO feed (id, name, url, user_id, last_fetched_at)
VALUES ($1,$2,$3,$4,$5)
RETURNING *;