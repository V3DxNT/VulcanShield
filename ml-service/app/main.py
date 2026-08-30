from fastapi import FastAPI
from app.schemas import PredictRequest, PredictResponse
from app.predictor import MLPredictor

app = FastAPI(title="VulcanShield ML Service", version="1.0.0")
predictor = None

@app.on_event("startup")
def startup_event():
    global predictor
    predictor = MLPredictor()

@app.get("/health")
def health():
    return {
        "status": "ok",
        "service": "ml-service",
        "version": "1.0.0",
    }

@app.post("/predict", response_model=PredictResponse)
def predict(req: PredictRequest):
    return predictor.predict(req)
