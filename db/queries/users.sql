-- name: ListUsers :many
SELECT * FROM users
ORDER BY email;

-- name: CreateUser :one
INSERT INTO users (
  email, password
) VALUES (
  $1, $2
)
RETURNING user_id, email, created_at;