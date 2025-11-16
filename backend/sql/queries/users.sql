-- name: CreateUser :one
INSERT INTO users (
  name, email, phone_number, password_hash
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserAuthByEmail :one
SELECT id, name, password_hash 
FROM users 
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, name, email, upi_id, phone_number, email_verified, created_at, version
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: VerifyUserEmail :exec
UPDATE users
SET 
  email_verified = TRUE,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateUserEmail :exec
UPDATE users
SET
    email = $1,
    updated_at = CURRENT_TIMESTAMP
where id = $2;

-- name: UpdateUserProfile :one
UPDATE users
SET 
  name = $1, 
  upi_id = $2,
  updated_at = CURRENT_TIMESTAMP,
  version = version + 1
WHERE id = $3 AND version = $4
RETURNING *;
