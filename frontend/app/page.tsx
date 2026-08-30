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
  const [kpi, setKpi] = useState({ total: 0, approved: 0, challenged: 0, blocked: 0, legitRevenue: 0, otpVerifiedRevenue: 0, otpRejectedRevenue: 0, fraudLossAvoided: 0 });

  useEffect(() => {
    const fetch_kpi = async () => {
      try {
        const res = await fetch('/api/v1/transactions?limit=500');
        if (res.ok) {
          const json = await res.json();
          const txs: any[] = json.data ?? [];
          const approved = txs.filter(t => t.status === 'APPROVED').length;
          const challenged = txs.filter(t => t.status === 'CHALLENGED').length;
          const blocked = txs.filter(t => t.status === 'BLOCKED').length;
          const approvedRevenue = txs.filter(t => t.status === 'APPROVED').reduce((sum, t) => sum + Number(t.amount || 0), 0);
          const otpVerifiedRevenue = txs.filter(t => t.challenge_status === 'VERIFIED').reduce((sum, t) => sum + Number(t.amount || 0), 0);
          const otpRejectedRevenue = txs.filter(t => t.challenge_status === 'FAILED' || t.challenge_status === 'EXPIRED').reduce((sum, t) => sum + Number(t.amount || 0), 0);
          const fraudLossAvoided = txs.filter(t => t.status === 'BLOCKED').reduce((sum, t) => sum + Number(t.amount || 0), 0);

          setKpi({
            total: txs.length,
            approved,
            challenged,
            blocked,
            legitRevenue: approvedRevenue + otpVerifiedRevenue,
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
          <h1 className="text-2xl font-semibold text-[#1d1d1f] tracking-tight">Risk Operations Dashboard</h1>
          <p className="text-sm text-[#6e6e73] mt-1">Real-time transaction monitoring · ML scoring · AI-powered investigation</p>
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
        />
        <ScenarioControl />

        <div className="grid grid-cols-1 gap-6">
          <TransactionStream
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
