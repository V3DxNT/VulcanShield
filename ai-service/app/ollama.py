import os
import httpx
import json
from typing import Dict, Any, Optional, List
from app.schemas import InvestigationResponse, EvidenceItem, SimilarCase

OLLAMA_URL = os.getenv("OLLAMA_URL", "http://ollama:11434")


def build_decision_summary(
    req: Any,
    customer_history: Optional[Dict[str, Any]] = None,
    device_profile: Optional[Dict[str, Any]] = None,
    ip_profile: Optional[Dict[str, Any]] = None,
    related_accounts: Optional[List[Any]] = None,
) -> str:
    """Return a deterministic, policy-first explanation for the transaction outcome."""
    customer_history = customer_history or {}
    device_profile = device_profile or {}
    ip_profile = ip_profile or {}
    related_accounts = related_accounts or []

    decision = str(getattr(req, "decision", "") or "").upper()
    status = str(getattr(req, "status", "") or "").upper()
    risk_score = getattr(req, "risk_score", 0) or 0
    amount = getattr(req, "amount", 0.0) or 0.0
    user_id = getattr(req, "user_id", "") or ""
    amount_threshold = float(customer_history.get("typical_max_amount", 0.0) or 0.0)

    if not decision:
        if "BLOCK" in status:
            decision = "BLOCK"
        elif "CHALLENGE" in status or "CHALLENGED" in status:
            decision = "CHALLENGE"
        else:
            decision = "ALLOW"

    device_risk = device_profile.get("is_emulator") or (device_profile.get("trust_score", 100) < 40)
    ip_risk = ip_profile.get("is_vpn") or (ip_profile.get("risk_score", 0) >= 60)
    abnormal_amount = amount_threshold > 0 and amount > amount_threshold * 2.0
    linked_graph = any(item.get("fraud_linked") for item in related_accounts)

    if decision == "BLOCK":
        bits = [
            f"• Policy blocked transaction {getattr(req, 'transaction_id', 'UNKNOWN')} for {user_id or 'this customer'}.",
            f"• The recorded risk score was {risk_score}/100. The final authorization decision was made by the policy engine, not by the raw ML score alone.",
        ]
        if device_risk:
            bits.append(f"• The device signal is suspicious: {device_profile.get('device_id', 'device')} looks untrusted or emulator-like.")
        elif ip_risk:
            bits.append(f"• The IP signal is suspicious: {ip_profile.get('ip_address', 'IP')} shows elevated fraud risk.")
        if abnormal_amount:
            bits.append(f"• The amount of ${amount:.2f} materially exceeds the customer baseline of ${amount_threshold:.2f}.")
        if linked_graph:
            bits.append("• Network links show related fraud behavior, which adds context to the current investigation.")
        if customer_history.get("previous_fraud_count", 0):
            bits.append(f"• This customer has {customer_history.get('previous_fraud_count')} prior fraud-related transaction(s), which strengthens the risk context.")
        bits.append("• Final outcome: BLOCK because the transaction exceeded the configured risk boundary and the evidence stack did not support approval.")
        return "\n".join(bits)

    if decision == "CHALLENGE":
        bits = [
            f"• Policy challenged transaction {getattr(req, 'transaction_id', 'UNKNOWN')} for {user_id or 'this customer'}.",
            f"• The risk score reached {risk_score}/100, so the policy engine required step-up verification before approval.",
        ]
        if device_risk:
            bits.append("• Device checks are weaker than the customer’s trusted baseline, which is why the challenge flow was triggered.")
        elif ip_risk:
            bits.append("• The IP profile is elevated enough to justify identity verification before the payment is allowed.")
        if abnormal_amount:
            bits.append(f"• The amount is unusually high for the customer, which increases the risk of a takeover or mule pattern.")
        if customer_history.get("last_transaction_status") in {"BLOCKED", "CANCELLED"}:
            bits.append("• The customer’s last transaction was already marked as risky or blocked, which is relevant user-specific context for this challenge decision.")
        bits.append("• Final outcome: CHALLENGE. If the OTP is verified, the decision can move to ALLOW; if it fails or expires, it remains blocked.")
        return "\n".join(bits)

    if status == "APPROVED" or "APPROVED" in status or "VERIFIED" in status:
        bits = [
            f"• Policy approved transaction {getattr(req, 'transaction_id', 'UNKNOWN')} after the customer completed the required verification flow.",
            f"• The transaction stayed within the policy posture despite earlier risk markers; the successful OTP outcome cleared the challenge state.",
        ]
        if device_risk:
            bits.append("• The device checks were reviewed, and the successful verification step was enough to reduce the risk to an acceptable level.")
        if ip_risk:
            bits.append("• The IP profile remains monitored, but the approved challenge outcome indicates the customer action was valid.")
        bits.append("• Final outcome: ALLOW because policy re-evaluated the transaction with the verified challenge result and accepted the step-up proof.")
        return "\n".join(bits)

    bits = [
        f"• Policy approved transaction {getattr(req, 'transaction_id', 'UNKNOWN')} for {user_id or 'this customer'}.",
        f"• The recorded risk score was {risk_score}/100, and the customer profile remained within the configured risk envelope for this transaction.",
    ]
    if device_risk:
        bits.append("• The device signal was reviewed but did not exceed the threshold that would require a block.")
    elif ip_risk:
        bits.append("• The IP signal was elevated but still below the policy threshold for a block.")
    if abnormal_amount:
        bits.append("• The amount stayed above the customer’s typical range, but it was still within the permitted policy posture.")
    if customer_history.get("previous_fraud_count", 0):
        bits.append(f"• Prior history for this customer includes {customer_history.get('previous_fraud_count')} prior risky event(s), but the current profile still qualified for approval under policy.")
    bits.append("• Final outcome: ALLOW because the deterministic policy engine treated these factors as acceptable for this customer profile.")
    return "\n".join(bits)


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
                fact=f"User has {fraud_neighbors} fraud-marked graph relationship(s); this is contextual evidence, not proof of current fraud.",
                severity="HIGH" if req.decision != "ALLOW" else "LOW"
            ))

    if not evidence:
        evidence.append(EvidenceItem(
            category="NORMAL_BEHAVIOR",
            fact="Transaction parameters align with historical customer behavior and trusted device/IP profile",
            severity="LOW"
        ))

    # Policy is authoritative. The investigator explains the final decision and any
    # step-up verification result rather than substituting a raw ML score.
    decision = (req.decision or "ALLOW").upper()
    if not decision:
        decision = "ALLOW"
    if decision == "BLOCK":
        risk_level = "CRITICAL"
        recommended_action = "BLOCK"
    elif decision == "CHALLENGE":
        risk_level = "HIGH"
        recommended_action = "CHALLENGE"
    else:
        risk_level = "LOW"
        recommended_action = "ALLOW"
    summary = build_decision_summary(
        req=req,
        customer_history=customer_history,
        device_profile=device_profile,
        ip_profile=ip_profile,
        related_accounts=related_accounts,
    )

    cases = [
        SimilarCase(
            case_id=c["case_id"],
            title=c["title"],
            relevance_score=c["relevance_score"]
        )
        for c in similar_cases[:2]
    ]

    inv_id = f"INV-{req.transaction_id}"
    llm_model = "rule-based-evidence-fallback"
    
    # Try calling Ollama LLM endpoint
    try:
        prompt = (
            f"You are a Fraud Analyst AI. Analyze transaction {req.transaction_id} for user {req.user_id}.\n"
            f"Amount: ${req.amount}, Risk Score: {req.risk_score}/100, Fraud Probability: {req.fraud_probability}.\n"
            f"Authoritative policy decision: {decision}. Transaction status: {req.status}.\n"
            f"Evidence: {[e.fact for e in evidence]}\n"
            "Return a 2-sentence formal investigation summary. Explain the supplied policy decision; do not recommend a different action or invent facts."
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
                    llm_model = "qwen2.5:7b-instruct"
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
        llm_model=llm_model,
    )
