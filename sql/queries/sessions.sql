-- name: CreateSession :exec
INSERT INTO sessions (
    token,
    user_id,
    expires_at
) VALUES (
             ?, ?, ?
         );

-- name: GetSession :one
SELECT * FROM sessions
WHERE token = ? LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE token = ?;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at < ?;