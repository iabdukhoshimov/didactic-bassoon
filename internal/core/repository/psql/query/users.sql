-- name: CreateUser :one
INSERT INTO "users" (
    "first_name",
    "last_name",
    "email",
    "password",
    "role",
    "user_status"
  )
VALUES ($1, $2, $3, $4, $5, 'ACTIVE')
RETURNING "id";

-- name: CreateUserWithRole :one
INSERT INTO "users" (
    "first_name",
    "last_name",
    "email",
    "password",
    "role",
    "user_status"
) VALUES ($1, $2, $3, $4, 'OWNER', $5)
RETURNING "id";


-- name: GetUser :one
select *
from "users"
where id = $1 and "deleted_at" is null; 

-- name: GetUserByProperty :one
select *
from "users"
where "email" = $1;

-- name: GetUserByEmail :one
select id,
  "first_name",
  "last_name",
  "password",
  "role",
  "user_status"
from "users"
where "email" = $1;

-- name: GetUsers :many
select *
from "users"
where "deleted_at" is null;

-- name: UpdateUser :exec
update "users"
set "first_name" = $2,
  "last_name" = $3,
  "updated_at" = now()
where id = $1;

-- name: UpdateUserStatus :exec
update "users"
set "user_status" = $2
where id = $1;

-- name: UpdateUserStatusByEmail :exec
update "users"
set "user_status" = $2
where email = $1;

-- name: DeleteUser :exec
update "users"
set "deleted_at" = now()
where id = $1;

