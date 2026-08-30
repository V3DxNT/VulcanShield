#!/usr/bin/env bash
# VulcanShield Clean Database Reset Script for Hackathon Demonstrations

set -e

# Load environment variables if .env exists
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

POSTGRES_CONTAINER="${POSTGRES_CONTAINER:-vulcanshield-postgres}"
POSTGRES_USER="${POSTGRES_USER:-vulcan}"
POSTGRES_DB="${POSTGRES_DB:-vulcanshield}"

echo "Resetting VulcanShield database in container: ${POSTGRES_CONTAINER}..."

# 1. Clear runtime/demo data only, while preserving seeded customer/device/IP data.
echo "Clearing live transaction, risk, challenge, and investigation state..."
docker exec -i "${POSTGRES_CONTAINER}" psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -c "TRUNCATE TABLE audit_events, investigation_evidence, investigations, otp_challenges, policy_decisions, risk_assessments, transactions, scenarios RESTART IDENTITY CASCADE;"

# 2. Clear Redis runtime state as well
echo "Flushing Redis demo cache..."
REDIS_CONTAINER="${REDIS_CONTAINER:-vulcanshield-redis}"
docker exec -i "${REDIS_CONTAINER}" redis-cli FLUSHALL >/dev/null 2>&1 || true

echo "VulcanShield demo state successfully reset. Seeded customer graph and risk baselines remain available."
