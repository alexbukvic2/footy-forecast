-- name: SearchPlayers :many
-- @escaped_query: wildcard-escaped term for ILIKE filtering
-- @raw_query:     original trimmed term for similarity ranking (must not be escaped)
SELECT p.id, p.name, t.name AS team_name, t.logo AS team_logo, p.tournament_id
FROM players p
JOIN teams t ON t.id = p.team_id
WHERE p.tournament_id = @tournament_id
  AND unaccent_immutable(p.name) ILIKE '%' || unaccent_immutable(@escaped_query) || '%' ESCAPE '\'
ORDER BY similarity(unaccent_immutable(p.name), unaccent_immutable(@raw_query)) DESC
LIMIT 5;
