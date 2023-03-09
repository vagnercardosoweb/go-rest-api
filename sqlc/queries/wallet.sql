-- name: GetWallets :many
SELECT
    wallets.id,
    wallets.name,
    wallets.sort_order,
    users.id AS "user_id",
    users.name AS "user_name"
FROM
    wallets
        INNER JOIN users ON wallets.user_id = users.id
LIMIT sqlc.arg('limit');

-- name: GetWalletById :one
SELECT
    id,
    name,
    sort_order,
    user_id
FROM wallets
WHERE id = sqlc.arg('id')
LIMIT 1;