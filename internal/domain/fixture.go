package domain

import (
	"time"

	"github.com/google/uuid"
)

// FixtureStatus represents the current state of a fixture.
type FixtureStatus string

// FixtureStatus values.
const (
	FixtureStatusUpcoming   FixtureStatus = "upcoming"
	FixtureStatusInProgress FixtureStatus = "in_progress"
	FixtureStatusFinished   FixtureStatus = "finished"
)

// Fixture represents a scheduled or completed football match.
type Fixture struct {
	ID               uuid.UUID
	ExternalID       int64
	TournamentID     uuid.UUID
	HomeTeamID       uuid.UUID
	AwayTeamID       uuid.UUID
	HomeTeamName     string
	AwayTeamName     string
	Round            string
	KickoffAt        time.Time
	Status           FixtureStatus
	PredictionLocked bool
	GoalsHome        *int
	GoalsAway        *int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// UserFixtureView pairs a fixture with the requesting user's score prediction (nil if none).
type UserFixtureView struct {
	Fixture    Fixture
	Prediction *ScorePrediction
}

// LeagueMemberPrediction is one member's score prediction for a fixture.
type LeagueMemberPrediction struct {
	UserID      uuid.UUID
	DisplayName string
	GoalsHome   *int
	GoalsAway   *int
	Points      *int
}

// LeagueFixtureView pairs a locked fixture with all league member predictions.
type LeagueFixtureView struct {
	Fixture     Fixture
	Predictions []LeagueMemberPrediction
}
