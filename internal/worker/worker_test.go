package worker_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
	"github.com/alexbukvic2/footy-forecast/internal/worker"
)

// ---- fakes ----

type fakeRepo struct {
	mu       sync.Mutex
	calls    map[string]int
	fixtures []domain.PollableFixture
	listErr  error

	groupComplete      bool
	roundComplete      bool
	groupStageComplete bool

	updateErr error
}

func newFakeRepo(fixtures []domain.PollableFixture) *fakeRepo {
	return &fakeRepo{
		calls:    make(map[string]int),
		fixtures: fixtures,
	}
}

func (r *fakeRepo) inc(name string) {
	r.mu.Lock()
	r.calls[name]++
	r.mu.Unlock()
}

func (r *fakeRepo) callCount(name string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.calls[name]
}

func (r *fakeRepo) ListPollableMatches(_ context.Context) ([]domain.PollableFixture, error) {
	r.inc("ListPollableMatches")
	return r.fixtures, r.listErr
}

func (r *fakeRepo) UpdateMatchAndRescoreLivePredictions(_ context.Context, _ domain.PollableFixture, _ worker.APIFixtureResult) error {
	r.inc("UpdateMatchAndRescoreLivePredictions")
	return r.updateErr
}

func (r *fakeRepo) UpdateGroupStandings(_ context.Context, _ uuid.UUID, _ string, _ []domain.StandingsEntry) error {
	r.inc("UpdateGroupStandings")
	return nil
}

func (r *fakeRepo) IsGroupComplete(_ context.Context, _ uuid.UUID, _ string) (bool, error) {
	r.inc("IsGroupComplete")
	return r.groupComplete, nil
}

func (r *fakeRepo) IsRoundComplete(_ context.Context, _ uuid.UUID, _ string) (bool, error) {
	r.inc("IsRoundComplete")
	return r.roundComplete, nil
}

func (r *fakeRepo) IsGroupStageComplete(_ context.Context, _ uuid.UUID) (bool, error) {
	r.inc("IsGroupStageComplete")
	return r.groupStageComplete, nil
}

func (r *fakeRepo) GetTeamByExternalID(_ context.Context, _ int64, _ uuid.UUID) (uuid.UUID, error) {
	r.inc("GetTeamByExternalID")
	return uuid.New(), nil
}

func (r *fakeRepo) GetPlayerByExternalID(_ context.Context, _ string, _ uuid.UUID) (uuid.UUID, error) {
	r.inc("GetPlayerByExternalID")
	return uuid.New(), nil
}

func (r *fakeRepo) SettleGroupWinnerPredictions(_ context.Context, _ uuid.UUID, _ string) error {
	r.inc("SettleGroupWinnerPredictions")
	return nil
}

func (r *fakeRepo) SettlePlayoffGroupPredictions(_ context.Context, _ uuid.UUID, _ string) error {
	r.inc("SettlePlayoffGroupPredictions")
	return nil
}

func (r *fakeRepo) SettlePlayoffWildcardPredictions(_ context.Context, _ uuid.UUID) error {
	r.inc("SettlePlayoffWildcardPredictions")
	return nil
}

func (r *fakeRepo) SettleGroupTopScorerPredictions(_ context.Context, _ uuid.UUID, _ string, _ uuid.UUID) error {
	r.inc("SettleGroupTopScorerPredictions")
	return nil
}

func (r *fakeRepo) SettleSemifinalistPredictions(_ context.Context, _ uuid.UUID) error {
	r.inc("SettleSemifinalistPredictions")
	return nil
}

func (r *fakeRepo) SettleTournamentWinnerPredictions(_ context.Context, _ uuid.UUID, _ uuid.UUID) error {
	r.inc("SettleTournamentWinnerPredictions")
	return nil
}

func (r *fakeRepo) SettleTopScorerPredictions(_ context.Context, _ uuid.UUID, _ uuid.UUID) error {
	r.inc("SettleTopScorerPredictions")
	return nil
}

type fakeAPI struct {
	mu           sync.Mutex
	callCount    int64
	result       worker.APIFixtureResult
	err          error
	standings    []worker.APIStandingsEntry
	standingsErr error
	topScorer    worker.APITopScorerResult
	topScorerErr error
}

func (a *fakeAPI) GetFixture(_ context.Context, _ int64) (worker.APIFixtureResult, error) {
	atomic.AddInt64(&a.callCount, 1)
	return a.result, a.err
}

func (a *fakeAPI) GetStandings(_ context.Context, _ int64, _ int) ([]worker.APIStandingsEntry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.standings, a.standingsErr
}

func (a *fakeAPI) GetGroupTopScorer(_ context.Context, _ int64, _ int, _ string) (worker.APITopScorerResult, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.topScorer, a.topScorerErr
}

func (a *fakeAPI) GetTournamentTopScorer(_ context.Context, _ int64, _ int) (worker.APITopScorerResult, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.topScorer, a.topScorerErr
}

type fakeClock struct{ t time.Time }

func (c fakeClock) Now() time.Time { return c.t }

func newWorker(repo worker.Repo, api worker.MatchAPI) *worker.Worker {
	return worker.New(repo, api, fakeClock{t: time.Now()}, nopLogger(), time.Minute)
}

// ---- helpers ----

func liveFixture() domain.PollableFixture {
	g := "A"
	return domain.PollableFixture{
		ID:                   uuid.New(),
		ExternalID:           1001,
		TournamentID:         uuid.New(),
		TournamentExternalID: 42,
		TournamentSeason:     2026,
		HomeTeamID:           uuid.New(),
		AwayTeamID:           uuid.New(),
		GroupLetter:          &g,
		Round:                "Group Stage - 1",
		Status:               domain.FixtureStatusInProgress,
		KickoffAt:            time.Now().Add(-45 * time.Minute),
	}
}

func finishedFixture(round string) domain.PollableFixture {
	f := liveFixture()
	f.Status = domain.FixtureStatusFinished
	f.Round = round
	f.GroupLetter = nil
	return f
}

// ---- tests ----

func TestMapAPIStatus(t *testing.T) {
	cases := []struct {
		short string
		want  domain.FixtureStatus
	}{
		{"NS", domain.FixtureStatusUpcoming},
		{"PST", domain.FixtureStatusUpcoming},
		{"1H", domain.FixtureStatusInProgress},
		{"HT", domain.FixtureStatusInProgress},
		{"2H", domain.FixtureStatusInProgress},
		{"ET", domain.FixtureStatusInProgress},
		{"BT", domain.FixtureStatusInProgress},
		{"P", domain.FixtureStatusInProgress},
		{"SUSP", domain.FixtureStatusInProgress},
		{"INT", domain.FixtureStatusInProgress},
		{"FT", domain.FixtureStatusFinished},
		{"AET", domain.FixtureStatusFinished},
		{"PEN", domain.FixtureStatusFinished},
		{"AWD", domain.FixtureStatusFinished},
		{"WO", domain.FixtureStatusFinished},
		{"CANC", domain.FixtureStatusCancelled},
		{"ABD", domain.FixtureStatusCancelled},
		{"UNKNOWN", domain.FixtureStatusUpcoming},
	}
	for _, tc := range cases {
		got := worker.MapAPIStatus(tc.short)
		if got != tc.want {
			t.Errorf("MapAPIStatus(%q) = %q, want %q", tc.short, got, tc.want)
		}
	}
}

func TestProcessSingleFixture_NoChange(t *testing.T) {
	gh, ga := 1, 0
	f := liveFixture()
	f.GoalsHome = &gh
	f.GoalsAway = &ga
	f.Status = domain.FixtureStatusInProgress

	repo := newFakeRepo(nil)
	api := &fakeAPI{result: worker.APIFixtureResult{
		StatusShort: "2H",
		GoalsHome:   &gh,
		GoalsAway:   &ga,
	}}
	w := newWorker(repo, api)
	w.Run(cancelledCtx()) //nolint:errcheck

	if repo.callCount("UpdateMatchAndRescoreLivePredictions") != 0 {
		t.Error("expected no DB write when state unchanged")
	}
}

func TestProcessSingleFixture_ScoreChange(t *testing.T) {
	old := 0
	newG := 1
	f := liveFixture()
	f.GoalsHome = &old
	f.GoalsAway = &old
	f.Status = domain.FixtureStatusInProgress

	repo := newFakeRepo([]domain.PollableFixture{f})
	api := &fakeAPI{result: worker.APIFixtureResult{
		StatusShort: "2H",
		GoalsHome:   &newG,
		GoalsAway:   &old,
	}}
	w := newWorker(repo, api)
	runOneTick(w)

	if repo.callCount("UpdateMatchAndRescoreLivePredictions") != 1 {
		t.Error("expected DB write on score change")
	}
}

func TestProcessSingleFixture_APIError(t *testing.T) {
	f := liveFixture()
	repo := newFakeRepo([]domain.PollableFixture{f})
	api := &fakeAPI{err: errors.New("timeout")}
	w := newWorker(repo, api)
	runOneTick(w)

	if repo.callCount("UpdateMatchAndRescoreLivePredictions") != 0 {
		t.Error("expected no DB write when API errors")
	}
}

func TestRunSettlement_Cancelled(t *testing.T) {
	f := liveFixture()
	repo := newFakeRepo([]domain.PollableFixture{f})
	api := &fakeAPI{result: worker.APIFixtureResult{StatusShort: "CANC"}}
	w := newWorker(repo, api)
	runOneTick(w)

	if repo.callCount("SettleGroupWinnerPredictions") != 0 {
		t.Error("expected no outright settlement for cancelled match")
	}
	if repo.callCount("UpdateMatchAndRescoreLivePredictions") != 1 {
		t.Error("expected score predictions to be zeroed on cancel")
	}
}

func TestRunSettlement_GroupMatchGroupNotDone(t *testing.T) {
	f := liveFixture() // has GroupLetter = "A"
	repo := newFakeRepo([]domain.PollableFixture{f})
	repo.groupComplete = false
	api := &fakeAPI{result: worker.APIFixtureResult{
		StatusShort: "FT",
		GoalsHome:   intPtr(1),
		GoalsAway:   intPtr(0),
	}}
	w := newWorker(repo, api)
	runOneTick(w)

	if repo.callCount("SettleGroupWinnerPredictions") != 0 {
		t.Error("must not settle group winner when group not done")
	}
	if repo.callCount("UpdateGroupStandings") != 1 {
		t.Error("standings should still be updated")
	}
}

func TestRunSettlement_GroupMatchGroupDone(t *testing.T) {
	f := liveFixture()
	repo := newFakeRepo([]domain.PollableFixture{f})
	repo.groupComplete = true
	repo.groupStageComplete = false
	api := &fakeAPI{
		result: worker.APIFixtureResult{
			StatusShort: "FT",
			GoalsHome:   intPtr(2),
			GoalsAway:   intPtr(1),
		},
		topScorer: worker.APITopScorerResult{PlayerExternalID: "999", Goals: 3},
	}
	w := newWorker(repo, api)
	runOneTick(w)

	if repo.callCount("SettleGroupWinnerPredictions") != 1 {
		t.Error("expected SettleGroupWinnerPredictions")
	}
	if repo.callCount("SettlePlayoffGroupPredictions") != 1 {
		t.Error("expected SettlePlayoffGroupPredictions")
	}
	if repo.callCount("SettleGroupTopScorerPredictions") != 1 {
		t.Error("expected SettleGroupTopScorerPredictions")
	}
	if repo.callCount("SettlePlayoffWildcardPredictions") != 0 {
		t.Error("must not settle wildcard when group stage not complete")
	}
}

func TestRunSettlement_LastGroupDone(t *testing.T) {
	f := liveFixture()
	repo := newFakeRepo([]domain.PollableFixture{f})
	repo.groupComplete = true
	repo.groupStageComplete = true
	api := &fakeAPI{
		result: worker.APIFixtureResult{
			StatusShort: "FT",
			GoalsHome:   intPtr(1),
			GoalsAway:   intPtr(0),
		},
		topScorer: worker.APITopScorerResult{PlayerExternalID: "1", Goals: 5},
	}
	w := newWorker(repo, api)
	runOneTick(w)

	if repo.callCount("SettlePlayoffWildcardPredictions") != 1 {
		t.Error("expected SettlePlayoffWildcardPredictions when all groups done")
	}
}

func TestRunSettlement_QFNotAllDone(t *testing.T) {
	f := finishedFixture("Quarter-finals")
	repo := newFakeRepo([]domain.PollableFixture{f})
	repo.roundComplete = false
	hw := true
	api := &fakeAPI{result: worker.APIFixtureResult{
		StatusShort: "FT",
		HomeWinner:  &hw,
		GoalsHome:   intPtr(2),
		GoalsAway:   intPtr(1),
	}}
	w := newWorker(repo, api)
	runOneTick(w)

	if repo.callCount("SettleSemifinalistPredictions") != 0 {
		t.Error("must not settle semifinalists when round not complete")
	}
}

func TestRunSettlement_AllQFsDone(t *testing.T) {
	f := finishedFixture("Quarter-finals")
	repo := newFakeRepo([]domain.PollableFixture{f})
	repo.roundComplete = true
	hw := true
	api := &fakeAPI{result: worker.APIFixtureResult{
		StatusShort: "FT",
		HomeWinner:  &hw,
		GoalsHome:   intPtr(2),
		GoalsAway:   intPtr(0),
	}}
	w := newWorker(repo, api)
	runOneTick(w)

	if repo.callCount("SettleSemifinalistPredictions") != 1 {
		t.Error("expected SettleSemifinalistPredictions when all QFs done")
	}
}

func TestRunSettlement_FinalConcluded(t *testing.T) {
	f := finishedFixture("Final")
	repo := newFakeRepo([]domain.PollableFixture{f})
	repo.roundComplete = false // not a QF round so irrelevant
	hw := true
	api := &fakeAPI{
		result: worker.APIFixtureResult{
			StatusShort: "FT",
			HomeWinner:  &hw,
			GoalsHome:   intPtr(3),
			GoalsAway:   intPtr(1),
		},
		topScorer: worker.APITopScorerResult{PlayerExternalID: "42", Goals: 7},
	}
	w := newWorker(repo, api)
	runOneTick(w)

	if repo.callCount("SettleTournamentWinnerPredictions") != 1 {
		t.Error("expected SettleTournamentWinnerPredictions")
	}
	if repo.callCount("SettleTopScorerPredictions") != 1 {
		t.Error("expected SettleTopScorerPredictions")
	}
}

func TestRunSettlement_FinalNilWinner(t *testing.T) {
	f := finishedFixture("Final")
	repo := newFakeRepo([]domain.PollableFixture{f})
	api := &fakeAPI{result: worker.APIFixtureResult{
		StatusShort: "FT",
		GoalsHome:   intPtr(1),
		GoalsAway:   intPtr(1),
		// HomeWinner and AwayWinner are nil
	}}
	w := newWorker(repo, api)
	runOneTick(w)

	if repo.callCount("SettleTournamentWinnerPredictions") != 0 {
		t.Error("must not settle winner when winnerTeamID is nil")
	}
}

func TestRun_ContextCancelled(t *testing.T) {
	repo := newFakeRepo(nil)
	api := &fakeAPI{}
	w := newWorker(repo, api)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := w.Run(ctx)
	if err != nil {
		t.Errorf("Run should return nil on context cancel, got %v", err)
	}
}

func TestRun_ListPollableMatchesError(t *testing.T) {
	repo := newFakeRepo(nil)
	repo.listErr = errors.New("db error")
	api := &fakeAPI{}
	w := newWorker(repo, api)

	runOneTick(w)

	if atomic.LoadInt64(&api.callCount) != 0 {
		t.Error("API should not be called when list fails")
	}
}

func TestTick_LiveFixtures_SemaphoreLimits(t *testing.T) {
	fixtures := make([]domain.PollableFixture, 10)
	for i := range fixtures {
		fixtures[i] = liveFixture()
		gh, ga := i, 0
		fixtures[i].GoalsHome = &gh
		fixtures[i].GoalsAway = &ga
		fixtures[i].Status = domain.FixtureStatusInProgress
	}

	// Just ensure all fixtures are polled.
	repo2 := newFakeRepo(fixtures)
	api2 := &fakeAPI{result: worker.APIFixtureResult{StatusShort: "1H"}} // same status → no change (in_progress same as fixture)
	w := newWorker(repo2, api2)
	runOneTick(w)

	if atomic.LoadInt64(&api2.callCount) != int64(len(fixtures)) {
		t.Errorf("expected %d API calls, got %d", len(fixtures), atomic.LoadInt64(&api2.callCount))
	}
}

// ---- helpers ----

func cancelledCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func runOneTick(w *worker.Worker) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Run worker until context cancels; this triggers one tick then exits.
	doneCh := make(chan struct{})
	go func() {
		w.Run(ctx) //nolint:errcheck
		close(doneCh)
	}()
	// Cancel after giving the tick time to run.
	time.Sleep(200 * time.Millisecond)
	cancel()
	<-doneCh
}

func intPtr(v int) *int { return &v }

func nopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
