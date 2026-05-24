-- name: GetFirstKickoffByTournament :one
SELECT kickoff_at
FROM fixtures
WHERE tournament_id = @tournament_id
ORDER BY kickoff_at ASC
LIMIT 1;
