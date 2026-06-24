package footballapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alexbukvic2/footy-forecast/internal/worker"
)

// Client implements worker.MatchAPI against api-sports.io v3.
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

// NewClient constructs a Client. If httpClient is nil, http.DefaultClient is used.
func NewClient(apiKey, baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{apiKey: apiKey, baseURL: baseURL, http: httpClient}
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

	var goalScorerIDs []int64
	for _, e := range item.Events {
		if e.Type == "Goal" && e.Detail != "Own Goal" {
			goalScorerIDs = append(goalScorerIDs, e.Player.ID)
		}
	}

	return worker.APIFixtureResult{
		ExternalID:            item.Fixture.ID,
		StatusShort:           item.Fixture.Status.Short,
		GoalsHome:             item.Goals.Home,
		GoalsAway:             item.Goals.Away,
		HomeWinner:            item.Teams.Home.Winner,
		AwayWinner:            item.Teams.Away.Winner,
		GoalScorerExternalIDs: goalScorerIDs,
	}, nil
}

// GetLeagueFixtures fetches all fixtures for a league/season and returns them.
// New fixtures not yet in our DB should be inserted; existing ones are left untouched.
func (c *Client) GetLeagueFixtures(
	ctx context.Context,
	externalLeagueID int64,
	season int,
) ([]worker.APILeagueFixtureResult, error) {
	url := c.baseURL + "/fixtures?season=" + strconv.Itoa(season) +
		"&league=" + strconv.FormatInt(externalLeagueID, 10)
	var resp fixtureResponse
	if err := c.get(ctx, url, &resp); err != nil {
		return nil, fmt.Errorf("get league fixtures league %d season %d: %w", externalLeagueID, season, err)
	}
	results := make([]worker.APILeagueFixtureResult, 0, len(resp.Response))
	for _, item := range resp.Response {
		kickoff, err := time.Parse(time.RFC3339, item.Fixture.Date)
		if err != nil {
			return nil, fmt.Errorf("parse kickoff %q for fixture %d: %w", item.Fixture.Date, item.Fixture.ID, err)
		}
		results = append(results, worker.APILeagueFixtureResult{
			ExternalID:         item.Fixture.ID,
			HomeTeamExternalID: item.Teams.Home.ID,
			AwayTeamExternalID: item.Teams.Away.ID,
			KickoffAt:          kickoff,
			StatusShort:        item.Fixture.Status.Short,
			Round:              item.League.Round,
			GoalsHome:          item.Goals.Home,
			GoalsAway:          item.Goals.Away,
		})
	}
	return results, nil
}

func (c *Client) get(ctx context.Context, url string, dest any) error {
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
