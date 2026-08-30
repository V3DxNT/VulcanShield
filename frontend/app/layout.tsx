import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'VulcanShield — Adaptive AI Risk Management',
  description: 'Real-time financial transaction risk monitoring, ML-powered fraud detection, and explainable AI investigations.',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet" />
      </head>
      <body className="bg-[#f5f5f7] text-[#1d1d1f] antialiased min-h-screen">
        {children}
      </body>
    </html>
  );
}
