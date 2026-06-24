// Command seed-player-stats fetches goal tallies from api-football.com for
// every team in the database and upserts them into players_stats.
//
// Required env vars:
//
//	API_FOOTBALL_KEY  – api-football.com API key
//	DATABASE_URL      – Postgres DSN
//
// The script is idempotent: re-runs update goal counts in place.
// Only players_stats is written; teams and players are read-only.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	apiBase  = "https://v3.football.api-sports.io"
	leagueID = 1
	season   = 2026
	// Pro plan: 300 r/m. ~60 r/m already consumed by other callers.
	// Cap this script at 200 r/m → 300 ms between requests.
	pageDelay = 300 * time.Millisecond
)

type config struct {
	apiKey      string
	databaseURL string
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

	apiKey := env("API_FOOTBALL_KEY")
	databaseURL := env("DATABASE_URL")

	if len(missing) > 0 {
		return config{}, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}
	return config{apiKey: apiKey, databaseURL: databaseURL}, nil
}

// ---------- API shapes ----------

type playersResp struct {
	Response []struct {
		Player struct {
			ID int64 `json:"id"`
		} `json:"player"`
		Statistics []struct {
			Goals struct {
				Total *int `json:"total"`
			} `json:"goals"`
		} `json:"statistics"`
	} `json:"response"`
	Paging struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"paging"`
}

// ---------- main ----------

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	if err := run(logger); err != nil {
		logger.Error("seed-player-stats failed", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.databaseURL)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	client := &http.Client{Timeout: 30 * time.Second}

	// Fetch all teams with their external_id.
	type team struct {
		externalID int64
		name       string
	}
	rows, err := pool.Query(ctx, `SELECT external_id, name FROM teams ORDER BY name`)
	if err != nil {
		return fmt.Errorf("query teams: %w", err)
	}
	var teams []team
	for rows.Next() {
		var t team
		if err := rows.Scan(&t.externalID, &t.name); err != nil {
			return fmt.Errorf("scan team row: %w", err)
		}
		teams = append(teams, t)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate teams: %w", err)
	}

	logger.Info("teams loaded", "count", len(teams))

	total := 0
	for i, t := range teams {
		if i > 0 {
			time.Sleep(pageDelay)
		}

		n, err := processTeam(ctx, logger, client, pool, cfg.apiKey, t.externalID, t.name)
		if err != nil {
			return fmt.Errorf("process team %q: %w", t.name, err)
		}
		total += n
	}

	logger.Info("done", "players_stats_upserted", total)
	return nil
}

func processTeam(
	ctx context.Context,
	logger *slog.Logger,
	client *http.Client,
	pool *pgxpool.Pool,
	apiKey string,
	teamExternalID int64,
	teamName string,
) (int, error) {
	upserted := 0

	for page := 1; ; page++ {
		if page > 1 {
			time.Sleep(pageDelay)
		}

		url := fmt.Sprintf(
			"%s/players?season=%d&league=%d&team=%d&page=%d",
			apiBase, season, leagueID, teamExternalID, page,
		)

		var resp playersResp
		if err := apiGet(client, apiKey, url, &resp); err != nil {
			return upserted, fmt.Errorf("fetch players page %d for team %d: %w", page, teamExternalID, err)
		}

		for _, r := range resp.Response {
			if len(r.Statistics) == 0 {
				continue
			}

			goalsTotal := r.Statistics[0].Goals.Total
			if goalsTotal == nil || *goalsTotal <= 0 {
				continue
			}

			// Resolve internal player ID from external_id.
			var playerID string
			err := pool.QueryRow(
				ctx,
				`SELECT id FROM players WHERE external_id = $1`,
				r.Player.ID,
			).Scan(&playerID)
			if err != nil {
				logger.Warn("player not found in DB, skipping",
					"external_id", r.Player.ID, "team", teamName)
				continue
			}

			_, err = pool.Exec(
				ctx,
				`INSERT INTO players_stats (player_id, goals)
				 VALUES ($1, $2)
				 ON CONFLICT (player_id) DO UPDATE SET goals = EXCLUDED.goals`,
				playerID, *goalsTotal,
			)
			if err != nil {
				return upserted, fmt.Errorf("upsert players_stats for player %d: %w", r.Player.ID, err)
			}
			upserted++
			logger.Info("upserted", "player_external_id", r.Player.ID, "team", teamName, "goals", *goalsTotal)
		}

		logger.Info("page done", "team", teamName, "page", page, "total_pages", resp.Paging.Total)

		if page >= resp.Paging.Total {
			break
		}
	}

	return upserted, nil
}

// apiGet makes an authenticated GET to the api-football API and decodes the
// JSON response body into dst.
func apiGet(client *http.Client, apiKey, url string, dst any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("x-apisports-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status %d for %s", resp.StatusCode, url)
	}

	if err := json.NewDecoder(resp.Body).Decode(dst); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
