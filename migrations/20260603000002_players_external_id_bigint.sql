-- +goose Up
-- +goose StatementBegin

ALTER TABLE players ALTER COLUMN external_id TYPE BIGINT USING external_id::bigint;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE players ALTER COLUMN external_id TYPE TEXT USING external_id::text;

-- +goose StatementEnd
