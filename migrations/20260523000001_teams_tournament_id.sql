-- +goose Up
-- +goose StatementBegin

ALTER TABLE teams
    ADD COLUMN tournament_id UUID NOT NULL REFERENCES tournaments(id);

CREATE INDEX teams_tournament_id_idx ON teams (tournament_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS teams_tournament_id_idx;
ALTER TABLE teams DROP COLUMN IF EXISTS tournament_id;

-- +goose StatementEnd
