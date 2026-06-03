-- name: ListPollableMatches :many
SELECT f.id, f.external_id, f.tournament_id,
       f.home_team_id, f.away_team_id,
       t.external_id  AS tournament_external_id,
       t.season       AS tournament_season,
       home_t.group_letter,
       f.round, f.status, f.kickoff_at,
       f.goals_home, f.goals_away,
       f.winner_team_id,
       COALESCE(f.last_polled_at, '0001-01-01 00:00:00 UTC'::timestamptz) AS last_polled_at
FROM fixtures f
JOIN tournaments t     ON t.id = f.tournament_id
JOIN teams       home_t ON home_t.id = f.home_team_id
WHERE f.is_demo = FALSE
  AND t.external_id IS NOT NULL
  AND (
    f.status = 'in_progress'
    OR (f.status = 'upcoming'  AND f.kickoff_at <= now() + INTERVAL '5 minutes')
    OR (f.status = 'finished'  AND f.updated_at >= now() - INTERVAL '24 hours')
  )
ORDER BY f.kickoff_at ASC;

-- name: CountIncompleteGroupFixtures :one
SELECT COUNT(*)::int AS n
FROM fixtures f
JOIN teams home_t ON home_t.id = f.home_team_id
WHERE f.tournament_id = @tournament_id
  AND home_t.group_letter = @group_letter
  AND f.status NOT IN ('finished', 'cancelled');

-- name: CountIncompleteRoundFixtures :one
SELECT COUNT(*)::int AS n
FROM fixtures f
WHERE f.tournament_id = @tournament_id
  AND f.round = @round
  AND f.winner_team_id IS NULL;

-- name: CountIncompleteGroupStageFixtures :one
SELECT COUNT(*)::int AS n
FROM fixtures f
JOIN teams home_t ON home_t.id = f.home_team_id
WHERE f.tournament_id = @tournament_id
  AND home_t.group_letter IS NOT NULL
  AND f.status NOT IN ('finished', 'cancelled');

-- name: GetTeamByExternalID :one
SELECT id FROM teams
WHERE external_id = @external_id
  AND tournament_id = @tournament_id
LIMIT 1;

-- name: GetPlayerByExternalID :one
SELECT id FROM players
WHERE external_id = @external_id
  AND tournament_id = @tournament_id;

-- name: UpdateFixtureStatus :exec
UPDATE fixtures
SET status     = @status,
    kickoff_at = @kickoff_at
WHERE id = @id;
