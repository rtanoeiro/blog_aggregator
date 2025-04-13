-- name: GetFeedNameFromUser :one
select
    feed.name
from feed_follows
inner join feed
    on feed.id = feed_follows.feed_id
where feed_follows.user_id = $1;