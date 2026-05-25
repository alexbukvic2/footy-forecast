-- +goose Up
-- +goose StatementBegin

CREATE TABLE player_outcomes (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tournament_id UUID NOT NULL REFERENCES tournaments(id),
    category      player_handicap_category NOT NULL,
    player_id     UUID NOT NULL REFERENCES players(id),
    recorded_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT player_outcomes_tournament_category_player_uq
        UNIQUE (tournament_id, category, player_id)
);

CREATE TABLE team_outcomes (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tournament_id UUID NOT NULL REFERENCES tournaments(id),
    category      team_handicap_category NOT NULL,
    team_id       UUID NOT NULL REFERENCES teams(id),
    recorded_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT team_outcomes_tournament_category_team_uq
        UNIQUE (tournament_id, category, team_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS team_outcomes;
DROP TABLE IF EXISTS player_outcomes;

-- +goose StatementEnd
