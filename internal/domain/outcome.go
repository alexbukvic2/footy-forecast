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
	RecordedAt   time.Time
}

// TeamOutcome records which team won a category in a tournament.
type TeamOutcome struct {
	ID           uuid.UUID
	TournamentID uuid.UUID
	Category     TeamHandicapCategory
	TeamID       uuid.UUID
	RecordedAt   time.Time
}
