-- name: GetAllUsers :many 
select name
from users
order by name asc;