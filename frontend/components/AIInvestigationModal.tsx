'use client';

import React, { useState, useEffect } from 'react';

interface AIInvestigationModalProps {
  transactionID: string;
  onClose: () => void;
}

export default function AIInvestigationModal({ transactionID, onClose }: AIInvestigationModalProps) {
  const [report, setReport] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchReport = async () => {
      setLoading(true);
      try {
        const res = await fetch(`/api/v1/investigations/${transactionID}`);
        if (res.ok) {
          const data = await res.json();
          setReport(data);
        }
      } catch (e) {
        // ignore
      } finally {
        setLoading(false);
      }
    };
    fetchReport();
  }, [transactionID]);

  return (
    <div className="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-2xl p-6 shadow-2xl max-h-[85vh] overflow-y-auto">
        <div className="flex justify-between items-center mb-4">
          <div>
            <h3 className="font-bold text-base text-white flex items-center gap-2">
              AI Fraud Analyst Investigation Report
              <span className="text-xs px-2 py-0.5 rounded bg-indigo-500/10 border border-indigo-500/30 text-indigo-400 font-mono">
                qwen2.5:7b-instruct
              </span>
            </h3>
            <p className="text-xs font-mono text-slate-400 mt-0.5">Transaction: {transactionID}</p>
          </div>
          <button onClick={onClose} className="text-slate-400 hover:text-white text-lg">×</button>
        </div>

        {loading ? (
          <div className="py-12 text-center text-slate-400 font-mono text-xs flex flex-col items-center gap-3">
            <div className="h-6 w-6 rounded-full border-2 border-blue-500 border-t-transparent animate-spin"></div>
            Analyzing behavioral context, device intelligence, and RAG playbooks...
          </div>
        ) : report ? (
          <div className="space-y-4">
            <div className="p-4 rounded-xl bg-slate-950 border border-slate-800">
              <div className="flex justify-between items-center mb-2">
                <span className="text-xs font-mono text-slate-400">Risk Assessment Level</span>
                <span className={`px-2.5 py-1 rounded text-xs font-mono font-bold ${
                  report.risk_level === 'CRITICAL' ? 'bg-rose-500/20 text-rose-400 border border-rose-500/40' :
                  report.risk_level === 'HIGH' ? 'bg-amber-500/20 text-amber-400 border border-amber-500/40' :
                  'bg-emerald-500/20 text-emerald-400 border border-emerald-500/40'
                }`}>
                  {report.risk_level}
                </span>
              </div>
              <p className="text-xs text-slate-200 leading-relaxed font-sans mt-2">{report.summary}</p>
            </div>

            <div>
              <h4 className="text-xs font-mono uppercase text-slate-400 mb-2">Extracted Evidence Signals</h4>
              <div className="space-y-2">
                {report.evidence?.map((item: any, idx: number) => (
                  <div key={idx} className="p-3 rounded-lg bg-slate-950 border border-slate-800 flex justify-between items-center">
                    <div>
                      <span className="text-[10px] font-mono text-blue-400 uppercase">{item.category}</span>
                      <p className="text-xs text-slate-300 mt-0.5">{item.fact}</p>
                    </div>
                    <span className="text-[10px] font-mono text-slate-400 px-2 py-0.5 rounded bg-slate-900">
                      {item.severity}
                    </span>
                  </div>
                ))}
              </div>
            </div>

            {report.similar_cases?.length > 0 && (
              <div>
                <h4 className="text-xs font-mono uppercase text-slate-400 mb-2">Retrieved RAG Playbooks & Attack Cases</h4>
                <div className="space-y-2">
                  {report.similar_cases.map((c: any, idx: number) => (
                    <div key={idx} className="p-3 rounded-lg bg-slate-950 border border-slate-800">
                      <p className="text-xs font-bold text-indigo-400">{c.title}</p>
                      <p className="text-[11px] font-mono text-slate-400 mt-0.5">Relevance score: {(c.relevance_score * 100).toFixed(0)}%</p>
                    </div>
                  ))}
                </div>
              </div>
            )}

            <div className="p-4 rounded-xl bg-blue-950/20 border border-blue-900/40 flex justify-between items-center">
              <div>
                <span className="text-[11px] font-mono text-slate-400 uppercase">Recommended Analyst Action</span>
                <p className="text-sm font-bold text-blue-400 font-mono mt-0.5">{report.recommended_action}</p>
              </div>
              <div className="text-right font-mono">
                <span className="text-[11px] text-slate-400">Confidence</span>
                <p className="text-sm font-bold text-white">{(report.confidence * 100).toFixed(0)}%</p>
              </div>
            </div>
          </div>
        ) : (
          <p className="text-xs text-slate-400 py-6 text-center">Failed to load investigation report.</p>
        )}
      </div>
    </div>
  );
}
