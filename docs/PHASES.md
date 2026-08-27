# PHASES.md

# VulcanShield Build Plan

PROJECT_SPEC.md is the source of truth.

This document defines the implementation order.

The agent MUST complete and validate each phase before proceeding
to the next phase.

Do NOT implement future phases prematurely.

---

# PHASE 0 — Repository Initialization

Goal:

Create the project structure and development environment.

Tasks:

- Initialize Git repository
- Create frontend directory
- Create backend directory
- Create ml-service directory
- Create ai-service directory
- Create database directory
- Create scripts directory
- Create Docker Compose
- Create environment configuration
- Create basic README

Validation:

- Repository structure exists
- Docker Compose configuration is valid
- No application logic yet

---

# PHASE 1 — Infrastructure

Goal:

Run the infrastructure locally.

Components:

- PostgreSQL
- pgvector
- Redis
- Kafka
- Ollama

Tasks:

- Docker Compose services
- Health checks
- Environment variables
- Network configuration
- Service connectivity

Validation:

- PostgreSQL reachable
- Redis reachable
- Kafka reachable
- Ollama reachable
- pgvector enabled

Do not build ML or AI logic yet.

---

# PHASE 2 — Database

Goal:

Create the persistent data model.

Implement:

- users
- transactions
- devices
- ips
- merchants
- user_devices
- user_ips
- risk_assessments
- policy_decisions
- otp_challenges
- fraud_relationships
- investigations
- investigation_evidence
- rag_documents
- embedding_records
- scenarios
- audit_events

Tasks:

- migrations
- indexes
- constraints
- seed system

Validation:

- migrations run successfully
- seed data loads
- relationships work

---

# PHASE 3 — Go Backend Foundation

Goal:

Create the central Go backend.

Implement:

- configuration
- database connection
- Redis connection
- Kafka producer
- Kafka consumer foundation
- HTTP server
- structured logging
- health endpoint
- transaction models

Validation:

GET /health works.

Go can connect to:

- PostgreSQL
- Redis
- Kafka

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