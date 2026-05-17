-- +goose Up
-- +goose StatementBegin

CREATE TYPE tournament_status AS ENUM ('upcoming', 'in_progress', 'concluded');

CREATE TABLE tournaments (
                             id              UUID PRIMARY KEY,
                             slug            TEXT NOT NULL UNIQUE,
                             name            TEXT NOT NULL,
                             status          tournament_status NOT NULL DEFAULT 'upcoming',
                             starts_at       TIMESTAMPTZ NOT NULL,
                             ends_at         TIMESTAMPTZ NOT NULL,
                             created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                             updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                             CONSTRAINT tournaments_dates_chk CHECK (ends_at > starts_at)
);

CREATE INDEX tournaments_status_idx ON tournaments (status);
CREATE INDEX tournaments_starts_at_idx ON tournaments (starts_at);

-- Trigger to keep updated_at fresh on every UPDATE.
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tournaments_set_updated_at
    BEFORE UPDATE ON tournaments
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS tournaments_set_updated_at ON tournaments;
DROP FUNCTION IF EXISTS set_updated_at();
DROP TABLE IF EXISTS tournaments;
DROP TYPE IF EXISTS tournament_status;

-- +goose StatementEnd
