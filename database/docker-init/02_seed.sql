-- VulcanShield Phase 2 Synthetic Seed Data
-- Seed file: database/seed/seed.sql

-- Clear existing seed data (clean restart)
TRUNCATE TABLE audit_events, scenarios, embedding_records, rag_documents, investigation_evidence, investigations, fraud_relationships, otp_challenges, policy_decisions, risk_assessments, transactions, user_ips, user_devices, merchants, ips, devices, users CASCADE;

-- 1. Seed Users (with authoritative per-user thresholds)
INSERT INTO users (user_id, name, email, risk_tolerance, challenge_threshold, block_threshold, account_age_days, trust_score, typical_min_amount, typical_max_amount)
VALUES 
    ('C1001', 'Alice Smith', 'alice@example.com', 'LOW', 55, 80, 450, 85, 15.00, 250.00),
    ('C1002', 'Bob Jones', 'bob@example.com', 'MEDIUM', 65, 85, 120, 60, 5.00, 100.00),
    ('C1003', 'Charlie Brown', 'charlie@example.com', 'HIGH', 75, 90, 30, 30, 50.00, 1500.00),
    ('C1004', 'Dana Lee', 'dana@example.com', 'LOW', 58, 82, 900, 92, 20.00, 400.00),
    ('C1005', 'Evan Patel', 'evan@example.com', 'MEDIUM', 63, 86, 280, 74, 10.00, 300.00),
    ('C1006', 'Fatima Khan', 'fatima@example.com', 'MEDIUM', 67, 88, 510, 81, 25.00, 600.00),
    ('C1007', 'Gabriel Chen', 'gabriel@example.com', 'LOW', 60, 84, 730, 89, 15.00, 350.00),
    ('C1008', 'Harper Wilson', 'harper@example.com', 'HIGH', 72, 91, 75, 42, 30.00, 900.00),
    ('C1009', 'Ishan Rao', 'ishan@example.com', 'MEDIUM', 66, 87, 190, 68, 8.00, 180.00),
    ('C1010', 'Jordan Kim', 'jordan@example.com', 'LOW', 57, 81, 1080, 95, 40.00, 700.00),
    ('C1011', 'Kai Morgan', 'kai@example.com', 'HIGH', 74, 92, 45, 35, 12.00, 220.00),
    ('C1012', 'Lena Ortiz', 'lena@example.com', 'MEDIUM', 64, 86, 365, 77, 18.00, 500.00)
ON CONFLICT (user_id) DO NOTHING;

-- 2. Seed Devices
INSERT INTO devices (device_id, fingerprint_hash, device_type, os, browser, trust_score, is_emulator)
VALUES 
    ('D204', 'fp_macbook_pro_991823a', 'desktop', 'macOS', 'Chrome', 90, FALSE),
    ('D205', 'fp_iphone_15_88172bc', 'mobile', 'iOS', 'Safari', 85, FALSE),
    ('D206', 'fp_android_emulator_99', 'mobile', 'Android', 'Chrome', 15, TRUE),
    ('D207', 'fp_pixel_8_1042', 'mobile', 'Android', 'Chrome', 82, FALSE),
    ('D208', 'fp_windows_11_5133', 'desktop', 'Windows', 'Edge', 78, FALSE),
    ('D209', 'fp_ipad_air_2291', 'tablet', 'iPadOS', 'Safari', 88, FALSE),
    ('D210', 'fp_linux_workstation_7734', 'desktop', 'Linux', 'Firefox', 72, FALSE),
    ('D211', 'fp_iphone_14_4412', 'mobile', 'iOS', 'Safari', 86, FALSE),
    ('D212', 'fp_galaxy_s24_9921', 'mobile', 'Android', 'Chrome', 79, FALSE),
    ('D213', 'fp_macbook_air_3488', 'desktop', 'macOS', 'Chrome', 91, FALSE),
    ('D214', 'fp_windows_laptop_6459', 'desktop', 'Windows', 'Chrome', 76, FALSE)
ON CONFLICT (device_id) DO NOTHING;

-- 3. Seed IPs
INSERT INTO ips (ip_address, country_code, city, isp, is_vpn, is_tor, is_proxy, risk_score)
VALUES 
    ('IP-17', 'US', 'New York', 'Comcast', FALSE, FALSE, FALSE, 10),
    ('IP-18', 'US', 'San Francisco', 'ATT Internet', FALSE, FALSE, FALSE, 15),
    ('IP-19', 'RU', 'Moscow', 'Unknown VPN Ltd', TRUE, FALSE, TRUE, 85),
    ('IP-20', 'IN', 'Mumbai', 'Jio Fiber', FALSE, FALSE, FALSE, 12),
    ('IP-21', 'GB', 'London', 'BT Group', FALSE, FALSE, FALSE, 18),
    ('IP-22', 'US', 'Austin', 'Spectrum', FALSE, FALSE, FALSE, 14),
    ('IP-23', 'CA', 'Toronto', 'Rogers', FALSE, FALSE, FALSE, 16),
    ('IP-24', 'DE', 'Berlin', 'Vodafone', FALSE, FALSE, FALSE, 20),
    ('IP-25', 'SG', 'Singapore', 'Singtel', FALSE, FALSE, FALSE, 11),
    ('IP-26', 'US', 'Chicago', 'Verizon', FALSE, FALSE, FALSE, 13),
    ('IP-27', 'AU', 'Sydney', 'Telstra', FALSE, FALSE, FALSE, 17),
    ('IP-28', 'BR', 'Sao Paulo', 'Claro', FALSE, FALSE, FALSE, 22)
ON CONFLICT (ip_address) DO NOTHING;

-- 4. Seed Merchants
INSERT INTO merchants (merchant_id, name, mcc, risk_category)
VALUES 
    ('M301', 'Acme Electronics', '5732', 'LOW'),
    ('M302', 'Global Luxury Jewels', '5944', 'HIGH'),
    ('M303', 'QuickPay Digital Giftcards', '5999', 'HIGH')
ON CONFLICT (merchant_id) DO NOTHING;

-- 5. Seed User-Device Associations
INSERT INTO user_devices (user_id, device_id, association_count)
VALUES 
    ('C1001', 'D204', 15),
    ('C1002', 'D205', 8),
    ('C1003', 'D206', 2),
    ('C1004', 'D207', 21), ('C1005', 'D208', 12), ('C1006', 'D209', 19),
    ('C1007', 'D210', 11), ('C1008', 'D211', 4), ('C1009', 'D212', 8),
    ('C1010', 'D213', 25), ('C1011', 'D214', 3), ('C1012', 'D207', 7)
ON CONFLICT (user_id, device_id) DO NOTHING;

-- 6. Seed User-IP Associations
INSERT INTO user_ips (user_id, ip_address, association_count)
VALUES 
    ('C1001', 'IP-17', 20),
    ('C1002', 'IP-18', 12),
    ('C1003', 'IP-19', 3),
    ('C1004', 'IP-20', 28), ('C1005', 'IP-21', 16), ('C1006', 'IP-22', 24),
    ('C1007', 'IP-23', 13), ('C1008', 'IP-24', 5), ('C1009', 'IP-25', 10),
    ('C1010', 'IP-26', 30), ('C1011', 'IP-27', 4), ('C1012', 'IP-28', 15)
ON CONFLICT (user_id, ip_address) DO NOTHING;

-- 7. Seed Fraud Relationships (Relational Graph)
INSERT INTO fraud_relationships (relationship_id, source_type, source_id, target_type, target_id, relationship_type, weight, fraud_linked)
VALUES 
    ('REL-001', 'DEVICE', 'D206', 'USER', 'C1003', 'SHARED_DEVICE', 0.95, TRUE),
    ('REL-002', 'IP', 'IP-19', 'USER', 'C1003', 'SHARED_IP', 0.85, TRUE)
ON CONFLICT (relationship_id) DO NOTHING;

-- 8. Seed RAG Documents (Domain Knowledge Base)
INSERT INTO rag_documents (document_id, title, category, content, metadata)
VALUES 
    ('RAG-DOC-001', 'Velocity Attack Patterns in Carding Operations', 'ATTACK_PATTERN', 'Rapid bursts of low to medium value transactions within 60 seconds originating from a shared device or IP address strongly indicate automated card testing or velocity attacks.', '{"severity": "HIGH", "author": "Fraud Ops"}'::jsonb),
    ('RAG-DOC-002', 'Account Takeover (ATO) Investigation Playbook', 'PLAYBOOK', 'When a user logs in from an unrecognized device/IP and immediately attempts a transaction exceeding 3x their historical average, verify OTP step-up challenge and analyze fraud neighbor connections.', '{"category": "ATO", "version": "1.2"}'::jsonb)
ON CONFLICT (document_id) DO NOTHING;

-- 9. Seed Embedding Records (Dummy 384-dimensional vector for all-MiniLM-L6-v2 compatibility testing)
-- A synthetic 384-element array converted to vector(384)
INSERT INTO embedding_records (embedding_id, document_id, chunk_text, embedding)
VALUES 
    ('EMB-001', 'RAG-DOC-001', 'Rapid bursts of low to medium value transactions within 60 seconds originating from a shared device', (SELECT ('[' || array_to_string(array_fill(0.05::float, ARRAY[384]), ',') || ']')::vector)),
    ('EMB-002', 'RAG-DOC-002', 'When a user logs in from an unrecognized device/IP and immediately attempts a transaction exceeding 3x average', (SELECT ('[' || array_to_string(array_fill(0.10::float, ARRAY[384]), ',') || ']')::vector))
ON CONFLICT (embedding_id) DO NOTHING;
