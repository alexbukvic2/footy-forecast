-- +goose Up
-- +goose StatementBegin

-- Seeded from API-Football for tournament_id=7006468e-21cb-433a-a98e-d6b0d51fe0d7 (world-cup-2026)
-- Idempotent: ON CONFLICT (external_id, tournament_id) updates name + team_id.

INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    730,
    'Thibaut Nicolas Marc Courtois',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162511,
    'Senne Lammens',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    340151,
    'Mike Louis Penders',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2920,
    'Timothy Castagne',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162007,
    'Maxim De Cuyper',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162141,
    'Koni De Winter',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    304228,
    'Zeno Koen Debast',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    69,
    'Brandon Mechele',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    264,
    'Thomas André A. Meunier',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    312964,
    'Nathan Ngoy',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    375974,
    'Joaquin Seys',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    204043,
    'Arthur Nicolas Theate',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    629,
    'Kevin De Bruyne',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162714,
    'Amadou Ba Zeund Georges Mvom Onana',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2120,
    'Nicolas Thierry Raskin',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1417,
    'Alexis Jesse Saelemaekers',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2926,
    'Youri Marion A. Tielemans',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    335056,
    'Diego Manuel Jadon da Silva Moreira',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    78,
    'Hans Vanaken',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20,
    'Axel Laurent Angel Lambert Witsel',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1422,
    'Jérémy Baffour Doku',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    340077,
    'Matias Fernandez-Pardo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    907,
    'Romelu Menama Lukaku Bolingoli',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25458,
    'Dodi Lukebakio Ngandoli',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    147859,
    'Charles Marc De Ketelaere',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1946,
    'Leandro Trossard',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    22221,
    'Mike Peterson Maignan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    347211,
    'Robin Risser Birckel',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    21628,
    'Brice Lauriche Samba',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2724,
    'Lucas Digne',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    161907,
    'Malo Gusto',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    33,
    'Lucas François Bernard Hernández',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47300,
    'Théo Bernard François Hernández',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1145,
    'Ibrahima Konaté',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1257,
    'Jules Olivier Koundé',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20995,
    'Maxence Guy Lacroix',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    22090,
    'William Alain André Gabriel Saliba',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1149,
    'Dayotchanculle Oswald Upamecano',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    156477,
    'Mathis Rayan Cherki',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2290,
    'N''Golo Kanté',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    22147,
    'Kouadio Emmanuel Koné',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    272,
    'Adrien Thibault Marie Rabiot-Provost',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1271,
    'Aurélien Djani Tchouaméni',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    336657,
    'Warren Zaïre-Emery',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    274300,
    'Maghnes Akliouche',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    161904,
    'Bradley Barcola',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25927,
    'Jean-Philippe Mateta',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    153,
    'Masour Ousmane Dembélé',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    343027,
    'Désiré Doué',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    278,
    'Kylian Mbappé Lottin',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    21509,
    'Marcus Lilian Thuram-Ulien',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19617,
    'Michael Akpovie Olise',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    524,
    'Dominik Kotarski',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1305,
    'Dominik Livaković',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    14268,
    'Ivor Pandur',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1902,
    'Duje Ćaleta-Car',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    30827,
    'Martin Erlić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    129033,
    'Joško Gvardiol',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1084,
    'Marin Pongračić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    125171,
    'Josip Stanišić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    14701,
    'Josip Šutalo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    387521,
    'Luka Vušković',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    295026,
    'Martin Baturina',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    284869,
    'Toni Fruk',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    14395,
    'Kristijan Jakić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2291,
    'Mateo Kovačić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    754,
    'Luka Modrić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1322,
    'Nikola Moro',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2763,
    'Mario Pašalić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    7332,
    'Luka Sučić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    348205,
    'Petar Sučić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    842,
    'Nikola Vlašić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    46746,
    'Ante Budimir',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    726,
    'Andrej Kramarić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    202696,
    'Igor Matanović',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    66055,
    'Petar Musa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    260865,
    'Marco Pašalić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    207,
    'Ivan Perišić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '3'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    158700,
    'Viktor Tobias Johansson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2851,
    'Bo Kristoffer Nordfeldt',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48033,
    'Jacob Mikael Widell Zetterström',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47903,
    'Hjalmar Ekdal',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47969,
    'Gabriel Johan Gudmundsson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    137976,
    'Isak Malcolm Kwaku Hien',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    137721,
    'Gustaf Johan Lagerbielke',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    889,
    'Victor Jörgen Nilsson Lindelöf',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47988,
    'Carl Anders Theodor Starfelt',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    161504,
    'Herman Nils Johansson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    8486,
    'Eric Anders Smith',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    226765,
    'Elliot Karl Stroud',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    198654,
    'Daniel Svensson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47696,
    'Alexander Olof Bernhardsson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    265820,
    'Yasin Abbas Ayari',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    347316,
    'Lucas Erik Holger Bergvall',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48002,
    'Erik Benjamin Nygren',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48047,
    'Jesper Kewe Karlström',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2860,
    'Kenneth Nlata Sema',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    30484,
    'Mattias Olof Svanberg',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    350850,
    'Besfort Zeneli',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    160925,
    'Taha Abdi Ali',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    153430,
    'Anthony David Junior Elanga',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18979,
    'Viktor Einar Gyökeres',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2864,
    'Alexander Isak',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    15683,
    'Håkan Gustaf Nilsson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    280,
    'Alisson Ramsés Becker',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    617,
    'Ederson Santana de Moraes',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2410,
    'Weverton Pereira da Silva',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    860,
    'Alex Sandro Lobo Silva',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    30497,
    'Gleison Bremer Silva Nascimento',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    22224,
    'Gabriel dos Santos Magalhães',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    30424,
    'Roger Ibañez da Silva',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    10124,
    'Leonardo Pereira',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    257,
    'Marcos Aoás Corrêa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    24866,
    'Douglas dos Santos Justino de Melo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    349001,
    'Wesley Vinícius França Lima',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    618,
    'Danilo Luiz da Silva',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    10135,
    'Bruno Guimarães Rodriguez Moura',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    747,
    'Carlos Henrique Casimiro',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    299,
    'Fábio Henrique Tavares',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1646,
    'Lucas Tolentino Coelho de Lima',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    275170,
    'Danilo dos Santos de Oliveira',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1165,
    'Matheus Santos Carneiro da Cunha',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    377122,
    'Endrick Felipe Moreira de Sousa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    265785,
    'Luiz Henrique André Rosa da Silva',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    127769,
    'Gabriel Teodoro Martinelli Silva',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    276,
    'Neymar da Silva Santos Júnior',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1496,
    'Raphael Dias Belloli',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    407806,
    'Rayan Vitor Simplício Rocha',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    196156,
    'Igor Thiago Nascimento Rodrigues',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    762,
    'Vinícius José Paixão de Oliveira Júnior',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '6'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    61895,
    'Santiago Andrés Mele Castanero',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    429,
    'Néstor Fernando Muslera Micol',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50077,
    'Sergio Ramón Rochet Álvarez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    101814,
    'Ronald Federico Araújo da Silva',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    135334,
    'Santiago Ignacio Bueno Sciutto',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51535,
    'Sebastián Enzo Cáceres Ramos',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    31,
    'José María Giménez de Vargas',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47254,
    'Mathías Olivera Miramontes',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1290,
    'Guillermo Varela Olivera',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51572,
    'Matías Nicolás Viña Susperreguy',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    863,
    'Rodrigo Bentancur Colmán',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    153083,
    'Emiliano Martínez Toranza',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162891,
    'Juan Manuel Sanabria Magole',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51494,
    'Manuel Ugarte Ribeiro',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    756,
    'Federico Santiago Valverde Dipetta',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    108563,
    'Rodrigo Zalazar Martínez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2612,
    'Giorgian Daniel de Arrascaeta Benedetti',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    5995,
    'Diego Nicolás de la Cruz Arcosa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51776,
    'Maximiliano Javier Araújo Vilches',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51603,
    'Agustín Canobbio Graviz',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    16482,
    'Rodrigo Sebastián Aguirre Soto',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51466,
    'Joaquín Piquerez Moreira',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51617,
    'Darwin Gabriel Núñez Ribeiro',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    70078,
    'Facundo Pellistri Rebollo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51618,
    'Paul Brian Rodríguez Bravo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51530,
    'Federico Sebastián Viñas Barboza',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '7'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2481,
    'Álvaro David Montero Perales',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    313,
    'David Ospina Ramírez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2482,
    'Camilo Andrés Vargas Gil',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    30,
    'Santiago Arias Naranjo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    13571,
    'Willer Emilio Ditta Pérez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1929,
    'Jhon Janer Lucumí Bonilla',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2483,
    'Deiver Andrés Machado Mena',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2484,
    'Yerry Fernando Mina González',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    64268,
    'Johan Andrés Mojica Palacio',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    13736,
    'Daniel Muñoz Mejía',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    324034,
    'Gustavo Adolfo Puerta Molano',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    168,
    'Davinson Sánchez Mina',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    289592,
    'Kevin Duvan Castaño Gil',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2490,
    'Jefferson Andrés Lerma Solís',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    13151,
    'Juan Camilo Portilla Orozco',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    13708,
    'Jhon Adolfo Arias Andrade',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    5994,
    'Jorge Andrés Carrascal Guardo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6005,
    'Juan Fernando Quintero Paniagua',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    195104,
    'Richard Ríos Montoya',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    517,
    'James David Rodríguez Rubio',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    13376,
    'Jaminton Leandro Campaz',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    24810,
    'Jhon Andrés Córdoba Copete',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47582,
    'Juan Camilo Hernández Suárez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2489,
    'Luis Fernando Díaz Marulanda',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    345748,
    'Carlos Andrés Gómez Hinestroza',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47237,
    'Luis Javier Suárez Charris',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '8'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    182718,
    'Joan García Pons',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19465,
    'David Raya Martin',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47270,
    'Unai Simón Mendibil',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    396623,
    'Pau Cubarsí Paredes',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47380,
    'Marc Cucurella Saseta',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    619,
    'Eric García Martret',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    563,
    'Alejandro Grimaldo García',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    622,
    'Aymeric Jean Louis Gérard Alph Laporte',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    753,
    'Marcos Llorente Moreno',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47519,
    'Pedro Antonio Porro Sauceda',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    295793,
    'Marc Pubill Pagès',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    296667,
    'Pablo Martín Páez Gavira',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47311,
    'Mikel Merino Zazón',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    133609,
    'Pedro González López',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44,
    'Rodrigo Hernández Cascante',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    328,
    'Fabián Ruiz Peña',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47315,
    'Martín Zubimendi Ibáñez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    182219,
    'Alejandro Baena Rodríguez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1323,
    'Daniel Olmo Carvajal',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47348,
    'Borja Iglesias Quintás',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    386828,
    'Lamine Yamal Nasraoui Ebana',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    338751,
    'Víctor Muñoz Villanueva',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47323,
    'Mikel Oyarzabal Ugarte',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    184226,
    'Yeremy Jesús Pino Santos',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    931,
    'Ferran Torres García',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    183799,
    'Nicholas Williams Arthuer',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '9'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19088,
    'Dean Bradley Henderson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2932,
    'Jordan Lee Pickford',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162489,
    'James Harrington Trafford',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18961,
    'Daniel Johnson Burn',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    67971,
    'Addji Keaninkin Marc-Israel Guéhi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19545,
    'Reece Lewis James',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19354,
    'Ezri Konsa Ngoyo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    158694,
    'Valentino Francisco Livramento',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    307123,
    'Nico O''Reilly',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    158698,
    'Jarell Amorin Quansah',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19235,
    'Diop Tehuti Djed-Hotep Spence',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    626,
    'John Stones',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    138908,
    'Elliot Anderson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    129718,
    'Jude Victor William Bellingham',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19586,
    'Eberechi Oluchi Eze',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    292,
    'Jordan Brian Henderson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    284322,
    'Kobbie Boateng Mainoo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2937,
    'Declan Rice',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19170,
    'Morgan Elliot Rogers',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    138787,
    'Anthony Michael Gordon',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    184,
    'Harry Edward Kane',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    136723,
    'Chukwunonso Tristan Madueke',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    909,
    'Marcus Rashford',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1460,
    'Bukayo Ayoyinka Temidayo Saka',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19974,
    'Ivan Benjamin Elijah Toney',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19366,
    'Oliver George Arthur Watkins',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '10'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2967,
    'Luis Ricardo Mejía Cajar',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2968,
    'Orlando Mosquera',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    57861,
    'César Jair Samudio Murillo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    7197,
    'Andrés Alberto Andrade Cedeño',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2975,
    'César Rodolfo Blackman Camarena',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    57739,
    'José Ángel Córdoba Chambers',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2970,
    'Éric Javier Davis Grajales',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2971,
    'Fidel Escobar Mendieta',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    328148,
    'Edgardo Isaac Fariña Wynter',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    57803,
    'Jorge Abdiel Gutiérrez Cornejo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    71230,
    'Carlos Miguel Harvey Cesneros',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    81659,
    'Roderick Alonso Miller Molina',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2973,
    'Michael Amir Murillo Bermúdez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    57667,
    'Jiovany Javier Ramos Díaz',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    57807,
    'Adalberto Eliécer Carrasquilla Alcázar',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2977,
    'Aníbal Casis Godoy Lemus',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    554208,
    'C. Martinez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2979,
    'José Luis Rodríguez Francis',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    46855,
    'Édgar Yoel Bárcenas Herrera',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    96615,
    'Ismael Díaz de León',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2978,
    'Alberto Abdiel Quintero Medina',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    57875,
    'César Augusto Yanis Velasco',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2983,
    'José Fajardo Nelson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    292396,
    'Azarías Emmanuel Londoño González',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    57910,
    'Tomás Abdiel Rodríguez Mena',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51648,
    'Cecilio Alfonso Waterman Ruíz',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '11'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    304782,
    'Tomoki Hayakawa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    33034,
    'Keisuke Osako',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    199578,
    'Zion Suzuki',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    38114,
    'Ko Itakura',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    32893,
    'Hiroki Ito',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    440,
    'Yuto Nagatomo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    33165,
    'Ayumu Seko',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    32887,
    'Yukinari Sugawara',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    351014,
    'Junnosuke Suzuki',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    32954,
    'Shogo Taniguchi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2597,
    'Takehiro Tomiyasu',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    32858,
    'Tsuyoshi Watanabe',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2598,
    'Ritsu Dōan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    8500,
    'Wataru Endo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2601,
    'Daichi Kamada',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    33889,
    'Kaishu Sano',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    32966,
    'Ao Tanaka',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1942,
    'Junya Ito',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    32862,
    'Takefusa Kubo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    33224,
    'Daizen Maeda',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    33321,
    'Keito Nakamura',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    199143,
    'Yuito Suzuki',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    375930,
    'Keisuke Goto',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    33289,
    'Koki Ogawa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    422572,
    'Kento Shiogai',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    72155,
    'Ayase Ueda',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '12'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    119853,
    'Mory Diaw',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20566,
    'Yehvann Diouf',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2986,
    'Édouard Osoque Mendy',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    409303,
    'El Hadji Malick Diouf',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    158121,
    'Ismail Joshua Jakobs',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    318,
    'Kalidou Koulibaly',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    313937,
    'Antoine Mendy',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25916,
    'Moussa Niakhaté',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    276184,
    'Mamadou Sarr',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    8450,
    'Abdoulaye Seck',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    81,
    'Krépin Diatta',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    374058,
    'Lamine Camara',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    41552,
    'Pathé Ismaël Ciss',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    327631,
    'Mouhamadou Habib Diarra',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2990,
    'Idrissa Gana Gueye',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20696,
    'Pape Alassane Gueye',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    630895,
    'Bara Ndiaye',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    237129,
    'Pape Matar Sarr',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    400948,
    'Assane Diao Diaoune',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    284072,
    'Cheikh Ahmadou Bamba Mbacké Dieng',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    283058,
    'Nicolas Jackson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    304,
    'Sadio Mané',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    446249,
    'Ibrahim Mbaye',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    14379,
    'Pape Cherif Ndiaye',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18592,
    'Iliman Cheikh Baroy Ndiaye',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2218,
    'Ismaïla Sarr',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '13'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    123468,
    'Marvin Keller',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25282,
    'Gregor Kobel',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1142,
    'Yvon Landry Mvogo Nganoma',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    5,
    'Manuel Obafemi Akanji',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162414,
    'Aurèle Florian Amenda',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48372,
    'Eray Ervin Cömert',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2803,
    'Nico Elvedi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    349344,
    'Luca Antony Jaquez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48489,
    'Miro Max Maria Muheim',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1631,
    'Ricardo Iván Rodríguez Araya',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48378,
    'Silvan Dominic Widmer',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    951,
    'Michel Aebischer',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2807,
    'Remo Marco Freuler',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    264705,
    'Ardon Jashari',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    957,
    'Mohameth Djibril Ibrahima Sow',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1464,
    'Granit Xhaka',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2810,
    'Denis Lemi Zakaria Lako Lado',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    952,
    'Christian Andreas Fassnacht',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48497,
    'Cedric Jan Itten',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    163032,
    'Fabian Rieder',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48471,
    'Rubén Estephan Vargas Martínez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    123469,
    'Mohamed Zeki Amdouni',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    421,
    'Breel Donald Embolo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    406244,
    'Johan Manzambi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48389,
    'Noah Arinzechukwu Okafor',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48648,
    'Dan Assane Ndoye',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '15'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    35769,
    'Carlos Acevedo López',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2098,
    'Francisco Guillermo Ochoa Magaña',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    270774,
    'José Raúl Rangel Aguilar',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2881,
    'Jesús Daniel Gallardo Vasconcelos',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2873,
    'César Jasib Montes Castro',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    127227,
    'Israel Reyes Romero',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2878,
    'Jorge Eduardo Sánchez Ramos',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    35544,
    'Johan Felipe Vázquez Ibarra',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2869,
    'Edson Omar Álvarez Velázquez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    390002,
    'Mateo Chávez García',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    313383,
    'Obed Gómez Vargas',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    35690,
    'Luis Gerardo Chávez Magallón',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    750,
    'Álvaro Fidalgo Fernández',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    212233,
    'Brian Gutiérrez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    266345,
    'Érik Antonio Lira Méndez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    482605,
    'Gilberto Rafael Mora Zambrano',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    35576,
    'Orbelín Pineda Alvarado',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    35970,
    'Luis Francisco Romo Barrón',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2879,
    'Roberto Carlos Alvarado Hernández',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    94562,
    'Santiago Tomás Giménez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    291713,
    'Armando González Alba',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    36111,
    'César Saúl Huerta Valera',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2887,
    'Raúl Alonso Jiménez Rodríguez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    36088,
    'Guillermo Martínez Ayala',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    35532,
    'Julián Andrés Quiñones Quiñones',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2889,
    'Ernesto Alexis Vega Rojas',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '16'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2890,
    'Hyeon-Woo Jo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2892,
    'Seung-Gyu Kim',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    34374,
    'Bum-Keun Song',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    547307,
    'Cho Wi-Je',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2897,
    'Min-Jae Kim',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2912,
    'Moon-Hwan Kim',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    34418,
    'Tae-Hyeon Kim',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    237218,
    'Han-Beom Lee',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    99211,
    'Jin-Seop Park',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    237220,
    'Tae-Seok Lee',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    197985,
    'Young-Woo Seol',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    280358,
    'Jens Castrop',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    304951,
    'Gi-Hyuk Lee',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    357286,
    'Jun-Ho Bae',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    237050,
    'Ji-Sung Eom',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2901,
    'In-Beom Hwang',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    34168,
    'Jin-Gyu Kim',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    34431,
    'Dong-Gyeong Lee',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2906,
    'Jae-Sung Lee',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    927,
    'Kang-In Lee',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2909,
    'Seung-Ho Paik',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    304958,
    'Hyun-Jun Yang',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    24888,
    'Hee-Chan Hwang',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    34211,
    'Gue-Sung Cho',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    34710,
    'Hyeon-Gyu Oh',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    186,
    'Heung-Min Son',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '17'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    353883,
    'Patrick Beach',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6870,
    'Paul Izzo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2741,
    'Mathew David Ryan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    225,
    'Aziz Eraltay Behich',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    337587,
    'Jordan Jacob Bos',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20457,
    'Cameron Robert Burgess',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    348568,
    'Alessandro Circati',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    33847,
    'Jason Kato Geria',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    426480,
    'Lucas Herrington',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20079,
    'Harry James Souttar',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2742,
    'Miloš Degenek',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6808,
    'Jacob Michael Italiano',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    153622,
    'Kai Clifton Trewin',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2749,
    'Jackson Alexander Irvine',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6832,
    'Cameron Peter Devlin',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    7050,
    'Aiden Connor O''Neill',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    441269,
    'Paul Michael Junior Okon-Engstler',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6904,
    'Connor Isaac Metcalfe',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    38123,
    'Ajdin Hrustić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    342035,
    'Cristian Volpato',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    338014,
    'Nestory Irankunda',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2751,
    'Mathew Allan Leckie',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2755,
    'Awer Bul Mabil',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    296645,
    'Tete Yengi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    198352,
    'Mohamed Touré',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    312459,
    'Nishan Velupillay',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '20'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2682,
    'Alireza Safar Beiranvand',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    29755,
    'Seyed Hossein Hosseini',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2681,
    'Seyed Payam Niazmand Ghader',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    532950,
    'D. Eiri',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2691,
    'Ramin Rezaeian Semeskandi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2685,
    'Ehsan Hajisafi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    136880,
    'Saleh Hardani',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2687,
    'Mohammad Hossein Kanani Zadegan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    29704,
    'Shoja Khalilzadeh',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2688,
    'Milad Mohammadi Keshmarzi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    533035,
    'A. Nemati',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    343405,
    'Aria Yousefi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2700,
    'Alireza Jahanbakhsh Jirandeh',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2683,
    'Roozbeh Cheshmi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19614,
    'Saeed Ezatolahi Afagh',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2699,
    'Sayed Saman Ghoddos',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    423753,
    'Amir Mohammad Razzaghinia',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    341844,
    'Mohammad Ghorbani',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    134217,
    'Mohammad Mohebi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2697,
    'Mehdi Torabi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    643918,
    'M. Ghaedi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    29937,
    'Amirhossein Hosseinzadeh',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    29720,
    'Ali Alipourghara',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    89982,
    'Shahriyar Moghanlou',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    37892,
    'Dennis-Yerai Eckert Ayensa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    42315,
    'Mehdi Taremi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '22'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    193288,
    'Nawaf Dhahi Faisal Al Shuweiti Al Aqidi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44449,
    'Ahmed Ali Al Kassar',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44411,
    'Mohammed Khalil Ibrahim Al Owais',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44594,
    'Saud Abdullah Salem Abdulhamid',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    403087,
    'Mohammed Waheeb Saeed Abu Al Shamat',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    134995,
    'Nawaf Mashari Abdulrahman Boushal',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44475,
    'Abdulelah Ali Awadh Al Amri',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    202381,
    'Moteb Saad Salim Al Naqi Al Harbi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44335,
    'Hassan Kadesh Yahya Mahboob',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44507,
    'Ali Mohammed Ali Lajami',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44367,
    'Ali Hassan Muhammad Majrashi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44362,
    'Hassan Mohammed Osama Al Tambakti',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    543059,
    'J. Thakri',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44339,
    'Nasser Essa Shafi Al Shardan Al Dawsari',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44349,
    'Mohamed Ibrahim Abdullah Kanno',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    269172,
    'Ziyad Mubarak Eid Al Marwani Al Johani',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    306380,
    'Musab Fahad Zaid Al Juwayr',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44315,
    'Abdullah Mohammed Hamza Al Ber Al Khaibari',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    593759,
    'Ala Al Haji',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44701,
    'Khalid Essa Muhammad Al Ghannam',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2639,
    'Sultan Ahmed Mohammed Mandash',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    147812,
    'Ayman Yahya Salem Ahmed',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44324,
    'Feras Tariq Nasser Al Brikan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44382,
    'Abdullah Abdulrahman Abdullah Al Hamdan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44551,
    'Saleh Khalid Mohammed Al Shehri',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44340,
    'Salem Mohammed Shafi Al Dawsari',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '23'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    702,
    'Oliver Baumann',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    497,
    'Manuel Peter Neuer',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    399,
    'Alexander Nübel',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25368,
    'Waldemar Riptsov Anton',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    280074,
    'Nathaniel Brown',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25158,
    'David Raum',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2285,
    'Antonio Rüdiger',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    26243,
    'Nico Cedric Schlotterbeck',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    972,
    'Jonathan Glao Tah',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    163189,
    'Malick Laye Thiaw',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    502,
    'Joshua Walter Kimmich',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    714,
    'Nadiem Amiri',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    511,
    'Leon Christoph Goretzka',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18970,
    'Pascal Groß',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    494131,
    'Lennart Karl',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    128533,
    'Jamie Leweling',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    181812,
    'Jamal Musiala',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    637,
    'Felix Kalu Nmecha',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    328033,
    'Aleksandar Pavlović',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    644,
    'Leroy Aziz Sané',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    137210,
    'Angelo Stiller',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    203224,
    'Florian Richard Wirtz',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    978,
    'Kai Lukas Havertz',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    158644,
    'Maximilian Beier',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    26475,
    'Deniz Undav',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    158054,
    'Nick Woltemade',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '25'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19599,
    'Damián Emiliano Martínez Romero',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2465,
    'Juan Agustín Musso',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47296,
    'Gerónimo Rulli',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6,
    'Leonardo Julián Balerdi Rosa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2467,
    'Lisandro Martínez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6231,
    'Facundo Axel Medina',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6503,
    'Nahuel Molina Lucero',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2468,
    'Gonzalo Ariel Montiel',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    624,
    'Nicolás Hernán Gonzalo Otamendi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    30776,
    'Cristian Gabriel Romero',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    529,
    'Nicolás Alejandro Tagliafico',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    319572,
    'Valentín Barco',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    5996,
    'Enzo Jeremías Fernández',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1578,
    'Giovani Lo Celso',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6716,
    'Alexis Mac Allister',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6002,
    'Exequiel Alejandro Palacios',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    271,
    'Leandro Daniel Paredes',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    26315,
    'Nicolás Iván González',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2472,
    'Rodrigo Javier De Paul',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6067,
    'Thiago Ezequiel Almada',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6009,
    'Julián Álvarez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    350037,
    'Nicolás Paz Martínez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    295513,
    'José Manuel Alberto López',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    217,
    'Lautaro Javier Martínez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    154,
    'Lionel Andrés Messi Cuccittini',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    323935,
    'Giuliano Simeone Baldini',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '26'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    369,
    'Diogo Meireles da Costa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1590,
    'José Pedro Malheiro de Sá',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    46672,
    'Rui Tiago Dantas da Silva',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    161939,
    'Tomás Lemos Araújo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    855,
    'João Pedro Cavaco Cancelo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    886,
    'José Diogo Dalot Teixeira',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    567,
    'Rúben dos Santos Gato Alves Dias',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    265595,
    'Gonçalo Bernardo Inácio',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    263482,
    'Nuno Alexandre Tavares Mendes',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    130,
    'Nélson Cabral Semedo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    336671,
    'Renato da Palma Veiga',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    190485,
    'Samuel de Almeida Costa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1485,
    'Bruno Miguel Borges Fernandes',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    335051,
    'João Pedro Gonçalves Neves',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2676,
    'Rúben Diogo da Silva Neves',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    41621,
    'Matheus Luiz Nunes',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    636,
    'Bernardo Mota Veiga de Carvalho e Silva',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    128384,
    'Vítor Machado Ferreira',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    161585,
    'Francisco Fernandes da Conceição',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    925,
    'Gonçalo Manuel Ganchinho Guedes',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    583,
    'João Félix Sequeira',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    22236,
    'Rafael Alexandre da Conceição Leão',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1864,
    'Pedro Lomba Neto',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    41585,
    'Gonçalo Matias Ramos',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    874,
    'Cristiano Ronaldo dos Santos Aveiro',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    41112,
    'Francisco António Machado Mota de Castro Trincão',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '27'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    49423,
    'Sabri Ben Hsan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    533394,
    'C. Abdelmouhib',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    49424,
    'Aymen Dahmen',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    49583,
    'Ali El Abdi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    393977,
    'Adam Arous',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    135059,
    'Mohamed  Amine Ben Hmida',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2945,
    'Dylan Daniel Mahmoud Bronn',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    533360,
    'R. Chikhaoui',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    375608,
    'Moutaz Neffati',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    163068,
    'Omar Rekik',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50030,
    'Montassar Omar Talbi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18942,
    'Yan Valery',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    49469,
    'Ahmed Mortadha Ben Ouanes',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    533295,
    'K. Ayari',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    310196,
    'Ismaël Seifallah Gharbi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25300,
    'Rani Khedira',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    180560,
    'Hannibal Mejbri',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    21587,
    'Ellyes Joris Skhiri',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    199310,
    'Anis Ben Slimane',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    67195,
    'Mohamed Belhadj Mahmoud',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    57518,
    'Sebastian Tounekti',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    42012,
    'Mohamed Elias Achouri',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2962,
    'Firas Chaouat',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    566059,
    'Rayan Elloumi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    344862,
    'Hazem Mastouri',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    323974,
    'Elias Saad',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '28'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2701,
    'Yassine Bounou',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2702,
    'Munir Mohand Mohamedi El Kajoui',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2703,
    'Ahmed Reda Tagnaouti',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    21694,
    'Nayef Aguerd',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18814,
    'Issa Laye Lucas Jean Diop',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    283252,
    'Zakaria El Ouahdi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    9,
    'Achraf Hakimi Mouh',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    326183,
    'Redouane Halhal',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    545,
    'Noussair Mazraoui',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    278898,
    'Chadi Riad Dnanou',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162451,
    'Anass Salah-Eddine',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    146772,
    'Youssef Belammari',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    74,
    'Sofyan Amrabat',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    438688,
    'Ayyoub Bouaddi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    277003,
    'Neil El Aynaoui',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    340573,
    'Bilal El Khannouss',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    415431,
    'Samir El Mourabet',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    129678,
    'Azzedine Ounahi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    161897,
    'Ismael Saibari Ben El Basra',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    336659,
    'Chemsdine Talbi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    369544,
    'Gessime Yassine',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    744,
    'Brahim Abdelkader Díaz',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    535046,
    'A. Amaimouni',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2722,
    'Ayoub El Kaabi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    181421,
    'Abdessamad Ezzalzouli',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    36579,
    'Soufiane Rahimi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '31'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    550469,
    'M. Alaa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    16797,
    'Mohamed El Sayed Mohamed El Sh Gomaa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    269174,
    'Mostafa Ahmed Abdelaziz Mohamed Shobeir',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    16831,
    'Al Mahdi Soliman',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    269621,
    'Hossam Abdelmaguid Abdelsalam Abdelmaguid',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    196343,
    'Mohamed Abdelmonem El Sayed Mohamed Ahmed',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2649,
    'Ahmed Mohamed Abou El Fotouh Mohamed',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2656,
    'Karim Hafez Ramadan Seif El Din',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2654,
    'Mohamed Hany Gamal Eldemerdash',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    16804,
    'Yasser Ahmed Ibrahim El Hanafi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    16805,
    'Rami Hisham Abdel Aziz Rabia',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    550371,
    'T. Alaa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    17269,
    'Emam Ashour Metwally Abdelghany',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    190575,
    'Marwan Attia Fahim Ghallab',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2660,
    'Nabil Emad Al El Mahdy',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    16813,
    'Hamdi Fathy Abdelhalim Abdul Fattah',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    16841,
    'Mohanad Mostafa Ahmed Abdelmonem',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    69196,
    'Mahmoud Saber Abdelmohsen Hassan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    395075,
    'Mostafa Mohamed Zaki Abdelraouf',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2664,
    'Mahmoud Ahmed Ibrahim Hassan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20844,
    'Haissem Hassan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    664079,
    'Ahmed Zizo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    550547,
    'H. Abdelkarim',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    70535,
    'Ibrahim Adel Ali Mohamed Hassan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    81573,
    'Omar Khaled Mohamed Abd Elsala Marmoush',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    306,
    'Mohamed Salah Hamed Mahrous Ghaly',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '32'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    269349,
    'Lukáš Horníček',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    138804,
    'Matěj Kovář',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    66347,
    'Jindřich Staněk',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    337740,
    'Štěpán Chaloupek',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1231,
    'Vladimír Coufal',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2252,
    'Tomáš Holeš',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162964,
    'Robin Hranáč',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    128793,
    'David Jurásek',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    66407,
    'Ladislav Krejčí',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1237,
    'Jaroslav Zelený',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    128772,
    'David Zima',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    66214,
    'David Douděra',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162194,
    'Lukáš Červ',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25348,
    'Vladimír Darida',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    66353,
    'Lukáš Provod',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    241,
    'Michal Sadílek',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    555437,
    'H. Sochurek',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    386837,
    'Alexandr Sojka',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1243,
    'Tomáš Souček',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    66387,
    'Pavel Šulc',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    818,
    'Tomáš Chorý',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    66275,
    'Mojmír Chytil',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    66019,
    'Adam Hložek',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    66340,
    'Jan Kuchta',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    794,
    'Patrik Schick',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    290212,
    'Denis Višinský',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '770'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    7598,
    'Patrick Pentz',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    7525,
    'Alexander Schlager',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    221605,
    'Florian Wiegele',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    126640,
    'David Affengruber',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    505,
    'David Olatukunbo Alaba',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25287,
    'Kevin Danso',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25314,
    'Marco Friedl',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    26240,
    'Philipp Lienhart',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25915,
    'Phillipp Mwene',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    711,
    'Stefan Posch',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    7090,
    'Michael Svoboda',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    715,
    'Christoph Baumgartner',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    138935,
    'Carney Chibueze Chukwuemeka',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    719,
    'Florian Grillitsch',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1157,
    'Konrad Laimer',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    7327,
    'Alexander Prass',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1159,
    'Marcel Sabitzer',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1095,
    'Xaver Schlager',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    7562,
    'Romano Christian Schmid',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    417,
    'Alessandro André Schöpf',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    7328,
    'Nicolas Seiwald',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    327895,
    'Paul Wanner',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    126642,
    'Patrick Wimmer',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18830,
    'Marko Arnautović',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25297,
    'Michael Gregoritsch',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    7722,
    'Saša Kalajdžić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '775'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50132,
    'Altay Bayındır',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    49866,
    'Uğurcan Çakır',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    49837,
    'Fehmi Mert Günok',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    62490,
    'Samet Akaydin',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    61837,
    'Abdülkerim Bardakcı',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    22222,
    'Mehmet Zeki Çelik',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    30521,
    'Merih Demiral',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50057,
    'Evren Eren Elmalı',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    26300,
    'Ozan Muhammed Kabak',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1361,
    'Ferdi Erenay Kadıoğlu',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1719,
    'Mert Müldür',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18776,
    'Çağlar Söyüncü',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1640,
    'Hakan Çalhanoğlu',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25448,
    'Kaan Ayhan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    37155,
    'Orkun Kökçü',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    24807,
    'Salih Özcan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    214463,
    'İsmail Yüksek',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    339887,
    'Can Yılmaz Uzun',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    134590,
    'Oğuz Aydın',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    291964,
    'Arda Güler',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    454,
    'Yunus Akgün',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    142959,
    'Muhammed Kerem Aktürkoğlu',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    388570,
    'Deniz Daniel Gül',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    49857,
    'İrfan Can Kahveci',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    339883,
    'Kenan Yıldız',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    63274,
    'Barış Alper Yılmaz',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '777'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19172,
    'Ørjan Håskjold Nyland',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    39082,
    'Egil Selvik',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    264378,
    'Sander Tangvik',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1119,
    'Kristoffer Vassbakk Ajer',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    39058,
    'Fredrik André Bjørkan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    180937,
    'Henrik Sælebakke Falchener',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    265444,
    'Sondre Klingen Langås',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    39254,
    'Torbjørn Lysaker Heggem',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18967,
    'Leo Skiri Østigård',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    39362,
    'Marcus Holmgren Pedersen',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    265782,
    'David Møller Wolfe',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    277930,
    'Thelo  Gerard Aasgaard',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    39043,
    'Fredrik Aursnes',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    39064,
    'Patrick Berg',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1934,
    'Sander Gard Bolin Berge',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    278133,
    'Oscar Bobb',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    37127,
    'Martin Ødegaard',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    36980,
    'Morten Thorsby',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    39143,
    'Kristian Thorstvedt',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    301528,
    'Andreas Rædergård Schjelderup',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    39073,
    'Jens Petter Hauge',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1100,
    'Erling Braut Haaland',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2032,
    'Jørgen Strand Larsen',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    314511,
    'Antonio Eromonsele Nordby Nusa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    8492,
    'Alexander Sørloth',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    24845,
    'Julian Ryerson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1090'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1106,
    'Craig Sinclair Gordon',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18933,
    'Angus Fraser James Gunn',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44937,
    'Liam Patrick Kelly',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19066,
    'Grant Campbell Hanley',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1111,
    'Jack William Hendry',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44871,
    'Aaron Buchanan Hickey',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19987,
    'Dominic John Hyam',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44811,
    'Scott Fraser McKenna',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    138417,
    'Nathan Kenneth Patterson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1115,
    'Anthony Ralston',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    289,
    'Andrew Henry Robertson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44865,
    'John Francis Souttar',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1117,
    'Kieran Tierney',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1125,
    'Ryan Christie',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44814,
    'Lewis Ferguson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    557460,
    'T. Fletcher',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19191,
    'John McGinn',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19077,
    'Kenneth McLean',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    903,
    'Scott Francis McTominay',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19524,
    'Ché Zach Everton Fred Adams',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    45307,
    'Lyndon John Dykes',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    343576,
    'Ben Gannon Doak',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    8794,
    'George David Eric Hirst',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    45175,
    'Lawrence Shankland',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    45078,
    'Ross Cameron Stewart',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    433272,
    'Findlay Curtis',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1108'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    446130,
    'Mladen Jurkas',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    9026,
    'Nikola Vasilj',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    108370,
    'Martin Zlomislić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    395589,
    'Nidal Čelik',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    7318,
    'Amar Dedić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    53517,
    'Dennis Hadžikadunić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1741,
    'Nikola Katić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1442,
    'Sead Kolašinac',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    271350,
    'Tarik Muharemović',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    76867,
    'Nihad Mujakić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    14301,
    'Stjepan Radeljić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162222,
    'Ivan Bašić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    70514,
    'Armin Gigović',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50006,
    'Amir Hadžiahmetović',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    340173,
    'Ermin Mahmić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1324,
    'Ivan Šunjić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    264094,
    'Benjamin Tahirović',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    322101,
    'Amar Memić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25129,
    'Dženis Burnić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    395559,
    'Kerim-Sam Alajbegović',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    329409,
    'Esmir Bajraktarević',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    314377,
    'Samed Baždar',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    46930,
    'Ermedin Demirović',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    790,
    'Edin Džeko',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    77037,
    'Jovo Lukić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    28382,
    'Haris Tabaković',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1113'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    26232,
    'Mark Flekken',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    194536,
    'Robin Roefs',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    129058,
    'Bart Verbruggen',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18861,
    'Nathan Benjamin Aké',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    226,
    'Denzel Justus Morris Dumfries',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    341642,
    'Jorrel Hato',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    38746,
    'Jurriën David Norman Timber',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    290,
    'Virgil van Dijk',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    38695,
    'Jan Paul van Hecke',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    152849,
    'Micky van de Ven',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    92993,
    'Mats Wieffer',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    542,
    'Ryan Jiro Gravenberch',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    36899,
    'Teun Koopmeiners',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    36902,
    'Tijjani Reijnders',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    36905,
    'Guus Berend Til',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    38747,
    'Quinten Ryan Crispito Timber',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    538,
    'Frenkie de Jong',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    30432,
    'Marten Elco de Roon',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    792,
    'Justin Dean Kluivert',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    38750,
    'Brian Ebenezer Adjei Brobbey',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    667,
    'Memphis Depay',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    247,
    'Cody Mathès Gakpo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    544,
    'Noa Noëll Lang',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    249,
    'Donyell Malen',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    37724,
    'Crysencio Jilbert Sylverio Cir Summerville',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25416,
    'Wout Weghorst',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1118'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    64190,
    'Yahia Fofana',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    277046,
    'Mohamed Koné',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    30393,
    'Alban-Marc Lafont',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    135068,
    'Emmanuel Elysee Djedje Agbadou Badobre',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    354753,
    'Ousmane Diomande',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    161747,
    'Guéla Maho Lewis Doué',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    22002,
    'Ghislain N''Clomande Konan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48119,
    'Kouakou Odilon Dorgeless Kossounou',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1807,
    'Obite Evan Ndicka',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20836,
    'Christopher Téa Operi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    30504,
    'Wilfried Stephane Singo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    30807,
    'Seko Mohamed Fofana',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    474591,
    'Christ Ravynel Inao Oulaï',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1642,
    'Franck Yannick Kessié',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    22149,
    'Ibrahim Sangaré',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    3243,
    'Jean Michaël Seri',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    263228,
    'Parfait Guiagon',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    387643,
    'Bazoumana Touré',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    301771,
    'Simon Adingra',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    275651,
    'Ange-Yoan Bonny',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    334429,
    'Oumar Diakité',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    157997,
    'Amad Diallo Traoré',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    513776,
    'Yan Diomande',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    137303,
    'Evann Guessand',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    3246,
    'Nicolas Pépé',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162707,
    'Sepe Elye Wahi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1501'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    480108,
    'Benjamin Asare',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    144709,
    'Joseph Tetteh Anang',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    3412,
    'Lawrence Ati Zigi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    369425,
    'Jonas Adjei Adjetey',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25341,
    'Derrick Luckassen',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    7578,
    'Gideon Mensah',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    196187,
    'Alidu Seidu',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    137223,
    'Jerome Osei Opoku',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    404172,
    'Kojo Peprah Oppong',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    15900,
    'Khalid Abdul Mumin Suleman',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    191240,
    'Marvin Senaya',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    21996,
    'Abdul Rahman Baba',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    21010,
    'Elisha Owusu',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    49,
    'Thomas Teye Partey',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    3608,
    'Kwasi Sibo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    475575,
    'Caleb Marfo Yirenkyi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19281,
    'Antoine Serlom Semenyo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    337426,
    'Augustine Boakye',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    303467,
    'Issahaku Abdul Fatawu',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    410016,
    'Prince Kwabena Adu',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    82090,
    'Solomon Brandon Michael Clarke Thomas-Asante',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    3428,
    'Jordan Pierre Ayew',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    411800,
    'Christopher Bonsu Baah',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47294,
    'Iñaki Arthuer Dannis Williams',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    350856,
    'Ernest Nuamah Appiah',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    199837,
    'Kamaldeen Sulemana',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1504'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    314509,
    'Matthieu Luka Epolo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48501,
    'Timothy Bruce Fayulu',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    24012,
    'Lionel Mpasi-Nzau',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    8445,
    'Buduka Dylan Batubinsika',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    24245,
    'Gédéon Kalulu Kyatengwa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    199767,
    'Steve Nkanu Kapuadi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    21098,
    'Joris Kayembe Ditu',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18816,
    'Fuka-Arthur Masuaku Kawela',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    375,
    'Chancel Mbemba Mangulu',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19182,
    'Axel Tuanzebe',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18846,
    'Aaron Wan-Bissaka',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1424,
    'Edouard Kayembe Kayembe',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    21101,
    'Samuel Moutoussamy',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    375598,
    'Ngal''ayel Mukau',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48555,
    'Charles Monginda Pickel',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    365331,
    'Noah Junior Sadiki',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44791,
    'Aaron Tshibola',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    129670,
    'Nathanaël Mbuku',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    8627,
    'Théo Bongonda Mbul''Ofeko Batomboat',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    3033,
    'Cédric Bakambu',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20674,
    'Simon Bokoté Banza',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    279482,
    'Brian Kibambe Cipenga',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    3034,
    'Meschack Elia Lina',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47545,
    'Gaël Romeo Kakuta Mambenga',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    179699,
    'Fiston Kalala Mayele',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20649,
    'Yoane Wissa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1508'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    46417,
    'Sipho Chaine',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    46245,
    'Stuart Ricardo Goss',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    3275,
    'Ronwen Hayden Williams',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    163041,
    'Bradley Cross',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    430078,
    'Samukelo Kabini',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    392387,
    'Olwethu Makhanya',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    510799,
    'Mbekezeli Mbokazi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    46334,
    'Aubrey Maphosa Modiba',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    46601,
    'Khuliso Johnson Mudau',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    474630,
    'Kulumani Ndamane',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    406752,
    'Ime Okon',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    46458,
    'Nkosinathi Sibisi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    330174,
    'Tholo Thabang Matuludi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    268710,
    'Jayden Oswin Adams',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    194430,
    'Thalente Mbatha',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    3287,
    'Teboho Mokoena',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    158433,
    'Sphephelo S''Miso Sithole',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    3289,
    'Themba Zwane',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    179893,
    'Oswin Reagan Appollis',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    98936,
    'Lyle Brent Foster',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    201354,
    'Sekotori Evidence Makgopa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    354831,
    'Thapelo Maseko',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    414149,
    'Relebohile Mofokeng',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    295977,
    'Tshepang Moremi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    127429,
    'Iqraam Rayners',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    359561,
    'Kamogelo Sebelebele',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1531'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    4414,
    'Oussama Benbot',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    304229,
    'Melvin Feyçal Mastil',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    732,
    'Luca Zinedine Zidane Fernández',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    342740,
    'Achref Abada',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    21138,
    'Rayan Aït-Nouri',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    4561,
    'Zinéddine Belaïd',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    303362,
    'Rafik Belghali',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2194,
    'Amir Selmane Ramy Bensebaïni',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20869,
    'Samir Sophian Chergui',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    271539,
    'Jaouen Djimmy Hadjam',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1567,
    'Aïssa Mandi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    4569,
    'Mohamed Amine Tougai',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    658,
    'Houssem-Eddine Chaâbane Aouar',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    409,
    'Nabil Bentaleb',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    4399,
    'Hicham Boudaoui',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    380587,
    'Ibrahim Maza',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    327599,
    'Yacine Titraoui',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    129047,
    'Ramiz Larbi Zerrouki',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    276670,
    'Farès Chaïbi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    200139,
    'Mohammed El Amine Amoura',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    292924,
    'Ahmed Nadir Benbouali',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    329163,
    'Adil Boulbina',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    334915,
    'Farès Ghedjemis',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    85041,
    'Amine Ferid Gouiri',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    326067,
    'Anis Hadj-Moussa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    635,
    'Riyad Karim Mahrez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1532'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    41764,
    'Márcio Salomão Brazão da Rosa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    15304,
    'Josimar José Évora Dias',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    163200,
    'Carlos Joaquim Antunes dos Santos',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    41608,
    'Edilson Alberto Monteiro Sanches Borges',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    21997,
    'Logan Evans Costa',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    308689,
    'Sidny Lopes Cabral',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    69260,
    'Roberto Carlos Lopes',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    22139,
    'Steven Moreira',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    332056,
    'Wagner Fabrício Cardoso Pina',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    268806,
    'Kelvin Spencer Pires',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2331,
    'Ianique dos Santos Tavares',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    265556,
    'Telmo Emanuel Gomes Arcanjo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    37435,
    'Deroy D''Encarnação Duarte',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    37436,
    'Laros d''Encarnação Duarte',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    142016,
    'João Paulo Moreira Fernandes',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    96512,
    'Kevin Lenini Gonçalves Pereira de Pina',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50745,
    'Jamiro Gregory Monteiro Alvarenga',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    128007,
    'Jair Semedo Monteiro',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44612,
    'Garry Mendes Rodrigues',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    22265,
    'Nuno Miguel da Costa Jóia',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1494,
    'Jovane Eduardo Borges Cabral',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    113580,
    'Willy Johnson Semedo Afonso',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    291024,
    'Hélio Sandro Oliveira Alves Varela',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50270,
    'Ryan Isaac Mendes da Graça',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    309182,
    'Gilson Benchimol Tavares',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    343287,
    'Dailon Rocha Livramento do Rosario',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1533'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    140607,
    'Yazid Moeen Hasan Abu Layla',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    163884,
    'Abdallah Ra''ed Mahmoud Al Fakhouri',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    163908,
    'Noureddin Zaid Khaleel Bani Ateyah',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    542822,
    'H. Abu Al Dahab',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    542710,
    'M. Abu Hasheesh',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    102538,
    'Mohammad Abualnadi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    53900,
    'Yazan Mousa Mahmoud Abu Al Arab',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    213978,
    'Saed Ahmad Salameh Al Rosan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    664028,
    'A. Badawi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    310835,
    'Abdallah Mousa Musallam Naseeb',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    72140,
    'Saleem Obaid',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    53902,
    'Ehsan Nabil Farhan Haddad',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    310785,
    'Mohannad Mahmoud Saleh Abu Taha',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    651096,
    'M. Al Daoud',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    542768,
    'Ibrahim Sa''deh',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    213980,
    'Nizar Mahmoud Ahmed Al Rashdan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    140609,
    'Noor Al Deen Mahmoud Ali Al Rawabdeh',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    104242,
    'Rajaei Ayed Fadel Hasan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    198211,
    'Amer Rasem Adel Jamous',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    123530,
    'Mahmoud Nayef Ahmad Al Mardi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    72142,
    'Mohammed Abu Zurayq',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    568556,
    'O. Al Fakhouri',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    575283,
    'A. Azaizeh',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    164026,
    'Ali Iyad Ali Olwan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    432841,
    'Ibrahim Mohammad Abdallah Sabra',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    15286,
    'Mousa Mohammad Mousa Sulaiman Al Tamari',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1548'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    197933,
    'Ahmed Basil Fadhil Al Fadhli',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    53886,
    'Jalal Hassan Hachim',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    123802,
    'Fahad Talib Raheem',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    145465,
    'Hussein Ali',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    271444,
    'Merchas Ghazi Salih Doski',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    15769,
    'Frans Dhia Jirjis Haddad',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    427187,
    'Mustafa Saadoun Al Korji',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    42261,
    'Rebin Ghareeb Solaka Adhamat',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    296373,
    'Zayed Tahseen',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    542849,
    'A. Yahya',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    295394,
    'Munaf Younus Hashim Al Tekreeti',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    453073,
    'Akam Hashim Rahman',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    47792,
    'Amir Fouad Aboud Al Ammari',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    282065,
    'Youssef Wali Faeq Amyn',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    140747,
    'Ibraheem Bayesh Kamil Al Kaabawi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    284295,
    'Zidane Aamar Iqbal',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    626479,
    'Z. Ismaeel',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48129,
    'Aimar Hazar Sher',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48025,
    'Kevin Enkido Yakob',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    542644,
    'A. Jasim',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    299813,
    'Ali Ibrahim Karim Al Hamadi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    265448,
    'Marko Hussein Farji',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    49451,
    'Aymen Hussein Ghadhban',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    542697,
    'Meme',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    229112,
    'Ahmed Qasem',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    542842,
    'A. Y. Hashim',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1567'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    72120,
    'Botirali Ergashev',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    73507,
    'Abduvohid Ne''matov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    73534,
    'Utkir Yusupov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    72122,
    'Khojiakbar Alijonov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    34121,
    'Rustamjon Ashurmatov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    73510,
    'Umar Eshmurodov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    416964,
    'Behruzjon Karimov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    360114,
    'Abduqodir Khusanov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    73514,
    'Sherzod Nasrullayev',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    532759,
    'F. Sayfiev',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    358427,
    'Jakhongir Urozov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    416952,
    'A. Abdullaev',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    65571,
    'Avazbek Ulmasaliev',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    53836,
    'Jaloliddin Masharipov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    73520,
    'Azizjon Gʻaniyev',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    65576,
    'Jamshid Iskanderov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    53835,
    'Odiljon Hamrobekov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    73522,
    'Akmal Mozgovoy',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50272,
    'Otabek Shukurov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    263676,
    'Abbosbek Fayzullaev',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    363723,
    'Sherzod Esanov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    72127,
    'Oston O''runov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    53834,
    'Dostonbek Khamdamov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    65584,
    'Azizbek Amonov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    72128,
    'Igor Sergeev',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    53535,
    'Eldor Azamat Shomurodov',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1568'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    42207,
    'Mahmud Ibrahim Abunada',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    42021,
    'Meshaal Aissa Mohammed Barsham',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    136503,
    'Salah Zakaria Mohamed Mousa Hassan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    175439,
    'Homam Elamin Mohamed Ahmed',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    42060,
    'Sultan Hussain Al Braik',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    542542,
    'A. Al Hussain',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    542548,
    'A. Al Oui',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    200981,
    'Jassem Gaber Abdulsallam',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2532,
    'Boualem Khoukhi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    42288,
    'Lucas Michel Mendes',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    366516,
    'Gueye Seydinaissa Laye',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2530,
    'Pedro Miguel Carvalho Deus Correia',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2535,
    'Assim Omer Al Haj Madibo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    42286,
    'Ahmed Mohammed Al Ganehi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2533,
    'Abdulaziz Hatem Mohammed Abdullah',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2537,
    'Karim Boudiaf',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2539,
    'Ahmed Fathi Abdoun',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2545,
    'Hassan Khalid Hassan Al Haydos',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    283174,
    'Mohamed Naceur Al Manai',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    542541,
    'Y. Abdurisag',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2544,
    'Akram Hassan Afif Yahya Afif',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2542,
    'Ahmed Alaaeldin Abdelmotaal',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2543,
    'Almoez Ali Zainalabedeen Moham Abdulla',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    42075,
    'Edmilson Junior Paulo da Silva',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    423737,
    'Tahsin Mohammed Jamshid',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    42089,
    'Mohammed Muntari',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '1569'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    535737,
    'R. Fernandez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    70852,
    'Orlando David Gill Noldin',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51733,
    'Gastón Hernán Olveira Echeverría',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6168,
    'Omar Federico Alderete Fernández',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2499,
    'Júnior Osmar Ignacio Alonso Mujica',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    195992,
    'Juan José Cáceres',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2500,
    'Fabián Cornelio Balbuena González',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2502,
    'Gustavo Raúl Gómez Portillo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    425020,
    'Alexandro Maidana Mendieta',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    35808,
    'Víctor Gustavo Velázquez Ramos',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    70732,
    'José María Canale Domínguez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    195107,
    'Damián Josué Bobadilla Benítez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    119128,
    'Gustavo Rubén Caballero González',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6236,
    'Adrián Andrés Cubas',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    305832,
    'Matías Galarza Fonda',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    278370,
    'Diego Alexander Gómez Amarilla',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    106485,
    'Maurício Magalhães Prado',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    70767,
    'Braian Óscar Ojeda Rodríguez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    196298,
    'Ramón Sosa Acosta',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2507,
    'Miguel Ángel Almirón Rejala',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2514,
    'Alejandro Sebastián Romero Gamarra',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    95460,
    'Alex Adrián Arce Barrios',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6483,
    'Gabriel Ávalos Stumpfs',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    70747,
    'Julio César Enciso Espínola',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    70670,
    'Isidro Miguel Pitta Saldívar',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2522,
    'Arnaldo Antonio Sanabria Ayala',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2380'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    16380,
    'Hernán Ismael Galíndez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    81224,
    'Wellington Moisés Ramírez Preciado',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    16642,
    'Gonzalo Roberto Valle Bustamante',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    46731,
    'Pervis Josué Estupiñán Tenorio',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    127817,
    'Piero Martín Hincapié Reyna',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    306940,
    'Yaimar Abel Medina Ortiz',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    354027,
    'Joel Leandro Ordóñez Guerrero',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    16367,
    'Willian Joel Pacho Tenorio',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2575,
    'Jackson Gabriel Porozo Vernaza',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2583,
    'Angelo Smit Preciado Quiñónez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    63964,
    'Félix Eduardo Torres Caicedo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    16470,
    'Jordy José Alcívar Macías',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    116117,
    'Moisés Isaac Caicedo Corozo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    338045,
    'Denil Daniel Castillo Preciado',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    16360,
    'Alan Steven Franco Palma',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    406303,
    'Ray Kendry Páez Andrade',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    280695,
    'Alan Steve Minda García',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    198347,
    'Anthony Lenín Valencia Bajaña',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    237078,
    'Pedro Jeampierre Vite Uca',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    361966,
    'Kevin José Rodríguez Cortez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    25414,
    'John Yeboah Zamora',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    311543,
    'Nilson David Angulo Ramírez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    16590,
    'Jordy Josué Caicedo Medina',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    350799,
    'Jeremy Alberto Arévalo Mera',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    16369,
    'Gonzalo Jordy Plata Jiménez',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    35533,
    'Enner Remberto Valencia Lastra',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2382'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    266606,
    'Christopher Keith Brady',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50728,
    'Matthew Andrew Geary Freese',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50999,
    'Matthew Charles Turner',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    38735,
    'Sergiño Gianni Dest',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    355994,
    'Alexander Michael Freeman',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50735,
    'Mark Irwin Robert Alexander McKenzie',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19023,
    'Timothy Michael Ream',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    126949,
    'Christopher Jeffrey Richards',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19549,
    'Antonee Robinson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50879,
    'Miles Gordon Robinson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50852,
    'Joseph Michael Scally',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50737,
    'Auston Levi-Jesaiah Trusty',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    362400,
    'Maximilian Arfsten',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    201713,
    'Sebastian Matthew Berhalter',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    415,
    'Weston James Earl McKennie',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    161921,
    'Giovanni Alejandro Reyna',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51114,
    'Cristian Roldan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162037,
    'Malik Leon Tillman',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1150,
    'Tyler Shaan Adams',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1138,
    'Timothy Tarpeh Weah',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    35885,
    'Alejandro Zendejas Saavedra',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50739,
    'Brenden Russell Aaronson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    138835,
    'Folarin Jerry Balogun',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    73868,
    'Ricardo Daniel Pepi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    17,
    'Christian Mate Pulišić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    427,
    'Haji Amir Wright',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2384'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    123742,
    'Josué Duverger',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    174768,
    'Alexandre Pierre',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    87789,
    'Johny Placide',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20850,
    'Carlens Jean Fedlaire Ruby Arcus',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    573613,
    'K. Thermoncy',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    12303,
    'Ricardo Ade',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    1411,
    'Hannes Delcroix',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20655,
    'Jean-Kévin Duverne',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    190747,
    'Martin Yves Expérience',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    102505,
    'Markhus Lacroix',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    275367,
    'Wilguens Paugain',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20665,
    'Jean-Ricner Bellegarde',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    338367,
    'Danley Jean-Jacques',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    20538,
    'Léverton Pierre',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    371050,
    'Olivier Woodensky Pierre',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    540857,
    'C. F. Sainte',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    194242,
    'Dominique Celidor Simon',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    174915,
    'Josué Casimir',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    128766,
    'Louicius Don Deedson',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50958,
    'Derrick Burckley Etienne Jr.',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    48535,
    'Yassin Enzo Fortuné',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    84087,
    'Wilson Isidor',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    21613,
    'Lenny Joseph',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    45020,
    'Duckens Moses Nazon',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    8601,
    'Frantzdy Pierrot',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162733,
    'Ruben Fritzner Providence',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '2386'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18110,
    'Maxime Teremoana Crocombe',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    94360,
    'Alexander Noah Paulsen',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    36777,
    'Michael Cornelis Woud',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    430835,
    'Tyler Bindon',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51149,
    'Michael Joseph Boxall',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6931,
    'Liberato Gianpaolo Cacace',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    94405,
    'Francis de Vries',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6932,
    'Callan Rennie Elliot',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    94346,
    'Timothy John Payne',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    94344,
    'Nando Zen Pijnaker',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51307,
    'Thomas Jefferson Smith',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    210165,
    'Finn Surman',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    405957,
    'Lachlan Bayliss',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    94541,
    'Joe Zen Robert Bell',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    180455,
    'Matthew Jimmy David Garbett',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    94322,
    'Callum William McCowatt',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6934,
    'Alex Arthur Rufer',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6935,
    'Sarpreet Singh',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    179862,
    'Marko Seufatu Nikola Stamenić',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    242,
    'Ryan Jared Thomas',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    94333,
    'Elijah Henry Just',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    179856,
    'Benjamin Craig Old',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6865,
    'Konstantinos Barbarouses',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    158688,
    'Jesse Randall',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    6938,
    'Benjamin Peter Waine',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18931,
    'Christopher Grant Wood',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '4673'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51274,
    'Maxime Crépeau',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    284554,
    'Owen Olamidayo Goodman',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51148,
    'Dayne Tristan St. Clair',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51295,
    'Derek Austin Cornelius',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    327738,
    'Luc Rollet De Fougerolles',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50816,
    'Richmond Mamah Laryea',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    78494,
    'Joel Robert Waterman',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    78547,
    'Alistair William Johnston',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18938,
    'Alfie Jones',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    407017,
    'Moïse Bombito Lumpungu',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    509,
    'Alphonso Boyle Davies',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    416901,
    'Niko Kristian Sigur',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50826,
    'Jacob Everett Shaffelburg',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    44798,
    'Liam Alan Millar',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50788,
    'Mathieu Choinière',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    35570,
    'Stephen Antunes Eustáquio',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    328046,
    'Ismaël Kenneth Jordan Koné',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    50817,
    'Jonathan Osorio',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    294824,
    'Nathan-Dylan Saliba',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    284061,
    'Marcelo Flores Dorrell',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    51016,
    'Tajon Trevor Buchanan',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    8489,
    'Jonathan Christian David',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    313353,
    'Promise Oluwatobi Emmanuel Akinpelu',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    2001,
    'Cyle Christopher Larin',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    362145,
    'Ali Ahmed',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    351587,
    'Tanitoluwa Oluwatimikhin Oluwaseyi',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5529'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    214129,
    'Tyrick Bodak',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    36967,
    'Trevor Irving Doornbusch',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    221,
    'Eloy Victor Room',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    37175,
    'Riechedly Bazoer',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    706,
    'Joshua Benjamin Brenet',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    36970,
    'Sherel Constancio Floranus',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    353808,
    'Deveron Fonville',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    36857,
    'Juriën Godfried Juan Gaari',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    194645,
    'Shurandy Ruggerio Sambo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    37066,
    'Roshon van Eijma',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    228,
    'Armando Obispo',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18997,
    'Leandro Jones Johan Bacuna',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    162199,
    'Ar''jany Martha',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    406846,
    'Tyrese Noslin',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    309763,
    'Livano Comenencia',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    37461,
    'Kevin Antonio Felida',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    906,
    'Tahith Jose Chong',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    36842,
    'Godfried Roemeratoe',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    19047,
    'Juninho Gracielo Bacuna',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    163220,
    'Jeremy Cornelis Jacobus Antonisse',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    37272,
    'Brandley Mack-Olien Kuwas',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    41627,
    'Kenji Joel Gorré',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    161884,
    'Sontje Hansen',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    38708,
    'Gervane Zjandric Adonnis Kastaneer',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    18981,
    'Jürgen Leonardo Locadia',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    195067,
    'Jearl Erwin Margaritha',
    (SELECT id FROM tournaments WHERE slug = 'world-cup-2026'),
    '5530'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM players WHERE tournament_id = (SELECT id FROM tournaments WHERE slug = 'world-cup-2026');

-- +goose StatementEnd
