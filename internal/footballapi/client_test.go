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

// TestClient_GetStandings verifies that the standings response decodes correctly
// and that every field in APIStandingsEntry carries a sensible value.
func TestClient_GetStandings(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	entries, err := c.GetStandings(ctx, testLeagueID, testSeason)
	require.NoError(t, err)
	require.NotEmpty(t, entries, "expected at least one standings entry")

	for _, e := range entries {
		assert.Positive(t, e.TeamExternalID, "TeamExternalID must be > 0")
		assert.Positive(t, e.Position, "Position must be > 0")
		assert.GreaterOrEqual(t, e.Points, 0, "Points must be >= 0")
		assert.Positive(t, e.Played, "Played must be > 0 for a completed season")
		assert.Equal(t, e.Played, e.Won+e.Drawn+e.Lost,
			"Played must equal Won+Drawn+Lost for team external_id=%d", e.TeamExternalID)
		assert.GreaterOrEqual(t, e.GoalsFor, 0)
		assert.GreaterOrEqual(t, e.GoalsAgainst, 0)
	}

	t.Logf("standings: %d entries; leader external_id=%d pts=%d played=%d",
		len(entries), entries[0].TeamExternalID, entries[0].Points, entries[0].Played)
}

// TestClient_GetTournamentTopScorer verifies that the top-scorers response decodes
// correctly and that the returned player has a non-empty ID and at least one goal.
func TestClient_GetTournamentTopScorer(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	result, err := c.GetTournamentTopScorer(ctx, testLeagueID, testSeason)
	require.NoError(t, err)

	assert.NotEmpty(t, result.PlayerExternalID, "PlayerExternalID must not be empty")
	assert.Positive(t, result.Goals, "top scorer must have scored at least one goal")

	t.Logf("top scorer: player_external_id=%s goals=%d", result.PlayerExternalID, result.Goals)
}

// TestClient_GetGroupTopScorer verifies the group-top-scorer path (which falls back
// to the full tournament top scorers endpoint for leagues that don't support
// group filtering).
func TestClient_GetGroupTopScorer(t *testing.T) {
	c := testClient(t)
	ctx := context.Background()

	result, err := c.GetGroupTopScorer(ctx, testLeagueID, testSeason, "A")
	require.NoError(t, err)

	assert.NotEmpty(t, result.PlayerExternalID)
	assert.Positive(t, result.Goals)

	t.Logf("group top scorer (league-wide fallback): player_external_id=%s goals=%d",
		result.PlayerExternalID, result.Goals)
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

	t.Logf("fixture %d: status=%s home_goals=%d away_goals=%d",
		fixtureID, result.StatusShort, *result.GoalsHome, *result.GoalsAway)
}
