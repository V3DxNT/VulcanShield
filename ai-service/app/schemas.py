from pydantic import BaseModel, Field
from typing import List, Dict, Any, Optional

class InvestigationRequest(BaseModel):
    transaction_id: str = Field(..., example="TX-1001")
    user_id: str = Field(..., example="C1001")
    device_id: str = Field(default="", example="D204")
    ip_address: str = Field(default="", example="IP-17")
    amount: float = Field(..., example=1500.00)
    risk_score: int = Field(default=85, example=85)
    fraud_probability: float = Field(default=0.85, example=0.85)
    anomaly_score: float = Field(default=0.90, example=0.90)
    status: str = Field(default="", example="APPROVED")
    decision: str = Field(default="", example="ALLOW")

class EvidenceItem(BaseModel):
    category: str = Field(..., example="DEVICE_INTELLIGENCE")
    fact: str = Field(..., example="Device D206 is flagged as an emulator")
    severity: str = Field(..., example="HIGH")

class SimilarCase(BaseModel):
    case_id: str = Field(..., example="RAG-DOC-001")
    title: str = Field(..., example="Velocity Attack Patterns in Carding Operations")
    relevance_score: float = Field(..., example=0.92)

class RetrievalTraceItem(BaseModel):
    source: str = Field(..., example="customer_history")
    query: str = Field(..., example="user C1001 previous fraud pattern")
    matched_documents: List[str] = Field(default_factory=list)
    relevance_score: float = Field(default=0.0, ge=0.0, le=1.0)

class InvestigationResponse(BaseModel):
    investigation_id: str
    transaction_id: str
    risk_level: str = Field(..., example="HIGH") # LOW, MEDIUM, HIGH, CRITICAL
    summary: str
    evidence: List[EvidenceItem] = Field(default_factory=list)
    similar_cases: List[SimilarCase] = Field(default_factory=list)
    recommended_action: str = Field(..., example="BLOCK_ACCOUNT") # ALLOW, CHALLENGE, BLOCK_ACCOUNT, MANUAL_REVIEW
    confidence: float = Field(..., ge=0.0, le=1.0)
    llm_model: str = Field(default="qwen2.5:7b-instruct")
    initial_risk_score: int = Field(default=0)
    final_risk_score: int = Field(default=0)
    retrieval_trace: List[RetrievalTraceItem] = Field(default_factory=list)
    reasoning_trace: List[str] = Field(default_factory=list)
