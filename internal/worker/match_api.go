package worker

import "context"

// MatchAPI abstracts the external football data provider.
type MatchAPI interface {
	GetFixture(ctx context.Context, externalFixtureID int64) (APIFixtureResult, error)
	GetStandings(ctx context.Context, externalLeagueID int64, season int) ([]APIStandingsEntry, error)
	GetGroupTopScorer(ctx context.Context, externalLeagueID int64, season int, groupLetter string) ([]APITopScorerResult, error)
	GetTournamentTopScorer(ctx context.Context, externalLeagueID int64, season int) ([]APITopScorerResult, error)
}

// APIFixtureResult is the data returned by GetFixture.
type APIFixtureResult struct {
	ExternalID  int64
	StatusShort string // raw API status: "1H", "HT", "FT", "CANC", etc.
	GoalsHome   *int
	GoalsAway   *int
	HomeWinner  *bool // teams.home.winner; handles ET/PEN
	AwayWinner  *bool
}

// APIStandingsEntry is one team row from GetStandings.
type APIStandingsEntry struct {
	TeamExternalID int64
	Position       int
	Points         int
	Played         int
	Won            int
	Drawn          int
	Lost           int
	GoalsFor       int
	GoalsAgainst   int
	Description    string // e.g. "Promotion - Championship (Group Stage: 1)", "Relegation"
}

// APITopScorerResult is returned by GetGroupTopScorer and GetTournamentTopScorer.
type APITopScorerResult struct {
	PlayerExternalID string // matches players.external_id
	Goals            int
}
