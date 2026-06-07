-- +goose Up
-- +goose StatementBegin

-- group_winner, playoff, semifinalist, and winner handicap points.
-- Idempotent: ON CONFLICT updates existing rows.
-- NOTE: team names must match exactly what was seeded from the API.
--       Verify with: SELECT name FROM teams ORDER BY name;

INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, 'group_winner'::team_handicap_category, v.points
FROM teams t
JOIN (VALUES
    -- Group A
    ('Mexico',                  5),
    ('Czech Republic',                12),
    ('South Korea',            13),
    ('South Africa',           38),
    -- Group B
    ('Switzerland',             5),
    ('Canada',                  8),
    ('Bosnia & Herzegovina',   16),
    ('Qatar',                  50),
    -- Group C
    ('Brazil',                  5),
    ('Morocco',                18),
    ('Scotland',               41),
    ('Haiti',                  50),
    -- Group D
    ('United States',           5),
    ('Türkiye',                 6),
    ('Paraguay',               11),
    ('Australia',              20),
    -- Group E
    ('Germany',                 5),
    ('Ecuador',                17),
    ('Ivory Coast',            26),
    ('Curaçao',                50),
    -- Group F
    ('Netherlands',             5),
    ('Japan',                  11),
    ('Sweden',                 16),
    ('Tunisia',                37),
    -- Group G
    ('Belgium',                 5),
    ('Egypt',                  18),
    ('Iran',                   32),
    ('New Zealand',            50),
    -- Group H
    ('Spain',                   5),
    ('Uruguay',                21),
    ('Saudi Arabia',           50),
    ('Cape Verde Islands',     50),
    -- Group I
    ('France',                  5),
    ('Norway',                 13),
    ('Senegal',                29),
    ('Iraq',                   50),
    -- Group J
    ('Argentina',               5),
    ('Austria',                17),
    ('Algeria',                30),
    ('Jordan',                 50),
    -- Group K
    ('Portugal',                5),
    ('Colombia',               10),
    ('Congo DR',               50),
    ('Uzbekistan',             50),
    -- Group L
    ('England',                 5),
    ('Croatia',                 8),
    ('Panama',                 50),
    ('Ghana',                  50)
) AS v(name, points) ON t.name = v.name
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, 'playoff'::team_handicap_category, v.points
FROM teams t
JOIN (VALUES
    -- Group A
    ('Mexico',                  5),
    ('Czechia',                12),
    ('South Korea',            12),
    ('South Africa',           20),
    -- Group B
    ('Switzerland',             5),
    ('Canada',                 10),
    ('Bosnia & Herzegovina',   20),
    ('Qatar',                  33),
    -- Group C
    ('Brazil',                  5),
    ('Morocco',                10),
    ('Scotland',               13),
    ('Haiti',                  45),
    -- Group D
    ('USA',                     5),
    ('Türkiye',                 6),
    ('Paraguay',                7),
    ('Australia',               9),
    -- Group E
    ('Germany',                 5),
    ('Ecuador',                11),
    ('Ivory Coast',            12),
    ('Curaçao',                50),
    -- Group F
    ('Netherlands',             5),
    ('Japan',                  12),
    ('Sweden',                 16),
    ('Tunisia',                24),
    -- Group G
    ('Belgium',                 5),
    ('Egypt',                  13),
    ('Iran',                   15),
    ('New Zealand',            25),
    -- Group H
    ('Spain',                   5),
    ('Uruguay',                10),
    ('Saudi Arabia',           19),
    ('Cape Verde',             29),
    -- Group I
    ('France',                  5),
    ('Norway',                 11),
    ('Senegal',                16),
    ('Iraq',                   27),
    -- Group J
    ('Argentina',               5),
    ('Austria',                12),
    ('Algeria',                13),
    ('Jordan',                 22),
    -- Group K
    ('Portugal',                5),
    ('Colombia',               11),
    ('DR Congo',               18),
    ('Uzbekistan',             30),
    -- Group L
    ('England',                 5),
    ('Croatia',                12),
    ('Ghana',                  16),
    ('Panama',                 32)
) AS v(name, points) ON t.name = v.name
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, 'semifinalist'::team_handicap_category, v.points
FROM teams t
JOIN (VALUES
    ('Spain',                   5),
    ('England',                 5),
    ('France',                  6),
    ('Argentina',               6),
    ('Brazil',                  7),
    ('Portugal',                7),
    ('Germany',                 9),
    ('Belgium',                11),
    ('Netherlands',            11),
    ('Norway',                 14),
    ('Colombia',               16),
    ('Uruguay',                16),
    ('United States',          18),
    ('Switzerland',            20),
    ('Mexico',                 20),
    ('Ecuador',                22),
    ('Croatia',                22),
    ('Japan',                  23),
    ('Morocco',                23),
    ('Austria',                27),
    ('Türkiye',                27),
    ('Sweden',                 30),
    ('Senegal',                30),
    ('Canada',                 30),
    ('Algeria',                34),
    ('Ivory Coast',            34),
    ('Paraguay',               34),
    ('Egypt',                  39),
    ('Ghana',                  43),
    ('Scotland',               48),
    ('South Korea',            50),
    ('Iran',                   50),
    ('Czechia',                50),
    ('Bosnia & Herzegovina',   50),
    ('DR Congo',               50),
    ('Saudi Arabia',           50),
    ('Tunisia',                50),
    ('South Africa',           50),
    ('Australia',              50),
    ('Panama',                 50),
    ('Uzbekistan',             50),
    ('Iraq',                   50),
    ('New Zealand',            50),
    ('Cape Verde',             50),
    ('Qatar',                  50),
    ('Jordan',                 50),
    ('Curaçao',                50),
    ('Haiti',                  50)
) AS v(name, points) ON t.name = v.name
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, 'winner'::team_handicap_category, v.points
FROM teams t
JOIN (VALUES
    ('Spain',                   5),
    ('France',                  5),
    ('England',                 7),
    ('Brazil',                  8),
    ('Argentina',               9),
    ('Portugal',               10),
    ('Germany',                13),
    ('Netherlands',            21),
    ('Norway',                 33),
    ('Belgium',                33),
    ('Colombia',               37),
    ('Morocco',                46),
    ('Uruguay',                46),
    ('United States',          50),
    ('Switzerland',            50),
    ('Japan',                  50),
    ('Mexico',                 50),
    ('Croatia',                50),
    ('Ecuador',                50),
    ('Senegal',                50),
    ('Sweden',                 50),
    ('Austria',                50),
    ('Türkiye',                50),
    ('Canada',                 50),
    ('Algeria',                50),
    ('Ivory Coast',            50),
    ('Paraguay',               50),
    ('Egypt',                  50),
    ('Scotland',               50),
    ('Bosnia & Herzegovina',   50),
    ('Ghana',                  50),
    ('Czechia',                50),
    ('South Korea',            50),
    ('Iran',                   50),
    ('Tunisia',                50),
    ('Cape Verde',             50),
    ('Uzbekistan',             50),
    ('Haiti',                  50),
    ('Panama',                 50),
    ('Curaçao',                50),
    ('Qatar',                  50),
    ('Saudi Arabia',           50),
    ('New Zealand',            50),
    ('Australia',              50),
    ('DR Congo',               50),
    ('Iraq',                   50),
    ('Jordan',                 50),
    ('South Africa',           50)
) AS v(name, points) ON t.name = v.name
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM team_handicap
WHERE category IN ('group_winner', 'playoff', 'semifinalist', 'winner');

-- +goose StatementEnd
