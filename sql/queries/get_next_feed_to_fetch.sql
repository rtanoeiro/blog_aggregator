-- name: GetNextFeedToFetch :one
select 
    url
from feed
order by last_fetched_at nulls FIRST
limit 1;