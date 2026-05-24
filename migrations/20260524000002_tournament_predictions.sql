-- +goose Up
-- +goose StatementBegin

CREATE TABLE player_predictions (
    id            UUID                     PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID                     NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tournament_id UUID                     NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    category      player_handicap_category NOT NULL,
    pick          UUID                     NOT NULL REFERENCES players(id),
    points        INT,
    created_at    TIMESTAMPTZ              NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ              NOT NULL DEFAULT now(),
    CONSTRAINT player_predictions_user_tournament_category_uq UNIQUE (user_id, tournament_id, category)
);

CREATE INDEX player_predictions_tournament_id_idx ON player_predictions (tournament_id);

CREATE TRIGGER set_updated_at_player_predictions
    BEFORE UPDATE ON player_predictions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE team_predictions (
    id            UUID                   PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID                   NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tournament_id UUID                   NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
    category      team_handicap_category NOT NULL,
    pick          UUID                   NOT NULL REFERENCES teams(id),
    points        INT,
    created_at    TIMESTAMPTZ            NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ            NOT NULL DEFAULT now(),
    CONSTRAINT team_predictions_user_tournament_category_uq UNIQUE (user_id, tournament_id, category)
);

CREATE INDEX team_predictions_tournament_id_idx ON team_predictions (tournament_id);

CREATE TRIGGER set_updated_at_team_predictions
    BEFORE UPDATE ON team_predictions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS set_updated_at_team_predictions   ON team_predictions;
DROP TRIGGER IF EXISTS set_updated_at_player_predictions ON player_predictions;
DROP TABLE IF EXISTS team_predictions;
DROP TABLE IF EXISTS player_predictions;

-- +goose StatementEnd
