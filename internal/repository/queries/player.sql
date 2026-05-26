-- name: GetPlayerByID :one
SELECT p.id, p.external_id, p.name, p.tournament_id, p.team_id,
       t.group_letter,
       p.created_at, p.updated_at
FROM players p
JOIN teams t ON t.id = p.team_id
WHERE p.id = @id;

-- name: SearchPlayers :many
-- @escaped_query: wildcard-escaped term for ILIKE filtering; empty string skips name filter
-- @raw_query:     original trimmed term for similarity ranking (must not be escaped)
-- @group_letter:  optional group filter; empty string means search all groups
-- @has_handicap:  when true, only return players with at least one handicap row
WITH top_players AS (
    SELECT p.id, p.name, t.name AS team_name, t.logo AS team_logo, t.group_letter, p.tournament_id
    FROM players p
    JOIN teams t ON t.id = p.team_id
    WHERE p.tournament_id = @tournament_id
      AND (@group_letter = '' OR t.group_letter = @group_letter)
      AND (@escaped_query = '' OR unaccent_immutable(p.name) ILIKE '%' || unaccent_immutable(@escaped_query) || '%' ESCAPE '\')
    ORDER BY CASE WHEN @escaped_query = '' THEN p.name ELSE NULL END,
             CASE WHEN @escaped_query <> '' THEN similarity(unaccent_immutable(p.name), unaccent_immutable(@raw_query)) ELSE NULL END DESC
)
SELECT tp.id, tp.name, tp.team_name, tp.team_logo, tp.group_letter, tp.tournament_id,
       coalesce(
           json_object_agg(ph.category, ph.points) FILTER (WHERE ph.category IS NOT NULL),
           '{}'::json
       )::text AS handicaps
FROM top_players tp
LEFT JOIN player_handicap ph ON ph.player_id = tp.id
GROUP BY tp.id, tp.name, tp.team_name, tp.team_logo, tp.group_letter, tp.tournament_id
HAVING NOT @has_handicap OR count(ph.player_id) > 0
ORDER BY CASE WHEN @escaped_query = '' THEN tp.name ELSE NULL END,
         CASE WHEN @escaped_query <> '' THEN similarity(unaccent_immutable(tp.name), unaccent_immutable(@raw_query)) ELSE NULL END DESC;
