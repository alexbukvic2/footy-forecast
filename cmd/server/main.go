// Command server runs the footy-forecast HTTP API.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexbukvic2/footy-forecast/internal/bedrock"
	"github.com/alexbukvic2/footy-forecast/internal/footballapi"
	"github.com/alexbukvic2/footy-forecast/internal/notification/expo"
	"github.com/alexbukvic2/footy-forecast/internal/notification/job"
	"github.com/alexbukvic2/footy-forecast/internal/repository"
	"github.com/alexbukvic2/footy-forecast/internal/worker"
	"github.com/joho/godotenv"

	"github.com/alexbukvic2/footy-forecast/internal/config"
	"github.com/alexbukvic2/footy-forecast/internal/db"
	"github.com/alexbukvic2/footy-forecast/internal/server"
)

func main() {
	// Load .env in dev. Silently ignore if missing — in prod, env comes from systemd.
	_ = godotenv.Load()

	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		logger.Error("fatal", "err", err)
		os.Exit(1)
	}
}

// startNotificationJobs launches the three notification background workers. Each
// runs on its own goroutine and follows the same Run(ctx) error pattern as worker.Worker.
func startNotificationJobs(ctx context.Context, pool *db.Pool, logger *slog.Logger) {
	expoClient := expo.NewClient()
	notifiers := []struct {
		name string
		run  func(context.Context) error
	}{
		{"pre_match notifier", job.NewPreMatchNotifier(pool, expoClient, logger).Run},
		{"matchday notifier", job.NewMatchdayNotifier(pool, expoClient, logger).Run},
		{"tournament_reminder notifier", job.NewTournamentReminderNotifier(pool, expoClient, logger).Run},
	}
	for _, n := range notifiers {
		go func() {
			if err := n.run(ctx); err != nil && !errors.Is(err, context.Canceled) {
				logger.Error(n.name+" stopped unexpectedly", "err", err)
			}
		}()
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()
	logger.Info("database connected")

	analyser, err := bedrock.NewClient(cfg.BedrockRegion, cfg.BedrockModelID)
	if err != nil {
		return fmt.Errorf("create bedrock analyser: %w", err)
	}

	workerRepo := repository.NewWorkerRepository(pool)
	apiClient := footballapi.NewClient(cfg.FootballAPIKey, cfg.FootballAPIBaseURL, nil)
	w := worker.New(
		workerRepo,
		apiClient,
		worker.RealClock{},
		logger,
		cfg.WorkerPollInterval,
		cfg.PredictionLockLeadMinutes,
		analyser,
	)

	workerErr := make(chan error, 1)
	go func() {
		if err := w.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			workerErr <- err
		}
		close(workerErr)
	}()

	// Notification jobs.
	startNotificationJobs(ctx, pool, logger)

	router := server.NewRouter(logger, pool, cfg)

	s := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("server starting", "addr", s.Addr, "env", cfg.Env)
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-serverErr:
		return fmt.Errorf("server: %w", err)
	case err := <-workerErr:
		return fmt.Errorf("worker: %w", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	logger.Info("server stopped cleanly")
	return nil
}
