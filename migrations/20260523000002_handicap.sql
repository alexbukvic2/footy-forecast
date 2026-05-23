-- +goose Up
-- +goose StatementBegin

CREATE TYPE player_handicap_category AS ENUM (
    'group_top_scorer',
    'total_top_scorer'
);

CREATE TYPE team_handicap_category AS ENUM (
    'group_winner',
    'playoff',
    'semifinalist',
    'winner'
);

CREATE TABLE player_handicap (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id   UUID        NOT NULL REFERENCES players(id),
    category    player_handicap_category NOT NULL,
    points      INTEGER     NOT NULL,
    CONSTRAINT player_handicap_player_category_uq UNIQUE (player_id, category)
);

CREATE TABLE team_handicap (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id     UUID        NOT NULL REFERENCES teams(id),
    category    team_handicap_category NOT NULL,
    points      INTEGER     NOT NULL,
    CONSTRAINT team_handicap_team_category_uq UNIQUE (team_id, category)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS team_handicap;
DROP TABLE IF EXISTS player_handicap;
DROP TYPE IF EXISTS team_handicap_category;
DROP TYPE IF EXISTS player_handicap_category;

-- +goose StatementEnd
