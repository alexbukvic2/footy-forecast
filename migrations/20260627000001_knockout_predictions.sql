-- +goose Up
-- +goose StatementBegin

-- Split goals_home/goals_away into regular-time and total (regular + ET, excluding penalties).
-- Scoring uses regular-time goals; total is for display.
ALTER TABLE fixtures
    RENAME COLUMN goals_home TO goals_home_regular;
ALTER TABLE fixtures
    RENAME COLUMN goals_away TO goals_away_regular;
ALTER TABLE fixtures
    ADD COLUMN goals_home INT CHECK (goals_home >= 0),
    ADD COLUMN goals_away INT CHECK (goals_away >= 0);

-- Back-fill: for existing rows, total = regular (all past matches were regular-time only).
UPDATE fixtures SET goals_home = goals_home_regular, goals_away = goals_away_regular;

-- Predicted winner for knockout fixtures (home or away team that the user thinks will advance).
-- NULL for group-stage predictions; required for knockout predictions.
ALTER TABLE score_predictions
    ADD COLUMN winner UUID REFERENCES teams(id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE score_predictions DROP COLUMN IF EXISTS winner;

ALTER TABLE fixtures DROP COLUMN IF EXISTS goals_away;
ALTER TABLE fixtures DROP COLUMN IF EXISTS goals_home;
ALTER TABLE fixtures RENAME COLUMN goals_away_regular TO goals_away;
ALTER TABLE fixtures RENAME COLUMN goals_home_regular TO goals_home;

-- +goose StatementEnd
