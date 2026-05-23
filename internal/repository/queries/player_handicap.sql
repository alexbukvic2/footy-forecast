-- name: GetPlayerHandicap :one
SELECT id, player_id, category, points
FROM player_handicap
WHERE player_id = @player_id
  AND category  = @category;
