'use client';
import { useState, useEffect, useCallback } from 'react';
import Navbar from '@/components/Navbar';

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

interface NodeInfo {
  id: string;
  type: 'USER' | 'DEVICE' | 'IP' | 'MERCHANT';
  fraudLinked: boolean;
  connections: number;
  edges: Relationship[];
}

function EntityBadge({ type }: { type: string }) {
  const map: Record<string, string> = {
    USER: 'bg-blue-50 text-blue-700 border-blue-200',
    DEVICE: 'bg-purple-50 text-purple-700 border-purple-200',
    IP: 'bg-amber-50 text-amber-700 border-amber-200',
    MERCHANT: 'bg-green-50 text-green-700 border-green-200',
  };
  return (
    <span className={`px-2 py-0.5 rounded text-xs font-medium border ${map[type] ?? 'bg-gray-100 text-gray-600 border-gray-200'}`}>
      {type}
    </span>
  );
}

function buildGraph(rels: Relationship[]) {
  const nodes = new Map<string, NodeInfo>();

  const ensureNode = (id: string, type: string) => {
    if (!nodes.has(id)) {
      nodes.set(id, { id, type: type as any, fraudLinked: false, connections: 0, edges: [] });
    }
  };

  rels.forEach(r => {
    ensureNode(r.source_id, r.source_type);
    ensureNode(r.target_id, r.target_type);

    const src = nodes.get(r.source_id)!;
    const tgt = nodes.get(r.target_id)!;
    src.connections++;
    tgt.connections++;
    src.edges.push(r);
    tgt.edges.push(r);

    if (r.fraud_linked) {
      src.fraudLinked = true;
      tgt.fraudLinked = true;
    }
  });

  return Array.from(nodes.values()).sort((a, b) => b.connections - a.connections);
}

function StatCard({ label, value, sub, color = 'text-[#1d1d1f]', bg = 'bg-white' }: any) {
  return (
    <div className={`${bg} rounded-2xl border border-[#d2d2d7] shadow-sm p-5`}>
      <p className="text-xs font-medium text-[#6e6e73] uppercase tracking-wide mb-1">{label}</p>
      <p className={`text-3xl font-semibold tabular-nums ${color}`}>{value}</p>
      {sub && <p className="text-xs text-[#a1a1a6] mt-1">{sub}</p>}
    </div>
  );
}

export default function NetworkPage() {
  const [relationships, setRelationships] = useState<Relationship[]>([]);
  const [nodes, setNodes] = useState<NodeInfo[]>([]);
  const [selectedNode, setSelectedNode] = useState<NodeInfo | null>(null);
  const [filterFraud, setFilterFraud] = useState(false);
  const [filterType, setFilterType] = useState('ALL');
  const [loading, setLoading] = useState(true);

  const fetchGraph = useCallback(async () => {
    try {
      const res = await fetch('/api/v1/graph/relationships?limit=200');
      if (res.ok) {
        const json = await res.json();
        const data: Relationship[] = json.data ?? [];
        setRelationships(data);
        setNodes(buildGraph(data));
      }
    } catch {}
    setLoading(false);
  }, []);

  useEffect(() => { fetchGraph(); }, [fetchGraph]);

  const fraudRels = relationships.filter(r => r.fraud_linked);
  const userNodes = nodes.filter(n => n.type === 'USER');
  const deviceNodes = nodes.filter(n => n.type === 'DEVICE');
  const ipNodes = nodes.filter(n => n.type === 'IP');

  const filteredNodes = nodes.filter(n => {
    if (filterFraud && !n.fraudLinked) return false;
    if (filterType !== 'ALL' && n.type !== filterType) return false;
    return true;
  });

  return (
    <div className="min-h-screen bg-[#f5f5f7]">
      <Navbar />
      <main className="max-w-7xl mx-auto px-6 py-8">
        <div className="flex items-start justify-between mb-8">
          <div>
            <h1 className="text-2xl font-semibold text-[#1d1d1f] tracking-tight">Fraud Network Graph</h1>
            <p className="text-sm text-[#6e6e73] mt-1">Entity relationships · Device intelligence · IP clustering · Fraud propagation</p>
          </div>
          <button onClick={fetchGraph} className="px-4 py-2 bg-white border border-[#d2d2d7] text-sm font-medium text-[#1d1d1f] rounded-xl hover:bg-[#f5f5f7] transition-colors">
            ↺ Refresh
          </button>
        </div>

        {/* Stats Row */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
          <StatCard label="Total Entities" value={nodes.length} sub="Unique graph nodes" />
          <StatCard label="Users" value={userNodes.length} sub="Monitored accounts" color="text-[#0071e3]" bg="bg-blue-50" />
          <StatCard label="Devices" value={deviceNodes.length} sub="Tracked device fingerprints" color="text-purple-600" bg="bg-purple-50" />
          <StatCard label="Fraud-Linked" value={fraudRels.length} sub="High-risk edge count" color="text-red-600" bg="bg-red-50" />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {/* Node List */}
          <div className="md:col-span-1 bg-white rounded-2xl border border-[#d2d2d7] shadow-sm">
            <div className="px-5 py-4 border-b border-[#e8e8ed]">
              <h2 className="font-semibold text-sm text-[#1d1d1f]">Entity Nodes</h2>
              <p className="text-xs text-[#6e6e73] mt-0.5">{filteredNodes.length} entities</p>
            </div>
            <div className="px-4 py-3 border-b border-[#e8e8ed] space-y-2">
              <div className="flex gap-1 flex-wrap">
                {['ALL', 'USER', 'DEVICE', 'IP', 'MERCHANT'].map(t => (
                  <button key={t} onClick={() => setFilterType(t)}
                    className={`px-2.5 py-1 text-xs rounded-lg font-medium transition-colors ${filterType === t ? 'bg-[#0071e3] text-white' : 'bg-[#f5f5f7] text-[#6e6e73] hover:bg-[#e8e8ed]'}`}>
                    {t}
                  </button>
                ))}
              </div>
              <button onClick={() => setFilterFraud(!filterFraud)}
                className={`w-full px-3 py-1.5 text-xs font-medium rounded-lg transition-colors text-left ${filterFraud ? 'bg-red-50 text-red-600 border border-red-200' : 'bg-[#f5f5f7] text-[#6e6e73] border border-[#e8e8ed]'}`}>
                {filterFraud ? '✓' : '○'} Fraud-Linked Only
              </button>
            </div>

            <div className="overflow-y-auto max-h-[500px]">
              {loading ? (
                <div className="p-6 text-center text-sm text-[#a1a1a6]">Loading graph data…</div>
              ) : filteredNodes.length === 0 ? (
                <div className="p-6 text-center text-sm text-[#a1a1a6]">No entities match filters.<br/>Start a scenario to generate network data.</div>
              ) : filteredNodes.map(node => (
                <div key={node.id}
                  onClick={() => setSelectedNode(selectedNode?.id === node.id ? null : node)}
                  className={`px-4 py-3.5 border-b border-[#f5f5f7] cursor-pointer hover:bg-[#f5f5f7] transition-colors ${selectedNode?.id === node.id ? 'bg-blue-50 border-l-2 border-l-[#0071e3]' : ''}`}>
                  <div className="flex items-center justify-between mb-1">
                    <span className="font-mono text-xs font-semibold text-[#1d1d1f]">{node.id}</span>
                    <EntityBadge type={node.type} />
                  </div>
                  <div className="flex items-center gap-3">
                    <span className="text-xs text-[#6e6e73]">{node.connections} connections</span>
                    {node.fraudLinked && (
                      <span className="text-xs px-1.5 py-0.5 rounded bg-red-100 text-red-600 border border-red-200">Fraud</span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Detail Panel */}
          <div className="md:col-span-2 space-y-4">
            {selectedNode ? (
              <>
                <div className="bg-white rounded-2xl border border-[#d2d2d7] shadow-sm p-6">
                  <div className="flex items-start justify-between mb-5">
                    <div>
                      <h2 className="font-semibold text-lg text-[#1d1d1f] font-mono">{selectedNode.id}</h2>
                      <div className="flex items-center gap-2 mt-1.5">
                        <EntityBadge type={selectedNode.type} />
                        {selectedNode.fraudLinked && (
                          <span className="px-2 py-0.5 rounded text-xs font-medium bg-red-50 text-red-600 border border-red-200">⚠ Fraud-Linked</span>
                        )}
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-3xl font-semibold text-[#1d1d1f]">{selectedNode.connections}</p>
                      <p className="text-xs text-[#6e6e73]">total connections</p>
                    </div>
                  </div>

                  {/* Per-Type Stats */}
                  <div className="grid grid-cols-3 gap-3 mb-4">
                    {[
                      { label: 'Device Links', count: selectedNode.edges.filter(e => e.relationship_type.includes('DEVICE')).length, color: 'text-purple-600' },
                      { label: 'IP Links', count: selectedNode.edges.filter(e => e.relationship_type.includes('IP')).length, color: 'text-amber-600' },
                      { label: 'Fraud Links', count: selectedNode.edges.filter(e => e.fraud_linked).length, color: 'text-red-600' },
                    ].map(stat => (
                      <div key={stat.label} className="p-3 rounded-xl bg-[#f5f5f7] border border-[#e8e8ed] text-center">
                        <p className={`text-xl font-semibold ${stat.color}`}>{stat.count}</p>
                        <p className="text-xs text-[#6e6e73] mt-0.5">{stat.label}</p>
                      </div>
                    ))}
                  </div>

                  <h3 className="text-xs font-medium text-[#6e6e73] uppercase tracking-wide mb-3">Connected Edges</h3>
                  <div className="space-y-2 max-h-[320px] overflow-y-auto">
                    {selectedNode.edges.map(r => (
                      <div key={r.relationship_id}
                        className={`flex items-center justify-between px-3.5 py-2.5 rounded-xl border ${r.fraud_linked ? 'bg-red-50 border-red-200' : 'bg-[#f5f5f7] border-[#e8e8ed]'}`}>
                        <div className="flex items-center gap-2 min-w-0">
                          <span className="font-mono text-xs text-[#0071e3] font-medium truncate">{r.source_id}</span>
                          <span className="text-xs text-[#a1a1a6] flex-shrink-0">→</span>
                          <span className="font-mono text-xs text-purple-600 font-medium truncate">{r.target_id}</span>
                        </div>
                        <div className="flex items-center gap-2 flex-shrink-0 ml-3">
                          <span className="text-xs text-[#6e6e73]">{r.relationship_type.replace(/_/g, ' ')}</span>
                          {r.fraud_linked && <span className="text-xs px-1.5 py-0.5 rounded bg-red-100 text-red-600 border border-red-200">Fraud</span>}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </>
            ) : (
              <div className="bg-white rounded-2xl border border-[#d2d2d7] shadow-sm">
                {/* Full Edge Table */}
                <div className="px-6 py-4 border-b border-[#e8e8ed]">
                  <h2 className="font-semibold text-sm text-[#1d1d1f]">All Network Edges</h2>
                  <p className="text-xs text-[#6e6e73] mt-0.5">Select a node from the left panel to see detailed analysis</p>
                </div>
                <div className="overflow-x-auto">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="bg-[#f5f5f7] border-b border-[#e8e8ed]">
                        {['Source', 'Relationship', 'Target', 'Weight', 'Status'].map(h => (
                          <th key={h} className="px-4 py-3 text-left text-xs font-medium text-[#6e6e73] uppercase tracking-wide">{h}</th>
                        ))}
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-[#f5f5f7]">
                      {loading ? (
                        <tr><td colSpan={5} className="px-4 py-8 text-center text-sm text-[#a1a1a6]">Loading…</td></tr>
                      ) : relationships.length === 0 ? (
                        <tr><td colSpan={5} className="px-4 py-12 text-center text-sm text-[#a1a1a6]">No relationships yet. Run a device_farm or ip_abuse scenario.</td></tr>
                      ) : relationships.map(r => (
                        <tr key={r.relationship_id} className={`hover:bg-[#f5f5f7] transition-colors ${r.fraud_linked ? 'bg-red-50/30' : ''}`}>
                          <td className="px-4 py-3 font-mono text-xs text-[#0071e3] font-medium">{r.source_id}</td>
                          <td className="px-4 py-3 text-xs text-[#6e6e73]">{r.relationship_type.replace(/_/g, ' ')}</td>
                          <td className="px-4 py-3 font-mono text-xs text-purple-600 font-medium">{r.target_id}</td>
                          <td className="px-4 py-3 text-xs text-[#6e6e73]">{r.weight?.toFixed(2) ?? '—'}</td>
                          <td className="px-4 py-3">
                            {r.fraud_linked
                              ? <span className="px-2 py-0.5 rounded-full text-xs font-medium bg-red-50 text-red-600 border border-red-200">Fraud Linked</span>
                              : <span className="px-2 py-0.5 rounded-full text-xs font-medium bg-green-50 text-green-600 border border-green-200">Clean</span>
                            }
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
