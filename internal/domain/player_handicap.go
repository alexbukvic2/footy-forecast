package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// PlayerHandicapCategory is the category for which a player handicap is assigned.
type PlayerHandicapCategory string

// Valid PlayerHandicapCategory values.
const (
	PlayerHandicapCategoryGroupTopScorer PlayerHandicapCategory = "group_top_scorer"
	PlayerHandicapCategoryTotalTopScorer PlayerHandicapCategory = "total_top_scorer"
)

// AllPlayerHandicapCategories lists every valid PlayerHandicapCategory.
var AllPlayerHandicapCategories = []PlayerHandicapCategory{
	PlayerHandicapCategoryGroupTopScorer,
	PlayerHandicapCategoryTotalTopScorer,
}

// ParsePlayerHandicapCategory returns a PlayerHandicapCategory for s, or ErrInvalid if unrecognised.
func ParsePlayerHandicapCategory(s string) (PlayerHandicapCategory, error) {
	switch PlayerHandicapCategory(s) {
	case PlayerHandicapCategoryGroupTopScorer, PlayerHandicapCategoryTotalTopScorer:
		return PlayerHandicapCategory(s), nil
	default:
		return "", fmt.Errorf("invalid category: %w", ErrInvalid)
	}
}

// PlayerHandicap holds the point assignment for a player in a given category.
type PlayerHandicap struct {
	ID       uuid.UUID
	PlayerID uuid.UUID
	Category PlayerHandicapCategory
	Points   int
}
