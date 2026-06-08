-- +goose Up
-- +goose StatementBegin

-- Seeded group_top_scorer handicap points for World Cup 2026.
-- Players matched from API-Football external_id; names in comments are the display names from the input.
-- Idempotent: ON CONFLICT updates existing rows.
--
-- Skipped (not found in players table for this tournament):
--   Hirving Lozano (Mexico)   — not in squad
--   Assan Ouedraogo (Germany) — not in squad
--   Takumi Minamino (Japan)   — not in squad

INSERT INTO player_handicap (player_id, category, points)
SELECT p.id, 'group_top_scorer'::player_handicap_category, v.points
FROM players p
JOIN tournaments tr ON p.tournament_id = tr.id
JOIN (VALUES
    -- Group A
    (35532,  4),  -- Julian Quinones (Mexico)
    (  794,  4),  -- Patrik Schick (Czech Republic)
    ( 2887,  4),  -- Raul Jimenez (Mexico)
    (94562,  4),  -- Santiago Gimenez (Mexico)
    (  186,  5),  -- Heung-Min Son (South Korea)
    (  927,  5),  -- Kang-In Lee (South Korea)
    (35576,  5),  -- Orbelin Pineda (Mexico)
    (24888,  6),  -- Hee-Chan Hwang (South Korea)
    (291713, 6),  -- Armando Gonzalez (Mexico)
    (98936,  8),  -- Lyle Foster (South Africa)
    -- Group B
    (  790,  4),  -- Edin Dzeko (Bosnia & Herzegovina)
    ( 8489,  5),  -- Jonathan David (Canada)
    (313353, 5),  -- Promise David (Canada) [db name: Promise Oluwatobi Emmanuel Akinpelu]
    (  421,  6),  -- Breel Embolo (Switzerland)
    ( 2544,  6),  -- Akram Afif (Qatar)
    ( 2001,  6),  -- Cyle Larin (Canada)
    (48648,  6),  -- Dan Ndoye (Switzerland)
    ( 1464,  6),  -- Granit Xhaka (Switzerland)
    (  509,  7),  -- Alphonso Davies (Canada)
    (395559, 8),  -- Kerim Alajbegovic (Bosnia & Herzegovina)
    -- Group C
    (  762,  4),  -- Vinicius Jr. (Brazil) [db name: Vinícius José Paixão de Oliveira Júnior]
    ( 1496,  4),  -- Raphinha (Brazil) [db name: Raphael Dias Belloli]
    (196156, 4),  -- Igor Thiago (Brazil)
    (  276,  4),  -- Neymar Jr (Brazil) [db name: Neymar da Silva Santos Júnior]
    ( 1165,  5),  -- Matheus Cunha (Brazil)
    (377122, 5),  -- Endrick (Brazil)
    (127769, 5),  -- Gabriel Martinelli (Brazil)
    (161897, 5),  -- Ismael Saibari (Morocco)
    (407806, 5),  -- Rayan Vitor (Brazil)
    ( 2722,  6),  -- Ayoub El Kaabi (Morocco)
    (  744,  6),  -- Brahim Diaz (Morocco)
    (  747,  6),  -- Casemiro (Brazil) [db name: Carlos Henrique Casimiro]
    (  903,  6),  -- Scott McTominay (Scotland)
    (45175,  6),  -- Lawrence Shankland (Scotland)
    (19524,  6),  -- Che Adams (Scotland)
    (45307,  7),  -- Lyndon Dykes (Scotland)
    ( 1125,  7),  -- Ryan Christie (Scotland)
    (    9,  7),  -- Achraf Hakimi (Morocco)
    (19191,  7),  -- John McGinn (Scotland)
    (84087,  8),  -- Wilson Isidor (Haiti)
    -- Group D
    (   17,  4),  -- Christian Pulisic (USA)
    (138835, 4),  -- Folarin Balogun (USA)
    (  427,  4),  -- Haji Wright (USA)
    (73868,  4),  -- Ricardo Pepi (USA)
    (291964, 5),  -- Arda Guler (Türkiye)
    (339887, 5),  -- Can Uzun (Türkiye)
    (339883, 5),  -- Kenan Yildiz (Türkiye)
    (142959, 5),  -- Kerem Akturkoglu (Türkiye)
    (161921, 6),  -- Giovanni Reyna (USA)
    ( 1640,  6),  -- Hakan Calhanoglu (Türkiye)
    (198352, 7),  -- Mohamed Toure (Australia)
    (70747,  7),  -- Julio Enciso (Paraguay)
    ( 2507,  7),  -- Miguel Almiron (Paraguay)
    (50739,  8),  -- Brenden Aaronson (USA)
    (338014, 8),  -- Nestory Irankunda (Australia)
    (296645, 8),  -- Tete Yengi (Australia)
    -- Group E
    (  978,  4),  -- Kai Havertz (Germany)
    (158054, 4),  -- Nick Woltemade (Germany)
    (203224, 4),  -- Florian Wirtz (Germany)
    (181812, 4),  -- Jamal Musiala (Germany)
    (  644,  5),  -- Leroy Sane (Germany)
    (26475,  5),  -- Deniz Undav (Germany)
    (35533,  5),  -- Enner Valencia (Ecuador)
    (  502,  6),  -- Joshua Kimmich (Germany)
    (162707, 6),  -- Elye Wahi (Ivory Coast)
    (16590,  6),  -- Jordy Caicedo (Ecuador)
    (157997, 7),  -- Amad Diallo (Ivory Coast)
    (275651, 7),  -- Ange Bonny (Ivory Coast)
    (16369,  7),  -- Gonzalo Plata (Ecuador)
    (513776, 7),  -- Yan Diomande (Ivory Coast)
    (137303, 7),  -- Evann Guessand (Ivory Coast)
    (311543, 7),  -- Nilson Angulo (Ecuador)
    (  906,  8),  -- Tahith Chong (Curaçao)
    -- Group F
    (  247,  4),  -- Cody Gakpo (Netherlands)
    (  249,  4),  -- Donyell Malen (Netherlands)
    (  667,  4),  -- Memphis Depay (Netherlands)
    (18979,  4),  -- Viktor Gyokeres (Sweden)
    ( 2864,  5),  -- Alexander Isak (Sweden)
    (72155,  5),  -- Ayase Ueda (Japan)
    (  544,  6),  -- Noa Lang (Netherlands)
    (37724,  6),  -- Crysencio Summerville (Netherlands)
    (33224,  6),  -- Daizen Maeda (Japan)
    (48002,  7),  -- Benjamin Nygren (Sweden)
    (  226,  7),  -- Denzel Dumfries (Netherlands)
    ( 2601,  8),  -- Daichi Kamada (Japan)
    (  290,  8),  -- Virgil van Dijk (Netherlands)
    -- Group G
    (  907,  4),  -- Romelu Lukaku (Belgium)
    (  629,  4),  -- Kevin De Bruyne (Belgium)
    (  306,  4),  -- Mohamed Salah (Egypt)
    (147859, 5),  -- Charles De Ketelaere (Belgium)
    ( 1422,  5),  -- Jeremy Doku (Belgium)
    ( 1946,  5),  -- Leandro Trossard (Belgium)
    (81573,  5),  -- Omar Marmoush (Egypt)
    (25458,  6),  -- Dodi Lukebakio (Belgium)
    (42315,  6),  -- Mehdi Taremi (Iran)
    (18931,  6),  -- Chris Wood (New Zealand)
    ( 6938,  8),  -- Ben Waine (New Zealand)
    -- Group H
    (47323,  4),  -- Mikel Oyarzabal (Spain)
    (386828, 4),  -- Lamine Yamal (Spain)
    (  931,  5),  -- Ferran Torres (Spain)
    ( 1323,  6),  -- Dani Olmo (Spain)
    (183799, 6),  -- Nico Williams (Spain)
    (51617,  6),  -- Darwin Nunez (Uruguay)
    (47311,  6),  -- Mikel Merino (Spain)
    (133609, 7),  -- Pedri (Spain) [db name: Pedro González López]
    (51618,  7),  -- Brian Rodriguez (Uruguay)
    (  328,  7),  -- Fabian Ruiz (Spain)
    (70078,  7),  -- Facundo Pellistri (Uruguay)
    (  756,  7),  -- Federico Valverde (Uruguay)
    (51530,  7),  -- Federico Vinas (Uruguay)
    (44551,  8),  -- Saleh Al Shehri (Saudi Arabia)
    (44340,  8),  -- Salem Al Dawsari (Saudi Arabia)
    ( 2612,  8),  -- Giorgian De Arrascaeta (Uruguay)
    -- Group I
    (  278,  4),  -- Kylian Mbappe (France)
    ( 1100,  5),  -- Erling Haaland (Norway)
    (  153,  5),  -- Ousmane Dembele (France)
    (19617,  6),  -- Michael Olise (France)
    (343027, 6),  -- Desire Doue (France)
    (25927,  7),  -- Jean-Philippe Mateta (France)
    (156477, 7),  -- Rayan Cherki (France)
    ( 8492,  7),  -- Alexander Sorloth (Norway)
    (161904, 7),  -- Bradley Barcola (France)
    (  304,  7),  -- Sadio Mane (Senegal)
    (283058, 7),  -- Nicolas Jackson (Senegal)
    ( 2032,  7),  -- Jorgen Larsen (Norway)
    (21509,  7),  -- Marcus Thuram (France)
    (278133, 7),  -- Oscar Bobb (Norway)
    ( 2218,  8),  -- Ismaila Sarr (Senegal)
    (37127,  8),  -- Martin Odegaard (Norway)
    (314511, 8),  -- Antonio Nusa (Norway)
    -- Group J
    (  154,  4),  -- Lionel Messi (Argentina)
    ( 6009,  5),  -- Julian Alvarez (Argentina)
    (  217,  5),  -- Lautaro Martinez (Argentina)
    (350037, 5),  -- Nicolas Paz (Argentina)
    ( 6067,  6),  -- Thiago Almada (Argentina)
    ( 5996,  7),  -- Enzo Fernandez (Argentina)
    (18830,  7),  -- Marko Arnautovic (Austria)
    (  635,  7),  -- Riyad Mahrez (Algeria)
    ( 1159,  7),  -- Marcel Sabitzer (Austria)
    (25297,  8),  -- Michael Gregoritsch (Austria)
    (200139, 8),  -- Mohamed Amoura (Algeria)
    -- Group K
    (  874,  4),  -- Cristiano Ronaldo (Portugal)
    (47237,  5),  -- Luis Suarez (Colombia)
    ( 1485,  5),  -- Bruno Fernandes (Portugal)
    (41585,  5),  -- Goncalo Ramos (Portugal)
    ( 2489,  5),  -- Luis Diaz (Colombia)
    ( 1864,  5),  -- Pedro Neto (Portugal)
    (22236,  6),  -- Rafael Leao (Portugal)
    (  517,  6),  -- James Rodriguez (Colombia)
    (  583,  6),  -- Joao Felix (Portugal)
    (47582,  7),  -- Cucho Hernandez (Colombia)
    (13708,  7),  -- Jhon Arias (Colombia)
    (53535,  8),  -- Eldor Shomurodov (Uzbekistan)
    -- Group L
    (  184,  4),  -- Harry Kane (England)
    ( 1460,  6),  -- Bukayo Saka (England)
    (19366,  6),  -- Ollie Watkins (England)
    (129718, 6),  -- Jude Bellingham (England)
    (  909,  6),  -- Marcus Rashford (England)
    (19170,  6),  -- Morgan Rogers (England)
    (19586,  6),  -- Eberechi Eze (England)
    (138787, 7),  -- Anthony Gordon (England)
    (19281,  7),  -- Antoine Semenyo (Ghana)
    (19974,  7),  -- Ivan Toney (England)
    (  726,  7),  -- Andrej Kramaric (Croatia)
    (307123, 8),  -- Nico O'Reilly (England)
    (  842,  8),  -- Nikola Vlasic (Croatia)
    ( 2937,  8)   -- Declan Rice (England)
) AS v(external_id, points) ON p.external_id = v.external_id
WHERE tr.slug = 'world-cup-2026'
ON CONFLICT (player_id, category) DO UPDATE SET points = EXCLUDED.points;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM player_handicap
WHERE category = 'group_top_scorer'
  AND player_id IN (
    SELECT p.id
    FROM players p
    JOIN tournaments tr ON p.tournament_id = tr.id
    WHERE tr.slug = 'world-cup-2026'
      AND p.external_id IN (
        35532, 794, 2887, 94562, 186, 927, 35576, 24888, 291713, 98936,
        790, 8489, 313353, 421, 2544, 2001, 48648, 1464, 509, 395559,
        762, 1496, 196156, 276, 1165, 377122, 127769, 161897, 407806,
        2722, 744, 747, 903, 45175, 19524, 45307, 1125, 9, 19191, 84087,
        17, 138835, 427, 73868, 291964, 339887, 339883, 142959, 161921,
        1640, 198352, 70747, 2507, 50739, 338014, 296645,
        978, 158054, 203224, 181812, 644, 26475, 35533, 502, 162707,
        16590, 157997, 275651, 16369, 513776, 137303, 311543, 906,
        247, 249, 667, 18979, 2864, 72155, 544, 37724, 33224, 48002,
        226, 2601, 290,
        907, 629, 306, 147859, 1422, 1946, 81573, 25458, 42315, 18931, 6938,
        47323, 386828, 931, 1323, 183799, 51617, 47311, 133609, 51618,
        328, 70078, 756, 51530, 44551, 44340, 2612,
        278, 1100, 153, 19617, 343027, 25927, 156477, 8492, 161904,
        304, 283058, 2032, 21509, 278133, 2218, 37127, 314511,
        154, 6009, 217, 350037, 6067, 5996, 18830, 635, 1159, 25297, 200139,
        874, 47237, 1485, 41585, 2489, 1864, 22236, 517, 583, 47582, 13708, 53535,
        184, 1460, 19366, 129718, 909, 19170, 19586, 138787, 19281,
        19974, 726, 307123, 842, 2937
      )
  );

-- +goose StatementEnd
