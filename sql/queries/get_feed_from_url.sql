-- name: GetFeedFromURL :one
select
    feed.id as feed_id,
    feed.name,
    feed.url,
    feed.user_id
from feed
where feed.url = $1;