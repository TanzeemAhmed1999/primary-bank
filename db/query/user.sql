-- name: GetUser :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (username, password, full_name, email)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET full_name = $2, email = $3, password = COALESCE($4, password), updated_at = now()
WHERE username = $1
RETURNING username, password, full_name, email, created_at, updated_at;
