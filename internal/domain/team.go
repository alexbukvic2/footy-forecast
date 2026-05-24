package domain

import "github.com/google/uuid"

// Team represents a football team scoped to a tournament.
type Team struct {
	ID           uuid.UUID
	Name         string
	Logo         string
	TournamentID uuid.UUID
}
