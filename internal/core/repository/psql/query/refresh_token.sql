-- name: CreateRefreshToken :one
insert into "refresh_token" (
        "user_id",
        "refresh_token",
        "expires_at"
    )
values ($1, $2, $3)
returning id;

-- name: GetRefreshToken :one
select *
from "refresh_token"
where id = $1;

-- name: UpdateRefreshToken :exec
update "refresh_token"
set "refresh_token" = $2,
    "expires_at" = $3
where id = $1;

-- name: DeleteRefreshToken :exec
delete from "refresh_token"
where refresh_token = $1;
