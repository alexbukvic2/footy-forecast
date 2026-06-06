#!/usr/bin/env bash
#
# gen-teams-migration.sh — fetch teams from API-Football and emit two goose migration files:
#   1. teams_seed        — INSERT tournament + INSERT teams (idempotent)
#   2. team_handicap_seed — template; replace NULL with points, leave NULL to skip
#
# Usage:
#   API_KEY=xxx \
#   LEAGUE_ID=1 \
#   SEASON=2026 \
#   TOURNAMENT_NAME="FIFA World Cup 2026" \
#   TOURNAMENT_SLUG="world-cup-2026" \
#   TOURNAMENT_STARTS_AT="2026-06-11T00:00:00Z" \
#   TOURNAMENT_ENDS_AT="2026-07-19T00:00:00Z" \
#   ./scripts/gen-teams-migration.sh
#
# Optional:
#   OUT_DIR=./migrations   (default: ./migrations)
#   SLUG=worldcup2026      (appended to migration filenames; default: league<LEAGUE_ID>)

set -euo pipefail

: "${API_KEY:?API_KEY is required}"
: "${LEAGUE_ID:?LEAGUE_ID is required}"
: "${SEASON:?SEASON is required}"
: "${TOURNAMENT_NAME:?TOURNAMENT_NAME is required}"
: "${TOURNAMENT_SLUG:?TOURNAMENT_SLUG is required}"
: "${TOURNAMENT_STARTS_AT:?TOURNAMENT_STARTS_AT is required}"
: "${TOURNAMENT_ENDS_AT:?TOURNAMENT_ENDS_AT is required}"

OUT_DIR="${OUT_DIR:-./migrations}"
SLUG="${SLUG:-league${LEAGUE_ID}}"

if ! command -v jq &>/dev/null; then
  echo "error: jq is required (brew install jq)" >&2
  exit 1
fi

echo "Fetching teams for league=${LEAGUE_ID} season=${SEASON}…"

RESPONSE=$(curl -fsSL \
  -H "x-apisports-key: ${API_KEY}" \
  "https://v3.football.api-sports.io/teams?league=${LEAGUE_ID}&season=${SEASON}")

ERRORS=$(echo "$RESPONSE" | jq '.errors | length')
if [[ "$ERRORS" -gt 0 ]]; then
  echo "API returned errors:" >&2
  echo "$RESPONSE" | jq '.errors' >&2
  exit 1
fi

TEAM_COUNT=$(echo "$RESPONSE" | jq '.results')
echo "Got ${TEAM_COUNT} teams."

TIMESTAMP=$(date -u +"%Y%m%d%H%M%S")
TIMESTAMP_NEXT=$(date -u -v+1S +"%Y%m%d%H%M%S" 2>/dev/null || date -u -d "+1 second" +"%Y%m%d%H%M%S")
TEAMS_FILE="${OUT_DIR}/${TIMESTAMP}_teams_seed_${SLUG}.sql"
HANDICAP_FILE="${OUT_DIR}/${TIMESTAMP_NEXT}_team_handicap_seed_${SLUG}.sql"

# ────────────────────────────────────────────────
# 1. teams_seed migration
# ────────────────────────────────────────────────
cat > "$TEAMS_FILE" <<SQL
-- +goose Up
-- +goose StatementBegin

-- Seeded from API-Football league=${LEAGUE_ID} season=${SEASON}
-- Tournament and teams are inserted idempotently:
--   ON CONFLICT (slug)  DO NOTHING  for tournaments
--   ON CONFLICT (name)  DO UPDATE   for teams (updates external_id, logo, tournament_id)

INSERT INTO tournaments (id, slug, name, external_id, season, starts_at, ends_at)
VALUES (
    gen_random_uuid(),
    '${TOURNAMENT_SLUG}',
    '${TOURNAMENT_NAME}',
    ${LEAGUE_ID},
    ${SEASON},
    '${TOURNAMENT_STARTS_AT}',
    '${TOURNAMENT_ENDS_AT}'
) ON CONFLICT (slug) DO NOTHING;

SQL

echo "$RESPONSE" | jq -r --arg slug "${TOURNAMENT_SLUG}" '
  .response[] |
  .team |
  "INSERT INTO teams (external_id, name, logo, tournament_id)\n" +
  "VALUES (\n" +
  "    " + (.id | tostring) + ",\n" +
  "    '"'"'" + (.name | gsub("'"'"'"; "'"'"''"'"'")) + "'"'"',\n" +
  "    '"'"'" + (.logo | gsub("'"'"'"; "'"'"''"'"'")) + "'"'"',\n" +
  "    (SELECT id FROM tournaments WHERE slug = '"'"'" + $slug + "'"'"')\n" +
  ") ON CONFLICT (name) DO UPDATE SET\n" +
  "    external_id   = EXCLUDED.external_id,\n" +
  "    logo          = EXCLUDED.logo,\n" +
  "    tournament_id = EXCLUDED.tournament_id;\n"
' >> "$TEAMS_FILE"

cat >> "$TEAMS_FILE" <<SQL

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Removing seeded data cascades to fixtures, predictions, handicaps, etc.
-- Run only if you are sure this is safe in the target environment.
DELETE FROM teams      WHERE tournament_id = (SELECT id FROM tournaments WHERE slug = '${TOURNAMENT_SLUG}');
DELETE FROM tournaments WHERE slug = '${TOURNAMENT_SLUG}';

-- +goose StatementEnd
SQL

echo "Written: ${TEAMS_FILE}"

# ────────────────────────────────────────────────
# 2. team_handicap_seed migration template
# ────────────────────────────────────────────────
cat > "$HANDICAP_FILE" <<'HEADER'
-- +goose Up
-- +goose StatementBegin

-- Team handicap seed template.
--
-- For each team, replace NULL with the integer point value for that category.
-- Leave NULL to skip that team/category pair — it will NOT be inserted.
-- Idempotent: ON CONFLICT updates the existing row's points.
--
-- Categories:
--   group_winner   — points for predicting the team wins its group
--   playoff        — points for predicting the team reaches the playoff round
--   semifinalist   — points for predicting the team reaches the semi-finals
--   winner         — points for predicting the team wins the tournament

HEADER

echo "$RESPONSE" | jq -r '
  .response[] |
  .team |
  "-- " + .name + " (" + (.code // "???") + ")\n" +
  "INSERT INTO team_handicap (team_id, category, points)\n" +
  "SELECT t.id, v.category, v.points\n" +
  "FROM   teams t\n" +
  "JOIN   (VALUES\n" +
  "    ('\''group_winner'\''::team_handicap_category,  NULL::int),\n" +
  "    ('\''playoff'\''::team_handicap_category,       NULL::int),\n" +
  "    ('\''semifinalist'\''::team_handicap_category,  NULL::int),\n" +
  "    ('\''winner'\''::team_handicap_category,        NULL::int)\n" +
  ") AS v(category, points) ON TRUE\n" +
  "WHERE  t.name = '"'"'" + (.name | gsub("'"'"'"; "'"'"''"'"'")) + "'"'"'\n" +
  "  AND  v.points IS NOT NULL\n" +
  "ON CONFLICT (team_id, category) DO UPDATE SET points = EXCLUDED.points;\n"
' >> "$HANDICAP_FILE"

cat >> "$HANDICAP_FILE" <<'FOOTER'
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- No automatic rollback for handicap data — delete manually if needed.

-- +goose StatementEnd
FOOTER

echo "Written: ${HANDICAP_FILE}"
echo
echo "Next steps:"
echo "  1. Review ${TEAMS_FILE} — tournament + teams are inserted idempotently."
echo "  2. Open ${HANDICAP_FILE} and replace NULL with points for each team/category."
echo "     Leave NULL to skip a team/category pair entirely."
echo "  3. Commit both files and push — goose will apply them on next deploy."
