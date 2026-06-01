package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
)

// ---------- fake service ----------

type fakeTournamentGroupTableSvc struct {
	listFn func(context.Context, uuid.UUID) ([]*domain.TournamentGroupEntry, error)
}

func (f *fakeTournamentGroupTableSvc) ListByTournament(ctx context.Context, id uuid.UUID) ([]*domain.TournamentGroupEntry, error) {
	return f.listFn(ctx, id)
}

// ---------- tests ----------

func TestTournamentGroupTable_ListGroupTable(t *testing.T) {
	tournamentID := uuid.Must(uuid.NewV7())

	t.Run("happy path returns 200 with entries", func(t *testing.T) {
		svc := &fakeTournamentGroupTableSvc{
			listFn: func(_ context.Context, _ uuid.UUID) ([]*domain.TournamentGroupEntry, error) {
				return []*domain.TournamentGroupEntry{
					{
						ID:           uuid.Must(uuid.NewV7()),
						TournamentID: tournamentID,
						TeamID:       uuid.Must(uuid.NewV7()),
						TeamName:     "Argentina",
						GroupLetter:  "A",
						Position:     1,
						Points:       9,
					},
				}, nil
			},
		}
		h := handler.NewTournamentGroupTable(silentLogger(), svc)
		rec := getWithPathValue(t, h.ListGroupTable, "/tournaments/"+tournamentID.String()+"/group-table", "tournamentId", tournamentID.String())

		require.Equal(t, http.StatusOK, rec.Code)
		var body []map[string]any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		require.Len(t, body, 1)
		require.Equal(t, "Argentina", body[0]["team_name"])
		require.Equal(t, "A", body[0]["group_letter"])
		require.Equal(t, float64(1), body[0]["position"])
		require.Equal(t, float64(9), body[0]["points"])
	})

	t.Run("empty table returns 200 empty array", func(t *testing.T) {
		svc := &fakeTournamentGroupTableSvc{
			listFn: func(_ context.Context, _ uuid.UUID) ([]*domain.TournamentGroupEntry, error) {
				return []*domain.TournamentGroupEntry{}, nil
			},
		}
		h := handler.NewTournamentGroupTable(silentLogger(), svc)
		rec := getWithPathValue(t, h.ListGroupTable, "/tournaments/"+tournamentID.String()+"/group-table", "tournamentId", tournamentID.String())

		require.Equal(t, http.StatusOK, rec.Code)
		var body []any
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		require.Empty(t, body)
	})

	t.Run("bad tournamentId returns 400", func(t *testing.T) {
		svc := &fakeTournamentGroupTableSvc{
			listFn: func(_ context.Context, _ uuid.UUID) ([]*domain.TournamentGroupEntry, error) {
				return nil, nil
			},
		}
		h := handler.NewTournamentGroupTable(silentLogger(), svc)
		req := httptest.NewRequest(http.MethodGet, "/tournaments/not-a-uuid/group-table", nil)
		req.SetPathValue("tournamentId", "not-a-uuid")
		rec := httptest.NewRecorder()
		h.ListGroupTable(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("tournament not found returns 404", func(t *testing.T) {
		svc := &fakeTournamentGroupTableSvc{
			listFn: func(_ context.Context, _ uuid.UUID) ([]*domain.TournamentGroupEntry, error) {
				return nil, domain.ErrNotFound
			},
		}
		h := handler.NewTournamentGroupTable(silentLogger(), svc)
		rec := getWithPathValue(t, h.ListGroupTable, "/tournaments/"+tournamentID.String()+"/group-table", "tournamentId", tournamentID.String())

		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("service error returns 500", func(t *testing.T) {
		svc := &fakeTournamentGroupTableSvc{
			listFn: func(_ context.Context, _ uuid.UUID) ([]*domain.TournamentGroupEntry, error) {
				return nil, errors.New("db down")
			},
		}
		h := handler.NewTournamentGroupTable(silentLogger(), svc)
		rec := getWithPathValue(t, h.ListGroupTable, "/tournaments/"+tournamentID.String()+"/group-table", "tournamentId", tournamentID.String())

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
