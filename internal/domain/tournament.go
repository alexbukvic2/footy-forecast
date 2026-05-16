package domain

import (
	"time"

	"github.com/google/uuid"
)

// TournamentStatus is the lifecycle state of a tournament.
type TournamentStatus string

// Tournament status values. These mirror the Postgres enum tournament_status.
const (
	TournamentStatusUpcoming   TournamentStatus = "upcoming"
	TournamentStatusInProgress TournamentStatus = "in_progress"
	TournamentStatusConcluded  TournamentStatus = "concluded"
)

// Valid reports whether s is one of the defined statuses.
func (s TournamentStatus) Valid() bool {
	switch s {
	case TournamentStatusUpcoming, TournamentStatusInProgress, TournamentStatusConcluded:
		return true
	}
	return false
}

// Tournament is a competition users can make predictions about (e.g., World Cup 2026).
type Tournament struct {
	ID        uuid.UUID
	Slug      string
	Name      string
	Status    TournamentStatus
	StartsAt  time.Time
	EndsAt    time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
