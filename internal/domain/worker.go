package domain

import (
	"time"

	"github.com/google/uuid"
)

// PollableFixture carries all fields the worker needs to poll and settle a fixture.
type PollableFixture struct {
	ID                   uuid.UUID
	ExternalID           int64
	TournamentID         uuid.UUID
	TournamentExternalID int64
	TournamentSeason     int
	HomeTeamID           uuid.UUID
	AwayTeamID           uuid.UUID
	GroupLetter          *string
	Round                string
	Status               FixtureStatus
	KickoffAt            time.Time
	GoalsHome            *int
	GoalsAway            *int
	WinnerTeamID         *uuid.UUID
	LastPolledAt         *time.Time
}

// StandingsEntry is one team's standings row resolved to our internal team UUID.
type StandingsEntry struct {
	TeamID       uuid.UUID
	Position     int
	Points       int
	Played       int
	Won          int
	Drawn        int
	Lost         int
	GoalsFor     int
	GoalsAgainst int
}
