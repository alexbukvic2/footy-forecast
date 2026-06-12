package worker

import (
	"context"
	"time"
)

// MatchAPI abstracts the external football data provider.
type MatchAPI interface {
	GetFixture(ctx context.Context, externalFixtureID int64) (APIFixtureResult, error)
	GetGroupTopScorer(ctx context.Context, externalLeagueID int64, season int, groupLetter string) ([]APITopScorerResult, error)
	GetTournamentTopScorer(ctx context.Context, externalLeagueID int64, season int) ([]APITopScorerResult, error)
	GetLeagueFixtures(ctx context.Context, externalLeagueID int64, season int) ([]APILeagueFixtureResult, error)
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

// APITopScorerResult is returned by GetGroupTopScorer and GetTournamentTopScorer.
type APITopScorerResult struct {
	PlayerExternalID int64 // matches players.external_id
	Goals            int
}

// APILeagueFixtureResult is one fixture row from GetLeagueFixtures.
type APILeagueFixtureResult struct {
	ExternalID         int64
	HomeTeamExternalID int64
	AwayTeamExternalID int64
	KickoffAt          time.Time
	StatusShort        string
	Round              string
	GoalsHome          *int
	GoalsAway          *int
}
