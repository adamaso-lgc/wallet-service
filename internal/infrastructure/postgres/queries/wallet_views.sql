-- name: UpsertWalletView :exec
INSERT INTO wallet_views (id, owner_id, balance, currency, status, created_at, updated_at)
VALUES (@id, @owner_id, @balance, @currency, @status, @created_at, @updated_at)
ON CONFLICT (id) DO UPDATE SET
    balance    = EXCLUDED.balance,
    currency   = EXCLUDED.currency,
    status     = EXCLUDED.status,
    updated_at = EXCLUDED.updated_at;

-- name: UpdateWalletViewBalance :exec
UPDATE wallet_views
SET balance = @balance, updated_at = @updated_at
WHERE id = @id;

-- name: UpdateWalletViewStatus :exec
UPDATE wallet_views
SET status = @status, updated_at = @updated_at
WHERE id = @id;

-- name: GetWalletView :one
SELECT id, owner_id, balance, currency, status, created_at, updated_at
FROM   wallet_views
WHERE  id = @id;

-- name: ListWalletViewsByOwner :many
SELECT id, owner_id, balance, currency, status, created_at, updated_at
FROM   wallet_views
WHERE  owner_id = @owner_id
ORDER  BY created_at ASC;
