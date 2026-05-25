-- name: GetTeamByID :one
SELECT id, name, logo, tournament_id, group_letter
FROM teams
WHERE id = @id;

-- name: ListGroupLettersByTournament :many
SELECT DISTINCT group_letter
FROM teams
WHERE tournament_id = @tournament_id
  AND group_letter IS NOT NULL
ORDER BY group_letter ASC;
