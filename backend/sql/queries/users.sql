-- name: CreateUser :one
INSERT INTO users (
  name, email, password_hash
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetUserAuthByEmail :one
-- Gets the user's ID and password hash for login
SELECT id, name, password_hash 
FROM users 
WHERE email = $1;


-- name: GetUserByID :one
-- Gets a user's public profile data
SELECT id, name, email, upi_id, email_verified, created_at, version
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
-- Checks if an email exists and gets user info
SELECT * FROM users
WHERE email = $1;

-- name: VerifyUserEmail :exec
UPDATE users
SET 
  email_verified = TRUE,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

--name: UpdateUserEmail :exec
Update users
Set
	email = $1,
	updated_at = CURRENT_TIMESTAMP
where id = $2

-- name: UpdateUserProfile :one
UPDATE users
SET 
  name = $1, 
  upi_id = $2,
  updated_at = CURRENT_TIMESTAMP,
  version = version + 1
WHERE id = $3 AND version = $4
RETURNING *;
