-- +goose Up
-- +goose StatementBegin

CREATE TABLE score_predictions (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id),
    fixture_id UUID        NOT NULL REFERENCES fixtures(id),
    goals_home INT         NOT NULL CHECK (goals_home >= 0),
    goals_away INT         NOT NULL CHECK (goals_away >= 0),
    points     INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT score_predictions_user_fixture_uq UNIQUE (user_id, fixture_id)
);

CREATE INDEX score_predictions_user_id_idx    ON score_predictions (user_id);
CREATE INDEX score_predictions_fixture_id_idx ON score_predictions (fixture_id);

CREATE TRIGGER score_predictions_set_updated_at
    BEFORE UPDATE ON score_predictions
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS score_predictions_set_updated_at ON score_predictions;
DROP TABLE IF EXISTS score_predictions;

-- +goose StatementEnd
