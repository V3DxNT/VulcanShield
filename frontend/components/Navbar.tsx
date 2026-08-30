'use client';

import React from 'react';

export default function Navbar() {
  return (
    <header className="border-b border-slate-800 bg-[#0d1322]/80 backdrop-blur sticky top-0 z-40 px-6 py-4 flex justify-between items-center">
      <div className="flex items-center gap-3">
        <div className="h-9 w-9 rounded-lg bg-gradient-to-tr from-blue-600 to-indigo-500 flex items-center justify-center font-bold text-lg text-white shadow-lg shadow-blue-500/20">
          V
        </div>
        <div>
          <h1 className="font-bold text-lg tracking-wide text-white flex items-center gap-2">
            VulcanShield
            <span className="text-xs px-2 py-0.5 rounded bg-blue-500/10 border border-blue-500/30 text-blue-400 font-mono">
              v1.0.0
            </span>
          </h1>
          <p className="text-xs text-slate-400">Adaptive AI Risk Management & Fraud Prevention Platform</p>
        </div>
      </div>
      <div className="flex items-center gap-4 text-xs font-mono text-slate-400">
        <div className="flex items-center gap-2 px-3 py-1.5 rounded-md bg-slate-900 border border-slate-800">
          <span className="h-2 w-2 rounded-full bg-emerald-400 animate-pulse"></span>
          <span>WEBSOCKET LIVE</span>
        </div>
      </div>
    </header>
  );
}
