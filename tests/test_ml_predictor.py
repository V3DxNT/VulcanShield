import importlib.util
import sys
import types
from pathlib import Path

base = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(base / 'ml-service'))

ml_app_dir = base / 'ml-service' / 'app'
app_pkg = types.ModuleType('app')
app_pkg.__path__ = [str(ml_app_dir)]
sys.modules['app'] = app_pkg

schemas_spec = importlib.util.spec_from_file_location('app.schemas', ml_app_dir / 'schemas.py')
schemas_module = importlib.util.module_from_spec(schemas_spec)
sys.modules['app.schemas'] = schemas_module
assert schemas_spec and schemas_spec.loader
schemas_spec.loader.exec_module(schemas_module)

predictor_spec = importlib.util.spec_from_file_location('app.predictor', ml_app_dir / 'predictor.py')
predictor_module = importlib.util.module_from_spec(predictor_spec)
sys.modules['app.predictor'] = predictor_module
assert predictor_spec and predictor_spec.loader
predictor_spec.loader.exec_module(predictor_module)

PredictRequest = schemas_module.PredictRequest
MLPredictor = predictor_module.MLPredictor


def test_predictor_flags_high_risk_transaction():
    predictor = MLPredictor()
    req = PredictRequest(
        transaction_id='TX-777',
        user_id='C1001',
        amount=5000.0,
        typical_max_amount=300.0,
        user_tx_count_60s=10,
        ip_tx_count_60s=8,
        device_tx_count_60s=7,
        trust_score=18,
        is_emulator=True,
        is_vpn=True,
    )

    resp = predictor.predict(req)
    assert 0.0 <= resp.fraud_probability <= 1.0
    assert 0.0 <= resp.anomaly_score <= 1.0
    assert resp.fraud_probability > 0.6
    assert resp.anomaly_score > 0.5
