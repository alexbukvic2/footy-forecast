#!/usr/bin/env bash
# Wipes the local dev database and resets the goose version table.
# Drops the public schema (all tables, types, extensions) and recreates it.
# After running this, apply migrations with: make migrate-up
set -euo pipefail

DB_URL="${DATABASE_URL:-postgres://footy:footy_dev_password@localhost:5432/footy_forecast?sslmode=disable}"

echo "Cleaning database: $DB_URL"

psql "$DB_URL" <<'SQL'
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO PUBLIC;
SQL

echo "Done. Run 'make migrate-up' to apply all migrations."
