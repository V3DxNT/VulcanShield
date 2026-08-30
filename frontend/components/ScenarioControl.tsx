'use client';

import React, { useState, useEffect } from 'react';

export default function ScenarioControl() {
  const [scenario, setScenario] = useState('normal');
  const [transactions, setTransactions] = useState(20);
  const [intervalMs, setIntervalMs] = useState(500);
  const [seed, setSeed] = useState(42);
  const [status, setStatus] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  const fetchStatus = async () => {
    try {
      const res = await fetch('/api/v1/scenarios/status');
      if (res.ok) {
        const data = await res.json();
        setStatus(data);
      }
    } catch (e) {
      // ignore
    }
  };

  useEffect(() => {
    fetchStatus();
    const timer = setInterval(fetchStatus, 2000);
    return () => clearInterval(timer);
  }, []);

  const handleStart = async () => {
    setLoading(true);
    try {
      await fetch('/api/v1/scenarios/start', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          scenario,
          transactions: Number(transactions),
          interval_ms: Number(intervalMs),
          seed: Number(seed),
        }),
      });
      await fetchStatus();
    } catch (e) {
      // ignore
    } finally {
      setLoading(false);
    }
  };

  const handleStop = async () => {
    setLoading(true);
    try {
      await fetch('/api/v1/scenarios/stop', { method: 'POST' });
      await fetchStatus();
    } catch (e) {
      // ignore
    } finally {
      setLoading(false);
    }
  };

  const isRunning = status && status.status === 'RUNNING';

  return (
    <div className="p-5 rounded-xl bg-slate-900/60 border border-slate-800 backdrop-blur mb-6">
      <div className="flex justify-between items-center mb-4">
        <h2 className="font-bold text-sm uppercase tracking-wider text-slate-300 flex items-center gap-2">
          Scenario Simulator Controls
          <span className="text-[10px] px-2 py-0.5 rounded bg-slate-800 text-slate-400 font-mono font-normal">
            Deterministic Engine
          </span>
        </h2>
        {isRunning ? (
          <span className="text-xs px-2.5 py-1 rounded-full bg-amber-500/10 border border-amber-500/30 text-amber-400 font-mono flex items-center gap-1.5">
            <span className="h-1.5 w-1.5 rounded-full bg-amber-400 animate-ping"></span>
            RUNNING ({status.generated_count || 0}/{status.transaction_count || 0})
          </span>
        ) : (
          <span className="text-xs px-2.5 py-1 rounded-full bg-slate-800 text-slate-400 font-mono">
            IDLE
          </span>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
        <div>
          <label className="block text-xs text-slate-400 mb-1">Scenario Type</label>
          <select
            value={scenario}
            onChange={(e) => setScenario(e.target.value)}
            disabled={isRunning}
            className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-xs text-white focus:outline-none focus:border-blue-500"
          >
            <option value="normal">Normal Traffic</option>
            <option value="velocity_attack">Velocity Attack (Carding Burst)</option>
            <option value="account_takeover">Account Takeover (ATO)</option>
            <option value="device_farm">Device Farm (Shared Device)</option>
            <option value="ip_abuse">IP Abuse (Shared VPN)</option>
            <option value="amount_anomaly">Amount Anomaly</option>
          </select>
        </div>

        <div>
          <label className="block text-xs text-slate-400 mb-1">Tx Count</label>
          <input
            type="number"
            value={transactions}
            onChange={(e) => setTransactions(Number(e.target.value))}
            disabled={isRunning}
            className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-xs text-white font-mono focus:outline-none focus:border-blue-500"
          />
        </div>

        <div>
          <label className="block text-xs text-slate-400 mb-1">Interval (ms)</label>
          <input
            type="number"
            value={intervalMs}
            onChange={(e) => setIntervalMs(Number(e.target.value))}
            disabled={isRunning}
            className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-xs text-white font-mono focus:outline-none focus:border-blue-500"
          />
        </div>

        <div>
          <label className="block text-xs text-slate-400 mb-1">Random Seed</label>
          <input
            type="number"
            value={seed}
            onChange={(e) => setSeed(Number(e.target.value))}
            disabled={isRunning}
            className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-xs text-white font-mono focus:outline-none focus:border-blue-500"
          />
        </div>
      </div>

      <div className="flex gap-3">
        {!isRunning ? (
          <button
            onClick={handleStart}
            disabled={loading}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white font-semibold text-xs rounded-lg shadow-lg shadow-blue-600/20 transition-all flex items-center gap-2"
          >
            Start Scenario
          </button>
        ) : (
          <button
            onClick={handleStop}
            disabled={loading}
            className="px-4 py-2 bg-rose-600 hover:bg-rose-500 text-white font-semibold text-xs rounded-lg shadow-lg shadow-rose-600/20 transition-all flex items-center gap-2"
          >
            Stop Scenario
          </button>
        )}
      </div>
    </div>
  );
}
