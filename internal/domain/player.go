package domain

import "github.com/google/uuid"

// Player represents a football player scoped to a tournament.
type Player struct {
	ID           uuid.UUID
	Name         string
	TournamentID uuid.UUID
	TeamID       uuid.UUID
}

// PlayerSearchResult carries player data together with the joined team columns
// returned by the search query.
type PlayerSearchResult struct {
	ID           uuid.UUID
	Name         string
	TournamentID uuid.UUID
	TeamName     string
	TeamLogo     string
}

// SearchPlayersInput carries the validated inputs for a player search.
type SearchPlayersInput struct {
	TournamentID uuid.UUID
	Query        string
}
