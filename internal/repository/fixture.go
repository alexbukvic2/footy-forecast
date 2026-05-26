package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
)

// FixtureRepository handles persistence of Fixture data.
type FixtureRepository struct {
	q *dbgen.Queries
}

// NewFixtureRepository constructs a FixtureRepository backed by pool.
func NewFixtureRepository(pool *db.Pool) *FixtureRepository {
	return &FixtureRepository{q: dbgen.New(pool)}
}

// GetFirstKickoffByTournament returns the earliest kickoff_at across all fixtures
// for the given tournament. Returns domain.ErrNotFound if no fixtures exist.
func (r *FixtureRepository) GetFirstKickoffByTournament(
	ctx context.Context,
	tournamentID uuid.UUID,
) (time.Time, error) {
	t, err := r.q.GetFirstKickoffByTournament(ctx, tournamentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return time.Time{}, fmt.Errorf("no fixtures for tournament %s: %w", tournamentID, domain.ErrNotFound)
		}
		return time.Time{}, fmt.Errorf("get first kickoff: %w", err)
	}
	return t, nil
}

// GetByID fetches a single fixture by ID.
// Returns domain.ErrNotFound if no row exists.
func (r *FixtureRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Fixture, error) {
	row, err := r.q.GetFixtureByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("fixture %s: %w", id, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("get fixture: %w", err)
	}
	return fixtureFromModel(row), nil
}

// ListByTournamentForUser returns all fixtures for a tournament paired with the user's predictions.
func (r *FixtureRepository) ListByTournamentForUser(ctx context.Context, tournamentID, userID uuid.UUID) ([]*domain.UserFixtureView, error) {
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
	for _, p := range preds {
		sp := scorePredictionFromModel(p)
		predByFixture[sp.FixtureID] = sp
	}

	out := make([]*domain.UserFixtureView, 0, len(fixtures))
	for _, f := range fixtures {
		fix := fixtureFromListRow(f)
		view := &domain.UserFixtureView{
			Fixture:    *fix,
			Prediction: predByFixture[f.ID],
		}
		out = append(out, view)
	}
	return out, nil
}

type memberPredictionJSON struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	GoalsHome   *int   `json:"goals_home"`
	GoalsAway   *int   `json:"goals_away"`
	Points      *int   `json:"points"`
}

// ListLockedByLeague returns all locked (in_progress or finished) fixtures for a league,
// along with each member's prediction. Only fixtures whose kickoff has passed are returned.
func (r *FixtureRepository) ListLockedByLeague(ctx context.Context, leagueID, requestingUserID uuid.UUID) ([]*domain.LeagueFixtureView, error) {
	rows, err := r.q.ListLockedFixturesByLeague(ctx, dbgen.ListLockedFixturesByLeagueParams{
		LeagueID:         leagueID,
		RequestingUserID: requestingUserID,
	})
	if err != nil {
		return nil, fmt.Errorf("list locked fixtures: %w", err)
	}

	out := make([]*domain.LeagueFixtureView, 0, len(rows))
	for _, row := range rows {
		fix := fixtureFromLockedRow(row)

		var members []memberPredictionJSON
		if err := json.Unmarshal([]byte(row.MemberPredictions), &members); err != nil {
			return nil, fmt.Errorf("unmarshal member predictions for fixture %s: %w", row.ID, err)
		}

		preds := make([]domain.LeagueMemberPrediction, 0, len(members))
		for _, m := range members {
			uid, err := uuid.Parse(m.UserID)
			if err != nil {
				return nil, fmt.Errorf("parse user_id in member predictions: %w", err)
			}
			preds = append(preds, domain.LeagueMemberPrediction{
				UserID:      uid,
				DisplayName: m.DisplayName,
				GoalsHome:   m.GoalsHome,
				GoalsAway:   m.GoalsAway,
				Points:      m.Points,
			})
		}

		out = append(out, &domain.LeagueFixtureView{
			Fixture:     *fix,
			Predictions: preds,
		})
	}
	return out, nil
}

func fixtureFromModel(f dbgen.GetFixtureByIDRow) *domain.Fixture {
	fix := &domain.Fixture{
		ID:               f.ID,
		ExternalID:       f.ExternalID,
		TournamentID:     f.TournamentID,
		HomeTeamID:       f.HomeTeamID,
		AwayTeamID:       f.AwayTeamID,
		Round:            f.Round,
		KickoffAt:        f.KickoffAt,
		Status:           domain.FixtureStatus(f.Status),
		PredictionLocked: f.PredictionLocked,
		CreatedAt:        f.CreatedAt,
		UpdatedAt:        f.UpdatedAt,
	}
	if f.GoalsHome != nil {
		v := int(*f.GoalsHome)
		fix.GoalsHome = &v
	}
	if f.GoalsAway != nil {
		v := int(*f.GoalsAway)
		fix.GoalsAway = &v
	}
	return fix
}

func fixtureFromListRow(f dbgen.ListFixturesByTournamentRow) *domain.Fixture {
	fix := &domain.Fixture{
		ID:               f.ID,
		ExternalID:       f.ExternalID,
		TournamentID:     f.TournamentID,
		HomeTeamID:       f.HomeTeamID,
		AwayTeamID:       f.AwayTeamID,
		HomeTeamName:     f.HomeTeamName,
		AwayTeamName:     f.AwayTeamName,
		Round:            f.Round,
		KickoffAt:        f.KickoffAt,
		Status:           domain.FixtureStatus(f.Status),
		PredictionLocked: f.PredictionLocked,
		CreatedAt:        f.CreatedAt,
		UpdatedAt:        f.UpdatedAt,
	}
	if f.GoalsHome != nil {
		v := int(*f.GoalsHome)
		fix.GoalsHome = &v
	}
	if f.GoalsAway != nil {
		v := int(*f.GoalsAway)
		fix.GoalsAway = &v
	}
	return fix
}

func fixtureFromLockedRow(f dbgen.ListLockedFixturesByLeagueRow) *domain.Fixture {
	fix := &domain.Fixture{
		ID:               f.ID,
		ExternalID:       f.ExternalID,
		TournamentID:     f.TournamentID,
		HomeTeamID:       f.HomeTeamID,
		AwayTeamID:       f.AwayTeamID,
		HomeTeamName:     f.HomeTeamName,
		AwayTeamName:     f.AwayTeamName,
		Round:            f.Round,
		KickoffAt:        f.KickoffAt,
		Status:           domain.FixtureStatus(f.Status),
		PredictionLocked: f.PredictionLocked,
		CreatedAt:        f.CreatedAt,
		UpdatedAt:        f.UpdatedAt,
	}
	if f.GoalsHome != nil {
		v := int(*f.GoalsHome)
		fix.GoalsHome = &v
	}
	if f.GoalsAway != nil {
		v := int(*f.GoalsAway)
		fix.GoalsAway = &v
	}
	return fix
}

func scorePredictionFromModel(sp dbgen.ScorePrediction) *domain.ScorePrediction {
	pred := &domain.ScorePrediction{
		ID:        sp.ID,
		UserID:    sp.UserID,
		FixtureID: sp.FixtureID,
		GoalsHome: int(sp.GoalsHome),
		GoalsAway: int(sp.GoalsAway),
		CreatedAt: sp.CreatedAt,
		UpdatedAt: sp.UpdatedAt,
	}
	if sp.Points != nil {
		v := int(*sp.Points)
		pred.Points = &v
	}
	return pred
}
