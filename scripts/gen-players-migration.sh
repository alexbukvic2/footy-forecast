#!/usr/bin/env bash
#
# gen-players-migration.sh — fetch squads from API-Football and emit a goose migration
#                            that seeds the players table for every team in a tournament.
#
# For each team in the given tournament the script:
#   1. Calls /players/squads?team={external_id} to enumerate player IDs.
#   2. Calls /players/profiles?player={id} for each player to get firstname + lastname.
#   3. Emits a single idempotent INSERT … ON CONFLICT migration file.
#
# Usage:
#   API_KEY=xxx \
#   TOURNAMENT_ID=7006468e-21cb-433a-a98e-d6b0d51fe0d7 \
#   TOURNAMENT_SLUG=world-cup-2026 \
#   DB_URL="postgres://user:pass@localhost/footy_forecast" \
#   ./scripts/gen-players-migration.sh
#
# Optional:
#   OUT_DIR=./migrations   (default: ./migrations)
#   SLUG=worldcup2026      (appended to migration filename; default: tournament_<short-id>)
#   RATE_SLEEP=1           (seconds to sleep between API calls; default: 1)

set -euo pipefail

: "${API_KEY:?API_KEY is required}"
: "${TOURNAMENT_ID:?TOURNAMENT_ID is required}"
: "${TOURNAMENT_SLUG:?TOURNAMENT_SLUG is required}"
: "${DB_URL:?DB_URL is required (used to read team external_ids)}"

OUT_DIR="${OUT_DIR:-./migrations}"
RATE_SLEEP="${RATE_SLEEP:-1}"
SHORT_ID="${TOURNAMENT_ID:0:8}"
SLUG="${SLUG:-tournament_${SHORT_ID}}"

for cmd in jq curl psql; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "error: $cmd is required" >&2
    exit 1
  fi
done

# ────────────────────────────────────────────────
# 1. Read team external_ids from DB
# ────────────────────────────────────────────────
echo "Reading teams for tournament ${TOURNAMENT_ID}…"

TEAMS_JSON=$(psql "$DB_URL" -t -A -c \
  "SELECT json_agg(json_build_object('internal_id', id, 'external_id', external_id, 'name', name))
   FROM teams
   WHERE tournament_id = '${TOURNAMENT_ID}';")

TEAM_COUNT=$(echo "$TEAMS_JSON" | jq 'length')
if [[ "$TEAM_COUNT" -eq 0 ]]; then
  echo "error: no teams found for tournament ${TOURNAMENT_ID}" >&2
  exit 1
fi
echo "Found ${TEAM_COUNT} teams."

# ────────────────────────────────────────────────
# 2. Fetch squads + profiles, accumulate rows
# ────────────────────────────────────────────────

# Temporary files to accumulate INSERT statements
TMP_INSERTS=$(mktemp)
trap 'rm -f "$TMP_INSERTS"' EXIT

TOTAL_PLAYERS=0

while IFS= read -r team; do
  INTERNAL_ID=$(echo "$team" | jq -r '.internal_id')
  EXTERNAL_ID=$(echo "$team" | jq -r '.external_id')
  TEAM_NAME=$(echo "$team"   | jq -r '.name')

  echo "  Team: ${TEAM_NAME} (external_id=${EXTERNAL_ID})…"

  SQUAD_RESP=$(curl -fsSL \
    -H "x-apisports-key: ${API_KEY}" \
    "https://v3.football.api-sports.io/players/squads?team=${EXTERNAL_ID}")

  ERRORS=$(echo "$SQUAD_RESP" | jq '.errors | length')
  if [[ "$ERRORS" -gt 0 ]]; then
    echo "    API error for squad (team=${EXTERNAL_ID}):" >&2
    echo "$SQUAD_RESP" | jq '.errors' >&2
    continue
  fi

  PLAYER_IDS=$(echo "$SQUAD_RESP" | jq -r '.response[0].players[].id // empty')
  PLAYER_COUNT=$(echo "$PLAYER_IDS" | grep -c '[0-9]' || true)
  echo "    ${PLAYER_COUNT} players in squad."

  while IFS= read -r PLAYER_ID; do
    [[ -z "$PLAYER_ID" ]] && continue

    sleep "$RATE_SLEEP"

    PROFILE_RESP=$(curl -fsSL \
      -H "x-apisports-key: ${API_KEY}" \
      "https://v3.football.api-sports.io/players/profiles?player=${PLAYER_ID}")

    PERRORS=$(echo "$PROFILE_RESP" | jq '.errors | length')
    if [[ "$PERRORS" -gt 0 ]]; then
      echo "    API error for player ${PLAYER_ID}, skipping." >&2
      continue
    fi

    # Compose full name from firstname + lastname; fall back to .name if missing
    FULL_NAME=$(echo "$PROFILE_RESP" | jq -r '
      .response[0].player |
      if (.firstname != null and .lastname != null)
      then (.firstname + " " + .lastname)
      else .name
      end
    ')

    if [[ -z "$FULL_NAME" || "$FULL_NAME" == "null" ]]; then
      echo "    No name for player ${PLAYER_ID}, skipping." >&2
      continue
    fi

    # Escape single quotes for SQL
    SAFE_NAME="${FULL_NAME//\'/\'\'}"
    SAFE_INTERNAL_ID="$INTERNAL_ID"

    cat >> "$TMP_INSERTS" <<SQL
INSERT INTO players (external_id, name, tournament_id, team_id)
VALUES (
    ${PLAYER_ID},
    '${SAFE_NAME}',
    '${TOURNAMENT_ID}',
    '${SAFE_INTERNAL_ID}'
) ON CONFLICT (external_id, tournament_id) DO UPDATE SET
    name    = EXCLUDED.name,
    team_id = EXCLUDED.team_id;
SQL

    TOTAL_PLAYERS=$(( TOTAL_PLAYERS + 1 ))
  done <<< "$PLAYER_IDS"

  sleep "$RATE_SLEEP"
done < <(echo "$TEAMS_JSON" | jq -c '.[]')

echo "Fetched ${TOTAL_PLAYERS} player profiles."

# ────────────────────────────────────────────────
# 3. Write migration file
# ────────────────────────────────────────────────
TIMESTAMP=$(date -u +"%Y%m%d%H%M%S")
OUT_FILE="${OUT_DIR}/${TIMESTAMP}_players_seed_${SLUG}.sql"

cat > "$OUT_FILE" <<SQL
-- +goose Up
-- +goose StatementBegin

-- Seeded from API-Football for tournament_id=${TOURNAMENT_ID} (${TOURNAMENT_SLUG})
-- Idempotent: ON CONFLICT (external_id, tournament_id) updates name + team_id.

SQL

cat "$TMP_INSERTS" >> "$OUT_FILE"

cat >> "$OUT_FILE" <<SQL

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM players WHERE tournament_id = '${TOURNAMENT_ID}';

-- +goose StatementEnd
SQL

echo "Written: ${OUT_FILE}"
echo
echo "Next steps:"
echo "  1. Review ${OUT_FILE}."
echo "  2. Commit the file — goose will apply it on next deploy."
