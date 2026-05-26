-- name: UpsertPlayerPredictionGroup :one
WITH ins AS (
    INSERT INTO player_predictions (user_id, tournament_id, category, pick, group_letter)
    VALUES (@user_id, @tournament_id, @category, @pick, @group_letter)
    ON CONFLICT (user_id, tournament_id, category, group_letter)
        WHERE group_letter IS NOT NULL
    DO UPDATE SET pick = EXCLUDED.pick, updated_at = now()
    RETURNING id, user_id, tournament_id, category, pick, group_letter, points, created_at, updated_at
)
SELECT ins.id, ins.user_id, ins.tournament_id, ins.category, ins.pick,
       ins.group_letter, p.name AS pick_name, ins.points, ins.created_at, ins.updated_at
FROM ins JOIN players p ON p.id = ins.pick;

-- name: UpsertPlayerPredictionNoGroup :one
WITH ins AS (
    INSERT INTO player_predictions (user_id, tournament_id, category, pick)
    VALUES (@user_id, @tournament_id, @category, @pick)
    ON CONFLICT (user_id, tournament_id, category)
        WHERE group_letter IS NULL
    DO UPDATE SET pick = EXCLUDED.pick, updated_at = now()
    RETURNING id, user_id, tournament_id, category, pick, group_letter, points, created_at, updated_at
)
SELECT ins.id, ins.user_id, ins.tournament_id, ins.category, ins.pick,
       ins.group_letter, p.name AS pick_name, ins.points, ins.created_at, ins.updated_at
FROM ins JOIN players p ON p.id = ins.pick;

-- name: ListPlayerPredictionsByTournamentForUser :many
SELECT pp.id, pp.user_id, pp.tournament_id, pp.category, pp.pick,
       pp.group_letter, p.name AS pick_name, pp.points, pp.created_at, pp.updated_at
FROM player_predictions pp
JOIN players p ON p.id = pp.pick
WHERE pp.tournament_id = @tournament_id AND pp.user_id = @user_id;

-- name: ListPlayerPredictionsByLeague :many
SELECT lm.user_id, pp.category, pp.group_letter, pp.pick AS player_id,
       p.name AS player_name, pp.points
FROM league_members lm
JOIN leagues l ON l.id = lm.league_id
JOIN player_predictions pp
    ON pp.user_id = lm.user_id AND pp.tournament_id = l.tournament_id
JOIN players p ON p.id = pp.pick
WHERE lm.league_id = @league_id;

-- name: UpsertTeamPredictionGroup :one
WITH ins AS (
    INSERT INTO team_predictions (user_id, tournament_id, category, pick, group_letter, slot_index)
    VALUES (@user_id, @tournament_id, @category, @pick, @group_letter, @slot_index)
    ON CONFLICT (user_id, tournament_id, category, group_letter, slot_index)
        WHERE group_letter IS NOT NULL
    DO UPDATE SET pick = EXCLUDED.pick, updated_at = now()
    RETURNING id, user_id, tournament_id, category, pick, group_letter, slot_index, points, created_at, updated_at
)
SELECT ins.id, ins.user_id, ins.tournament_id, ins.category, ins.pick,
       ins.group_letter, ins.slot_index, t.name AS pick_name, ins.points, ins.created_at, ins.updated_at
FROM ins JOIN teams t ON t.id = ins.pick;

-- name: UpsertTeamPredictionNoGroup :one
WITH ins AS (
    INSERT INTO team_predictions (user_id, tournament_id, category, pick, slot_index)
    VALUES (@user_id, @tournament_id, @category, @pick, @slot_index)
    ON CONFLICT (user_id, tournament_id, category, slot_index)
        WHERE group_letter IS NULL
    DO UPDATE SET pick = EXCLUDED.pick, updated_at = now()
    RETURNING id, user_id, tournament_id, category, pick, group_letter, slot_index, points, created_at, updated_at
)
SELECT ins.id, ins.user_id, ins.tournament_id, ins.category, ins.pick,
       ins.group_letter, ins.slot_index, t.name AS pick_name, ins.points, ins.created_at, ins.updated_at
FROM ins JOIN teams t ON t.id = ins.pick;

-- name: ListTeamPredictionsByTournamentForUser :many
SELECT tp.id, tp.user_id, tp.tournament_id, tp.category, tp.pick,
       tp.group_letter, tp.slot_index, t.name AS pick_name,
       tp.points, tp.created_at, tp.updated_at
FROM team_predictions tp
JOIN teams t ON t.id = tp.pick
WHERE tp.tournament_id = @tournament_id AND tp.user_id = @user_id;

-- name: ListTeamPredictionsByLeague :many
SELECT lm.user_id, tp.category, tp.group_letter, tp.slot_index,
       tp.pick AS team_id, t.name AS team_name, tp.points
FROM league_members lm
JOIN leagues l ON l.id = lm.league_id
JOIN team_predictions tp
    ON tp.user_id = lm.user_id AND tp.tournament_id = l.tournament_id
JOIN teams t ON t.id = tp.pick
WHERE lm.league_id = @league_id;

-- name: CountPlayoffWildcards :one
SELECT COUNT(*)::int FROM team_predictions
WHERE tournament_id = @tournament_id
  AND user_id       = @user_id
  AND category      = 'playoff'
  AND slot_index    = 2;

-- name: DeletePlayerPredictionGroup :exec
DELETE FROM player_predictions
WHERE user_id       = @user_id
  AND tournament_id = @tournament_id
  AND category      = @category
  AND group_letter  = @group_letter;

-- name: DeletePlayerPredictionNoGroup :exec
DELETE FROM player_predictions
WHERE user_id       = @user_id
  AND tournament_id = @tournament_id
  AND category      = @category
  AND group_letter  IS NULL;

-- name: DeleteTeamPredictionGroup :exec
DELETE FROM team_predictions
WHERE user_id       = @user_id
  AND tournament_id = @tournament_id
  AND category      = @category
  AND group_letter  = @group_letter
  AND slot_index    = @slot_index;

-- name: DeleteTeamPredictionNoGroup :exec
DELETE FROM team_predictions
WHERE user_id       = @user_id
  AND tournament_id = @tournament_id
  AND category      = @category
  AND group_letter  IS NULL
  AND slot_index    = @slot_index;

-- name: ListLeagueMembersForPredictions :many
SELECT lm.user_id, u.display_name
FROM league_members lm
JOIN users u ON u.id = lm.user_id
WHERE lm.league_id = @league_id
ORDER BY u.display_name ASC;
