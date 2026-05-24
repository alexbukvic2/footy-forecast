-- name: UpsertPlayerPrediction :one
WITH ins AS (
    INSERT INTO player_predictions (user_id, tournament_id, category, pick)
    VALUES (@user_id, @tournament_id, @category, @pick)
    ON CONFLICT (user_id, tournament_id, category) DO UPDATE
        SET pick       = EXCLUDED.pick,
            updated_at = now()
    RETURNING id, user_id, tournament_id, category, pick, points, created_at, updated_at
)
SELECT ins.id, ins.user_id, ins.tournament_id, ins.category, ins.pick,
       p.name AS pick_name, ins.points, ins.created_at, ins.updated_at
FROM ins
JOIN players p ON p.id = ins.pick;

-- name: ListPlayerPredictionsByTournamentForUser :many
SELECT pp.id, pp.user_id, pp.tournament_id, pp.category, pp.pick,
       p.name AS pick_name, pp.points, pp.created_at, pp.updated_at
FROM player_predictions pp
JOIN players p ON p.id = pp.pick
WHERE pp.tournament_id = @tournament_id AND pp.user_id = @user_id;

-- name: ListPlayerPredictionsByLeague :many
SELECT lm.user_id, pp.category, pp.pick AS player_id,
       p.name AS player_name, pp.points
FROM league_members lm
JOIN leagues l ON l.id = lm.league_id
JOIN player_predictions pp ON pp.user_id = lm.user_id AND pp.tournament_id = l.tournament_id
JOIN players p ON p.id = pp.pick
WHERE lm.league_id = @league_id;

-- name: UpsertTeamPrediction :one
WITH ins AS (
    INSERT INTO team_predictions (user_id, tournament_id, category, pick)
    VALUES (@user_id, @tournament_id, @category, @pick)
    ON CONFLICT (user_id, tournament_id, category) DO UPDATE
        SET pick       = EXCLUDED.pick,
            updated_at = now()
    RETURNING id, user_id, tournament_id, category, pick, points, created_at, updated_at
)
SELECT ins.id, ins.user_id, ins.tournament_id, ins.category, ins.pick,
       t.name AS pick_name, ins.points, ins.created_at, ins.updated_at
FROM ins
JOIN teams t ON t.id = ins.pick;

-- name: ListTeamPredictionsByTournamentForUser :many
SELECT tp.id, tp.user_id, tp.tournament_id, tp.category, tp.pick,
       t.name AS pick_name, tp.points, tp.created_at, tp.updated_at
FROM team_predictions tp
JOIN teams t ON t.id = tp.pick
WHERE tp.tournament_id = @tournament_id AND tp.user_id = @user_id;

-- name: ListTeamPredictionsByLeague :many
SELECT lm.user_id, tp.category, tp.pick AS team_id,
       t.name AS team_name, tp.points
FROM league_members lm
JOIN leagues l ON l.id = lm.league_id
JOIN team_predictions tp ON tp.user_id = lm.user_id AND tp.tournament_id = l.tournament_id
JOIN teams t ON t.id = tp.pick
WHERE lm.league_id = @league_id;

-- name: ListLeagueMembersForPredictions :many
SELECT lm.user_id, u.display_name
FROM league_members lm
JOIN users u ON u.id = lm.user_id
WHERE lm.league_id = @league_id
ORDER BY u.display_name ASC;
