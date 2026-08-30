import os
import joblib
import numpy as np
import pandas as pd
from sklearn.ensemble import IsolationForest
from xgboost import XGBClassifier

def generate_synthetic_dataset(n_samples=2000, seed=42):
    np.random.seed(seed)
    
    # Feature distribution matching VulcanShield Phase 2 baseline entities
    amount = np.random.exponential(scale=100.0, size=n_samples) + 5.0
    typical_max_amount = np.random.choice([100.0, 250.0, 1500.0], size=n_samples)
    user_tx_count_60s = np.random.poisson(lam=1.5, size=n_samples)
    ip_tx_count_60s = np.random.poisson(lam=1.2, size=n_samples)
    device_tx_count_60s = np.random.poisson(lam=1.1, size=n_samples)
    trust_score = np.random.choice([30, 60, 85, 90], size=n_samples)
    is_emulator = np.random.choice([0, 1], p=[0.9, 0.1], size=n_samples)
    is_vpn = np.random.choice([0, 1], p=[0.85, 0.15], size=n_samples)
    
    # Calculate amount ratio
    amount_ratio = amount / (typical_max_amount + 1e-5)
    
    # Synthetic target label (fraud: 1, normal: 0)
    fraud_prob_raw = (
        (amount_ratio > 3.0).astype(float) * 0.4 +
        (user_tx_count_60s > 5).astype(float) * 0.35 +
        (is_emulator == 1).astype(float) * 0.3 +
        (is_vpn == 1).astype(float) * 0.2 +
        (trust_score < 40).astype(float) * 0.25
    )
    is_fraud = (fraud_prob_raw + np.random.normal(0, 0.1, n_samples) > 0.45).astype(int)
    
    df = pd.DataFrame({
        "amount": amount,
        "typical_max_amount": typical_max_amount,
        "amount_ratio": amount_ratio,
        "user_tx_count_60s": user_tx_count_60s,
        "ip_tx_count_60s": ip_tx_count_60s,
        "device_tx_count_60s": device_tx_count_60s,
        "trust_score": trust_score,
        "is_emulator": is_emulator,
        "is_vpn": is_vpn,
        "is_fraud": is_fraud,
    })
    return df

def train_and_save_models():
    models_dir = os.path.join(os.path.dirname(__file__), "models")
    os.makedirs(models_dir, exist_ok=True)
    
    df = generate_synthetic_dataset()
    feature_cols = [
        "amount", "typical_max_amount", "amount_ratio",
        "user_tx_count_60s", "ip_tx_count_60s", "device_tx_count_60s",
        "trust_score", "is_emulator", "is_vpn"
    ]
    
    X = df[feature_cols]
    y = df["is_fraud"]
    
    # 1. XGBoost Classifier (Fraud Probability)
    xgb = XGBClassifier(
        n_estimators=50,
        max_depth=4,
        learning_rate=0.1,
        random_state=42,
        eval_metric="logloss"
    )
    xgb.fit(X, y)
    
    # 2. Isolation Forest (Anomaly Score)
    iso = IsolationForest(
        n_estimators=50,
        contamination=0.08,
        random_state=42
    )
    iso.fit(X)
    
    xgb_path = os.path.join(models_dir, "xgboost_model.joblib")
    iso_path = os.path.join(models_dir, "isolation_forest.joblib")
    
    joblib.dump(xgb, xgb_path)
    joblib.dump(iso, iso_path)
    print(f"Successfully trained and saved XGBoost ({xgb_path}) & IsolationForest ({iso_path})")

if __name__ == "__main__":
    train_and_save_models()
