package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/repository/dbgen"
)

// PlayerPredictionRepository handles persistence of PlayerPrediction aggregates.
type PlayerPredictionRepository struct {
	q *dbgen.Queries
}

// NewPlayerPredictionRepository constructs a PlayerPredictionRepository backed by pool.
func NewPlayerPredictionRepository(pool *db.Pool) *PlayerPredictionRepository {
	return &PlayerPredictionRepository{q: dbgen.New(pool)}
}

// UpsertPlayer inserts or updates a player prediction.
// On conflict (user, tournament, category) it updates pick and updated_at — points excluded.
func (r *PlayerPredictionRepository) UpsertPlayer(
	ctx context.Context,
	in domain.UpsertPlayerPredictionInput,
) (*domain.PlayerPrediction, error) {
	row, err := r.q.UpsertPlayerPrediction(ctx, dbgen.UpsertPlayerPredictionParams{
		UserID:       in.UserID,
		TournamentID: in.TournamentID,
		Category:     dbgen.PlayerHandicapCategory(in.Category),
		Pick:         in.Pick,
	})
	if err != nil {
		return nil, fmt.Errorf("upsert player prediction: %w", err)
	}
	return playerPredictionFromUpsertRow(row), nil
}

// ListPlayersByTournamentForUser returns all player predictions for a user in a tournament.
func (r *PlayerPredictionRepository) ListPlayersByTournamentForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) ([]*domain.PlayerPrediction, error) {
	rows, err := r.q.ListPlayerPredictionsByTournamentForUser(ctx, dbgen.ListPlayerPredictionsByTournamentForUserParams{
		TournamentID: tournamentID,
		UserID:       userID,
	})
	if err != nil {
		return nil, fmt.Errorf("list player predictions: %w", err)
	}
	out := make([]*domain.PlayerPrediction, 0, len(rows))
	for _, row := range rows {
		p := &domain.PlayerPrediction{
			ID:           row.ID,
			UserID:       row.UserID,
			TournamentID: row.TournamentID,
			Category:     domain.PlayerHandicapCategory(row.Category),
			Pick:         row.Pick,
			PickName:     row.PickName,
			CreatedAt:    row.CreatedAt,
			UpdatedAt:    row.UpdatedAt,
		}
		if row.Points != nil {
			v := int(*row.Points)
			p.Points = &v
		}
		out = append(out, p)
	}
	return out, nil
}

// ListPlayersByLeague returns player prediction rows for league members who have predicted.
// Members without a prediction are absent; the service merges with the full member list.
func (r *PlayerPredictionRepository) ListPlayersByLeague(
	ctx context.Context,
	leagueID uuid.UUID,
) ([]*domain.PlayerLeaguePick, error) {
	rows, err := r.q.ListPlayerPredictionsByLeague(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list player predictions by league: %w", err)
	}
	out := make([]*domain.PlayerLeaguePick, 0, len(rows))
	for _, row := range rows {
		p := &domain.PlayerLeaguePick{
			UserID:     row.UserID,
			Category:   domain.PlayerHandicapCategory(row.Category),
			PlayerID:   row.PlayerID,
			PlayerName: row.PlayerName,
		}
		if row.Points != nil {
			v := int(*row.Points)
			p.Points = &v
		}
		out = append(out, p)
	}
	return out, nil
}

func playerPredictionFromUpsertRow(row dbgen.UpsertPlayerPredictionRow) *domain.PlayerPrediction {
	p := &domain.PlayerPrediction{
		ID:           row.ID,
		UserID:       row.UserID,
		TournamentID: row.TournamentID,
		Category:     domain.PlayerHandicapCategory(row.Category),
		Pick:         row.Pick,
		PickName:     row.PickName,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	if row.Points != nil {
		v := int(*row.Points)
		p.Points = &v
	}
	return p
}

// TeamPredictionRepository handles persistence of TeamPrediction aggregates.
type TeamPredictionRepository struct {
	q *dbgen.Queries
}

// NewTeamPredictionRepository constructs a TeamPredictionRepository backed by pool.
func NewTeamPredictionRepository(pool *db.Pool) *TeamPredictionRepository {
	return &TeamPredictionRepository{q: dbgen.New(pool)}
}

// UpsertTeam inserts or updates a team prediction.
// On conflict (user, tournament, category) it updates pick and updated_at — points excluded.
func (r *TeamPredictionRepository) UpsertTeam(
	ctx context.Context,
	in domain.UpsertTeamPredictionInput,
) (*domain.TeamPrediction, error) {
	row, err := r.q.UpsertTeamPrediction(ctx, dbgen.UpsertTeamPredictionParams{
		UserID:       in.UserID,
		TournamentID: in.TournamentID,
		Category:     dbgen.TeamHandicapCategory(in.Category),
		Pick:         in.Pick,
	})
	if err != nil {
		return nil, fmt.Errorf("upsert team prediction: %w", err)
	}
	return teamPredictionFromUpsertRow(row), nil
}

// ListTeamsByTournamentForUser returns all team predictions for a user in a tournament.
func (r *TeamPredictionRepository) ListTeamsByTournamentForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) ([]*domain.TeamPrediction, error) {
	rows, err := r.q.ListTeamPredictionsByTournamentForUser(ctx, dbgen.ListTeamPredictionsByTournamentForUserParams{
		TournamentID: tournamentID,
		UserID:       userID,
	})
	if err != nil {
		return nil, fmt.Errorf("list team predictions: %w", err)
	}
	out := make([]*domain.TeamPrediction, 0, len(rows))
	for _, row := range rows {
		p := &domain.TeamPrediction{
			ID:           row.ID,
			UserID:       row.UserID,
			TournamentID: row.TournamentID,
			Category:     domain.TeamHandicapCategory(row.Category),
			Pick:         row.Pick,
			PickName:     row.PickName,
			CreatedAt:    row.CreatedAt,
			UpdatedAt:    row.UpdatedAt,
		}
		if row.Points != nil {
			v := int(*row.Points)
			p.Points = &v
		}
		out = append(out, p)
	}
	return out, nil
}

// ListTeamsByLeague returns team prediction rows for league members who have predicted.
// Members without a prediction are absent; the service merges with the full member list.
func (r *TeamPredictionRepository) ListTeamsByLeague(
	ctx context.Context,
	leagueID uuid.UUID,
) ([]*domain.TeamLeaguePick, error) {
	rows, err := r.q.ListTeamPredictionsByLeague(ctx, leagueID)
	if err != nil {
		return nil, fmt.Errorf("list team predictions by league: %w", err)
	}
	out := make([]*domain.TeamLeaguePick, 0, len(rows))
	for _, row := range rows {
		p := &domain.TeamLeaguePick{
			UserID:   row.UserID,
			Category: domain.TeamHandicapCategory(row.Category),
			TeamID:   row.TeamID,
			TeamName: row.TeamName,
		}
		if row.Points != nil {
			v := int(*row.Points)
			p.Points = &v
		}
		out = append(out, p)
	}
	return out, nil
}

func teamPredictionFromUpsertRow(row dbgen.UpsertTeamPredictionRow) *domain.TeamPrediction {
	p := &domain.TeamPrediction{
		ID:           row.ID,
		UserID:       row.UserID,
		TournamentID: row.TournamentID,
		Category:     domain.TeamHandicapCategory(row.Category),
		Pick:         row.Pick,
		PickName:     row.PickName,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	if row.Points != nil {
		v := int(*row.Points)
		p.Points = &v
	}
	return p
}
