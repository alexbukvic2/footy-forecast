-- +goose Up
-- +goose StatementBegin

ALTER TABLE tournament_group_table
    ADD COLUMN description TEXT NOT NULL DEFAULT '';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE tournament_group_table DROP COLUMN description;

-- +goose StatementEnd
