-- +goose NO TRANSACTION
-- ALTER TYPE ... ADD VALUE cannot be used in the same transaction as statements
-- that reference the new value (Postgres SQLSTATE 55P04). NO TRANSACTION makes
-- goose run each statement outside a transaction so ADD VALUE commits before
-- SET DEFAULT tries to use it.

-- +goose Up
ALTER TYPE user_status ADD VALUE 'pending_profile';

-- New users start in pending_profile and must confirm their display name before
-- being fully active. ON CONFLICT (cognito_sub) in UpsertUser does not touch
-- status, so existing active users are unaffected.
ALTER TABLE users ALTER COLUMN status SET DEFAULT 'pending_profile';

-- +goose Down
-- Revert the default. Rows already in pending_profile are promoted to active so
-- the enum value becomes unused. Postgres does not support removing enum values
-- without recreating the type, so pending_profile is intentionally left in place.
UPDATE users SET status = 'active' WHERE status = 'pending_profile';
ALTER TABLE users ALTER COLUMN status SET DEFAULT 'active';
