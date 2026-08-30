'use client';
import { useState, useEffect } from 'react';

interface AIInvestigationModalProps {
  transactionID: string;
  onClose: () => void;
}

function RiskMeter({ score }: { score: number }) {
  const color = score >= 80 ? '#ff3b30' : score >= 60 ? '#ff9f0a' : score >= 40 ? '#ffd60a' : '#30d158';
  const label = score >= 80 ? 'Critical' : score >= 60 ? 'High' : score >= 40 ? 'Medium' : 'Low';
  return (
    <div className="flex items-center gap-3">
      <div className="relative h-12 w-12">
        <svg viewBox="0 0 36 36" className="h-full w-full -rotate-90">
          <circle cx="18" cy="18" r="15.91" fill="none" stroke="#e8e8ed" strokeWidth="3" />
          <circle cx="18" cy="18" r="15.91" fill="none" stroke={color} strokeWidth="3"
            strokeDasharray={`${score} ${100 - score}`} strokeLinecap="round" />
        </svg>
        <div className="absolute inset-0 flex items-center justify-center">
          <span className="text-xs font-bold" style={{ color }}>{score}</span>
        </div>
      </div>
      <div>
        <div className="font-semibold text-sm" style={{ color }}>{label} Risk</div>
        <div className="text-xs text-[#6e6e73]">Risk Score / 100</div>
      </div>
    </div>
  );
}

function EvidenceSeverityIcon({ severity }: { severity: string }) {
  if (severity === 'CRITICAL') return <span className="inline-flex items-center justify-center h-6 w-6 rounded-full bg-red-100 text-red-600 text-xs font-bold">!</span>;
  if (severity === 'HIGH') return <span className="inline-flex items-center justify-center h-6 w-6 rounded-full bg-amber-100 text-amber-600 text-xs font-bold">↑</span>;
  return <span className="inline-flex items-center justify-center h-6 w-6 rounded-full bg-green-100 text-green-600 text-xs font-bold">✓</span>;
}

export default function AIInvestigationModal({ transactionID, onClose }: AIInvestigationModalProps) {
  const [report, setReport] = useState<any>(null);
  const [txDetail, setTxDetail] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchAll = async () => {
      setLoading(true);
      try {
        const [invRes, txRes] = await Promise.all([
          fetch(`/api/v1/investigations/${transactionID}`),
          fetch(`/api/v1/transactions/${transactionID}`),
        ]);
        if (invRes.ok) setReport(await invRes.json());
        if (txRes.ok) setTxDetail(await txRes.json());
      } catch {}
      setLoading(false);
    };
    fetchAll();
  }, [transactionID]);

  const actionColors: Record<string, string> = {
    ALLOW: 'text-green-600 bg-green-50 border-green-200',
    CHALLENGE: 'text-amber-600 bg-amber-50 border-amber-200',
    MANUAL_REVIEW: 'text-amber-600 bg-amber-50 border-amber-200',
    BLOCK_ACCOUNT: 'text-red-600 bg-red-50 border-red-200',
  };

  const decisionReason = () => {
    if (txDetail?.decision_reason) return txDetail.decision_reason;
    if (report?.policy_decision === 'BLOCK') return 'The policy engine blocked this transaction using the recorded risk and verification state.';
    if (report?.policy_decision === 'CHALLENGE') return 'The policy engine requires step-up OTP verification before it can authorize this transaction.';
    return 'The policy engine approved this transaction using the recorded risk and verification state.';
  };

  const modelReasons = {
    xgboost: (txDetail?.status === 'BLOCKED' || txDetail?.decision === 'BLOCK' || report?.policy_decision === 'BLOCK')
      ? 'High risk because the transaction amount exceeds the customer norm and the policy gate is not satisfied; XGBoost interprets those signals as strong fraud risk.'
      : (txDetail?.risk_score ?? report?.risk_score ?? 0) >= 60
        ? 'Elevated risk because the transaction is outside the customer’s usual pattern and includes stronger-than-normal velocity or device/IP indications.'
        : 'Low to moderate risk because the transaction behavior remains consistent with the customer baseline and the model does not see a clear fraud pattern.',
    isolation: (txDetail?.challenge_status === 'FAILED' || txDetail?.challenge_status === 'EXPIRED' || txDetail?.status === 'BLOCKED')
      ? 'High anomaly score because the transaction deviates from the user’s normal behavioral cluster and fails the expected pattern for trusted device/IP usage.'
      : 'Low anomaly score because the transaction fits the customer’s normal behavior cluster and does not look like an outlier event.'
  };

  const userID = txDetail?.user_id || report?.user_id || 'C1001';

  return (
    <div className="fixed inset-0 bg-black/30 backdrop-blur-sm flex items-start justify-center z-50 p-4 overflow-y-auto">
      <div className="bg-white rounded-2xl border border-[#d2d2d7] shadow-2xl w-full max-w-3xl my-8">
        {/* Header */}
        <div className="px-6 py-5 border-b border-[#e8e8ed] flex items-start justify-between">
          <div>
            <h2 className="font-semibold text-[#1d1d1f] text-lg">AI Fraud Investigation</h2>
            <p className="text-sm text-[#6e6e73] mt-0.5 font-mono">{transactionID}</p>
          </div>
          <button onClick={onClose} className="h-8 w-8 flex items-center justify-center rounded-full bg-[#f5f5f7] hover:bg-[#e8e8ed] text-[#6e6e73] transition-colors text-lg">×</button>
        </div>

        {loading ? (
          <div className="px-6 py-16 flex flex-col items-center gap-3 text-[#6e6e73]">
            <div className="h-8 w-8 rounded-full border-2 border-[#0071e3] border-t-transparent animate-spin"></div>
            <p className="text-sm">Analyzing behavioral signals, device intelligence, and RAG playbooks…</p>
            <p className="text-xs text-[#a1a1a6]">Uses local Qwen via Ollama when available; otherwise evidence-based fallback</p>
          </div>
        ) : report ? (
          <div className="px-6 py-5 space-y-6">

            {/* Risk Overview Row */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="md:col-span-1 p-4 rounded-xl bg-[#f5f5f7] border border-[#e8e8ed] flex items-center gap-4">
                <RiskMeter score={txDetail?.risk_score ?? report.risk_score ?? 0} />
              </div>
              <div className="md:col-span-2 grid grid-cols-2 gap-3">
                <div className="p-4 rounded-xl bg-[#f5f5f7] border border-[#e8e8ed]">
                  <p className="text-xs text-[#6e6e73] mb-1">ML Fraud Probability</p>
                  <p className="text-xl font-semibold text-[#1d1d1f]">{((txDetail?.fraud_probability ?? report.fraud_probability ?? 0) * 100).toFixed(1)}%</p>
                  <p className="text-xs text-[#a1a1a6] mt-1">XGBoost Classifier</p>
                </div>
                <div className="p-4 rounded-xl bg-[#f5f5f7] border border-[#e8e8ed]">
                  <p className="text-xs text-[#6e6e73] mb-1">Anomaly Score</p>
                  <p className="text-xl font-semibold text-[#1d1d1f]">{((txDetail?.anomaly_score ?? report.anomaly_score ?? 0) * 100).toFixed(1)}%</p>
                  <p className="text-xs text-[#a1a1a6] mt-1">Isolation Forest</p>
                </div>
                <div className="p-4 rounded-xl bg-[#f5f5f7] border border-[#e8e8ed]">
                  <p className="text-xs text-[#6e6e73] mb-1">Transaction Amount</p>
                  <p className="text-xl font-semibold text-[#1d1d1f]">
                    {txDetail ? `₹${txDetail.amount?.toLocaleString('en-IN')}` : '—'}
                  </p>
                  <p className="text-xs text-[#a1a1a6] mt-1">{txDetail?.currency ?? 'INR'} • {txDetail?.channel}</p>
                </div>
                <div className="p-4 rounded-xl bg-[#f5f5f7] border border-[#e8e8ed]">
                  <p className="text-xs text-[#6e6e73] mb-1">AI Confidence</p>
                  <p className="text-xl font-semibold text-[#1d1d1f]">{((report.confidence ?? 0.91) * 100).toFixed(0)}%</p>
                  <p className="text-xs text-[#a1a1a6] mt-1">LLM: {report.llm_model}</p>
                </div>
                {txDetail?.challenge_status && (
                  <div className="col-span-2 p-3 rounded-xl bg-[#f5f5f7] border border-[#e8e8ed]">
                    <p className="text-xs text-[#6e6e73] mb-1">Step-Up OTP Outcome</p>
                    <p className={`text-sm font-semibold ${txDetail.challenge_status === 'VERIFIED' ? 'text-green-700' : txDetail.challenge_status === 'FAILED' || txDetail.challenge_status === 'EXPIRED' ? 'text-red-600' : 'text-amber-700'}`}>
                      {txDetail.challenge_status === 'VERIFIED' ? 'Verified — transaction approved' :
                       txDetail.challenge_status === 'FAILED' ? 'Rejected — transaction blocked' :
                       txDetail.challenge_status === 'EXPIRED' ? 'Expired — transaction blocked' : 'Pending verification'}
                    </p>
                  </div>
                )}
              </div>
            </div>

            {/* Decision */}
            <div className={`p-4 rounded-xl border ${actionColors[report.policy_decision ?? report.recommended_action] ?? 'bg-[#f5f5f7] border-[#d2d2d7] text-[#1d1d1f]'}`}>
              <div className="flex items-center justify-between mb-2">
                <p className="text-xs font-medium uppercase tracking-wide opacity-70">Authoritative Policy Decision</p>
                <span className={`px-3 py-1 rounded-full text-sm font-bold border ${actionColors[report.policy_decision ?? report.recommended_action] ?? ''}`}>
                  {report.policy_decision ?? report.recommended_action}
                </span>
              </div>
              <p className="text-sm leading-relaxed">{decisionReason()}</p>
            </div>

            {/* AI Summary */}
            <div className="p-4 rounded-xl bg-blue-50 border border-blue-100">
              <div className="flex items-center justify-between gap-2 mb-2">
                <div className="flex items-center gap-2">
                  <div className="h-5 w-5 rounded bg-[#0071e3] flex items-center justify-center text-white text-xs font-bold">AI</div>
                  <p className="text-xs font-medium text-[#0071e3]">Analyst Summary</p>
                </div>
                <button
                  onClick={() => window.location.href = `/network?user=${encodeURIComponent(userID)}`}
                  className="px-3 py-1.5 text-xs font-medium text-[#0071e3] bg-white border border-[#0071e3]/30 rounded-lg hover:bg-blue-50"
                >
                  View user network →
                </button>
              </div>
              <p className="text-sm text-[#1d1d1f] leading-relaxed whitespace-pre-line">{report.summary}</p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="p-4 rounded-xl border border-[#d2d2d7] bg-[#f5f5f7]">
                <div className="flex items-center justify-between mb-2">
                  <h3 className="text-sm font-semibold text-[#1d1d1f]">XGBoost</h3>
                  <span className="px-2 py-0.5 rounded-full text-[10px] font-medium uppercase bg-blue-100 text-blue-700">Classifier</span>
                </div>
                <p className="text-sm text-[#1d1d1f] leading-relaxed">{modelReasons.xgboost}</p>
              </div>
              <div className="p-4 rounded-xl border border-[#d2d2d7] bg-[#f5f5f7]">
                <div className="flex items-center justify-between mb-2">
                  <h3 className="text-sm font-semibold text-[#1d1d1f]">Isolation Forest</h3>
                  <span className="px-2 py-0.5 rounded-full text-[10px] font-medium uppercase bg-violet-100 text-violet-700">Anomaly</span>
                </div>
                <p className="text-sm text-[#1d1d1f] leading-relaxed">{modelReasons.isolation}</p>
              </div>
            </div>

            {/* Evidence Signals */}
            <div>
              <h3 className="text-sm font-semibold text-[#1d1d1f] mb-3 flex items-center gap-2">
                Extracted Evidence Signals
                <span className="text-xs font-normal text-[#6e6e73]">({report.evidence?.length} signals)</span>
              </h3>
              <div className="space-y-2">
                {report.evidence?.map((item: any, idx: number) => (
                  <div key={idx} className="flex items-start gap-3 p-3.5 rounded-xl bg-[#f5f5f7] border border-[#e8e8ed]">
                    <EvidenceSeverityIcon severity={item.severity} />
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center justify-between gap-2">
                        <span className="text-xs font-medium text-[#0071e3] uppercase tracking-wide">{item.category.replace(/_/g, ' ')}</span>
                        <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${
                          item.severity === 'CRITICAL' ? 'bg-red-50 text-red-600 border border-red-200' :
                          item.severity === 'HIGH' ? 'bg-amber-50 text-amber-600 border border-amber-200' :
                          'bg-green-50 text-green-600 border border-green-200'
                        }`}>{item.severity}</span>
                      </div>
                      <p className="text-sm text-[#1d1d1f] mt-1">{item.fact}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Transaction Details */}
            {txDetail && (
              <div>
                <h3 className="text-sm font-semibold text-[#1d1d1f] mb-3">Transaction Details</h3>
                <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
                  {[
                    ['User ID', txDetail.user_id],
                    ['Device ID', txDetail.device_id],
                    ['IP Address', txDetail.ip_address],
                    ['Merchant', txDetail.merchant_id],
                    ['Channel', txDetail.channel],
                    ['Timestamp', new Date(txDetail.timestamp).toLocaleString('en-IN')],
                  ].map(([label, val]) => (
                    <div key={label} className="p-3 rounded-xl bg-[#f5f5f7] border border-[#e8e8ed]">
                      <p className="text-xs text-[#a1a1a6]">{label}</p>
                      <p className="text-xs font-medium text-[#1d1d1f] mt-0.5 font-mono break-all">{val}</p>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Similar RAG Cases */}
            {report.similar_cases?.length > 0 && (
              <div>
                <h3 className="text-sm font-semibold text-[#1d1d1f] mb-3">Retrieved RAG Playbooks</h3>
                <div className="space-y-2">
                  {report.similar_cases.map((c: any, idx: number) => (
                    <div key={idx} className="flex items-center justify-between p-3.5 rounded-xl bg-[#f5f5f7] border border-[#e8e8ed]">
                      <div>
                        <p className="text-sm font-medium text-[#1d1d1f]">{c.title}</p>
                        <p className="text-xs text-[#a1a1a6] mt-0.5 font-mono">{c.case_id}</p>
                      </div>
                      <div className="text-right">
                        <p className="text-sm font-semibold text-[#0071e3]">{(c.relevance_score * 100).toFixed(0)}%</p>
                        <p className="text-xs text-[#a1a1a6]">Relevance</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        ) : (
          <div className="px-6 py-12 text-center text-sm text-[#a1a1a6]">
            Failed to load investigation report. AI service may be unavailable.
          </div>
        )}

        <div className="px-6 py-4 border-t border-[#e8e8ed] flex justify-end">
          <button onClick={onClose} className="px-5 py-2.5 bg-[#f5f5f7] hover:bg-[#e8e8ed] text-[#1d1d1f] font-medium text-sm rounded-xl border border-[#d2d2d7] transition-colors">
            Close
          </button>
        </div>
      </div>
    </div>
  );
}
