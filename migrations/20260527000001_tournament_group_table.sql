-- +goose Up
-- +goose StatementBegin

CREATE TABLE tournament_group_table (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    tournament_id UUID        NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    team_id       UUID        NOT NULL REFERENCES teams(id)       ON DELETE CASCADE,
    group_letter  CHAR(1)     NOT NULL,
    position      SMALLINT    NOT NULL,
    points        SMALLINT    NOT NULL,
    played        SMALLINT    NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT tournament_group_table_tournament_team_uq UNIQUE (tournament_id, team_id)
);

CREATE INDEX tournament_group_table_tournament_id_idx ON tournament_group_table (tournament_id);

CREATE TRIGGER set_updated_at_tournament_group_table
    BEFORE UPDATE ON tournament_group_table
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS set_updated_at_tournament_group_table ON tournament_group_table;
DROP TABLE IF EXISTS tournament_group_table;

-- +goose StatementEnd
