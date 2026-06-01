-- name: ListPlayerOutcomesByTournament :many
SELECT po.id, po.tournament_id, po.category, po.recorded_at,
       p.id AS player_id, p.name AS player_name,
       t.id AS team_id, t.name AS team_name
FROM player_outcomes po
JOIN players p ON p.id = po.player_id
JOIN teams   t ON t.id = p.team_id
WHERE po.tournament_id = @tournament_id
ORDER BY po.category ASC, p.name ASC;

-- name: ListTeamOutcomesByTournament :many
SELECT to_.id, to_.tournament_id, to_.category, to_.recorded_at,
       t.id AS team_id, t.name AS team_name
FROM team_outcomes to_
JOIN teams t ON t.id = to_.team_id
WHERE to_.tournament_id = @tournament_id
ORDER BY to_.category ASC, t.name ASC;
