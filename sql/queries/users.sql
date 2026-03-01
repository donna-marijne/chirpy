-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, hashed_password)
values (
	gen_random_uuid(),
	now(),
	now(),
	$1,
	$2
)
returning *;

-- name: GetUserByEmail :one
select *
from users
where email = $1;

-- name: UpdateUser :one
update users
set
	updated_at = now(),
	email = $2,
	hashed_password = $3
where id = $1
returning *;

-- name: UpdateUserSetChirpyRed :one
update users
set
	updated_at = now(),
	is_chirpy_red = $2
where id = $1
returning *;

-- name: DeleteUsers :exec
delete from users;
