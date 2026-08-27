# AGENTS.md

# AI Risk Management Platform
## Coding-Agent Rules, Boundaries, Architecture Invariants, and Development Contract

---

# Source of Truth

The following files define the project:

1. `PROJECT_SPEC.md` — authoritative system specification
2. `docs/PHASES.md` — authoritative implementation sequence
3. `AGENTS.md` — rules for AI coding agents

Priority:

PROJECT_SPEC.md > docs/PHASES.md > AGENTS.md > agent assumptions

If an implementation decision is not specified, ask before making
a major architectural decision.

Do not invent new infrastructure, services, databases, frameworks,
or architectural patterns without justification.

## 1. Purpose of This File

PROJECT DOCUMENTATION RULE

PROJECT_SPEC.md is the canonical and authoritative specification
for the entire project.

Before implementing or modifying any feature:
1. Read PROJECT_SPEC.md.
2. Determine which component owns the responsibility.
3. Do not duplicate responsibility across services.
4. Do not introduce new infrastructure unless absolutely necessary.
5. Do not create additional large specification documents.
6. Do not redesign the architecture without explicit instruction.
7. Prefer the simplest implementation that satisfies PROJECT_SPEC.md.
8. Keep implementation modular and production-quality.
9. When uncertain, follow PROJECT_SPEC.md over assumptions.

---

# 2. Project Identity

Project name:

**Adaptive AI Risk Management Platform**

Primary purpose:

Build a realistic, demonstrable transaction-risk-management system for an AI Risk Management hackathon.

The platform simulates financial transactions and evaluates them using:

- Real-time behavioral signals
- Velocity detection
- Device intelligence
- IP intelligence
- Fraud-network relationships
- Supervised machine learning
- Unsupervised anomaly detection
- Explainable AI
- Customer-aware risk policies
- Step-up verification
- Retrieval-augmented generation
- Local LLM investigation
- Event-driven processing

The system MUST be designed as a coherent risk-management platform rather than as a collection of unrelated AI demos.

---

# 3. Core Architectural Principle

The following statement is a project invariant:

> **ML predicts. Policy decides. AI investigates and explains.**

A second core principle is:

> **The fraud graph reveals relationships.**

These principles MUST NOT be violated.

---

# 4. Responsibilities of the Major Components

## 4.1 Next.js

Next.js is responsible for:

- Dashboard UI
- Transaction monitoring
- Live transaction stream
- Risk visualization
- Transaction detail pages
- Fraud-network graph visualization
- Scenario controls
- Challenge/OTP UI
- AI investigation UI
- Risk timeline
- System health/status visualization
- Human-readable explanations

Next.js MUST NOT:

- Directly access PostgreSQL
- Directly access Redis
- Directly access Kafka
- Execute ML models
- Execute the LLM
- Contain fraud-decision logic that belongs to the backend
- Become the source of truth for risk decisions

Frontend decisions MUST be derived from backend responses/events.

---

# 5. Go Backend Responsibilities

Go is the primary backend and orchestration layer.

Go is responsible for:

- HTTP API
- WebSocket/SSE real-time communication
- Transaction lifecycle
- Scenario orchestration
- Transaction generation
- Redis interaction
- Kafka interaction
- PostgreSQL interaction
- Risk-pipeline orchestration
- Policy-engine execution
- Challenge/OTP lifecycle
- AI investigation orchestration
- Audit trail creation
- Service-to-service communication
- Authentication/authorization if implemented
- Configuration
- Error handling

Go is the primary owner of the transaction lifecycle.

The Go backend MUST NOT delegate the entire business logic to the LLM.

---

# 6. FastAPI Responsibilities

FastAPI is the Python AI/ML service.

FastAPI is responsible for:

- ML inference
- ML feature processing where appropriate
- XGBoost inference
- Isolation Forest inference
- Risk aggregation
- SHAP explanation generation
- AI investigation orchestration
- RAG retrieval
- LLM interaction
- AI tool execution where appropriate

FastAPI MUST expose typed, documented APIs.

FastAPI MUST NOT become the main transaction orchestration service.

Go remains the system orchestrator.

---

# 7. ML Responsibilities

The ML layer predicts risk.

It MAY:

- Calculate fraud probability
- Calculate anomaly score
- Generate model-derived risk scores
- Generate feature importance
- Generate SHAP explanations
- Consume historical behavioral features
- Consume real-time features
- Consume fraud-network-derived numerical features

The ML layer MUST NOT:

- Directly authorize a payment
- Directly block a customer
- Directly request OTP
- Directly decide the final transaction outcome
- Use an LLM as a hidden decision-maker

ML output should be an input to the policy engine.

---

# 8. Policy Engine Responsibilities

The policy engine is the authoritative authorization-decision layer.

It decides:

- ALLOW
- CHALLENGE
- BLOCK

based on:

- ML risk score
- Customer risk profile
- Transaction context
- Business rules
- Device trust
- IP signals
- Challenge state
- Step-up verification results
- Configured thresholds

The policy engine MUST be deterministic and auditable.

The same input state and policy configuration should produce the same decision.

The policy engine MUST NOT depend on an LLM for authorization.

---

# 9. AI/LLM Responsibilities

The LLM is an investigation and explanation system.

It may:

- Investigate suspicious transactions
- Summarize evidence
- Explain risk factors
- Retrieve relevant fraud cases
- Retrieve investigation playbooks
- Analyze customer history supplied through tools
- Analyze device relationships
- Analyze IP relationships
- Analyze fraud-network evidence
- Compare a case with historical patterns
- Recommend investigation actions
- Produce analyst-facing reports

The LLM MUST NOT be the authoritative payment decision-maker.

The LLM MUST NOT directly override:

- Policy Engine decisions
- Risk thresholds
- OTP requirements
- Authorization rules

The LLM may provide recommendations, but deterministic policy remains authoritative.

---

# 10. Redis Responsibilities

Redis is the real-time state and velocity engine.

Redis may maintain:

- Transaction velocity
- Short-term counters
- Sliding windows
- IP-to-account relationships
- Device-to-account relationships
- Recent transaction state
- Challenge state
- Temporary OTP state
- Short-lived risk signals

Redis MUST NOT be treated as the durable source of truth.

Important historical records MUST be stored in PostgreSQL.

---

# 11. Kafka Responsibilities

Kafka is the event-streaming layer.

Kafka is used for:

- Decoupling transaction generation from downstream processing
- Event propagation
- Risk-processing events
- Audit events
- AI investigation events
- Analytics events
- Future extensibility

Kafka MUST NOT replace PostgreSQL.

Kafka events should represent meaningful domain events.

Example:

```text
transaction.created
risk.evaluated
challenge.created
challenge.completed
transaction.decisioned
ai.investigation.requested
ai.investigation.completed

```


# 12. PostgreSQL Responsibilities

PostgreSQL is the durable source of truth.

It stores:

Customers
Customer risk profiles
Merchants
Devices
IP intelligence
Transactions
Risk assessments
Risk factors
Decisions
Challenges
Fraud cases
Scenario runs
AI investigations
Audit information
RAG documents
Embeddings if pgvector is used

PostgreSQL MUST remain authoritative for historical state.


# 13. Fraud Network Responsibilities

The fraud network represents relationships between:

Customers
Devices
IP addresses
Merchants
Payment instruments
Locations
Transactions

The graph MUST NOT exist solely as a visualization.

The system SHOULD derive graph features such as:
device_account_count
ip_account_count
fraud_neighbor_count
shared_device_accounts
shared_ip_accounts
These graph-derived features MAY feed the ML layer.

The frontend MUST provide a graph visualization for transaction investigation.

# 14. RAG Responsibilities

RAG belongs primarily to the AI/LLM investigation layer.

RAG provides semantic/domain knowledge such as:

Historical fraud cases
Fraud patterns
Investigation playbooks
Risk-management policies
Attack-pattern descriptions
Analyst guidance

RAG MUST NOT be used as a replacement for real-time transactional state.

Real-time customer/device/IP information SHOULD be retrieved through deterministic tools or structured queries.

# 15. Historical Context vs RAG

The following distinction MUST be preserved.

Structured historical context

Used by ML:

customer_average_amount
transaction_count_24h
transactions_last_60s
last_transaction_risk
previous_block_count
device_age
ip_account_count
Semantic historical knowledge

Used by RAG:

historical fraud cases
fraud investigation procedures
attack patterns
analyst playbooks
risk policy documents

Do NOT blindly put all transaction history into the vector database.
# 16. OTP / Step-Up Verification

The system will implement a simulated OTP challenge.

An external OTP provider is NOT required.

The backend will:

Generate OTP
Store it temporarily
Start a 60-second expiration window
Display the challenge in the frontend
Accept user input
Validate OTP
Record success/failure/expiry
Generate a new risk state
Re-run policy evaluation

The OTP should be visibly simulated for the hackathon.

Example:

OTP: 381924

The frontend may display a demo OTP or provide a safe demo mechanism.

Do not implement real financial authentication.

# 17. Adaptive Risk Policy

Do NOT hardcode:

65 = challenge
85 = block
globally.

Instead:

ML → standardized risk score 0-100

Policy Engine:
risk score
+
customer risk profile
+
transaction context
+
policy configuration
→ decision

Customer profiles may include:

risk_tolerance
challenge_threshold
block_threshold
account_age
trusted_devices
typical_transaction_range
historical_behavior

Thresholds should represent a risk policy profile rather than arbitrary user-provided values.

# 18. Risk Score Contract

The ML layer MUST produce a normalized risk score:

0 - 100

Interpretation:

0   = extremely low risk
100 = extremely high risk

The exact model probabilities MUST NOT be directly exposed as final authorization decisions.

Example:

{
  "fraud_probability": 0.78,
  "anomaly_score": 0.91,
  "risk_score": 83
}
# 19. Risk Timeline

Every transaction should have an auditable risk timeline.

Example:

TRANSACTION CREATED
        ↓
REDIS SIGNALS
        ↓
ML EVALUATION
        ↓
RISK SCORE
        ↓
POLICY DECISION
        ↓
CHALLENGE
        ↓
OTP SUCCESS
        ↓
UPDATED RISK
        ↓
FINAL POLICY DECISION
        ↓
AI INVESTIGATION

The frontend should be able to visualize this timeline.

# 20. Scenario Engine

The system must contain a reusable scenario engine.

Do NOT implement each attack as an unrelated generator.

Use:

BaseTransactionGenerator
        ↓
Scenario
        ↓
Scenario modifiers
        ↓
Transactions

Required initial scenarios:

Normal Traffic
Velocity Attack
Account Takeover
Device Reuse

Optional scenarios:

Fraud Ring
High Value Anomaly
Card Testing

The scenario engine should support:

Number of transactions
Duration
Speed
Seed
Customer selection
Scenario type
# 21. Simulation Requirements

The simulator should support:

100-200 transactions

per demonstration.

Transactions should appear to arrive in real time.

The simulator may:

Emit some transactions immediately
Wait approximately 1 second between others
Produce bursts
Produce attack phases
Produce normal phases

The system MUST support deterministic seeded scenarios for reproducible demonstrations.

# 22. Demo Reliability Rule

The demo is more important than maximum randomness.

For every scenario:

Use a known seed
Produce predictable risk progression
Ensure expected outcomes occur
Avoid random behavior that could cause a failed demonstration

The system should have:

START SCENARIO
STOP SCENARIO
RESET SCENARIO
# 23. No Fake Claims

The application MUST NOT claim:

Production fraud detection accuracy
Real payment authorization
Real banking connectivity
Real OTP delivery
Real-world fraud prevention guarantees
Regulatory certification

The project is a research/hackathon prototype.

If synthetic data is used, the UI/documentation should identify it as synthetic.

# 24. Security Rules

Never:

Hardcode secrets
Commit API keys
Commit passwords
Commit tokens
Store secrets in frontend code

Use environment variables.

Provide:

.env.example

with safe placeholder values.

# 25. Privacy Rules

All demo data must be synthetic.

Do not use:

Real customer data
Real payment card numbers
Real personally identifiable information
Real authentication credentials

Use synthetic identifiers.

Examples:

C1001
D204
IP-17
TX-10982
M301
# 26. Technology Constraints

Primary stack:

Frontend
Next.js
React
TypeScript
Tailwind CSS
Backend
Go
AI/ML
Python
FastAPI
XGBoost
scikit-learn
SHAP
LLM
Ollama
Local open-source instruction model

The exact model is specified in 12_OLLAMA.md.

Database
PostgreSQL
pgvector where required
Real-time state
Redis
Event streaming
Kafka
Infrastructure
Docker
Docker Compose
# 27. Infrastructure Simplicity Rule

Do not introduce unnecessary infrastructure.

Avoid creating:

Kubernetes
Service mesh
API gateway products
Multiple databases
Multiple queues
Multiple LLM providers
Neo4j unless explicitly approved
Cloud services unless explicitly approved

The project must remain buildable within the hackathon timeline.
# 28. Monorepo Rule

The project MUST use a monorepo structure.

Recommended:

/frontend
/backend
/ai
/database
/docs
/simulator

Do not create separate repositories for each service.

# 29. Service Boundary Rule

The following service boundaries are fixed:

Next.js
    ↓
Go Backend
    ↓
FastAPI AI/ML
    ↓
Ollama

Infrastructure:

PostgreSQL
Redis
Kafka

The agent MUST NOT convert every internal module into a separate microservice.

# 30. Go Internal Modules

The Go backend should have logical modules such as:

api
transaction
scenario
risk
policy
challenge
redis
kafka
postgres
ai
websocket

These are internal packages/modules.

They do NOT automatically imply independent network services.

# 31. API Design

APIs must:

Use JSON
Have explicit schemas
Validate input
Return meaningful errors
Use appropriate HTTP status codes
Be versionable where appropriate
Avoid leaking internal errors

Use:

/api/v1/...

unless another convention is explicitly specified.

# 32. Frontend Rule

The frontend must not duplicate backend business logic.

Bad:

if (risk > 85) {
    decision = "BLOCK";
}

Good:

decision = transaction.decision;

The backend owns the decision.

33. AI Tool Rule

AI tools must be:

Explicitly defined
Typed
Read-only by default
Auditable

Examples:

get_customer_history
get_device_profile
get_ip_profile
get_related_accounts
get_similar_fraud_cases
get_customer_risk_profile

The LLM MUST NOT receive unrestricted database access.

# 34. AI Output Rule

AI investigation output MUST be structured.

Preferred:

{
  "summary": "...",
  "risk_level": "HIGH",
  "evidence": [],
  "similar_cases": [],
  "recommendation": "ESCALATE",
  "confidence": 0.88
}

Do not rely exclusively on free-form text.

# 35. AI Hallucination Rule

The LLM must distinguish between:

Retrieved evidence
Model inference
Recommendation

The LLM must not invent:

Transaction history
Fraud cases
Device relationships
IP relationships
Customer attributes

If evidence is unavailable, it must say so.

# 36. Error Handling

Every service must gracefully handle:

Database unavailable
Redis unavailable
Kafka unavailable
FastAPI unavailable
Ollama unavailable
Invalid transaction
ML inference error
Invalid OTP
OTP expiry
Scenario failure

A single AI failure MUST NOT destroy the core transaction pipeline.

# 37. Degraded Mode

The system should ideally remain useful if the AI investigator is unavailable.

For example:

ML works
Policy works
Transaction decisions work
AI investigation unavailable

The UI should show:

AI Investigator unavailable

rather than failing the entire transaction.

# 38. Testing Rules

Every major subsystem must have at least basic tests.

Required:

Policy tests
Risk aggregation tests
Scenario generator tests
Redis velocity tests
OTP tests
API tests
ML inference smoke test
AI output schema test

# 39. Documentation Rules

Every major subsystem must have a corresponding specification.

Required documentation:

docs/
├── 00_PROJECT_OVERVIEW.md
├── 01_ARCHITECTURE.md
├── 02_DATABASE.md
├── 03_REDIS.md
├── 04_KAFKA.md
├── 05_GO_BACKEND.md
├── 06_SCENARIO_ENGINE.md
├── 07_ML_SYSTEM.md
├── 08_FASTAPI.md
├── 09_FRAUD_GRAPH.md
├── 10_AI_INVESTIGATOR.md
├── 11_RAG.md
├── 12_OLLAMA.md
├── 13_POLICY_ENGINE.md
├── 14_CHALLENGE_OTP.md
├── 15_NEXTJS.md
├── 16_WEBSOCKET_EVENTS.md
├── 17_DOCKER.md
├── 18_TESTING.md
└── 19_DEMO_SCRIPT.md

# 40. Implementation Order

The agent should implement in this order unless explicitly instructed otherwise:

Phase 1

Infrastructure:

Docker
PostgreSQL
Redis
Kafka
Phase 2

Core Go backend:

transaction
scenario engine
database
Redis
Kafka
Phase 3

Frontend:

dashboard
transaction list
live stream
Phase 4

ML:

feature pipeline
XGBoost
Isolation Forest
SHAP
FastAPI
Phase 5

Policy:

risk profiles
thresholds
ALLOW
CHALLENGE
BLOCK
Phase 6

Challenge:

OTP
60-second timer
verification
risk update
re-evaluation
Phase 7

Fraud graph:

relationship extraction
graph features
graph visualization
Phase 8

AI:

tools
RAG
Ollama
investigation
Phase 9

Polish:

animations
risk timeline
logs
graphs
scenario UX
demo reliability

# 41. Agent Workflow

Before changing code:

Read relevant specification.
Inspect existing repository.
Determine whether implementation already exists.
Do not duplicate functionality.
Implement the smallest compliant change.
Run tests.
Run formatting/linting.
Verify integration.
Report changed files.
Report tests performed

# 42. Architecture Change Rule

The agent MUST NOT silently change:

Technology stack
Service boundaries
Database choice
Message broker
ML architecture
LLM architecture
Policy ownership
Decision authority

If a change appears necessary:

Explain the problem.
Explain the proposed change.
Explain the impact.
Wait for explicit approval.

# 43. Dependency Rule

Before introducing a new dependency, ask:

Is it necessary?
Can the existing stack solve the problem?
Does it add infrastructure?
Does it increase hackathon risk?
Is it consistent with the architecture?

Prefer existing dependencies.

# 44. Code Quality Rule

Prefer:

Simple code
Clear names
Small functions
Explicit interfaces
Strong typing
Meaningful errors
Testable modules

Avoid:

Overengineering
Excessive abstraction
Premature optimization
Unnecessary design patterns
Generated boilerplate with no purpose

# 45. Demo-First Principle

This is a hackathon project.

The system must be:

Reliable > complicated
Observable > clever
Explainable > opaque
Demonstrable > theoretical

Every feature should have a visible purpose in the demo.

# 46. Final Invariants

The following statements MUST remain true:


```text
ML predicts.
Policy decides.
AI investigates and explains.
Fraud graph reveals relationships.
Redis provides real-time state.
PostgreSQL provides durable state.
Kafka provides event streaming.
Go orchestrates the transaction lifecycle.
FastAPI hosts AI/ML workloads.
Ollama hosts the local LLM.
RAG provides semantic domain knowledge.
Structured tools provide live customer context.
OTP provides step-up verification.
Challenge results feed back into policy evaluation.
The LLM cannot directly authorize transactions.
Synthetic data is used.
The system must remain demonstrable and reproducible.

```