'use client';
import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';
import NetworkGraph, { GraphEdge } from '@/components/NetworkGraph';

export default function GraphVisualizer() {
  const [relationships, setRelationships] = useState<GraphEdge[]>([]);
  const [userID, setUserID] = useState('');

  const fetchGraph = useCallback(async () => {
    try {
      const res = await fetch('/api/v1/graph/relationships?limit=200');
      if (res.ok) {
        const json = await res.json();
        const data: GraphEdge[] = json.data ?? [];
        setRelationships(data);
        const firstUser = data.find(edge => edge.source_type === 'USER')?.source_id ?? data.find(edge => edge.target_type === 'USER')?.target_id ?? '';
        setUserID(current => current || firstUser);
      }
    } catch {}
  }, []);

  useEffect(() => {
    fetchGraph();
    const t = setInterval(fetchGraph, 3000);
    return () => clearInterval(t);
  }, [fetchGraph]);

  const users = Array.from(new Set(relationships.flatMap(r => [r.source_type === 'USER' ? r.source_id : '', r.target_type === 'USER' ? r.target_id : '']).filter(Boolean))).sort();
  const userRelationships = relationships.filter(r => r.source_id === userID || r.target_id === userID);
  const fraudCount = userRelationships.filter(r => r.fraud_linked).length;

  return (
    <div className="bg-white rounded-2xl border border-[#d2d2d7] shadow-sm">
      <div className="px-6 py-4 border-b border-[#e8e8ed] flex items-center justify-between">
        <div>
          <h2 className="font-semibold text-[#1d1d1f] text-base">Fraud Network Graph</h2>
          <p className="text-xs text-[#6e6e73] mt-0.5">Individual user neighborhood · {userRelationships.length} direct edges · {fraudCount} fraud-linked</p>
        </div>
        <div className="flex items-center gap-2">
          <select value={userID} onChange={e => setUserID(e.target.value)} className="rounded-lg border border-[#d2d2d7] bg-white px-2.5 py-1.5 text-xs font-mono text-[#1d1d1f]">
            {users.map(id => <option key={id} value={id}>{id}</option>)}
          </select>
          <Link href="/network" className="px-3.5 py-1.5 text-xs font-medium text-[#0071e3] border border-[#0071e3]/30 bg-blue-50 rounded-lg hover:bg-blue-100 transition-colors">User network →</Link>
        </div>
      </div>
      <div className="p-4">
        <NetworkGraph edges={userRelationships} selectedId={userID} height={300} />
      </div>
    </div>
  );
}
