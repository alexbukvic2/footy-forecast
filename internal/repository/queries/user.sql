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
RETURNING id, cognito_sub, email, display_name, status, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, cognito_sub, email, display_name, status, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByCognitoSub :one
SELECT id, cognito_sub, email, display_name, status, created_at, updated_at
FROM users
WHERE cognito_sub = $1;
