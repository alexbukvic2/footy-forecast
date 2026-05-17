//go:build integration

package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/alexbukvic2/footy-forecast/internal/db"
)

// startPostgres starts a fresh Postgres container, runs all migrations against it,
// and returns a connected pool. The container is terminated when the test ends.
//
// Each call gets its own container so tests are fully isolated. Slower per test,
// but no cross-test state leakage.
func startPostgres(t *testing.T) *db.Pool {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	container, err := tcpostgres.Run(
		ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("footy_forecast_test"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err, "start postgres container")

	t.Cleanup(
		func() {
			// Use a fresh context — t.Cleanup may run after the test context expired.
			stopCtx, stopCancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer stopCancel()
			if err := container.Terminate(stopCtx); err != nil {
				t.Logf("terminate postgres container: %v", err)
			}
		},
	)

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	// Run migrations.
	require.NoError(t, runMigrations(dsn), "run migrations")

	pool, err := db.New(ctx, dsn)
	require.NoError(t, err, "connect pool")

	t.Cleanup(pool.Close)

	return pool
}

// runMigrations applies all goose migrations against the given DSN.
//
// Uses goose's Provider API rather than the legacy global-state functions
// (SetDialect, Up) — the legacy ones share package-level state and race
// when called from parallel tests.
func runMigrations(dsn string) error {
	// Find the migrations dir relative to this file.
	_, here, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(here), "..", "..", "migrations")

	dbConn, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer dbConn.Close()

	provider, err := goose.NewProvider(goose.DialectPostgres, dbConn, os.DirFS(migrationsDir))
	if err != nil {
		return fmt.Errorf("new provider: %w", err)
	}

	if _, err := provider.Up(context.Background()); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

// truncate clears the given tables. Useful at the start of each test if you
// don't want a fresh container.
func truncate(
	t *testing.T,
	pool *db.Pool,
	tables ...string,
) {
	t.Helper()
	for _, table := range tables {
		_, err := pool.Exec(context.Background(), fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		require.NoError(t, err, "truncate %s", table)
	}
}
