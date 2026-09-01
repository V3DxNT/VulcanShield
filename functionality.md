# VulcanShield Functional Architecture

## 1. System Purpose

VulcanShield is a synthetic, demo-first AI risk management platform built to simulate real-time financial fraud monitoring. It combines transaction generation, behavioral analysis, risk scoring, policy enforcement, graph-based fraud relationships, and AI-driven investigation into a single operational system.

The project follows the core architecture invariant:

- ML predicts
- Policy decides
- AI investigates and explains

This separation keeps the application explainable and demonstrable without pretending to be a production banking system.

---

## 2. High-Level Architecture

```text
Browser / User
   |
   v
Next.js Frontend
   |
   | REST API + WebSocket updates
   v
Go Backend
   |
   | - orchestrates transaction lifecycle
   | - stores and reads durable state
   | - manages Redis and Kafka
   | - calls ML and AI services
   v
+--------------------------+
|      Shared Services      |
|  PostgreSQL  Redis  Kafka |
+--------------------------+
   |
   +------> ML Service (Python / FastAPI)
   |             - feature extraction
   |             - XGBoost scoring
   |             - Isolation Forest scoring
   |
   +------> AI Service (Python / FastAPI)
   |             - customer/device/IP context retrieval
   |             - graph-aware evidence gathering
   |             - prompt creation for LLM
   |             - structured investigation output
   |
   +------> Ollama (local LLM)
                 - investigation reasoning
                 - explanation generation
                 - analyst summary
```

---

## 3. Service Breakdown

### 3.1 Frontend: Next.js + TypeScript + Tailwind

The frontend is the visual layer of the application. It displays:

- transaction stream and live risk activity
- dashboard KPIs and KPI cards
- scenario controls for demo simulation
- risk timeline and decision explanation
- fraud graph/network visualization
- OTP and challenge flow UI
- AI investigation modal with evidence and prompt visibility

It does not decide risk. It consumes backend responses and displays them clearly for operators and demo viewers.

Why it matters:

- makes the platform observable in real time
- keeps the user interface focused on interpretation rather than logic
- communicates risk, graph relationships, and AI reasoning in a demo-friendly way

---

### 3.2 Go Backend: system orchestrator

The Go backend is the central coordination layer. It manages:

- API endpoints
- transaction lifecycle orchestration
- scenario generation and scenario control
- Redis interaction for velocity and temporary state
- Kafka event producers and event flow
- PostgreSQL persistence for durable records
- policy engine integration
- AI service orchestration
- WebSocket streaming to the frontend

This is the main runtime engine that keeps the application coherent.

Why it matters:

- it provides a single source of operational truth for the transaction pipeline
- it keeps business logic centralized and deterministic
- it isolates the frontend from low-level data and infrastructure concerns

---

### 3.3 ML Service: Python FastAPI model layer

The ML service performs risk scoring and anomaly detection. It is responsible for:

- building input feature vectors
- running XGBoost-based fraud probability estimation
- running Isolation Forest anomaly scoring
- returning structured risk outputs with explanation-friendly values
- feeding the backend risk pipeline

This is where model inference happens. It does not decide whether a payment is allowed or blocked. It predicts risk.

Why it matters:

- produces measurable risk probabilities and anomaly scores
- gives the policy engine a quantitative factor to evaluate
- makes the risk logic explainable and auditable

---

### 3.4 AI Service: contextual investigation layer

The AI service is responsible for investigation and explanation. It gathers structured evidence such as:

- recent customer transaction history
- customer risk profile
- device and IP intelligence
- related account and fraud-network context
- previous suspicious activity history
- similar historical patterns or case context

It then assembles the prompt sent to the local LLM and returns a structured investigation payload.

Why it matters:

- reduces hallucination by grounding the LLM in retrievable, traceable evidence
- keeps the LLM focused on explanation instead of authorization
- turns raw data into analyst-friendly investigation summaries

---

### 3.5 Ollama: local LLM inference

Ollama runs the local model used for investigation. It receives a prompt built from the retrieved customer, device, IP, and transaction context and produces the final summary and risk narrative.

It does not make the payment decision.

Why it matters:

- keeps the system fully local and demo-friendly
- avoids dependence on external cloud APIs
- allows the LLM to explain suspicious patterns using structured context

---

### 3.6 PostgreSQL: durable history and state

PostgreSQL is the durable source of truth for the application.

It stores:

- customers and profiles
- merchants and transaction data
- risk assessments and scores
- challenge and OTP records
- graph relationships
- historical investigation data
- scenario state and audit context

Why it matters:

- retains long-term, reliable state
- supports graph analysis and historical retrieval
- provides the durable foundation necessary for reproducible demos and audits

---

### 3.7 Redis: real-time velocity and ephemeral state

Redis supports short-lived, high-speed state for:

- transaction velocity windows
- recent activity counters
- temporary challenge data
- stateful risk signals
- quick event lookups during live processing

Why it matters:

- improves responsiveness for real-time monitoring
- supports fast detection of bursts or suspicious rapid behavior
- keeps high-frequency state separate from durable storage

---

### 3.8 Kafka: event-driven communication

Kafka decouples services by carrying important events such as:

- transaction creation
- risk evaluation
- challenge creation
- decisioning
- investigation requests
- analytics and event propagation

Why it matters:

- allows event-driven processing without tight coupling
- supports scalable future extension without redesigning the app
- keeps the backend and downstream services modular

---

## 4. How the Services Connect to Each Other

### Request flow

1. The frontend loads the dashboard and subscribes to updates.
2. The Go backend receives requests or orchestrates scenario generation.
3. Transactions are created and recorded in PostgreSQL.
4. Redis captures short-term velocity and real-time behavioral state.
5. Kafka emits transaction and risk events.
6. The Go backend calls the ML service for risk and anomaly scores.
7. The backend evaluates the policy engine and decides ALLOW, CHALLENGE, or BLOCK.
8. If challenge is involved, the OTP flow is executed and re-evaluated.
9. The AI service gathers transaction context, customer history, graph signals, and device/IP intelligence.
10. The AI service sends the structured prompt to Ollama.
11. Ollama produces an explanation and investigation summary.
12. The result is returned to the frontend for display and the graph and transaction UI are refreshed.

### Data boundary

- Redis = fast transient state
- PostgreSQL = durable historical truth
- Kafka = event streaming and decoupling
- ML = predictive scoring
- Policy = final decision authority
- AI = explanation and investigation

This keeps responsibilities clean and prevent the LLM from acting as an unreviewed decision-maker.

---

## 5. Redis, Kafka, and PostgreSQL/pgvector in Depth

### 5.1 Redis and velocity attacks

Redis is used for short-lived, high-speed behavioural signals. In this project it does not store the canonical transaction record. Instead it keeps rolling, time-windowed counters and sorted sets that answer questions such as:

- How many transactions did this customer make in the last 60 seconds?
- How many transactions came from this IP in the last minute?
- How many transactions originated from this device in the previous minute?

The implementation uses Redis sorted sets keyed by entity, such as:

- `velocity:user:<user_id>`
- `velocity:ip:<ip_address>`
- `velocity:device:<device_id>`

Each transaction is recorded with a timestamp as its score and a composite member string such as `transaction_id:amount`. When a new transaction arrives, the system calls `RecordTransaction` and then immediately asks `GetVelocitySignals` for the count in the last 60 seconds. Those numbers directly feed the risk evaluator, where a burst of transactions adds extra risk through the velocityScore calculation.

Why this matters for a velocity attack:

- a card tester or bot can send many low-value transactions in a short window
- Redis detects that burst in real time because it is optimized for fast time-windowed lookups
- the backend combines that Redis signal with ML fraud probability and anomaly score
- the final risk score then decides whether the transaction should be allowed, challenged, or blocked

This is exactly why Redis is used for real-time velocity checks rather than for long-term reporting. It is not the permanent system of record; it is the fast detection engine.

### 5.2 Kafka event flow

Kafka is used to decouple the system and propagate meaningful domain events. The generator emits events such as:

- transaction creation
- risk evaluation
- policy decision
- challenge creation
- AI investigation requests

This gives the system a modular event pipeline rather than a single tightly-coupled code path. In practice, Kafka keeps the pipeline observable and extensible: one service may create or emit domain events, while downstream consumers or monitoring components can react without needing to know the full internal implementation.

The important architectural point is that Kafka is not the durable source of truth. PostgreSQL remains the canonical transactional store. Kafka is the message bus that carries state transitions and downstream processing signals.

### 5.3 PostgreSQL and pgvector for durable and semantic memory

PostgreSQL stores the durable state of the platform, including:

- users and policy thresholds
- transactions and risk results
- challenge history and OTP state
- fraud relationship graph edges
- investigations
- RAG document collection
- vector embeddings used for semantic search

The critical extension here is `pgvector`, which lets the project store and query embeddings in PostgreSQL rather than creating a separate vector database. In this project, `embedding_records` stores 384-dimensional embeddings for documents in `rag_documents`, and the AI service performs the semantic retrieval using the vector distance operator `<=>`.

This means:

- the RAG documents remain in the same relational database as other transactional data
- the AI service can search for semantically related fraud-playbook content
- the LLM gets a grounded, retrievable context layer to explain the case without treating the raw transaction and policy engine as a free-form guess

The system intentionally uses an architecture where:

- structured data (customer history, device profile, IP profile, fraud graph) is retrieved with deterministic SQL queries
- semantic knowledge (playbooks, fraud patterns, attack descriptions) is retrieved with pgvector search

This is the clean separation between operational truth and long-term domain knowledge.

---

## 6. Why the Application Is Stronger

This architecture makes the application stronger because it combines the right tools for the right purpose:

- Next.js gives a clear operational dashboard and real-time observability
- Go gives a reliable orchestration layer and API foundation
- Python ML models provide measurable fraud predictions
- Policy logic ensures deterministic and auditable decisions
- PostgreSQL preserves durable state and historical context
- Redis handles real-time velocity and challenge responsiveness
- Kafka supports decoupling and future extensibility
- AI service and Ollama provide explainable investigation without forcing the LLM to make payment decisions
- Graph relationships expose suspicious connection patterns among users, devices, and IPs

The result is not just a demo; it is a structured fraud-risk platform with layered intelligence and explainability.

---

## 6. Demo Integrity and Constraints

This project uses synthetic data only. It is designed for demonstration and education, not for real banking or regulatory deployment.

The architecture intentionally preserves:

- synthetic data boundaries
- explainability over black-box behavior
- auditable investigation flows
- deterministic scenario generation for repeatable demos
- a clear separation between model prediction and policy enforcement

---

## 7. Summary

VulcanShield is stronger because each service has a clear responsibility and operates in a connected system that mirrors a realistic fraud-risk workflow:

- frontend visualizes the story
- backend orchestrates the business flow
- ML predicts risk
- policy decides outcomes
- graph reveals relationships
- AI explains suspicious signals
- infrastructure supports speed, durability, and real-time operations

This combination turns the platform into a cohesive, explainable, and demonstrable risk intelligence system.
