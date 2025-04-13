-- name: GetFeeds :many
select
    feed.name,
    feed.url,
    feed.user_id,
    users.name as username
from feed
inner join users
    on users.id = feed.user_id;