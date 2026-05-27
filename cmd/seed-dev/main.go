// Command seed-dev loads fixture JSON files and seeds a local/dev database with
// deterministic test data: 6 teams (2 groups), 12 players, 6 fixtures, 5 users,
// handicaps, and predictions at various levels of completeness.
//
// Required env vars (can also be provided via .env):
//
//	DATABASE_URL   – Postgres DSN
//	TOURNAMENT_ID  – UUID of an existing tournament row
//	LEAGUE_ID      – UUID of an existing league row (seeded users are added as members)
package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

//go:embed data/teams.json
var teamsData []byte

//go:embed data/players.json
var playersData []byte

//go:embed data/fixtures.json
var fixturesData []byte

//go:embed data/users.json
var usersData []byte

//go:embed data/handicaps.json
var handicapsData []byte

//go:embed data/predictions.json
var predictionsData []byte

// ---------- config ----------

type config struct {
	databaseURL  string
	tournamentID uuid.UUID
	leagueID     uuid.UUID
}

func loadConfig() (config, error) {
	var missing []string
	env := func(key string) string {
		v := os.Getenv(key)
		if v == "" {
			missing = append(missing, key)
		}
		return v
	}

	databaseURL := env("DATABASE_URL")
	tournamentStr := env("TOURNAMENT_ID")
	leagueStr := env("LEAGUE_ID")

	if len(missing) > 0 {
		return config{}, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	tournamentID, err := uuid.Parse(tournamentStr)
	if err != nil {
		return config{}, fmt.Errorf("TOURNAMENT_ID must be a UUID: %w", err)
	}
	leagueID, err := uuid.Parse(leagueStr)
	if err != nil {
		return config{}, fmt.Errorf("LEAGUE_ID must be a UUID: %w", err)
	}

	return config{databaseURL: databaseURL, tournamentID: tournamentID, leagueID: leagueID}, nil
}

// ---------- seed data types ----------

type seedTeam struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	GroupLetter string `json:"group_letter"`
}

type seedPlayer struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	Name       string `json:"name"`
	TeamID     string `json:"team_id"`
}

type seedFixture struct {
	ID         string `json:"id"`
	ExternalID int64  `json:"external_id"`
	HomeTeamID string `json:"home_team_id"`
	AwayTeamID string `json:"away_team_id"`
	Round      string `json:"round"`
	KickoffAt  string `json:"kickoff_at"`
	Status     string `json:"status"`
	GoalsHome  *int32 `json:"goals_home"`
	GoalsAway  *int32 `json:"goals_away"`
	Locked     bool   `json:"locked"`
}

type seedUser struct {
	ID          string `json:"id"`
	CognitoSub  string `json:"cognito_sub"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type seedHandicaps struct {
	PlayerHandicap []struct {
		PlayerID string `json:"player_id"`
		Category string `json:"category"`
		Points   int    `json:"points"`
	} `json:"player_handicap"`
	TeamHandicap []struct {
		TeamID   string `json:"team_id"`
		Category string `json:"category"`
		Points   int    `json:"points"`
	} `json:"team_handicap"`
}

type seedPredictions struct {
	PlayerPredictions []struct {
		UserID      string  `json:"user_id"`
		Category    string  `json:"category"`
		Pick        string  `json:"pick"`
		GroupLetter *string `json:"group_letter"`
		Points      *int    `json:"points"`
	} `json:"player_predictions"`
	TeamPredictions []struct {
		UserID      string  `json:"user_id"`
		Category    string  `json:"category"`
		Pick        string  `json:"pick"`
		GroupLetter *string `json:"group_letter"`
		SlotIndex   int16   `json:"slot_index"`
		Points      *int    `json:"points"`
	} `json:"team_predictions"`
	ScorePredictions []struct {
		UserID    string `json:"user_id"`
		FixtureID string `json:"fixture_id"`
		GoalsHome int    `json:"goals_home"`
		GoalsAway int    `json:"goals_away"`
		Points    *int   `json:"points"`
	} `json:"score_predictions"`
}

// ---------- main ----------

func main() {
	_ = godotenv.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	if err := run(logger); err != nil {
		logger.Error("seed-dev failed", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	var teams []seedTeam
	var players []seedPlayer
	var fixtures []seedFixture
	var users []seedUser
	var handicaps seedHandicaps
	var predictions seedPredictions

	for _, p := range []struct {
		src  []byte
		dst  any
		name string
	}{
		{teamsData, &teams, "teams.json"},
		{playersData, &players, "players.json"},
		{fixturesData, &fixtures, "fixtures.json"},
		{usersData, &users, "users.json"},
		{handicapsData, &handicaps, "handicaps.json"},
		{predictionsData, &predictions, "predictions.json"},
	} {
		if err = json.Unmarshal(p.src, p.dst); err != nil {
			return fmt.Errorf("parse %s: %w", p.name, err)
		}
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.databaseURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	logger.Info("truncating tables")
	if err := truncateAll(ctx, pool); err != nil {
		return fmt.Errorf("truncate: %w", err)
	}
	logger.Info("tables truncated")

	teamIDs := extractIDs(teams, func(t seedTeam) string { return t.ID })
	playerIDs := extractIDs(players, func(p seedPlayer) string { return p.ID })
	fixtureIDs := extractIDs(fixtures, func(f seedFixture) string { return f.ID })
	userIDs := extractIDs(users, func(u seedUser) string { return u.ID })

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	logger.Info("cleaning existing seed data")
	if err = clean(ctx, tx, teamIDs, playerIDs, fixtureIDs, userIDs, cfg.leagueID); err != nil {
		return fmt.Errorf("clean: %w", err)
	}

	steps := []struct {
		label string
		fn    func() error
	}{
		{"teams", func() error { return insertTeams(ctx, tx, teams, cfg.tournamentID) }},
		{"players", func() error { return insertPlayers(ctx, tx, players, cfg.tournamentID) }},
		{"fixtures", func() error { return insertFixtures(ctx, tx, fixtures, cfg.tournamentID) }},
		{"users", func() error { return insertUsers(ctx, tx, users) }},
		{"league_members", func() error { return insertLeagueMembers(ctx, tx, users, cfg.leagueID) }},
		{"player_handicap", func() error { return insertPlayerHandicap(ctx, tx, handicaps) }},
		{"team_handicap", func() error { return insertTeamHandicap(ctx, tx, handicaps) }},
		{"player_predictions", func() error {
			return insertPlayerPredictions(ctx, tx, predictions, cfg.tournamentID)
		}},
		{"team_predictions", func() error {
			return insertTeamPredictions(ctx, tx, predictions, cfg.tournamentID)
		}},
		{"score_predictions", func() error { return insertScorePredictions(ctx, tx, predictions) }},
	}
	for _, s := range steps {
		logger.Info("seeding " + s.label)
		if err = s.fn(); err != nil {
			return fmt.Errorf("seed %s: %w", s.label, err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	logger.Info(
		"seed-dev complete",
		"teams", len(teams),
		"players", len(players),
		"fixtures", len(fixtures),
		"users", len(users),
		"team_predictions", len(predictions.TeamPredictions),
		"player_predictions", len(predictions.PlayerPredictions),
		"score_predictions", len(predictions.ScorePredictions),
	)
	return nil
}

// ---------- insert helpers ----------

func insertTeams(
	ctx context.Context,
	tx pgx.Tx,
	teams []seedTeam,
	tournamentID uuid.UUID,
) error {
	for _, t := range teams {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO teams (id, name, logo, tournament_id, group_letter)
			 VALUES ($1, $2, $3, $4, $5)`,
			t.ID, t.Name, t.Logo, tournamentID, t.GroupLetter,
		); err != nil {
			return fmt.Errorf("team %q: %w", t.Name, err)
		}
	}
	return nil
}

func insertPlayers(
	ctx context.Context,
	tx pgx.Tx,
	players []seedPlayer,
	tournamentID uuid.UUID,
) error {
	for _, p := range players {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO players (id, external_id, name, tournament_id, team_id)
			 VALUES ($1, $2, $3, $4, $5)`,
			p.ID, p.ExternalID, p.Name, tournamentID, p.TeamID,
		); err != nil {
			return fmt.Errorf("player %q: %w", p.Name, err)
		}
	}
	return nil
}

func insertFixtures(
	ctx context.Context,
	tx pgx.Tx,
	fixtures []seedFixture,
	tournamentID uuid.UUID,
) error {
	for _, f := range fixtures {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO fixtures
			   (id, external_id, tournament_id, home_team_id, away_team_id,
			    round, kickoff_at, status, goals_home, goals_away, prediction_locked)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8::fixture_status, $9, $10, $11)`,
			f.ID, f.ExternalID, tournamentID,
			f.HomeTeamID, f.AwayTeamID,
			f.Round, f.KickoffAt, f.Status, f.GoalsHome, f.GoalsAway, f.Locked,
		); err != nil {
			return fmt.Errorf("fixture external_id=%d: %w", f.ExternalID, err)
		}
	}
	return nil
}

func insertUsers(
	ctx context.Context,
	tx pgx.Tx,
	users []seedUser,
) error {
	for _, u := range users {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO users (id, cognito_sub, email, display_name)
			 VALUES ($1, $2, $3, $4)`,
			u.ID, u.CognitoSub, u.Email, u.DisplayName,
		); err != nil {
			return fmt.Errorf("user %q: %w", u.Email, err)
		}
	}
	return nil
}

func insertLeagueMembers(
	ctx context.Context,
	tx pgx.Tx,
	users []seedUser,
	leagueID uuid.UUID,
) error {
	for _, u := range users {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO league_members (league_id, user_id, role) VALUES ($1, $2, 'member')`,
			leagueID, u.ID,
		); err != nil {
			return fmt.Errorf("league_member user=%s: %w", u.ID, err)
		}
	}
	return nil
}

func insertPlayerHandicap(
	ctx context.Context,
	tx pgx.Tx,
	h seedHandicaps,
) error {
	for _, ph := range h.PlayerHandicap {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO player_handicap (player_id, category, points)
			 VALUES ($1, $2::player_handicap_category, $3)`,
			ph.PlayerID, ph.Category, ph.Points,
		); err != nil {
			return fmt.Errorf("player_handicap player=%s category=%s: %w", ph.PlayerID, ph.Category, err)
		}
	}
	return nil
}

func insertTeamHandicap(
	ctx context.Context,
	tx pgx.Tx,
	h seedHandicaps,
) error {
	for _, th := range h.TeamHandicap {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO team_handicap (team_id, category, points)
			 VALUES ($1, $2::team_handicap_category, $3)`,
			th.TeamID, th.Category, th.Points,
		); err != nil {
			return fmt.Errorf("team_handicap team=%s category=%s: %w", th.TeamID, th.Category, err)
		}
	}
	return nil
}

func insertPlayerPredictions(
	ctx context.Context,
	tx pgx.Tx,
	p seedPredictions,
	tournamentID uuid.UUID,
) error {
	for _, pp := range p.PlayerPredictions {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO player_predictions
			   (user_id, tournament_id, category, pick, group_letter, points)
			 VALUES ($1, $2, $3::player_handicap_category, $4, $5, $6)`,
			pp.UserID, tournamentID, pp.Category, pp.Pick, pp.GroupLetter, pp.Points,
		); err != nil {
			return fmt.Errorf("player_prediction user=%s category=%s: %w", pp.UserID, pp.Category, err)
		}
	}
	return nil
}

func insertTeamPredictions(
	ctx context.Context,
	tx pgx.Tx,
	p seedPredictions,
	tournamentID uuid.UUID,
) error {
	for _, tp := range p.TeamPredictions {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO team_predictions
			   (user_id, tournament_id, category, pick, group_letter, slot_index, points)
			 VALUES ($1, $2, $3::team_handicap_category, $4, $5, $6, $7)`,
			tp.UserID, tournamentID, tp.Category, tp.Pick, tp.GroupLetter, tp.SlotIndex, tp.Points,
		); err != nil {
			return fmt.Errorf(
				"team_prediction user=%s category=%s slot=%d: %w",
				tp.UserID, tp.Category, tp.SlotIndex, err,
			)
		}
	}
	return nil
}

func insertScorePredictions(
	ctx context.Context,
	tx pgx.Tx,
	p seedPredictions,
) error {
	for _, sp := range p.ScorePredictions {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO score_predictions (user_id, fixture_id, goals_home, goals_away, points)
			 VALUES ($1, $2, $3, $4, $5)`,
			sp.UserID, sp.FixtureID, sp.GoalsHome, sp.GoalsAway, sp.Points,
		); err != nil {
			return fmt.Errorf("score_prediction user=%s fixture=%s: %w", sp.UserID, sp.FixtureID, err)
		}
	}
	return nil
}

// ---------- cleanup ----------

// clean removes all rows inserted by this seed so the command is idempotent.
// Deletions are ordered to respect FK constraints.
func clean(
	ctx context.Context,
	tx pgx.Tx,
	teamIDs, playerIDs, fixtureIDs, userIDs []uuid.UUID,
	leagueID uuid.UUID,
) error {
	steps := []struct {
		label string
		query string
		args  []any
	}{
		{"score_predictions", `DELETE FROM score_predictions WHERE user_id = ANY($1)`, []any{userIDs}},
		{"team_predictions", `DELETE FROM team_predictions WHERE user_id = ANY($1)`, []any{userIDs}},
		{"player_predictions", `DELETE FROM player_predictions WHERE user_id = ANY($1)`, []any{userIDs}},
		{"league_members", `DELETE FROM league_members WHERE user_id = ANY($1) AND league_id = $2`, []any{userIDs, leagueID}},
		{"player_handicap", `DELETE FROM player_handicap WHERE player_id = ANY($1)`, []any{playerIDs}},
		{"team_handicap", `DELETE FROM team_handicap WHERE team_id = ANY($1)`, []any{teamIDs}},
		{"fixtures", `DELETE FROM fixtures WHERE id = ANY($1)`, []any{fixtureIDs}},
		{"players", `DELETE FROM players WHERE id = ANY($1)`, []any{playerIDs}},
		{"teams", `DELETE FROM teams WHERE id = ANY($1)`, []any{teamIDs}},
		{"users", `DELETE FROM users WHERE id = ANY($1)`, []any{userIDs}},
	}
	for _, s := range steps {
		if _, err := tx.Exec(ctx, s.query, s.args...); err != nil {
			return fmt.Errorf("clean %s: %w", s.label, err)
		}
	}
	return nil
}

// extractIDs maps a slice to a []uuid.UUID using the provided key function.
func extractIDs[T any](
	slice []T,
	key func(T) string,
) []uuid.UUID {
	out := make([]uuid.UUID, len(slice))
	for i, item := range slice {
		out[i] = uuid.MustParse(key(item))
	}
	return out
}

func truncateAll(
	ctx context.Context,
	pool *pgxpool.Pool,
) error {
	_, err := pool.Exec(
		ctx, `
		TRUNCATE
			score_predictions,
			player_predictions,
			team_predictions,
			player_outcomes,
			team_outcomes,
			player_handicap,
			team_handicap,
			league_members,
			fixtures,
			players,
			teams
		CASCADE
	`,
	)
	return err
}
