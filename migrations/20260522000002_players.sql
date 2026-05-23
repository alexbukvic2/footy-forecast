-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;

-- unaccent() is STABLE, not IMMUTABLE, so it cannot be used directly in an
-- index expression. This wrapper declares IMMUTABLE so Postgres accepts it.
CREATE OR REPLACE FUNCTION unaccent_immutable(text)
    RETURNS text
    LANGUAGE sql IMMUTABLE PARALLEL SAFE STRICT
AS $$ SELECT unaccent($1) $$;

CREATE TABLE teams (
    id         UUID  PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT  NOT NULL UNIQUE,
    logo       TEXT  NOT NULL DEFAULT ''
);

CREATE TABLE players (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id   TEXT        NOT NULL,
    name          TEXT        NOT NULL,
    tournament_id UUID        NOT NULL REFERENCES tournaments(id),
    team_id       UUID        NOT NULL REFERENCES teams(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT players_external_id_tournament_uq UNIQUE (external_id, tournament_id),
    CONSTRAINT players_name_length_chk CHECK (char_length(name) BETWEEN 1 AND 200)
);

CREATE INDEX players_tournament_id_idx ON players (tournament_id);
CREATE INDEX players_name_trgm_idx     ON players USING GIN (unaccent_immutable(name) gin_trgm_ops);

CREATE TRIGGER players_set_updated_at
    BEFORE UPDATE ON players
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS players_set_updated_at ON players;
DROP TABLE IF EXISTS players;
DROP TABLE IF EXISTS teams;
DROP FUNCTION IF EXISTS unaccent_immutable(text);
-- Extensions left in place intentionally.

-- +goose StatementEnd
