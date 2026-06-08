-- +goose Up
-- +goose StatementBegin

CREATE TABLE player_handicap_defaults (
    id            UUID                     PRIMARY KEY DEFAULT gen_random_uuid(),
    tournament_id UUID                     NOT NULL REFERENCES tournaments(id),
    category      player_handicap_category NOT NULL,
    default_points INTEGER                 NOT NULL,
    CONSTRAINT player_handicap_defaults_tournament_category_uq UNIQUE (tournament_id, category)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS player_handicap_defaults;

-- +goose StatementEnd
