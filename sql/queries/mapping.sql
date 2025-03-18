-- name: SaveMapping :exec
INSERT INTO mapping (id, short_url, long_url, userId)
VALUES (?, ?, ?, ?);

-- name: GetLongURLByShortURLAndUserID :one
SELECT long_url
FROM mapping
WHERE short_url = ? AND userId = ?;

-- name: GetLongURLByShortURL :one
SELECT long_url
FROM mapping
WHERE short_url = ?;

-- name: GetMappingsByUserID :many
SELECT id, short_url, long_url
FROM mapping
WHERE userId = ?;