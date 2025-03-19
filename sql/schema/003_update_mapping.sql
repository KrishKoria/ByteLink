-- +goose Up
DROP TABLE mapping;

CREATE TABLE mappings (
id UUID PRIMARY KEY,
url_id UUID NOT NULL,
short_url TEXT NOT NULL,
user_id UUID NULL,  -- make this nullable for guest users, if needed
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
FOREIGN KEY (url_id) REFERENCES urls(id),
UNIQUE(user_id, url_id)  -- ensures a user doesn't get duplicate mappings for the same URL
);