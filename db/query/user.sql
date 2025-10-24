-- name: GetAllUsers :many
SELECT id, "name", username, email, "password",phone_number, "address", "role", created_at, updated_at
FROM users
WHERE deleted_at IS NULL;

-- name: GetUserByUsername :one
SELECT id, "name", username, email, "password",phone_number, "address", "role", created_at, updated_at
FROM users
WHERE username = $1 AND deleted_at IS NULL;

-- name: GetUserById :one
SELECT id, "name", username, email, "password",phone_number, "address", "role", created_at, updated_at
FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateUser :execresult
INSERT INTO users (id, name, username, email, password, phone_number, "address", role)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: UpdateUser :execresult
UPDATE users
SET
    "name" = $2,
    username = $3,
    email = $4,
    "password" = $5,
    "role" = $6,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteUser :execresult
UPDATE users
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;
