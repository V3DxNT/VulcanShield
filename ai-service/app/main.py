from fastapi import FastAPI
from app.schemas import InvestigationRequest, InvestigationResponse
from app.tools import get_db_pool, get_customer_history, get_device_profile, get_ip_profile, get_related_accounts
from app.rag import get_similar_fraud_cases
from app.ollama import generate_llm_investigation

app = FastAPI(title="VulcanShield AI Service", version="1.0.0")
db_pool = None

@app.on_event("startup")
async def startup_event():
    global db_pool
    try:
        db_pool = await get_db_pool()
    except Exception as e:
        print(f"Warning: PostgreSQL DB pool error: {e}")

@app.get("/health")
def health():
    return {
        "status": "ok",
        "service": "ai-service",
        "version": "1.0.0",
    }

@app.post("/ai/investigate", response_model=InvestigationResponse)
async def investigate(req: InvestigationRequest):
    customer_history = {}
    device_profile = {}
    ip_profile = {}
    related_accounts = []
    similar_cases = []

    if db_pool:
        try:
            customer_history = await get_customer_history(db_pool, req.user_id)
            device_profile = await get_device_profile(db_pool, req.device_id)
            ip_profile = await get_ip_profile(db_pool, req.ip_address)
            related_accounts = await get_related_accounts(db_pool, req.user_id)
            similar_cases = await get_similar_fraud_cases(db_pool, f"Transaction {req.transaction_id}")
        except Exception as e:
            print(f"Tool execution warning: {e}")

    return await generate_llm_investigation(
        req=req,
        customer_history=customer_history,
        device_profile=device_profile,
        ip_profile=ip_profile,
        related_accounts=related_accounts,
        similar_cases=similar_cases,
    )
