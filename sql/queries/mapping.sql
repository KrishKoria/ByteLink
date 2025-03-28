-- name: SaveMapping :exec
INSERT INTO mappings (id, url_id, short_url, user_id) VALUES (?, ?, ?, ?);

-- name: GetMappingByShortURLAndUserID :one
SELECT m.id, m.short_url, u.long_url, m.user_id
FROM mappings m
JOIN urls u ON m.url_id = u.id
WHERE m.short_url = ? AND m.user_id = ?
LIMIT 1;

-- name: GetMappingsByUserID :many
SELECT m.short_url, u.long_url
FROM mappings m
JOIN urls u ON m.url_id = u.id
WHERE m.user_id = ?;

-- name: DeleteMapping :exec
DELETE FROM mappings
WHERE short_url = ? AND user_id = ?;

-- name: GetMappingByUserIDAndUrlID :one
SELECT id, short_url, url_id, user_id
FROM mappings
WHERE user_id = ? AND url_id = ?
LIMIT 1;

-- name: DeleteMappingByShortURLAndUserID :exec
DELETE FROM mappings
WHERE short_url = ? AND user_id = ?;

-- name: CountMappingsByUrlID :one
SELECT COUNT(*) FROM mappings WHERE url_id = ?;

-- name: IncrementClickCount :exec
UPDATE mappings
SET click_count = click_count + 1
WHERE short_url = ?;