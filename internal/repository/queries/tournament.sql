-- name: CreateTournament :one
INSERT INTO tournaments (id, slug, name, starts_at, ends_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, slug, name, status, starts_at, ends_at, created_at, updated_at, external_id, season;

-- name: GetTournamentByID :one
SELECT id, slug, name, status, starts_at, ends_at, created_at, updated_at, external_id, season
FROM tournaments
WHERE id = $1;

-- name: GetTournamentBySlug :one
SELECT id, slug, name, status, starts_at, ends_at, created_at, updated_at, external_id, season
FROM tournaments
WHERE slug = $1;

-- name: ListTournaments :many
SELECT id, slug, name, status, starts_at, ends_at, created_at, updated_at, external_id, season
FROM tournaments
ORDER BY starts_at DESC;
