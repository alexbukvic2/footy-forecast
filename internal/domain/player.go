package domain

import "github.com/google/uuid"

// Player represents a football player scoped to a tournament.
type Player struct {
	ID           uuid.UUID
	Name         string
	TournamentID uuid.UUID
	TeamID       uuid.UUID
	GroupLetter  *string
}

// PlayerSearchResult carries player data together with the joined team columns
// returned by the search query.
type PlayerSearchResult struct {
	ID           uuid.UUID
	Name         string
	TournamentID uuid.UUID
	TeamName     string
	TeamLogo     string
	GroupLetter  *string
	Handicaps    map[PlayerHandicapCategory]int
}

// SearchPlayersInput carries the validated inputs for a player search.
type SearchPlayersInput struct {
	TournamentID uuid.UUID
	Query        string
	GroupLetter  *string // optional; nil means search all groups
	HasHandicap  bool    // when true, only return players with at least one handicap row
}
