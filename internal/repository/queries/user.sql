-- name: UpsertUser :one
INSERT INTO users (id, cognito_sub, email, display_name)
VALUES ($1, $2, $3, $4)
ON CONFLICT (cognito_sub) DO UPDATE
    SET email        = EXCLUDED.email,
        display_name = CASE
                           WHEN users.display_name = '' THEN EXCLUDED.display_name
                           ELSE users.display_name
                       END,
        updated_at   = now()
RETURNING id, cognito_sub, email, display_name, status, created_at, updated_at, timezone, silent_from, silent_until;

-- name: GetUserByID :one
SELECT id, cognito_sub, email, display_name, status, created_at, updated_at, timezone, silent_from, silent_until
FROM users
WHERE id = $1;

-- name: GetUserByCognitoSub :one
SELECT id, cognito_sub, email, display_name, status, created_at, updated_at, timezone, silent_from, silent_until
FROM users
WHERE cognito_sub = $1;

-- name: UpdateDisplayName :one
UPDATE users
SET display_name = $2,
    status       = CASE WHEN status = 'pending_profile' THEN 'active'::user_status ELSE status END
WHERE id = $1
RETURNING id, cognito_sub, email, display_name, status, created_at, updated_at, timezone, silent_from, silent_until;

-- name: UpdateTimezone :one
UPDATE users
SET timezone     = $2,
    silent_from  = $3,
    silent_until = $4,
    updated_at   = now()
WHERE id = $1
RETURNING id, cognito_sub, email, display_name, status, created_at, updated_at, timezone, silent_from, silent_until;
