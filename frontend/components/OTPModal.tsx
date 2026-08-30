'use client';

import React, { useState } from 'react';

interface OTPModalProps {
  challengeID: string;
  transactionID: string;
  onClose: () => void;
  onSuccess: () => void;
}

export default function OTPModal({ challengeID, transactionID, onClose, onSuccess }: OTPModalProps) {
  const [otpCode, setOtpCode] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<any>(null);

  const handleVerify = async () => {
    setLoading(true);
    setResult(null);
    try {
      const res = await fetch(`/api/v1/challenges/${challengeID}/verify`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ otp_code: otpCode }),
      });
      const data = await res.json();
      setResult(data);
      if (res.ok && data.status === 'VERIFIED') {
        setTimeout(() => {
          onSuccess();
          onClose();
        }, 1500);
      }
    } catch (e) {
      setResult({ message: 'Verification request failed' });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-md p-6 shadow-2xl">
        <div className="flex justify-between items-center mb-4">
          <h3 className="font-bold text-base text-white flex items-center gap-2">
            Step-Up Authentication Challenge
          </h3>
          <button onClick={onClose} className="text-slate-400 hover:text-white text-lg">×</button>
        </div>

        <p className="text-xs text-slate-400 mb-4">
          Transaction <span className="font-mono text-blue-400">{transactionID}</span> triggered a step-up OTP requirement (60s expiration).
        </p>

        <div className="mb-4">
          <label className="block text-xs font-mono text-slate-300 mb-1">Enter 6-Digit OTP Code</label>
          <input
            type="text"
            maxLength={6}
            value={otpCode}
            onChange={(e) => setOtpCode(e.target.value)}
            placeholder="381924"
            className="w-full bg-slate-950 border border-slate-800 rounded-xl px-4 py-3 font-mono text-lg text-center tracking-widest text-white focus:outline-none focus:border-amber-500"
          />
        </div>

        {result && (
          <div className={`p-3 rounded-lg text-xs font-mono mb-4 ${
            result.status === 'VERIFIED' ? 'bg-emerald-500/10 border border-emerald-500/30 text-emerald-400' : 'bg-rose-500/10 border border-rose-500/30 text-rose-400'
          }`}>
            {result.message}
          </div>
        )}

        <div className="flex gap-3">
          <button
            onClick={onClose}
            className="flex-1 py-2.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl text-xs font-semibold"
          >
            Cancel
          </button>
          <button
            onClick={handleVerify}
            disabled={loading || otpCode.length < 6}
            className="flex-1 py-2.5 bg-amber-500 hover:bg-amber-400 text-slate-950 font-bold rounded-xl text-xs shadow-lg shadow-amber-500/20"
          >
            Submit OTP
          </button>
        </div>
      </div>
    </div>
  );
}
