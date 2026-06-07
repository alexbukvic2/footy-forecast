-- +goose Up
-- +goose StatementBegin

-- Seeded from API-Football league=1 season=2026
-- Tournament and teams are inserted idempotently:
--   ON CONFLICT (slug)  DO NOTHING  for tournaments
--   ON CONFLICT (name)  DO UPDATE   for teams (updates external_id, logo, tournament_id, group_letter)

INSERT INTO tournaments (id, slug, name, external_id, season, starts_at, ends_at)
VALUES (
    gen_random_uuid(),
    'world-cup-2026',
    'World Cup 2026',
    1,
    2026,
    '2026-06-11T00:00:00Z',
    '2026-07-19T00:00:00Z'
) ON CONFLICT (slug) DO NOTHING;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1,
    'Belgium',
    'https://media.api-sports.io/football/teams/1.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'G'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    2,
    'France',
    'https://media.api-sports.io/football/teams/2.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'I'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    3,
    'Croatia',
    'https://media.api-sports.io/football/teams/3.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'L'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    5,
    'Sweden',
    'https://media.api-sports.io/football/teams/5.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'F'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    6,
    'Brazil',
    'https://media.api-sports.io/football/teams/6.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'C'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    7,
    'Uruguay',
    'https://media.api-sports.io/football/teams/7.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'H'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    8,
    'Colombia',
    'https://media.api-sports.io/football/teams/8.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'K'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    9,
    'Spain',
    'https://media.api-sports.io/football/teams/9.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'H'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    10,
    'England',
    'https://media.api-sports.io/football/teams/10.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'L'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    11,
    'Panama',
    'https://media.api-sports.io/football/teams/11.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'L'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    12,
    'Japan',
    'https://media.api-sports.io/football/teams/12.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'F'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    13,
    'Senegal',
    'https://media.api-sports.io/football/teams/13.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'I'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    15,
    'Switzerland',
    'https://media.api-sports.io/football/teams/15.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'B'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    16,
    'Mexico',
    'https://media.api-sports.io/football/teams/16.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'A'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    17,
    'South Korea',
    'https://media.api-sports.io/football/teams/17.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'A'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    20,
    'Australia',
    'https://media.api-sports.io/football/teams/20.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'D'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    22,
    'Iran',
    'https://media.api-sports.io/football/teams/22.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'G'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    23,
    'Saudi Arabia',
    'https://media.api-sports.io/football/teams/23.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'H'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    25,
    'Germany',
    'https://media.api-sports.io/football/teams/25.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'E'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    26,
    'Argentina',
    'https://media.api-sports.io/football/teams/26.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'J'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    27,
    'Portugal',
    'https://media.api-sports.io/football/teams/27.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'K'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    28,
    'Tunisia',
    'https://media.api-sports.io/football/teams/28.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'F'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    31,
    'Morocco',
    'https://media.api-sports.io/football/teams/31.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'C'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    32,
    'Egypt',
    'https://media.api-sports.io/football/teams/32.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'G'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    770,
    'Czech Republic',
    'https://media.api-sports.io/football/teams/770.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'A'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    775,
    'Austria',
    'https://media.api-sports.io/football/teams/775.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'J'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    777,
    'Türkiye',
    'https://media.api-sports.io/football/teams/777.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'D'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1090,
    'Norway',
    'https://media.api-sports.io/football/teams/1090.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'I'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1108,
    'Scotland',
    'https://media.api-sports.io/football/teams/1108.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'C'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1113,
    'Bosnia & Herzegovina',
    'https://media.api-sports.io/football/teams/1113.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'B'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1118,
    'Netherlands',
    'https://media.api-sports.io/football/teams/1118.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'F'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1501,
    'Ivory Coast',
    'https://media.api-sports.io/football/teams/1501.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'E'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1504,
    'Ghana',
    'https://media.api-sports.io/football/teams/1504.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'L'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1508,
    'Congo DR',
    'https://media.api-sports.io/football/teams/1508.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'K'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1531,
    'South Africa',
    'https://media.api-sports.io/football/teams/1531.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'A'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1532,
    'Algeria',
    'https://media.api-sports.io/football/teams/1532.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'J'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1533,
    'Cape Verde Islands',
    'https://media.api-sports.io/football/teams/1533.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'H'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1548,
    'Jordan',
    'https://media.api-sports.io/football/teams/1548.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'J'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1567,
    'Iraq',
    'https://media.api-sports.io/football/teams/1567.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'I'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1568,
    'Uzbekistan',
    'https://media.api-sports.io/football/teams/1568.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'K'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    1569,
    'Qatar',
    'https://media.api-sports.io/football/teams/1569.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'B'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    2380,
    'Paraguay',
    'https://media.api-sports.io/football/teams/2380.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'D'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    2382,
    'Ecuador',
    'https://media.api-sports.io/football/teams/2382.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'E'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    2384,
    'USA',
    'https://media.api-sports.io/football/teams/2384.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'D'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    2386,
    'Haiti',
    'https://media.api-sports.io/football/teams/2386.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'C'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    4673,
    'New Zealand',
    'https://media.api-sports.io/football/teams/4673.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'G'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    5529,
    'Canada',
    'https://media.api-sports.io/football/teams/5529.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'B'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;

INSERT INTO teams (external_id, name, logo, tournament_id, group_letter)
VALUES (
    5530,
    'Curaçao',
    'https://media.api-sports.io/football/teams/5530.png',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    'E'
) ON CONFLICT (name) DO UPDATE SET
    external_id   = EXCLUDED.external_id,
    logo          = EXCLUDED.logo,
    tournament_id = EXCLUDED.tournament_id,
    group_letter  = EXCLUDED.group_letter;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Removing seeded data cascades to fixtures, predictions, handicaps, etc.
-- Run only if you are sure this is safe in the target environment.
DELETE FROM teams       WHERE tournament_id = (SELECT id FROM tournaments WHERE slug = 'world-cup-2026');
DELETE FROM tournaments WHERE slug = 'world-cup-2026';

-- +goose StatementEnd
