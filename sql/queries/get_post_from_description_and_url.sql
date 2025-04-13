-- name: CheckIfPostExists :one
select
    title,
    url,
    description
from posts
where description = $1
    and url = $2;