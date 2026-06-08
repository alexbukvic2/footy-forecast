-- +goose Up
-- +goose StatementBegin

INSERT INTO player_handicap_defaults (tournament_id, category, default_points)
SELECT id, 'group_top_scorer'::player_handicap_category,  10 FROM tournaments WHERE slug = 'world-cup-2026'
UNION ALL
SELECT id, 'total_top_scorer'::player_handicap_category,  20 FROM tournaments WHERE slug = 'world-cup-2026'
ON CONFLICT (tournament_id, category) DO UPDATE SET default_points = EXCLUDED.default_points;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM player_handicap_defaults
WHERE tournament_id = (SELECT id FROM tournaments WHERE slug = 'world-cup-2026')
  AND category IN ('group_top_scorer', 'total_top_scorer');

-- +goose StatementEnd
