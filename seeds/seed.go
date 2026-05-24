// seed loads JSON fixtures from seeds/data/ into the database.
// Run manually: go run ./seeds [--db <url>]
// It is idempotent — every INSERT uses ON CONFLICT DO NOTHING.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbURL := flag.String("db", os.Getenv("DATABASE_URL"), "postgres connection URL")
	flag.Parse()

	if *dbURL == "" {
		slog.Error("DATABASE_URL not set and --db not provided")
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, *dbURL)
	if err != nil {
		slog.Error("connect", "err", err)
		os.Exit(1)
	}

	steps := []struct {
		name string
		fn   func(context.Context, *pgxpool.Pool) (int, error)
	}{
		{"tournament", seedTournaments},
		{"teams", seedTeams},
		{"players", seedPlayers},
		{"users", seedUsers},
		{"leagues", seedLeagues},
		{"league_members", seedLeagueMembers},
		{"team_handicap", seedTeamHandicap},
		{"player_handicap", seedPlayerHandicap},
		{"fixtures", seedFixtures},
		{"score_predictions", seedScorePredictions},
	}

	for _, s := range steps {
		n, err := s.fn(ctx, pool)
		if err != nil {
			pool.Close()
			slog.Error("seed failed", "step", s.name, "err", err)
			os.Exit(1)
		}
		slog.Info("seeded", "step", s.name, "inserted", n)
	}
	pool.Close()
}

// ── helpers ──────────────────────────────────────────────────────────────────

func readJSON[T any](filename string) ([]T, error) {
	path := filepath.Join("seeds", "data", filename)
	f, err := os.Open(path) //nolint:gosec // path is constructed from a trusted static prefix
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", filename, err)
	}
	defer func() { _ = f.Close() }()
	var rows []T
	if err := json.NewDecoder(f).Decode(&rows); err != nil {
		return nil, fmt.Errorf("decode %s: %w", filename, err)
	}
	return rows, nil
}

// ── seeders ──────────────────────────────────────────────────────────────────

type tournament struct {
	ID       string `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
}

func seedTournaments(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	rows, err := readJSON[tournament]("tournament.json")
	if err != nil {
		return 0, err
	}
	var n int
	for _, r := range rows {
		startsAt, err := time.Parse(time.RFC3339, r.StartsAt)
		if err != nil {
			return 0, fmt.Errorf("parse starts_at %q: %w", r.StartsAt, err)
		}
		endsAt, err := time.Parse(time.RFC3339, r.EndsAt)
		if err != nil {
			return 0, fmt.Errorf("parse ends_at %q: %w", r.EndsAt, err)
		}
		tag, err := pool.Exec(ctx, `
			INSERT INTO tournaments (id, slug, name, status, starts_at, ends_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO NOTHING`,
			r.ID, r.Slug, r.Name, r.Status, startsAt, endsAt,
		)
		if err != nil {
			return 0, fmt.Errorf("tournament %s: %w", r.ID, err)
		}
		n += int(tag.RowsAffected())
	}
	return n, nil
}

type team struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Logo         string `json:"logo"`
	TournamentID string `json:"tournament_id"`
}

func seedTeams(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	rows, err := readJSON[team]("teams.json")
	if err != nil {
		return 0, err
	}
	var n int
	for _, r := range rows {
		tag, err := pool.Exec(ctx, `
			INSERT INTO teams (id, name, logo, tournament_id)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO NOTHING`,
			r.ID, r.Name, r.Logo, r.TournamentID,
		)
		if err != nil {
			return 0, fmt.Errorf("team %s: %w", r.ID, err)
		}
		n += int(tag.RowsAffected())
	}
	return n, nil
}

type player struct {
	ID           string `json:"id"`
	ExternalID   string `json:"external_id"`
	Name         string `json:"name"`
	TournamentID string `json:"tournament_id"`
	TeamID       string `json:"team_id"`
}

func seedPlayers(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	rows, err := readJSON[player]("players.json")
	if err != nil {
		return 0, err
	}
	var n int
	for _, r := range rows {
		tag, err := pool.Exec(ctx, `
			INSERT INTO players (id, external_id, name, tournament_id, team_id)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO NOTHING`,
			r.ID, r.ExternalID, r.Name, r.TournamentID, r.TeamID,
		)
		if err != nil {
			return 0, fmt.Errorf("player %s: %w", r.ID, err)
		}
		n += int(tag.RowsAffected())
	}
	return n, nil
}

type user struct {
	ID          string `json:"id"`
	CognitoSub  string `json:"cognito_sub"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
}

func seedUsers(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	rows, err := readJSON[user]("users.json")
	if err != nil {
		return 0, err
	}
	var n int
	for _, r := range rows {
		tag, err := pool.Exec(ctx, `
			INSERT INTO users (id, cognito_sub, email, display_name, status)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO NOTHING`,
			r.ID, r.CognitoSub, r.Email, r.DisplayName, r.Status,
		)
		if err != nil {
			return 0, fmt.Errorf("user %s: %w", r.ID, err)
		}
		n += int(tag.RowsAffected())
	}
	return n, nil
}

type league struct {
	ID           string `json:"id"`
	TournamentID string `json:"tournament_id"`
	OwnerID      string `json:"owner_id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
}

func seedLeagues(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	rows, err := readJSON[league]("leagues.json")
	if err != nil {
		return 0, err
	}
	var n int
	for _, r := range rows {
		tag, err := pool.Exec(ctx, `
			INSERT INTO leagues (id, tournament_id, owner_id, name, code)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO NOTHING`,
			r.ID, r.TournamentID, r.OwnerID, r.Name, r.Code,
		)
		if err != nil {
			return 0, fmt.Errorf("league %s: %w", r.ID, err)
		}
		n += int(tag.RowsAffected())
	}
	return n, nil
}

type leagueMember struct {
	LeagueID string `json:"league_id"`
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
}

func seedLeagueMembers(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	rows, err := readJSON[leagueMember]("league_members.json")
	if err != nil {
		return 0, err
	}
	var n int
	for _, r := range rows {
		tag, err := pool.Exec(ctx, `
			INSERT INTO league_members (league_id, user_id, role)
			VALUES ($1, $2, $3)
			ON CONFLICT (league_id, user_id) DO NOTHING`,
			r.LeagueID, r.UserID, r.Role,
		)
		if err != nil {
			return 0, fmt.Errorf("league_member %s/%s: %w", r.LeagueID, r.UserID, err)
		}
		n += int(tag.RowsAffected())
	}
	return n, nil
}

type teamHandicap struct {
	ID       string `json:"id"`
	TeamID   string `json:"team_id"`
	Category string `json:"category"`
	Points   int    `json:"points"`
}

func seedTeamHandicap(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	rows, err := readJSON[teamHandicap]("team_handicap.json")
	if err != nil {
		return 0, err
	}
	var n int
	for _, r := range rows {
		tag, err := pool.Exec(ctx, `
			INSERT INTO team_handicap (id, team_id, category, points)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (team_id, category) DO NOTHING`,
			r.ID, r.TeamID, r.Category, r.Points,
		)
		if err != nil {
			return 0, fmt.Errorf("team_handicap %s: %w", r.ID, err)
		}
		n += int(tag.RowsAffected())
	}
	return n, nil
}

type playerHandicap struct {
	ID       string `json:"id"`
	PlayerID string `json:"player_id"`
	Category string `json:"category"`
	Points   int    `json:"points"`
}

func seedPlayerHandicap(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	rows, err := readJSON[playerHandicap]("player_handicap.json")
	if err != nil {
		return 0, err
	}
	var n int
	for _, r := range rows {
		tag, err := pool.Exec(ctx, `
			INSERT INTO player_handicap (id, player_id, category, points)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (player_id, category) DO NOTHING`,
			r.ID, r.PlayerID, r.Category, r.Points,
		)
		if err != nil {
			return 0, fmt.Errorf("player_handicap %s: %w", r.ID, err)
		}
		n += int(tag.RowsAffected())
	}
	return n, nil
}

type fixture struct {
	ID           string `json:"id"`
	ExternalID   int64  `json:"external_id"`
	TournamentID string `json:"tournament_id"`
	HomeTeamID   string `json:"home_team_id"`
	AwayTeamID   string `json:"away_team_id"`
	KickoffAt    string `json:"kickoff_at"`
	Status       string `json:"status"`
	GoalsHome    *int   `json:"goals_home"`
	GoalsAway    *int   `json:"goals_away"`
}

func seedFixtures(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	rows, err := readJSON[fixture]("fixtures.json")
	if err != nil {
		return 0, err
	}
	var n int
	for _, r := range rows {
		kickoff, err := time.Parse(time.RFC3339, r.KickoffAt)
		if err != nil {
			return 0, fmt.Errorf("parse kickoff_at %q: %w", r.KickoffAt, err)
		}
		tag, err := pool.Exec(ctx, `
			INSERT INTO fixtures (id, external_id, tournament_id, home_team_id, away_team_id, kickoff_at, status, goals_home, goals_away)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (id) DO NOTHING`,
			r.ID, r.ExternalID, r.TournamentID, r.HomeTeamID, r.AwayTeamID,
			kickoff, r.Status, r.GoalsHome, r.GoalsAway,
		)
		if err != nil {
			return 0, fmt.Errorf("fixture %s: %w", r.ID, err)
		}
		n += int(tag.RowsAffected())
	}
	return n, nil
}

type scorePrediction struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	FixtureID string `json:"fixture_id"`
	GoalsHome int    `json:"goals_home"`
	GoalsAway int    `json:"goals_away"`
	Points    *int   `json:"points"`
}

func seedScorePredictions(ctx context.Context, pool *pgxpool.Pool) (int, error) {
	rows, err := readJSON[scorePrediction]("score_predictions.json")
	if err != nil {
		return 0, err
	}
	var n int
	for _, r := range rows {
		tag, err := pool.Exec(ctx, `
			INSERT INTO score_predictions (id, user_id, fixture_id, goals_home, goals_away, points)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO NOTHING`,
			r.ID, r.UserID, r.FixtureID, r.GoalsHome, r.GoalsAway, r.Points,
		)
		if err != nil {
			return 0, fmt.Errorf("score_prediction %s: %w", r.ID, err)
		}
		n += int(tag.RowsAffected())
	}
	return n, nil
}
