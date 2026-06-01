-- name: ListGroupTableByTournament :many
SELECT tgt.id, tgt.tournament_id, tgt.team_id, t.name AS team_name,
       tgt.group_letter, tgt.position, tgt.points, tgt.played, tgt.created_at, tgt.updated_at
FROM tournament_group_table tgt
JOIN teams t ON t.id = tgt.team_id
WHERE tgt.tournament_id = @tournament_id
ORDER BY tgt.group_letter ASC, tgt.position ASC;
