'use client';

import React, { useState, useEffect } from 'react';

interface Transaction {
  transaction_id: string;
  user_id: string;
  device_id: string;
  ip_address: string;
  merchant_id: string;
  amount: number;
  currency: string;
  channel: string;
  status: string;
  timestamp: string;
  risk_score?: number;
}

interface StreamProps {
  onSelectTx: (tx: Transaction) => void;
  onOpenOTP: (challengeID: string, txID: string) => void;
}

export default function TransactionStream({ onSelectTx, onOpenOTP }: StreamProps) {
  const [transactions, setTransactions] = useState<Transaction[]>([]);

  const fetchTxList = async () => {
    try {
      const res = await fetch('/api/v1/transactions?limit=25');
      if (res.ok) {
        const json = await res.json();
        if (json.data) {
          setTransactions(json.data);
        }
      }
    } catch (e) {
      // ignore
    }
  };

  useEffect(() => {
    fetchTxList();
    const interval = setInterval(fetchTxList, 1500);

    // WebSocket real-time connection
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.hostname}:8080/api/v1/ws`;
    
    let ws: WebSocket | null = null;
    try {
      ws = new WebSocket(wsUrl);
      ws.onmessage = (event) => {
        try {
          const payload = JSON.parse(event.data);
          if (payload.event_type === 'transaction_created') {
            setTransactions((prev) => [payload.data, ...prev.slice(0, 24)]);
          }
        } catch (err) {
          // ignore
        }
      };
    } catch (err) {
      // ignore
    }

    return () => {
      clearInterval(interval);
      if (ws) ws.close();
    };
  }, []);

  const getStatusBadge = (status: string, tx: Transaction) => {
    switch (status) {
      case 'APPROVED':
        return <span className="px-2 py-0.5 rounded text-[11px] font-mono bg-emerald-500/10 border border-emerald-500/30 text-emerald-400">APPROVED</span>;
      case 'CHALLENGED':
        return (
          <button
            onClick={(e) => {
              e.stopPropagation();
              onOpenOTP(`CH-${tx.transaction_id}`, tx.transaction_id);
            }}
            className="px-2 py-0.5 rounded text-[11px] font-mono bg-amber-500/20 border border-amber-500/40 text-amber-300 hover:bg-amber-500/30 flex items-center gap-1"
          >
            CHALLENGED (OTP)
          </button>
        );
      case 'BLOCKED':
        return <span className="px-2 py-0.5 rounded text-[11px] font-mono bg-rose-500/10 border border-rose-500/30 text-rose-400">BLOCKED</span>;
      default:
        return <span className="px-2 py-0.5 rounded text-[11px] font-mono bg-slate-800 text-slate-400">PENDING</span>;
    }
  };

  return (
    <div className="p-5 rounded-xl bg-slate-900/60 border border-slate-800 backdrop-blur">
      <h2 className="font-bold text-sm uppercase tracking-wider text-slate-300 mb-4 flex justify-between items-center">
        <span>Live Transaction Stream</span>
        <span className="text-xs font-mono font-normal text-slate-500">Auto-refreshing</span>
      </h2>

      <div className="overflow-x-auto">
        <table className="w-full text-left text-xs font-mono">
          <thead className="text-[11px] text-slate-500 uppercase border-b border-slate-800">
            <tr>
              <th className="pb-3">Transaction ID</th>
              <th className="pb-3">User ID</th>
              <th className="pb-3">Device / IP</th>
              <th className="pb-3">Amount</th>
              <th className="pb-3">Status</th>
              <th className="pb-3">Timestamp</th>
              <th className="pb-3 text-right">Action</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800/50">
            {transactions.length === 0 ? (
              <tr>
                <td colSpan={7} className="py-8 text-center text-slate-500">
                  No transactions generated yet. Start a scenario to stream traffic.
                </td>
              </tr>
            ) : (
              transactions.map((tx) => (
                <tr
                  key={tx.transaction_id}
                  onClick={() => onSelectTx(tx)}
                  className="hover:bg-slate-800/40 cursor-pointer transition-colors"
                >
                  <td className="py-3 font-semibold text-blue-400">{tx.transaction_id}</td>
                  <td className="py-3 text-slate-300">{tx.user_id}</td>
                  <td className="py-3 text-slate-400">
                    <div>{tx.device_id}</div>
                    <div className="text-[10px] text-slate-500">{tx.ip_address}</div>
                  </td>
                  <td className="py-3 font-bold text-slate-200">${tx.amount?.toFixed(2)}</td>
                  <td className="py-3">{getStatusBadge(tx.status, tx)}</td>
                  <td className="py-3 text-slate-500 text-[10px]">
                    {new Date(tx.timestamp).toLocaleTimeString()}
                  </td>
                  <td className="py-3 text-right">
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        onSelectTx(tx);
                      }}
                      className="px-2.5 py-1 bg-slate-800 hover:bg-slate-700 text-slate-200 rounded text-[11px] transition-colors"
                    >
                      Investigate
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
