import os
import logging
import httpx
import json
from typing import Dict, Any, Optional, List
from app.schemas import InvestigationResponse, EvidenceItem, SimilarCase

logger = logging.getLogger(__name__)


def _default_ollama_url() -> str:
    """Prefer host-local Ollama by default, but allow Dockerized services to reach the host via host.docker.internal."""
    if os.path.exists("/.dockerenv"):
        return "http://host.docker.internal:11434"
    return "http://localhost:11434"


OLLAMA_URL = os.getenv("OLLAMA_URL") or os.getenv("OLLAMA_HOST") or _default_ollama_url()
GROQ_API_KEY = (os.getenv("GROQ_API_KEY") or "").strip()
GROQ_MODEL = os.getenv("GROQ_MODEL", "llama-3.3-70b-versatile")
DEFAULT_LLM_PROVIDER = (os.getenv("LLM_PROVIDER") or "ollama").lower()


def detect_best_ollama_model(models: Optional[List[str]]) -> str:
    """Choose a working Ollama model from the installed list."""
    preferred = [
        "llama3.1:8b",
        "llama-3.1:8b",
        "llama3.1",
        "llama3.2",
        "qwen2.5:7b-instruct",
        "qwen2.5:latest",
        "qwen2.5",
        "mistral",
        "mistral:latest",
        "deepseek-r1",
    ]
    available = [m.strip() for m in (models or []) if isinstance(m, str) and m.strip()]
    if not available:
        return preferred[0]
    for candidate in preferred:
        if candidate in available:
            return candidate
    return available[0]


async def fetch_available_ollama_models() -> List[str]:
    """Ask Ollama which models are currently available."""
    try:
        async with httpx.AsyncClient(timeout=2.0) as client:
            resp = await client.get(f"{OLLAMA_URL}/api/tags")
            if resp.status_code != 200:
                return []
            payload = resp.json()
            models = payload.get("models", [])
            return [item.get("name", "") for item in models if isinstance(item, dict)]
    except Exception:
        return []


def resolve_llm_provider(
    ollama_available: bool,
    groq_api_key: Optional[str] = None,
    preferred_provider: Optional[str] = None,
    available_models: Optional[List[str]] = None,
) -> Dict[str, str]:
    """Choose the active LLM provider, preferring Ollama locally and allowing Groq as a fallback when configured."""
    override = (preferred_provider or DEFAULT_LLM_PROVIDER or "ollama").lower()
    key = (groq_api_key or GROQ_API_KEY or "").strip()
    available = [m.strip() for m in (available_models or []) if isinstance(m, str) and m.strip()]

    if override == "groq" and key:
        return {"provider": "groq", "model": GROQ_MODEL}
    if override == "ollama" and ollama_available:
        return {"provider": "ollama", "model": detect_best_ollama_model(available)}
    if key:
        return {"provider": "groq", "model": GROQ_MODEL}
    return {"provider": "rule-based", "model": "rule-based-evidence-fallback"}


def build_risk_progression(req: Any, decision: str) -> Dict[str, int]:
    """Return the initial and policy-adjusted final risk scores for the investigation."""
    initial = max(0, min(100, int(getattr(req, "risk_score", 0) or 0)))
    final = initial
    if decision == "BLOCK":
        final = min(initial + 8, 100)
    elif decision == "ALLOW":
        final = max(initial - 8, 0)
    elif decision == "CHALLENGE":
        final = max(initial, 60)
    return {"initial": initial, "final": final}


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
    last_tx_status = customer_history.get("last_transaction_status")
    prior_fraud_count = int(customer_history.get("previous_fraud_count", 0) or 0)
    recent_transactions = customer_history.get("recent_transactions") or []
    successful_history = [
        tx for tx in recent_transactions
        if str(tx.get("status", "")).upper() in {"APPROVED", "VERIFIED"}
    ]

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

    if last_tx_status in {"APPROVED", "VERIFIED"} and decision in {"BLOCK", "CHALLENGE"} and recent_transactions:
        user_context = (
            f"• User-specific retrieval context: the customer’s last transaction was {last_tx_status.upper()} and the historical sequence shows {len(successful_history)} successful approval(s) before this event. "
            f"This current transaction differs from that trusted pattern because it reached ₹{amount:,.2f}, which exceeds the customer’s historical max of ₹{amount_threshold:,.2f}."
        )
    elif prior_fraud_count > 0 and last_tx_status in {"BLOCKED", "CANCELLED"}:
        user_context = (
            f"• User-specific retrieval context: this customer’s last transaction was {last_tx_status.upper()}, "
            "so the current decision is evaluated against the same historical fraud pattern before the next transaction is accepted."
        )
    elif prior_fraud_count > 0:
        user_context = (
            f"• User-specific retrieval context: this customer has {prior_fraud_count} prior risky event(s), which the policy and investigation layer treats as relevant historical context."
        )
    else:
        user_context = "• User-specific retrieval context: no recent prior fraud pattern was detected, so this decision relies mostly on the current transaction and profile checks."

    if decision == "BLOCK":
        bits = [
            f"• Policy blocked transaction {getattr(req, 'transaction_id', 'UNKNOWN')} for {user_id or 'this customer'}.",
            f"• The recorded risk score was {risk_score}/100. The final authorization decision was made by the policy engine, not by the raw ML score alone.",
            user_context,
        ]
        if device_risk:
            bits.append(f"• The device signal is suspicious: {device_profile.get('device_id', 'device')} looks untrusted or emulator-like.")
        elif ip_risk:
            bits.append(f"• The IP signal is suspicious: {ip_profile.get('ip_address', 'IP')} shows elevated fraud risk.")
        if abnormal_amount:
            bits.append(f"• The amount of ₹{amount:,.2f} materially exceeds the customer baseline of ₹{amount_threshold:,.2f}.")
        if last_tx_status in {"APPROVED", "VERIFIED"} and recent_transactions:
            bits.append(
                "• Historical sequence: this customer had multiple recent approved payments, but the current payment broke the prior approval pattern because the amount moved far outside the user’s normal behavior."
            )
        if linked_graph:
            bits.append("• Network links show related fraud behavior, which adds context to the current investigation.")
        if prior_fraud_count:
            bits.append(f"• This customer has {prior_fraud_count} prior fraud-related transaction(s), which strengthens the risk context.")
        bits.append("• Final outcome: BLOCK because the transaction exceeded the configured risk boundary and the evidence stack did not support approval.")
        return "\n".join(bits)

    if decision == "CHALLENGE":
        bits = [
            f"• Policy challenged transaction {getattr(req, 'transaction_id', 'UNKNOWN')} for {user_id or 'this customer'}.",
            f"• The risk score reached {risk_score}/100, so the policy engine required step-up verification before approval.",
            user_context,
        ]
        if device_risk:
            bits.append("• Device checks are weaker than the customer’s trusted baseline, which is why the challenge flow was triggered.")
        elif ip_risk:
            bits.append("• The IP profile is elevated enough to justify identity verification before the payment is allowed.")
        if abnormal_amount:
            bits.append(f"• The amount is unusually high for the customer, which increases the risk of a takeover or mule pattern.")
        if last_tx_status in {"APPROVED", "VERIFIED"} and recent_transactions:
            bits.append(
                "• Historical sequence: even though the user’s recent payments were accepted, this new transaction is materially different and therefore the policy is requiring a challenge before authorizing it."
            )
        if last_tx_status in {"BLOCKED", "CANCELLED"}:
            bits.append("• The customer’s last transaction was already marked as risky or blocked, which is relevant user-specific context for this challenge decision.")
        bits.append("• Final outcome: CHALLENGE. If the OTP is verified, the decision can move to ALLOW; if it fails or expires, it remains blocked.")
        return "\n".join(bits)

    if status == "APPROVED" or "APPROVED" in status or "VERIFIED" in status:
        bits = [
            f"• Policy approved transaction {getattr(req, 'transaction_id', 'UNKNOWN')} after the customer completed the required verification flow.",
            f"• The transaction stayed within the policy posture despite earlier risk markers; the successful OTP outcome cleared the challenge state.",
            user_context,
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
        user_context,
    ]
    if device_risk:
        bits.append("• The device signal was reviewed but did not exceed the threshold that would require a block.")
    elif ip_risk:
        bits.append("• The IP signal was elevated but still below the policy threshold for a block.")
    if abnormal_amount:
        bits.append("• The amount stayed above the customer’s typical range, but it was still within the permitted policy posture.")
    if prior_fraud_count:
        bits.append(f"• Prior history for this customer includes {prior_fraud_count} prior risky event(s), but the current profile still qualified for approval under policy.")
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
            fact=f"Transaction amount ₹{req.amount:,.2f} significantly exceeds the customer’s historical max of ₹{typical_max:,.2f}",
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

    risk_progression = build_risk_progression(req, decision)
    summary = build_decision_summary(
        req=req,
        customer_history=customer_history,
        device_profile=device_profile,
        ip_profile=ip_profile,
        related_accounts=related_accounts,
    )
    summary = (
        f"• Initial risk: {risk_progression['initial']}/100\n"
        f"• Final risk after policy review: {risk_progression['final']}/100\n\n"
        + summary
    )

    cases = [
        SimilarCase(
            case_id=c["case_id"],
            title=c["title"],
            relevance_score=c["relevance_score"]
        )
        for c in similar_cases[:2]
    ]
    retrieval_trace = [{
        "source": "customer_history",
        "query": f"user {req.user_id} previous fraud pattern and last transaction status",
        "matched_documents": [
            f"last_transaction_status={customer_history.get('last_transaction_status') or 'unknown'}",
            f"previous_fraud_count={customer_history.get('previous_fraud_count', 0)}",
            f"historical_max_amount={customer_history.get('typical_max_amount', 'unknown')}",
            f"recent_tx_count={len(customer_history.get('recent_transactions', []))}",
        ],
        "relevance_score": 0.94,
    }]
    if similar_cases:
        retrieval_trace.append({
            "source": "rag",
            "query": f"{req.user_id} {req.amount} abnormal transaction pattern comparison",
            "matched_documents": [case["title"] for case in similar_cases[:2]],
            "relevance_score": max((case.get("relevance_score", 0.0) for case in similar_cases[:2]), default=0.0),
        })
    elif similar_cases is not None:
        retrieval_trace.append({
            "source": "rag",
            "query": f"{req.user_id} {req.amount} abnormal transaction pattern comparison",
            "matched_documents": ["No stored RAG match; deterministic policy evidence used instead."],
            "relevance_score": 0.0,
        })

    reasoning_trace = [
        f"Initial risk score = {risk_progression['initial']}/100 from the ML model and velocity signals.",
        f"Structured evidence reviewed: device trust={device_profile.get('trust_score', 'n/a')}, IP risk={ip_profile.get('risk_score', 'n/a')}, historical max amount={customer_history.get('typical_max_amount', 'n/a')}.",
        f"Policy decision = {decision}; final risk score = {risk_progression['final']}/100 after challenge or verification outcome was considered.",
    ]

    recent_history = customer_history.get("recent_transactions") or []
    recent_history_text = " | ".join(
        f"{tx.get('status', 'UNKNOWN')}:{tx.get('amount', 0):,.2f}"
        for tx in recent_history[:5]
    ) if recent_history else "no recent transactions available"

    inv_id = f"INV-{req.transaction_id}"
    llm_model = "rule-based-evidence-fallback"
    llm_prompt = ""
    confidence = 0.58 + min(0.2, 0.05 * max(0, len(evidence))) + (0.1 if similar_cases else 0.0) + (0.08 if decision in {"CHALLENGE", "BLOCK"} else 0.0)
    confidence = round(min(0.97, confidence), 2)

    try:
        available_models = await fetch_available_ollama_models()
        ollama_available = bool(available_models)
        provider_config = resolve_llm_provider(
            ollama_available=ollama_available,
            available_models=available_models,
        )
        provider = provider_config["provider"]
        llm_model = provider_config["model"]

        if provider == "ollama":
            selected_model = detect_best_ollama_model(available_models)
            llm_model = selected_model
            llm_prompt = (
                f"You are a Fraud Analyst AI. Analyze transaction {req.transaction_id} for user {req.user_id}.\n"
                f"Amount: ₹{req.amount:,.2f}, Risk Score: {req.risk_score}/100, Fraud Probability: {req.fraud_probability}, Anomaly Score: {req.anomaly_score}.\n"
                f"Authoritative policy decision: {decision}. Transaction status: {req.status}.\n"
                f"Customer retrieval context: last_transaction_status={customer_history.get('last_transaction_status', 'unknown')}; previous_fraud_count={customer_history.get('previous_fraud_count', 0)}; historical_max_amount=₹{customer_history.get('typical_max_amount', 'unknown')}; recent_transaction_history={recent_history_text}.\n"
                f"Device profile: device_id={device_profile.get('device_id', 'unknown')}; trust_score={device_profile.get('trust_score', 'n/a')}; is_emulator={device_profile.get('is_emulator', False)}.\n"
                f"IP profile: ip_address={ip_profile.get('ip_address', 'unknown')}; risk_score={ip_profile.get('risk_score', 'n/a')}; is_vpn={ip_profile.get('is_vpn', False)}.\n"
                f"Evidence: {[e.fact for e in evidence]}\n"
                "Return a 5-bullet structured explanation that explains why this transaction succeeded or failed under the policy engine using only the supplied evidence and the authoritative decision. Do not invent facts or claim hidden details."
            )
            async with httpx.AsyncClient(timeout=10.0) as client:
                res = await client.post(
                    f"{OLLAMA_URL}/api/generate",
                    json={"model": selected_model, "prompt": llm_prompt, "stream": False}
                )
                payload = res.json() if res.headers.get("content-type", "").startswith("application/json") else {}
                if res.status_code != 200:
                    error_text = payload.get("error") if isinstance(payload, dict) else str(payload)
                    raise RuntimeError(f"Ollama generation failed with status {res.status_code}: {error_text}")

                llm_out = payload.get("response", "").strip()
                if not llm_out:
                    raise RuntimeError("Ollama returned an empty response for the investigation prompt.")
                summary = llm_out
        elif provider == "groq":
            llm_prompt = (
                f"Transaction {req.transaction_id} for user {req.user_id}.\n"
                f"Amount: ₹{req.amount:,.2f}; Risk score: {req.risk_score}; Fraud probability: {req.fraud_probability}; Anomaly score: {req.anomaly_score}; Policy decision: {decision}; Status: {req.status}.\n"
                f"Last transaction status: {customer_history.get('last_transaction_status')}; previous_fraud_count: {customer_history.get('previous_fraud_count', 0)}.\n"
                f"Recent customer history: {recent_history_text}.\n"
                f"Evidence: {[e.fact for e in evidence]}\n"
                "Return a 5-bullet structured explanation grounded in the provided evidence and the policy decision."
            )
            groq_api_key = (os.getenv("GROQ_API_KEY") or "").strip()
            if not groq_api_key:
                raise RuntimeError("GROQ_API_KEY is not configured, so the Groq fallback cannot be used.")

            messages = [
                {
                    "role": "system",
                    "content": "You are a fraud analyst. Use only the supplied evidence and historical context. Do not invent transaction facts. Explain if the policy decision was ALLOW, CHALLENGE, or BLOCK and reference the user history accurately.",
                },
                {
                    "role": "user",
                    "content": (
                        f"Transaction {req.transaction_id} for user {req.user_id}.\n"
                        f"Amount: ₹{req.amount:,.2f}; Risk score: {req.risk_score}; Fraud probability: {req.fraud_probability}; Policy decision: {decision}; Status: {req.status}.\n"
                        f"Last transaction status: {customer_history.get('last_transaction_status')}; previous_fraud_count: {customer_history.get('previous_fraud_count', 0)}.\n"
                        f"Evidence: {[e.fact for e in evidence]}\n"
                        "Return a 5-bullet structured explanation that is grounded in the provided evidence and the policy decision."
                    ),
                },
            ]
            async with httpx.AsyncClient(timeout=15.0) as client:
                res = await client.post(
                    "https://api.groq.com/openai/v1/chat/completions",
                    headers={
                        "Authorization": f"Bearer {groq_api_key}",
                        "Content-Type": "application/json",
                    },
                    json={
                        "model": GROQ_MODEL,
                        "messages": messages,
                        "temperature": 0.3,
                    },
                )
                payload = res.json() if res.headers.get("content-type", "").startswith("application/json") else {}
                if res.status_code != 200:
                    error_text = payload.get("error", {}).get("message", str(payload)) if isinstance(payload, dict) else str(payload)
                    raise RuntimeError(f"Groq generation failed with status {res.status_code}: {error_text}")

                choices = payload.get("choices") or []
                if not choices:
                    raise RuntimeError("Groq returned no completion choices for the investigation prompt.")
                llm_out = ((choices[0].get("message") or {}).get("content") or "").strip()
                if not llm_out:
                    raise RuntimeError("Groq returned an empty investigation response.")
                summary = llm_out
        else:
            raise RuntimeError("No LLM provider is configured for investigation generation.")
    except Exception as exc:
        logger.exception("AI generation failed; falling back to deterministic summary: %s", exc)

    return InvestigationResponse(
        investigation_id=inv_id,
        transaction_id=req.transaction_id,
        risk_level=risk_level,
        summary=summary,
        evidence=evidence,
        similar_cases=cases,
        recommended_action=recommended_action,
        confidence=confidence,
        llm_model=llm_model,
        llm_prompt=llm_prompt,
        initial_risk_score=risk_progression["initial"],
        final_risk_score=risk_progression["final"],
        retrieval_trace=retrieval_trace,
        reasoning_trace=reasoning_trace,
    )
