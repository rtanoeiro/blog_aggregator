-- name: GetUsers :many 
select name
from users
order by name asc;