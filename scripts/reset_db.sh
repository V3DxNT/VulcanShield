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

# 1. Execute Down Migration (Teardown)
echo "Executing 000001_init_schema.down.sql..."
docker exec -i "${POSTGRES_CONTAINER}" psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" < database/migrations/000001_init_schema.down.sql

# 2. Execute Up Migration (Schema Creation)
echo "Executing 000001_init_schema.up.sql..."
docker exec -i "${POSTGRES_CONTAINER}" psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" < database/migrations/000001_init_schema.up.sql

# 3. Load Synthetic Seed Data
echo "Executing database/seed/seed.sql..."
docker exec -i "${POSTGRES_CONTAINER}" psql -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" < database/seed/seed.sql

echo "VulcanShield database successfully reset and seeded!"
