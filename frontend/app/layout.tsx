import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'VulcanShield — Adaptive AI Risk Management Platform',
  description: 'Real-time financial transaction risk monitoring, fraud detection, and explainable AI investigations.',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="antialiased bg-[#090d16] text-slate-100 min-h-screen">
        {children}
      </body>
    </html>
  );
}
