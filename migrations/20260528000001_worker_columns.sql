-- +goose Up

-- +goose StatementBegin
ALTER TYPE fixture_status ADD VALUE IF NOT EXISTS 'cancelled';
-- +goose StatementEnd

-- +goose StatementBegin

ALTER TABLE fixtures
    ADD COLUMN is_demo        BOOLEAN     NOT NULL DEFAULT FALSE,
    ADD COLUMN winner_team_id UUID        REFERENCES teams(id),
    ADD COLUMN last_polled_at TIMESTAMPTZ;

ALTER TABLE tournaments
    ADD COLUMN external_id BIGINT,
    ADD COLUMN season      SMALLINT;

ALTER TABLE teams
    ADD COLUMN external_id BIGINT;

ALTER TABLE score_predictions
    ADD COLUMN scored_at TIMESTAMPTZ;

ALTER TABLE player_predictions
    ADD COLUMN scored_at TIMESTAMPTZ;

ALTER TABLE team_predictions
    ADD COLUMN scored_at TIMESTAMPTZ;

CREATE INDEX score_predictions_fixture_id_idx ON score_predictions (fixture_id);

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP INDEX IF EXISTS score_predictions_fixture_id_idx;

ALTER TABLE team_predictions   DROP COLUMN IF EXISTS scored_at;
ALTER TABLE player_predictions DROP COLUMN IF EXISTS scored_at;
ALTER TABLE score_predictions  DROP COLUMN IF EXISTS scored_at;

ALTER TABLE teams        DROP COLUMN IF EXISTS external_id;
ALTER TABLE tournaments  DROP COLUMN IF EXISTS season,
                         DROP COLUMN IF EXISTS external_id;
ALTER TABLE fixtures     DROP COLUMN IF EXISTS last_polled_at,
                         DROP COLUMN IF EXISTS winner_team_id,
                         DROP COLUMN IF EXISTS is_demo;

-- Note: enum values cannot be removed in PostgreSQL; 'cancelled' remains in fixture_status.

-- +goose StatementEnd
