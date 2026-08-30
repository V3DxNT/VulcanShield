'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const navItems = [
  { label: 'Dashboard', href: '/' },
  { label: 'Network Graph', href: '/network' },
];

export default function Navbar() {
  const path = usePathname();
  return (
    <header className="bg-white/80 backdrop-blur-xl border-b border-[#d2d2d7] sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-6 py-3 flex justify-between items-center">
        <div className="flex items-center gap-8">
          <div className="flex items-center gap-2.5">
            <div className="h-8 w-8 rounded-xl bg-[#0071e3] flex items-center justify-center text-white font-bold text-sm shadow-sm">
              
            </div>
            <span className="font-semibold text-[#1d1d1f] text-sm tracking-tight">VulcanShield</span>
          </div>
          <nav className="flex items-center gap-1">
            {navItems.map(item => (
              <Link
                key={item.href}
                href={item.href}
                className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                  path === item.href
                    ? 'bg-[#f5f5f7] text-[#1d1d1f]'
                    : 'text-[#6e6e73] hover:text-[#1d1d1f] hover:bg-[#f5f5f7]'
                }`}
              >
                {item.label}
              </Link>
            ))}
          </nav>
        </div>
        <div className="flex items-center gap-2 text-xs font-medium text-[#6e6e73] bg-[#f5f5f7] px-3 py-1.5 rounded-full border border-[#d2d2d7]">
          <span className="h-1.5 w-1.5 rounded-full bg-green-500 animate-pulse"></span>
          <span>Live System</span>
        </div>
      </div>
    </header>
  );
}
