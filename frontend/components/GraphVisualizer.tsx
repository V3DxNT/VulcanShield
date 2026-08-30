'use client';
import { useState, useEffect, useCallback } from 'react';
import Link from 'next/link';

interface Relationship {
  relationship_id: string;
  source_type: string;
  source_id: string;
  target_type: string;
  target_id: string;
  relationship_type: string;
  weight: number;
  fraud_linked: boolean;
}

interface NodeStats {
  userId: string;
  deviceCount: number;
  ipCount: number;
  txCount: number;
  fraudLinked: boolean;
  devices: string[];
  ips: string[];
}

function buildNodeStats(rels: Relationship[]): NodeStats[] {
  const map = new Map<string, NodeStats>();

  rels.forEach(r => {
    if (r.source_type === 'USER' && !map.has(r.source_id)) {
      map.set(r.source_id, { userId: r.source_id, deviceCount: 0, ipCount: 0, txCount: 0, fraudLinked: false, devices: [], ips: [] });
    }
    if (r.target_type === 'USER' && !map.has(r.target_id)) {
      map.set(r.target_id, { userId: r.target_id, deviceCount: 0, ipCount: 0, txCount: 0, fraudLinked: false, devices: [], ips: [] });
    }

    const user = r.source_type === 'USER' ? map.get(r.source_id) : r.target_type === 'USER' ? map.get(r.target_id) : null;
    if (user) {
      if (r.relationship_type === 'DEVICE_SHARED' || r.relationship_type === 'SHARED_DEVICE') {
        const other = r.source_type === 'DEVICE' ? r.source_id : r.target_id;
        if (!user.devices.includes(other)) { user.devices.push(other); user.deviceCount++; }
      }
      if (r.relationship_type === 'IP_SHARED' || r.relationship_type === 'SHARED_IP') {
        const other = r.source_type === 'IP' ? r.source_id : r.target_id;
        if (!user.ips.includes(other)) { user.ips.push(other); user.ipCount++; }
      }
      if (r.fraud_linked) user.fraudLinked = true;
    }
  });

  return Array.from(map.values());
}

export default function GraphVisualizer() {
  const [relationships, setRelationships] = useState<Relationship[]>([]);
  const [nodeStats, setNodeStats] = useState<NodeStats[]>([]);
  const [selected, setSelected] = useState<Relationship | null>(null);

  const fetchGraph = useCallback(async () => {
    try {
      const res = await fetch('/api/v1/graph/relationships?limit=50');
      if (res.ok) {
        const json = await res.json();
        const data = json.data ?? [];
        setRelationships(data);
        setNodeStats(buildNodeStats(data));
      }
    } catch {}
  }, []);

  useEffect(() => { fetchGraph(); }, [fetchGraph]);

  const fraudCount = relationships.filter(r => r.fraud_linked).length;
  const userNodes = Array.from(new Set(relationships.flatMap(r => [r.source_id, r.target_id])));

  return (
    <div className="bg-white rounded-2xl border border-[#d2d2d7] shadow-sm">
      <div className="px-6 py-4 border-b border-[#e8e8ed] flex items-center justify-between">
        <div>
          <h2 className="font-semibold text-[#1d1d1f] text-base">Fraud Network Graph</h2>
          <p className="text-xs text-[#6e6e73] mt-0.5">{relationships.length} edges · {userNodes.length} unique nodes · {fraudCount} fraud-linked</p>
        </div>
        <Link href="/network"
          className="px-3.5 py-1.5 text-xs font-medium text-[#0071e3] border border-[#0071e3]/30 bg-blue-50 rounded-lg hover:bg-blue-100 transition-colors">
          Full Network View →
        </Link>
      </div>

      <div className="p-6">
        {relationships.length === 0 ? (
          <p className="text-sm text-center text-[#a1a1a6] py-8">No relationships extracted yet. Start a device_farm or ip_abuse scenario to generate graph data.</p>
        ) : (
          <div className="space-y-4">
            {/* Summary Tiles */}
            <div className="grid grid-cols-3 gap-3 mb-4">
              <div className="p-3 rounded-xl bg-[#f5f5f7] border border-[#e8e8ed] text-center">
                <p className="text-2xl font-semibold text-[#1d1d1f]">{userNodes.length}</p>
                <p className="text-xs text-[#6e6e73] mt-0.5">Unique Nodes</p>
              </div>
              <div className="p-3 rounded-xl bg-[#f5f5f7] border border-[#e8e8ed] text-center">
                <p className="text-2xl font-semibold text-[#1d1d1f]">{relationships.length}</p>
                <p className="text-xs text-[#6e6e73] mt-0.5">Total Edges</p>
              </div>
              <div className="p-3 rounded-xl bg-red-50 border border-red-100 text-center">
                <p className="text-2xl font-semibold text-red-600">{fraudCount}</p>
                <p className="text-xs text-red-500 mt-0.5">Fraud-Linked</p>
              </div>
            </div>

            {/* Edge List */}
            <div className="max-h-56 overflow-y-auto space-y-2">
              {relationships.slice(0, 20).map(r => (
                <div key={r.relationship_id}
                  onClick={() => setSelected(selected?.relationship_id === r.relationship_id ? null : r)}
                  className={`flex items-center gap-3 px-3.5 py-2.5 rounded-xl border cursor-pointer transition-colors ${
                    r.fraud_linked ? 'bg-red-50 border-red-200 hover:bg-red-100' : 'bg-[#f5f5f7] border-[#e8e8ed] hover:bg-[#eeeef0]'
                  } ${selected?.relationship_id === r.relationship_id ? 'ring-2 ring-[#0071e3]' : ''}`}>
                  <span className={`font-mono text-xs font-semibold ${r.fraud_linked ? 'text-red-600' : 'text-[#0071e3]'}`}>{r.source_id}</span>
                  <span className="text-xs text-[#a1a1a6] flex-shrink-0">─ {r.relationship_type.replace(/_/g, ' ')} →</span>
                  <span className={`font-mono text-xs font-semibold ${r.fraud_linked ? 'text-red-600' : 'text-purple-600'}`}>{r.target_id}</span>
                  {r.fraud_linked && <span className="ml-auto text-xs font-medium px-2 py-0.5 rounded-full bg-red-100 text-red-600 border border-red-200">Fraud</span>}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
