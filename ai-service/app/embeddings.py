"""Optional local sentence-transformer embeddings for pgvector retrieval.

The service stays available with its deterministic lexical retrieval path when
the local embedding package/model has not been provisioned yet.
"""

import asyncio
import os
from typing import List, Optional

_model = None
_model_failed = False
MODEL_NAME = os.getenv("EMBEDDING_MODEL", "all-MiniLM-L6-v2")


def _encode_sync(text: str) -> Optional[List[float]]:
    global _model, _model_failed
    if _model_failed:
        return None
    if os.getenv("ENABLE_SEMANTIC_RAG", "false").lower() != "true":
        # Enable only after the local all-MiniLM-L6-v2 package/model has been
        # provisioned.  Lexical retrieval remains deterministic for the demo.
        return None
    try:
        if _model is None:
            from sentence_transformers import SentenceTransformer
            _model = SentenceTransformer(MODEL_NAME)
        return _model.encode(text, normalize_embeddings=True).tolist()
    except Exception as exc:
        # RAG can fall back to lexical retrieval.  Core policy and ML do not
        # depend on a downloaded embedding model being available.
        print(f"Embedding model unavailable; using lexical RAG fallback: {exc}")
        _model_failed = True
        return None


async def embed_query(text: str) -> Optional[List[float]]:
    return await asyncio.to_thread(_encode_sync, text)
