Phase 1 has been completed and validated.

The infrastructure is working:

- PostgreSQL + pgvector
- Redis
- Zookeeper
- Kafka
- Ollama

Phase 1 has been committed and pushed.

We are now preparing for PHASE 2 — DATABASE.

Before modifying anything:

1. Read AGENTS.md.
2. Read PROJECT_SPEC.md.
3. Read docs/PHASES.md.
4. Read docs/DECISIONS.md.
5. Inspect the current repository and Git state.
6. Inspect the existing PostgreSQL infrastructure configuration.

Do NOT modify files yet.

Analyze ONLY Phase 2.

Provide:

1. The exact database tables required by PROJECT_SPEC.md.
2. Every table's purpose.
3. Primary keys and foreign keys.
4. Important indexes.
5. Constraints.
6. Relationships between entities.
7. Which fields are needed for risk scoring.
8. Which fields are needed for velocity/risk signals.
9. Which fields are needed for the fraud network graph.
10. Which fields are needed for AI investigation.
11. Which data will eventually be embedded for RAG.
12. How pgvector should be used.
13. The migration strategy.
14. How the database will be reset for clean demonstrations.
15. The exact files you intend to create/modify.
16. How you will validate the database.

IMPORTANT:

Do not implement the database yet.

Do not implement Go code.
Do not implement Redis logic.
Do not implement Kafka producers/consumers.
Do not implement ML.
Do not implement FastAPI.
Do not implement RAG.
Do not implement the LLM.
Do not implement the frontend.

I want ONLY the Phase 2 database implementation plan.

If PROJECT_SPEC.md contains ambiguity, identify it rather than
inventing a new architecture.

Wait for my approval before making changes.