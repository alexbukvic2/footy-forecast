#!/bin/bash
#
# Resume the footy-forecast infrastructure from paused state.
#
# Steps:
#   1. Start the EC2 instance
#   2. Allocate and associate a fresh Elastic IP
#   3. Print the new IP and DNS instructions
#   4. Wait for the app to become healthy
#
# Required env: AWS_PROFILE
set -euo pipefail

: "${AWS_PROFILE:?AWS_PROFILE must be set (e.g. hexa)}"
REGION="${AWS_REGION:-eu-central-1}"
INSTANCE_TAG="${INSTANCE_TAG:-footy-forecast}"
DOMAIN="${DOMAIN:-footy-forecast.hexagonon.net}"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
NC='\033[0m'

info()  { printf "${GREEN}==>${NC} %s\n" "$*"; }
warn()  { printf "${YELLOW}!!!${NC} %s\n" "$*" >&2; }
fail()  { printf "${RED}ERR${NC} %s\n" "$*" >&2; exit 1; }

info "Resuming footy-forecast infrastructure"
info "Profile: $AWS_PROFILE | Region: $REGION"

# --- 1. Find the instance --------------------------------------------------

INSTANCE_ID=$(aws ec2 describe-instances \
  --filters "Name=tag:Name,Values=${INSTANCE_TAG}" \
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

# --- 2. Start the instance -------------------------------------------------

if [ "$INSTANCE_STATE" = "running" ]; then
  warn "Instance is already running. Skipping start."
else
  info "Starting instance"
  aws ec2 start-instances \
    --instance-ids "$INSTANCE_ID" \
    --profile "$AWS_PROFILE" \
    --region "$REGION" \
    --output text > /dev/null

  info "Waiting for instance to enter 'running' state"
  aws ec2 wait instance-running \
    --instance-ids "$INSTANCE_ID" \
    --profile "$AWS_PROFILE" \
    --region "$REGION"
  info "Instance is running"
fi

# --- 3. Check for existing EIP, or allocate a new one ---------------------

ALLOCATION_ID=$(aws ec2 describe-addresses \
  --filters "Name=instance-id,Values=${INSTANCE_ID}" \
  --query 'Addresses[0].AllocationId' \
  --output text \
  --profile "$AWS_PROFILE" \
  --region "$REGION")

if [ "$ALLOCATION_ID" != "None" ] && [ -n "$ALLOCATION_ID" ]; then
  PUBLIC_IP=$(aws ec2 describe-addresses \
    --allocation-ids "$ALLOCATION_ID" \
    --query 'Addresses[0].PublicIp' \
    --output text \
    --profile "$AWS_PROFILE" \
    --region "$REGION")
  warn "EIP already attached: $PUBLIC_IP (allocation $ALLOCATION_ID). Reusing."
else
  info "Allocating new Elastic IP"
  ALLOCATION_ID=$(aws ec2 allocate-address \
    --domain vpc \
    --tag-specifications "ResourceType=elastic-ip,Tags=[{Key=Project,Value=footy-forecast}]" \
    --query 'AllocationId' \
    --output text \
    --profile "$AWS_PROFILE" \
    --region "$REGION")

  PUBLIC_IP=$(aws ec2 describe-addresses \
    --allocation-ids "$ALLOCATION_ID" \
    --query 'Addresses[0].PublicIp' \
    --output text \
    --profile "$AWS_PROFILE" \
    --region "$REGION")

  info "Allocated EIP $PUBLIC_IP (allocation $ALLOCATION_ID)"

  info "Associating EIP with instance"
  aws ec2 associate-address \
    --instance-id "$INSTANCE_ID" \
    --allocation-id "$ALLOCATION_ID" \
    --allow-reassociation \
    --profile "$AWS_PROFILE" \
    --region "$REGION" \
    --output text > /dev/null

  info "EIP associated"
fi

# --- 4. Manual DNS reminder ------------------------------------------------

cat <<EOF

${BOLD}${YELLOW}=== ACTION REQUIRED: Update DNS at Cloudflare ===${NC}

  Domain: ${DOMAIN}
  Type:   A
  Value:  ${PUBLIC_IP}
  Proxy:  DNS only (grey cloud)
  TTL:    Auto

  Until DNS is updated, the app will not be reachable at ${DOMAIN}.
  TLS will resume once DNS resolves to the new IP and Caddy renews
  its cert (the existing cert may still work for ~80 days regardless).

EOF

read -p "Press Enter once DNS is updated to continue health-checking, or Ctrl-C to skip: "

# --- 5. Wait for SSH ------------------------------------------------------

info "Waiting for SSH (port 22) to accept connections at ${PUBLIC_IP}"
for i in {1..30}; do
  if nc -z -w 3 "$PUBLIC_IP" 22 2>/dev/null; then
    info "SSH reachable after ${i} attempts"
    break
  fi
  printf "."
  sleep 2
  if [ "$i" = 30 ]; then
    warn "SSH not reachable after 60s. The instance may still be booting; check manually."
  fi
done

# --- 6. Wait for HTTPS readiness ------------------------------------------

info "Waiting for ${DOMAIN}/health/ready to respond 200"
for i in {1..60}; do
  if curl -fsS --max-time 5 "https://${DOMAIN}/health/ready" > /dev/null 2>&1; then
    info "App is healthy"
    break
  fi
  printf "."
  sleep 5
  if [ "$i" = 60 ]; then
    warn "App did not become healthy in 5 minutes."
    warn "Common causes: DNS not propagated yet, security group SSH rule out of date."
    warn "SSH to the box and check 'sudo systemctl status footy-forecast'."
  fi
done

# --- Final summary --------------------------------------------------------

cat <<EOF

${GREEN}==> Resume complete${NC}

  Instance:  ${INSTANCE_ID} (running)
  Public IP: ${PUBLIC_IP}
  URL:       https://${DOMAIN}

  Don't forget: the security group's SSH rule may still point to your old
  laptop IP. If you can't SSH, update the rule:
    aws ec2 authorize-security-group-ingress \\
      --group-name footy-forecast-sg \\
      --protocol tcp --port 22 \\
      --cidr <YOUR-CURRENT-IP>/32 \\
      --profile $AWS_PROFILE --region $REGION

EOF
