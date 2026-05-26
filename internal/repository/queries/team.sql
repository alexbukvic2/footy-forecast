-- name: GetTeamByID :one
SELECT id, name, logo, tournament_id, group_letter
FROM teams
WHERE id = @id;

-- name: ListGroupLettersByTournament :many
SELECT DISTINCT group_letter
FROM teams
WHERE tournament_id = @tournament_id
  AND group_letter IS NOT NULL
ORDER BY group_letter ASC;

-- name: ListTeamsWithHandicapsByTournament :many
SELECT
    t.id,
    t.name,
    t.logo,
    t.tournament_id,
    t.group_letter,
    th.category AS handicap_category,
    th.points   AS handicap_points
FROM teams t
LEFT JOIN team_handicap th ON th.team_id = t.id
WHERE t.tournament_id = @tournament_id
ORDER BY t.name ASC, th.category ASC;

