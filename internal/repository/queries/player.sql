-- name: GetPlayerByID :one
SELECT p.id, p.external_id, p.name, p.tournament_id, p.team_id,
       t.group_letter,
       p.created_at, p.updated_at
FROM players p
JOIN teams t ON t.id = p.team_id
WHERE p.id = @id;

-- name: SearchPlayers :many
-- @escaped_query: wildcard-escaped term for ILIKE filtering
-- @raw_query:     original trimmed term for similarity ranking (must not be escaped)
WITH top_players AS (
    SELECT p.id, p.name, t.name AS team_name, t.logo AS team_logo, p.tournament_id
    FROM players p
    JOIN teams t ON t.id = p.team_id
    WHERE p.tournament_id = @tournament_id
      AND unaccent_immutable(p.name) ILIKE '%' || unaccent_immutable(@escaped_query) || '%' ESCAPE '\'
    ORDER BY similarity(unaccent_immutable(p.name), unaccent_immutable(@raw_query)) DESC
    LIMIT 5
)
SELECT tp.id, tp.name, tp.team_name, tp.team_logo, tp.tournament_id,
       coalesce(
           json_object_agg(ph.category, ph.points) FILTER (WHERE ph.category IS NOT NULL),
           '{}'::json
       )::text AS handicaps
FROM top_players tp
LEFT JOIN player_handicap ph ON ph.player_id = tp.id
GROUP BY tp.id, tp.name, tp.team_name, tp.team_logo, tp.tournament_id
ORDER BY similarity(unaccent_immutable(tp.name), unaccent_immutable(@raw_query)) DESC;
