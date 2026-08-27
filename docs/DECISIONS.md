# VulcanShield Architecture Decisions

## Decision 001 — Backend

Go is the primary transaction/risk orchestration backend.

Reason:
High concurrency, low latency, strong typing and good fit for
real-time transaction processing.

---

## Decision 002 — Real-Time State

Redis is used for short-lived real-time risk signals.

Primary use:
Sliding-window velocity detection.

---

## Decision 003 — Event Streaming

Kafka is used for asynchronous event processing.

---

## Decision 004 — ML

ML predicts risk.

Policy Engine makes the final deterministic decision.

---

## Decision 005 — AI

AI investigates and explains suspicious transactions.

AI must not directly authorize or reject transactions.

---

## Decision 006 — LLM

The LLM runs locally rather than relying on an external inference API.

---

## Decision 007 — Database

PostgreSQL is the primary persistent database.

pgvector is used for semantic retrieval.

---

## Decision 008 — Project Strategy

Working end-to-end functionality has priority over additional
infrastructure or complexity.