-- +goose Up
-- +goose StatementBegin

CREATE TYPE league_member_role AS ENUM ('owner', 'member');

CREATE TABLE leagues (
    id            UUID        PRIMARY KEY,
    tournament_id UUID        NOT NULL REFERENCES tournaments(id),
    owner_id      UUID        NOT NULL REFERENCES users(id),
    name          TEXT        NOT NULL,
    code          TEXT        NOT NULL UNIQUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT leagues_name_length_chk CHECK (char_length(name) BETWEEN 1 AND 100)
);

CREATE INDEX leagues_tournament_id_idx ON leagues (tournament_id);
CREATE INDEX leagues_owner_id_idx      ON leagues (owner_id);
CREATE INDEX leagues_code_idx          ON leagues (code);

CREATE TABLE league_members (
    league_id UUID               NOT NULL REFERENCES leagues(id) ON DELETE CASCADE,
    user_id   UUID               NOT NULL REFERENCES users(id),
    role      league_member_role NOT NULL DEFAULT 'member',
    joined_at TIMESTAMPTZ        NOT NULL DEFAULT now(),
    PRIMARY KEY (league_id, user_id)
);

CREATE INDEX league_members_user_id_idx ON league_members (user_id);

CREATE TRIGGER leagues_set_updated_at
    BEFORE UPDATE ON leagues
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS leagues_set_updated_at ON leagues;
DROP TABLE IF EXISTS league_members;
DROP TABLE IF EXISTS leagues;
DROP TYPE IF EXISTS league_member_role;

-- +goose StatementEnd
