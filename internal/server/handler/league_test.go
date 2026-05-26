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
	"github.com/alexbukvic2/footy-forecast/internal/server/ctxutil"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
)

// ---------- fake service ----------

type fakeLeagueService struct {
	createLeagueFn       func(context.Context, uuid.UUID, domain.CreateLeagueInput) (*domain.League, error)
	getLeagueFn          func(context.Context, uuid.UUID, uuid.UUID) (*domain.League, []*domain.LeagueMember, error)
	listLeaguesForUserFn func(context.Context, uuid.UUID) ([]*domain.LeagueSummary, error)
	updateLeagueNameFn   func(context.Context, uuid.UUID, uuid.UUID, string) (*domain.League, error)
	deleteLeagueFn       func(context.Context, uuid.UUID, uuid.UUID) error
	regenerateCodeFn     func(context.Context, uuid.UUID, uuid.UUID) (*domain.League, error)
	joinLeagueFn         func(context.Context, string, uuid.UUID) (*domain.League, error)
	removeMemberFn       func(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) error
}

func (f *fakeLeagueService) CreateLeague(ctx context.Context, userID uuid.UUID, in domain.CreateLeagueInput) (*domain.League, error) {
	return f.createLeagueFn(ctx, userID, in)
}
func (f *fakeLeagueService) GetLeague(ctx context.Context, leagueID, requesterID uuid.UUID) (*domain.League, []*domain.LeagueMember, error) {
	return f.getLeagueFn(ctx, leagueID, requesterID)
}
func (f *fakeLeagueService) ListLeaguesForUser(ctx context.Context, userID uuid.UUID) ([]*domain.LeagueSummary, error) {
	return f.listLeaguesForUserFn(ctx, userID)
}
func (f *fakeLeagueService) UpdateLeagueName(ctx context.Context, leagueID, requesterID uuid.UUID, name string) (*domain.League, error) {
	return f.updateLeagueNameFn(ctx, leagueID, requesterID, name)
}
func (f *fakeLeagueService) DeleteLeague(ctx context.Context, leagueID, requesterID uuid.UUID) error {
	return f.deleteLeagueFn(ctx, leagueID, requesterID)
}
func (f *fakeLeagueService) RegenerateCode(ctx context.Context, leagueID, requesterID uuid.UUID) (*domain.League, error) {
	return f.regenerateCodeFn(ctx, leagueID, requesterID)
}
func (f *fakeLeagueService) JoinLeague(ctx context.Context, code string, userID uuid.UUID) (*domain.League, error) {
	return f.joinLeagueFn(ctx, code, userID)
}
func (f *fakeLeagueService) RemoveMember(ctx context.Context, leagueID, targetUserID, requesterID uuid.UUID) error {
	return f.removeMemberFn(ctx, leagueID, targetUserID, requesterID)
}

// ---------- helpers ----------

func authedRequest(r *http.Request) *http.Request {
	u := domain.User{ID: uuid.New(), Email: "test@example.com", Status: domain.UserStatusActive}
	return r.WithContext(ctxutil.WithUser(r.Context(), u))
}

func postJSONAuthed(t *testing.T, h http.HandlerFunc, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, path, mustReader(t, raw))
	req.Header.Set("Content-Type", "application/json")
	req = authedRequest(req)
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

func deleteWithPathValues(t *testing.T, h http.HandlerFunc, path string, kvs ...string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodDelete, path, nil)
	req = authedRequest(req)
	for i := 0; i+1 < len(kvs); i += 2 {
		req.SetPathValue(kvs[i], kvs[i+1])
	}
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

func getAuthedWithPathValue(t *testing.T, h http.HandlerFunc, path, key, value string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.SetPathValue(key, value)
	req = authedRequest(req)
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

func patchJSONAuthedWithPathValue(t *testing.T, h http.HandlerFunc, path, key, value string, body any) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPatch, path, mustReader(t, raw))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(key, value)
	req = authedRequest(req)
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

func postJSONAuthedWithPathValue(t *testing.T, h http.HandlerFunc, path, key, value string, body any) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, path, mustReader(t, raw))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue(key, value)
	req = authedRequest(req)
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

func mustReader(t *testing.T, b []byte) *bytesReader {
	t.Helper()
	return &bytesReader{b: b}
}

// bytesReader wraps []byte as io.Reader for use with httptest.
type bytesReader struct {
	b   []byte
	pos int
}

func (r *bytesReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.b) {
		return 0, nil
	}
	n := copy(p, r.b[r.pos:])
	r.pos += n
	return n, nil
}

// ---------- fixtures ----------

func aLeague() *domain.League {
	return &domain.League{
		ID:           uuid.New(),
		TournamentID: uuid.New(),
		OwnerID:      uuid.New(),
		Name:         "Test League",
		Code:         "TESTCODE",
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
}

func aLeagueMember(leagueID uuid.UUID, role domain.LeagueMemberRole) *domain.LeagueMember {
	return &domain.LeagueMember{
		LeagueID: leagueID,
		UserID:   uuid.New(),
		Role:     role,
		JoinedAt: time.Now().UTC(),
	}
}

// ---------- Create ----------

func TestLeague_Create(t *testing.T) {
	t.Parallel()

	t.Run("returns 201 with league", func(t *testing.T) {
		t.Parallel()
		league := aLeague()
		svc := &fakeLeagueService{
			createLeagueFn: func(context.Context, uuid.UUID, domain.CreateLeagueInput) (*domain.League, error) {
				return league, nil
			},
		}
		h := handler.NewLeague(silentLogger(), svc)

		body := map[string]any{"tournament_id": uuid.New().String(), "name": "My League"}
		rec := postJSONAuthed(t, h.Create, "/leagues", body)

		require.Equal(t, http.StatusCreated, rec.Code)
		var resp map[string]any
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Equal(t, league.ID.String(), resp["id"])
		require.Equal(t, "TESTCODE", resp["code"])
	})

	t.Run("returns 400 on malformed JSON", func(t *testing.T) {
		t.Parallel()
		h := handler.NewLeague(silentLogger(), &fakeLeagueService{})
		req := httptest.NewRequest(http.MethodPost, "/leagues", strReader("not json"))
		req = authedRequest(req)
		rec := httptest.NewRecorder()
		h.Create(rec, req)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 400 on invalid tournament_id", func(t *testing.T) {
		t.Parallel()
		h := handler.NewLeague(silentLogger(), &fakeLeagueService{})
		body := map[string]any{"tournament_id": "not-a-uuid", "name": "My League"}
		rec := postJSONAuthed(t, h.Create, "/leagues", body)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 400 when service returns ErrInvalid", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			createLeagueFn: func(context.Context, uuid.UUID, domain.CreateLeagueInput) (*domain.League, error) {
				return nil, domain.ErrInvalid
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		body := map[string]any{"tournament_id": uuid.New().String(), "name": ""}
		rec := postJSONAuthed(t, h.Create, "/leagues", body)
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 404 when tournament not found", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			createLeagueFn: func(context.Context, uuid.UUID, domain.CreateLeagueInput) (*domain.League, error) {
				return nil, domain.ErrNotFound
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		body := map[string]any{"tournament_id": uuid.New().String(), "name": "My League"}
		rec := postJSONAuthed(t, h.Create, "/leagues", body)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

// ---------- Get ----------

func TestLeague_Get(t *testing.T) {
	t.Parallel()

	t.Run("returns 200 with league and members", func(t *testing.T) {
		t.Parallel()
		league := aLeague()
		members := []*domain.LeagueMember{aLeagueMember(league.ID, domain.LeagueMemberRoleOwner)}
		svc := &fakeLeagueService{
			getLeagueFn: func(context.Context, uuid.UUID, uuid.UUID) (*domain.League, []*domain.LeagueMember, error) {
				return league, members, nil
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		rec := getAuthedWithPathValue(t, h.Get, "/leagues/"+league.ID.String(), "id", league.ID.String())

		require.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]any
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Equal(t, league.ID.String(), resp["id"])
		require.NotNil(t, resp["members"])
	})

	t.Run("returns 404 for non-member", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			getLeagueFn: func(context.Context, uuid.UUID, uuid.UUID) (*domain.League, []*domain.LeagueMember, error) {
				return nil, nil, domain.ErrNotFound
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		id := uuid.New()
		rec := getAuthedWithPathValue(t, h.Get, "/leagues/"+id.String(), "id", id.String())
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("returns 400 on invalid id", func(t *testing.T) {
		t.Parallel()
		h := handler.NewLeague(silentLogger(), &fakeLeagueService{})
		rec := getAuthedWithPathValue(t, h.Get, "/leagues/bad-uuid", "id", "bad-uuid")
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// ---------- List ----------

func TestLeague_List(t *testing.T) {
	t.Parallel()

	t.Run("returns 200 with leagues and my_position", func(t *testing.T) {
		t.Parallel()
		l1 := aLeague()
		l2 := aLeague()
		svc := &fakeLeagueService{
			listLeaguesForUserFn: func(context.Context, uuid.UUID) ([]*domain.LeagueSummary, error) {
				return []*domain.LeagueSummary{
					{League: l1, MyPosition: 2},
					{League: l2, MyPosition: 1},
				}, nil
			},
		}
		h := handler.NewLeague(silentLogger(), svc)

		req := httptest.NewRequest(http.MethodGet, "/leagues", nil)
		req = authedRequest(req)
		rec := httptest.NewRecorder()
		h.List(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp struct {
			Leagues []map[string]any `json:"leagues"`
		}
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Len(t, resp.Leagues, 2)
		require.Equal(t, float64(2), resp.Leagues[0]["my_position"])
		require.Equal(t, float64(1), resp.Leagues[1]["my_position"])
	})

	t.Run("returns empty list when no leagues", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			listLeaguesForUserFn: func(context.Context, uuid.UUID) ([]*domain.LeagueSummary, error) {
				return []*domain.LeagueSummary{}, nil
			},
		}
		h := handler.NewLeague(silentLogger(), svc)

		req := httptest.NewRequest(http.MethodGet, "/leagues", nil)
		req = authedRequest(req)
		rec := httptest.NewRecorder()
		h.List(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("returns 500 on service error", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			listLeaguesForUserFn: func(context.Context, uuid.UUID) ([]*domain.LeagueSummary, error) {
				return nil, context.DeadlineExceeded
			},
		}
		h := handler.NewLeague(silentLogger(), svc)

		req := httptest.NewRequest(http.MethodGet, "/leagues", nil)
		req = authedRequest(req)
		rec := httptest.NewRecorder()
		h.List(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

// ---------- UpdateName ----------

func TestLeague_UpdateName(t *testing.T) {
	t.Parallel()

	t.Run("returns 200 on success", func(t *testing.T) {
		t.Parallel()
		league := aLeague()
		league.Name = "Renamed"
		svc := &fakeLeagueService{
			updateLeagueNameFn: func(context.Context, uuid.UUID, uuid.UUID, string) (*domain.League, error) {
				return league, nil
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		rec := patchJSONAuthedWithPathValue(t, h.UpdateName, "/leagues/"+league.ID.String(), "id", league.ID.String(), map[string]any{"name": "Renamed"})

		require.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]any
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Equal(t, "Renamed", resp["name"])
	})

	t.Run("returns 403 when not owner", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			updateLeagueNameFn: func(context.Context, uuid.UUID, uuid.UUID, string) (*domain.League, error) {
				return nil, domain.ErrForbidden
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		id := uuid.New()
		rec := patchJSONAuthedWithPathValue(t, h.UpdateName, "/leagues/"+id.String(), "id", id.String(), map[string]any{"name": "X"})
		require.Equal(t, http.StatusForbidden, rec.Code)
	})
}

// ---------- Delete ----------

func TestLeague_Delete(t *testing.T) {
	t.Parallel()

	t.Run("returns 204 on success", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			deleteLeagueFn: func(context.Context, uuid.UUID, uuid.UUID) error { return nil },
		}
		h := handler.NewLeague(silentLogger(), svc)
		id := uuid.New()
		rec := deleteWithPathValues(t, h.Delete, "/leagues/"+id.String(), "id", id.String())
		require.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("returns 403 when not owner", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			deleteLeagueFn: func(context.Context, uuid.UUID, uuid.UUID) error {
				return domain.ErrForbidden
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		id := uuid.New()
		rec := deleteWithPathValues(t, h.Delete, "/leagues/"+id.String(), "id", id.String())
		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("returns 404 when not found", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			deleteLeagueFn: func(context.Context, uuid.UUID, uuid.UUID) error {
				return domain.ErrNotFound
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		id := uuid.New()
		rec := deleteWithPathValues(t, h.Delete, "/leagues/"+id.String(), "id", id.String())
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

// ---------- RegenerateCode ----------

func TestLeague_RegenerateCode(t *testing.T) {
	t.Parallel()

	t.Run("returns 200 with updated league", func(t *testing.T) {
		t.Parallel()
		league := aLeague()
		league.Code = "NEWCODE1"
		svc := &fakeLeagueService{
			regenerateCodeFn: func(context.Context, uuid.UUID, uuid.UUID) (*domain.League, error) {
				return league, nil
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		rec := postJSONAuthedWithPathValue(t, h.RegenerateCode, "/leagues/"+league.ID.String()+"/code", "id", league.ID.String(), nil)

		require.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]any
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Equal(t, "NEWCODE1", resp["code"])
	})

	t.Run("returns 403 when not owner", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			regenerateCodeFn: func(context.Context, uuid.UUID, uuid.UUID) (*domain.League, error) {
				return nil, domain.ErrForbidden
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		id := uuid.New()
		rec := postJSONAuthedWithPathValue(t, h.RegenerateCode, "/leagues/"+id.String()+"/code", "id", id.String(), nil)
		require.Equal(t, http.StatusForbidden, rec.Code)
	})
}

// ---------- Join ----------

func TestLeague_Join(t *testing.T) {
	t.Parallel()

	t.Run("returns 200 on success", func(t *testing.T) {
		t.Parallel()
		league := aLeague()
		svc := &fakeLeagueService{
			joinLeagueFn: func(context.Context, string, uuid.UUID) (*domain.League, error) {
				return league, nil
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		rec := postJSONAuthed(t, h.Join, "/leagues/join", map[string]any{"code": "TESTCODE"})

		require.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]any
		require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
		require.Equal(t, league.ID.String(), resp["id"])
	})

	t.Run("returns 409 when already a member", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			joinLeagueFn: func(context.Context, string, uuid.UUID) (*domain.League, error) {
				return nil, domain.ErrConflict
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		rec := postJSONAuthed(t, h.Join, "/leagues/join", map[string]any{"code": "TESTCODE"})
		require.Equal(t, http.StatusConflict, rec.Code)
	})

	t.Run("returns 404 on bad code", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			joinLeagueFn: func(context.Context, string, uuid.UUID) (*domain.League, error) {
				return nil, domain.ErrNotFound
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		rec := postJSONAuthed(t, h.Join, "/leagues/join", map[string]any{"code": "BADCODE1"})
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("returns 400 when code missing", func(t *testing.T) {
		t.Parallel()
		h := handler.NewLeague(silentLogger(), &fakeLeagueService{})
		rec := postJSONAuthed(t, h.Join, "/leagues/join", map[string]any{"code": ""})
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

// ---------- RemoveMember ----------

func TestLeague_RemoveMember(t *testing.T) {
	t.Parallel()

	t.Run("returns 204 on success", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			removeMemberFn: func(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) error { return nil },
		}
		h := handler.NewLeague(silentLogger(), svc)
		leagueID := uuid.New()
		memberID := uuid.New()
		rec := deleteWithPathValues(t, h.RemoveMember, "/leagues/"+leagueID.String()+"/members/"+memberID.String(),
			"id", leagueID.String(), "userId", memberID.String())
		require.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("returns 400 when owner tries to self-remove", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			removeMemberFn: func(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) error {
				return domain.ErrInvalid
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		leagueID := uuid.New()
		ownerID := uuid.New()
		rec := deleteWithPathValues(t, h.RemoveMember, "/leagues/"+leagueID.String()+"/members/"+ownerID.String(),
			"id", leagueID.String(), "userId", ownerID.String())
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns 403 when non-owner removes other member", func(t *testing.T) {
		t.Parallel()
		svc := &fakeLeagueService{
			removeMemberFn: func(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) error {
				return domain.ErrForbidden
			},
		}
		h := handler.NewLeague(silentLogger(), svc)
		leagueID := uuid.New()
		memberID := uuid.New()
		rec := deleteWithPathValues(t, h.RemoveMember, "/leagues/"+leagueID.String()+"/members/"+memberID.String(),
			"id", leagueID.String(), "userId", memberID.String())
		require.Equal(t, http.StatusForbidden, rec.Code)
	})
}

// ---------- helpers (test-package-private) ----------

func strReader(s string) *bytesReader {
	return &bytesReader{b: []byte(s)}
}
