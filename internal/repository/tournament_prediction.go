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
// Routes to the group-scoped query when GroupLetter is set, otherwise the no-group query.
func (r *PlayerPredictionRepository) UpsertPlayer(
	ctx context.Context,
	in domain.UpsertPlayerPredictionInput,
) (*domain.PlayerPrediction, error) {
	if in.GroupLetter != nil {
		row, err := r.q.UpsertPlayerPredictionGroup(ctx, dbgen.UpsertPlayerPredictionGroupParams{
			UserID:       in.UserID,
			TournamentID: in.TournamentID,
			Category:     dbgen.PlayerHandicapCategory(in.Category),
			Pick:         in.Pick,
			GroupLetter:  in.GroupLetter,
		})
		if err != nil {
			if isUniqueViolation(err) {
				return nil, fmt.Errorf("slot already taken: %w", domain.ErrConflict)
			}
			return nil, fmt.Errorf("upsert player prediction: %w", err)
		}
		return playerPredictionFromGroupRow(row), nil
	}
	row, err := r.q.UpsertPlayerPredictionNoGroup(ctx, dbgen.UpsertPlayerPredictionNoGroupParams{
		UserID:       in.UserID,
		TournamentID: in.TournamentID,
		Category:     dbgen.PlayerHandicapCategory(in.Category),
		Pick:         in.Pick,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("slot already taken: %w", domain.ErrConflict)
		}
		return nil, fmt.Errorf("upsert player prediction: %w", err)
	}
	return playerPredictionFromNoGroupRow(row), nil
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
			GroupLetter:  row.GroupLetter,
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
			UserID:      row.UserID,
			Category:    domain.PlayerHandicapCategory(row.Category),
			GroupLetter: row.GroupLetter,
			PlayerID:    row.PlayerID,
			PlayerName:  row.PlayerName,
		}
		if row.Points != nil {
			v := int(*row.Points)
			p.Points = &v
		}
		out = append(out, p)
	}
	return out, nil
}

func playerPredictionFromGroupRow(row dbgen.UpsertPlayerPredictionGroupRow) *domain.PlayerPrediction {
	p := &domain.PlayerPrediction{
		ID:           row.ID,
		UserID:       row.UserID,
		TournamentID: row.TournamentID,
		Category:     domain.PlayerHandicapCategory(row.Category),
		Pick:         row.Pick,
		PickName:     row.PickName,
		GroupLetter:  row.GroupLetter,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	if row.Points != nil {
		v := int(*row.Points)
		p.Points = &v
	}
	return p
}

func playerPredictionFromNoGroupRow(row dbgen.UpsertPlayerPredictionNoGroupRow) *domain.PlayerPrediction {
	p := &domain.PlayerPrediction{
		ID:           row.ID,
		UserID:       row.UserID,
		TournamentID: row.TournamentID,
		Category:     domain.PlayerHandicapCategory(row.Category),
		Pick:         row.Pick,
		PickName:     row.PickName,
		GroupLetter:  row.GroupLetter,
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
// Routes to the group-scoped query when GroupLetter is set, otherwise the no-group query.
// Returns domain.ErrConflict if the unique slot is already taken.
func (r *TeamPredictionRepository) UpsertTeam(
	ctx context.Context,
	in domain.UpsertTeamPredictionInput,
) (*domain.TeamPrediction, error) {
	if in.GroupLetter != nil {
		row, err := r.q.UpsertTeamPredictionGroup(ctx, dbgen.UpsertTeamPredictionGroupParams{
			UserID:       in.UserID,
			TournamentID: in.TournamentID,
			Category:     dbgen.TeamHandicapCategory(in.Category),
			Pick:         in.Pick,
			GroupLetter:  in.GroupLetter,
			SlotIndex:    int16(in.SlotIndex), //nolint:gosec
		})
		if err != nil {
			if isUniqueViolation(err) {
				return nil, fmt.Errorf("slot already taken: %w", domain.ErrConflict)
			}
			return nil, fmt.Errorf("upsert team prediction: %w", err)
		}
		return teamPredictionFromGroupRow(row), nil
	}
	row, err := r.q.UpsertTeamPredictionNoGroup(ctx, dbgen.UpsertTeamPredictionNoGroupParams{
		UserID:       in.UserID,
		TournamentID: in.TournamentID,
		Category:     dbgen.TeamHandicapCategory(in.Category),
		Pick:         in.Pick,
		SlotIndex:    int16(in.SlotIndex), //nolint:gosec
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("slot already taken: %w", domain.ErrConflict)
		}
		return nil, fmt.Errorf("upsert team prediction: %w", err)
	}
	return teamPredictionFromNoGroupRow(row), nil
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
			GroupLetter:  row.GroupLetter,
			SlotIndex:    int(row.SlotIndex),
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
			UserID:      row.UserID,
			Category:    domain.TeamHandicapCategory(row.Category),
			GroupLetter: row.GroupLetter,
			SlotIndex:   int(row.SlotIndex),
			TeamID:      row.TeamID,
			TeamName:    row.TeamName,
		}
		if row.Points != nil {
			v := int(*row.Points)
			p.Points = &v
		}
		out = append(out, p)
	}
	return out, nil
}

// CountPlayoffWildcards returns the number of playoff wildcard (slot_index=2) picks
// a user has for a given tournament.
func (r *TeamPredictionRepository) CountPlayoffWildcards(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) (int, error) {
	n, err := r.q.CountPlayoffWildcards(ctx, dbgen.CountPlayoffWildcardsParams{
		TournamentID: tournamentID,
		UserID:       userID,
	})
	if err != nil {
		return 0, fmt.Errorf("count playoff wildcards: %w", err)
	}
	return int(n), nil
}

func teamPredictionFromGroupRow(row dbgen.UpsertTeamPredictionGroupRow) *domain.TeamPrediction {
	p := &domain.TeamPrediction{
		ID:           row.ID,
		UserID:       row.UserID,
		TournamentID: row.TournamentID,
		Category:     domain.TeamHandicapCategory(row.Category),
		Pick:         row.Pick,
		PickName:     row.PickName,
		GroupLetter:  row.GroupLetter,
		SlotIndex:    int(row.SlotIndex),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	if row.Points != nil {
		v := int(*row.Points)
		p.Points = &v
	}
	return p
}

func teamPredictionFromNoGroupRow(row dbgen.UpsertTeamPredictionNoGroupRow) *domain.TeamPrediction {
	p := &domain.TeamPrediction{
		ID:           row.ID,
		UserID:       row.UserID,
		TournamentID: row.TournamentID,
		Category:     domain.TeamHandicapCategory(row.Category),
		Pick:         row.Pick,
		PickName:     row.PickName,
		GroupLetter:  row.GroupLetter,
		SlotIndex:    int(row.SlotIndex),
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	if row.Points != nil {
		v := int(*row.Points)
		p.Points = &v
	}
	return p
}
