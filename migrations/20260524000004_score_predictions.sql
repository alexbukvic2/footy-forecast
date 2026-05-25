-- +goose Up
-- +goose StatementBegin

CREATE TABLE score_predictions (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    fixture_id  UUID        NOT NULL REFERENCES fixtures(id) ON DELETE CASCADE,
    goals_home  INT         NOT NULL CHECK (goals_home >= 0),
    goals_away  INT         NOT NULL CHECK (goals_away >= 0),
    points      INT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT score_predictions_user_fixture_uq UNIQUE (user_id, fixture_id)
);

CREATE TRIGGER set_updated_at_score_predictions
    BEFORE UPDATE ON score_predictions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS set_updated_at_score_predictions ON score_predictions;
DROP TABLE   IF EXISTS score_predictions;

-- +goose StatementEnd
