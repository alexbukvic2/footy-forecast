package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
)

// fakeService implements handler.TournamentService.
type fakeService struct {
	createFn func(
		context.Context,
		domain.CreateTournamentInput,
	) (*domain.Tournament, error)
	getByIDFn func(
		context.Context,
		uuid.UUID,
	) (*domain.Tournament, error)
	getBySlugFn func(
		context.Context,
		string,
	) (*domain.Tournament, error)
	listFn func(context.Context) ([]*domain.Tournament, error)
}

func (f *fakeService) Create(
	ctx context.Context,
	in domain.CreateTournamentInput,
) (*domain.Tournament, error) {
	return f.createFn(ctx, in)
}
func (f *fakeService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Tournament, error) {
	return f.getByIDFn(ctx, id)
}
func (f *fakeService) GetBySlug(
	ctx context.Context,
	slug string,
) (*domain.Tournament, error) {
	return f.getBySlugFn(ctx, slug)
}
func (f *fakeService) List(ctx context.Context) ([]*domain.Tournament, error) {
	return f.listFn(ctx)
}

// silentLogger returns a logger that discards everything, for use in tests.
func silentLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{Level: slog.LevelError + 1}))
}

func TestTournament_Create(t *testing.T) {
	t.Parallel()

	t.Run(
		"returns 201 with the created tournament", func(t *testing.T) {
			t.Parallel()
			created := &domain.Tournament{
				ID:       uuid.New(),
				Slug:     "world-cup-2026",
				Name:     "FIFA World Cup 2026",
				Status:   domain.TournamentStatusUpcoming,
				StartsAt: time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
				EndsAt:   time.Date(2026, 7, 19, 0, 0, 0, 0, time.UTC),
			}
			svc := &fakeService{
				createFn: func(
					context.Context,
					domain.CreateTournamentInput,
				) (*domain.Tournament, error) {
					return created, nil
				},
			}
			h := handler.NewTournament(silentLogger(), svc)

			body := `{"slug":"world-cup-2026","name":"FIFA World Cup 2026","starts_at":"2026-06-11T00:00:00Z","ends_at":"2026-07-19T00:00:00Z"}`
			req := httptest.NewRequest(http.MethodPost, "/tournaments", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h.Create(rec, req)

			require.Equal(t, http.StatusCreated, rec.Code)
			require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var resp map[string]any
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			require.Equal(t, created.ID.String(), resp["id"])
			require.Equal(t, "world-cup-2026", resp["slug"])
		},
	)

	t.Run(
		"returns 400 on malformed JSON", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournament(silentLogger(), &fakeService{})

			req := httptest.NewRequest(http.MethodPost, "/tournaments", strings.NewReader("not json"))
			rec := httptest.NewRecorder()

			h.Create(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 400 when service returns ErrInvalid", func(t *testing.T) {
			t.Parallel()
			svc := &fakeService{
				createFn: func(
					context.Context,
					domain.CreateTournamentInput,
				) (*domain.Tournament, error) {
					return nil, domain.ErrInvalid
				},
			}
			h := handler.NewTournament(silentLogger(), svc)

			body := `{"slug":"x","name":"","starts_at":"2026-06-11T00:00:00Z","ends_at":"2026-07-19T00:00:00Z"}`
			req := httptest.NewRequest(http.MethodPost, "/tournaments", strings.NewReader(body))
			rec := httptest.NewRecorder()

			h.Create(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 409 when service returns ErrConflict", func(t *testing.T) {
			t.Parallel()
			svc := &fakeService{
				createFn: func(
					context.Context,
					domain.CreateTournamentInput,
				) (*domain.Tournament, error) {
					return nil, domain.ErrConflict
				},
			}
			h := handler.NewTournament(silentLogger(), svc)

			body := `{"slug":"world-cup-2026","name":"FIFA","starts_at":"2026-06-11T00:00:00Z","ends_at":"2026-07-19T00:00:00Z"}`
			req := httptest.NewRequest(http.MethodPost, "/tournaments", strings.NewReader(body))
			rec := httptest.NewRecorder()

			h.Create(rec, req)
			require.Equal(t, http.StatusConflict, rec.Code)
		},
	)
}

func TestTournament_GetByID(t *testing.T) {
	t.Parallel()

	t.Run(
		"returns 200 with the tournament", func(t *testing.T) {
			t.Parallel()
			id := uuid.New()
			found := &domain.Tournament{
				ID:     id,
				Slug:   "world-cup-2026",
				Name:   "FIFA World Cup 2026",
				Status: domain.TournamentStatusUpcoming,
			}
			svc := &fakeService{
				getByIDFn: func(
					_ context.Context,
					gotID uuid.UUID,
				) (*domain.Tournament, error) {
					require.Equal(t, id, gotID)
					return found, nil
				},
			}
			h := handler.NewTournament(silentLogger(), svc)

			req := httptest.NewRequest(http.MethodGet, "/tournaments/"+id.String(), nil)
			req.SetPathValue("id", id.String())
			rec := httptest.NewRecorder()

			h.GetByID(rec, req)
			require.Equal(t, http.StatusOK, rec.Code)
		},
	)

	t.Run(
		"returns 400 on invalid UUID", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournament(silentLogger(), &fakeService{})

			req := httptest.NewRequest(http.MethodGet, "/tournaments/not-a-uuid", nil)
			req.SetPathValue("id", "not-a-uuid")
			rec := httptest.NewRecorder()

			h.GetByID(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 404 when service returns ErrNotFound", func(t *testing.T) {
			t.Parallel()
			svc := &fakeService{
				getByIDFn: func(
					context.Context,
					uuid.UUID,
				) (*domain.Tournament, error) {
					return nil, domain.ErrNotFound
				},
			}
			h := handler.NewTournament(silentLogger(), svc)

			id := uuid.New()
			req := httptest.NewRequest(http.MethodGet, "/tournaments/"+id.String(), nil)
			req.SetPathValue("id", id.String())
			rec := httptest.NewRecorder()

			h.GetByID(rec, req)
			require.Equal(t, http.StatusNotFound, rec.Code)
		},
	)
}

func TestTournament_List(t *testing.T) {
	t.Parallel()

	svc := &fakeService{
		listFn: func(context.Context) ([]*domain.Tournament, error) {
			return []*domain.Tournament{
				{ID: uuid.New(), Slug: "world-cup-2026", Name: "FIFA"},
				{ID: uuid.New(), Slug: "euro-2028", Name: "Euro"},
			}, nil
		},
	}
	h := handler.NewTournament(silentLogger(), svc)

	req := httptest.NewRequest(http.MethodGet, "/tournaments", nil)
	rec := httptest.NewRecorder()

	h.List(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var resp struct {
		Tournaments []map[string]any `json:"tournaments"`
	}
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	require.Len(t, resp.Tournaments, 2)
}
