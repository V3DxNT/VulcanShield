Phase 2 has been completed, validated, committed, and pushed.

We are now beginning PHASE 3 — GO CORE BACKEND.

Do NOT implement Phase 3 yet.

First read:

1. AGENTS.md
2. PROJECT_SPEC.md
3. docs/PHASES.md
4. docs/DECISIONS.md
5. Existing database migrations
6. Existing docker-compose.yml
7. Existing repository structure

Then produce ONLY an implementation plan for Phase 3.

The objective of Phase 3 is to establish the Go backend as the
central application orchestrator.

Phase 3 must initially cover:

- Go project initialization
- configuration management
- HTTP server
- routing
- middleware
- structured logging
- error handling
- PostgreSQL connection/repository layer
- Redis connection/client layer
- Kafka producer/client foundation
- health/readiness endpoints
- graceful shutdown
- clean package structure
- basic application configuration

Do NOT implement:

- ML
- FastAPI
- XGBoost
- Isolation Forest
- LLM
- Ollama integration
- RAG
- embeddings
- policy engine
- OTP
- fraud graph logic
- frontend
- scenario generation

Those belong to later phases.

The Go service must be designed so these future components can
be integrated cleanly without rewriting the core.

The plan must specify:

1. Exact folder structure.
2. Exact Go packages.
3. Responsibilities of each package.
4. HTTP API conventions.
5. Configuration variables.
6. PostgreSQL connection strategy.
7. Redis connection strategy.
8. Kafka connection strategy.
9. Logging strategy.
10. Error handling strategy.
11. Health/readiness endpoint behavior.
12. Graceful shutdown behavior.
13. Docker integration.
14. Testing strategy.
15. Exact files to create/modify.

Do not modify files.

Do not implement anything.

Wait for my approval after presenting the plan.