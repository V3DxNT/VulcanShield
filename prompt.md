You are the implementation agent for the VulcanShield project.

IMPORTANT:
Before doing anything, read these files completely:

1. AGENTS.md
2. PROJECT_SPEC.md
3. docs/PHASES.md

PROJECT_SPEC.md is approximately 2600 lines long.

You MUST read the complete file. Do not skip sections because of its
length.

However:

- Do not reproduce the specification in your response.
- Do not summarize all 2600 lines unnecessarily.
- Use the files as the authoritative source of truth.
- Do not invent architecture that conflicts with them.

==================================================
CURRENT TASK
==================================================

We are starting PHASE 0.

You MUST work ONLY on PHASE 0.

Do NOT implement Phase 1 or any later phase.

Do NOT build:

- Redis logic
- Kafka pipelines
- ML models
- FastAPI ML service
- LLM
- RAG
- Fraud graph
- Policy engine
- OTP
- Dashboard functionality
- Transaction generation
- AI investigation

Those belong to later phases.

==================================================
BEFORE MODIFYING FILES
==================================================

First inspect the existing repository.

Tell me:

1. What files currently exist.
2. What the current Git state is.
3. What Phase 0 requires according to docs/PHASES.md.
4. What parts of Phase 0 are already satisfied.
5. What files/directories you intend to create.
6. What technologies/tools Phase 0 requires.
7. How you will validate Phase 0.
8. Any contradictions between AGENTS.md, PROJECT_SPEC.md,
   and docs/PHASES.md.

Do not modify files yet.

Wait for my approval.

==================================================
IMPLEMENTATION RULES
==================================================

When I approve the plan:

- Implement Phase 0 only.
- Keep the architecture consistent with PROJECT_SPEC.md.
- Do not add unnecessary dependencies.
- Do not introduce new frameworks.
- Do not introduce Kubernetes.
- Do not introduce cloud infrastructure.
- Do not introduce additional databases.
- Do not introduce unnecessary microservices.
- Do not create placeholder implementations for future phases.
- Do not implement future functionality "for convenience".

Every dependency added must have a clear reason.

Prefer a simple local development environment.

==================================================
QUALITY REQUIREMENTS
==================================================

This project is being built for a technical AI/Risk Management
hackathon and internship evaluation.

Therefore:

- Code must be production-quality where practical.
- Configuration must be explicit.
- Services must have clear boundaries.
- Errors must be handled properly.
- Logs should be useful.
- Environment variables should be used for configuration.
- Secrets must never be committed.
- Do not hardcode credentials.
- Do not leave broken placeholder code.
- Do not silently ignore errors.

==================================================
VALIDATION
==================================================

After implementation:

1. Run the relevant formatters.
2. Run static checks/linting where applicable.
3. Build the relevant components.
4. Start the required infrastructure/services.
5. Verify health checks.
6. Verify connectivity between components where Phase 0
   requires it.
7. Fix any errors you introduced.
8. Report exactly what was tested.

Do not claim something works unless you actually tested it.

==================================================
GIT SAFETY
==================================================

Do not:

- reset the repository
- delete existing work
- rewrite Git history
- force push
- change the remote
- commit automatically

If you believe a destructive operation is necessary,
STOP and ask me first.

==================================================
COMPLETION
==================================================

When Phase 0 is complete, report:

1. Files created/modified.
2. Dependencies added.
3. Commands executed.
4. Tests/checks performed.
5. Results.
6. Any known issues.
7. Exact instructions for running the project.
8. Whether Phase 0 is fully complete.

Do not proceed to Phase 1.

Wait for my explicit approval before continuing.