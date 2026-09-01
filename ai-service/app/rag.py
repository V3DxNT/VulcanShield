import asyncpg
from typing import List, Dict, Any, Optional
from app.embeddings import embed_query

async def get_similar_fraud_cases(pool: asyncpg.Pool, query_text: str) -> List[Dict[str, Any]]:
    """Hybrid global RAG: semantic pgvector search, then lexical fallback.

    Retrieved documents are contextual guidance only; structured database
    tools remain the source of truth for the current transaction.
    """
    embedding: Optional[List[float]] = await embed_query(query_text)
    async with pool.acquire() as conn:
                                                                              
                                                                               
        rows = await conn.fetch(
            """
            SELECT document_id, title, category, content,
                   ts_rank(to_tsvector('english', title || ' ' || content),
                           websearch_to_tsquery('english', $1)) AS relevance_score
            FROM rag_documents
            WHERE to_tsvector('english', title || ' ' || content)
                  @@ websearch_to_tsquery('english', $1)
            ORDER BY relevance_score DESC, document_id
            LIMIT 3
            """,
            query_text,
        )
        if not rows and embedding:
            vector = "[" + ",".join(f"{value:.8f}" for value in embedding) + "]"
            rows = await conn.fetch(
                """
                SELECT d.document_id, d.title, d.category, d.content,
                       1 - (e.embedding <=> $1::vector) AS relevance_score
                FROM embedding_records e
                JOIN rag_documents d ON d.document_id = e.document_id
                ORDER BY e.embedding <=> $1::vector
                LIMIT 3
                """,
                vector,
            )
        return [
            {
                "case_id": r["document_id"],
                "title": r["title"],
                "category": r["category"],
                "content": r["content"][:200],
                "relevance_score": max(0.0, min(1.0, float(r["relevance_score"]))),
            }
            for r in rows
        ]
