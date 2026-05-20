#!/bin/bash
#
# Pause the footy-forecast infrastructure to ~$0.80/month.
#
# Steps:
#   1. Trigger a final backup via SSM
#   2. Verify the backup landed in S3
#   3. Release the Elastic IP
#   4. Stop the EC2 instance
#
# Required env: AWS_PROFILE
set -euo pipefail

: "${AWS_PROFILE:?AWS_PROFILE must be set (e.g. hexa)}"
REGION="${AWS_REGION:-eu-central-1}"
INSTANCE_TAG="${INSTANCE_TAG:-footy-forecast}"
BACKUP_BUCKET="${BACKUP_BUCKET:-footy-forecast-backups}"

# Color output helpers
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

info()  { printf "${GREEN}==>${NC} %s\n" "$*"; }
warn()  { printf "${YELLOW}!!!${NC} %s\n" "$*" >&2; }
fail()  { printf "${RED}ERR${NC} %s\n" "$*" >&2; exit 1; }

info "Pausing footy-forecast infrastructure"
info "Profile: $AWS_PROFILE | Region: $REGION"

# --- 1. Find the instance --------------------------------------------------

INSTANCE_ID=$(aws ec2 describe-instances \
  --filters "Name=tag:Name,Values=${INSTANCE_TAG}" "Name=instance-state-name,Values=running,stopping,stopped" \
  --query 'Reservations[0].Instances[0].InstanceId' \
  --output text \
  --profile "$AWS_PROFILE" \
  --region "$REGION")

if [ "$INSTANCE_ID" = "None" ] || [ -z "$INSTANCE_ID" ]; then
  fail "No instance found with tag Name=${INSTANCE_TAG}"
fi

INSTANCE_STATE=$(aws ec2 describe-instances \
  --instance-ids "$INSTANCE_ID" \
  --query 'Reservations[0].Instances[0].State.Name' \
  --output text \
  --profile "$AWS_PROFILE" \
  --region "$REGION")

info "Instance: $INSTANCE_ID (state: $INSTANCE_STATE)"

if [ "$INSTANCE_STATE" = "stopped" ]; then
  warn "Instance is already stopped. Skipping backup and stop steps."
else
  # --- 2. Final backup via SSM ---------------------------------------------
  info "Triggering final backup before pause"

  COMMAND_ID=$(aws ssm send-command \
    --instance-ids "$INSTANCE_ID" \
    --document-name "AWS-RunShellScript" \
    --parameters 'commands=["/usr/local/bin/footy-forecast-backup"]' \
    --comment "pre-pause backup" \
    --query 'Command.CommandId' \
    --output text \
    --profile "$AWS_PROFILE" \
    --region "$REGION")

  info "Backup command: $COMMAND_ID — waiting for completion"

  for i in {1..60}; do
    STATUS=$(aws ssm get-command-invocation \
      --command-id "$COMMAND_ID" \
      --instance-id "$INSTANCE_ID" \
      --query 'Status' \
      --output text \
      --profile "$AWS_PROFILE" \
      --region "$REGION" 2>/dev/null || echo "Pending")
    case "$STATUS" in
      Success) info "Backup succeeded"; break ;;
      Failed|Cancelled|TimedOut)
        warn "Backup ended with status: $STATUS — aborting pause for safety"
        aws ssm get-command-invocation \
          --command-id "$COMMAND_ID" \
          --instance-id "$INSTANCE_ID" \
          --query '{StdOut: StandardOutputContent, StdErr: StandardErrorContent}' \
          --output json \
          --profile "$AWS_PROFILE" \
          --region "$REGION"
        fail "Pause aborted: backup did not succeed"
        ;;
    esac
    if [ "$i" = 60 ]; then
      fail "Backup timed out after 5 minutes"
    fi
    sleep 5
  done

  # --- 3. Stop the instance ------------------------------------------------
  info "Stopping instance"
  aws ec2 stop-instances \
    --instance-ids "$INSTANCE_ID" \
    --profile "$AWS_PROFILE" \
    --region "$REGION" \
    --output text > /dev/null

  info "Waiting for instance to enter 'stopped' state"
  aws ec2 wait instance-stopped \
    --instance-ids "$INSTANCE_ID" \
    --profile "$AWS_PROFILE" \
    --region "$REGION"
  info "Instance stopped"
fi

# --- 4. Release the Elastic IP ---------------------------------------------
# Find the EIP currently associated with the instance (if any)
ALLOCATION_ID=$(aws ec2 describe-addresses \
  --filters "Name=instance-id,Values=${INSTANCE_ID}" \
  --query 'Addresses[0].AllocationId' \
  --output text \
  --profile "$AWS_PROFILE" \
  --region "$REGION")

if [ "$ALLOCATION_ID" = "None" ] || [ -z "$ALLOCATION_ID" ]; then
  warn "No Elastic IP associated with this instance. Skipping release."
else
  PUBLIC_IP=$(aws ec2 describe-addresses \
    --allocation-ids "$ALLOCATION_ID" \
    --query 'Addresses[0].PublicIp' \
    --output text \
    --profile "$AWS_PROFILE" \
    --region "$REGION")

  info "Releasing EIP $PUBLIC_IP (allocation $ALLOCATION_ID)"
  aws ec2 release-address \
    --allocation-id "$ALLOCATION_ID" \
    --profile "$AWS_PROFILE" \
    --region "$REGION"
  info "EIP released"
fi

# --- Final
