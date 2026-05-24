-- name: GetFixtureByID :one
SELECT id, external_id, tournament_id, home_team_id, away_team_id,
       kickoff_at, status, goals_home, goals_away, created_at, updated_at
FROM fixtures
WHERE id = @id;

-- name: ListFixturesByTournament :many
SELECT id, external_id, tournament_id, home_team_id, away_team_id,
       kickoff_at, status, goals_home, goals_away, created_at, updated_at
FROM fixtures
WHERE tournament_id = @tournament_id
ORDER BY kickoff_at;

-- name: ListLockedFixturesByLeague :many
SELECT f.id, f.external_id, f.tournament_id, f.home_team_id, f.away_team_id,
       f.kickoff_at, f.status, f.goals_home, f.goals_away, f.created_at, f.updated_at,
       coalesce(
           json_agg(
               json_build_object(
                   'user_id',     lm.user_id,
                   'display_name', u.display_name,
                   'goals_home',  sp.goals_home,
                   'goals_away',  sp.goals_away,
                   'points',      sp.points
               ) ORDER BY (lm.user_id = @requesting_user_id) DESC, u.display_name ASC
           ),
           '[]'::json
       ) AS member_predictions
FROM fixtures f
JOIN leagues l ON l.tournament_id = f.tournament_id
JOIN league_members lm ON lm.league_id = l.id
JOIN users u ON u.id = lm.user_id
LEFT JOIN score_predictions sp
    ON sp.fixture_id = f.id AND sp.user_id = lm.user_id
WHERE l.id = @league_id
  AND f.status IN ('in_progress', 'finished')
GROUP BY f.id
ORDER BY f.kickoff_at;
