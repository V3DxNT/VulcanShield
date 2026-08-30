'use client';
import { useState, useEffect } from 'react';
import Navbar from '@/components/Navbar';
import KPICards from '@/components/KPICards';
import ScenarioControl from '@/components/ScenarioControl';
import TransactionStream from '@/components/TransactionStream';
import GraphVisualizer from '@/components/GraphVisualizer';
import AIInvestigationModal from '@/components/AIInvestigationModal';

export default function Home() {
  const [selectedTx, setSelectedTx] = useState<any>(null);
  const [activeFilter, setActiveFilter] = useState<'all' | 'approved' | 'challenged' | 'blocked'>('all');
  const [kpi, setKpi] = useState({ total: 0, approved: 0, challenged: 0, blocked: 0, legitRevenue: 0, otpVerifiedRevenue: 0, otpRejectedRevenue: 0, fraudLossAvoided: 0 });

  useEffect(() => {
    const fetch_kpi = async () => {
      try {
        const res = await fetch('/api/v1/transactions?limit=500', { cache: 'no-store' });
        if (res.ok) {
          const json = await res.json();
          const txs: any[] = json.data ?? [];

          const approved = txs.filter(t => t.status === 'APPROVED' || t.challenge_status === 'VERIFIED').length;
          const challenged = txs.filter(t => t.status === 'CHALLENGED' || t.challenge_status === 'PENDING').length;
          const blocked = txs.filter(t => t.status === 'BLOCKED' || ['FAILED', 'EXPIRED'].includes(t.challenge_status ?? '')).length;
          const otpVerified = txs.filter(t => t.challenge_status === 'VERIFIED').length;
          const otpRejected = txs.filter(t => ['FAILED', 'EXPIRED'].includes(t.challenge_status ?? '')).length;

          const approvedRevenue = txs.filter(t => t.status === 'APPROVED' || t.challenge_status === 'VERIFIED').reduce((sum, t) => sum + Number(t.amount || 0), 0);
          const otpVerifiedRevenue = txs.filter(t => t.challenge_status === 'VERIFIED').reduce((sum, t) => sum + Number(t.amount || 0), 0);
          const otpRejectedRevenue = txs.filter(t => ['FAILED', 'EXPIRED'].includes(t.challenge_status ?? '')).reduce((sum, t) => sum + Number(t.amount || 0), 0);
          const fraudLossAvoided = txs.filter(t => t.status === 'BLOCKED' || ['FAILED', 'EXPIRED'].includes(t.challenge_status ?? '')).reduce((sum, t) => sum + Number(t.amount || 0), 0);

          setKpi({
            total: txs.length,
            approved,
            challenged,
            blocked,
            legitRevenue: approvedRevenue,
            otpVerifiedRevenue,
            otpRejectedRevenue,
            fraudLossAvoided,
          });
        }
      } catch {}
    };
    fetch_kpi();
    const t = setInterval(fetch_kpi, 3000);
    return () => clearInterval(t);
  }, []);

  return (
    <div className="min-h-screen bg-[#f5f5f7]">
      <Navbar />
      <main className="max-w-7xl mx-auto px-6 py-8">
        <div className="mb-8">
          <div className="flex flex-col lg:flex-row lg:items-end lg:justify-between gap-4">
            <div>
              <p className="text-[11px] font-semibold uppercase tracking-[0.18em] text-[#6e6e73]">Fraud operations</p>
              <h1 className="mt-2 text-3xl font-semibold text-[#1d1d1f] tracking-tight">Risk Operations Console</h1>
              <p className="text-sm text-[#6e6e73] mt-2">Live transaction monitoring · ML scoring · policy enforcement · AI investigation</p>
            </div>
            <div className="flex flex-wrap gap-2">
              {[
                { label: 'Live feed', value: 'Streaming' },
                { label: 'Policy mode', value: 'Deterministic' },
                { label: 'LLM state', value: 'Verified' },
              ].map(item => (
                <div key={item.label} className="rounded-full border border-[#d2d2d7] bg-white px-3 py-1.5 shadow-sm">
                  <span className="text-[10px] uppercase tracking-[0.14em] text-[#6e6e73]">{item.label}</span>
                  <div className="text-xs font-semibold text-[#1d1d1f] mt-0.5">{item.value}</div>
                </div>
              ))}
            </div>
          </div>
        </div>

        <KPICards
          totalTx={kpi.total}
          approvedTx={kpi.approved}
          challengedTx={kpi.challenged}
          blockedTx={kpi.blocked}
          legitRevenue={kpi.legitRevenue}
          otpVerifiedRevenue={kpi.otpVerifiedRevenue}
          otpRejectedRevenue={kpi.otpRejectedRevenue}
          fraudLossAvoided={kpi.fraudLossAvoided}
          onCategoryClick={setActiveFilter}
        />
        <ScenarioControl />

        <div className="grid grid-cols-1 gap-6">
          <TransactionStream
            activeFilter={activeFilter}
            onFilterChange={setActiveFilter}
            onSelectTx={tx => setSelectedTx(tx)}
          />
          <GraphVisualizer />
        </div>
      </main>

      {selectedTx && (
        <AIInvestigationModal
          transactionID={selectedTx.transaction_id}
          onClose={() => setSelectedTx(null)}
        />
      )}
    </div>
  );
}
