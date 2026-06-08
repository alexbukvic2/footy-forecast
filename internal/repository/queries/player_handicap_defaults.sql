-- name: GetPlayerHandicapDefaults :many
SELECT category, default_points
FROM player_handicap_defaults
WHERE tournament_id = @tournament_id;
