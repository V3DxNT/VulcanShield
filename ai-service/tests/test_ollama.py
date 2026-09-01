import unittest

from app.ollama import build_decision_summary, detect_best_ollama_model, resolve_llm_provider


class OllamaSummaryTests(unittest.TestCase):
    def test_user_specific_last_fraud_context_is_included(self):
        req = type('Req', (), {
            'transaction_id': 'TX-9001',
            'user_id': 'C1008',
            'decision': 'ALLOW',
            'status': 'APPROVED',
            'risk_score': 62,
            'amount': 2600,
        })()

        customer_history = {
            'typical_max_amount': 2200,
            'previous_fraud_count': 1,
            'last_transaction_status': 'BLOCKED',
        }

        summary = build_decision_summary(req, customer_history=customer_history)
        self.assertIn('last transaction', summary.lower())
        self.assertIn('fraud', summary.lower())

    def test_ollama_model_selection_prefers_available_model(self):
        models = [
            'llama3.2:3b',
            'qwen2.5:7b-instruct',
            'mistral:latest',
        ]
        self.assertEqual(detect_best_ollama_model(models), 'qwen2.5:7b-instruct')

    def test_resolve_llm_provider_uses_available_ollama_model(self):
        provider = resolve_llm_provider(True, available_models=['llama3.2:3b', 'qwen2.5:7b-instruct'])
        self.assertEqual(provider['provider'], 'ollama')
        self.assertEqual(provider['model'], 'qwen2.5:7b-instruct')


if __name__ == '__main__':
    unittest.main()
