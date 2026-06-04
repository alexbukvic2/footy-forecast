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
	"sort"
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
		League struct {
			Name      string `json:"name"`
			Logo      string `json:"logo"`
			Season    int    `json:"season"`
			Standings [][]struct {
				Team struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
					Logo string `json:"logo"`
				} `json:"team"`
				Points int    `json:"points"`
				Group  string `json:"group"`
				Status string `json:"status"`
			} `json:"standings"`
		} `json:"league"`
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

	logger.Info("truncating tables")
	if err := truncateAll(ctx, pool); err != nil {
		return fmt.Errorf("truncate: %w", err)
	}
	logger.Info("tables truncated")

	logger.Info("importing teams", "league", cfg.leagueID, "season", cfg.season)
	teamMap, err := importTeams(ctx, logger, client, pool, cfg)
	if err != nil {
		return fmt.Errorf("import teams: %w", err)
	}
	logger.Info("teams done", "count", len(teamMap))

	logger.Info("importing fixtures", "league", cfg.leagueID, "season", cfg.season)
	n, err := importFixtures(ctx, logger, client, pool, cfg, teamMap)
	if err != nil {
		return fmt.Errorf("import fixtures: %w", err)
	}
	logger.Info("fixtures done", "upserted", n)

	logger.Info("importing players", "league", cfg.leagueID, "season", cfg.season, "tournament_id", cfg.tournamentID)
	//n, err = importPlayers(ctx, logger, client, pool, cfg, teamMap)
	//if err != nil {
	//	return fmt.Errorf("import players: %w", err)
	//}
	//logger.Info("players done", "inserted", n)
	return nil
}

// ---------- import logic ----------

func truncateAll(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		TRUNCATE
			score_predictions,
			player_predictions,
			team_predictions,
			player_outcomes,
			team_outcomes,
			player_handicap,
			team_handicap,
			league_members,
			leagues,
			fixtures,
			players,
			teams,
			users
		CASCADE
	`)
	return err
}

// importTeams fetches all teams for the league+season, upserts them (updating
// the logo URL on re-runs), and returns a map of api-football team ID → our UUID.
func importTeams(
	ctx context.Context,
	logger *slog.Logger,
	client *http.Client,
	pool *pgxpool.Pool,
	cfg config,
) (map[int]uuid.UUID, error) {
	url := fmt.Sprintf("%s/standings?league=%d&season=%d", apiBase, cfg.leagueID, cfg.season)

	var resp teamsResp
	if err := apiGet(client, cfg.apiKey, url, &resp); err != nil {
		return nil, fmt.Errorf("fetch teams: %w", err)
	}

	teamMap := make(map[int]uuid.UUID, len(resp.Response))
	for _, group := range resp.Response[0].League.Standings {
		for _, r := range group {
			if r.Team.Name == "" {
				logger.Warn("team with empty name skipped", "api_id", r.Team.ID)
				continue
			}
			var id uuid.UUID
			groupLetter := strings.Split(r.Group, " ")[1] // "Group A" → "A"
			err := pool.QueryRow(
				ctx,
				`INSERT INTO teams (name, logo, tournament_id, group_letter)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (name) DO UPDATE SET logo = EXCLUDED.logo
			 RETURNING id`,
				r.Team.Name, r.Team.Logo, cfg.tournamentID, groupLetter,
			).Scan(&id)
			if err != nil {
				return nil, fmt.Errorf("upsert team %q: %w", r.Team.Name, err)
			}
			teamMap[r.Team.ID] = id
			logger.Info("team upserted", "name", r.Team.Name, "api_id", r.Team.ID, "id", id)
		}
	}
	for _, r := range resp.Response[0].League.Standings[0] {
		if r.Team.Name == "" {
			logger.Warn("team with empty name skipped", "api_id", r.Team.ID)
			continue
		}
		var id uuid.UUID
		groupLetter := strings.Split(r.Group, " ")[1] // "Group A" → "A"
		err := pool.QueryRow(
			ctx,
			`INSERT INTO teams (name, logo, tournament_id, group_letter)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (name) DO UPDATE SET logo = EXCLUDED.logo
			 RETURNING id`,
			r.Team.Name, r.Team.Logo, cfg.tournamentID, groupLetter,
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

			tag, err := pool.Exec(
				ctx,
				`INSERT INTO players (external_id, name, tournament_id, team_id)
				 VALUES ($1, $2, $3, $4)
				 ON CONFLICT (external_id, tournament_id) DO NOTHING`,
				int64(r.Player.ID), name, cfg.tournamentID, teamID,
			)
			if err != nil {
				return inserted, fmt.Errorf("insert player %q (external_id %d): %w", name, r.Player.ID, err)
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

// importFixtures fetches all fixtures for the league+season and upserts them.
// The first 20 valid fixtures by kickoff are back-dated (date only, time preserved)
// so they appear to have already happened, with realistic goals and status=finished.
// Re-runs are idempotent. Fixtures whose home or away team is not in teamMap are skipped.
func importFixtures(
	ctx context.Context,
	logger *slog.Logger,
	client *http.Client,
	pool *pgxpool.Pool,
	cfg config,
	teamMap map[int]uuid.UUID,
) (int, error) {
	url := fmt.Sprintf("%s/fixtures?league=%d&season=%d", apiBase, cfg.leagueID, cfg.season)

	var resp fixtures
	if err := apiGet(client, cfg.apiKey, url, &resp); err != nil {
		return 0, fmt.Errorf("fetch fixtures: %w", err)
	}

	sort.Slice(resp.Response, func(i, j int) bool {
		return resp.Response[i].Fixture.Date.Before(resp.Response[j].Fixture.Date)
	})

	const pastCount = 20
	realisticScores := [pastCount][2]int32{
		{1, 0}, {2, 1}, {0, 0}, {3, 1}, {1, 1},
		{2, 0}, {0, 1}, {1, 2}, {2, 2}, {3, 0},
		{0, 2}, {1, 3}, {4, 1}, {2, 1}, {0, 1},
		{1, 0}, {3, 2}, {2, 0}, {1, 1}, {0, 0},
	}

	// Compute how many days to shift back so the 20th valid fixture lands yesterday.
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	var dateShift int
	{
		n := 0
		for _, r := range resp.Response {
			if _, ok := teamMap[r.Teams.Home.Id]; !ok {
				continue
			}
			if _, ok := teamMap[r.Teams.Away.Id]; !ok {
				continue
			}
			n++
			if n == pastCount {
				d := r.Fixture.Date.UTC()
				lastDate := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
				if days := int(lastDate.Sub(today).Hours() / 24); days >= 0 {
					dateShift = days + 1
				}
				break
			}
		}
	}

	upserted := 0
	validIdx := 0
	for _, r := range resp.Response {
		homeID, ok := teamMap[r.Teams.Home.Id]
		if !ok {
			logger.Warn("home team not in import set, skipping fixture",
				"fixture_id", r.Fixture.Id, "api_team_id", r.Teams.Home.Id, "name", r.Teams.Home.Name)
			continue
		}
		awayID, ok := teamMap[r.Teams.Away.Id]
		if !ok {
			logger.Warn("away team not in import set, skipping fixture",
				"fixture_id", r.Fixture.Id, "api_team_id", r.Teams.Away.Id, "name", r.Teams.Away.Name)
			continue
		}

		kickoff := r.Fixture.Date
		status := apiStatusToFixtureStatus(r.Fixture.Status.Short)
		var goalsHome, goalsAway *int32

		switch {
		case validIdx < pastCount:
			kickoff = kickoff.AddDate(0, 0, -dateShift)
			status = "finished"
			h, a := realisticScores[validIdx][0], realisticScores[validIdx][1]
			goalsHome, goalsAway = &h, &a
		case status == "":
			logger.Warn("unrecognised fixture status, skipping",
				"fixture_id", r.Fixture.Id, "status_short", r.Fixture.Status.Short, "status_long", r.Fixture.Status.Long)
			continue
		default:
			goalsHome = jsonNumToInt32(r.Goals.Home)
			goalsAway = jsonNumToInt32(r.Goals.Away)
		}
		validIdx++

		tag, err := pool.Exec(
			ctx,
			`INSERT INTO fixtures (external_id, tournament_id, home_team_id, away_team_id, kickoff_at, status, goals_home, goals_away)
			 VALUES ($1, $2, $3, $4, $5, $6::fixture_status, $7, $8)
			 ON CONFLICT (external_id) DO UPDATE
			 SET tournament_id = EXCLUDED.tournament_id,
			     home_team_id  = EXCLUDED.home_team_id,
			     away_team_id  = EXCLUDED.away_team_id,
			     kickoff_at    = EXCLUDED.kickoff_at,
			     status        = EXCLUDED.status,
			     goals_home    = EXCLUDED.goals_home,
			     goals_away    = EXCLUDED.goals_away`,
			r.Fixture.Id, cfg.tournamentID, homeID, awayID,
			kickoff, status, goalsHome, goalsAway,
		)
		if err != nil {
			return upserted, fmt.Errorf("upsert fixture %d: %w", r.Fixture.Id, err)
		}
		if tag.RowsAffected() > 0 {
			upserted++
		}
		logger.Info("fixture upserted",
			"fixture_id", r.Fixture.Id,
			"home", r.Teams.Home.Name, "away", r.Teams.Away.Name,
			"kickoff", kickoff, "status", status)
	}

	return upserted, nil
}

// apiStatusToFixtureStatus maps api-football short status codes to our DB enum.
// Returns "" for codes that should not be imported (cancelled, abandoned).
func apiStatusToFixtureStatus(short string) string {
	switch short {
	case "NS", "TBD", "PST":
		return "upcoming"
	case "1H", "HT", "2H", "ET", "BT", "P", "SUSP", "INT", "LIVE":
		return "in_progress"
	case "FT", "AET", "PEN", "AWD", "WO":
		return "finished"
	default:
		return ""
	}
}

// jsonNumToInt32 converts a JSON number (decoded as float64) to *int32.
// Returns nil for null or non-numeric values.
func jsonNumToInt32(v any) *int32 {
	if v == nil {
		return nil
	}
	f, ok := v.(float64)
	if !ok {
		return nil
	}
	n := int32(f)
	return &n
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
