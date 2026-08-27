# VulcanShield — Adaptive AI Risk Management Platform

VulcanShield is a financial transaction risk management platform designed for real-time risk assessment, velocity signal detection, deterministic policy authorization, fraud-network analysis, and explainable AI investigation.

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
Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

### Docker Compose
Validate container configurations:
```bash
docker compose config
```

---

## Current Status

* **Phase 0 — Repository Initialization**: Completed. Base monorepo structure, environment configuration, and docker-compose specification established.
