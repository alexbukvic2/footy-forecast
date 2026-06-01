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
	) (bool, []*domain.PlayerPredictionView, error)
	ListMyTeamPredictionsFn func(
		context.Context,
		uuid.UUID,
		uuid.UUID,
	) (bool, []*domain.TeamPredictionView, error)
	listLeagueGroupFn func(
		context.Context,
		uuid.UUID,
		uuid.UUID,
		string,
	) (*domain.LeagueGroupPredictions, error)
	listLeaguePlayoffFn func(
		context.Context,
		uuid.UUID,
		uuid.UUID,
	) (*domain.LeaguePlayoffPredictions, error)
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
) (bool, []*domain.PlayerPredictionView, error) {
	return f.listMyPlayersFn(ctx, tournamentID, userID)
}
func (f *fakeTournamentPredictionSvc) ListTeamPredictionsForUser(
	ctx context.Context,
	tournamentID, userID uuid.UUID,
) (bool, []*domain.TeamPredictionView, error) {
	return f.ListMyTeamPredictionsFn(ctx, tournamentID, userID)
}
func (f *fakeTournamentPredictionSvc) ListLeagueGroupPredictions(
	ctx context.Context,
	leagueID, userID uuid.UUID,
	groupLetter string,
) (*domain.LeagueGroupPredictions, error) {
	if f.listLeagueGroupFn != nil {
		return f.listLeagueGroupFn(ctx, leagueID, userID, groupLetter)
	}
	return &domain.LeagueGroupPredictions{}, nil
}
func (f *fakeTournamentPredictionSvc) ListLeaguePlayoffPredictions(
	ctx context.Context,
	leagueID, userID uuid.UUID,
) (*domain.LeaguePlayoffPredictions, error) {
	if f.listLeaguePlayoffFn != nil {
		return f.listLeaguePlayoffFn(ctx, leagueID, userID)
	}
	return &domain.LeaguePlayoffPredictions{}, nil
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
				) (bool, []*domain.PlayerPredictionView, error) {
					return false, views, nil
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			rec := getAuthedWithPathValues(
				t, h.ListMyPlayerPredictions, "/tournaments/"+tournamentID.String()+"/predictions/players",
				"tournamentId", tournamentID.String(),
			)

			require.Equal(t, http.StatusOK, rec.Code)
			var resp map[string]any
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			require.Equal(t, false, resp["locked"])
			preds := resp["predictions"].([]any)
			require.Len(t, preds, 2)
			require.Equal(t, "group_top_scorer", preds[0].(map[string]any)["category"])
			require.Equal(t, playerID.String(), preds[0].(map[string]any)["player_id"])
			require.Nil(t, preds[1].(map[string]any)["player_id"])
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
				) (bool, []*domain.TeamPredictionView, error) {
					return false, views, nil
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			rec := getAuthedWithPathValues(
				t, h.ListMyTeamPredictions, "/tournaments/"+tournamentID.String()+"/predictions/teams",
				"tournamentId", tournamentID.String(),
			)

			require.Equal(t, http.StatusOK, rec.Code)
			var resp map[string]any
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			require.Equal(t, false, resp["locked"])
			preds := resp["predictions"].([]any)
			require.Len(t, preds, 1)
			require.Equal(t, "winner", preds[0].(map[string]any)["category"])
			require.Nil(t, preds[0].(map[string]any)["team_id"])
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

// ---------- ListLeagueGroupPredictions ----------

func TestTournamentPrediction_ListLeagueGroupPredictions(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()
	userID := uuid.New()

	t.Run(
		"returns 200 with group predictions", func(t *testing.T) {
			t.Parallel()
			playerID := uuid.New()
			pName := "Lionel Messi"
			teamID := uuid.New()
			tName := "Argentina"
			groupLetter := "A"
			result := &domain.LeagueGroupPredictions{
				Group: "A",
				TeamPredictions: []*domain.LeagueTeamCategoryView{
					{
						Category:    domain.TeamHandicapCategoryGroupWinner,
						GroupLetter: &groupLetter,
						SlotIndex:   0,
						Predictions: []domain.LeagueMemberTeamPick{
							{UserID: userID, DisplayName: "Alice", TeamID: &teamID, TeamName: &tName},
						},
					},
					{
						Category:    domain.TeamHandicapCategoryPlayoff,
						GroupLetter: &groupLetter,
						SlotIndex:   0,
						Predictions: []domain.LeagueMemberTeamPick{
							{UserID: userID, DisplayName: "Alice", TeamID: &teamID, TeamName: &tName},
						},
					},
					{
						Category:    domain.TeamHandicapCategoryPlayoff,
						GroupLetter: &groupLetter,
						SlotIndex:   1,
						Predictions: []domain.LeagueMemberTeamPick{
							{UserID: userID, DisplayName: "Alice"},
						},
					},
					{
						Category:    domain.TeamHandicapCategoryPlayoff,
						GroupLetter: &groupLetter,
						SlotIndex:   2,
						Predictions: []domain.LeagueMemberTeamPick{
							{UserID: userID, DisplayName: "Alice"},
						},
					},
				},
				PlayerPredictions: []*domain.LeaguePlayerCategoryView{
					{
						Category:    domain.PlayerHandicapCategoryGroupTopScorer,
						GroupLetter: &groupLetter,
						Predictions: []domain.LeagueMemberPlayerPick{
							{UserID: userID, DisplayName: "Alice", PlayerID: &playerID, PlayerName: &pName},
						},
					},
				},
			}
			svc := &fakeTournamentPredictionSvc{
				listLeagueGroupFn: func(
					_ context.Context,
					_, _ uuid.UUID,
					_ string,
				) (*domain.LeagueGroupPredictions, error) {
					return result, nil
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			req := httptest.NewRequest(http.MethodGet, "/leagues/"+leagueID.String()+"/predictions/groups?group=A", nil)
			req = authedRequest(req)
			req.SetPathValue("leagueId", leagueID.String())
			rec := httptest.NewRecorder()
			h.ListLeagueGroupPredictions(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
			var resp map[string]any
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			require.Equal(t, "A", resp["group"])
			teamPreds := resp["team_predictions"].([]any)
			require.Len(t, teamPreds, 4)
			require.Equal(t, "group_winner", teamPreds[0].(map[string]any)["category"])
			require.Equal(t, "playoff", teamPreds[1].(map[string]any)["category"])
			require.Equal(t, "playoff", teamPreds[2].(map[string]any)["category"])
			require.Equal(t, "playoff", teamPreds[3].(map[string]any)["category"])
			playerPreds := resp["player_predictions"].([]any)
			require.Len(t, playerPreds, 1)
			require.Equal(t, "group_top_scorer", playerPreds[0].(map[string]any)["category"])
		},
	)

	t.Run(
		"returns 400 when group param missing", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			req := httptest.NewRequest(http.MethodGet, "/leagues/"+leagueID.String()+"/predictions/groups", nil)
			req = authedRequest(req)
			req.SetPathValue("leagueId", leagueID.String())
			rec := httptest.NewRecorder()
			h.ListLeagueGroupPredictions(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 400 when group param too long", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			req := httptest.NewRequest(http.MethodGet, "/leagues/"+leagueID.String()+"/predictions/groups?group=AB", nil)
			req = authedRequest(req)
			req.SetPathValue("leagueId", leagueID.String())
			rec := httptest.NewRecorder()
			h.ListLeagueGroupPredictions(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 400 on invalid league ID", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			req := httptest.NewRequest(http.MethodGet, "/leagues/bad-uuid/predictions/groups?group=A", nil)
			req = authedRequest(req)
			req.SetPathValue("leagueId", "bad-uuid")
			rec := httptest.NewRecorder()
			h.ListLeagueGroupPredictions(rec, req)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 401 when unauthenticated", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			req := httptest.NewRequest(http.MethodGet, "/leagues/"+leagueID.String()+"/predictions/groups?group=A", nil)
			req.SetPathValue("leagueId", leagueID.String())
			rec := httptest.NewRecorder()
			h.ListLeagueGroupPredictions(rec, req)
			require.Equal(t, http.StatusUnauthorized, rec.Code)
		},
	)

	t.Run(
		"returns 403 when not a member", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentPredictionSvc{
				listLeagueGroupFn: func(
					_ context.Context,
					_, _ uuid.UUID,
					_ string,
				) (*domain.LeagueGroupPredictions, error) {
					return nil, domain.ErrForbidden
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			req := httptest.NewRequest(http.MethodGet, "/leagues/"+leagueID.String()+"/predictions/groups?group=A", nil)
			req = authedRequest(req)
			req.SetPathValue("leagueId", leagueID.String())
			rec := httptest.NewRecorder()
			h.ListLeagueGroupPredictions(rec, req)
			require.Equal(t, http.StatusForbidden, rec.Code)
		},
	)

	t.Run(
		"returns 404 when league not found", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentPredictionSvc{
				listLeagueGroupFn: func(
					_ context.Context,
					_, _ uuid.UUID,
					_ string,
				) (*domain.LeagueGroupPredictions, error) {
					return nil, domain.ErrNotFound
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			req := httptest.NewRequest(http.MethodGet, "/leagues/"+leagueID.String()+"/predictions/groups?group=A", nil)
			req = authedRequest(req)
			req.SetPathValue("leagueId", leagueID.String())
			rec := httptest.NewRecorder()
			h.ListLeagueGroupPredictions(rec, req)
			require.Equal(t, http.StatusNotFound, rec.Code)
		},
	)
}

// ---------- ListLeaguePlayoffPredictions ----------

func TestTournamentPrediction_ListLeaguePlayoffPredictions(t *testing.T) {
	t.Parallel()

	leagueID := uuid.New()

	t.Run(
		"returns 200 with playoff predictions", func(t *testing.T) {
			t.Parallel()
			teamID := uuid.New()
			tName := "Brazil"
			playerID := uuid.New()
			pName := "Neymar"
			result := &domain.LeaguePlayoffPredictions{
				TeamPredictions: []*domain.LeagueTeamCategoryView{
					{
						Category:    domain.TeamHandicapCategoryWinner,
						GroupLetter: nil,
						SlotIndex:   0,
						Predictions: []domain.LeagueMemberTeamPick{
							{UserID: uuid.New(), DisplayName: "Alice", TeamID: &teamID, TeamName: &tName},
						},
					},
				},
				PlayerPredictions: []*domain.LeaguePlayerCategoryView{
					{
						Category:    domain.PlayerHandicapCategoryTotalTopScorer,
						GroupLetter: nil,
						Predictions: []domain.LeagueMemberPlayerPick{
							{UserID: uuid.New(), DisplayName: "Alice", PlayerID: &playerID, PlayerName: &pName},
						},
					},
				},
			}
			svc := &fakeTournamentPredictionSvc{
				listLeaguePlayoffFn: func(
					_ context.Context,
					_, _ uuid.UUID,
				) (*domain.LeaguePlayoffPredictions, error) {
					return result, nil
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			rec := getAuthedWithPathValues(
				t, h.ListLeaguePlayoffPredictions, "/leagues/"+leagueID.String()+"/predictions/playoff",
				"leagueId", leagueID.String(),
			)

			require.Equal(t, http.StatusOK, rec.Code)
			var resp map[string]any
			require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
			teamPreds := resp["team_predictions"].([]any)
			require.Len(t, teamPreds, 1)
			require.Equal(t, "winner", teamPreds[0].(map[string]any)["category"])
			playerPreds := resp["player_predictions"].([]any)
			require.Len(t, playerPreds, 1)
			require.Equal(t, "total_top_scorer", playerPreds[0].(map[string]any)["category"])
		},
	)

	t.Run(
		"returns 400 on invalid league ID", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			rec := getAuthedWithPathValues(
				t, h.ListLeaguePlayoffPredictions, "/leagues/bad-uuid/predictions/playoff",
				"leagueId", "bad-uuid",
			)
			require.Equal(t, http.StatusBadRequest, rec.Code)
		},
	)

	t.Run(
		"returns 401 when unauthenticated", func(t *testing.T) {
			t.Parallel()
			h := handler.NewTournamentPrediction(silentLogger(), &fakeTournamentPredictionSvc{})
			req := httptest.NewRequest(http.MethodGet, "/leagues/"+leagueID.String()+"/predictions/playoff", nil)
			req.SetPathValue("leagueId", leagueID.String())
			rec := httptest.NewRecorder()
			h.ListLeaguePlayoffPredictions(rec, req)
			require.Equal(t, http.StatusUnauthorized, rec.Code)
		},
	)

	t.Run(
		"returns 403 when not a member", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentPredictionSvc{
				listLeaguePlayoffFn: func(
					_ context.Context,
					_, _ uuid.UUID,
				) (*domain.LeaguePlayoffPredictions, error) {
					return nil, domain.ErrForbidden
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			rec := getAuthedWithPathValues(
				t, h.ListLeaguePlayoffPredictions, "/leagues/"+leagueID.String()+"/predictions/playoff",
				"leagueId", leagueID.String(),
			)
			require.Equal(t, http.StatusForbidden, rec.Code)
		},
	)

	t.Run(
		"returns 404 when league not found", func(t *testing.T) {
			t.Parallel()
			svc := &fakeTournamentPredictionSvc{
				listLeaguePlayoffFn: func(
					_ context.Context,
					_, _ uuid.UUID,
				) (*domain.LeaguePlayoffPredictions, error) {
					return nil, domain.ErrNotFound
				},
			}
			h := handler.NewTournamentPrediction(silentLogger(), svc)
			rec := getAuthedWithPathValues(
				t, h.ListLeaguePlayoffPredictions, "/leagues/"+leagueID.String()+"/predictions/playoff",
				"leagueId", leagueID.String(),
			)
			require.Equal(t, http.StatusNotFound, rec.Code)
		},
	)
}
