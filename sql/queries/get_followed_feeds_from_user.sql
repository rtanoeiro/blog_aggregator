-- name: GetFollowedFeedsFromUser :many
SELECT
    feed.id as feed_id,
    feed.name,
    feed.url,
    feed.user_id
FROM feed_follows
INNER JOIN feed
    ON feed.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;