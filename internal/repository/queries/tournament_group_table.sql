-- name: ListGroupTableByTournament :many
SELECT tgt.*, t.name AS team_name
FROM tournament_group_table tgt
JOIN teams t ON t.id = tgt.team_id
WHERE tgt.tournament_id = @tournament_id
ORDER BY tgt.group_letter ASC, tgt.position ASC;
