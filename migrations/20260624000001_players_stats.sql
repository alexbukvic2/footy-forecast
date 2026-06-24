-- +goose Up
-- +goose StatementBegin
CREATE TABLE players_stats (
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id),
    goals     INT  NOT NULL DEFAULT 0,
    CONSTRAINT players_stats_player_uq UNIQUE (player_id)
);

CREATE INDEX players_stats_player_id_idx ON players_stats (player_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS players_stats;
-- +goose StatementEnd
