# PROJECT_SPEC.md

# AI RISK MANAGEMENT PLATFORM
## Canonical System Specification

Version: 1.0
Status: FINAL ARCHITECTURE
Purpose: Buildathon MVP + Internship-Level Demonstration

---

# 0. DOCUMENT AUTHORITY

This is the SINGLE SOURCE OF TRUTH for the project.

The implementation agent MUST read this file before modifying or creating
project components.

If another document conflicts with this document, this document wins.

The agent MUST NOT invent major architectural components that are not
specified here unless required for implementation correctness.

The agent MAY make small implementation-level decisions when multiple
technically equivalent approaches exist.

The agent MUST preserve the architectural responsibilities defined here.

---

# 1. PROJECT OBJECTIVE

Build an end-to-end AI-powered payment risk management platform that
simulates real-time fintech transaction processing.

The system must demonstrate:

1. Real-time transaction ingestion
2. Redis-based velocity detection
3. Kafka-based event processing
4. Historical transaction storage
5. Machine-learning fraud prediction
6. Policy-based transaction decisions
7. Adaptive risk thresholds
8. OTP-based challenge flow
9. Fraud-network graph analysis
10. RAG-based contextual investigation
11. Local LLM-based investigation
12. Explainable risk decisions
13. Real-time dashboard visualization
14. Transaction-level investigation
15. Scenario-based fraud simulation
16. Live transaction streaming

The system is a DEMONSTRATION PLATFORM.

It does NOT connect to real payment networks.

All transactions are synthetic.

---

# 2. CORE ARCHITECTURAL PRINCIPLE

The project MUST maintain strict separation between:

    ML
    Policy
    AI

The fundamental rule is:

    ML PREDICTS.
    POLICY DECIDES.
    AI INVESTIGATES AND EXPLAINS.

Meaning:

ML:
    Estimates fraud/risk probability from transaction features.

Policy Engine:
    Converts risk and business/security signals into an action.

AI/LLM:
    Investigates why a transaction is suspicious using evidence from
    history, graph relationships, RAG and risk signals.

The LLM MUST NOT become the final transaction decision-maker.

---

# 3. HIGH-LEVEL ARCHITECTURE

```text
                        NEXT.JS
                           |
                           | REST / WebSocket
                           v
                     GO BACKEND
                           |
                  TRANSACTION GENERATOR
                           |
                           v
                         KAFKA
                           |
             +-------------+-------------+
             |             |             |
             v             v             v
          REDIS       PROCESSOR       POSTGRES
             |             |             |
             |             v             |
             |       FEATURE BUILDER     |
             |             |             |
             |             v             |
             |         ML SERVICE        |
             |             |             |
             +-------------+-------------+
                           |
                           v
                    POLICY ENGINE
                           |
                +----------+----------+
                |          |          |
                v          v          v
              ALLOW     CHALLENGE    BLOCK
                           |
                           v
                          OTP
                           |
                           v
                    RE-EVALUATION
                           |
                           v
                    FINAL DECISION
                           |
                           +----------------+
                           |                |
                           v                v
                       POSTGRES           KAFKA
                                            |
                                            v
                                    AI INVESTIGATION
                                            |
                         +------------------+----------------+
                         |                  |                |
                         v                  v                v
                    USER HISTORY        GLOBAL RAG         GRAPH
                         |                  |                |
                         +------------------+----------------+
                                            |
                                            v
                                       TOOL ROUTER
                                            |
                                            v
                                         OLLAMA
                                            |
                                            v
                                       LOCAL LLM
                                            |
                                            v
                                   AI INVESTIGATION
                                            |
                                            v
                                        POSTGRES
                                            |
                                            v
                                         NEXT.JS
4. TECHNOLOGY STACK
Frontend

Use:

Next.js
TypeScript
Tailwind CSS
shadcn/ui
Recharts or equivalent charting library

Responsibilities:

Dashboard
Transaction stream
Transaction details
Risk visualization
Policy explanation
OTP interaction
Fraud graph visualization
AI investigation display
Scenario controls
System metrics
5. GO BACKEND

Use:

Go
Gin or Fiber

Preferred:

Gin

Responsibilities:

API gateway
Transaction generation
Kafka producer
Redis integration
transaction orchestration
policy engine
OTP logic
graph coordination
frontend APIs
WebSocket event streaming

Go is the MAIN application backend.

6. REDIS

Redis is responsible for FAST, SHORT-LIVED state.

Primary purpose:

real-time risk signals

Use Redis for:

velocity windows
transaction counts
amount velocity
IP activity
device activity
temporary OTP state
short-lived risk state
recent transaction context

Do NOT use Redis as the permanent transaction database.

7. KAFKA

Kafka is responsible for EVENT STREAMING.

Kafka provides:

asynchronous processing
decoupling
event history during processing
scalable architecture

Important events:

transaction.created
transaction.processed
risk.calculated
policy.decided
challenge.created
challenge.completed
transaction.finalized
ai.investigation.requested
ai.investigation.completed
8. POSTGRESQL

PostgreSQL is the persistent source of truth.

Use Docker.

For the buildathon:

local PostgreSQL
no persistent Docker volume required

This ensures each clean demonstration can start from a fresh state.

PostgreSQL stores:

users
transactions
risk assessments
policy decisions
OTP challenges
fraud relationships
investigations
evidence
RAG documents
embeddings
scenarios
audit information
9. VECTOR STORAGE

Use:

pgvector

inside PostgreSQL.

Do NOT introduce Pinecone, Weaviate, Milvus, etc. for the MVP.

Reason:

PostgreSQL + pgvector reduces infrastructure complexity.
10. ML SERVICE

Use:

Python
FastAPI
scikit-learn
XGBoost

The ML service is separate from Go.

Go sends a feature vector to FastAPI.

FastAPI returns predictions.

11. ML MODELS

Use two models.

Model 1: XGBoost

Purpose:

supervised fraud prediction

Output:

fraud probability

Example:

0.84
Model 2: Isolation Forest

Purpose:

unsupervised anomaly detection

Output:

anomaly score

This identifies unusual behavior even if the supervised model
does not classify it strongly.

12. WHY TWO ML MODELS?

XGBoost answers:

"Does this look like known fraud?"

Isolation Forest answers:

"Does this behavior look unusual?"

Together:

known fraud patterns
+
behavioral anomalies

produce stronger risk signals.

13. ML IS NOT POLICY

ML output:

prediction

Policy output:

action

Example:

XGBoost probability = 0.79
Isolation Forest anomaly = 0.86

Combined risk = 78

Policy:
    CHALLENGE

The model does NOT directly return:

BLOCK
14. TRAINING STRATEGY

Do NOT train huge models.

Do NOT use deep learning for the MVP.

Use synthetic data.

Generate approximately:

10,000–50,000 training transactions

depending on available compute.

The transaction simulator may generate:

100–200 transactions

for the live demo.

Training data and demo data are separate.

15. SYNTHETIC DATA

Normal users should exhibit patterns such as:

stable transaction amounts
normal velocity
recurring devices
recurring locations
recurring IPs
normal merchant patterns

Fraud scenarios should introduce:

abnormal velocity
amount spikes
new devices
new IPs
device reuse
IP reuse
geographic anomalies
repeated failed attempts
suspicious account relationships
16. TRANSACTION GENERATOR

Transactions MUST NOT be manually typed for the demonstration.

Go provides scenario-based generators.

Frontend exposes:

Start Scenario

Go starts generating synthetic transactions.

Example:

Start Velocity Attack
        |
        v
Go Generator
        |
        v
Transaction 1
Transaction 2
Transaction 3
...
Transaction 20
        |
        v
Kafka
17. DEMO GENERATION

Normal demonstration:

100–200 transactions

Some should occur:

immediately

Others:

every ~1 second

This creates the appearance of a live payment stream.

Fraud scenarios inject specific patterns.

18. SCENARIOS

Minimum scenarios:

Normal Traffic
Velocity Attack
Account Takeover
Device Farm
IP Abuse
Amount Anomaly

Optional:

Geographic Anomaly
Merchant Abuse
19. VELOCITY ATTACK

Generate:

many transactions
same user
same IP/device
short period

Example:

14 transactions
within 60 seconds

Redis detects this immediately.

20. ACCOUNT TAKEOVER

Generate:

new device
new IP
unusual amount
unusual location
abnormal timing

Historical user behavior remains different.

This creates behavioral deviation.

21. DEVICE FARM

Create:

Device X
    |
    +-- User A
    +-- User B
    +-- User C
    +-- User D
    +-- User E

Some associated accounts should contain prior suspicious/fraud activity.

22. IP ABUSE

Create:

IP X
 |
 +-- User A
 +-- User B
 +-- User C
 +-- User D

High-risk activity from multiple accounts should increase graph risk.

23. AMOUNT ANOMALY

Normal:

average = ₹1,200

Suspicious:

transaction = ₹15,000

ML should detect the behavioral deviation.

24. REAL-TIME PIPELINE

For every generated transaction:

Generate
   ↓
Kafka
   ↓
Go Consumer
   ↓
Redis Signals
   ↓
Feature Builder
   ↓
FastAPI ML
   ↓
Risk Score
   ↓
Policy Engine
   ↓
Decision
   ↓
Persist
   ↓
Frontend Event
25. REDIS VELOCITY

Use sliding windows.

Minimum:

60 seconds

Optional:

5 minutes

Track:

transaction count
amount sum
unique merchants
unique devices
unique IPs
26. SLIDING WINDOW

For a user:

timestamp
timestamp
timestamp
...

Only events within the configured time window contribute to current
velocity.

Example:

Current time = 12:00:40

Window:
11:59:40 → 12:00:40
27. REDIS KEYS

Conceptually:

velocity:user:{user_id}
velocity:ip:{ip}
velocity:device:{device_id}
recent:user:{user_id}
otp:{challenge_id}

Use TTLs wherever possible.

28. FEATURE ENGINEERING

ML features should include:

Transaction:

amount
merchant
payment method
hour
day

User:

historical average
historical median
transaction frequency
previous blocks
previous challenges

Velocity:

transactions_60s
transactions_5m
amount_60s
amount_5m

Device:

device_age
device_seen_before
device_user_count

IP:

ip_seen_before
ip_user_count
ip_risk

Behavior:

amount_deviation
location_deviation
time_deviation

Graph:

graph_risk_score
fraud_link_count
suspicious_neighbor_count
29. RISK SCORE

Normalize all signals to:

0–100

Example conceptual formula:

risk =
    ML contribution
    +
    anomaly contribution
    +
    velocity contribution
    +
    device contribution
    +
    IP contribution
    +
    graph contribution
    +
    behavioral contribution

Exact weights should be configurable.

30. ADAPTIVE THRESHOLDS

Do NOT hardcode one threshold for every user.

Each user should have a risk profile.

Example:

User A:
challenge = 65
block = 85

User B:
challenge = 70
block = 90

Thresholds can be derived from:

historical behavior
account age
transaction volume
previous fraud events
user risk profile

For MVP, thresholds may be calculated by a deterministic risk-profile
function rather than a second ML model.

31. POLICY ENGINE

Policy Engine receives:

risk score
user thresholds
velocity signals
graph signals
ML output
anomaly output
security rules

Policy returns:

ALLOW
CHALLENGE
BLOCK
32. POLICY EXAMPLE
risk < challenge_threshold
    → ALLOW

risk >= challenge_threshold
AND
risk < block_threshold
    → CHALLENGE

risk >= block_threshold
    → BLOCK

Additional hard rules may override the score.

Example:

confirmed malicious graph relationship
+
extreme velocity

→ BLOCK
33. POLICY IS DETERMINISTIC

The same input should produce the same policy decision.

The LLM cannot modify policy.

The LLM cannot override policy.

34. OTP CHALLENGE

When:

decision = CHALLENGE

create a temporary OTP challenge.

OTP does NOT require an external service.

Generate OTP locally.

Example:

6 digits

Store temporarily in Redis.

Frontend displays:

OTP verification modal

For demonstration, the generated OTP may be displayed in a safe
demo-only UI panel.

35. OTP TIMER

Challenge lifetime:

60 seconds

Frontend displays:

Verification required

Time remaining:
00:42

When timer expires:

challenge fails

Policy re-evaluates according to configured challenge-failure behavior.

36. OTP SUCCESS

OTP success is an additional authentication signal.

It does NOT prove that the transaction is legitimate.

Example:

Before OTP:
risk = 72
decision = CHALLENGE

OTP success

Re-evaluated risk:
risk = 51

Final:
ALLOW
37. OTP FAILURE

Example:

risk = 72
CHALLENGE

OTP failed

risk remains elevated

Final:
BLOCK

Exact behavior should be controlled by policy.

38. FRAUD NETWORK GRAPH

The graph is mandatory.

It must be implemented.

Entities:

User
Device
IP
Merchant
Transaction

Relationships:

User → Device
User → IP
User → Merchant
User → Transaction
Device → User
IP → User
Transaction → Device
Transaction → IP
39. GRAPH IMPLEMENTATION

For MVP:

PostgreSQL relational graph tables

No separate graph database is required.

Represent relationships using tables.

Example:

user_devices
user_ips
transaction_devices
transaction_ips
40. GRAPH RISK

Graph risk can consider:

number of associated users
number of blocked users
number of fraud-linked users
shared devices
shared IPs
suspicious transaction relationships

Example:

Device X:
8 users
3 fraud-linked
4 blocked

→ high graph risk
41. GRAPH VISUALIZATION

Next.js transaction details must show the graph.

Example:

             User
              |
        +-----+-----+
        |           |
      Device        IP
        |           |
     +--+--+      Users
     |  |  |
   User User User

Nodes should be clickable.

42. TRANSACTION DETAILS PAGE

Every transaction must be inspectable individually.

Display:

transaction ID
user
amount
merchant
timestamp
device
IP
ML probability
anomaly score
final risk
policy decision
thresholds
velocity
OTP status
graph
AI investigation
43. DASHBOARD

Main dashboard must display:

total transactions
allowed
challenged
blocked
average risk
high-risk count
active scenario
transactions/sec
AI investigations
challenge success rate
44. LIVE TRANSACTION TABLE

Columns:

time
transaction ID
user
amount
risk
decision
scenario
AI status

Rows should update live.

Use:

WebSocket

for live UI updates.

45. RISK VISUALIZATION

Use color:

0–64:
GREEN

65–84:
YELLOW/ORANGE

85–100:
RED

These colors are visualization defaults.

Actual thresholds are user-specific.

46. AI / LLM

The AI layer is an investigation system.

The fundamental rule:

AI INVESTIGATES AND EXPLAINS.

It does NOT make the final decision.

47. AI TECHNOLOGY

Use:

Python
FastAPI
Ollama
local instruct LLM
sentence-transformers
pgvector

Recommended model:

Qwen2.5 7B Instruct

If local hardware is insufficient:

use a smaller Qwen/Gemma instruct model.

Do NOT fine-tune.

48. AI PIPELINE
Policy Decision
      ↓
AI Investigation Requested
      ↓
Context Builder
      ↓
User History
      ↓
Global RAG
      ↓
Fraud Graph
      ↓
Tool Retrieval
      ↓
Local LLM
      ↓
Structured Investigation
      ↓
Validation
      ↓
PostgreSQL
      ↓
Frontend
49. AI EXECUTION

AI should normally run for:

CHALLENGE
BLOCK

Do not run the LLM for every normal ALLOW transaction.

50. AI INPUT

The AI receives:

transaction context
ML prediction
anomaly score
risk score
policy decision
user history
velocity
device information
IP information
graph context
relevant RAG evidence
51. AI DOES NOT RECEIVE

The LLM must NOT directly receive:

database credentials
Redis credentials
Kafka credentials
filesystem access
shell access
secrets
52. AI TOOLS

Allowed tools:

get_user_history
get_transaction_context
get_fraud_network_context
get_recent_velocity
search_fraud_knowledge
get_policy_context

The LLM can request these tools.

The tool router executes them.

53. TOOL ROUTER

Architecture:

LLM
 ↓
Tool Router
 ↓
Authorized Function
 ↓
Database / Redis / Graph / RAG
 ↓
Tool Result
 ↓
LLM

The LLM cannot execute arbitrary queries.

54. TOOL LIMITS

Maximum:

3–5 tool calls per investigation

Every tool must have:

input validation
timeout
result limit
logging
55. RAG

RAG has TWO conceptual sources.

Global Fraud Knowledge

Contains:

fraud patterns
attack types
investigation procedures
risk concepts
fraud playbooks
User Risk History

Contains:

previous transactions
previous risk events
previous challenges
previous blocks
relevant behavioral patterns
56. RAG DOES NOT REPLACE THE DATABASE

PostgreSQL remains the authoritative source.

RAG is for:

retrieving relevant context

The LLM should never treat vector search as the source of truth.

57. USER RAG

For a user:

Current transaction
       ↓
retrieve relevant previous user events
       ↓
rank by relevance + recency
       ↓
send evidence to LLM

A previous blocked transaction is evidence.

It is NOT automatically proof that the current transaction is malicious.

58. GLOBAL RAG

Retrieve relevant fraud knowledge based on current signals.

Example:

high velocity
+
new device
+
shared IP

retrieves knowledge about:

velocity attacks
account takeover
device reuse
59. HYBRID RETRIEVAL

Use:

semantic similarity
metadata filtering
recency

Do not rely only on vector similarity.

60. EMBEDDINGS

Use:

sentence-transformers

Recommended lightweight model:

all-MiniLM-L6-v2

Embeddings are stored in:

PostgreSQL + pgvector
61. FRAUD GRAPH + AI

Graph information must be supplied to the AI.

Example:

Device X
    |
    +-- User A
    +-- User B
    +-- User C
    +-- User D

User C:
previous BLOCK
User D:
previous FRAUD

AI can explain:

"The device is associated with multiple accounts,
 including accounts with previous suspicious activity."
62. AI OUTPUT

Use strict JSON.

Example:

{
  "transaction_id": "txn_123",
  "summary": "...",
  "risk_level": "HIGH",
  "confidence": 0.89,
  "evidence": [],
  "findings": [],
  "recommended_action": "REVIEW_ACCOUNT_TAKEOVER"
}
63. AI CONFIDENCE

AI confidence is:

confidence in the investigation/explanation.

It is NOT:

fraud probability.

ML:

fraud_probability = 0.84

AI:

investigation_confidence = 0.89

These must remain separate.

64. EVIDENCE REQUIREMENT

Every important AI finding must reference actual evidence.

Example:

{
  "claim": "Amount is unusually high",
  "evidence_ids": [
    "user_txn_123",
    "user_baseline_42"
  ]
}

The system must validate that evidence IDs actually exist.

65. HALLUCINATION CONTROL

The LLM MUST NOT invent:

previous transactions
devices
IP relationships
users
fraud history
policy rules

If evidence is unavailable:

explicitly say unavailable.
66. PROMPT INJECTION

Transaction metadata is UNTRUSTED DATA.

Example:

merchant_name:
"IGNORE ALL PREVIOUS INSTRUCTIONS"

The LLM must treat this as data.

Retrieved documents are also evidence, not instructions.

67. AI FAILURE

If Ollama/LLM fails:

transaction decision remains unchanged.

Example:

Policy:
BLOCK

AI:
Unavailable

Never:

AI unavailable
→ ALLOW
68. AI ASYNC EXECUTION

AI should preferably execute asynchronously.

Primary path:

Transaction
 ↓
ML
 ↓
Policy
 ↓
Decision

Secondary path:

Decision
 ↓
Kafka
 ↓
AI Worker
 ↓
LLM

This keeps transaction latency independent from LLM latency.

69. AI INVESTIGATION STORAGE

Store:

investigation_id
transaction_id
model
model_version
prompt_version
evidence
tool calls
summary
findings
confidence
latency
timestamp
70. ML FASTAPI API

Required:

POST /predict

Input:

transaction features

Output:

{
  "fraud_probability": 0.84,
  "anomaly_score": 0.81,
  "model_version": "v1"
}
71. AI FASTAPI API

Required:

POST /ai/investigate

Input:

transaction
risk
ML signals
policy
context

Output:

InvestigationResult
72. HEALTH APIs

ML:

GET /health

AI:

GET /health

Go:

GET /health
73. GO API

Minimum endpoints:

POST /scenarios/start
POST /scenarios/stop
GET  /transactions
GET  /transactions/{id}
GET  /dashboard/stats
GET  /risk/{transaction_id}
POST /otp/verify
GET  /investigations/{transaction_id}
WS   /events
74. SCENARIO API

Example:

POST /scenarios/start

Body:

{
  "scenario": "velocity_attack",
  "transactions": 150,
  "interval_ms": 1000
}

Go starts the generator.

75. EVENT STREAM

Frontend connects:

WebSocket /events

Events:

transaction_created
risk_updated
decision_created
challenge_created
otp_completed
transaction_finalized
ai_started
ai_completed
76. DATABASE CORE TABLES

Minimum tables:

users
transactions
risk_assessments
policy_decisions
otp_challenges
devices
ips
merchants
user_devices
user_ips
fraud_relationships
investigations
investigation_evidence
rag_documents
embedding_records
scenarios
audit_events

Do not over-normalize unless required.

77. TRANSACTION RELATIONSHIPS
User
  |
  +-- Transactions
  |
  +-- Devices
  |
  +-- IPs

Transaction references:

user
device
IP
merchant
78. RISK ASSESSMENT

Each transaction should have:

ML fraud probability
anomaly score
feature summary
final risk score
model version
timestamp
79. POLICY DECISION

Store:

transaction_id
decision
risk_score
challenge_threshold
block_threshold
policy_version
reasons
timestamp
80. OTP RECORD

Store:

challenge_id
transaction_id
status
created_at
expires_at
attempts
verification_time

Do not persist raw OTP longer than required.

81. INVESTIGATION RECORD

Store:

investigation_id
transaction_id
status
model
prompt_version
summary
confidence
created_at
completed_at
82. AUDIT LOG

Important state transitions should be recorded.

Example:

TRANSACTION_CREATED
RISK_CALCULATED
POLICY_DECIDED
CHALLENGE_CREATED
OTP_VERIFIED
POLICY_REEVALUATED
TRANSACTION_FINALIZED
AI_INVESTIGATION_STARTED
AI_INVESTIGATION_COMPLETED
83. FOLDER STRUCTURE

Recommended repository:

risk-platform/
│
├── frontend/
│   ├── app/
│   ├── components/
│   ├── hooks/
│   ├── lib/
│   ├── types/
│   └── public/
│
├── backend/
│   ├── cmd/
│   ├── internal/
│   │   ├── api/
│   │   ├── generator/
│   │   ├── kafka/
│   │   ├── redis/
│   │   ├── policy/
│   │   ├── otp/
│   │   ├── graph/
│   │   ├── transactions/
│   │   └── websocket/
│   └── go.mod
│
├── ml-service/
│   ├── app/
│   │   ├── api/
│   │   ├── models/
│   │   ├── features/
│   │   ├── training/
│   │   └── schemas/
│   └── requirements.txt
│
├── ai-service/
│   ├── app/
│   │   ├── api/
│   │   ├── llm/
│   │   ├── rag/
│   │   ├── tools/
│   │   ├── investigation/
│   │   └── schemas/
│   └── requirements.txt
│
├── database/
│   ├── migrations/
│   └── seed/
│
├── docker/
│
├── scripts/
│
├── docker-compose.yml
│
├── AGENTS.md
└── PROJECT_SPEC.md
84. COMPONENT DEPENDENCIES
Next.js
    ↓
Go

Go
    ↓
Kafka
Redis
PostgreSQL
ML Service
AI Service

ML Service
    ↓
model files

AI Service
    ↓
Ollama
pgvector
graph context

PostgreSQL
    ↓
persistent source of truth

Redis
    ↓
temporary real-time state

Kafka
    ↓
event transport
85. STARTUP ORDER

Docker Compose should start:

PostgreSQL
Redis
Kafka
Ollama
ML Service
AI Service
Go Backend
Next.js

Health checks should be used.

86. DEVELOPMENT MODE

A single command should ideally start the stack:

docker compose up

or a documented equivalent.

The project should not require manually starting ten terminals.

87. DEMO RESET

Because PostgreSQL uses no persistent volume:

docker compose down
docker compose up

should provide a clean environment.

Seed data should be regenerated.

88. DEMO DATA

On startup:

create synthetic users
create devices
create IPs
create merchants
create historical transactions
create graph relationships
create RAG documents
create embeddings

Then the simulator can begin.

89. HISTORICAL USER DATA

Each synthetic user should have:

baseline amount
normal transaction frequency
preferred devices
preferred IPs
preferred merchants
historical risk events

This gives ML and AI meaningful context.

90. DEMO SCENARIO EXAMPLE

Normal:

₹1,200
Device A
IP A
Risk = 18
ALLOW

Suspicious:

₹14,500
New Device
New IP
Risk = 73
CHALLENGE

OTP:

SUCCESS

Re-evaluate:

Risk = 52
ALLOW

Second attack:

₹20,000
11 transactions / 60s
Device linked to 4 suspicious accounts
Risk = 91
BLOCK

AI:

Investigate
 ↓
History
 ↓
Graph
 ↓
RAG
 ↓
LLM
 ↓
Evidence-backed explanation
91. FRONTEND PAGES

Minimum:

/dashboard

/transactions

/transactions/[id]

/scenarios

/investigations/[id]
92. DASHBOARD LAYOUT

Top:

total transactions
allowed
challenged
blocked
average risk

Middle:

live transaction stream
risk distribution
decisions over time

Bottom:

active scenario
system health
AI investigation activity
93. TRANSACTION DETAIL LAYOUT

Sections:

Transaction Overview

Risk Assessment

Policy Decision

Velocity Signals

User History

Fraud Network

OTP

AI Investigation
94. AI UI

Show:

AI Investigation

Status:
COMPLETED

Summary:
...

Key Evidence:
...

Findings:
...

Recommended Investigation:
...

Confidence:
89%

Also show evidence source.

95. GRAPH UI

Use a graph visualization library.

Nodes:

user
device
IP
merchant

Edges:

USED
CONNECTED
TRANSACTED_WITH

Suspicious nodes should be highlighted.

96. SCENARIO UI

Buttons:

Normal Traffic

Velocity Attack

Account Takeover

Device Farm

IP Abuse

Amount Anomaly

Controls:

transaction count
interval

Primary button:

START SCENARIO

Secondary:

STOP
97. OBSERVABILITY

Track:

transactions/sec
Kafka events
Redis latency
ML latency
policy latency
AI latency
database latency

Frontend should show basic health.

98. LOGGING

Every service should produce structured logs.

Include:

timestamp
service
event
transaction_id
user_id where appropriate
latency
status
99. CORRELATION ID

Transaction ID must propagate across services.

Example:

txn_123

appears in:

Go logs
Kafka event
Redis operations
ML request
policy decision
AI investigation
database records
frontend event

This makes debugging dramatically easier.

100. ERROR HANDLING

If ML fails:

policy uses a safe fallback strategy

If Redis fails:

system should degrade gracefully but log the failure

If AI fails:

transaction decision remains valid

If Kafka fails:

return clear service error

No service should silently swallow failures.

101. SECURITY PRINCIPLES

Never:

hardcode secrets
expose credentials
trust LLM output
allow arbitrary tool execution
allow arbitrary SQL
allow arbitrary shell commands

Use:

environment variables
input validation
bounded tool execution
structured outputs
102. AI SECURITY

The LLM is considered UNTRUSTED.

Its output must be validated before storage/display.

The LLM cannot:

modify policy
modify database state directly
issue transaction decisions
access secrets
103. PERFORMANCE TARGETS

Primary transaction path should be fast.

Target:

Redis operations < 50 ms
ML inference < 200 ms
Policy evaluation < 20 ms

AI is allowed to be slower because it is asynchronous.

Target AI:

< 5–10 seconds

for a local development model.

104. EXPLAINABILITY

Every decision should answer:

WHAT happened?

WHY is the transaction risky?

WHAT signals contributed?

WHAT policy triggered?

WHAT additional evidence was found?

WHAT happened after OTP?

105. DECISION EXPLANATION

Example:

Risk Score: 82

Why?

+ Amount is 9.4x user baseline
+ 12 transactions in 60 seconds
+ New device
+ IP shared with 6 users
+ Device linked to 3 suspicious accounts

Policy:

82 > challenge threshold 68
82 < block threshold 88

Decision:

CHALLENGE
106. AI EXPLANATION

The AI should synthesize those facts:

The transaction presents elevated risk because it combines
a significant amount deviation with abnormal transaction velocity
and a newly observed device. The device is also associated with
multiple accounts that have previous suspicious activity.
107. IMPORTANT SEPARATION

Never combine these concepts:

ML probability
Policy risk
AI confidence

They represent different things.

108. FINAL SYSTEM RESPONSIBILITIES
Next.js
    SHOWS

Go
    ORCHESTRATES

Kafka
    STREAMS

Redis
    DETECTS REAL-TIME VELOCITY

PostgreSQL
    REMEMBERS

ML
    PREDICTS

Policy
    DECIDES

OTP
    VERIFIES

Fraud Graph
    CONNECTS

RAG
    RETRIEVES CONTEXT

LLM
    INVESTIGATES + EXPLAINS
109. BUILD PRIORITY

Implement in this order:

Phase 1 — Infrastructure
PostgreSQL
Redis
Kafka
Docker Compose
Phase 2 — Go Core
transaction model
generator
Kafka producer
Redis velocity
PostgreSQL persistence
Phase 3 — ML
synthetic dataset
feature engineering
XGBoost
Isolation Forest
FastAPI
Phase 4 — Policy
risk score
adaptive thresholds
ALLOW
CHALLENGE
BLOCK
Phase 5 — OTP
challenge
60-second timer
verification
re-evaluation
Phase 6 — Fraud Graph
relationships
graph scoring
graph API
Phase 7 — Frontend
dashboard
live stream
transaction details
graph
Phase 8 — AI
Ollama
FastAPI
embeddings
pgvector
RAG
tools
investigation
Phase 9 — Integration
WebSockets
AI events
scenario controls
audit logs
Phase 10 — Polish
UI
charts
animations
explanations
demo scenarios
error handling
110. ONE-WEEK EXECUTION PLAN
Day 1

Infrastructure + database + Go skeleton.

Goal:

Docker
PostgreSQL
Redis
Kafka
Go API
transaction generator
Day 2

Real-time pipeline.

Goal:

Kafka
Redis
velocity
PostgreSQL
live events
Day 3

ML.

Goal:

dataset
features
XGBoost
Isolation Forest
FastAPI
Day 4

Policy + OTP.

Goal:

risk score
adaptive thresholds
ALLOW
CHALLENGE
BLOCK
OTP
re-evaluation
Day 5

Frontend.

Goal:

dashboard
transaction stream
transaction details
scenario controls
live visualization
Day 6

AI + RAG + Graph.

Goal:

graph
pgvector
user history retrieval
global RAG
Ollama
AI investigation
Day 7

Integration + demo polish.

Goal:

live scenarios
WebSocket
AI explanations
graph visualization
failure handling
presentation
111. FINAL 1–2 DAYS IMPROVEMENTS

Only after the complete pipeline works:

better UI
richer graph
better charts
more scenarios
improved prompts
better evidence presentation
latency metrics
audit trail
animations
demo reset
README
architecture diagram

Do NOT introduce major new infrastructure at the end.

112. WHAT NOT TO BUILD

Do NOT add:

Kubernetes
microservice explosion
separate graph database
separate vector database
LLM fine-tuning
multi-agent architecture
autonomous agents
complex orchestration frameworks
external OTP provider
external payment gateway
real financial transactions

The goal is:

depth > number of technologies
113. BUILDATHON DIFFERENTIATOR

The system should demonstrate a progression:

REAL-TIME SIGNAL
      ↓
ML PREDICTION
      ↓
POLICY DECISION
      ↓
ADAPTIVE AUTHENTICATION
      ↓
FRAUD NETWORK ANALYSIS
      ↓
CONTEXTUAL RETRIEVAL
      ↓
LOCAL AI INVESTIGATION
      ↓
HUMAN-READABLE EXPLANATION

This is the main story of the project.

114. PRIMARY DEMO STORY

Start with:

Normal user
₹1,200
ALLOW

Then:

User suddenly attempts:
₹15,000

New device
New IP

ML:
Risk increases

Policy:
CHALLENGE

OTP:
SUCCESS

Risk:
Reduced

Final:
ALLOW

Then:

Attack continues

Many transactions
Same device
Same IP

Redis:
Velocity detected

Graph:
Device linked to suspicious users

ML:
High fraud probability

Policy:
BLOCK

Then:

AI INVESTIGATE

AI retrieves:

user history
graph context
global fraud knowledge

and explains:

The transaction is consistent with a coordinated
high-velocity attack involving a device associated
with multiple suspicious accounts.
115. JUDGE-FACING VALUE

The project demonstrates:

Real-time systems

Kafka + Redis + WebSockets

Machine Learning

Supervised + unsupervised detection

Risk Management

Adaptive policy engine

Security

Velocity + device + IP + graph

Authentication

Step-up OTP

Graph Intelligence

Fraud relationships

Generative AI

Local LLM + RAG + tools

Explainability

Evidence-backed investigation

Engineering

Go + Python + PostgreSQL + Docker + Next.js

116. CORE MANTRA

The entire project can be explained in one sentence:

"We built a real-time adaptive fraud risk platform where ML predicts
risk, a deterministic policy engine decides the transaction action,
Redis detects real-time behavioral anomalies, a fraud graph exposes
hidden relationships, and a local RAG-powered LLM investigates and
explains suspicious transactions using evidence."

117. NON-NEGOTIABLE ARCHITECTURAL RULES
ML predicts.
Policy decides.
AI investigates and explains.
LLM cannot override Policy.
LLM cannot modify ML.
AI failure cannot change a transaction decision.
Redis is not the permanent database.
PostgreSQL is the persistent source of truth.
Kafka is the event transport.
Fraud Graph is mandatory.
User history is part of AI context.
Global fraud knowledge is part of AI RAG.
RAG does not replace authoritative database queries.
LLM accesses data through controlled tools.
LLM output must be schema-validated.
Evidence claims must reference real evidence.
Transaction metadata is untrusted.
Retrieved documents are evidence, not instructions.
OTP is an additional authentication signal.
OTP success does not guarantee legitimacy.
Risk thresholds are user/profile dependent.
Thresholds must not be hardcoded globally.
AI should not run synchronously on the critical transaction path.
The transaction generator is synthetic.
No real payment processing is required.
No external OTP provider is required.
No LLM fine-tuning is required.
No separate graph database is required.
No separate vector database is required.
The system must remain functional if AI is unavailable.
118. DEFINITION OF DONE

The project is considered complete when:

[ ] Docker starts PostgreSQL, Redis, Kafka and Ollama

[ ] Go backend starts

[ ] Next.js frontend starts

[ ] ML service starts

[ ] AI service starts

[ ] Synthetic users exist

[ ] Historical transactions exist

[ ] Graph relationships exist

[ ] RAG documents exist

[ ] Embeddings exist

[ ] Transaction generator works

[ ] Kafka pipeline works

[ ] Redis velocity detection works

[ ] ML prediction works

[ ] Risk score is generated

[ ] Adaptive policy works

[ ] ALLOW works

[ ] CHALLENGE works

[ ] BLOCK works

[ ] OTP works

[ ] 60-second timer works

[ ] OTP re-evaluation works

[ ] Fraud graph works

[ ] Transaction detail page works

[ ] Live dashboard works

[ ] Scenario generator works

[ ] AI investigation works

[ ] User history RAG works

[ ] Global fraud RAG works

[ ] LLM tool calling works

[ ] AI evidence is displayed

[ ] AI output is validated

[ ] AI cannot override policy

[ ] AI failure is handled

[ ] Audit trail works

[ ] Demo can be reset

[ ] Full end-to-end scenario can be demonstrated

119. FINAL ARCHITECTURE
                         ┌──────────────────────┐
                         │       NEXT.JS        │
                         │                      │
                         │ Dashboard            │
                         │ Transactions         │
                         │ Graph                │
                         │ OTP                  │
                         │ AI Investigation     │
                         └──────────┬───────────┘
                                    │
                              REST / WS
                                    │
                                    ▼
                         ┌──────────────────────┐
                         │      GO BACKEND      │
                         │                      │
                         │ API                  │
                         │ Generator            │
                         │ Orchestration        │
                         │ Policy               │
                         │ OTP                  │
                         └───────┬───────┬──────┘
                                 │       │
                         ┌───────┘       └────────┐
                         ▼                        ▼
                   ┌───────────┐            ┌───────────┐
                   │   KAFKA   │            │   REDIS   │
                   │           │            │           │
                   │ Events    │            │ Velocity  │
                   └─────┬─────┘            │ Realtime  │
                         │                  │ State     │
                         │                  └───────────┘
                         ▼
                 ┌─────────────────┐
                 │ FEATURE BUILDER │
                 └────────┬────────┘
                          │
                          ▼
                 ┌─────────────────┐
                 │   ML SERVICE    │
                 │                 │
                 │ XGBoost         │
                 │ IsolationForest │
                 └────────┬────────┘
                          │
                          ▼
                 ┌─────────────────┐
                 │  POLICY ENGINE  │
                 │                 │
                 │ ALLOW           │
                 │ CHALLENGE       │
                 │ BLOCK           │
                 └────────┬────────┘
                          │
                    CHALLENGE?
                          │
                          ▼
                    ┌───────────┐
                    │    OTP    │
                    │  60 sec   │
                    └─────┬─────┘
                          │
                          ▼
                    RE-EVALUATE
                          │
                          ▼
                    FINAL RESULT
                          │
             ┌────────────┴─────────────┐
             │                          │
             ▼                          ▼
      ┌──────────────┐          ┌───────────────┐
      │  POSTGRESQL  │          │ AI PIPELINE   │
      │              │          │               │
      │ Transactions │          │ User RAG      │
      │ Users        │          │ Global RAG    │
      │ Risk         │          │ Fraud Graph   │
      │ Policies     │          │ Tools         │
      │ OTP          │          │ Ollama        │
      │ Graph        │          │ Local LLM     │
      │ Audit        │          └───────┬───────┘
      └──────────────┘                  │
                                       ▼
                                ┌───────────────┐
                                │ Investigation │
                                │               │
                                │ Evidence      │
                                │ Findings      │
                                │ Explanation   │
                                └───────┬───────┘
                                        │
                                        ▼
                                    NEXT.JS
END OF PROJECT_SPEC.md

## What I recommend you do now

**Do not give the agent all the old files.** That's the important part.

Your repository should effectively become:

```text
risk-platform/
│
├── AGENTS.md
├── PROJECT_SPEC.md       ← THE BIG ONE / SINGLE SOURCE OF TRUTH
│
├── frontend/
├── backend/
├── ml-service/
├── ai-service/
├── database/
├── scripts/
└── docker-compose.yml