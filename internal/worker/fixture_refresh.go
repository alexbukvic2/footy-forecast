package worker

import (
	"context"

	"github.com/google/uuid"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

// refreshAllActiveTournamentFixtures fetches and inserts missing fixtures for
// every tournament that has an external API ID. Called once at startup.
func (w *Worker) refreshAllActiveTournamentFixtures(ctx context.Context) {
	tournaments, err := w.repo.ListActiveTournaments(ctx)
	if err != nil {
		w.logger.Warn("worker: list active tournaments for fixture refresh", "err", err)
		return
	}
	for _, t := range tournaments {
		w.refreshFixturesForTournament(ctx, t.ExternalID, t.Season, t.ID)
	}
}

// refreshFixturesForTournament fetches the full fixture list for a tournament
// from the API and inserts any fixtures not yet in the DB.
func (w *Worker) refreshFixturesForTournament(
	ctx context.Context,
	externalLeagueID int64,
	season int,
	tournamentID uuid.UUID,
) {
	apiFixtures, err := w.api.GetLeagueFixtures(ctx, externalLeagueID, season)
	if err != nil {
		w.logger.Warn("worker: get league fixtures", "league_id", externalLeagueID, "season", season, "err", err)
		return
	}

	toInsert := make([]domain.NewFixture, 0, len(apiFixtures))
	for _, af := range apiFixtures {
		homeID, err := w.repo.GetTeamByExternalID(ctx, af.HomeTeamExternalID, tournamentID)
		if err != nil {
			w.logger.Warn("worker: resolve home team for fixture refresh",
				"external_id", af.HomeTeamExternalID, "fixture_external_id", af.ExternalID, "err", err)
			continue
		}
		awayID, err := w.repo.GetTeamByExternalID(ctx, af.AwayTeamExternalID, tournamentID)
		if err != nil {
			w.logger.Warn("worker: resolve away team for fixture refresh",
				"external_id", af.AwayTeamExternalID, "fixture_external_id", af.ExternalID, "err", err)
			continue
		}
		toInsert = append(toInsert, domain.NewFixture{
			ExternalID: af.ExternalID,
			HomeTeamID: homeID,
			AwayTeamID: awayID,
			KickoffAt:  af.KickoffAt,
			Status:     MapAPIStatus(af.StatusShort),
			Round:      af.Round,
			GoalsHome:  af.GoalsHome,
			GoalsAway:  af.GoalsAway,
		})
	}

	if err := w.repo.InsertMissingFixtures(ctx, tournamentID, toInsert); err != nil {
		w.logger.Error("worker: insert missing fixtures", "league_id", externalLeagueID, "err", err)
	}
}
