import os
import joblib
import numpy as np
import pandas as pd
from app.schemas import PredictRequest, PredictResponse

class MLPredictor:
    def __init__(self):
        base_dir = os.path.dirname(os.path.dirname(__file__))
        models_dir = os.path.join(base_dir, "models")
        
        xgb_path = os.path.join(models_dir, "xgboost_model.joblib")
        iso_path = os.path.join(models_dir, "isolation_forest.joblib")
        
        if not os.path.exists(xgb_path) or not os.path.exists(iso_path):
            print("Models not found; generating synthetic training data and fitting models...")
            from train import train_and_save_models
            train_and_save_models()
            
        self.xgb = joblib.load(xgb_path)
        self.iso = joblib.load(iso_path)
        self.feature_cols = [
            "amount", "typical_max_amount", "amount_ratio",
            "user_tx_count_60s", "ip_tx_count_60s", "device_tx_count_60s",
            "trust_score", "is_emulator", "is_vpn"
        ]

    def predict(self, req: PredictRequest) -> PredictResponse:
        amount_ratio = req.amount / (req.typical_max_amount + 1e-5)
        
        features_dict = {
            "amount": req.amount,
            "typical_max_amount": req.typical_max_amount,
            "amount_ratio": amount_ratio,
            "user_tx_count_60s": req.user_tx_count_60s,
            "ip_tx_count_60s": req.ip_tx_count_60s,
            "device_tx_count_60s": req.device_tx_count_60s,
            "trust_score": req.trust_score,
            "is_emulator": 1 if req.is_emulator else 0,
            "is_vpn": 1 if req.is_vpn else 0,
        }
        
        df = pd.DataFrame([features_dict])[self.feature_cols]
        
        # 1. Fraud Probability from XGBoost
        fraud_prob = float(self.xgb.predict_proba(df)[0][1])
        
        # 2. Anomaly Score from Isolation Forest (normalized 0.0 - 1.0)
        raw_iso_score = float(self.iso.score_samples(df)[0]) # Higher is less anomalous
        # Transform score_samples [-0.8, -0.2] to anomaly_score [0.0, 1.0]
        anomaly_score = float(np.clip(1.0 - (raw_iso_score + 0.8) / 0.6, 0.0, 1.0))
        
        return PredictResponse(
            transaction_id=req.transaction_id,
            fraud_probability=round(fraud_prob, 4),
            anomaly_score=round(anomaly_score, 4),
            model_version="v1.0",
            feature_snapshot=features_dict,
        )
