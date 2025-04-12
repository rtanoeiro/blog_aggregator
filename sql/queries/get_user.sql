-- name: GetUser :one
SELECT
    id,
    updated_at,
    created_at,
    name
FROM users
WHERE name = $1;
