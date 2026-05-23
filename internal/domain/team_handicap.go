package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// TeamHandicapCategory is the category for which a team handicap is assigned.
type TeamHandicapCategory string

// Valid TeamHandicapCategory values.
const (
	TeamHandicapCategoryGroupWinner  TeamHandicapCategory = "group_winner"
	TeamHandicapCategoryPlayoff      TeamHandicapCategory = "playoff"
	TeamHandicapCategorySemifinalist TeamHandicapCategory = "semifinalist"
	TeamHandicapCategoryWinner       TeamHandicapCategory = "winner"
)

// ParseTeamHandicapCategory returns a TeamHandicapCategory for s, or ErrInvalid if unrecognised.
func ParseTeamHandicapCategory(s string) (TeamHandicapCategory, error) {
	switch TeamHandicapCategory(s) {
	case TeamHandicapCategoryGroupWinner, TeamHandicapCategoryPlayoff,
		TeamHandicapCategorySemifinalist, TeamHandicapCategoryWinner:
		return TeamHandicapCategory(s), nil
	default:
		return "", fmt.Errorf("invalid category: %w", ErrInvalid)
	}
}

// TeamHandicap holds the point assignment for a team in a given category.
type TeamHandicap struct {
	ID       uuid.UUID
	TeamID   uuid.UUID
	Category TeamHandicapCategory
	Points   int
}
