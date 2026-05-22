-- name: CreateLeague :one
INSERT INTO leagues (id, tournament_id, owner_id, name, code)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, tournament_id, owner_id, name, code, created_at, updated_at;

-- name: GetLeagueByID :one
SELECT id, tournament_id, owner_id, name, code, created_at, updated_at
FROM leagues
WHERE id = $1;

-- name: GetLeagueByCode :one
SELECT id, tournament_id, owner_id, name, code, created_at, updated_at
FROM leagues
WHERE code = $1;

-- name: ListLeaguesForUser :many
SELECT l.id, l.tournament_id, l.owner_id, l.name, l.code, l.created_at, l.updated_at
FROM leagues l
JOIN league_members lm ON lm.league_id = l.id
WHERE lm.user_id = $1
ORDER BY l.created_at DESC;

-- name: UpdateLeagueName :one
UPDATE leagues
SET name = $2
WHERE id = $1
RETURNING id, tournament_id, owner_id, name, code, created_at, updated_at;

-- name: UpdateLeagueCode :one
UPDATE leagues
SET code = $2
WHERE id = $1
RETURNING id, tournament_id, owner_id, name, code, created_at, updated_at;

-- name: DeleteLeague :exec
DELETE FROM leagues
WHERE id = $1;

-- name: AddLeagueMember :one
INSERT INTO league_members (league_id, user_id, role)
VALUES ($1, $2, $3)
RETURNING league_id, user_id, role, joined_at;

-- name: RemoveLeagueMember :execresult
DELETE FROM league_members
WHERE league_id = $1 AND user_id = $2;

-- name: GetLeagueMember :one
SELECT league_id, user_id, role, joined_at
FROM league_members
WHERE league_id = $1 AND user_id = $2;

-- name: ListLeagueMembersForLeague :many
SELECT league_id, user_id, role, joined_at
FROM league_members
WHERE league_id = $1
ORDER BY joined_at ASC;
