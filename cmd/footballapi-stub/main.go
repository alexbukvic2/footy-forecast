// Command footballapi-stub runs a local HTTP stub for the api-sports.io football
// API. It serves the exact JSON shape that internal/footballapi.Client parses,
// and exposes a /control prefix for updating fixture state on the fly.
//
// Usage:
//
//	go run ./cmd/footballapi-stub          # listens on :9001 by default
//	STUB_PORT=9002 go run ./cmd/footballapi-stub
//
// Then set FOOTBALL_API_BASE_URL=http://localhost:9001 in .env and start the
// server normally. Advance a fixture's state with:
//
//	curl -s -X PUT localhost:9001/control/fixtures/123 \
//	  -H 'Content-Type: application/json' \
//	  -d '{"status":"1H","goals_home":0,"goals_away":0}'
//
//	curl -s -X PUT localhost:9001/control/fixtures/123 \
//	  -d '{"status":"FT","goals_home":2,"goals_away":1,"home_winner":true}'
//
//	curl -s localhost:9001/control/fixtures      # inspect current state
package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// --- in-memory state ---

type fixtureState struct {
	ID         int64  `json:"id"`
	Status     string `json:"status"`     // e.g. "NS", "1H", "HT", "FT"
	GoalsHome  *int   `json:"goals_home"` // nil until match starts
	GoalsAway  *int   `json:"goals_away"`
	HomeWinner *bool  `json:"home_winner"` // nil unless finished
	AwayWinner *bool  `json:"away_winner"`
}

type standingsState struct {
	Teams []standingsTeamState `json:"teams"`
}

type standingsTeamState struct {
	TeamID      int64  `json:"team_id"`
	Rank        int    `json:"rank"`
	Points      int    `json:"points"`
	Played      int    `json:"played"`
	Won         int    `json:"won"`
	Drawn       int    `json:"drawn"`
	Lost        int    `json:"lost"`
	GoalsFor    int    `json:"goals_for"`
	Against     int    `json:"goals_against"`
	Description string `json:"description"` // e.g. "Promotion - Championship (Group Stage: 1)"
}

type topScorerState struct {
	PlayerID int64 `json:"player_id"`
	Goals    int   `json:"goals"`
}

type leagueSeason struct {
	LeagueID int64
	Season   int
}

type store struct {
	mu         sync.RWMutex
	fixtures   map[int64]fixtureState
	standings  map[leagueSeason]standingsState
	topScorers map[leagueSeason]topScorerState
}

func newStore() *store {
	return &store{
		fixtures:   make(map[int64]fixtureState),
		standings:  make(map[leagueSeason]standingsState),
		topScorers: make(map[leagueSeason]topScorerState),
	}
}

// --- api-sports.io JSON shapes (mirrors internal/footballapi/types.go) ---

type apiFixtureResponse struct {
	Response []apiFixtureItem `json:"response"`
}

type apiFixtureItem struct {
	Fixture apiFixtureDetails `json:"fixture"`
	Goals   apiGoals          `json:"goals"`
	Teams   apiTeams          `json:"teams"`
}

type apiFixtureDetails struct {
	ID     int64            `json:"id"`
	Status apiFixtureStatus `json:"status"`
}

type apiFixtureStatus struct {
	Short string `json:"short"`
}

type apiGoals struct {
	Home *int `json:"home"`
	Away *int `json:"away"`
}

type apiTeams struct {
	Home apiTeamResult `json:"home"`
	Away apiTeamResult `json:"away"`
}

type apiTeamResult struct {
	Winner *bool `json:"winner"`
}

type apiStandingsResponse struct {
	Response []apiStandingsItem `json:"response"`
}

type apiStandingsItem struct {
	League apiStandingsLeague `json:"league"`
}

type apiStandingsLeague struct {
	Standings [][]apiStandingsTeam `json:"standings"`
}

type apiStandingsTeam struct {
	Rank        int              `json:"rank"`
	Team        apiTeamRef       `json:"team"`
	Points      int              `json:"points"`
	All         apiStandingsStat `json:"all"`
	Description string           `json:"description"`
}

type apiTeamRef struct {
	ID int64 `json:"id"`
}

type apiStandingsStat struct {
	Played int          `json:"played"`
	Win    int          `json:"win"`
	Draw   int          `json:"draw"`
	Lose   int          `json:"lose"`
	Goals  apiStatGoals `json:"goals"`
}

type apiStatGoals struct {
	For     int `json:"for"`
	Against int `json:"against"`
}

type apiTopScorersResponse struct {
	Response []apiTopScorerItem `json:"response"`
}

type apiTopScorerItem struct {
	Player apiPlayerRef  `json:"player"`
	Stats  []apiGoalStat `json:"statistics"`
}

type apiPlayerRef struct {
	ID int64 `json:"id"`
}

type apiGoalStat struct {
	Goals apiGoalTotal `json:"goals"`
}

type apiGoalTotal struct {
	Total *int `json:"total"`
}

// --- handlers ---

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encode response", "err", err)
	}
}

func (s *store) handleGetFixture(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	f, ok := s.fixtures[id]
	s.mu.RUnlock()

	if !ok {
		writeJSON(w, apiFixtureResponse{Response: []apiFixtureItem{}})
		return
	}

	item := apiFixtureItem{
		Fixture: apiFixtureDetails{
			ID:     f.ID,
			Status: apiFixtureStatus{Short: f.Status},
		},
		Goals: apiGoals{Home: f.GoalsHome, Away: f.GoalsAway},
		Teams: apiTeams{
			Home: apiTeamResult{Winner: f.HomeWinner},
			Away: apiTeamResult{Winner: f.AwayWinner},
		},
	}
	writeJSON(w, apiFixtureResponse{Response: []apiFixtureItem{item}})
}

func (s *store) handleGetStandings(w http.ResponseWriter, r *http.Request) {
	leagueID, err := strconv.ParseInt(r.URL.Query().Get("league"), 10, 64)
	if err != nil {
		http.Error(w, "bad league", http.StatusBadRequest)
		return
	}
	season, err := strconv.Atoi(r.URL.Query().Get("season"))
	if err != nil {
		http.Error(w, "bad season", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	st, ok := s.standings[leagueSeason{leagueID, season}]
	s.mu.RUnlock()

	if !ok || len(st.Teams) == 0 {
		writeJSON(w, apiStandingsResponse{Response: []apiStandingsItem{}})
		return
	}

	teams := make([]apiStandingsTeam, len(st.Teams))
	for i, t := range st.Teams {
		teams[i] = apiStandingsTeam{
			Rank:        t.Rank,
			Team:        apiTeamRef{ID: t.TeamID},
			Points:      t.Points,
			Description: t.Description,
			All: apiStandingsStat{
				Played: t.Played,
				Win:    t.Won,
				Draw:   t.Drawn,
				Lose:   t.Lost,
				Goals:  apiStatGoals{For: t.GoalsFor, Against: t.Against},
			},
		}
	}

	writeJSON(w, apiStandingsResponse{
		Response: []apiStandingsItem{{League: apiStandingsLeague{Standings: [][]apiStandingsTeam{teams}}}},
	})
}

func (s *store) handleGetTopScorers(w http.ResponseWriter, r *http.Request) {
	leagueID, err := strconv.ParseInt(r.URL.Query().Get("league"), 10, 64)
	if err != nil {
		http.Error(w, "bad league", http.StatusBadRequest)
		return
	}
	season, err := strconv.Atoi(r.URL.Query().Get("season"))
	if err != nil {
		http.Error(w, "bad season", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	ts, ok := s.topScorers[leagueSeason{leagueID, season}]
	s.mu.RUnlock()

	if !ok {
		writeJSON(w, apiTopScorersResponse{Response: []apiTopScorerItem{}})
		return
	}

	goals := ts.Goals
	writeJSON(w, apiTopScorersResponse{
		Response: []apiTopScorerItem{{
			Player: apiPlayerRef{ID: ts.PlayerID},
			Stats:  []apiGoalStat{{Goals: apiGoalTotal{Total: &goals}}},
		}},
	})
}

// --- control handlers ---

func (s *store) handleControlListFixtures(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	out := make([]fixtureState, 0, len(s.fixtures))
	for _, f := range s.fixtures {
		out = append(out, f)
	}
	s.mu.RUnlock()
	writeJSON(w, out)
}

func (s *store) handleControlSetFixture(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}

	var body fixtureState
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body: "+err.Error(), http.StatusBadRequest)
		return
	}
	body.ID = id

	s.mu.Lock()
	s.fixtures[id] = body
	s.mu.Unlock()

	slog.Info("fixture updated",
		"id", id,
		"status", body.Status,
		"goals_home", intPtrStr(body.GoalsHome),
		"goals_away", intPtrStr(body.GoalsAway),
	)
	writeJSON(w, body)
}

func (s *store) handleControlSetStandings(w http.ResponseWriter, r *http.Request) {
	leagueID, err := strconv.ParseInt(r.PathValue("league"), 10, 64)
	if err != nil {
		http.Error(w, "bad league", http.StatusBadRequest)
		return
	}
	season, err := strconv.Atoi(r.PathValue("season"))
	if err != nil {
		http.Error(w, "bad season", http.StatusBadRequest)
		return
	}

	var body standingsState
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body: "+err.Error(), http.StatusBadRequest)
		return
	}

	key := leagueSeason{leagueID, season}
	s.mu.Lock()
	s.standings[key] = body
	s.mu.Unlock()

	slog.Info("standings updated", "league", leagueID, "season", season, "teams", len(body.Teams))
	writeJSON(w, body)
}

func (s *store) handleControlSetTopScorer(w http.ResponseWriter, r *http.Request) {
	leagueID, err := strconv.ParseInt(r.PathValue("league"), 10, 64)
	if err != nil {
		http.Error(w, "bad league", http.StatusBadRequest)
		return
	}
	season, err := strconv.Atoi(r.PathValue("season"))
	if err != nil {
		http.Error(w, "bad season", http.StatusBadRequest)
		return
	}

	var body topScorerState
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body: "+err.Error(), http.StatusBadRequest)
		return
	}

	key := leagueSeason{leagueID, season}
	s.mu.Lock()
	s.topScorers[key] = body
	s.mu.Unlock()

	slog.Info("top scorer updated", "league", leagueID, "season", season, "player", body.PlayerID, "goals", body.Goals)
	writeJSON(w, body)
}

func intPtrStr(p *int) string {
	if p == nil {
		return "nil"
	}
	return strconv.Itoa(*p)
}

// --- main ---

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	port := os.Getenv("STUB_PORT")
	if port == "" {
		port = "9001"
	}

	s := newStore()
	mux := http.NewServeMux()

	// Football API routes (what footballapi.Client calls)
	mux.HandleFunc("GET /fixtures", s.handleGetFixture)
	mux.HandleFunc("GET /standings", s.handleGetStandings)
	mux.HandleFunc("GET /players/topscorers", s.handleGetTopScorers)

	// Control routes
	mux.HandleFunc("GET /control/fixtures", s.handleControlListFixtures)
	mux.HandleFunc("PUT /control/fixtures/{id}", s.handleControlSetFixture)
	mux.HandleFunc("PUT /control/standings/{league}/{season}", s.handleControlSetStandings)
	mux.HandleFunc("PUT /control/topscorer/{league}/{season}", s.handleControlSetTopScorer)

	addr := ":" + port
	slog.Info("footballapi-stub listening",
		"addr", addr,
		"tip", fmt.Sprintf("set FOOTBALL_API_BASE_URL=http://localhost:%s in .env", port),
	)
	slog.Info("control endpoints",
		"list_fixtures", "GET  /control/fixtures",
		"set_fixture", "PUT  /control/fixtures/{id}",
		"set_standings", "PUT  /control/standings/{league}/{season}",
		"set_topscorer", "PUT  /control/topscorer/{league}/{season}",
	)
	slog.Info("fixture status codes",
		"upcoming", "NS TBD PST CANC ABD",
		"in_progress", "1H HT 2H ET BT P SUSP INT",
		"finished", "FT AET PEN AWD WO",
	)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}
