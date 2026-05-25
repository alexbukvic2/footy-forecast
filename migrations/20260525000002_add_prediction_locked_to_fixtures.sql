-- +goose Up
-- +goose StatementBegin

ALTER TABLE fixtures ADD COLUMN prediction_locked BOOLEAN NOT NULL DEFAULT false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE fixtures DROP COLUMN prediction_locked;

-- +goose StatementEnd
