-- name: UpsertScorePrediction :one
INSERT INTO score_predictions (user_id, fixture_id, goals_home, goals_away)
VALUES (@user_id, @fixture_id, @goals_home, @goals_away)
ON CONFLICT (user_id, fixture_id) DO UPDATE
    SET goals_home = EXCLUDED.goals_home,
        goals_away = EXCLUDED.goals_away,
        updated_at = now()
RETURNING id, user_id, fixture_id, goals_home, goals_away, points, created_at, updated_at;

-- name: ListPredictionsByUserAndTournament :many
SELECT sp.id, sp.user_id, sp.fixture_id, sp.goals_home, sp.goals_away, sp.points, sp.created_at, sp.updated_at, sp.scored_at
FROM score_predictions sp
WHERE sp.user_id = @user_id
  AND sp.fixture_id IN (SELECT id FROM fixtures WHERE tournament_id = @tournament_id);
