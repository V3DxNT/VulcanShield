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

# PHASE 4 — Transaction Generator

Goal:

Create synthetic payment traffic.

Implement:

- normal traffic generator
- scenario engine
- transaction generation
- configurable transaction count
- configurable interval
- scenario start
- scenario stop

Scenarios:

- normal
- velocity attack
- account takeover
- device farm
- IP abuse
- amount anomaly

Validation:

Starting a scenario produces transactions.

Transactions appear in PostgreSQL and Kafka.

---

# PHASE 5 — Redis Risk Signals

Goal:

Implement real-time behavioral signals.

Implement:

- user velocity
- IP velocity
- device velocity
- amount velocity
- sliding 60-second window
- optional 5-minute window
- recent transaction state

Validation:

Velocity attack produces visibly different Redis signals.

Redis state expires correctly.

---

# PHASE 6 — ML Service

Goal:

Create the fraud prediction system.

Technology:

- Python
- FastAPI
- XGBoost
- Isolation Forest
- scikit-learn

Tasks:

- synthetic training dataset
- feature engineering
- model training
- model serialization
- prediction API
- model versioning

Output:

- fraud probability
- anomaly score

Validation:

POST /predict returns valid predictions.

---

# PHASE 7 — Feature Pipeline

Goal:

Connect Go risk signals to ML.

Pipeline:

Transaction

↓

Redis signals

↓

Historical PostgreSQL features

↓

Feature builder

↓

ML Service

↓

Prediction

Validation:

Every transaction receives:

- fraud probability
- anomaly score
- feature snapshot

---

# PHASE 8 — Risk Scoring

Goal:

Combine ML and behavioral signals.

Implement:

- normalized risk score
- velocity contribution
- behavioral contribution
- device contribution
- IP contribution
- graph contribution

Output:

Risk score 0–100.

Validation:

Risk score changes appropriately between normal and attack scenarios.

---

# PHASE 9 — Policy Engine

Goal:

Convert risk into action.

Implement:

- user risk profiles
- adaptive thresholds
- ALLOW
- CHALLENGE
- BLOCK
- deterministic rules
- hard security overrides

Important:

ML predicts.

Policy decides.

LLM cannot override Policy.

Validation:

Different users can have different thresholds.

---

# PHASE 10 — OTP Challenge

Goal:

Implement step-up authentication.

Implement:

- challenge creation
- locally generated OTP
- Redis OTP state
- 60-second expiry
- verification
- attempt tracking
- policy re-evaluation

Validation:

CHALLENGE → OTP → re-evaluation → final decision.

---

# PHASE 11 — Fraud Network Graph

Goal:

Implement relationship-based fraud detection.

Entities:

- users
- devices
- IPs
- merchants
- transactions

Implement:

- relationship storage
- graph traversal
- suspicious relationship detection
- graph risk score

Validation:

Shared device/IP scenarios expose connected accounts.

---

# PHASE 12 — Next.js Dashboard

Goal:

Create the operational risk dashboard.

Pages:

- dashboard
- transactions
- transaction details
- scenarios
- investigations

Components:

- KPI cards
- live transaction table
- risk visualization
- decision distribution
- transaction timeline
- graph visualization
- OTP interface
- AI investigation panel

Validation:

Frontend successfully displays backend data.

---

# PHASE 13 — Real-Time Frontend

Goal:

Make the system feel live.

Implement:

- WebSocket connection
- transaction events
- risk events
- decision events
- OTP events
- AI investigation events

Validation:

Starting a scenario visibly changes the dashboard in real time.

---

# PHASE 14 — AI Service

Goal:

Build the local AI investigation system.

Technology:

- Python
- FastAPI
- Ollama
- local instruct model

Recommended:

Qwen2.5 7B Instruct

Fallback:

smaller local instruct model if required.

Implement:

- investigation endpoint
- context builder
- structured output
- model integration

Validation:

AI can investigate a suspicious transaction locally.

---

# PHASE 15 — RAG

Goal:

Give AI contextual knowledge.

Implement:

- global fraud knowledge
- user historical context
- embeddings
- pgvector
- semantic retrieval
- metadata filtering
- recency

Validation:

AI investigation retrieves relevant evidence.

---

# PHASE 16 — AI Tools

Goal:

Allow controlled contextual investigation.

Tools:

- get_user_history
- get_transaction_context
- get_fraud_network_context
- get_recent_velocity
- search_fraud_knowledge
- get_policy_context

Implement:

- tool router
- validation
- timeout
- bounded results
- logging

Validation:

LLM can request tools but cannot execute arbitrary code or SQL.

---

# PHASE 17 — AI + Fraud Graph

Goal:

Give the AI network-level context.

AI should receive:

- connected users
- shared devices
- shared IPs
- suspicious relationships
- previous fraud-linked accounts

Validation:

AI explanation references real graph evidence.

---

# PHASE 18 — AI + Transaction Pipeline

Goal:

Connect AI investigation to suspicious transactions.

Pipeline:

Policy Decision

↓

CHALLENGE/BLOCK

↓

Kafka event

↓

AI Worker

↓

Context Retrieval

↓

Tools

↓

RAG

↓

Ollama

↓

Investigation

↓

PostgreSQL

↓

WebSocket

↓

Frontend

Validation:

A suspicious transaction automatically receives an AI investigation.

AI failure must not affect the policy decision.

---

# PHASE 19 — Evidence and Explainability

Goal:

Make the system judge-friendly.

Display:

- risk score
- ML probability
- anomaly score
- policy thresholds
- policy reasons
- velocity
- user history
- graph evidence
- AI findings
- AI confidence

Every AI claim must reference evidence.

Validation:

A judge can inspect WHY a transaction was blocked/challenged.

---

# PHASE 20 — Observability

Goal:

Demonstrate engineering quality.

Display:

- transaction throughput
- ML latency
- policy latency
- AI latency
- Redis latency
- Kafka activity
- service health

Implement:

- correlation IDs
- structured logs
- audit events

---

# PHASE 21 — End-to-End Testing

Test:

1. Normal transaction
2. Velocity attack
3. Account takeover
4. Device farm
5. IP abuse
6. Amount anomaly
7. OTP success
8. OTP failure
9. OTP expiry
10. AI failure
11. ML failure
12. Redis failure
13. Kafka failure

---

# PHASE 22 — Demo Optimization

Create one polished demonstration flow.

Sequence:

Normal transaction

↓

ALLOW

↓

Suspicious transaction

↓

ML detects elevated risk

↓

Policy:

CHALLENGE

↓

OTP

↓

Risk decreases

↓

ALLOW

↓

Attack continues

↓

Redis detects velocity

↓

Graph exposes suspicious relationships

↓

ML risk becomes HIGH

↓

Policy:

BLOCK

↓

AI investigation starts

↓

RAG + user history + graph

↓

Local LLM explanation

↓

Evidence shown on dashboard

---

# PHASE 23 — FINAL POLISH

Only after all functionality works:

- improve UI
- improve charts
- improve animations
- improve graph visualization
- improve AI explanation
- improve scenario realism
- improve error messages
- improve README
- create architecture diagram
- create demo script

Do not add major infrastructure.

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