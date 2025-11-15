-- name: CreateWallet :one
INSERT INTO wallets (
  user_id
) VALUES (
  $1
)
RETURNING *;

-- name: GetBalanceById :one
SELECT *
FROM wallets
WHERE user_id = $1;

-- name: GetWalletByUserID :one
SELECT * FROM wallets
WHERE user_id = $1
FOR UPDATE;

-- name: CreditWallet :one
UPDATE wallets 
SET 
  balance = balance + $1,
  lifetime_earned = lifetime_earned + $1
WHERE user_id = $2
RETURNING *;

-- name: DebitWallet :one
UPDATE wallets 
SET 
  balance = balance - $1,
  lifetime_spent = lifetime_spent + $1
WHERE user_id = $2
RETURNING *;
