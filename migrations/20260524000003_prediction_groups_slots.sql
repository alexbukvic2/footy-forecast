-- +goose Up
-- +goose StatementBegin

-- 1. Group letter on teams
ALTER TABLE teams
    ADD COLUMN group_letter VARCHAR(1);

-- 2. player_predictions: add group_letter, replace unique constraint
ALTER TABLE player_predictions
    ADD COLUMN group_letter VARCHAR(1);

ALTER TABLE player_predictions
    DROP CONSTRAINT player_predictions_user_tournament_category_uq;

CREATE UNIQUE INDEX player_predictions_group_uq
    ON player_predictions (user_id, tournament_id, category, group_letter)
    WHERE group_letter IS NOT NULL;

CREATE UNIQUE INDEX player_predictions_no_group_uq
    ON player_predictions (user_id, tournament_id, category)
    WHERE group_letter IS NULL;

-- 3. team_predictions: add group_letter + slot_index, replace unique constraint
ALTER TABLE team_predictions
    ADD COLUMN group_letter VARCHAR(1),
    ADD COLUMN slot_index   SMALLINT NOT NULL DEFAULT 0;

ALTER TABLE team_predictions
    DROP CONSTRAINT team_predictions_user_tournament_category_uq;

-- group-scoped (group_winner, playoff): unique per (category, group, slot)
CREATE UNIQUE INDEX team_predictions_group_uq
    ON team_predictions (user_id, tournament_id, category, group_letter, slot_index)
    WHERE group_letter IS NOT NULL;

-- tournament-scoped (semifinalist, winner): unique per (category, slot)
CREATE UNIQUE INDEX team_predictions_no_group_uq
    ON team_predictions (user_id, tournament_id, category, slot_index)
    WHERE group_letter IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS team_predictions_no_group_uq;
DROP INDEX IF EXISTS team_predictions_group_uq;
ALTER TABLE team_predictions
    DROP COLUMN IF EXISTS slot_index,
    DROP COLUMN IF EXISTS group_letter;
ALTER TABLE team_predictions
    ADD CONSTRAINT team_predictions_user_tournament_category_uq
        UNIQUE (user_id, tournament_id, category);

DROP INDEX IF EXISTS player_predictions_no_group_uq;
DROP INDEX IF EXISTS player_predictions_group_uq;
ALTER TABLE player_predictions
    DROP COLUMN IF EXISTS group_letter;
ALTER TABLE player_predictions
    ADD CONSTRAINT player_predictions_user_tournament_category_uq
        UNIQUE (user_id, tournament_id, category);

ALTER TABLE teams DROP COLUMN IF EXISTS group_letter;

-- +goose StatementEnd
