-- name: CreateRefreshToken :one
insert into refresh_tokens (
	token,
	created_at,
	updated_at,
	user_id,
	expires_at,
	revoked_at
)
values (
	$1,
	now(),
	now(),
	$2,
	$3,
	null
)
returning *;

-- name: GetRefreshToken :one
select token, created_at, updated_at, user_id, expires_at, revoked_at
from refresh_tokens
where token = $1;

-- name: RevokeRefreshToken :one
update refresh_tokens
set updated_at = now(), revoked_at = now()
where token = $1
returning *;

