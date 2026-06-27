-- name: GetFirstKickoffByTournament :one
SELECT kickoff_at
FROM fixtures
WHERE tournament_id = @tournament_id
ORDER BY kickoff_at ASC
LIMIT 1;

-- name: GetFixtureByID :one
SELECT f.id, f.external_id, f.tournament_id, f.home_team_id, f.away_team_id,
       home_t.name AS home_team_name, away_t.name AS away_team_name,
       CASE WHEN f.round LIKE '%Group%' THEN home_t.group_letter ELSE NULL END AS group_letter,
       f.round, f.kickoff_at, f.status, f.prediction_locked,
       f.goals_home_regular, f.goals_away_regular, f.goals_home, f.goals_away,
       f.winner_team_id, f.created_at, f.updated_at
FROM fixtures f
JOIN teams home_t ON home_t.id = f.home_team_id
JOIN teams away_t ON away_t.id = f.away_team_id
WHERE f.id = @id;

-- name: ListFixturesByTournament :many
SELECT f.id, f.external_id, f.tournament_id, f.home_team_id, f.away_team_id,
       home_t.name AS home_team_name, away_t.name AS away_team_name,
       CASE WHEN f.round LIKE '%Group%' THEN home_t.group_letter ELSE NULL END AS group_letter,
       f.round, f.kickoff_at, f.status, f.prediction_locked,
       f.goals_home_regular, f.goals_away_regular, f.goals_home, f.goals_away,
       f.winner_team_id, f.created_at, f.updated_at
FROM fixtures f
JOIN teams home_t ON home_t.id = f.home_team_id
JOIN teams away_t ON away_t.id = f.away_team_id
WHERE f.tournament_id = @tournament_id ORDER BY f.kickoff_at;

-- name: GetLockedFixtureDatesByLeague :many
SELECT DISTINCT f.kickoff_at::date AS fixture_date
FROM fixtures f
JOIN leagues l ON l.tournament_id = f.tournament_id
WHERE l.id = @league_id
  AND f.prediction_locked = true
  AND f.kickoff_at::date <= now()::date
ORDER BY fixture_date DESC;

-- name: ListLockedFixturesByLeague :many
SELECT f.id, f.external_id, f.tournament_id, f.home_team_id, f.away_team_id,
       home_t.name AS home_team_name, away_t.name AS away_team_name,
       CASE WHEN f.round LIKE '%Group%' THEN home_t.group_letter ELSE NULL END AS group_letter,
       f.round, f.kickoff_at, f.status, f.prediction_locked,
       f.goals_home_regular, f.goals_away_regular, f.goals_home, f.goals_away,
       f.winner_team_id, f.created_at, f.updated_at,
       coalesce(json_agg(json_build_object(
           'user_id', lm.user_id,
           'display_name', u.display_name,
           'goals_home', sp.goals_home,
           'goals_away', sp.goals_away,
           'winner', sp.winner,
           'points', sp.points
       ) ORDER BY (lm.user_id = @requesting_user_id) DESC, u.display_name ASC), '[]'::json)::text
       AS member_predictions
FROM fixtures f
JOIN teams home_t ON home_t.id = f.home_team_id
JOIN teams away_t ON away_t.id = f.away_team_id
JOIN leagues l ON l.tournament_id = f.tournament_id
JOIN league_members lm ON lm.league_id = l.id
JOIN users u ON u.id = lm.user_id
LEFT JOIN score_predictions sp ON sp.fixture_id = f.id AND sp.user_id = lm.user_id
WHERE l.id = @league_id AND f.status IN ('in_progress', 'finished', 'cancelled')
GROUP BY f.id, home_t.name, away_t.name, home_t.group_letter, f.round
ORDER BY f.kickoff_at DESC;

-- name: ListLockedFixturesByLeagueAndDates :many
SELECT f.id, f.external_id, f.tournament_id, f.home_team_id, f.away_team_id,
       home_t.name AS home_team_name, away_t.name AS away_team_name,
       CASE WHEN f.round LIKE '%Group%' THEN home_t.group_letter ELSE NULL END AS group_letter,
       f.round, f.kickoff_at, f.status, f.prediction_locked,
       f.goals_home_regular, f.goals_away_regular, f.goals_home, f.goals_away,
       f.winner_team_id, f.created_at, f.updated_at,
       coalesce(json_agg(json_build_object(
           'user_id', lm.user_id,
           'display_name', u.display_name,
           'goals_home', sp.goals_home,
           'goals_away', sp.goals_away,
           'winner', sp.winner,
           'points', sp.points
       ) ORDER BY (lm.user_id = @requesting_user_id) DESC, u.display_name ASC), '[]'::json)::text
       AS member_predictions
FROM fixtures f
JOIN teams home_t ON home_t.id = f.home_team_id
JOIN teams away_t ON away_t.id = f.away_team_id
JOIN leagues l ON l.tournament_id = f.tournament_id
JOIN league_members lm ON lm.league_id = l.id
JOIN users u ON u.id = lm.user_id
LEFT JOIN score_predictions sp ON sp.fixture_id = f.id AND sp.user_id = lm.user_id
WHERE l.id = @league_id
  AND f.prediction_locked = true
  AND f.kickoff_at::date = ANY(@dates::date[])
GROUP BY f.id, home_t.name, away_t.name, home_t.group_letter, f.round
ORDER BY f.kickoff_at DESC;
