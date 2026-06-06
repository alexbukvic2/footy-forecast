-- +goose Up
-- +goose StatementBegin

-- The original teams table had UNIQUE(name) before tournament_id was added.
-- Now that each tournament manages its own set of teams, uniqueness must be
-- scoped to (name, tournament_id) so different tournaments can share team names.
ALTER TABLE teams DROP CONSTRAINT teams_name_key;
ALTER TABLE teams ADD CONSTRAINT teams_name_tournament_uq UNIQUE (name, tournament_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE teams DROP CONSTRAINT teams_name_tournament_uq;
ALTER TABLE teams ADD CONSTRAINT teams_name_key UNIQUE (name);

-- +goose StatementEnd
