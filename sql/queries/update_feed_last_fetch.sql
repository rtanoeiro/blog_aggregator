-- name: UpdateFetchDate :one
update feed
set last_fetched_at = $1
where id = $2

RETURNING *;