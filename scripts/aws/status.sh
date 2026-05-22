#!/bin/bash
#
# Show current state of the footy-forecast AWS infrastructure.
set -euo pipefail

: "${AWS_PROFILE:?AWS_PROFILE must be set}"
REGION="${AWS_REGION:-eu-central-1}"
INSTANCE_TAG="${INSTANCE_TAG:-footy-forecast}"
DOMAIN="${DOMAIN:-footy-forecast.hexagonon.net}"

BOLD='\033[1m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

printf "${BOLD}footy-forecast status${NC}\n"
printf "Profile: %s | Region: %s\n\n" "$AWS_PROFILE" "$REGION"

# Instance
INSTANCE=$(aws ec2 describe-instances \
  --filters "Name=tag:Name,Values=${INSTANCE_TAG}" \
  --query 'Reservations[0].Instances[0].{Id:InstanceId,State:State.Name,Type:InstanceType,PublicIp:PublicIpAddress}' \
  --output json \
  --profile "$AWS_PROFILE" \
  --region "$REGION")

INSTANCE_ID=$(echo "$INSTANCE" | jq -r '.Id // empty')
INSTANCE_STATE=$(echo "$INSTANCE" | jq -r '.State // empty')
INSTANCE_IP=$(echo "$INSTANCE" | jq -r '.PublicIp // "—"')

if [ -z "$INSTANCE_ID" ]; then
  printf "  Instance:  ${RED}not found${NC}\n"
else
  case "$INSTANCE_STATE" in
    running) STATE_COLOR=$GREEN ;;
    stopped) STATE_COLOR=$YELLOW ;;
    *)       STATE_COLOR=$RED ;;
  esac
  printf "  Instance:  %s (${STATE_COLOR}%s${NC}) %s\n" "$INSTANCE_ID" "$INSTANCE_STATE" "$INSTANCE_IP"
fi

# EIP
EIP=$(aws ec2 describe-addresses \
  --filters "Name=instance-id,Values=${INSTANCE_ID}" \
  --query 'Addresses[0].PublicIp' \
  --output text \
  --profile "$AWS_PROFILE" \
  --region "$REGION" 2>/dev/null || echo "None")

if [ "$EIP" = "None" ] || [ -z "$EIP" ]; then
  printf "  EIP:       ${YELLOW}none attached${NC}\n"
else
  printf "  EIP:       %s\n" "$EIP"
fi

# DNS
DNS_IP=$(dig +short "$DOMAIN" 2>/dev/null | tail -1)
if [ -z "$DNS_IP" ]; then
  printf "  DNS:       ${YELLOW}not resolving${NC}\n"
elif [ "$DNS_IP" = "$EIP" ]; then
  printf "  DNS:       ${GREEN}%s → matches EIP${NC}\n" "$DNS_IP"
else
  printf "  DNS:       ${RED}%s → does NOT match EIP${NC}\n" "$DNS_IP"
fi

# Health check (only if running and DNS matches)
if [ "$INSTANCE_STATE" = "running" ] && [ "$DNS_IP" = "$EIP" ]; then
  if curl -fsS --max-time 5 "https://${DOMAIN}/health/ready" > /dev/null 2>&1; then
    printf "  Health:    ${GREEN}/health/ready 200 OK${NC}\n"
  else
    printf "  Health:    ${RED}/health/ready not responding${NC}\n"
  fi
fi

# Recent backups
LATEST_BACKUP=$(aws s3 ls "s3://footy-forecast-backups/daily/" --recursive --profile "$AWS_PROFILE" --region "$REGION" 2>/dev/null \
  | sort -r | head -1)
if [ -n "$LATEST_BACKUP" ]; then
  printf "  Backup:    %s\n" "$(echo "$LATEST_BACKUP" | awk '{print $1, $2, "—", $4}')"
else
  printf "  Backup:    ${YELLOW}none found${NC}\n"
fi
