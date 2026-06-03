package footballapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alexbukvic2/footy-forecast/internal/worker"
)

// Client implements worker.MatchAPI against api-sports.io v3.
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

// NewClient constructs a Client. If httpClient is nil, http.DefaultClient is used.
func NewClient(
	apiKey, baseURL string,
	httpClient *http.Client,
) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		http:    httpClient,
	}
}

var _ worker.MatchAPI = (*Client)(nil)

// GetFixture fetches the current state of a single fixture from the API.
func (c *Client) GetFixture(
	ctx context.Context,
	externalFixtureID int64,
) (worker.APIFixtureResult, error) {
	url := c.baseURL + "/fixtures?id=" + strconv.FormatInt(externalFixtureID, 10)
	var resp fixtureResponse
	if err := c.get(ctx, url, &resp); err != nil {
		return worker.APIFixtureResult{}, fmt.Errorf("get fixture %d: %w", externalFixtureID, err)
	}
	if len(resp.Response) == 0 {
		return worker.APIFixtureResult{}, fmt.Errorf("get fixture %d: empty response", externalFixtureID)
	}
	item := resp.Response[0]
	return worker.APIFixtureResult{
		ExternalID:  item.Fixture.ID,
		StatusShort: item.Fixture.Status.Short,
		GoalsHome:   item.Goals.Home,
		GoalsAway:   item.Goals.Away,
		HomeWinner:  item.Teams.Home.Winner,
		AwayWinner:  item.Teams.Away.Winner,
	}, nil
}

// GetStandings fetches the current standings for a league/season.
func (c *Client) GetStandings(
	ctx context.Context,
	externalLeagueID int64,
	season int,
) ([]worker.APIStandingsEntry, error) {
	url := c.baseURL + "/standings?league=" + strconv.FormatInt(externalLeagueID, 10) +
		"&season=" + strconv.Itoa(season)
	var resp standingsResponse
	if err := c.get(ctx, url, &resp); err != nil {
		return nil, fmt.Errorf("get standings league %d season %d: %w", externalLeagueID, season, err)
	}
	if len(resp.Response) == 0 {
		return nil, fmt.Errorf("get standings league %d season %d: empty response", externalLeagueID, season)
	}

	var entries []worker.APIStandingsEntry
	for _, league := range resp.Response {
		for _, group := range league.League.Standings {
			for _, t := range group {
				entries = append(
					entries, worker.APIStandingsEntry{
						TeamExternalID: t.Team.ID,
						Position:       t.Rank,
						Points:         t.Points,
						Played:         t.All.Played,
						Won:            t.All.Win,
						Drawn:          t.All.Draw,
						Lost:           t.All.Lose,
						GoalsFor:       t.All.Goals.For,
						GoalsAgainst:   t.All.Goals.Against,
					},
				)
			}
		}
	}
	return entries, nil
}

// GetGroupTopScorer returns the top scorer for a specific group stage group.
// api-sports.io does not have a group-filtered top scorers endpoint, so this
// fetches all top scorers for the league/season and returns the first result.
// See plan open question 6.
func (c *Client) GetGroupTopScorer(
	ctx context.Context,
	externalLeagueID int64,
	season int,
	_ string,
) (worker.APITopScorerResult, error) {
	return c.getTopScorer(ctx, externalLeagueID, season)
}

// GetTournamentTopScorer returns the overall top scorer for the tournament.
func (c *Client) GetTournamentTopScorer(
	ctx context.Context,
	externalLeagueID int64,
	season int,
) (worker.APITopScorerResult, error) {
	return c.getTopScorer(ctx, externalLeagueID, season)
}

func (c *Client) getTopScorer(
	ctx context.Context,
	externalLeagueID int64,
	season int,
) (worker.APITopScorerResult, error) {
	url := c.baseURL + "/players/topscorers?league=" + strconv.FormatInt(externalLeagueID, 10) +
		"&season=" + strconv.Itoa(season)
	var resp topScorersResponse
	if err := c.get(ctx, url, &resp); err != nil {
		return worker.APITopScorerResult{}, fmt.Errorf(
			"get top scorers league %d season %d: %w",
			externalLeagueID,
			season,
			err,
		)
	}
	if len(resp.Response) == 0 {
		return worker.APITopScorerResult{}, fmt.Errorf(
			"get top scorers league %d season %d: empty response",
			externalLeagueID,
			season,
		)
	}
	item := resp.Response[0]
	goals := 0
	if len(item.Stats) > 0 && item.Stats[0].Goals.Total != nil {
		goals = *item.Stats[0].Goals.Total
	}
	return worker.APITopScorerResult{
		PlayerExternalID: strconv.FormatInt(item.Player.ID, 10),
		Goals:            goals,
	}, nil
}

func (c *Client) get(
	ctx context.Context,
	url string,
	dest interface{},
) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("x-apisports-key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
