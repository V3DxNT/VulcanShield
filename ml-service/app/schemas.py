from pydantic import BaseModel, Field
from typing import Dict, Any, Optional

class PredictRequest(BaseModel):
    transaction_id: str = Field(..., example="TX-1001")
    user_id: str = Field(..., example="C1001")
    amount: float = Field(..., example=120.50)
    typical_max_amount: float = Field(default=250.0, example=250.0)
    user_tx_count_60s: int = Field(default=1, example=3)
    ip_tx_count_60s: int = Field(default=1, example=3)
    device_tx_count_60s: int = Field(default=1, example=1)
    trust_score: int = Field(default=85, example=85)
    is_emulator: bool = Field(default=False, example=False)
    is_vpn: bool = Field(default=False, example=False)

class PredictResponse(BaseModel):
    transaction_id: str
    fraud_probability: float = Field(..., ge=0.0, le=1.0)
    anomaly_score: float = Field(..., ge=0.0, le=1.0)
    model_version: str = Field(default="v1.0")
    feature_snapshot: Dict[str, Any] = Field(default_factory=dict)
