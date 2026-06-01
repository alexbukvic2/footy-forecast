package domain

import (
	"time"

	"github.com/google/uuid"
)

// PlayerOutcome records which player won a category in a tournament.
type PlayerOutcome struct {
	ID           uuid.UUID
	TournamentID uuid.UUID
	Category     PlayerHandicapCategory
	PlayerID     uuid.UUID
	PlayerName   string
	TeamID       uuid.UUID
	TeamName     string
	RecordedAt   time.Time
}

// TeamOutcome records which team won a category in a tournament.
type TeamOutcome struct {
	ID           uuid.UUID
	TournamentID uuid.UUID
	Category     TeamHandicapCategory
	TeamID       uuid.UUID
	TeamName     string
	RecordedAt   time.Time
}

// TournamentOutcomes bundles all recorded outcomes for a tournament.
type TournamentOutcomes struct {
	PlayerOutcomes []*PlayerOutcome
	TeamOutcomes   []*TeamOutcome
}
