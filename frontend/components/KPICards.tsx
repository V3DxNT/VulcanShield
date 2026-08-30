'use client';

interface KPIProps {
  totalTx: number;
  approvedTx: number;
  challengedTx: number;
  blockedTx: number;
}

export function formatINR(amount: number) {
  return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 }).format(amount);
}

export default function KPICards({ totalTx, approvedTx, challengedTx, blockedTx }: KPIProps) {
  const cards = [
    {
      label: 'Total Transactions',
      value: totalTx,
      sub: 'Real-time volume',
      color: 'text-[#1d1d1f]',
      bg: 'bg-white',
      border: 'border-[#d2d2d7]',
    },
    {
      label: 'Approved',
      value: approvedTx,
      sub: 'Risk below threshold',
      color: 'text-green-600',
      bg: 'bg-green-50',
      border: 'border-green-100',
    },
    {
      label: 'Challenged (OTP)',
      value: challengedTx,
      sub: 'Step-up verification',
      color: 'text-amber-600',
      bg: 'bg-amber-50',
      border: 'border-amber-100',
    },
    {
      label: 'Blocked',
      value: blockedTx,
      sub: 'High risk blocked',
      color: 'text-red-500',
      bg: 'bg-red-50',
      border: 'border-red-100',
    },
  ];

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
      {cards.map((c, i) => (
        <div key={i} className={`p-5 rounded-2xl border ${c.bg} ${c.border} shadow-sm`}>
          <p className="text-xs font-medium text-[#6e6e73] uppercase tracking-wide mb-1">{c.label}</p>
          <p className={`text-3xl font-semibold ${c.color} tabular-nums`}>{c.value.toLocaleString('en-IN')}</p>
          <p className="text-xs text-[#a1a1a6] mt-1">{c.sub}</p>
        </div>
      ))}
    </div>
  );
}
