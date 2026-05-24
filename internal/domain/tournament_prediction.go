package domain

import (
	"time"

	"github.com/google/uuid"
)

// PlayerPrediction is a user's pick for a player category in a tournament.
type PlayerPrediction struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	TournamentID uuid.UUID
	Category     PlayerHandicapCategory
	Pick         uuid.UUID
	PickName     string
	Points       *int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UpsertPlayerPredictionInput carries the caller-supplied fields for upserting a player prediction.
type UpsertPlayerPredictionInput struct {
	UserID       uuid.UUID
	TournamentID uuid.UUID
	Category     PlayerHandicapCategory
	Pick         uuid.UUID
}

// PlayerPredictionView is one category row in the personal listing.
// Prediction is nil when the user has not yet predicted for this category.
type PlayerPredictionView struct {
	Category   PlayerHandicapCategory
	Prediction *PlayerPrediction
}

// LeagueMemberPlayerPick is one member's pick within a league category view.
type LeagueMemberPlayerPick struct {
	UserID      uuid.UUID
	DisplayName string
	PlayerID    *uuid.UUID
	PlayerName  *string
	Points      *int
}

// LeaguePlayerCategoryView is one category row in the league listing.
type LeaguePlayerCategoryView struct {
	Category    PlayerHandicapCategory
	Predictions []LeagueMemberPlayerPick
}

// TeamPrediction is a user's pick for a team category in a tournament.
type TeamPrediction struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	TournamentID uuid.UUID
	Category     TeamHandicapCategory
	Pick         uuid.UUID
	PickName     string
	Points       *int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UpsertTeamPredictionInput carries the caller-supplied fields for upserting a team prediction.
type UpsertTeamPredictionInput struct {
	UserID       uuid.UUID
	TournamentID uuid.UUID
	Category     TeamHandicapCategory
	Pick         uuid.UUID
}

// TeamPredictionView is one category row in the personal listing.
// Prediction is nil when the user has not yet predicted for this category.
type TeamPredictionView struct {
	Category   TeamHandicapCategory
	Prediction *TeamPrediction
}

// LeagueMemberTeamPick is one member's pick within a league category view.
type LeagueMemberTeamPick struct {
	UserID      uuid.UUID
	DisplayName string
	TeamID      *uuid.UUID
	TeamName    *string
	Points      *int
}

// LeagueTeamCategoryView is one category row in the league listing.
type LeagueTeamCategoryView struct {
	Category    TeamHandicapCategory
	Predictions []LeagueMemberTeamPick
}

// PlayerLeaguePick is a raw row from the league player predictions query.
// Used internally by the service to build LeaguePlayerCategoryView.
type PlayerLeaguePick struct {
	UserID     uuid.UUID
	Category   PlayerHandicapCategory
	PlayerID   uuid.UUID
	PlayerName string
	Points     *int
}

// TeamLeaguePick is a raw row from the league team predictions query.
// Used internally by the service to build LeagueTeamCategoryView.
type TeamLeaguePick struct {
	UserID   uuid.UUID
	Category TeamHandicapCategory
	TeamID   uuid.UUID
	TeamName string
	Points   *int
}

// LeagueMemberDisplay carries a league member's user ID and display name.
type LeagueMemberDisplay struct {
	UserID      uuid.UUID
	DisplayName string
}
