# VulcanShield — Adaptive AI Risk Management Platform

VulcanShield is a financial transaction risk management platform designed for real-time risk assessment, velocity signal detection, deterministic policy authorization, fraud-network analysis, and explainable AI investigation.

## Architecture Overview

![Vulcan Shield logo](./assets/vulcan-shield.svg)

- Code: vedx.dev
- YouTube: [placeholder](https://youtube.com/placeholder)

Services and responsibilities

- backend: Go service that orchestrates the transaction lifecycle, policy engine, Redis/Kafka integration, and Postgres persistence. It is the authoritative decision-maker (policy decides).
- frontend: Next.js + React UI for dashboards, live transaction stream, investigation modal, and scenario controls.
- ml-service: Python FastAPI service hosting the XGBoost and Isolation Forest models (ML predicts). It receives structured feature vectors and returns `fraud_probability` and `anomaly_score`.
- ai-service: Python FastAPI investigator that performs structured retrieval (RAG-like customer-history fetch) and calls the configured LLM (Ollama/Groq) or falls back to deterministic rule-based summaries when the LLM is unavailable (AI investigates and explains).
- postgres: Durable source of truth for users, transactions, risk assessments, and RAG documents (pgvector for semantic retrieval).
- redis: Real-time counters and velocity engine for temporary state, sliding windows, and OTP/challenge state.
- kafka: Event streaming for decoupling transaction generation, risk evaluation, and downstream consumers.

Problems addressed

- Real-time velocity detection: handled by Redis to avoid expensive DB lookups for high-frequency checks.
- Deterministic policy: Policy engine ensures auditable ALLOW / CHALLENGE / BLOCK decisions separate from ML output.
- Explainability: AI investigator combines structured retrieval + LLM to produce analyst-friendly explanations; falls back to deterministic summaries when models are unavailable.

Notes

- The ML models in `ml-service` rely on training data. The project uses the internal scenario generator to produce synthetic training and demo data, but the models must be trained before expecting meaningful predictions. (Do not train now — the system will use fallback rules or pre-saved models if present.)


---

## Core Architectural Invariants

> **ML predicts. Policy decides. AI investigates and explains.**

* **ML Layer**: Predicts risk probability and anomaly scores from behavioral and transactional features.
* **Policy Engine**: Makes deterministic, auditable decisions (`ALLOW`, `CHALLENGE`, `BLOCK`).
* **AI / LLM Layer**: Investigates suspicious transactions, provides semantic context (RAG), and explains evidence without overriding policy decisions.
* **Fraud Network**: Discovers hidden relationships across users, devices, IPs, and merchants.

---

## Repository Structure

```text
.
├── backend/          # Go primary transaction orchestration & HTTP/WebSocket backend
├── frontend/         # Next.js real-time monitoring dashboard UI
├── ml-service/       # Python FastAPI ML inference service (XGBoost, Isolation Forest)
├── ai-service/       # Python FastAPI AI service (Ollama local LLM, RAG, AI Tools)
├── database/         # PostgreSQL schema migrations and seed scripts
├── scripts/          # Helper scripts and automation utilities
├── docs/             # System specifications, decision records, and phase guides
│   ├── PHASES.md     # Authoritative phase execution roadmap
│   └── DECISIONS.md  # Architecture decision records (ADRs)
├── AGENTS.md         # Development guidelines for AI coding agents
├── PROJECT_SPEC.md   # Canonical system specification (Single Source of Truth)
├── docker-compose.yml# Container orchestration specification
├── .env.example      # Environment variable template
└── README.md         # Repository documentation
```

---

## Development Setup

### Environment Variables
Copy `.env.example` to `.env` (note: `.env` is ignored by `.gitignore`):
```bash
cp .env.example .env
```

### Infrastructure Management (Phase 1)
Start the services that do not require Ollama in Docker:
```bash
docker compose up -d postgres redis zookeeper kafka backend ml-service ai-service frontend
```

Check infrastructure health status:
```bash
docker compose ps
```

Stop infrastructure services:
```bash
docker compose down
```

### Local LLM Model Setup (host-installed Ollama)
Use a native Ollama install on your machine rather than the Docker Ollama container. The app will talk to the local service on port `11434`.

Install Ollama on macOS:
```bash
brew install ollama
```

Start the local service:
```bash
ollama serve
```

Pull the model once on the host:
```bash
ollama pull qwen2.5:7b-instruct
```

If the app is running inside Docker, point it to the host service with:
```bash
OLLAMA_URL=http://host.docker.internal:11434
```

When running the AI service directly on your machine, use:
```bash
OLLAMA_URL=http://localhost:11434
```

---

## Current Status

* **Phase 0 — Repository Initialization**: Completed.
* **Phase 1 — Infrastructure**: Completed. Local infrastructure (PostgreSQL with pgvector, Redis, Kafka + Zookeeper, Ollama) is configured, validated, and runnable via Docker Compose.

