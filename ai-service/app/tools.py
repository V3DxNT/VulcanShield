import os
import asyncpg
from typing import Dict, Any, List

DATABASE_URL = os.getenv("DATABASE_URL", "postgres://vulcan:vulcanpass@postgres:5432/vulcanshield")

async def get_db_pool():
    return await asyncpg.create_pool(DATABASE_URL)

async def get_customer_history(pool: asyncpg.Pool, user_id: str) -> Dict[str, Any]:
    async with pool.acquire() as conn:
        user = await conn.fetchrow(
            "SELECT user_id, name, risk_tolerance, challenge_threshold, block_threshold, account_age_days, trust_score, typical_min_amount, typical_max_amount FROM users WHERE user_id = $1",
            user_id
        )
        if not user:
            return {}
        
        tx_count = await conn.fetchval(
            "SELECT COUNT(*) FROM transactions WHERE user_id = $1", user_id
        )
        blocks = await conn.fetchval(
            "SELECT COUNT(*) FROM transactions WHERE user_id = $1 AND status = 'BLOCKED'", user_id
        )
        
        return {
            "user_id": user["user_id"],
            "name": user["name"],
            "trust_score": user["trust_score"],
            "account_age_days": user["account_age_days"],
            "typical_max_amount": float(user["typical_max_amount"]),
            "historical_tx_count": tx_count,
            "previous_block_count": blocks,
        }

async def get_device_profile(pool: asyncpg.Pool, device_id: str) -> Dict[str, Any]:
    if not device_id:
        return {}
    async with pool.acquire() as conn:
        device = await conn.fetchrow(
            "SELECT device_id, device_type, os, browser, trust_score, is_emulator FROM devices WHERE device_id = $1",
            device_id
        )
        if not device:
            return {}
        return {
            "device_id": device["device_id"],
            "device_type": device["device_type"],
            "os": device["os"],
            "trust_score": device["trust_score"],
            "is_emulator": device["is_emulator"],
        }

async def get_ip_profile(pool: asyncpg.Pool, ip_address: str) -> Dict[str, Any]:
    if not ip_address:
        return {}
    async with pool.acquire() as conn:
        ip = await conn.fetchrow(
            "SELECT ip_address, country_code, city, isp, is_vpn, is_tor, is_proxy, risk_score FROM ips WHERE ip_address = $1",
            ip_address
        )
        if not ip:
            return {}
        return {
            "ip_address": ip["ip_address"],
            "country_code": ip["country_code"],
            "is_vpn": ip["is_vpn"],
            "is_proxy": ip["is_proxy"],
            "risk_score": ip["risk_score"],
        }

async def get_related_accounts(pool: asyncpg.Pool, user_id: str) -> List[Dict[str, Any]]:
    async with pool.acquire() as conn:
        rows = await conn.fetch(
            "SELECT source_id, target_id, relationship_type, weight, fraud_linked FROM fraud_relationships WHERE source_id = $1 OR target_id = $1",
            user_id
        )
        return [
            {
                "source_id": r["source_id"],
                "target_id": r["target_id"],
                "relationship_type": r["relationship_type"],
                "fraud_linked": r["fraud_linked"],
            }
            for r in rows
        ]
