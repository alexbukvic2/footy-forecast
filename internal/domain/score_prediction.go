package domain

import (
	"time"

	"github.com/google/uuid"
)

// ScorePrediction is a user's exact-score prediction for a fixture.
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

// UpsertScorePredictionInput carries the caller-supplied fields for upserting a score prediction.
type UpsertScorePredictionInput struct {
	UserID    uuid.UUID
	FixtureID uuid.UUID
	GoalsHome int
	GoalsAway int
}
