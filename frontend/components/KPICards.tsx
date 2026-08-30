'use client';

import React from 'react';

interface KPIProps {
  totalTx: number;
  approvedTx: number;
  challengedTx: number;
  blockedTx: number;
}

export default function KPICards({ totalTx, approvedTx, challengedTx, blockedTx }: KPIProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
      <div className="p-4 rounded-xl bg-slate-900/60 border border-slate-800 backdrop-blur">
        <p className="text-xs font-mono uppercase tracking-wider text-slate-400">Total Transactions</p>
        <p className="text-2xl font-bold text-white mt-1 font-mono">{totalTx}</p>
        <div className="mt-2 text-[11px] text-slate-500">Real-time Stream Volume</div>
      </div>
      <div className="p-4 rounded-xl bg-emerald-950/20 border border-emerald-900/40 backdrop-blur">
        <p className="text-xs font-mono uppercase tracking-wider text-emerald-400">Approved (ALLOW)</p>
        <p className="text-2xl font-bold text-emerald-400 mt-1 font-mono">{approvedTx}</p>
        <div className="mt-2 text-[11px] text-emerald-600">Risk Score &lt; Challenge Threshold</div>
      </div>
      <div className="p-4 rounded-xl bg-amber-950/20 border border-amber-900/40 backdrop-blur">
        <p className="text-xs font-mono uppercase tracking-wider text-amber-400">Challenged (OTP)</p>
        <p className="text-2xl font-bold text-amber-400 mt-1 font-mono">{challengedTx}</p>
        <div className="mt-2 text-[11px] text-amber-600">Step-Up Verification Triggered</div>
      </div>
      <div className="p-4 rounded-xl bg-rose-950/20 border border-rose-900/40 backdrop-blur">
        <p className="text-xs font-mono uppercase tracking-wider text-rose-400">Blocked (BLOCK)</p>
        <p className="text-2xl font-bold text-rose-400 mt-1 font-mono">{blockedTx}</p>
        <div className="mt-2 text-[11px] text-rose-600">High Risk Threshold Exceeded</div>
      </div>
    </div>
  );
}
