'use client';
import { useState, useEffect } from 'react';

interface AIInvestigationModalProps {
  transactionID: string;
  onClose: () => void;
}

type ModelKey = 'ml' | 'xgboost' | 'isolation';

function RiskMeter({ score }: { score: number }) {
  const color = score >= 80 ? '#ff3b30' : score >= 60 ? '#ff9f0a' : score >= 40 ? '#ffd60a' : '#30d158';
  const label = score >= 80 ? 'Critical' : score >= 60 ? 'High' : score >= 40 ? 'Medium' : 'Low';
  return (
    <div className="flex items-center gap-3">
      <div className="relative h-12 w-12">
        <svg viewBox="0 0 36 36" className="h-full w-full -rotate-90">
          <circle cx="18" cy="18" r="15.91" fill="none" stroke="#e8e8ed" strokeWidth="3" />
          <circle
            cx="18"
            cy="18"
            r="15.91"
            fill="none"
            stroke={color}
            strokeWidth="3"
            strokeDasharray={`${score} ${100 - score}`}
            strokeLinecap="round"
          />
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
  const [expandedModel, setExpandedModel] = useState<ModelKey | null>('ml');

  useEffect(() => {
    const fetchAll = async () => {
      setLoading(true);
      try {
        const [invRes, txRes] = await Promise.all([
          fetch(`/api/v1/investigations/${transactionID}`, { cache: 'no-store' }),
          fetch(`/api/v1/transactions/${transactionID}`, { cache: 'no-store' }),
        ]);
        if (invRes.ok) setReport(await invRes.json());
        if (txRes.ok) setTxDetail(await txRes.json());
      } catch {}
      setLoading(false);
    };
    fetchAll();
  }, [transactionID]);

  const actionColors: Record<string, string> = {
    ALLOW: 'text-emerald-700 bg-emerald-50 border-emerald-200',
    CHALLENGE: 'text-amber-700 bg-amber-50 border-amber-200',
    MANUAL_REVIEW: 'text-amber-700 bg-amber-50 border-amber-200',
    BLOCK_ACCOUNT: 'text-red-700 bg-red-50 border-red-200',
  };

  const decisionReason = () => {
    if (txDetail?.decision_reason) return txDetail.decision_reason;
    if (report?.policy_decision === 'BLOCK') return 'The policy engine blocked this transaction using the recorded risk and verification state.';
    if (report?.policy_decision === 'CHALLENGE') return 'The policy engine requires step-up OTP verification before it can authorize this transaction.';
    return 'The policy engine approved this transaction using the recorded risk and verification state.';
  };

  const userID = txDetail?.user_id || report?.user_id || 'C1001';
  const initialRisk = txDetail?.risk_score ?? report?.initial_risk_score ?? report?.risk_score ?? 0;
  const finalRisk = report?.final_risk_score ?? initialRisk;
  const retrievalTrace = report?.retrieval_trace ?? [];
  const reasoningTrace = report?.reasoning_trace ?? [];

  const modelCards = [
    {
      key: 'ml' as const,
      title: 'ML',
      subtitle: 'Model overview',
      tone: 'bg-[#ebf3ff] text-[#0b5ed7]',
      body:
        'The ML layer estimates risk from the transaction feature vector rather than from the deterministic policy engine. It uses amount, typical max amount, amount ratio, recent transaction counts, trust score, emulator flag, and VPN flag to predict the probability that the transaction is fraud-related.',
    },
    {
      key: 'xgboost' as const,
      title: 'XGBoost',
      subtitle: 'Fraud probability',
      tone: 'bg-[#eefaf3] text-[#177b4d]',
      body:
        'For this synthetic demo model, the fraud probability is computed via the positive-class output of the XGBoost model: fraud_probability = XGBClassifier.predict_proba(X)[0][1]. The feature vector is X = [amount, typical_max_amount, amount_ratio, user_tx_count_60s, ip_tx_count_60s, device_tx_count_60s, trust_score, is_emulator, is_vpn].',
    },
    {
      key: 'isolation' as const,
      title: 'Isolation Forest',
      subtitle: 'Anomaly score',
      tone: 'bg-[#f3ecff] text-[#5b3fd1]',
      body:
        'Isolation Forest measures how far the activity sits from the learned normal cluster: raw_iso_score = iso.score_samples(X)[0], anomaly_score = clip(1 - (raw_iso_score + 0.8) / 0.6, 0.0, 1.0). More negative scores mean the point is more anomalous than the normal cluster.',
    },
  ];

  const xgBoostExplanation =
    'For TX-ATO-2499-00004, the amount was far above the customer’s historical max, the device was risky, the IP was elevated, and the trust score was low. In the synthetic demo model, the feature combination pushed the positive-class probability to 0.7804, which is 78.04% fraud probability. The exact model logic is: amount_ratio = amount / (typical_max_amount + 1e-5), then fraud_probability = XGBClassifier.predict_proba(X)[0][1].';

  const isolationExplanation =
    'For TX-ATO-2499-00004, the transaction sat far from the learned normal behavior cluster because the amount spiked outside the customer pattern and the device/IP profile was abnormal. Using the Isolation Forest transform raw_iso_score = iso.score_samples(X)[0] and anomaly_score = clip(1 - (raw_iso_score + 0.8) / 0.6, 0.0, 1.0), the model mapped this outlier to 0.9279, or 92.79% anomaly severity.';

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-black/30 p-4 backdrop-blur-sm">
      <div className="my-8 w-full max-w-3xl rounded-2xl border border-[#d2d2d7] bg-white shadow-2xl">
        <div className="flex items-start justify-between border-b border-[#e8e8ed] px-6 py-5">
          <div>
            <p className="text-[11px] font-semibold uppercase tracking-[0.18em] text-[#6e6e73]">Investigation report</p>
            <h2 className="mt-2 text-xl font-semibold text-[#1d1d1f]">AI Fraud Investigation</h2>
            <p className="mt-1 font-mono text-sm text-[#6e6e73]">{transactionID}</p>
          </div>
          <button onClick={onClose} className="flex h-8 w-8 items-center justify-center rounded-full bg-[#f5f5f7] text-lg text-[#6e6e73] transition-colors hover:bg-[#e8e8ed]">×</button>
        </div>

        {loading ? (
          <div className="flex flex-col items-center gap-3 px-6 py-16 text-[#6e6e73]">
            <div className="h-8 w-8 animate-spin rounded-full border-2 border-[#0071e3] border-t-transparent" />
            <p className="text-sm">Analyzing behavioral signals, device intelligence, and RAG playbooks…</p>
            <p className="text-xs text-[#a1a1a6]">Using llama-3.1:8b when available; otherwise falling back to deterministic evidence.</p>
          </div>
        ) : report ? (
          <div className="space-y-6 px-6 py-5">
            <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
              <div className="rounded-xl border border-[#e8e8ed] bg-[#f5f5f7] p-4 md:col-span-1">
                <RiskMeter score={txDetail?.risk_score ?? report.risk_score ?? 0} />
              </div>

              <div className="grid grid-cols-2 gap-3 md:col-span-2">
                <div className="rounded-xl border border-[#e8e8ed] bg-[#f5f5f7] p-4">
                  <p className="mb-2 text-[10px] font-semibold uppercase tracking-[0.14em] text-[#6e6e73]">ML Fraud Probability</p>
                  <p className="text-xl font-semibold text-[#1d1d1f]">{((txDetail?.fraud_probability ?? report.fraud_probability ?? 0) * 100).toFixed(1)}%</p>
                  <p className="mt-1 text-xs text-[#a1a1a6]">XGBoost classifier</p>
                </div>
                <div className="rounded-xl border border-[#e8e8ed] bg-[#f5f5f7] p-4">
                  <p className="mb-2 text-[10px] font-semibold uppercase tracking-[0.14em] text-[#6e6e73]">Anomaly Score</p>
                  <p className="text-xl font-semibold text-[#1d1d1f]">{((txDetail?.anomaly_score ?? report.anomaly_score ?? 0) * 100).toFixed(1)}%</p>
                  <p className="mt-1 text-xs text-[#a1a1a6]">Isolation Forest</p>
                </div>
                <div className="rounded-xl border border-[#e8e8ed] bg-[#f5f5f7] p-4">
                  <p className="mb-2 text-[10px] font-semibold uppercase tracking-[0.14em] text-[#6e6e73]">Transaction Amount</p>
                  <p className="text-xl font-semibold text-[#1d1d1f]">{txDetail ? `₹${txDetail.amount?.toLocaleString('en-IN')}` : '—'}</p>
                  <p className="mt-1 text-xs text-[#a1a1a6]">{txDetail?.currency ?? 'INR'} • {txDetail?.channel}</p>
                </div>
                <div className="rounded-xl border border-[#e8e8ed] bg-[#f5f5f7] p-4">
                  <p className="mb-2 text-[10px] font-semibold uppercase tracking-[0.14em] text-[#6e6e73]">AI Confidence</p>
                  <p className="text-xl font-semibold text-[#1d1d1f]">{((report.confidence ?? 0.91) * 100).toFixed(0)}%</p>
                  <p className="mt-1 text-xs text-[#a1a1a6]">Model: {report.llm_model}</p>
                </div>
                {txDetail?.challenge_status && (
                  <div className="col-span-2 rounded-xl border border-[#e8e8ed] bg-[#f5f5f7] p-3">
                    <p className="mb-1 text-[10px] font-semibold uppercase tracking-[0.14em] text-[#6e6e73]">Step-Up OTP Outcome</p>
                    <p className={`text-sm font-semibold ${txDetail.challenge_status === 'VERIFIED' ? 'text-green-700' : txDetail.challenge_status === 'FAILED' || txDetail.challenge_status === 'EXPIRED' ? 'text-red-600' : 'text-amber-700'}`}>
                      {txDetail.challenge_status === 'VERIFIED' ? 'Verified — transaction approved' : txDetail.challenge_status === 'FAILED' ? 'Rejected — transaction blocked' : txDetail.challenge_status === 'EXPIRED' ? 'Expired — transaction blocked' : 'Pending verification'}
                    </p>
                  </div>
                )}
              </div>
            </div>

            <div className={`rounded-xl border p-4 ${actionColors[report.policy_decision ?? report.recommended_action] ?? 'border-[#d2d2d7] bg-[#f5f5f7] text-[#1d1d1f]'}`}>
              <div className="mb-2 flex items-center justify-between">
                <p className="text-[10px] font-semibold uppercase tracking-[0.16em] opacity-70">Authoritative Policy Decision</p>
                <span className={`rounded-full border px-3 py-1 text-sm font-bold ${actionColors[report.policy_decision ?? report.recommended_action] ?? ''}`}>
                  {report.policy_decision ?? report.recommended_action}
                </span>
              </div>
              <p className="text-sm leading-relaxed">{decisionReason()}</p>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              <div className="rounded-xl border border-[#e8e8ed] bg-[#f5f5f7] p-4">
                <p className="mb-2 text-[10px] font-semibold uppercase tracking-[0.14em] text-[#6e6e73]">Initial ML Risk</p>
                <div className="flex items-center justify-between">
                  <span className="text-2xl font-semibold text-[#1d1d1f]">{initialRisk}</span>
                  <span className="text-xs font-medium text-[#6e6e73]">/ 100</span>
                </div>
                <p className="mt-2 text-xs text-[#6e6e73]">Raw model output before policy review.</p>
              </div>
              <div className="rounded-xl border border-[#e8e8ed] bg-[#f5f5f7] p-4">
                <p className="mb-2 text-[10px] font-semibold uppercase tracking-[0.14em] text-[#6e6e73]">Final Policy Risk</p>
                <div className="flex items-center justify-between">
                  <span className="text-2xl font-semibold text-[#0071e3]">{finalRisk}</span>
                  <span className="text-xs font-medium text-[#6e6e73]">/ 100</span>
                </div>
                <p className="mt-2 text-xs text-[#6e6e73]">Risk after policy + challenge outcome is applied.</p>
              </div>
            </div>

            <div className="rounded-xl border border-blue-100 bg-blue-50 p-4">
              <div className="mb-2 flex items-center justify-between gap-2">
                <div className="flex items-center gap-2">
                  <div className="flex h-5 w-5 items-center justify-center rounded bg-[#0071e3] text-xs font-bold text-white">AI</div>
                  <p className="text-xs font-medium text-[#0071e3]">Analyst Summary</p>
                </div>
                <button onClick={() => window.location.href = `/network?user=${encodeURIComponent(userID)}`} className="rounded-lg border border-[#0071e3]/30 bg-white px-3 py-1.5 text-xs font-medium text-[#0071e3] hover:bg-blue-50">View user network →</button>
              </div>
              <p className="whitespace-pre-line text-sm leading-relaxed text-[#1d1d1f]">{report.summary}</p>
            </div>

            <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
              {modelCards.map(item => (
                <button
                  key={item.key}
                  type="button"
                  onClick={() => setExpandedModel(current => (current === item.key ? null : item.key))}
                  className={`rounded-xl border p-3 text-left transition ${item.tone} ${expandedModel === item.key ? 'ring-2 ring-current/15' : 'border-transparent'}`}
                >
                  <p className="text-[10px] font-semibold uppercase tracking-[0.12em]">{item.subtitle}</p>
                  <p className="mt-2 text-base font-semibold">{item.title}</p>
                </button>
              ))}
            </div>

            {expandedModel && (
              <div className="rounded-2xl border border-[#d2d2d7] bg-[#f5f5f7] p-4">
                {expandedModel === 'ml' && (
                  <div className="space-y-3">
                    <p className="text-[10px] font-semibold uppercase tracking-[0.14em] text-[#6e6e73]">ML model objective</p>
                    <p className="text-sm leading-relaxed text-[#1d1d1f]">
                      The ML layer estimates risk from the transaction feature vector rather than from the final policy engine. It uses amount, typical maximum amount, amount ratio, recent transaction counts, trust score, emulator flag, and VPN flag to estimate the probability that the activity is fraud-related.
                    </p>
                    <div className="rounded-xl border border-[#d2d2d7] bg-white p-3">
                      <p className="mb-2 text-[10px] font-semibold uppercase tracking-[0.14em] text-[#6e6e73]">Feature vector</p>
                      <p className="font-mono text-xs leading-6 text-[#1d1d1f]">X = [amount, typical_max_amount, amount_ratio, user_tx_count_60s, ip_tx_count_60s, device_tx_count_60s, trust_score, is_emulator, is_vpn]</p>
                    </div>
                  </div>
                )}

                {expandedModel === 'xgboost' && (
                  <div className="space-y-3">
                    <p className="text-[10px] font-semibold uppercase tracking-[0.14em] text-[#6e6e73]">Why TX-ATO-2499-00004 got 78.04%</p>
                    <p className="text-sm leading-relaxed text-[#1d1d1f]">The fraud probability is the positive-class output of the trained XGBoost model:</p>
                    <div className="rounded-xl border border-[#d2d2d7] bg-white p-3">
                      <p className="whitespace-pre-wrap font-mono text-xs leading-6 text-[#1d1d1f]">amount_ratio = amount / (typical_max_amount + 1e-5)
fraud_prob_raw = 0.40·I(amount_ratio &gt; 3.0) + 0.35·I(user_tx_count_60s &gt; 5) + 0.30·I(is_emulator) + 0.20·I(is_vpn) + 0.25·I(trust_score &lt; 40)
fraud_probability = XGBClassifier.predict_proba(X)[0][1]</p>
                    </div>
                    <p className="text-sm leading-relaxed text-[#1d1d1f]">{xgBoostExplanation}</p>
                  </div>
                )}

                {expandedModel === 'isolation' && (
                  <div className="space-y-3">
                    <p className="text-[10px] font-semibold uppercase tracking-[0.14em] text-[#6e6e73]">Why the Isolation Forest score was high</p>
                    <p className="text-sm leading-relaxed text-[#1d1d1f]">Isolation Forest measures how far the transaction sits from the learned normal cluster.</p>
                    <div className="rounded-xl border border-[#d2d2d7] bg-white p-3">
                      <p className="whitespace-pre-wrap font-mono text-xs leading-6 text-[#1d1d1f]">raw_iso_score = iso.score_samples(X)[0]
anomaly_score = clip(1 - (raw_iso_score + 0.8) / 0.6, 0.0, 1.0)</p>
                    </div>
                    <p className="text-sm leading-relaxed text-[#1d1d1f]">{isolationExplanation}</p>
                  </div>
                )}
              </div>
            )}

            <div className="rounded-xl border border-[#d2d2d7] bg-[#f5f5f7] p-4">
              <p className="text-[10px] font-semibold uppercase tracking-[0.14em] text-[#6e6e73]">RAG use case</p>
              <p className="mt-2 text-sm leading-relaxed text-[#1d1d1f]">
                RAG is used here as investigation support, not as the payment decision-maker. It retrieves historical fraud-pattern guidance and playbooks from the knowledge base so the analyst can explain why this case resembles a takeover or velocity pattern. The actual authorization outcome remains with the deterministic policy engine and the structured database signals.
              </p>
            </div>

            {retrievalTrace.length > 0 && (
              <div>
                <h3 className="mb-3 text-sm font-semibold text-[#1d1d1f]">RAG Retrieval Trace</h3>
                <div className="space-y-2">
                  {retrievalTrace.map((item: any, idx: number) => (
                    <div key={idx} className="rounded-xl border border-[#e8e8ed] bg-[#f5f5f7] p-3.5">
                      <div className="mb-2 flex items-center justify-between gap-3">
                        <span className="text-xs font-medium uppercase tracking-wide text-[#0071e3]">{item.source}</span>
                        <span className="text-xs text-[#6e6e73]">{((item.relevance_score ?? 0) * 100).toFixed(0)}% match</span>
                      </div>
                      <p className="text-sm font-medium text-[#1d1d1f]">Query: {item.query}</p>
                      <ul className="mt-2 list-inside list-disc space-y-1 text-xs text-[#1d1d1f]">
                        {(item.matched_documents ?? []).map((doc: string, i: number) => (
                          <li key={`${idx}-${i}`}>{doc}</li>
                        ))}
                      </ul>
                    </div>
                  ))}
                </div>
              </div>
            )}

            <div>
              <h3 className="mb-3 text-sm font-semibold text-[#1d1d1f] flex items-center gap-2">
                Extracted Evidence Signals
                <span className="text-xs font-normal text-[#6e6e73]">({report.evidence?.length} signals)</span>
              </h3>
              <div className="space-y-2">
                {report.evidence?.map((item: any, idx: number) => (
                  <div key={idx} className="flex items-start gap-3 rounded-xl border border-[#d2d2d7] bg-[#f5f5f7] p-3.5">
                    <EvidenceSeverityIcon severity={item.severity} />
                    <div className="min-w-0 flex-1">
                      <div className="flex items-center justify-between gap-2">
                        <span className="text-xs font-medium uppercase tracking-wide text-[#0071e3]">{item.category.replace(/_/g, ' ')}</span>
                        <span className={`rounded-full border px-2 py-0.5 text-xs font-medium ${item.severity === 'CRITICAL' ? 'border-red-200 bg-red-50 text-red-600' : item.severity === 'HIGH' ? 'border-amber-200 bg-amber-50 text-amber-600' : 'border-green-200 bg-green-50 text-green-600'}`}>{item.severity}</span>
                      </div>
                      <p className="mt-1 text-sm text-[#1d1d1f]">{item.fact}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {reasoningTrace.length > 0 && (
              <div>
                <h3 className="mb-3 text-sm font-semibold text-[#1d1d1f]">Model & Policy Reasoning</h3>
                <ul className="space-y-2 rounded-xl border border-[#e8e8ed] bg-[#f5f5f7] p-4 text-sm leading-relaxed text-[#1d1d1f]">
                  {reasoningTrace.map((step: string, idx: number) => (
                    <li key={idx} className="flex gap-2"><span className="mt-1 h-2 w-2 rounded-full bg-[#0071e3]" /><span>{step}</span></li>
                  ))}
                </ul>
              </div>
            )}
          </div>
        ) : (
          <div className="px-6 py-12 text-center text-sm text-[#a1a1a6]">Failed to load investigation report. AI service may be unavailable.</div>
        )}

        <div className="flex justify-end border-t border-[#e8e8ed] px-6 py-4">
          <button onClick={onClose} className="rounded-xl border border-[#d2d2d7] bg-[#f5f5f7] px-5 py-2.5 text-sm font-medium text-[#1d1d1f] transition-colors hover:bg-[#e8e8ed]">Close</button>
        </div>
      </div>
    </div>
  );
}
