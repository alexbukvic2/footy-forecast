package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/server/handler"
)

// ---------- fake service ----------

type fakeTournamentPredictionSvc struct {
	bulkUpsertPlayersFn func(
		context.Context,
		uuid.UUID,
		uuid.UUID,
		[]domain.BulkPlayerPredictionItem,
	) ([]*domain.PlayerPredictionView, error)
	bulkUpsertTeamsFn func(
		context.Context,
		uuid.UUID,
		uuid.UUID,
		[]domain.BulkTeamPredictionItem,
	) ([]*domain.TeamPredictionView, error)
	listMyPlayersFn func(
		context.Context,
		uuid.UUID,
		uuid.UUID,
	) ([]*domain.PlayerPredictionView, error)
	ListMyTeamPredictionsFn func(
		context.Context,
		uuid.UUID,
		uuid.UUID,
	) ([]*domain.TeamPredictionView, error)
	ListLeaguePlayerPredictionsFn func(
		context.Context,
		uuid.UUID,
		uuid.UUID,
	) ([]*domain.LeaguePlayerCategoryView, error)
	listLeagueTeamsFn func(
		context.Context,
		uuid.UUID,
		uuid.UUID,
	) ([]*domain.LeagueTeamCategoryView, error)
}

func (f *fakeTournamentPredictionSvc) BulkUpsertPlayerPredictions(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
	items []domain.BulkPlayerPredictionItem,
) ([]*domain.PlayerPredictionView, error) {
	return f.bulkUpsertPlayersFn(ctx, tournamentID, userID, items)
}
func (f *fakeTournamentPredictionSvc) BulkUpsertTeamPredictions(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
	items []domain.BulkTeamPredictionItem,
) ([]*domain.TeamPredictionView, error) {
	return f.bulkUpsertTeamsFn(ctx, tournamentID, userID, items)
}
func (f *fakeTournamentPredictionSvc) ListPlayerPredictionsForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) ([]*domain.PlayerPredictionView, error) {
	return f.listMyPlayersFn(ctx, tournamentID, userID)
}
func (f *fakeTournamentPredictionSvc) ListTeamPredictionsForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) ([]*domain.TeamPredictionView, error) {
	return f.ListMyTeamPredictionsFn(ctx, tournamentID, userID)
}
func (f *fakeTournamentPredictionSvc) ListLeaguePlayerPredictions(
	ctx context.Context,
	leagueID, userID uuid.UUID,
) ([]*domain.LeaguePlayerCategoryView, error) {
	return f.ListLeaguePlayerPredictionsFn(ctx, leagueID, userID)
}
func (f *fakeTournamentPredictionSvc) ListLeagueTeamPredictions(
	ctx context.Context,
	leagueID, userID uuid.UUID,
) ([]*domain.LeagueTeamCategoryView, error) {
	return f.listLeagueTeamsFn(ctx, leagueID, userID)
}

// ---------- helpers ----------

func putJSONAuthedWithPathValues(
	t *testing.T,
	h http.HandlerFunc,
	path string,
	kvs []string,
	body any,
) *httptest.ResponseRecorder {
	t.Helper()
	raw, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPut, path, strings.NewReader(string(raw)))
	req.Header.Set("Content-Type", "application/json")
	req = authedRequest(req)
	for i := 0; i+1 < len(kvs); i += 2 {
		req.SetPathValue(kvs[i], kvs[i+1])
	}
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

func getAuthedWithPathValues(
	t *testing.T,
	h http.HandlerFunc,
	path string,
	kvs ...string,
) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req = authedRequest(req)
	for i := 0; i+1 < len(kvs); i += 2 {
		req.SetPathValue(kvs[i], kvs[i+1])
	}
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

// ---------- fixtures ----------

func aPlayerPredictionView(
	tournamentID uuid.UUID,
	cat domain.PlayerHandicapCategory,
	playerID uuid.UUID,
) *domain.PlayerPredictionView {
	return &domain.PlayerPredictionView{
		Category: cat,
		Prediction: &domain.PlayerPrediction{
			ID:           uuid.New(),
			UserID:       uuid.New(),
			TournamentID: tournamentID,
			Category:     cat,
			Pick:         playerID,
			PickName:     "Test Player",
		},
	}
}

func aTeamPredictionView(
	tournamentID uuid.UUID,
	cat domain.TeamHandicapCategory,
	teamID uuid.UUID,
) *domain.TeamPredictionView {
	return &domain.TeamPredictionView{
		Category: cat,
		Prediction: &domain.TeamPrediction{
			ID:           uuid.New(),
			UserID:       uuid.New(),
			TournamentID: tournamentID,
			Category:     cat,
			Pick:         teamID,
			PickName:     "Test Team",
		},
	}
}

// ---------- BulkUpsertPlayerPredictions ----------

func TestTournamentPrediction_BulkUpsertPlayers(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	playerID := uuid.New()

	t.Run(
		"returns 200 on success", func(t *testing.T) {
			t.Parallel()
			view := aPlayerPredictionView(tournamentID, domain.PlayerHandicapCategoryTotalTopScorer, playerID)
			svc := &fakeTournamentPredictionSvc{
				bulkUpsertPlayersFn: func(
					_ context.Context,
					_, _ uuid.UUID,
					_ []domain.BulkPlayerPredictionItem,
				) ([]*domain.PlayerPredictionView, error) {
					return []*domain.PlayerPredictionView{view}, nil
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			body := []map[string]any{
				{"category": "total_top_scorer", "player_id": playerID.String()},
			}
			rec := putJSONAuthedWithPathValues(
				t, h.BulkUpsertPlayerPredictions, "/tournaments/"+tournamentID.String()+"/predictions/players",
				[]string{"tournamentId", tournamentID.String()},
				body,
			)
			require.Equal(t, http.StatusOK, rec.Code)
			var resp []map[string]any
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			require.Len(t, resp, 1)
			require.Equal(t, "total_top_scorer", resp[0]["category"])
			require.Equal(t, playerID.String(), resp[0]["player_id"])
		},
	)

	t.Run(
		"accepts null player_id to clear slot", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentPredictionSvc{
				bulkUpsertPlayersFn: func(
					_ context.Context,
					_, _ uuid.UUID,
					items []domain.BulkPlayerPredictionItem,
				) ([]*domain.PlayerPredictionView, error) {
					require.Len(t, items, 1)
					require.Nil(t, items[0].PlayerID)
					return []*domain.PlayerPredictionView{}, nil
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			body := []map[string]any{
				{"category": "total_top_scorer", "player_id": nil},
			}
			rec := putJSONAuthedWithPathValues(
				t, h.BulkUpsertPlayerPredictions, "/tournaments/"+tournamentID.String()+"/predictions/players",
				[]string{"tournamentId", tournamentID.String()},
				body,
			)
			require.Equal(t, http.StatusOK, rec.Code)
		},
	)

	t.Run(
		"returns 400 on malformed JSON", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			req := httptest.NewRequest(http.MethodPut, "/tournaments/"+tournamentID.String()+"/predictions/players", strings.NewReader("not-json"))
			req.Header.Set("Content-Type", "application/json")
			req = authedRequest(req)
			req.SetPathValue("tournamentId", tournamentID.String())
			rec := httptest.NewRecorder()
			h.BulkUpsertPlayerPredictions(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 400 on invalid category", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			body := []map[string]any{{"category": "bad_cat", "player_id": playerID.String()}}
			rec := putJSONAuthedWithPathValues(
				t, h.BulkUpsertPlayerPredictions, "/tournaments/"+tournamentID.String()+"/predictions/players",
				[]string{"tournamentId", tournamentID.String()},
				body,
			)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 400 on invalid player_id", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			body := []map[string]any{{"category": "total_top_scorer", "player_id": "not-a-uuid"}}
			rec := putJSONAuthedWithPathValues(
				t, h.BulkUpsertPlayerPredictions, "/tournaments/"+tournamentID.String()+"/predictions/players",
				[]string{"tournamentId", tournamentID.String()},
				body,
			)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 403 when locked", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentPredictionSvc{
				bulkUpsertPlayersFn: func(_ context.Context, _, _ uuid.UUID, _ []domain.BulkPlayerPredictionItem) ([]*domain.PlayerPredictionView, error) {
					return nil, domain.ErrForbidden
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			body := []map[string]any{{"category": "total_top_scorer", "player_id": playerID.String()}}
			rec := putJSONAuthedWithPathValues(
				t, h.BulkUpsertPlayerPredictions, "/tournaments/"+tournamentID.String()+"/predictions/players",
				[]string{"tournamentId", tournamentID.String()},
				body,
			)
			require.Equal(t, http.StatusForbidden, rec.Code)
		},
	)

	t.Run(
		"returns 404 when player not found", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentPredictionSvc{
				bulkUpsertPlayersFn: func(_ context.Context, _, _ uuid.UUID, _ []domain.BulkPlayerPredictionItem) ([]*domain.PlayerPredictionView, error) {
					return nil, domain.ErrNotFound
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			body := []map[string]any{{"category": "total_top_scorer", "player_id": playerID.String()}}
			rec := putJSONAuthedWithPathValues(
				t, h.BulkUpsertPlayerPredictions, "/tournaments/"+tournamentID.String()+"/predictions/players",
				[]string{"tournamentId", tournamentID.String()},
				body,
			)
			require.Equal(t, http.StatusNotFound, rec.Code)
		},
	)

	t.Run(
		"returns 401 when unauthenticated", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			raw, _ := json.Marshal([]map[string]any{{"category": "total_top_scorer", "player_id": playerID.String()}})
			req := httptest.NewRequest(http.MethodPut, "/tournaments/"+tournamentID.String()+"/predictions/players", strings.NewReader(string(raw)))
			req.Header.Set("Content-Type", "application/json")
			req.SetPathValue("tournamentId", tournamentID.String())
			rec := httptest.NewRecorder()
			h.BulkUpsertPlayerPredictions(rec, req)
			require.Equal(t, http.StatusUnauthorized, rec.Code)
		},
	)
}

// ---------- BulkUpsertTeamPredictions ----------

func TestTournamentPrediction_BulkUpsertTeams(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()
	teamID := uuid.New()

	t.Run(
		"returns 200 on success", func(t *testing.T) {
			t.Parallel()
			view := aTeamPredictionView(tournamentID, domain.TeamHandicapCategoryWinner, teamID)
			svc := &fakeTournamentPredictionSvc{
				bulkUpsertTeamsFn: func(
					_ context.Context,
					_, _ uuid.UUID,
					_ []domain.BulkTeamPredictionItem,
				) ([]*domain.TeamPredictionView, error) {
					return []*domain.TeamPredictionView{view}, nil
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			body := []map[string]any{
				{"category": "winner", "slot_index": 0, "team_id": teamID.String()},
			}
			rec := putJSONAuthedWithPathValues(
				t, h.BulkUpsertTeamPredictions, "/tournaments/"+tournamentID.String()+"/predictions/teams",
				[]string{"tournamentId", tournamentID.String()},
				body,
			)
			require.Equal(t, http.StatusOK, rec.Code)
			var resp []map[string]any
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			require.Len(t, resp, 1)
			require.Equal(t, "winner", resp[0]["category"])
			require.Equal(t, teamID.String(), resp[0]["team_id"])
		},
	)

	t.Run(
		"accepts null team_id to clear slot", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentPredictionSvc{
				bulkUpsertTeamsFn: func(
					_ context.Context,
					_, _ uuid.UUID,
					items []domain.BulkTeamPredictionItem,
				) ([]*domain.TeamPredictionView, error) {
					require.Len(t, items, 1)
					require.Nil(t, items[0].TeamID)
					return []*domain.TeamPredictionView{}, nil
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			body := []map[string]any{
				{"category": "winner", "slot_index": 0, "team_id": nil},
			}
			rec := putJSONAuthedWithPathValues(
				t, h.BulkUpsertTeamPredictions, "/tournaments/"+tournamentID.String()+"/predictions/teams",
				[]string{"tournamentId", tournamentID.String()},
				body,
			)
			require.Equal(t, http.StatusOK, rec.Code)
		},
	)

	t.Run(
		"returns 400 on malformed JSON", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			req := httptest.NewRequest(http.MethodPut, "/tournaments/"+tournamentID.String()+"/predictions/teams", strings.NewReader("not-json"))
			req.Header.Set("Content-Type", "application/json")
			req = authedRequest(req)
			req.SetPathValue("tournamentId", tournamentID.String())
			rec := httptest.NewRecorder()
			h.BulkUpsertTeamPredictions(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 400 on invalid category", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			body := []map[string]any{{"category": "bad_cat", "slot_index": 0, "team_id": teamID.String()}}
			rec := putJSONAuthedWithPathValues(
				t, h.BulkUpsertTeamPredictions, "/tournaments/"+tournamentID.String()+"/predictions/teams",
				[]string{"tournamentId", tournamentID.String()},
				body,
			)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 400 on invalid team_id", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			body := []map[string]any{{"category": "winner", "slot_index": 0, "team_id": "not-a-uuid"}}
			rec := putJSONAuthedWithPathValues(
				t, h.BulkUpsertTeamPredictions, "/tournaments/"+tournamentID.String()+"/predictions/teams",
				[]string{"tournamentId", tournamentID.String()},
				body,
			)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 401 when unauthenticated", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			raw, _ := json.Marshal([]map[string]any{{"category": "winner", "slot_index": 0, "team_id": teamID.String()}})
			req := httptest.NewRequest(http.MethodPut, "/tournaments/"+tournamentID.String()+"/predictions/teams", strings.NewReader(string(raw)))
			req.Header.Set("Content-Type", "application/json")
			req.SetPathValue("tournamentId", tournamentID.String())
			rec := httptest.NewRecorder()
			h.BulkUpsertTeamPredictions(rec, req)
			require.Equal(t, http.StatusUnauthorized, rec.Code)
		},
	)

	t.Run(
		"returns 403 when locked", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentPredictionSvc{
				bulkUpsertTeamsFn: func(_ context.Context, _, _ uuid.UUID, _ []domain.BulkTeamPredictionItem) ([]*domain.TeamPredictionView, error) {
					return nil, domain.ErrForbidden
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			body := []map[string]any{{"category": "winner", "slot_index": 0, "team_id": teamID.String()}}
			rec := putJSONAuthedWithPathValues(
				t, h.BulkUpsertTeamPredictions, "/tournaments/"+tournamentID.String()+"/predictions/teams",
				[]string{"tournamentId", tournamentID.String()},
				body,
			)
			require.Equal(t, http.StatusForbidden, rec.Code)
		},
	)

	t.Run(
		"returns 404 when team not found", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentPredictionSvc{
				bulkUpsertTeamsFn: func(_ context.Context, _, _ uuid.UUID, _ []domain.BulkTeamPredictionItem) ([]*domain.TeamPredictionView, error) {
					return nil, domain.ErrNotFound
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			body := []map[string]any{{"category": "winner", "slot_index": 0, "team_id": teamID.String()}}
			rec := putJSONAuthedWithPathValues(
				t, h.BulkUpsertTeamPredictions, "/tournaments/"+tournamentID.String()+"/predictions/teams",
				[]string{"tournamentId", tournamentID.String()},
				body,
			)
			require.Equal(t, http.StatusNotFound, rec.Code)
		},
	)
}

// ---------- ListMyPlayerPredictions ----------

func TestTournamentPrediction_ListMyPlayers(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()

	t.Run(
		"returns 200 with all categories", func(t *testing.T) {
			t.Parallel()
			playerID := uuid.New()
			views := []*domain.PlayerPredictionView{
				{
					Category: domain.PlayerHandicapCategoryGroupTopScorer,
					Prediction: &domain.PlayerPrediction{
						ID:       uuid.New(),
						Pick:     playerID,
						PickName: "Lionel Messi",
					},
				},
				{
					Category:   domain.PlayerHandicapCategoryTotalTopScorer,
					Prediction: nil,
				},
			}
			svc := &fakeTournamentPredictionSvc{
				listMyPlayersFn: func(
					_ context.Context,
					_, _ uuid.UUID,
				) ([]*domain.PlayerPredictionView, error) {
					return views, nil
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			rec := getAuthedWithPathValues(
				t, h.ListMyPlayerPredictions, "/tournaments/"+tournamentID.String()+"/predictions/players",
				"tournamentId", tournamentID.String(),
			)

			require.Equal(t, http.StatusOK, rec.Code)
			var resp []map[string]any
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			require.Len(t, resp, 2)
			require.Equal(t, "group_top_scorer", resp[0]["category"])
			require.Equal(t, playerID.String(), resp[0]["player_id"])
			require.Nil(t, resp[1]["player_id"])
		},
	)

	t.Run(
		"returns 400 on invalid tournament ID", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			rec := getAuthedWithPathValues(
				t, h.ListMyPlayerPredictions, "/tournaments/bad-uuid/predictions/players",
				"tournamentId", "bad-uuid",
			)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 401 when unauthenticated", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			req := httptest.NewRequest(
				http.MethodGet,
				"/tournaments/"+tournamentID.String()+"/predictions/players",
				nil,
			)
			req.SetPathValue("tournamentId", tournamentID.String())
			rec := httptest.NewRecorder()
			h.ListMyPlayerPredictions(rec, req)
			require.Equal(t, http.StatusUnauthorized, rec.Code)
		},
	)
}

// ---------- ListMyTeamPredictions ----------

func TestTournamentPrediction_ListMyTeamPredictions(t *testing.T) {
	t.Parallel()

	tournamentID := uuid.New()

	t.Run(
		"returns 200 with all categories (none predicted)", func(t *testing.T) {
			t.Parallel()
			views := []*domain.TeamPredictionView{
				{Category: domain.TeamHandicapCategoryWinner, Prediction: nil},
			}
			svc := &fakeTournamentPredictionSvc{
				ListMyTeamPredictionsFn: func(
					_ context.Context,
					_, _ uuid.UUID,
				) ([]*domain.TeamPredictionView, error) {
					return views, nil
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			rec := getAuthedWithPathValues(
				t, h.ListMyTeamPredictions, "/tournaments/"+tournamentID.String()+"/predictions/teams",
				"tournamentId", tournamentID.String(),
			)

			require.Equal(t, http.StatusOK, rec.Code)
			var resp []map[string]any
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			require.Len(t, resp, 1)
			require.Equal(t, "winner", resp[0]["category"])
			require.Nil(t, resp[0]["team_id"])
		},
	)

	t.Run(
		"returns 401 when unauthenticated", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			req := httptest.NewRequest(http.MethodGet, "/tournaments/"+tournamentID.String()+"/predictions/teams", nil)
			req.SetPathValue("tournamentId", tournamentID.String())
			rec := httptest.NewRecorder()
			h.ListMyTeamPredictions(rec, req)
			require.Equal(t, http.StatusUnauthorized, rec.Code)
		},
	)
}

// ---------- ListLeaguePlayerPredictions ----------

func TestTournamentPrediction_ListLeaguePlayerPredictions(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()
	userID := uuid.New()

	t.Run(
		"returns 200 with grouped predictions", func(t *testing.T) {
			t.Parallel()
			playerID := uuid.New()
			pName := "Lionel Messi"
			views := []*domain.LeaguePlayerCategoryView{
				{
					Category: domain.PlayerHandicapCategoryGroupTopScorer,
					Predictions: []domain.LeagueMemberPlayerPick{
						{UserID: userID, DisplayName: "Alice", PlayerID: &playerID, PlayerName: &pName},
					},
				},
			}
			svc := &fakeTournamentPredictionSvc{
				ListLeaguePlayerPredictionsFn: func(
					_ context.Context,
					_, _ uuid.UUID,
				) ([]*domain.LeaguePlayerCategoryView, error) {
					return views, nil
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			rec := getAuthedWithPathValues(
				t, h.ListLeaguePlayerPredictions, "/leagues/"+leagueID.String()+"/predictions/players",
				"leagueId", leagueID.String(),
			)

			require.Equal(t, http.StatusOK, rec.Code)
			var resp []map[string]any
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			require.Len(t, resp, 1)
			require.Equal(t, "group_top_scorer", resp[0]["category"])
			preds := resp[0]["predictions"].([]any)
			require.Len(t, preds, 1)
			member := preds[0].(map[string]any)
			require.Equal(t, playerID.String(), member["player_id"])
		},
	)

	t.Run(
		"returns 403 when not a member", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentPredictionSvc{
				ListLeaguePlayerPredictionsFn: func(
					_ context.Context,
					_, _ uuid.UUID,
				) ([]*domain.LeaguePlayerCategoryView, error) {
					return nil, domain.ErrForbidden
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			rec := getAuthedWithPathValues(
				t, h.ListLeaguePlayerPredictions, "/leagues/"+leagueID.String()+"/predictions/players",
				"leagueId", leagueID.String(),
			)
			require.Equal(t, http.StatusForbidden, rec.Code)
		},
	)

	t.Run(
		"returns 403 before lock", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentPredictionSvc{
				ListLeaguePlayerPredictionsFn: func(
					_ context.Context,
					_, _ uuid.UUID,
				) ([]*domain.LeaguePlayerCategoryView, error) {
					return nil, domain.ErrForbidden
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			rec := getAuthedWithPathValues(
				t, h.ListLeaguePlayerPredictions, "/leagues/"+leagueID.String()+"/predictions/players",
				"leagueId", leagueID.String(),
			)
			require.Equal(t, http.StatusForbidden, rec.Code)
		},
	)

	t.Run(
		"returns 400 on invalid league ID", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			rec := getAuthedWithPathValues(
				t, h.ListLeaguePlayerPredictions, "/leagues/bad-uuid/predictions/players",
				"leagueId", "bad-uuid",
			)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 401 when unauthenticated", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			req := httptest.NewRequest(http.MethodGet, "/leagues/"+leagueID.String()+"/predictions/players", nil)
			req.SetPathValue("leagueId", leagueID.String())
			rec := httptest.NewRecorder()
			h.ListLeaguePlayerPredictions(rec, req)
			require.Equal(t, http.StatusUnauthorized, rec.Code)
		},
	)
}

// ---------- ListLeagueTeamPredictions ----------

func TestTournamentPrediction_ListLeagueTeams(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()

	t.Run(
		"returns 200 with grouped predictions", func(t *testing.T) {
			t.Parallel()
			teamID := uuid.New()
			tName := "Brazil"
			views := []*domain.LeagueTeamCategoryView{
				{
					Category: domain.TeamHandicapCategoryWinner,
					Predictions: []domain.LeagueMemberTeamPick{
						{UserID: uuid.New(), DisplayName: "Alice", TeamID: &teamID, TeamName: &tName},
					},
				},
			}
			svc := &fakeTournamentPredictionSvc{
				listLeagueTeamsFn: func(
					_ context.Context,
					_, _ uuid.UUID,
				) ([]*domain.LeagueTeamCategoryView, error) {
					return views, nil
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			rec := getAuthedWithPathValues(
				t, h.ListLeagueTeamPredictions, "/leagues/"+leagueID.String()+"/predictions/teams",
				"leagueId", leagueID.String(),
			)

			require.Equal(t, http.StatusOK, rec.Code)
			var resp []map[string]any
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			require.Len(t, resp, 1)
			require.Equal(t, "winner", resp[0]["category"])
			preds := resp[0]["predictions"].([]any)
			require.Len(t, preds, 1)
			member := preds[0].(map[string]any)
			require.Equal(t, teamID.String(), member["team_id"])
		},
	)

	t.Run(
		"returns 403 when not a member", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentPredictionSvc{
				listLeagueTeamsFn: func(
					_ context.Context,
					_, _ uuid.UUID,
				) ([]*domain.LeagueTeamCategoryView, error) {
					return nil, domain.ErrForbidden
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			rec := getAuthedWithPathValues(
				t, h.ListLeagueTeamPredictions, "/leagues/"+leagueID.String()+"/predictions/teams",
				"leagueId", leagueID.String(),
			)
			require.Equal(t, http.StatusForbidden, rec.Code)
		},
	)

	t.Run(
		"returns 401 when unauthenticated", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			req := httptest.NewRequest(http.MethodGet, "/leagues/"+leagueID.String()+"/predictions/teams", nil)
			req.SetPathValue("leagueId", leagueID.String())
			rec := httptest.NewRecorder()
			h.ListLeagueTeamPredictions(rec, req)
			require.Equal(t, http.StatusUnauthorized, rec.Code)
		},
	)
}
