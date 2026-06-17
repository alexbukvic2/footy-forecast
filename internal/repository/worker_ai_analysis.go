package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
)

// ListLeaguesForFixture returns IDs of all leagues whose tournament contains the fixture.
func (r *WorkerRepository) ListLeaguesForFixture(
	ctx context.Context,
	fixtureID uuid.UUID,
) ([]uuid.UUID, error) {
	ids, err := r.q.ListLeaguesForFixture(ctx, fixtureID)
	if err != nil {
		return nil, fmt.Errorf("list leagues for fixture: %w", err)
	}
	return ids, nil
}

// GetFixtureAnalysisInput returns fixture + league-member prediction data needed for the AI prompt.
func (r *WorkerRepository) GetFixtureAnalysisInput(
	ctx context.Context,
	fixtureID uuid.UUID,
	leagueID uuid.UUID,
) (domain.FixtureAnalysisInput, error) {
	row, err := r.q.GetFixtureAnalysisInput(ctx, dbgen.GetFixtureAnalysisInputParams{
		FixtureID: fixtureID,
		LeagueID:  leagueID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.FixtureAnalysisInput{}, fmt.Errorf("fixture %s league %s: %w", fixtureID, leagueID, domain.ErrNotFound)
		}
		return domain.FixtureAnalysisInput{}, fmt.Errorf("get fixture analysis input: %w", err)
	}

	var rawPreds []workerAnalysisPredictionJSON
	if err := json.Unmarshal([]byte(row.Predictions), &rawPreds); err != nil {
		return domain.FixtureAnalysisInput{}, fmt.Errorf("unmarshal predictions for fixture %s: %w", fixtureID, err)
	}

	preds := make([]domain.AnalysisPrediction, 0, len(rawPreds))
	for _, p := range rawPreds {
		preds = append(preds, domain.AnalysisPrediction{
			DisplayName: p.DisplayName,
			GoalsHome:   p.GoalsHome,
			GoalsAway:   p.GoalsAway,
			Points:      p.Points,
		})
	}

	input := domain.FixtureAnalysisInput{
		HomeTeamName: row.HomeTeamName,
		AwayTeamName: row.AwayTeamName,
		Round:        row.Round,
		GroupLetter:  row.GroupLetter,
		Predictions:  preds,
	}
	if row.GoalsHome != nil {
		v := int(*row.GoalsHome)
		input.GoalsHome = &v
	}
	if row.GoalsAway != nil {
		v := int(*row.GoalsAway)
		input.GoalsAway = &v
	}

	return input, nil
}

// UpsertFixtureAnalysis inserts or updates the AI analysis for a fixture+league pair.
func (r *WorkerRepository) UpsertFixtureAnalysis(
	ctx context.Context,
	fixtureID uuid.UUID,
	leagueID uuid.UUID,
	analysis string,
) error {
	if err := r.q.UpsertFixtureAnalysis(ctx, dbgen.UpsertFixtureAnalysisParams{
		FixtureID: fixtureID,
		LeagueID:  leagueID,
		Analysis:  analysis,
	}); err != nil {
		return fmt.Errorf("upsert fixture analysis: %w", err)
	}
	return nil
}

type workerAnalysisPredictionJSON struct {
	DisplayName string `json:"display_name"`
	GoalsHome   *int   `json:"goals_home"`
	GoalsAway   *int   `json:"goals_away"`
	Points      *int   `json:"points"`
}
