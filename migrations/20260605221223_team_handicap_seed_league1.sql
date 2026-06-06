-- +goose Up
-- +goose StatementBegin

-- Team handicap seed template.
--
-- For each team, replace NULL with the integer point value for that category.
-- Leave NULL to skip that team/category pair — it will NOT be inserted.
-- Idempotent: ON CONFLICT updates the existing row's points.
--
-- Categories:
--   group_winner   — points for predicting the team wins its group
--   playoff        — points for predicting the team reaches the playoff round
--   semifinalist   — points for predicting the team reaches the semi-finals
--   winner         — points for predicting the team wins the tournament

-- Belgium (BEL)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Belgium'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- France (FRA)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'France'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Croatia (CRO)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Croatia'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Sweden (SWE)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Sweden'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Brazil (BRA)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Brazil'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Uruguay (URU)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Uruguay'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Colombia (COL)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Colombia'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Spain (ESP)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Spain'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- England (ENG)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'England'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Panama (PAN)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Panama'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Japan (JPN)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Japan'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Senegal (SEN)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Senegal'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Switzerland (SUI)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Switzerland'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Mexico (MEX)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Mexico'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- South Korea (KOR)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'South Korea'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Australia (AUS)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Australia'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Iran (IRN)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Iran'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Saudi Arabia (KSA)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Saudi Arabia'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Germany (GER)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Germany'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Argentina (ARG)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Argentina'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Portugal (POR)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Portugal'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Tunisia (TUN)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Tunisia'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Morocco (MAR)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Morocco'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Egypt (EGY)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Egypt'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Czech Republic (CZE)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Czech Republic'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Austria (AUT)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Austria'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Türkiye (TUR)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Türkiye'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Norway (NOR)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Norway'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Scotland (SCO)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Scotland'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Bosnia & Herzegovina (BIH)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Bosnia & Herzegovina'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Netherlands (NED)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Netherlands'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Ivory Coast (CIV)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Ivory Coast'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Ghana (GHA)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Ghana'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Congo DR (CGO)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Congo DR'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- South Africa (RSA)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'South Africa'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Algeria (ALG)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Algeria'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Cape Verde Islands (CPV)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Cape Verde Islands'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Jordan (JOR)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Jordan'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Iraq (IRQ)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Iraq'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Uzbekistan (UZB)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Uzbekistan'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Qatar (QAT)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Qatar'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Paraguay (PAR)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Paraguay'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Ecuador (ECU)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Ecuador'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- USA (USA)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'USA'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Haiti (HAI)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Haiti'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- New Zealand (NZL)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'New Zealand'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Canada (CAN)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Canada'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- Curaçao (CUR)
INSERT INTO team_handicap (team_id, category, points)
SELECT t.id, v.category, v.points
FROM   teams t
JOIN   (VALUES
    ('group_winner'::team_handicap_category,  NULL::int),
    ('playoff'::team_handicap_category,       NULL::int),
    ('semifinalist'::team_handicap_category,  NULL::int),
    ('winner'::team_handicap_category,        NULL::int)
) AS v(category, points) ON TRUE
WHERE  t.name = 'Curaçao'
  AND  v.points IS NOT NULL
ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- No automatic rollback for handicap data — delete manually if needed.

-- +goose StatementEnd
