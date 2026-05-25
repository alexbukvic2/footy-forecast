package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
)

// ---------- tournament-specific fakes ----------

// fakeTournamentService implements handler.TournamentService.
type fakeTournamentService struct {
	createFn func(
		context.Context,
		domain.CreateTournamentInput,
	) (*domain.Tournament, error)
	getByIDFn func(
		context.Context,
		uuid.UUID,
	) (*domain.Tournament, error)
	listFn func(context.Context) ([]*domain.Tournament, error)
}

func (f *fakeTournamentService) Create(
	ctx context.Context,
	in domain.CreateTournamentInput,
) (*domain.Tournament, error) {
	return f.createFn(ctx, in)
}
func (f *fakeTournamentService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Tournament, error) {
	return f.getByIDFn(ctx, id)
}
func (f *fakeTournamentService) List(ctx context.Context) ([]*domain.Tournament, error) {
	return f.listFn(ctx)
}

// ---------- tournament-specific fixtures ----------

// createTournamentReq mirrors the JSON shape clients send to POST /tournaments.
type createTournamentReq struct {
	Slug     string    `json:"slug"`
	Name     string    `json:"name"`
	StartsAt time.Time `json:"starts_at"`
	EndsAt   time.Time `json:"ends_at"`
}

// validCreateTournamentReq returns a request that should always be accepted.
// Tests mutate one field at a time to exercise specific paths.
func validCreateTournamentReq() createTournamentReq {
	return createTournamentReq{
		Slug:     "world-cup-2026",
		Name:     "FIFA World Cup 2026",
		StartsAt: time.Date(2026, 6, 11, 0, 0, 0, 0, time.UTC),
		EndsAt:   time.Date(2026, 7, 19, 0, 0, 0, 0, time.UTC),
	}
}

// ---------- tests ----------

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
			svc := &fakeTournamentService{
				createFn: func(
					context.Context,
					domain.CreateTournamentInput,
				) (*domain.Tournament, error) {
					return created, nil
				},
			}
			h := handler.NewTournament(silentLogger(), svc)

			rec := postJSON(t, h.Create, "/tournaments", validCreateTournamentReq())

			require.Equal(t, http.StatusCreated, rec.Code)
			require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var resp map[string]any
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			require.Equal(t, created.ID.String(), resp["id"])
			require.Equal(t, "world-cup-2026", resp["slug"])
		},
	)

	t.Run(
		"passes request through to service", func(t *testing.T) {
			t.Parallel()
			var receivedInput domain.CreateTournamentInput
			svc := &fakeTournamentService{
				createFn: func(
					_ context.Context,
					in domain.CreateTournamentInput,
				) (*domain.Tournament, error) {
					receivedInput = in
					return &domain.Tournament{ID: uuid.New(), Slug: in.Slug, Name: in.Name}, nil
				},
			}
			h := handler.NewTournament(silentLogger(), svc)

			rec := postJSON(t, h.Create, "/tournaments", validCreateTournamentReq())

			require.Equal(t, http.StatusCreated, rec.Code)
			require.Equal(t, "world-cup-2026", receivedInput.Slug)
			require.Equal(t, "FIFA World Cup 2026", receivedInput.Name)
		},
	)

	t.Run(
		"returns 400 on malformed JSON", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournament(silentLogger(), &fakeTournamentService{})
			rec := postRaw(t, h.Create, "/tournaments", "not json")
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 400 when service returns ErrInvalid", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentService{
				createFn: func(
					context.Context,
					domain.CreateTournamentInput,
				) (*domain.Tournament, error) {
					return nil, domain.ErrInvalid
				},
			}
			h := handler.NewTournament(silentLogger(), svc)

			req := validCreateTournamentReq()
			req.Name = ""
			rec := postJSON(t, h.Create, "/tournaments", req)

			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 409 when service returns ErrConflict", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentService{
				createFn: func(
					context.Context,
					domain.CreateTournamentInput,
				) (*domain.Tournament, error) {
					return nil, domain.ErrConflict
				},
			}
			h := handler.NewTournament(silentLogger(), svc)

			rec := postJSON(t, h.Create, "/tournaments", validCreateTournamentReq())

			require.Equal(t, http.StatusConflict, rec.Code)
		},
	)

	t.Run(
		"returns 500 on unexpected service error", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentService{
				createFn: func(
					context.Context,
					domain.CreateTournamentInput,
				) (*domain.Tournament, error) {
					return nil, context.DeadlineExceeded
				},
			}
			h := handler.NewTournament(silentLogger(), svc)

			rec := postJSON(t, h.Create, "/tournaments", validCreateTournamentReq())

			require.Equal(t, http.StatusInternalServerError, rec.Code)
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
			svc := &fakeTournamentService{
				getByIDFn: func(
					_ context.Context,
					gotID uuid.UUID,
				) (*domain.Tournament, error) {
					require.Equal(t, id, gotID)
					return found, nil
				},
			}
			h := handler.NewTournament(silentLogger(), svc)

			rec := getWithPathValue(t, h.GetByID, "/tournaments/"+id.String(), "id", id.String())

			require.Equal(t, http.StatusOK, rec.Code)

			var resp map[string]any
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			require.Equal(t, id.String(), resp["id"])
		},
	)

	t.Run(
		"returns 400 on invalid UUID", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournament(silentLogger(), &fakeTournamentService{})

			rec := getWithPathValue(t, h.GetByID, "/tournaments/not-a-uuid", "id", "not-a-uuid")

			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 404 when service returns ErrNotFound", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentService{
				getByIDFn: func(
					context.Context,
					uuid.UUID,
				) (*domain.Tournament, error) {
					return nil, domain.ErrNotFound
				},
			}
			h := handler.NewTournament(silentLogger(), svc)

			id := uuid.New()
			rec := getWithPathValue(t, h.GetByID, "/tournaments/"+id.String(), "id", id.String())

			require.Equal(t, http.StatusNotFound, rec.Code)
		},
	)
}

func TestTournament_List(t *testing.T) {
	t.Parallel()

	t.Run(
		"returns 200 with all tournaments", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentService{
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
		},
	)

	t.Run(
		"returns empty list when no tournaments", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentService{
				listFn: func(context.Context) ([]*domain.Tournament, error) {
					return []*domain.Tournament{}, nil
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
			require.Empty(t, resp.Tournaments)
		},
	)

	t.Run(
		"returns 500 when service errors", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentService{
				listFn: func(context.Context) ([]*domain.Tournament, error) {
					return nil, context.DeadlineExceeded
				},
			}
			h := handler.NewTournament(silentLogger(), svc)

			req := httptest.NewRequest(http.MethodGet, "/tournaments", nil)
			rec := httptest.NewRecorder()
			h.List(rec, req)

			require.Equal(t, http.StatusInternalServerError, rec.Code)
		},
	)
}
