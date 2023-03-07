-- name: GetWallets :many
SELECT
    wallets.id,
    wallets.name,
    wallets.sort_order,
    JSON_BUILD_OBJECT(
      'id', users.id,
      'name', users.name
        ) AS "user"
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