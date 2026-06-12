-- +goose Up
-- +goose StatementBegin

-- 1. Shift existing playoff predictions up to make room at slot_index 0.
--    Two-step via negatives avoids unique-constraint violations from row-by-row
--    checking: 0,1,2 → -1,-2,-3 → 1,2,3 (no intermediate collisions).
UPDATE team_predictions
SET slot_index = -(slot_index + 1)
WHERE category = 'playoff';

UPDATE team_predictions
SET slot_index = -slot_index
WHERE category = 'playoff';

-- 2. Clone group_winner records into playoff at slot_index 0
INSERT INTO team_predictions (user_id, tournament_id, category, pick, group_letter, slot_index, points)
SELECT user_id, tournament_id, 'playoff', pick, group_letter, 0, NULL
FROM team_predictions
WHERE category = 'group_winner';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- 1. Remove the cloned playoff rows (slot_index 0 were seeded from group_winner)
DELETE FROM team_predictions
WHERE category = 'playoff' AND slot_index = 0;

-- 2. Shift remaining playoff predictions back down (same two-step trick).
UPDATE team_predictions
SET slot_index = -(slot_index - 1)
WHERE category = 'playoff';

UPDATE team_predictions
SET slot_index = -slot_index
WHERE category = 'playoff';

-- +goose StatementEnd
