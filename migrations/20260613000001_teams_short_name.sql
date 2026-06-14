-- +goose Up
ALTER TABLE teams ADD COLUMN short_name TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE teams DROP COLUMN short_name;
