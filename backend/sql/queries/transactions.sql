-- name: CreateTransaction :one
INSERT INTO wallet_transactions (
    user_id,
    amount,
    transaction_type,
    transaction_status,
    related_user_id,
    razorpay_order_id,
    razorpay_payment_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;


-- name: CreditWalletBalance :exec
UPDATE wallets
SET balance = balance + $2
WHERE user_id = $1;


-- name: GetTransactionAmount :one
SELECT amount FROM wallet_transactions
WHERE id = $1;

-- name: GetTransactionByOrderID :one
SELECT * FROM wallet_transactions
WHERE razorpay_order_id = $1
LIMIT 1;

-- name: UpdateTransactionOrderID :one
UPDATE wallet_transactions
SET razorpay_order_id = $2
WHERE id = $1
RETURNING *;

-- name: UpdateTransactionStatus :one
UPDATE wallet_transactions
SET transaction_status = $2, razorpay_payment_id = $3
WHERE id = $1
RETURNING *;

