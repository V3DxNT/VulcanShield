'use client';
import { useState, useEffect } from 'react';

const SCENARIOS = [
  {
    value: 'normal',
    label: 'Normal Traffic',
    fn: 'Not an attack. Establishes the conversion baseline.',
    fp: 'ALLOW on in-profile spend so legitimate checkout is not blocked (false-positive control).',
  },
  {
    value: 'velocity_attack',
    label: 'Velocity Attack (Carding Burst)',
    fn: 'Redis 60s windows catch card-testing bursts a single-txn amount rule would miss (false-negative control).',
    fp: 'Policy uses CHALLENGE before BLOCK so isolated legitimate retries are not auto-declined.',
  },
  {
    value: 'account_takeover',
    label: 'Account Takeover (ATO)',
    fn: 'New device/IP + amount spike is scored by ML and challenged — stolen-session payouts are not approved silently.',
    fp: 'OTP step-up lets the real customer recover the payment instead of a hard decline.',
  },
  {
    value: 'device_farm',
    label: 'Device Farm (Shared Device)',
    fn: 'Graph USED edges show many accounts on one device — rings that look like separate customers in isolation.',
    fp: 'Only fraud-linked / emulator paths elevate; a trusted device is not treated as a farm.',
  },
  {
    value: 'ip_abuse',
    label: 'IP Abuse (Shared VPN)',
    fn: 'Shared VPN CONNECTED edges plus IP velocity catch mule clusters behind one exit node.',
    fp: 'Per-user thresholds (C1003 high vs C1001 low) avoid treating every VPN hop as a block.',
  },
  {
    value: 'amount_anomaly',
    label: 'Amount Anomaly',
    fn: 'Isolation Forest + typical_max_amount catch out-of-profile payouts (false-negative vs static MCC rules).',
    fp: 'CHALLENGE band (between user challenge and block thresholds) preserves GMV on unusual but legitimate tickets.',
  },
];

export default function ScenarioControl() {
  const [scenario, setScenario] = useState('normal');
  const [transactions, setTransactions] = useState(30);
  const [intervalMs, setIntervalMs] = useState(400);
  const [seed, setSeed] = useState(42);
  const [status, setStatus] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  const fetchStatus = async () => {
    try {
      const res = await fetch('/api/v1/scenarios/status');
      if (res.ok) setStatus(await res.json());
    } catch {}
  };

  useEffect(() => {
    fetchStatus();
    const t = setInterval(fetchStatus, 2000);
    return () => clearInterval(t);
  }, []);

  const isRunning = status?.status === 'RUNNING';

  const handleStart = async () => {
    setLoading(true);
    try {
      await fetch('/api/v1/scenarios/start', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ scenario, transactions: +transactions, interval_ms: +intervalMs, seed: +seed }),
      });
      await fetchStatus();
    } finally { setLoading(false); }
  };

  const handleStop = async () => {
    setLoading(true);
    try {
      await fetch('/api/v1/scenarios/stop', { method: 'POST' });
      await fetchStatus();
    } finally { setLoading(false); }
  };

  return (
    <div className="bg-white rounded-2xl border border-[#d2d2d7] shadow-sm p-6 mb-6">
      <div className="flex items-center justify-between mb-5">
        <div>
          <div className="flex flex-wrap items-center gap-2 mb-2">
            <span className="rounded-full border border-blue-200 bg-blue-50 px-2.5 py-1 text-[10px] font-semibold uppercase tracking-[0.14em] text-blue-700">Fresh reset</span>
            <span className="rounded-full border border-emerald-200 bg-emerald-50 px-2.5 py-1 text-[10px] font-semibold uppercase tracking-[0.14em] text-emerald-700">Synthetic baseline</span>
          </div>
          <h2 className="font-semibold text-[#1d1d1f] text-base">Scenario Simulator</h2>
          <p className="text-xs text-[#6e6e73] mt-0.5">Fresh demo data — establish a clean baseline before replaying risk conditions.</p>
        </div>
        {isRunning ? (
          <div className="flex items-center gap-2 px-3 py-1.5 bg-amber-50 border border-amber-200 rounded-full">
            <span className="h-2 w-2 rounded-full bg-amber-400 animate-ping"></span>
            <span className="text-xs font-medium text-amber-700">
              Running {status.generated_count || 0}/{status.transaction_count || 0}
            </span>
          </div>
        ) : (
          <div className="flex items-center gap-2 px-3 py-1.5 bg-[#f5f5f7] border border-[#d2d2d7] rounded-full">
            <span className="h-2 w-2 rounded-full bg-[#a1a1a6]"></span>
            <span className="text-xs font-medium text-[#6e6e73]">Idle</span>
          </div>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-5">
        <div>
          <label className="block text-xs font-medium text-[#6e6e73] mb-1.5">Scenario Type</label>
          <select
            value={scenario}
            onChange={e => setScenario(e.target.value)}
            disabled={isRunning}
            className="w-full bg-[#f5f5f7] border border-[#d2d2d7] rounded-xl px-3 py-2.5 text-sm text-[#1d1d1f] focus:outline-none focus:border-[#0071e3] focus:ring-2 focus:ring-[#0071e3]/20 disabled:opacity-50"
          >
            {SCENARIOS.map(s => (
              <option key={s.value} value={s.value}>{s.label}</option>
            ))}
          </select>
        </div>
        <div>
          <label className="block text-xs font-medium text-[#6e6e73] mb-1.5">Tx Count</label>
          <input type="number" value={transactions} onChange={e => setTransactions(+e.target.value)} disabled={isRunning}
            className="w-full bg-[#f5f5f7] border border-[#d2d2d7] rounded-xl px-3 py-2.5 text-sm text-[#1d1d1f] font-mono focus:outline-none focus:border-[#0071e3] focus:ring-2 focus:ring-[#0071e3]/20 disabled:opacity-50" />
        </div>
        <div>
          <label className="block text-xs font-medium text-[#6e6e73] mb-1.5">Interval (ms)</label>
          <input type="number" value={intervalMs} onChange={e => setIntervalMs(+e.target.value)} disabled={isRunning}
            className="w-full bg-[#f5f5f7] border border-[#d2d2d7] rounded-xl px-3 py-2.5 text-sm text-[#1d1d1f] font-mono focus:outline-none focus:border-[#0071e3] focus:ring-2 focus:ring-[#0071e3]/20 disabled:opacity-50" />
        </div>
        <div>
          <label className="block text-xs font-medium text-[#6e6e73] mb-1.5">Random Seed</label>
          <input type="number" value={seed} onChange={e => setSeed(+e.target.value)} disabled={isRunning}
            className="w-full bg-[#f5f5f7] border border-[#d2d2d7] rounded-xl px-3 py-2.5 text-sm text-[#1d1d1f] font-mono focus:outline-none focus:border-[#0071e3] focus:ring-2 focus:ring-[#0071e3]/20 disabled:opacity-50" />
        </div>
      </div>

      <div className="flex gap-3">
        {!isRunning ? (
          <button onClick={handleStart} disabled={loading}
            className="px-5 py-2.5 bg-[#0071e3] hover:bg-[#0077ed] text-white font-medium text-sm rounded-xl shadow-sm transition-all disabled:opacity-50">
            Start Scenario
          </button>
        ) : (
          <button onClick={handleStop} disabled={loading}
            className="px-5 py-2.5 bg-red-500 hover:bg-red-600 text-white font-medium text-sm rounded-xl shadow-sm transition-all disabled:opacity-50">
            Stop Scenario
          </button>
        )}
      </div>

      {SCENARIOS.filter(s => s.value === scenario).map(s => (
        <div key={s.value} className="mt-4 grid grid-cols-1 md:grid-cols-2 gap-3">
          <div className="rounded-xl border border-red-100 bg-red-50/70 p-3.5">
            <p className="text-[11px] font-semibold uppercase tracking-wide text-red-600 mb-1">False negatives (missed fraud / chargebacks)</p>
            <p className="text-xs text-[#1d1d1f] leading-relaxed">{s.fn}</p>
          </div>
          <div className="rounded-xl border border-green-100 bg-green-50/70 p-3.5">
            <p className="text-[11px] font-semibold uppercase tracking-wide text-green-700 mb-1">False positives (blocked good checkout)</p>
            <p className="text-xs text-[#1d1d1f] leading-relaxed">{s.fp}</p>
          </div>
        </div>
      ))}
    </div>
  );
}
