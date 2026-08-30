'use client';

import React, { useState, useEffect } from 'react';

interface Rel {
  relationship_id: string;
  source_type: string;
  source_id: string;
  target_type: string;
  target_id: string;
  relationship_type: string;
  fraud_linked: boolean;
}

export default function GraphVisualizer() {
  const [relationships, setRelationships] = useState<Rel[]>([]);

  useEffect(() => {
    const fetchGraph = async () => {
      try {
        const res = await fetch('/api/v1/graph/relationships?limit=20');
        if (res.ok) {
          const json = await res.json();
          if (json.data) {
            setRelationships(json.data);
          }
        }
      } catch (e) {
        // ignore
      }
    };
    fetchGraph();
  }, []);

  return (
    <div className="p-5 rounded-xl bg-slate-900/60 border border-slate-800 backdrop-blur mb-6">
      <h2 className="font-bold text-sm uppercase tracking-wider text-slate-300 mb-4 flex justify-between items-center">
        <span>Fraud Network Graph Visualizer</span>
        <span className="text-xs font-mono font-normal text-slate-500">Relational Graph Engine</span>
      </h2>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="p-4 rounded-lg bg-slate-950 border border-slate-800 font-mono text-xs max-h-60 overflow-y-auto">
          <p className="text-[11px] text-slate-500 uppercase mb-3">Graph Edges & Fraud Links</p>
          {relationships.length === 0 ? (
            <p className="text-slate-500">No graph relationships extracted yet.</p>
          ) : (
            relationships.map((r) => (
              <div key={r.relationship_id} className="p-2 mb-2 rounded bg-slate-900 border border-slate-800 flex justify-between items-center">
                <div>
                  <span className="text-blue-400 font-bold">{r.source_id}</span>
                  <span className="text-slate-500 mx-2">→ ({r.relationship_type}) →</span>
                  <span className="text-indigo-400 font-bold">{r.target_id}</span>
                </div>
                {r.fraud_linked ? (
                  <span className="px-2 py-0.5 rounded text-[10px] bg-rose-500/20 text-rose-400 border border-rose-500/30">FRAUD LINKED</span>
                ) : (
                  <span className="px-2 py-0.5 rounded text-[10px] bg-slate-800 text-slate-400">CLEAN</span>
                )}
              </div>
            ))
          )}
        </div>

        <div className="p-4 rounded-lg bg-slate-950 border border-slate-800 font-mono text-xs flex flex-col justify-center">
          <h3 className="text-slate-300 font-bold mb-2">Graph Intelligence Insights</h3>
          <p className="text-slate-400 text-xs leading-relaxed">
            The relational graph engine tracks entities (users, devices, IPs, merchants) and detects shared high-risk nodes across accounts to prevent coordinated device-farm and IP-abuse carding attacks.
          </p>
        </div>
      </div>
    </div>
  );
}
