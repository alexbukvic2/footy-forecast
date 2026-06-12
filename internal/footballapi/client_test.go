//go:build integration

package footballapi

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testLeagueID / testSeason are the fixtures the user nominated for mapping validation.
const (
	testLeagueID int64 = 1
	testSeason         = 2022
)

// testClient returns a real API client, skipping the test if FOOTBALL_API_KEY is unset.
func testClient(t *testing.T) *Client {
	t.Helper()
	key := os.Getenv("API_FOOTBALL_KEY")
	if key == "" {
		t.Skip("API_FOOTBALL_KEY not set — skipping live API test")
	}
	baseURL := os.Getenv("FOOTBALL_API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://v3.football.api-sports.io"
	}
	return NewClient(key, baseURL, nil)
}

// TestClient_GetTournamentTopScorer verifies that the top-scorers response decodes
// correctly and that all returned players share the same (highest) goal count.
func TestClient_GetTournamentTopScorer(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	results, err := c.GetTournamentTopScorer(ctx, testLeagueID, testSeason)
	require.NoError(t, err)
	require.NotEmpty(t, results, "expected at least one top scorer")

	topGoals := results[0].Goals
	for _, r := range results {
		assert.NotEmpty(t, r.PlayerExternalID, "PlayerExternalID must not be empty")
		assert.Positive(t, r.Goals, "top scorer must have scored at least one goal")
		assert.Equal(t, topGoals, r.Goals, "all returned players must share the top goal count")
	}

	t.Logf("top scorers: count=%d goals=%d first_player=%d", len(results), topGoals, results[0].PlayerExternalID)
}

// TestClient_GetGroupTopScorer verifies the group-top-scorer path (which falls back
// to the full tournament top scorers endpoint for leagues that don't support
// group filtering).
func TestClient_GetGroupTopScorer(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	results, err := c.GetGroupTopScorer(ctx, testLeagueID, testSeason, "A")
	require.NoError(t, err)
	require.NotEmpty(t, results, "expected at least one top scorer")

	topGoals := results[0].Goals
	for _, r := range results {
		assert.NotEmpty(t, r.PlayerExternalID)
		assert.Positive(t, r.Goals)
		assert.Equal(t, topGoals, r.Goals, "all returned players must share the top goal count")
	}

	t.Logf(
		"group top scorer (league-wide fallback): count=%d goals=%d first_player=%d",
		len(results), topGoals, results[0].PlayerExternalID,
	)
}

// TestClient_GetFixture first fetches any finished fixture for league 1 / season 2022
// using the private get() helper (white-box, same package), then calls GetFixture
// with that ID to validate the full struct mapping.
func TestClient_GetFixture(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	// Discover a known-finished fixture ID via the raw endpoint.
	var discovery struct {
		Response []struct {
			Fixture struct {
				ID int64 `json:"id"`
			} `json:"fixture"`
		} `json:"response"`
	}
	discoverURL := c.baseURL + "/fixtures?league=1&season=2022&status=FT"
	require.NoError(t, c.get(ctx, discoverURL, &discovery), "discover fixtures")
	require.NotEmpty(t, discovery.Response, "expected at least one finished fixture for league=1 season=2022")

	fixtureID := discovery.Response[0].Fixture.ID
	require.Positive(t, fixtureID, "discovered fixture ID must be > 0")

	result, err := c.GetFixture(ctx, fixtureID)
	require.NoError(t, err)

	assert.Equal(t, fixtureID, result.ExternalID, "ExternalID must round-trip")
	assert.NotEmpty(t, result.StatusShort, "StatusShort must not be empty")
	// Finished fixtures must have non-nil goal counts.
	require.NotNil(t, result.GoalsHome, "GoalsHome must not be nil for a finished fixture")
	require.NotNil(t, result.GoalsAway, "GoalsAway must not be nil for a finished fixture")
	assert.GreaterOrEqual(t, *result.GoalsHome, 0)
	assert.GreaterOrEqual(t, *result.GoalsAway, 0)
	// Exactly one winner flag should be set (or both nil if draw).
	if result.HomeWinner != nil || result.AwayWinner != nil {
		homeWins := result.HomeWinner != nil && *result.HomeWinner
		awayWins := result.AwayWinner != nil && *result.AwayWinner
		assert.False(t, homeWins && awayWins, "HomeWinner and AwayWinner must not both be true")
	}

	t.Logf(
		"fixture %d: status=%s home_goals=%d away_goals=%d",
		fixtureID, result.StatusShort, *result.GoalsHome, *result.GoalsAway,
	)
}
