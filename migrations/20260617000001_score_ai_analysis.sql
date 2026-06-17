-- +goose Up
CREATE TABLE score_ai_analysis (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    fixture_id UUID        NOT NULL REFERENCES fixtures(id),
    league_id  UUID        NOT NULL REFERENCES leagues(id),
    analysis   TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (fixture_id, league_id)
);

-- +goose Down
DROP TABLE score_ai_analysis;
