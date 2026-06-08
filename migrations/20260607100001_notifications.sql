-- +goose Up
ALTER TABLE users ADD COLUMN timezone TEXT NOT NULL DEFAULT 'UTC';
ALTER TABLE users ADD COLUMN silent_from TIME;
ALTER TABLE users ADD COLUMN silent_until TIME;

CREATE TABLE push_tokens (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (token)
);
CREATE INDEX push_tokens_user_id_idx ON push_tokens (user_id);

CREATE TABLE notification_preferences (
    user_id  UUID    NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type     TEXT    NOT NULL,
    enabled  BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (user_id, type)
);

CREATE TABLE notification_log (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type         TEXT        NOT NULL,
    reference_id TEXT        NOT NULL,
    sent_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, type, reference_id)
);
CREATE INDEX notification_log_sent_at_idx ON notification_log (sent_at);

-- +goose Down
DROP TABLE IF EXISTS notification_log;
DROP TABLE IF EXISTS notification_preferences;
DROP TABLE IF EXISTS push_tokens;
ALTER TABLE users DROP COLUMN IF EXISTS silent_until;
ALTER TABLE users DROP COLUMN IF EXISTS silent_from;
ALTER TABLE users DROP COLUMN IF EXISTS timezone;
