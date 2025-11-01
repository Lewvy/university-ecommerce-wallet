-- name: CreateTransaction :one
INSERT INTO wallet_transactions (
  user_id, amount, transaction_status, transaction_type, metadata
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateTransactionStatus :one
UPDATE wallet_transactions
SET 
  transaction_status = $1,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $2
RETURNING *;

-- name: GetTransaction :one
SELECT * FROM wallet_transactions
WHERE id = $1 AND user_id = $2;

-- name: ListInitialTransactionsForUser :many
SELECT id, amount, transaction_type, transaction_status, created_at 
FROM wallet_transactions 
WHERE user_id = $1 
ORDER BY 
    created_at DESC, 
    id DESC
LIMIT $2;

-- name: ListNextTransactionsForUser :many
SELECT id, amount, transaction_type, transaction_status, created_at 
FROM wallet_transactions 
WHERE 
    user_id = $1 
    AND 
    (created_at, id) < ($2, $3) 
ORDER BY 
    created_at DESC, 
    id DESC
LIMIT $4;
