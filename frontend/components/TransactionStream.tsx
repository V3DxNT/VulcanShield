'use client';
import { useState, useEffect, useCallback } from 'react';

function formatINR(amount: number) {
  return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 }).format(amount);
}

function StatusBadge({ status, challengeStatus }: { status: string; challengeStatus?: string }) {
  if (challengeStatus === 'VERIFIED') return <span className="px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-50 text-green-700 border border-green-200">OTP accepted</span>;
  if (challengeStatus === 'FAILED' || challengeStatus === 'EXPIRED') return <span className="px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-50 text-red-600 border border-red-200">OTP rejected</span>;
  if (status === 'APPROVED') return <span className="px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-50 text-green-700 border border-green-200">Approved</span>;
  if (status === 'BLOCKED') return <span className="px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-50 text-red-600 border border-red-200">Blocked</span>;
  if (status === 'CHALLENGED' || challengeStatus === 'PENDING') return (
    <span className="px-2.5 py-0.5 rounded-full text-xs font-medium bg-amber-50 text-amber-700 border border-amber-200 flex items-center gap-1">
      <span className="h-1.5 w-1.5 rounded-full bg-amber-500 animate-pulse"></span>
      Challenge pending
    </span>
  );
  return <span className="px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-500 border border-gray-200">Pending</span>;
}

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
  challenge_status?: string;
  timestamp: string;
}

interface Props {
  onSelectTx: (tx: Transaction) => void;
}

export default function TransactionStream({ onSelectTx }: Props) {
  const [txs, setTxs] = useState<Transaction[]>([]);
  const [filter, setFilter] = useState('all');
  const [userQuery, setUserQuery] = useState('');

  const fetchTxList = useCallback(async () => {
    try {
      const res = await fetch('/api/v1/transactions?limit=50');
      if (res.ok) {
        const json = await res.json();
        if (json.data) setTxs(json.data);
      }
    } catch {}
  }, []);

  useEffect(() => {
    fetchTxList();
    const interval = setInterval(fetchTxList, 1500);

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.hostname}:8080/api/v1/ws`;
    let ws: WebSocket | null = null;
    try {
      ws = new WebSocket(wsUrl);
      ws.onmessage = event => {
        try {
          const payload = JSON.parse(event.data);
          if (payload.event_type === 'transaction_created') {
            setTxs(prev => [payload.data, ...prev.slice(0, 49)]);
          }
        } catch {}
      };
    } catch {}

    return () => { clearInterval(interval); ws?.close(); };
  }, [fetchTxList]);

  const filtered = txs.filter(tx => {
    const isApproved = tx.status === 'APPROVED' || tx.challenge_status === 'VERIFIED';
    const isChallenged = tx.status === 'CHALLENGED' || ['PENDING', 'VERIFIED', 'FAILED', 'EXPIRED'].includes(tx.challenge_status ?? '');
    const isBlocked = tx.status === 'BLOCKED' || tx.challenge_status === 'FAILED' || tx.challenge_status === 'EXPIRED';

    const matchesStatus = filter === 'all'
      ? true
      : filter === 'approved'
        ? isApproved
        : filter === 'challenged'
          ? isChallenged
          : isBlocked;

    const matchesUser = userQuery.trim() === '' || tx.user_id.toLowerCase().includes(userQuery.toLowerCase());
    return matchesStatus && matchesUser;
  });

  return (
    <div className="bg-white rounded-2xl border border-[#d2d2d7] shadow-sm overflow-hidden">
      <div className="px-6 py-4 border-b border-[#e8e8ed] flex items-center justify-between">
        <div>
          <h2 className="font-semibold text-[#1d1d1f] text-base">Live Transaction Stream</h2>
          <p className="text-xs text-[#6e6e73] mt-0.5">{txs.length} transactions • auto-refreshing</p>
        </div>
        <div className="flex items-center gap-2 flex-wrap">
          <div className="flex gap-1 bg-[#f5f5f7] p-1 rounded-lg border border-[#e8e8ed]">
            {['all', 'approved', 'challenged', 'blocked'].map(f => (
              <button key={f} onClick={() => setFilter(f)}
                className={`px-3 py-1 text-xs font-medium rounded-md transition-colors capitalize ${
                  filter === f ? 'bg-white text-[#1d1d1f] shadow-sm border border-[#d2d2d7]' : 'text-[#6e6e73] hover:text-[#1d1d1f]'
                }`}>
                {f}
              </button>
            ))}
          </div>
          <input
            value={userQuery}
            onChange={e => setUserQuery(e.target.value)}
            placeholder="Search customer ID"
            className="w-44 rounded-lg border border-[#d2d2d7] bg-white px-2.5 py-1.5 text-xs text-[#1d1d1f] placeholder:text-[#a1a1a6]"
          />
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left text-sm">
          <thead>
            <tr className="border-b border-[#e8e8ed] bg-[#f5f5f7]">
              {['Transaction', 'User', 'Device / IP', 'Amount', 'Channel', 'Status', 'Time', ''].map(h => (
                <th key={h} className="px-4 py-3 text-xs font-medium text-[#6e6e73] uppercase tracking-wide first:pl-6 last:pr-6">{h}</th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-[#f5f5f7]">
            {filtered.length === 0 ? (
              <tr>
                <td colSpan={8} className="px-6 py-12 text-center text-sm text-[#a1a1a6]">
                  No transactions yet. Start a scenario to stream live data.
                </td>
              </tr>
            ) : filtered.map(tx => (
              <tr key={tx.transaction_id} onClick={() => onSelectTx(tx)}
                className="hover:bg-[#f5f5f7] cursor-pointer transition-colors group">
                <td className="pl-6 pr-4 py-3.5 font-mono text-xs font-medium text-[#0071e3]">{tx.transaction_id}</td>
                <td className="px-4 py-3.5 text-xs font-medium text-[#1d1d1f]">{tx.user_id}</td>
                <td className="px-4 py-3.5">
                  <div className="text-xs text-[#1d1d1f] font-mono">{tx.device_id}</div>
                  <div className="text-xs text-[#a1a1a6] font-mono">{tx.ip_address}</div>
                </td>
                <td className="px-4 py-3.5 font-semibold text-sm text-[#1d1d1f]">
                  {tx.currency === 'INR' ? `₹${tx.amount?.toLocaleString('en-IN')}` : `$${tx.amount?.toFixed(2)}`}
                </td>
                <td className="px-4 py-3.5 text-xs text-[#6e6e73]">{tx.channel}</td>
                <td className="px-4 py-3.5">
                  <StatusBadge status={tx.status} challengeStatus={tx.challenge_status} />
                </td>
                <td className="px-4 py-3.5 text-xs text-[#a1a1a6]">
                  {new Date(tx.timestamp).toLocaleTimeString('en-IN')}
                </td>
                <td className="pr-6 pl-4 py-3.5 text-right">
                  <button onClick={e => { e.stopPropagation(); onSelectTx(tx); }}
                    className="opacity-0 group-hover:opacity-100 transition-opacity px-3 py-1.5 bg-[#0071e3] text-white rounded-lg text-xs font-medium">
                    Investigate
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
