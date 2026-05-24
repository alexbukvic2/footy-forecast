package domain

import (
	"time"

	"github.com/google/uuid"
)

// FixtureStatus represents the lifecycle state of a fixture.
type FixtureStatus string

// Fixture status values. These mirror the Postgres enum fixture_status.
const (
	FixtureStatusUpcoming   FixtureStatus = "upcoming"
	FixtureStatusInProgress FixtureStatus = "in_progress"
	FixtureStatusFinished   FixtureStatus = "finished"
)

// Fixture is a scheduled match between two teams.
type Fixture struct {
	ID           uuid.UUID
	ExternalID   int64
	TournamentID uuid.UUID
	HomeTeamID   uuid.UUID
	AwayTeamID   uuid.UUID
	KickoffAt    time.Time
	Status       FixtureStatus
	GoalsHome    *int
	GoalsAway    *int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserFixtureView pairs a fixture with the requesting user's prediction (if any).
type UserFixtureView struct {
	Fixture    Fixture
	Prediction *ScorePrediction
}

// LeagueMemberPrediction is one member's prediction within a league fixture view.
type LeagueMemberPrediction struct {
	UserID      uuid.UUID
	DisplayName string
	GoalsHome   *int
	GoalsAway   *int
	Points      *int
}

// LeagueFixtureView pairs a locked fixture with all member predictions.
type LeagueFixtureView struct {
	Fixture     Fixture
	Predictions []LeagueMemberPrediction
}
