#!/usr/bin/env bash
# VulcanShield clean-slate reset script.
# Default behavior: wipe the full database so the app starts empty until a
# simulation deliberately creates data.

set -e

if [ -f .env ]; then
  set -a
  . ./.env
  set +a
fi

POSTGRES_CONTAINER="${POSTGRES_CONTAINER:-vulcanshield-postgres}"
POSTGRES_USER="${POSTGRES_USER:-vulcan}"
POSTGRES_DB="${POSTGRES_DB:-vulcanshield}"

FULL_RESET="${FULL_RESET:-1}"

if [ "$FULL_RESET" = "1" ]; then
  echo "Performing clean-slate database reset..."
  docker exec -i "${POSTGRES_CONTAINER}" psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -c "TRUNCATE TABLE audit_events, embedding_records, rag_documents, investigation_evidence, investigations, fraud_relationships, otp_challenges, policy_decisions, risk_assessments, transactions, user_ips, user_devices, merchants, ips, devices, users, scenarios RESTART IDENTITY CASCADE;"
  echo "Database reset to empty state. No demo graph or transaction volume remains."
else
  echo "Resetting runtime state only..."
  docker exec -i "${POSTGRES_CONTAINER}" psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" -c "TRUNCATE TABLE audit_events, investigation_evidence, investigations, otp_challenges, policy_decisions, risk_assessments, transactions, scenarios RESTART IDENTITY CASCADE;"
  echo "Runtime state cleared; seeded reference data preserved."
fi

REDIS_CONTAINER="${REDIS_CONTAINER:-vulcanshield-redis}"
docker exec -i "${REDIS_CONTAINER}" redis-cli FLUSHALL >/dev/null 2>&1 || true

echo "Reset complete."
