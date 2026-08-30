import os
import httpx
import json
from typing import Dict, Any
from app.schemas import InvestigationResponse, EvidenceItem, SimilarCase

OLLAMA_URL = os.getenv("OLLAMA_URL", "http://ollama:11434")

async def generate_llm_investigation(
    req: Any,
    customer_history: Dict[str, Any],
    device_profile: Dict[str, Any],
    ip_profile: Dict[str, Any],
    related_accounts: list,
    similar_cases: list,
) -> InvestigationResponse:
    
    # 1. Build Evidence List from Structured Tools
    evidence = []
    
    if device_profile.get("is_emulator"):
        evidence.append(EvidenceItem(
            category="DEVICE_INTELLIGENCE",
            fact=f"Device {device_profile.get('device_id')} is flagged as an emulator (trust score: {device_profile.get('trust_score')})",
            severity="HIGH"
        ))
        
    if ip_profile.get("is_vpn") or ip_profile.get("risk_score", 0) > 50:
        evidence.append(EvidenceItem(
            category="IP_INTELLIGENCE",
            fact=f"IP {ip_profile.get('ip_address')} is a high-risk proxy/VPN (risk score: {ip_profile.get('risk_score')})",
            severity="HIGH"
        ))

    typical_max = customer_history.get("typical_max_amount", 250.0)
    if req.amount > typical_max * 2.0:
        evidence.append(EvidenceItem(
            category="BEHAVIORAL_ANOMALY",
            fact=f"Transaction amount ${req.amount:.2f} significantly exceeds historical max of ${typical_max:.2f}",
            severity="HIGH"
        ))

    if related_accounts:
        fraud_neighbors = sum(1 for r in related_accounts if r.get("fraud_linked"))
        if fraud_neighbors > 0:
            evidence.append(EvidenceItem(
                category="FRAUD_GRAPH",
                fact=f"User is linked to {fraud_neighbors} confirmed fraud entity node(s)",
                severity="CRITICAL"
            ))

    if not evidence:
        evidence.append(EvidenceItem(
            category="NORMAL_BEHAVIOR",
            fact="Transaction parameters align with historical customer behavior and trusted device/IP profile",
            severity="LOW"
        ))

    # Determine Risk Level & Action
    if req.risk_score >= 80 or any(e.severity == "CRITICAL" for e in evidence):
        risk_level = "CRITICAL"
        recommended_action = "BLOCK_ACCOUNT"
        summary = f"High-risk anomaly detected. Transaction of ${req.amount:.2f} generated a risk score of {req.risk_score}/100 with multiple high-severity signals."
    elif req.risk_score >= 60:
        risk_level = "HIGH"
        recommended_action = "MANUAL_REVIEW"
        summary = f"Suspicious transaction pattern. Risk score {req.risk_score}/100 exceeds challenge threshold."
    else:
        risk_level = "LOW"
        recommended_action = "ALLOW"
        summary = f"Low risk transaction (${req.amount:.2f}). Parameters match baseline customer profile."

    cases = [
        SimilarCase(
            case_id=c["case_id"],
            title=c["title"],
            relevance_score=c["relevance_score"]
        )
        for c in similar_cases[:2]
    ]

    inv_id = f"INV-{req.transaction_id}"
    
    # Try calling Ollama LLM endpoint
    try:
        prompt = (
            f"You are a Fraud Analyst AI. Analyze transaction {req.transaction_id} for user {req.user_id}.\n"
            f"Amount: ${req.amount}, Risk Score: {req.risk_score}/100, Fraud Probability: {req.fraud_probability}.\n"
            f"Evidence: {[e.fact for e in evidence]}\n"
            f"Return a 2-sentence formal investigation summary."
        )
        async with httpx.AsyncClient(timeout=3.0) as client:
            res = await client.post(
                f"{OLLAMA_URL}/api/generate",
                json={"model": "qwen2.5:7b-instruct", "prompt": prompt, "stream": False}
            )
            if res.status_code == 200:
                llm_out = res.json().get("response", "").strip()
                if llm_out:
                    summary = llm_out
    except Exception:
        pass # Degrade to rule-based summary if Ollama offline

    return InvestigationResponse(
        investigation_id=inv_id,
        transaction_id=req.transaction_id,
        risk_level=risk_level,
        summary=summary,
        evidence=evidence,
        similar_cases=cases,
        recommended_action=recommended_action,
        confidence=0.91,
        llm_model="qwen2.5:7b-instruct",
    )
