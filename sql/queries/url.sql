-- name: SaveURL :exec
INSERT INTO urls (id, long_url) VALUES (?, ?);

-- name: GetURLById :one
SELECT * FROM urls WHERE id = ? LIMIT 1;

-- name: GetURLIdByLongURL :one
SELECT id FROM urls WHERE long_url = ? LIMIT 1;

-- name: GetLongURLByShortURLAndUserID :one
SELECT u.long_url
FROM mappings m
JOIN urls u ON m.url_id = u.id
WHERE m.short_url = ? AND m.user_id = ?
LIMIT 1;

-- name: GetLongURLByShortURL :one
SELECT u.long_url
FROM mappings m
JOIN urls u ON m.url_id = u.id
WHERE m.short_url = ?
LIMIT 1;
