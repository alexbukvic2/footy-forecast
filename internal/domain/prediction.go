package domain

import (
	"time"

	"github.com/google/uuid"
)

// ScorePrediction is a user's predicted score for a fixture.
type ScorePrediction struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FixtureID uuid.UUID
	GoalsHome int
	GoalsAway int
	Points    *int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UpsertScorePredictionInput carries the fields needed to create or update a score prediction.
type UpsertScorePredictionInput struct {
	UserID    uuid.UUID
	FixtureID uuid.UUID
	GoalsHome int
	GoalsAway int
}
