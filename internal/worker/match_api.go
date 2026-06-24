package worker

import (
	"context"
	"time"
)

// MatchAPI abstracts the external football data provider.
type MatchAPI interface {
	GetFixture(ctx context.Context, externalFixtureID int64) (APIFixtureResult, error)
	GetLeagueFixtures(ctx context.Context, externalLeagueID int64, season int) ([]APILeagueFixtureResult, error)
}

// APIFixtureResult is the data returned by GetFixture.
type APIFixtureResult struct {
	ExternalID            int64
	StatusShort           string // raw API status: "1H", "HT", "FT", "CANC", etc.
	GoalsHome             *int
	GoalsAway             *int
	HomeWinner            *bool // teams.home.winner; handles ET/PEN
	AwayWinner            *bool
	GoalScorerExternalIDs []int64 // one entry per goal scored, excluding own goals
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
