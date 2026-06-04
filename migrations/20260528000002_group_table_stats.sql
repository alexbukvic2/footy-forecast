-- +goose Up
-- +goose StatementBegin

ALTER TABLE tournament_group_table
    ADD COLUMN won           SMALLINT NOT NULL DEFAULT 0,
    ADD COLUMN drawn         SMALLINT NOT NULL DEFAULT 0,
    ADD COLUMN lost          SMALLINT NOT NULL DEFAULT 0,
    ADD COLUMN goals_for     SMALLINT NOT NULL DEFAULT 0,
    ADD COLUMN goals_against SMALLINT NOT NULL DEFAULT 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE tournament_group_table
    DROP COLUMN IF EXISTS goals_against,
    DROP COLUMN IF EXISTS goals_for,
    DROP COLUMN IF EXISTS lost,
    DROP COLUMN IF EXISTS drawn,
    DROP COLUMN IF EXISTS won;

-- +goose StatementEnd
