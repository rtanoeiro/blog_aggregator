-- name: InsertFeed :one
INSERT INTO feed (name, url, user_id)
VALUES ($1,$2,$3)
RETURNING *;