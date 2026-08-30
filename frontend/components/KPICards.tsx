'use client';

interface KPIProps {
  totalTx: number;
  approvedTx: number;
  challengedTx: number;
  blockedTx: number;
  legitRevenue?: number;
  otpVerifiedRevenue?: number;
  otpRejectedRevenue?: number;
  fraudLossAvoided?: number;
  onCategoryClick?: (category: 'all' | 'approved' | 'challenged' | 'blocked') => void;
}

export function formatINR(amount: number) {
  return new Intl.NumberFormat('en-IN', { style: 'currency', currency: 'INR', maximumFractionDigits: 0 }).format(amount);
}

export default function KPICards({ totalTx, approvedTx, challengedTx, blockedTx, legitRevenue = 0, otpVerifiedRevenue = 0, otpRejectedRevenue = 0, fraudLossAvoided = 0, onCategoryClick }: KPIProps) {
  const cards = [
    {
      key: 'all' as const,
      label: 'Total Transactions',
      value: totalTx,
      sub: 'Live volume',
      color: 'text-[#1d1d1f]',
      bg: 'bg-white',
      border: 'border-[#d2d2d7]',
    },
    {
      key: 'approved' as const,
      label: 'Valid Revenue',
      value: formatINR(legitRevenue),
      sub: 'Final approved value',
      color: 'text-emerald-700',
      bg: 'bg-emerald-50',
      border: 'border-emerald-200',
    },
    {
      key: 'challenged' as const,
      label: 'Challenge → Accepted',
      value: formatINR(otpVerifiedRevenue),
      sub: 'OTP verified / accepted',
      color: 'text-amber-700',
      bg: 'bg-amber-50',
      border: 'border-amber-200',
    },
    {
      key: 'blocked' as const,
      label: 'Challenge → Rejected',
      value: formatINR(otpRejectedRevenue),
      sub: 'OTP failed / expired',
      color: 'text-red-700',
      bg: 'bg-red-50',
      border: 'border-red-200',
    },
    {
      key: 'approved' as const,
      label: 'Approved',
      value: approvedTx,
      sub: 'Final status approved',
      color: 'text-emerald-700',
      bg: 'bg-emerald-50',
      border: 'border-emerald-200',
    },
    {
      key: 'challenged' as const,
      label: 'Challenged',
      value: challengedTx,
      sub: 'Awaiting step-up verification',
      color: 'text-orange-700',
      bg: 'bg-orange-50',
      border: 'border-orange-200',
    },
    {
      key: 'blocked' as const,
      label: 'Blocked',
      value: blockedTx,
      sub: 'Final status blocked',
      color: 'text-red-700',
      bg: 'bg-red-50',
      border: 'border-red-200',
    },
    {
      key: 'all' as const,
      label: 'Fraud Loss Avoided',
      value: formatINR(fraudLossAvoided),
      sub: 'Protected by policy',
      color: 'text-violet-700',
      bg: 'bg-violet-50',
      border: 'border-violet-200',
    },
  ];

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4 mb-6">
      {cards.map((c, i) => (
        <button
          key={`${c.label}-${i}`}
          type="button"
          onClick={() => onCategoryClick?.(c.key)}
          className={`p-5 rounded-2xl border ${c.bg} ${c.border} shadow-sm text-left transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md`}
        >
          <div className="flex items-center justify-between gap-3 mb-3">
            <p className="text-[10px] font-semibold uppercase tracking-[0.16em] text-[#6e6e73]">{c.label}</p>
            <span className="h-2.5 w-2.5 rounded-full bg-[#1d1d1f]/10" />
          </div>
          <p className={`text-2xl font-semibold ${c.color} tabular-nums`}>{typeof c.value === 'number' ? c.value.toLocaleString('en-IN') : c.value}</p>
          <p className="text-xs text-[#6e6e73] mt-2">{c.sub}</p>
        </button>
      ))}
    </div>
  );
}
