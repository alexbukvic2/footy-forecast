#!/usr/bin/env bash
#
# db-pull-prod.sh — Dump the production database via SSM and restore it locally.
#
# What it does:
#   1. Runs pg_dump on the EC2 box via SSM Run Command.
#   2. Uploads the gzipped dump to a temporary S3 object.
#   3. Downloads the dump locally.
#   4. Wipes the local database (DROP/CREATE SCHEMA public).
#   5. Restores via psql.
#   6. Deletes the temporary S3 object.
#
# Prerequisites:
#   - AWS CLI with profile "hexa" (or set AWS_PROFILE)
#   - The EC2 instance must be running
#   - psql available locally
#
# Usage:
#   ./scripts/db-pull-prod.sh
#   AWS_PROFILE=hexa ./scripts/db-pull-prod.sh
#
set -euo pipefail

AWS_PROFILE="${AWS_PROFILE:-hexa}"
REGION="${AWS_REGION:-eu-central-1}"
INSTANCE_TAG="${INSTANCE_TAG:-footy-forecast}"
BACKUP_BUCKET="${BACKUP_BUCKET:-footy-forecast-backups}"
LOCAL_DB_URL="${DATABASE_URL:-postgres://footy:footy_dev_password@localhost:5432/footy_forecast?sslmode=disable}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() { printf "${GREEN}==>${NC} %s\n" "$*"; }
warn() { printf "${YELLOW}!!!${NC} %s\n" "$*" >&2; }
fail() { printf "${RED}ERR${NC} %s\n" "$*" >&2; exit 1; }

# Cleanup on exit so the local dump and S3 temp file don't linger on failure
DUMP_S3_KEY=""
LOCAL_DUMP=""
cleanup() {
  [ -n "$LOCAL_DUMP" ] && rm -f "$LOCAL_DUMP"
  if [ -n "$DUMP_S3_KEY" ]; then
    aws s3 rm "s3://${BACKUP_BUCKET}/${DUMP_S3_KEY}" \
      --profile "$AWS_PROFILE" --region "$REGION" 2>/dev/null || true
  fi
}
trap cleanup EXIT

# --- safety prompt -----------------------------------------------------------

warn "This will WIPE your local database and replace it with production data."
warn "Local DB: $LOCAL_DB_URL"
read -r -p "Continue? [y/N] " CONFIRM
CONFIRM_LOWER=$(echo "$CONFIRM" | tr '[:upper:]' '[:lower:]')
[[ "$CONFIRM_LOWER" == "y" ]] || { echo "Aborted."; exit 0; }

# --- 1. Find EC2 instance ----------------------------------------------------

info "Looking up EC2 instance (tag Name=${INSTANCE_TAG})"

INSTANCE_ID=$(aws ec2 describe-instances \
  --filters \
    "Name=tag:Name,Values=${INSTANCE_TAG}" \
    "Name=instance-state-name,Values=running" \
  --query 'Reservations[0].Instances[0].InstanceId' \
  --output text \
  --profile "$AWS_PROFILE" \
  --region "$REGION")

if [ "$INSTANCE_ID" = "None" ] || [ -z "$INSTANCE_ID" ]; then
  fail "No running EC2 instance found with tag Name=${INSTANCE_TAG}. Is the instance started?"
fi

info "Found instance: $INSTANCE_ID"

# --- 2. Run pg_dump on EC2 and upload to S3 ----------------------------------

TIMESTAMP=$(date -u +%Y%m%d-%H%M%S)
DUMP_S3_KEY="tmp/prod-dump-${TIMESTAMP}.sql.gz"
LOCAL_DUMP="/tmp/prod-dump-${TIMESTAMP}.sql.gz"

info "Running pg_dump on $INSTANCE_ID -> s3://${BACKUP_BUCKET}/${DUMP_S3_KEY}"

# pg_dump flags:
#   --no-owner  — don't emit ALTER OWNER (prod uses footy_app; local uses footy)
#   --no-acl    — don't emit GRANT/REVOKE (same reason)
#   -Fp         — plain SQL output (pipe-friendly)
#
# The SSM parameters JSON is built with python3 to avoid shell quoting issues.
SSM_PARAMS=$(python3 - <<PYEOF
import json
bucket = "${BACKUP_BUCKET}"
key    = "${DUMP_S3_KEY}"
cmds = [
    "#!/bin/bash",
    "set -euo pipefail",
    "export PGPASSWORD=\$(aws ssm get-parameter --name /footy-forecast/prod/postgres-password --with-decryption --region eu-central-1 --query 'Parameter.Value' --output text)",
    f"pg_dump -h 127.0.0.1 -U footy_app -d footy_forecast --no-owner --no-acl -Fp | gzip | aws s3 cp - s3://{bucket}/{key} --region eu-central-1",
]
print(json.dumps({"commands": cmds}))
PYEOF
)

COMMAND_ID=$(aws ssm send-command \
  --instance-ids "$INSTANCE_ID" \
  --document-name "AWS-RunShellScript" \
  --parameters "$SSM_PARAMS" \
  --timeout-seconds 300 \
  --comment "db-pull-prod: pg_dump to S3" \
  --profile "$AWS_PROFILE" \
  --region "$REGION" \
  --query 'Command.CommandId' \
  --output text)

info "SSM command dispatched: $COMMAND_ID"
info "Waiting for pg_dump to complete (up to 5 minutes)..."

for i in $(seq 1 60); do
  STATUS=$(aws ssm get-command-invocation \
    --command-id "$COMMAND_ID" \
    --instance-id "$INSTANCE_ID" \
    --profile "$AWS_PROFILE" \
    --region "$REGION" \
    --query 'Status' \
    --output text 2>/dev/null || echo "Pending")

  case "$STATUS" in
    Success)
      info "pg_dump complete"
      break
      ;;
    Failed|Cancelled|TimedOut|Undeliverable|Terminated)
      STDERR=$(aws ssm get-command-invocation \
        --command-id "$COMMAND_ID" \
        --instance-id "$INSTANCE_ID" \
        --profile "$AWS_PROFILE" \
        --region "$REGION" \
        --query 'StandardErrorContent' \
        --output text 2>/dev/null || echo "(no stderr)")
      fail "SSM command $STATUS. Stderr: $STDERR"
      ;;
    *)
      printf "."
      sleep 5
      ;;
  esac

  if [ "$i" -eq 60 ]; then
    fail "Timed out waiting for SSM command to complete."
  fi
done

echo ""

# --- 3. Download dump locally ------------------------------------------------

info "Downloading dump from S3"
aws s3 cp "s3://${BACKUP_BUCKET}/${DUMP_S3_KEY}" "$LOCAL_DUMP" \
  --profile "$AWS_PROFILE" \
  --region "$REGION"

info "Downloaded to $LOCAL_DUMP"

# --- 4. Wipe local database --------------------------------------------------

info "Wiping local database"
psql "$LOCAL_DB_URL" <<'SQL'
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
GRANT ALL ON SCHEMA public TO PUBLIC;
SQL

# --- 5. Restore --------------------------------------------------------------

info "Restoring production data to local database"
gunzip -c "$LOCAL_DUMP" | psql --quiet "$LOCAL_DB_URL"

# cleanup is handled by the EXIT trap
info "Done. Production data is now in your local database."
