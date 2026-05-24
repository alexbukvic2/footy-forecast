-- +goose Up
-- +goose StatementBegin

CREATE TYPE fixture_status AS ENUM ('upcoming', 'in_progress', 'finished');

CREATE TABLE fixtures (
    id            UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id   BIGINT         NOT NULL UNIQUE,
    tournament_id UUID           NOT NULL REFERENCES tournaments(id),
    home_team_id  UUID           NOT NULL REFERENCES teams(id),
    away_team_id  UUID           NOT NULL REFERENCES teams(id),
    kickoff_at    TIMESTAMPTZ    NOT NULL,
    status        fixture_status NOT NULL DEFAULT 'upcoming',
    goals_home    INT            CHECK (goals_home >= 0),
    goals_away    INT            CHECK (goals_away >= 0),
    created_at    TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX fixtures_tournament_id_idx ON fixtures (tournament_id);
CREATE INDEX fixtures_kickoff_at_idx    ON fixtures (kickoff_at);

CREATE TRIGGER fixtures_set_updated_at
    BEFORE UPDATE ON fixtures
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS fixtures_set_updated_at ON fixtures;
DROP TABLE IF EXISTS fixtures;
DROP TYPE IF EXISTS fixture_status;

-- +goose StatementEnd
