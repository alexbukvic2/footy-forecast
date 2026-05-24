package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
)

// FixtureRepository handles persistence of Fixture aggregates.
type FixtureRepository struct {
	pool *db.Pool
	q    *dbgen.Queries
}

// NewFixtureRepository constructs a FixtureRepository backed by pool.
func NewFixtureRepository(pool *db.Pool) *FixtureRepository {
	return &FixtureRepository{pool: pool, q: dbgen.New(pool)}
}

// GetByID fetches a fixture by its UUID.
// Returns domain.ErrNotFound if no row exists.
func (r *FixtureRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Fixture, error) {
	row, err := r.q.GetFixtureByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("fixture %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get fixture: %w", err)
	}
	return fixtureFromRow(row), nil
}

// ListByTournamentForUser returns all fixtures for a tournament paired with the
// requesting user's prediction. Uses two queries merged in Go to avoid nullable
// UUID scan issues with pgx/v5 when sqlc generates non-pointer types.
func (r *FixtureRepository) ListByTournamentForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) ([]*domain.UserFixtureView, error) {
	fixtures, err := r.q.ListFixturesByTournament(ctx, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("list fixtures: %w", err)
	}

	preds, err := r.q.ListPredictionsByUserAndTournament(ctx, dbgen.ListPredictionsByUserAndTournamentParams{
		UserID:       userID,
		TournamentID: tournamentID,
	})
	if err != nil {
		return nil, fmt.Errorf("list predictions: %w", err)
	}

	predByFixture := make(map[uuid.UUID]*domain.ScorePrediction, len(preds))
	for i := range preds {
		p := scorePredictionFromRow(preds[i])
		predByFixture[p.FixtureID] = p
	}

	out := make([]*domain.UserFixtureView, 0, len(fixtures))
	for _, f := range fixtures {
		v := &domain.UserFixtureView{Fixture: *fixtureFromRow(f)}
		if p, ok := predByFixture[f.ID]; ok {
			v.Prediction = p
		}
		out = append(out, v)
	}
	return out, nil
}

// ListLockedByLeague returns in-progress and finished fixtures for a league,
// each paired with all member predictions ordered requesting-user-first then
// alphabetically by display name.
func (r *FixtureRepository) ListLockedByLeague(
	ctx context.Context,
	leagueID, requestingUserID uuid.UUID,
) ([]*domain.LeagueFixtureView, error) {
	rows, err := r.q.ListLockedFixturesByLeague(ctx, dbgen.ListLockedFixturesByLeagueParams{
		LeagueID:         leagueID,
		RequestingUserID: requestingUserID,
	})
	if err != nil {
		return nil, fmt.Errorf("list locked fixtures: %w", err)
	}

	out := make([]*domain.LeagueFixtureView, 0, len(rows))
	for _, row := range rows {
		preds, err := parseMemberPredictions(row.MemberPredictions)
		if err != nil {
			return nil, fmt.Errorf("parse member predictions for fixture %s: %w", row.ID, err)
		}
		out = append(out, &domain.LeagueFixtureView{
			Fixture:     *fixtureFromLockedRow(row),
			Predictions: preds,
		})
	}
	return out, nil
}

// memberPredictionJSON is the JSON shape produced by the ListLockedFixturesByLeague query.
type memberPredictionJSON struct {
	UserID      uuid.UUID `json:"user_id"`
	DisplayName string    `json:"display_name"`
	GoalsHome   *int      `json:"goals_home"`
	GoalsAway   *int      `json:"goals_away"`
	Points      *int      `json:"points"`
}

func parseMemberPredictions(raw interface{}) ([]domain.LeagueMemberPrediction, error) {
	b, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	var items []memberPredictionJSON
	if err := json.Unmarshal(b, &items); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	out := make([]domain.LeagueMemberPrediction, len(items))
	for i, item := range items {
		out[i] = domain.LeagueMemberPrediction{
			UserID:      item.UserID,
			DisplayName: item.DisplayName,
			GoalsHome:   item.GoalsHome,
			GoalsAway:   item.GoalsAway,
			Points:      item.Points,
		}
	}
	return out, nil
}

func fixtureFromRow(row dbgen.Fixture) *domain.Fixture {
	f := &domain.Fixture{
		ID:           row.ID,
		ExternalID:   row.ExternalID,
		TournamentID: row.TournamentID,
		HomeTeamID:   row.HomeTeamID,
		AwayTeamID:   row.AwayTeamID,
		KickoffAt:    row.KickoffAt,
		Status:       domain.FixtureStatus(row.Status),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	if row.GoalsHome != nil {
		v := int(*row.GoalsHome)
		f.GoalsHome = &v
	}
	if row.GoalsAway != nil {
		v := int(*row.GoalsAway)
		f.GoalsAway = &v
	}
	return f
}

func fixtureFromLockedRow(row dbgen.ListLockedFixturesByLeagueRow) *domain.Fixture {
	f := &domain.Fixture{
		ID:           row.ID,
		ExternalID:   row.ExternalID,
		TournamentID: row.TournamentID,
		HomeTeamID:   row.HomeTeamID,
		AwayTeamID:   row.AwayTeamID,
		KickoffAt:    row.KickoffAt,
		Status:       domain.FixtureStatus(row.Status),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	if row.GoalsHome != nil {
		v := int(*row.GoalsHome)
		f.GoalsHome = &v
	}
	if row.GoalsAway != nil {
		v := int(*row.GoalsAway)
		f.GoalsAway = &v
	}
	return f
}
