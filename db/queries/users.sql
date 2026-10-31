-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByGoogleID :one
SELECT * FROM users WHERE google_id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (email, name, picture, google_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateEmailUser :one
INSERT INTO users (email, name, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET name = $2, picture = $3, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateUserHourlyWorth :one
UPDATE users
SET hourly_worth = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
