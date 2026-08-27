-- VulcanShield Phase 2 Schema Initialization Migration
-- Up Migration: 000001_init_schema.up.sql

-- Enable pgvector extension for vector embeddings
CREATE EXTENSION IF NOT EXISTS vector;

-- 1. users: Customer profiles, baselines, and authoritative policy thresholds
CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    risk_tolerance VARCHAR(32) NOT NULL DEFAULT 'MEDIUM', -- Contextual metadata only
    challenge_threshold INT NOT NULL DEFAULT 60 CHECK (challenge_threshold BETWEEN 0 AND 100),
    block_threshold INT NOT NULL DEFAULT 85 CHECK (block_threshold BETWEEN 0 AND 100),
    account_age_days INT NOT NULL DEFAULT 0 CHECK (account_age_days >= 0),
    trust_score INT NOT NULL DEFAULT 50 CHECK (trust_score BETWEEN 0 AND 100),
    typical_min_amount NUMERIC(12, 2) NOT NULL DEFAULT 10.00 CHECK (typical_min_amount >= 0),
    typical_max_amount NUMERIC(12, 2) NOT NULL DEFAULT 500.00 CHECK (typical_max_amount >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_user_thresholds CHECK (challenge_threshold < block_threshold)
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 2. devices: Device intelligence registry
CREATE TABLE IF NOT EXISTS devices (
    device_id VARCHAR(64) PRIMARY KEY,
    fingerprint_hash VARCHAR(128) NOT NULL UNIQUE,
    device_type VARCHAR(32) NOT NULL DEFAULT 'web',
    os VARCHAR(64) NOT NULL DEFAULT 'unknown',
    browser VARCHAR(64) NOT NULL DEFAULT 'unknown',
    trust_score INT NOT NULL DEFAULT 50 CHECK (trust_score BETWEEN 0 AND 100),
    is_emulator BOOLEAN NOT NULL DEFAULT FALSE,
    first_seen_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_devices_fingerprint ON devices(fingerprint_hash);

-- 3. ips: IP address intelligence registry
CREATE TABLE IF NOT EXISTS ips (
    ip_address VARCHAR(45) PRIMARY KEY, -- IPv4 / IPv6 or synthetic IP identifier
    country_code VARCHAR(8) NOT NULL DEFAULT 'US',
    city VARCHAR(128) NOT NULL DEFAULT 'Unknown',
    isp VARCHAR(128) NOT NULL DEFAULT 'Unknown',
    is_vpn BOOLEAN NOT NULL DEFAULT FALSE,
    is_tor BOOLEAN NOT NULL DEFAULT FALSE,
    is_proxy BOOLEAN NOT NULL DEFAULT FALSE,
    risk_score INT NOT NULL DEFAULT 0 CHECK (risk_score BETWEEN 0 AND 100),
    first_seen_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ips_country ON ips(country_code);

-- 4. merchants: Merchant registry
CREATE TABLE IF NOT EXISTS merchants (
    merchant_id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mcc VARCHAR(8) NOT NULL, -- Merchant Category Code
    risk_category VARCHAR(32) NOT NULL DEFAULT 'LOW', -- LOW, MEDIUM, HIGH
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 5. user_devices: Relational junction connecting users and devices
CREATE TABLE IF NOT EXISTS user_devices (
    user_id VARCHAR(64) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    device_id VARCHAR(64) NOT NULL REFERENCES devices(device_id) ON DELETE CASCADE,
    association_count INT NOT NULL DEFAULT 1 CHECK (association_count > 0),
    first_used_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, device_id)
);

-- 6. user_ips: Relational junction connecting users and IP addresses
CREATE TABLE IF NOT EXISTS user_ips (
    user_id VARCHAR(64) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    ip_address VARCHAR(45) NOT NULL REFERENCES ips(ip_address) ON DELETE CASCADE,
    association_count INT NOT NULL DEFAULT 1 CHECK (association_count > 0),
    first_used_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, ip_address)
);

-- 7. transactions: Synthetic financial transaction stream and lifecycle status
CREATE TABLE IF NOT EXISTS transactions (
    transaction_id VARCHAR(64) PRIMARY KEY,
    user_id VARCHAR(64) NOT NULL REFERENCES users(user_id),
    device_id VARCHAR(64) NOT NULL REFERENCES devices(device_id),
    ip_address VARCHAR(45) NOT NULL REFERENCES ips(ip_address),
    merchant_id VARCHAR(64) NOT NULL REFERENCES merchants(merchant_id),
    amount NUMERIC(12, 2) NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    channel VARCHAR(32) NOT NULL DEFAULT 'WEB',
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'CHALLENGED', 'BLOCKED', 'CANCELLED')),
    timestamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_transactions_user_time ON transactions(user_id, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_device ON transactions(device_id);
CREATE INDEX IF NOT EXISTS idx_transactions_ip ON transactions(ip_address);
CREATE INDEX IF NOT EXISTS idx_transactions_timestamp ON transactions(timestamp DESC);

-- 8. risk_assessments: ML predictions for XGBoost and Isolation Forest models
CREATE TABLE IF NOT EXISTS risk_assessments (
    assessment_id VARCHAR(64) PRIMARY KEY,
    transaction_id VARCHAR(64) NOT NULL REFERENCES transactions(transaction_id) ON DELETE CASCADE,
    fraud_probability FLOAT NOT NULL CHECK (fraud_probability BETWEEN 0.0 AND 1.0),
    anomaly_score FLOAT NOT NULL CHECK (anomaly_score BETWEEN 0.0 AND 1.0),
    fraud_model_version VARCHAR(32) NOT NULL DEFAULT 'xgboost-v1',
    anomaly_model_version VARCHAR(32) NOT NULL DEFAULT 'isoforest-v1',
    risk_score INT NOT NULL CHECK (risk_score BETWEEN 0 AND 100),
    feature_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_risk_assessments_tx ON risk_assessments(transaction_id);

-- 9. policy_decisions: Authoritative deterministic policy engine decisions
CREATE TABLE IF NOT EXISTS policy_decisions (
    decision_id VARCHAR(64) PRIMARY KEY,
    transaction_id VARCHAR(64) NOT NULL REFERENCES transactions(transaction_id) ON DELETE CASCADE,
    decision VARCHAR(32) NOT NULL CHECK (decision IN ('ALLOW', 'CHALLENGE', 'BLOCK')),
    risk_score INT NOT NULL CHECK (risk_score BETWEEN 0 AND 100),
    challenge_threshold INT NOT NULL CHECK (challenge_threshold BETWEEN 0 AND 100),
    block_threshold INT NOT NULL CHECK (block_threshold BETWEEN 0 AND 100),
    policy_version VARCHAR(32) NOT NULL DEFAULT 'v1.0',
    reasons JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_policy_decisions_tx ON policy_decisions(transaction_id);
CREATE INDEX IF NOT EXISTS idx_policy_decisions_decision ON policy_decisions(decision);

-- 10. otp_challenges: Step-up authentication verification state
CREATE TABLE IF NOT EXISTS otp_challenges (
    challenge_id VARCHAR(64) PRIMARY KEY,
    transaction_id VARCHAR(64) NOT NULL REFERENCES transactions(transaction_id) ON DELETE CASCADE,
    otp_code_hash VARCHAR(128) NOT NULL, -- Plaintext OTP is NEVER stored
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'VERIFIED', 'EXPIRED', 'FAILED')),
    attempts INT NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    max_attempts INT NOT NULL DEFAULT 3 CHECK (max_attempts > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    verified_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_otp_challenges_tx ON otp_challenges(transaction_id);
CREATE INDEX IF NOT EXISTS idx_otp_challenges_status ON otp_challenges(status);

-- 11. fraud_relationships: PostgreSQL relational graph data model
CREATE TABLE IF NOT EXISTS fraud_relationships (
    relationship_id VARCHAR(64) PRIMARY KEY,
    source_type VARCHAR(32) NOT NULL CHECK (source_type IN ('USER', 'DEVICE', 'IP', 'MERCHANT')),
    source_id VARCHAR(64) NOT NULL,
    target_type VARCHAR(32) NOT NULL CHECK (target_type IN ('USER', 'DEVICE', 'IP', 'MERCHANT')),
    target_id VARCHAR(64) NOT NULL,
    relationship_type VARCHAR(64) NOT NULL, -- e.g. SHARED_DEVICE, SHARED_IP, SUSPICIOUS_TRANSFER
    weight FLOAT NOT NULL DEFAULT 1.0,
    fraud_linked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_fraud_rel_source ON fraud_relationships(source_id, source_type);
CREATE INDEX IF NOT EXISTS idx_fraud_rel_target ON fraud_relationships(target_id, target_type);

-- 12. investigations: Async AI investigation reports
CREATE TABLE IF NOT EXISTS investigations (
    investigation_id VARCHAR(64) PRIMARY KEY,
    transaction_id VARCHAR(64) NOT NULL REFERENCES transactions(transaction_id) ON DELETE CASCADE,
    status VARCHAR(32) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'COMPLETED', 'FAILED')),
    llm_model VARCHAR(64) NOT NULL DEFAULT 'qwen2.5:7b-instruct',
    prompt_version VARCHAR(32) NOT NULL DEFAULT 'v1.0',
    summary TEXT,
    risk_level VARCHAR(32) CHECK (risk_level IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL')),
    confidence FLOAT CHECK (confidence BETWEEN 0.0 AND 1.0),
    recommendation VARCHAR(32) CHECK (recommendation IN ('ALLOW', 'CHALLENGE', 'BLOCK', 'ESCALATE')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_investigations_tx ON investigations(transaction_id);

-- 13. investigation_evidence: Structured evidence cited by AI investigation
CREATE TABLE IF NOT EXISTS investigation_evidence (
    evidence_id VARCHAR(64) PRIMARY KEY,
    investigation_id VARCHAR(64) NOT NULL REFERENCES investigations(investigation_id) ON DELETE CASCADE,
    evidence_type VARCHAR(64) NOT NULL, -- e.g. HISTORICAL_TX, VELOCITY_SIGNAL, GRAPH_RELATIONSHIP, RAG_DOCUMENT
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    relevance_score FLOAT NOT NULL DEFAULT 1.0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_evidence_investigation ON investigation_evidence(investigation_id);

-- 14. rag_documents: Knowledge base for semantic retrieval
CREATE TABLE IF NOT EXISTS rag_documents (
    document_id VARCHAR(64) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    category VARCHAR(64) NOT NULL, -- e.g. PLAYBOOK, FRAUD_CASE, POLICY_RULE, ATTACK_PATTERN
    content TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_rag_docs_category ON rag_documents(category);

-- 15. embedding_records: Vector storage for semantic RAG search using 384-dim embeddings
CREATE TABLE IF NOT EXISTS embedding_records (
    embedding_id VARCHAR(64) PRIMARY KEY,
    document_id VARCHAR(64) NOT NULL REFERENCES rag_documents(document_id) ON DELETE CASCADE,
    chunk_text TEXT NOT NULL,
    embedding vector(384) NOT NULL, -- Configured for all-MiniLM-L6-v2 384-dimensional vectors
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_embedding_records_doc ON embedding_records(document_id);
CREATE INDEX IF NOT EXISTS idx_embedding_records_vec ON embedding_records USING hnsw (embedding vector_cosine_ops);

-- 16. scenarios: Scenario execution runs
CREATE TABLE IF NOT EXISTS scenarios (
    scenario_id VARCHAR(64) PRIMARY KEY,
    scenario_type VARCHAR(64) NOT NULL, -- e.g. NORMAL, VELOCITY_ATTACK, ACCOUNT_TAKEOVER, DEVICE_REUSE
    transaction_count INT NOT NULL DEFAULT 0,
    seed BIGINT NOT NULL DEFAULT 42,
    status VARCHAR(32) NOT NULL DEFAULT 'RUNNING' CHECK (status IN ('RUNNING', 'COMPLETED', 'STOPPED')),
    started_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMPTZ
);

-- 17. audit_events: State transition audit log for transaction investigation timeline
CREATE TABLE IF NOT EXISTS audit_events (
    event_id BIGSERIAL PRIMARY KEY,
    transaction_id VARCHAR(64) REFERENCES transactions(transaction_id) ON DELETE SET NULL,
    event_type VARCHAR(64) NOT NULL, -- e.g. TRANSACTION_CREATED, RISK_CALCULATED, POLICY_DECIDED, OTP_VERIFIED
    details JSONB NOT NULL DEFAULT '{}'::jsonb,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_audit_tx ON audit_events(transaction_id);
CREATE INDEX IF NOT EXISTS idx_audit_timestamp ON audit_events(timestamp DESC);
