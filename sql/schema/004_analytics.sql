
-- +goose Up
ALTER TABLE mappings ADD COLUMN click_count INTEGER DEFAULT 0;