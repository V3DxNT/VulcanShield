import unittest

from app.ollama import build_decision_summary


class DummyRequest:
    def __init__(self):
        self.transaction_id = "TX-1009"
        self.user_id = "C1001"
        self.amount = 18500.0
        self.risk_score = 82
        self.decision = "BLOCK"
        self.status = "BLOCKED"


class InvestigationContextTest(unittest.TestCase):
    def test_failure_after_prior_approvals_describes_user_history(self):
        req = DummyRequest()
        history = {
            "typical_max_amount": 7000.0,
            "previous_fraud_count": 0,
            "last_transaction_status": "APPROVED",
            "recent_transactions": [
                {"status": "APPROVED", "amount": 3200.0},
                {"status": "APPROVED", "amount": 4500.0},
                {"status": "APPROVED", "amount": 5200.0},
            ],
        }
        summary = build_decision_summary(req, customer_history=history)
        self.assertIn("last transaction was approved", summary.lower())
        self.assertIn("historical sequence", summary.lower())
        self.assertIn("₹18,500.00", summary)


if __name__ == "__main__":
    unittest.main()
