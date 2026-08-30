'use client';
import { useState } from 'react';

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
        setTimeout(() => { onSuccess(); onClose(); }, 1500);
      }
    } catch {
      setResult({ message: 'Verification request failed' });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/30 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl border border-[#d2d2d7] shadow-2xl w-full max-w-sm">
        <div className="px-6 py-5 border-b border-[#e8e8ed] flex justify-between items-center">
          <div>
            <h3 className="font-semibold text-[#1d1d1f] text-base">Step-Up Verification</h3>
            <p className="text-xs text-[#6e6e73] mt-0.5 font-mono">{transactionID}</p>
          </div>
          <button onClick={onClose} className="h-8 w-8 flex items-center justify-center rounded-full bg-[#f5f5f7] hover:bg-[#e8e8ed] text-[#6e6e73] transition-colors text-lg">×</button>
        </div>

        <div className="px-6 py-5">
          <div className="flex items-start gap-3 p-3.5 rounded-xl bg-amber-50 border border-amber-200 mb-5">
            <div className="h-8 w-8 rounded-full bg-amber-100 flex items-center justify-center text-amber-600 font-bold text-sm flex-shrink-0">!</div>
            <div>
              <p className="text-sm font-medium text-amber-800">OTP Required</p>
              <p className="text-xs text-amber-700 mt-0.5">This transaction was flagged as suspicious. Enter the 6-digit OTP to proceed (60s expiration).</p>
            </div>
          </div>

          <label className="block text-xs font-medium text-[#6e6e73] mb-1.5">6-Digit OTP Code</label>
          <input
            type="text" maxLength={6} value={otpCode}
            onChange={e => setOtpCode(e.target.value.replace(/\D/g, ''))}
            placeholder="• • • • • •"
            className="w-full bg-[#f5f5f7] border border-[#d2d2d7] rounded-xl px-4 py-3.5 text-2xl font-mono tracking-[0.5em] text-center text-[#1d1d1f] focus:outline-none focus:border-[#0071e3] focus:ring-2 focus:ring-[#0071e3]/20 mb-4"
          />

          {result && (
            <div className={`p-3 rounded-xl text-sm mb-4 ${result.status === 'VERIFIED' ? 'bg-green-50 border border-green-200 text-green-700' : 'bg-red-50 border border-red-200 text-red-700'}`}>
              {result.status === 'VERIFIED' ? '✓ ' : '✗ '}{result.message}
            </div>
          )}

          <div className="flex gap-3">
            <button onClick={onClose} className="flex-1 py-2.5 bg-[#f5f5f7] hover:bg-[#e8e8ed] text-[#1d1d1f] font-medium text-sm rounded-xl border border-[#d2d2d7] transition-colors">Cancel</button>
            <button onClick={handleVerify} disabled={loading || otpCode.length < 6}
              className="flex-1 py-2.5 bg-[#0071e3] hover:bg-[#0077ed] text-white font-medium text-sm rounded-xl transition-colors disabled:opacity-50">
              {loading ? 'Verifying…' : 'Verify OTP'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
