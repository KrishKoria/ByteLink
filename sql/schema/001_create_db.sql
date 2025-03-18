-- +goose Up
CREATE TABLE mapping (
 id UUID PRIMARY KEY,
 short_url TEXT NOT NULL,
 long_url TEXT NOT NULL,
 userId UUID NOT NULL
);

-- +goose Down
DROP TABLE mapping;