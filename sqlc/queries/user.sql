-- name: GetUserByEmailInLogin :one
SELECT
    id,
    name,
    email,
    password_hash
FROM users
WHERE LOWER(email) = LOWER(sqlc.arg('email'))
LIMIT 1;

-- name: CreateUser :exec
INSERT INTO
    users (name, email, birth_date, code_to_invite, password_hash, token_to_confirm_email, confirmed_email_at,
           login_blocked_until, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;

-- name: GetUsers :many
SELECT
    id,
    name,
    email,
    birth_date,
    code_to_invite,
    token_to_confirm_email,
    confirmed_email_at,
    login_blocked_until
FROM users
LIMIT sqlc.arg('limit');
