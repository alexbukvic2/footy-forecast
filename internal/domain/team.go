package domain

import "github.com/google/uuid"

// Team represents a football team scoped to a tournament.
type Team struct {
	ID           uuid.UUID
	Name         string
	ShortName    string
	Logo         string
	TournamentID uuid.UUID
	GroupLetter  *string
}

// TeamHandicapItem is a single handicap entry embedded in a team response.
type TeamHandicapItem struct {
	Category TeamHandicapCategory
	Points   int
}

// TeamWithHandicaps is a Team enriched with all its handicap assignments.
type TeamWithHandicaps struct {
	Team
	Handicaps []TeamHandicapItem
}
