-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    name,
    password,
    provider,
    email_verified,
    created_at,
    updated_at
) VALUES (
             ?, ?, ?, ?, ?, ?, ?, ?
         )
RETURNING *;


-- name: CreateOAuthUser :one
INSERT INTO users (
    id,
    email,
    name,
    provider,
    provider_id,
    email_verified,
    created_at,
    updated_at
) VALUES (
             ?, ?, ?, ?, ?, ?, ?, ?
         )
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ? LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    name = COALESCE(?, name),
    email = COALESCE(?, email),
    password = COALESCE(?, password),
    email_verified = COALESCE(?, email_verified),
    updated_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: GetUserByProviderID :one
SELECT * FROM users
WHERE provider = ? AND provider_id = ? LIMIT 1;