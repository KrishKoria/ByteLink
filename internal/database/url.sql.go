// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: url.sql

package database

import (
	"context"
	"database/sql"
)

const deleteURLByID = `-- name: DeleteURLByID :exec
DELETE FROM urls WHERE id = ?
`

func (q *Queries) DeleteURLByID(ctx context.Context, id interface{}) error {
	_, err := q.db.ExecContext(ctx, deleteURLByID, id)
	return err
}

const getLongURLByShortURL = `-- name: GetLongURLByShortURL :one
SELECT u.long_url
FROM mappings m
JOIN urls u ON m.url_id = u.id
WHERE m.short_url = ?
LIMIT 1
`

func (q *Queries) GetLongURLByShortURL(ctx context.Context, shortUrl string) (string, error) {
	row := q.db.QueryRowContext(ctx, getLongURLByShortURL, shortUrl)
	var long_url string
	err := row.Scan(&long_url)
	return long_url, err
}

const getLongURLByShortURLAndUserID = `-- name: GetLongURLByShortURLAndUserID :one
SELECT u.long_url
FROM mappings m
JOIN urls u ON m.url_id = u.id
WHERE m.short_url = ? AND m.user_id = ?
LIMIT 1
`

type GetLongURLByShortURLAndUserIDParams struct {
	ShortUrl string
	UserID   interface{}
}

func (q *Queries) GetLongURLByShortURLAndUserID(ctx context.Context, arg GetLongURLByShortURLAndUserIDParams) (string, error) {
	row := q.db.QueryRowContext(ctx, getLongURLByShortURLAndUserID, arg.ShortUrl, arg.UserID)
	var long_url string
	err := row.Scan(&long_url)
	return long_url, err
}

const getOrphanedURLs = `-- name: GetOrphanedURLs :many
SELECT u.id FROM urls u
LEFT JOIN mappings m ON u.id = m.url_id
WHERE m.id IS NULL
`

func (q *Queries) GetOrphanedURLs(ctx context.Context) ([]interface{}, error) {
	rows, err := q.db.QueryContext(ctx, getOrphanedURLs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []interface{}
	for rows.Next() {
		var id interface{}
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getURLById = `-- name: GetURLById :one
SELECT id, long_url, created_at FROM urls WHERE id = ? LIMIT 1
`

func (q *Queries) GetURLById(ctx context.Context, id interface{}) (Url, error) {
	row := q.db.QueryRowContext(ctx, getURLById, id)
	var i Url
	err := row.Scan(&i.ID, &i.LongUrl, &i.CreatedAt)
	return i, err
}

const getURLIdByLongURL = `-- name: GetURLIdByLongURL :one
SELECT id FROM urls WHERE long_url = ? LIMIT 1
`

func (q *Queries) GetURLIdByLongURL(ctx context.Context, longUrl string) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, getURLIdByLongURL, longUrl)
	var id interface{}
	err := row.Scan(&id)
	return id, err
}

const getURLStatsForUser = `-- name: GetURLStatsForUser :one
SELECT m.short_url, u.long_url, m.click_count
FROM mappings m
JOIN urls u ON m.url_id = u.id
WHERE m.short_url = ? AND m.user_id = ?
`

type GetURLStatsForUserParams struct {
	ShortUrl string
	UserID   interface{}
}

type GetURLStatsForUserRow struct {
	ShortUrl   string
	LongUrl    string
	ClickCount sql.NullInt64
}

func (q *Queries) GetURLStatsForUser(ctx context.Context, arg GetURLStatsForUserParams) (GetURLStatsForUserRow, error) {
	row := q.db.QueryRowContext(ctx, getURLStatsForUser, arg.ShortUrl, arg.UserID)
	var i GetURLStatsForUserRow
	err := row.Scan(&i.ShortUrl, &i.LongUrl, &i.ClickCount)
	return i, err
}

const saveURL = `-- name: SaveURL :exec
INSERT INTO urls (id, long_url) VALUES (?, ?)
`

type SaveURLParams struct {
	ID      interface{}
	LongUrl string
}

func (q *Queries) SaveURL(ctx context.Context, arg SaveURLParams) error {
	_, err := q.db.ExecContext(ctx, saveURL, arg.ID, arg.LongUrl)
	return err
}
