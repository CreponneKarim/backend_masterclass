-- name: CreateAccount :one
INSERT INTO account (
  owner,
  balance,
  curency
) VALUES (
  $1, $2, $3
)
RETURNING *;