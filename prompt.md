Phase 0 has been completed, validated, committed, and pushed.

Current repository state is the result of:

Phase 0 — Repository Initialization

Now we are beginning:

PHASE 1 — Infrastructure

Before modifying anything:

1. Read AGENTS.md.
2. Read PROJECT_SPEC.md.
3. Read docs/PHASES.md.
4. Read docs/DECISIONS.md if it exists.
5. Inspect the current repository.
6. Inspect the existing docker-compose.yml and environment
   configuration created during Phase 0.

Do NOT assume that Phase 0 is correct merely because it completed.
Verify the current repository state.

==================================================
PHASE 1 OBJECTIVE
==================================================

The objective of Phase 1 is to establish a reliable local
infrastructure environment for VulcanShield.

The infrastructure components are:

1. PostgreSQL
2. pgvector
3. Redis
4. Kafka
5. Ollama

These services must be runnable locally using Docker/Docker Compose
where appropriate.

==================================================
IMPORTANT BOUNDARY
==================================================

Implement PHASE 1 ONLY.

Do NOT implement:

- Go backend
- transaction generator
- Redis velocity logic
- Kafka transaction processing
- ML service
- FastAPI ML endpoints
- XGBoost
- Isolation Forest
- policy engine
- OTP
- fraud graph
- Next.js application
- WebSockets
- RAG pipeline
- embeddings generation
- AI investigation
- LLM tools

Phase 1 establishes infrastructure only.

==================================================
POSTGRESQL
==================================================

Configure PostgreSQL for local development.

Requirements:

- PostgreSQL must be accessible by other services.
- Configuration must use environment variables.
- Do not hardcode credentials.
- PostgreSQL must support pgvector.
- Use a suitable pgvector-enabled PostgreSQL image.
- Configure health checks.
- Keep the database local and disposable.
- Do not introduce persistent production infrastructure.

Do NOT implement application tables yet unless explicitly required
by PROJECT_SPEC.md for infrastructure initialization.

Database schema belongs to the database phase.

==================================================
REDIS
==================================================

Configure Redis for local development.

Requirements:

- Redis must start successfully.
- Redis must be reachable from the local environment.
- Configure health checks.
- Use environment configuration.
- Do not implement velocity detection yet.
- Do not implement OTP yet.

Redis logic belongs to later phases.

==================================================
KAFKA
==================================================

Configure Kafka for local development.

Requirements:

- Kafka must start successfully.
- The configuration must be appropriate for local development.
- Other local services must be able to connect to Kafka.
- Configure health checks where practical.
- Keep the configuration as simple as possible.

Do not implement producers or consumers yet.

Do not create application event-processing logic yet.

==================================================
OLLAMA
==================================================

Configure Ollama as the local LLM runtime.

Requirements:

- Ollama must be runnable locally.
- It must be reachable from the AI service environment later.
- Do not implement the AI service yet.
- Do not create RAG.
- Do not create prompts.
- Do not add application-level LLM orchestration.

If pulling a model during Docker startup would make the development
environment slow or unreliable, do NOT automatically download a large
model.

Instead, document the exact model setup command in the README.

==================================================
DOCKER COMPOSE
==================================================

Update docker-compose.yml as required.

Requirements:

- Clear service names.
- Internal networking.
- Health checks where practical.
- Environment variables instead of hardcoded secrets.
- Reasonable restart behavior for local development.
- Avoid unnecessary infrastructure.
- Avoid Kubernetes.
- Avoid cloud services.
- Avoid additional databases.

Do not add infrastructure that is not required by PROJECT_SPEC.md.

==================================================
ENVIRONMENT
==================================================

Update .env.example with safe development placeholders.

Do NOT create or commit .env.

Ensure .gitignore continues to exclude:

.env
secrets
credentials
local model/cache data where appropriate
logs
build artifacts
IDE files

==================================================
VALIDATION
==================================================

After implementation, actually test the infrastructure.

At minimum verify:

1. docker compose configuration is valid.
2. PostgreSQL starts.
3. PostgreSQL accepts connections.
4. pgvector is available.
5. Redis starts.
6. Redis accepts connections.
7. Kafka starts.
8. Kafka is reachable.
9. Ollama starts.
10. Ollama is reachable.

Where practical, test connectivity using actual commands rather
than only checking container status.

If something fails, debug and fix it.

Do not claim success based only on configuration inspection.

==================================================
RESOURCE CONSTRAINT
==================================================

This is a local hackathon project.

Do not over-engineer the infrastructure.

Prefer:

simple
reliable
reproducible
easy to restart

over:

distributed production-grade infrastructure.

==================================================
COMPLETION
==================================================

When Phase 1 is complete, report:

1. Files modified.
2. Infrastructure services created.
3. Docker images used.
4. Environment variables added.
5. Commands used to start the infrastructure.
6. Health-check results.
7. Connectivity test results.
8. Any known limitations.
9. Exact commands I should use to start/stop the infrastructure.

Do NOT proceed to Phase 2.

Do NOT commit changes automatically.

Wait for my approval after reporting completion.