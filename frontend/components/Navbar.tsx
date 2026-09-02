'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const navItems = [
  { label: 'Dashboard', href: '/' },
  { label: 'Network Graph', href: '/network' },
];

function ShieldLogo() {
  return (
    <div className="relative h-10 w-10 shrink-0 rounded-[14px] bg-gradient-to-br from-[#0d67d1] via-[#0071e3] to-[#0b47a5] shadow-[0_10px_28px_rgba(0,113,227,0.32)] ring-2 ring-white/80">
      <div className="absolute inset-0 rounded-[14px] bg-[radial-gradient(circle_at_top,_rgba(255,255,255,0.38),_transparent_42%)]" />
      <div className="absolute inset-0 flex items-center justify-center text-[0.9rem] font-black tracking-[-0.06em] text-white drop-shadow-sm">
        VS
      </div>
    </div>
  );
}

export default function Navbar() {
  const path = usePathname();
  return (
    <header className="sticky top-4 z-50 mx-auto max-w-7xl px-4">
      <div className="flex items-center justify-between rounded-2xl border border-[#d2d2d7] bg-white/80 px-4 py-2.5 shadow-[0_12px_34px_rgba(15,23,42,0.08)] backdrop-blur-xl">
        <div className="flex items-center gap-8">
          <div className="flex items-center gap-3">
            <ShieldLogo />
            <div>
              <span className="block text-sm font-semibold tracking-tight text-[#1d1d1f]">VulcanShield</span>
              <span className="text-[10px] uppercase tracking-[0.18em] text-[#6e6e73]">Fraud ops</span>
            </div>
          </div>
          <nav className="flex items-center gap-1">
            {navItems.map(item => (
              <Link
                key={item.href}
                href={item.href}
                className={`rounded-xl px-3 py-1.5 text-sm font-medium transition-colors ${
                  path === item.href
                    ? 'bg-[#f5f5f7] text-[#1d1d1f]'
                    : 'text-[#6e6e73] hover:bg-[#f5f5f7] hover:text-[#1d1d1f]'
                }`}
              >
                {item.label}
              </Link>
            ))}
          </nav>
        </div>
        <div className="flex items-center gap-2 rounded-full border border-[#d2d2d7] bg-[#f5f5f7] px-3 py-1.5 text-xs font-medium text-[#6e6e73]">
          <span className="h-1.5 w-1.5 rounded-full bg-green-500 animate-pulse"></span>
          <span>Online</span>
        </div>
      </div>
    </header>
  );
}
