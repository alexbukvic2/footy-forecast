-- +goose Up
-- +goose StatementBegin
ALTER TABLE tournaments ADD COLUMN predictions_locked BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tournaments DROP COLUMN IF EXISTS predictions_locked;
-- +goose StatementEnd
