-- name: GetTeamHandicap :one
SELECT id, team_id, category, points
FROM team_handicap
WHERE team_id  = @team_id
  AND category = @category;
