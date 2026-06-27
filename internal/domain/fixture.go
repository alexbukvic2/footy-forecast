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
	FixtureStatusCancelled  FixtureStatus = "cancelled"
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
	Group            *string // nil for knockout fixtures
	Round            string
	KickoffAt        time.Time
	Status           FixtureStatus
	PredictionLocked bool
	// GoalsHomeRegular/GoalsAwayRegular reflect 90-minute score only (frozen once ET begins).
	GoalsHomeRegular *int
	GoalsAwayRegular *int
	// GoalsHome/GoalsAway include ET goals (but not penalty shootout).
	GoalsHome    *int
	GoalsAway    *int
	WinnerTeamID *uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	Winner      *uuid.UUID
	Points      *int
}

// LeagueFixtureView pairs a locked fixture with all league member predictions.
type LeagueFixtureView struct {
	Fixture     Fixture
	Predictions []LeagueMemberPrediction
}
