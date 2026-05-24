-- name: GetTeamByID :one
SELECT id, name, logo, tournament_id
FROM teams
WHERE id = @id;
