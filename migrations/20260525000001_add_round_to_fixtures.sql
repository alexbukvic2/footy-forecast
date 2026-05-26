-- +goose Up
-- +goose StatementBegin

ALTER TABLE fixtures ADD COLUMN round TEXT NOT NULL DEFAULT '';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE fixtures DROP COLUMN round;

-- +goose StatementEnd
