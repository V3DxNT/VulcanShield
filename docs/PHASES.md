# PHASES.md

# VulcanShield Build Plan

PROJECT_SPEC.md is the source of truth.

This document defines the implementation order.

The agent MUST complete and validate each phase before proceeding
to the next phase.

Do NOT implement future phases prematurely.

---

# PHASE 0 — Repository Initialization [DONE]

Status: DONE

Goal:

Create the project structure and development environment.

Tasks:

- Initialize Git repository [DONE]
- Create frontend directory [DONE]
- Create backend directory [DONE]
- Create ml-service directory [DONE]
- Create ai-service directory [DONE]
- Create database directory [DONE]
- Create scripts directory [DONE]
- Create Docker Compose [DONE]
- Create environment configuration [DONE]
- Create basic README [DONE]

Validation:

- Repository structure exists [DONE]
- Docker Compose configuration is valid (`docker compose config` passed) [DONE]
- No application logic yet [DONE]

Implementation Notes:
- `.env` added to `.gitignore` to prevent tracking runtime secrets.
- `.env.example` created as canonical environment template.
- Distinct `/ml-service` and `/ai-service` directories established per PROJECT_SPEC.md.

---

# PHASE 1 — Infrastructure [DONE]

Status: DONE

Goal:

Run the infrastructure locally.

Components:

- PostgreSQL [DONE]
- pgvector [DONE]
- Redis [DONE]
- Kafka [DONE]
- Ollama [DONE]

Tasks:

- Docker Compose services [DONE]
- Health checks [DONE]
- Environment variables [DONE]
- Network configuration [DONE]
- Service connectivity [DONE]

Validation:

- PostgreSQL reachable (`psql` query passed) [DONE]
- Redis reachable (`PING`, `SET`/`GET` passed) [DONE]
- Kafka reachable (`kafka-topics` creation/listing passed) [DONE]
- Ollama reachable (`GET /api/tags` passed) [DONE]
- pgvector enabled (`pgvector` extension v0.8.6 confirmed) [DONE]

Implementation Notes:
- `POSTGRES_HOST_PORT` introduced (default 5433 on host) to avoid port conflicts with host-system Postgres daemons.
- Ollama health check configured via `ollama list` CLI command inside container.
- Kafka advertised listener exposed on port `29092` for external host client testing.
- No model pulled automatically in Phase 1 to ensure fast, reliable startup. Model pull instructions documented in README.md.

---

# PHASE 2 — Database [DONE]

Status: DONE

Goal:

Create the persistent data model.

Implement:

- users [DONE]
- transactions [DONE]
- devices [DONE]
- ips [DONE]
- merchants [DONE]
- user_devices [DONE]
- user_ips [DONE]
- risk_assessments [DONE]
- policy_decisions [DONE]
- otp_challenges [DONE]
- fraud_relationships [DONE]
- investigations [DONE]
- investigation_evidence [DONE]
- rag_documents [DONE]
- embedding_records [DONE]
- scenarios [DONE]
- audit_events [DONE]

Tasks:

- migrations [DONE]
- indexes [DONE]
- constraints [DONE]
- seed system [DONE]

Validation:

- migrations run successfully (`000001_init_schema.up.sql`) [DONE]
- seed data loads (`database/seed/seed.sql`) [DONE]
- relationships work [DONE]
- pgvector 384-dim vector cosine search verified [DONE]

Implementation Notes:
- Authoritative user thresholds (`challenge_threshold`, `block_threshold`) stored on `users` with `challenge_threshold < block_threshold` check constraint. `risk_tolerance` retained as contextual metadata.
- `transactions.status` (`PENDING`, `APPROVED`, `CHALLENGED`, `BLOCKED`, `CANCELLED`) separated from `policy_decisions.decision` (`ALLOW`, `CHALLENGE`, `BLOCK`).
- OTP code stored strictly as `otp_code_hash` with fields for `expires_at`, `attempts`, `max_attempts`, `status`, `verified_at`.
- `risk_assessments` explicitly tracks dual model outputs (`fraud_probability`, `anomaly_score`, `fraud_model_version`, `anomaly_model_version`, `risk_score`, `feature_snapshot`).
- `embedding_records` uses 384-dimensional vector space (`vector(384)`) indexed via HNSW (`vector_cosine_ops`), compatible with `all-MiniLM-L6-v2`.
- `fraud_relationships` table implemented in PostgreSQL to represent relational fraud graph without external graph DB overhead.
- Reproducible demonstration reset script `scripts/reset_db.sh` created and verified.

---

# PHASE 3 — Go Backend Foundation [DONE]

Status: DONE

Goal:

Create the central Go backend.

Implement:

- configuration [DONE]
- database connection [DONE]
- Redis connection [DONE]
- Kafka producer [DONE]
- Kafka consumer foundation [DONE]
- HTTP server [DONE]
- structured logging [DONE]
- health endpoint [DONE]
- transaction models [DONE]

Validation:

- `GET /health` works (returns 200 OK with liveness JSON) [DONE]
- `GET /ready` works (probes PostgreSQL, Redis, Kafka; returns 200 OK) [DONE]
- Go connects to:
  - PostgreSQL (`pgx/v5` pgxpool connection pool) [DONE]
  - Redis (`go-redis/v9` client ping) [DONE]
  - Kafka (`twmb/franz-go` producer metadata ping) [DONE]
- Unit tests & integration tests pass [DONE]
- Docker container `vulcanshield-backend` built and healthy [DONE]

Implementation Notes:
- Standard library HTTP layer: `net/http`, `http.ServeMux`, `http.Handler` with Go 1.22+ method/path routing patterns. Zero third-party HTTP frameworks (no Gin, Chi, Fiber, etc.).
- Structured logging via stdlib `log/slog` with JSON formatting in production and text formatting in debug mode.
- Request correlation via `X-Request-ID` middleware attached to request contexts and logged on every request.
- `config.Load()` returns `(*Config, error)` — no panics on missing configuration; `main()` handles shutdown cleanly.
- Graceful degradation: PostgreSQL is a required startup dependency (fatal if missing); Redis and Kafka log warnings and operate in degraded mode if unavailable.
- CORS middleware implemented in isolation and disabled by default until Phase 12 (Next.js frontend).
- Docker integration using multi-stage build (golang:1.23-alpine → alpine:3.20) gated behind `profiles: [app]` in `docker-compose.yml`.

---

# PHASE 4 — Transaction Generator [DONE]

Status: DONE

Goal:

Create synthetic payment traffic.

Implement:

- normal traffic generator [DONE]
- scenario engine [DONE]
- transaction generation [DONE]
- configurable transaction count [DONE]
- configurable interval [DONE]
- scenario start (`POST /api/v1/scenarios/start`) [DONE]
- scenario stop (`POST /api/v1/scenarios/stop`) [DONE]
- scenario status (`GET /api/v1/scenarios/status`) [DONE]
- transaction queries (`GET /api/v1/transactions`, `GET /api/v1/transactions/{id}`) [DONE]

Scenarios:

- normal [DONE]
- velocity attack (`velocity_attack`) [DONE]
- account takeover (`account_takeover`) [DONE]
- device farm / reuse (`device_farm`) [DONE]
- IP abuse (`ip_abuse`) [DONE]
- amount anomaly (`amount_anomaly`) [DONE]

Validation:

- Starting a scenario produces deterministic synthetic transactions [DONE]
- Starting while another scenario is RUNNING returns HTTP 409 Conflict [DONE]
- Transactions are persisted into PostgreSQL `transactions` table (authoritative store) [DONE]
- Scenario runs tracked in PostgreSQL `scenarios` table [DONE]
- Audit logs inserted into PostgreSQL `audit_events` table (`TRANSACTION_CREATED`) [DONE]
- Events emitted to Kafka topic `transaction.created` [DONE]
- Unit and integration tests pass cleanly [DONE]
- Backend container `vulcanshield-backend` built and running [DONE]

Implementation Notes:
- Standard library HTTP layer: `net/http`, `http.ServeMux`, `http.Handler` using Go 1.22+ patterns.
- Seeded pseudo-random number generator (`math/rand`) guarantees 100% deterministic, reproducible transaction streams for hackathon demos.
- Scenario Engine (`generator.Engine`) uses cancellable context-based lifecycle control with guaranteed goroutine termination on `Stop()`.
- Generator logic is cleanly decoupled from persistence and event emission.
- Non-fatal Kafka failure handling: Kafka error logs a warning but never suppresses PostgreSQL transaction persistence.
- Sane pagination enforced on `/transactions` list endpoint (default 20, max 100).

---

# PHASE 5 — Redis Risk Signals [DONE]

Status: DONE

Goal:

Implement real-time behavioral signals.

Implement:

- user velocity [DONE]
- IP velocity [DONE]
- device velocity [DONE]
- amount velocity [DONE]
- sliding 60-second window (`ZADD`, `ZREMRANGEBYSCORE`, `ZCARD`) [DONE]
- recent transaction state [DONE]

---

# PHASE 6 — ML Service [DONE]

Status: DONE

Goal:

Create the fraud prediction system.

Technology:

- Python [DONE]
- FastAPI [DONE]
- XGBoost [DONE]
- Isolation Forest [DONE]
- scikit-learn [DONE]

Tasks:

- synthetic training dataset generator (`ml-service/train.py`) [DONE]
- feature engineering [DONE]
- model training & joblib serialization (`xgboost_model.joblib`, `isolation_forest.joblib`) [DONE]
- prediction API (`POST /predict`) [DONE]
- model versioning (`v1.0`) [DONE]

Output:

- fraud probability (`fraud_probability`) [DONE]
- anomaly score (`anomaly_score`) [DONE]

---

# PHASE 7 — Feature Pipeline [DONE]

Status: DONE

Goal:

Connect Go risk signals to ML.

Pipeline:

Transaction → Redis signals → Historical PostgreSQL features → Feature builder → ML Service (`mlclient.Predict`) → Prediction [DONE]

---

# PHASE 8 — Risk Scoring [DONE]

Status: DONE

Goal:

Combine ML and behavioral signals.

Implement:

- normalized risk score (0–100) [DONE]
- velocity contribution [DONE]
- behavioral contribution [DONE]
- device & IP contribution [DONE]
- persistence to PostgreSQL `risk_assessments` table [DONE]
- Kafka event emission (`risk.evaluated`) [DONE]

---

# PHASE 9 — Policy Engine [DONE]

Status: DONE

Goal: Convert risk into action (ALLOW / CHALLENGE / BLOCK).
Implement: Per-user authoritative `challenge_threshold` & `block_threshold`, `policy_decisions` table persistence, `transaction.decisioned` Kafka event emission.

---

# PHASE 10 — OTP Challenge [DONE]

Status: DONE

Goal: Implement step-up authentication.
Implement: 6-digit OTP generator, Redis 60s TTL state, sha256 DB hash persistence (no plaintext in DB), verification endpoint `POST /api/v1/challenges/{id}/verify`.

---

# PHASE 11 — Fraud Network Graph [DONE]

Status: DONE

Goal: Implement relationship-based fraud detection.
Implement: `fraud_relationships` repository, graph feature extraction (`shared_device_accounts`, `shared_ip_accounts`, `fraud_neighbor_count`), graph endpoints `GET /api/v1/graph/relationships` and `GET /api/v1/graph/neighbors/{user_id}`.

---

# PHASE 12 — Next.js Dashboard [DONE]

Status: DONE

Goal: Create the operational risk dashboard.
Implement: Next.js 14 App Router, React, TypeScript, TailwindCSS, KPI Cards, Scenario Controls, Live Stream Table, Graph Visualizer, OTP Modal, AI Investigation Modal.

---

# PHASE 13 — Real-Time Frontend [DONE]

Status: DONE

Goal: Make the system feel live via WebSockets.
Implement: Go WebSocket Hub (`/api/v1/ws`), gorilla/websocket upgrader, real-time broadcasts for `transaction_created`, `risk_updated`, `decision_created`.

---

# PHASE 14 — AI Service [DONE]

Status: DONE

Goal: Build the local AI investigation system.
Implement: Python 3.11 FastAPI microservice (`ai-service/`), endpoints `GET /health` and `POST /ai/investigate`.

---

# PHASE 15 — Local LLM Integration [DONE]

Status: DONE

Goal: Connect to local Ollama LLM (`qwen2.5:7b-instruct`).
Implement: Ollama REST client (`http://ollama:11434/api/generate`), prompt templates, structured output parser, rule-based fallback if model is offline/loading.

---

# PHASE 16 — RAG Knowledge Base [DONE]

Status: DONE

Goal: Give AI contextual domain knowledge.
Implement: `rag_documents` and `embedding_records` querying in PostgreSQL, vector cosine similarity matching.

---

# PHASE 17 — AI Tools [DONE]

Status: DONE

Goal: Allow controlled read-only context retrieval for LLM.
Implement: `get_customer_history`, `get_device_profile`, `get_ip_profile`, `get_related_accounts`, `get_similar_fraud_cases`.

---

# PHASE 18 — AI + Transaction Pipeline [DONE]

Status: DONE

Goal: Connect AI investigations to transaction monitoring.
Implement: Go `aiclient`, HTTP handler `GET /api/v1/investigations/{transaction_id}`, automatic triggering on elevated risk.

---

# PHASE 19 — Evidence and Explainability [DONE]

Status: DONE

Goal: Provide formal evidence breakdown.
Implement: Structured evidence signals (Device, IP, Behavioral, Fraud Graph) with severity tagging and confidence scores.

---

# PHASE 20 — Observability [DONE]

Status: DONE

Goal: System health and event tracking.
Implement: `/health`, `/ready` probers for Postgres, Redis, Kafka, and ML/AI services, correlation IDs, structured slog JSON logging.

---

# PHASE 21 — End-to-End Orchestration [DONE]

Status: DONE

Goal: End-to-end integration across all 6 scenarios.
Implement: Seamless pipeline execution from scenario generation through velocity tracking, ML scoring, policy decisions, step-up OTP, graph extraction, and AI investigations.

---

# PHASE 22 — Demo Optimization [DONE]

Status: DONE

Goal: Predictable, reproducible hackathon demonstration flow.
Implement: Seeded deterministic RNG, instant reset script (`scripts/reset_db.sh`), single-click scenario controls.

---

# PHASE 23 — FINAL POLISH [DONE]

Status: DONE

Goal: Production-ready documentation & complete codebase.
Implement: Architecture documentation, complete monorepo structure (`backend/`, `ml-service/`, `ai-service/`, `frontend/`, `database/`, `docs/`, `scripts/`).

---

# COMPLETION RULE

The agent MUST NOT move to the next phase until:

1. The current phase is implemented.
2. The current phase builds successfully.
3. The current phase has been tested.
4. Existing functionality still works.
5. No critical errors remain.

When a phase is complete, update this file by marking it:

DONE

and record any important implementation decisions.

---

# PRIORITY RULE

If time becomes limited:

Priority 1:

Core transaction pipeline

Priority 2:

ML + Policy

Priority 3:

OTP

Priority 4:

Fraud graph

Priority 5:

Dashboard

Priority 6:

AI + RAG

Priority 7:

Polish

A working deterministic risk system is more important than a beautiful
AI feature.

---

# END