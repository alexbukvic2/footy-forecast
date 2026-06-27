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
	GoalsHome            *int // total goals including ET (what the API returns); used for change detection
	GoalsAway            *int
	WinnerTeamID         *uuid.UUID
	LastPolledAt         *time.Time
}

// ActiveTournament holds the identifiers needed to refresh fixtures from the API.
type ActiveTournament struct {
	ID         uuid.UUID
	ExternalID int64
	Season     int
}

// NewFixture is a fixture sourced from the API to be inserted into the DB.
type NewFixture struct {
	ExternalID int64
	HomeTeamID uuid.UUID
	AwayTeamID uuid.UUID
	KickoffAt  time.Time
	Status     FixtureStatus
	Round      string
	GoalsHome  *int
	GoalsAway  *int
}

// StandingsEntry is one team's standings row resolved to our internal team UUID.
type StandingsEntry struct {
	TeamID       uuid.UUID
	Group        *string // nil for non-group-stage standings; "A", "B", … for group stages
	Position     int
	Points       int
	Played       int
	Won          int
	Drawn        int
	Lost         int
	GoalsFor     int
	GoalsAgainst int
	Description  string
}

// GroupFixture is a group-stage fixture whose score is already known.
type GroupFixture struct {
	HomeTeamID uuid.UUID
	AwayTeamID uuid.UUID
	GoalsHome  int
	GoalsAway  int
}
