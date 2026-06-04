package domain

import (
	"time"

	"github.com/google/uuid"
)

// TournamentGroupEntry is one team's standing in a tournament group.
type TournamentGroupEntry struct {
	ID           uuid.UUID
	TournamentID uuid.UUID
	TeamID       uuid.UUID
	TeamName     string
	GroupLetter  string
	Position     int
	Points       int
	Played       int
	Won          int
	Drawn        int
	Lost         int
	GoalsFor     int
	GoalsAgainst int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
