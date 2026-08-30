'use client';

import React, { useState } from 'react';
import Navbar from '@/components/Navbar';
import KPICards from '@/components/KPICards';
import ScenarioControl from '@/components/ScenarioControl';
import TransactionStream from '@/components/TransactionStream';
import GraphVisualizer from '@/components/GraphVisualizer';
import OTPModal from '@/components/OTPModal';
import AIInvestigationModal from '@/components/AIInvestigationModal';

export default function Home() {
  const [selectedTx, setSelectedTx] = useState<any>(null);
  const [otpChallenge, setOtpChallenge] = useState<{ challengeID: string; txID: string } | null>(null);

  return (
    <div className="min-h-screen bg-[#090d16] text-slate-100 flex flex-col">
      <Navbar />

      <main className="flex-1 max-w-7xl w-full mx-auto p-6 space-y-6">
        <KPICards
          totalTx={25}
          approvedTx={18}
          challengedTx={4}
          blockedTx={3}
        />

        <ScenarioControl />

        <div className="grid grid-cols-1 gap-6">
          <TransactionStream
            onSelectTx={(tx) => setSelectedTx(tx)}
            onOpenOTP={(challengeID, txID) => setOtpChallenge({ challengeID, txID })}
          />

          <GraphVisualizer />
        </div>
      </main>

      {otpChallenge && (
        <OTPModal
          challengeID={otpChallenge.challengeID}
          transactionID={otpChallenge.txID}
          onClose={() => setOtpChallenge(null)}
          onSuccess={() => setOtpChallenge(null)}
        />
      )}

      {selectedTx && (
        <AIInvestigationModal
          transactionID={selectedTx.transaction_id}
          onClose={() => setSelectedTx(null)}
        />
      )}
    </div>
  );
}
