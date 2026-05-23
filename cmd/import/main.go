// Command import fetches teams and players from api-football.com and seeds
// the local database. Run once per tournament before accepting predictions.
//
// Required env vars (see .env.import.example):
//
//	API_FOOTBALL_KEY  – api-football.com API key
//	DATABASE_URL      – Postgres DSN
//	LEAGUE_ID         – api-football league integer ID
//	SEASON            – four-digit season year (e.g. 2022)
//	TOURNAMENT_ID     – UUID of the tournament row already present in the DB
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	apiBase      = "https://v3.football.api-sports.io"
	maxNameRunes = 200
	// api-football free tier: 10 req/min. 7 s between pages keeps us well clear.
	pageDelay = 7 * time.Second
)

// ---------- config ----------

type config struct {
	apiKey       string
	databaseURL  string
	leagueID     int
	season       int
	tournamentID uuid.UUID
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
	leagueStr := env("LEAGUE_ID")
	seasonStr := env("SEASON")
	tournamentStr := env("TOURNAMENT_ID")

	if len(missing) > 0 {
		return config{}, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	leagueID, err := strconv.Atoi(leagueStr)
	if err != nil {
		return config{}, fmt.Errorf("LEAGUE_ID must be an integer: %w", err)
	}
	season, err := strconv.Atoi(seasonStr)
	if err != nil {
		return config{}, fmt.Errorf("SEASON must be an integer: %w", err)
	}
	tournamentID, err := uuid.Parse(tournamentStr)
	if err != nil {
		return config{}, fmt.Errorf("TOURNAMENT_ID must be a UUID: %w", err)
	}

	return config{
		apiKey:       apiKey,
		databaseURL:  databaseURL,
		leagueID:     leagueID,
		season:       season,
		tournamentID: tournamentID,
	}, nil
}

// ---------- API response shapes ----------

type teamsResp struct {
	Response []struct {
		Team struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Logo string `json:"logo"`
		} `json:"team"`
	} `json:"response"`
}

type playersResp struct {
	Response []struct {
		Player struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"player"`
		Statistics []struct {
			Team struct {
				ID int `json:"id"`
			} `json:"team"`
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
		logger.Error("import failed", "err", err)
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

	logger.Info("importing teams", "league", cfg.leagueID, "season", cfg.season)
	teamMap, err := importTeams(ctx, logger, client, pool, cfg)
	if err != nil {
		return fmt.Errorf("import teams: %w", err)
	}
	logger.Info("teams done", "count", len(teamMap))

	logger.Info("importing players", "league", cfg.leagueID, "season", cfg.season, "tournament_id", cfg.tournamentID)
	n, err := importPlayers(ctx, logger, client, pool, cfg, teamMap)
	if err != nil {
		return fmt.Errorf("import players: %w", err)
	}
	logger.Info("players done", "inserted", n)
	return nil
}

// ---------- import logic ----------

// importTeams fetches all teams for the league+season, upserts them (updating
// the logo URL on re-runs), and returns a map of api-football team ID → our UUID.
func importTeams(
	ctx context.Context,
	logger *slog.Logger,
	client *http.Client,
	pool *pgxpool.Pool,
	cfg config,
) (map[int]uuid.UUID, error) {
	url := fmt.Sprintf("%s/teams?league=%d&season=%d", apiBase, cfg.leagueID, cfg.season)

	var resp teamsResp
	if err := apiGet(client, cfg.apiKey, url, &resp); err != nil {
		return nil, fmt.Errorf("fetch teams: %w", err)
	}

	teamMap := make(map[int]uuid.UUID, len(resp.Response))
	for _, r := range resp.Response {
		if r.Team.Name == "" {
			logger.Warn("team with empty name skipped", "api_id", r.Team.ID)
			continue
		}
		var id uuid.UUID
		err := pool.QueryRow(
			ctx,
			`INSERT INTO teams (name, logo)
			 VALUES ($1, $2)
			 ON CONFLICT (name) DO UPDATE SET logo = EXCLUDED.logo
			 RETURNING id`,
			r.Team.Name, r.Team.Logo,
		).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("upsert team %q: %w", r.Team.Name, err)
		}
		teamMap[r.Team.ID] = id
		logger.Info("team upserted", "name", r.Team.Name, "api_id", r.Team.ID, "id", id)
	}
	return teamMap, nil
}

// importPlayers paginates through the players API, inserts each player linked to
// a tournament row, and returns the total number of rows inserted. Players whose
// team is not in teamMap (i.e. not in the league) are skipped with a warning.
// Re-runs are idempotent: ON CONFLICT (external_id, tournament_id) DO NOTHING.
func importPlayers(
	ctx context.Context,
	logger *slog.Logger,
	client *http.Client,
	pool *pgxpool.Pool,
	cfg config,
	teamMap map[int]uuid.UUID,
) (int, error) {
	inserted := 0

	for page := 1; ; page++ {
		url := fmt.Sprintf(
			"%s/players?league=%d&season=%d&page=%d",
			apiBase, cfg.leagueID, cfg.season, page,
		)

		var resp playersResp
		if err := apiGet(client, cfg.apiKey, url, &resp); err != nil {
			return inserted, fmt.Errorf("fetch players page %d: %w", page, err)
		}

		for _, r := range resp.Response {
			if len(r.Statistics) == 0 {
				logger.Warn("player has no statistics, skipping", "player_id", r.Player.ID)
				continue
			}

			apiTeamID := r.Statistics[0].Team.ID
			teamID, ok := teamMap[apiTeamID]
			if !ok {
				logger.Warn(
					"player's team not in import set, skipping",
					"player_id", r.Player.ID, "api_team_id", apiTeamID,
				)
				continue
			}

			name := r.Player.Name
			if utf8.RuneCountInString(name) == 0 {
				logger.Warn("player has empty name, skipping", "player_id", r.Player.ID)
				continue
			}
			if utf8.RuneCountInString(name) > maxNameRunes {
				logger.Warn("player name too long, truncating", "player_id", r.Player.ID, "name", name)
				runes := []rune(name)
				name = string(runes[:maxNameRunes])
			}

			externalID := strconv.Itoa(r.Player.ID)
			tag, err := pool.Exec(
				ctx,
				`INSERT INTO players (external_id, name, tournament_id, team_id)
				 VALUES ($1, $2, $3, $4)
				 ON CONFLICT (external_id, tournament_id) DO NOTHING`,
				externalID, name, cfg.tournamentID, teamID,
			)
			if err != nil {
				return inserted, fmt.Errorf("insert player %q (external_id %s): %w", name, externalID, err)
			}
			if tag.RowsAffected() > 0 {
				inserted++
			}
		}

		logger.Info("page processed", "page", page, "total_pages", resp.Paging.Total, "inserted_so_far", inserted)

		if page >= resp.Paging.Total {
			break
		}
		time.Sleep(pageDelay)
	}

	return inserted, nil
}

// ---------- HTTP ----------

// apiGet makes an authenticated GET to the api-football API and decodes the
// JSON response body into dst.
func apiGet(
	client *http.Client,
	apiKey, url string,
	dst any,
) error {
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
