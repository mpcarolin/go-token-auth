-- name: ListUsers :many
SELECT * FROM users
ORDER BY email;

-- name: UserExistsByEmail :one
SELECT EXISTS (
    SELECT 1 FROM users
    WHERE email = $1
);

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (
  email, password
) VALUES (
  $1, $2
)
RETURNING user_id, email, created_at;