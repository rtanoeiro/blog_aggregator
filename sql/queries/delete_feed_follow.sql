-- name: Unfollow :one
WITH url_to_unfollow as (
    SELECT
        feed_follows.id,
        feed_follows.user_id
    from feed
    inner join feed_follows
        on feed_follows.feed_id = feed.id
    where feed.url = $1
        and feed_follows.user_id = $2
)

DELETE FROM feed_follows
WHERE id = (select id from url_to_unfollow)
and user_id = (select user_id from url_to_unfollow)

RETURNING *;