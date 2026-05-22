-- +goose Up
-- +goose StatementBegin

CREATE TYPE user_status AS ENUM ('active', 'suspended');

CREATE TABLE users (
    id            UUID        PRIMARY KEY,
    cognito_sub   TEXT        NOT NULL UNIQUE,
    email         TEXT        NOT NULL UNIQUE,
    display_name  TEXT        NOT NULL DEFAULT '',
    status        user_status NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX users_email_idx ON users (email);

-- Reuse set_updated_at() created in the tournaments migration.
CREATE TRIGGER users_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS users_set_updated_at ON users;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS user_status;

-- NOTE: set_updated_at() is NOT dropped here — tournaments still depends on it.

-- +goose StatementEnd
