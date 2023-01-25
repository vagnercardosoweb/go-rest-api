-- name: GetAllTest :many
SELECT * FROM test LIMIT $1 OFFSET $2;
