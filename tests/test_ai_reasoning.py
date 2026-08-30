from types import SimpleNamespace
import importlib.util
from pathlib import Path

base = Path(__file__).resolve().parents[1]
module_path = base / 'ai-service' / 'app' / 'ollama.py'
spec = importlib.util.spec_from_file_location('ai_app_ollama', module_path)
module = importlib.util.module_from_spec(spec)
assert spec and spec.loader
spec.loader.exec_module(module)
build_decision_summary = module.build_decision_summary


def test_build_decision_summary_for_blocked_transaction():
    req = SimpleNamespace(
        transaction_id="TX-905",
        user_id="C1009",
        amount=3200,
        risk_score=88,
        fraud_probability=0.89,
        anomaly_score=0.91,
        status="BLOCKED",
        decision="BLOCK",
    )

    summary = build_decision_summary(
        req=req,
        device_profile={"device_id": "D204", "is_emulator": True, "trust_score": 22},
        ip_profile={"ip_address": "IP-17", "is_vpn": True, "risk_score": 91},
        customer_history={"typical_max_amount": 1200},
        related_accounts=[{"fraud_linked": True}],
    )

    assert "blocked" in summary.lower()
    assert "device" in summary.lower() or "ip" in summary.lower()


def test_build_decision_summary_for_otp_verified_transaction():
    req = SimpleNamespace(
        transaction_id="TX-911",
        user_id="C1004",
        amount=1900,
        risk_score=71,
        fraud_probability=0.74,
        anomaly_score=0.77,
        status="APPROVED",
        decision="ALLOW",
    )

    summary = build_decision_summary(
        req=req,
        device_profile={"device_id": "D110", "is_emulator": False, "trust_score": 86},
        ip_profile={"ip_address": "IP-22", "is_vpn": False, "risk_score": 29},
        customer_history={"typical_max_amount": 1500},
        related_accounts=[],
    )

    assert "otp" in summary.lower()
    assert "approved" in summary.lower()
