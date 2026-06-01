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
	GroupLetter  *string
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
	GroupLetter  *string
}

// PlayerPredictionView is one category row in the personal listing.
// Prediction is nil when the user has not yet predicted for this category.
type PlayerPredictionView struct {
	Category    PlayerHandicapCategory
	GroupLetter *string
	Prediction  *PlayerPrediction
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
	GroupLetter *string
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
	GroupLetter  *string
	SlotIndex    int
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
	GroupLetter  *string
	SlotIndex    int
}

// TeamPredictionView is one category row in the personal listing.
// Prediction is nil when the user has not yet predicted for this category.
type TeamPredictionView struct {
	Category    TeamHandicapCategory
	GroupLetter *string
	SlotIndex   int
	Prediction  *TeamPrediction
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
	GroupLetter *string
	SlotIndex   int
	Predictions []LeagueMemberTeamPick
}

// PlayerLeaguePick is a raw row from the league player predictions query.
// Used internally by the service to build LeaguePlayerCategoryView.
type PlayerLeaguePick struct {
	UserID      uuid.UUID
	Category    PlayerHandicapCategory
	GroupLetter *string
	PlayerID    uuid.UUID
	PlayerName  string
	Points      *int
}

// TeamLeaguePick is a raw row from the league team predictions query.
// Used internally by the service to build LeagueTeamCategoryView.
type TeamLeaguePick struct {
	UserID      uuid.UUID
	Category    TeamHandicapCategory
	GroupLetter *string
	SlotIndex   int
	TeamID      uuid.UUID
	TeamName    string
	Points      *int
}

// BulkPlayerPredictionItem is one item in a bulk upsert request.
// PlayerID nil means clear (delete) the slot.
type BulkPlayerPredictionItem struct {
	Category    PlayerHandicapCategory
	GroupLetter *string
	PlayerID    *uuid.UUID
}

// BulkTeamPredictionItem is one item in a bulk upsert request.
// TeamID nil means clear (delete) the slot.
type BulkTeamPredictionItem struct {
	Category    TeamHandicapCategory
	GroupLetter *string
	SlotIndex   int
	TeamID      *uuid.UUID
}

// LeagueMemberDisplay carries a league member's user ID and display name.
type LeagueMemberDisplay struct {
	UserID      uuid.UUID
	DisplayName string
}

// LeagueGroupPredictions is the response shape for a specific group's predictions.
type LeagueGroupPredictions struct {
	Group             string
	TeamPredictions   []*LeagueTeamCategoryView
	PlayerPredictions []*LeaguePlayerCategoryView
}

// LeaguePlayoffPredictions is the response shape for knockout/outright predictions.
type LeaguePlayoffPredictions struct {
	TeamPredictions   []*LeagueTeamCategoryView
	PlayerPredictions []*LeaguePlayerCategoryView
}
