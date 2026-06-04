import { useState } from "react";

export function Nav() {
  const [open, setOpen] = useState(false);

  return (
    <div className="fixed top-0 inset-x-0 z-50 px-4 pt-4">
      {/* Floating glass panel — Raycast style */}
      <nav
        className="max-w-[1200px] mx-auto rounded-2xl flex items-center justify-between px-6 sm:px-8 h-[60px] md:h-[68px]"
        style={{
          backdropFilter: 'blur(12px)',
          WebkitBackdropFilter: 'blur(12px)',
          border: '0.8px solid rgba(255, 255, 255, 0.06)',
          boxShadow: 'rgba(255, 255, 255, 0.12) 0px 1px 1px 0px inset, 0 4px 24px -4px rgba(0, 0, 0, 0.3)',
          background: 'rgba(10, 10, 15, 0.6)',
        }}
      >
        {/* Logo */}
        <a href="#" className="flex items-center gap-2.5">
          <div className="w-7 h-7 bg-gradient-to-br from-sky-400 to-blue-600 rounded-lg flex items-center justify-center shadow-glow-sm">
            <span className="text-white text-xs font-bold">V</span>
          </div>
          <span className="font-semibold text-sm text-white/90">Volo</span>
        </a>

        {/* Desktop links */}
        <div className="hidden md:flex items-center gap-8 text-[13px] text-white/40 font-medium">
          <a href="#features" className="hover:text-white/80 transition-colors">Features</a>
          <a href="#how-it-works" className="hover:text-white/80 transition-colors">How it works</a>
          <a href="https://github.com/OPetricevic/volo" className="hover:text-white/80 transition-colors">GitHub</a>
          <a href="https://github.com/OPetricevic/volo#readme" className="hover:text-white/80 transition-colors">Docs</a>
        </div>

        {/* CTA + Mobile toggle */}
        <div className="flex items-center gap-3">
          <a
            href="#download"
            className="hidden sm:inline-flex px-5 py-2 rounded-xl text-[13px] font-medium text-white/80 hover:text-white transition-all"
            style={{
              background: 'rgba(255, 255, 255, 0.06)',
              border: '0.8px solid rgba(255, 255, 255, 0.08)',
            }}
          >
            Get Volo
          </a>

          {/* Mobile hamburger */}
          <button
            onClick={() => setOpen(!open)}
            className="md:hidden w-8 h-8 flex items-center justify-center text-white/40 hover:text-white/80 transition"
            aria-label="Toggle menu"
            aria-expanded={open}
          >
            {open ? (
              <svg width="18" height="18" viewBox="0 0 18 18" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round">
                <line x1="4" y1="4" x2="14" y2="14" />
                <line x1="14" y1="4" x2="4" y2="14" />
              </svg>
            ) : (
              <svg width="18" height="18" viewBox="0 0 18 18" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round">
                <line x1="3" y1="5" x2="15" y2="5" />
                <line x1="3" y1="9" x2="15" y2="9" />
                <line x1="3" y1="13" x2="15" y2="13" />
              </svg>
            )}
          </button>
        </div>
      </nav>

      {/* Mobile dropdown — also glass */}
      {open && (
        <div
          className="md:hidden max-w-[1200px] mx-auto mt-2 rounded-2xl px-6 py-5 flex flex-col gap-4"
          style={{
            backdropFilter: 'blur(12px)',
            WebkitBackdropFilter: 'blur(12px)',
            border: '0.8px solid rgba(255, 255, 255, 0.06)',
            boxShadow: 'rgba(255, 255, 255, 0.12) 0px 1px 1px 0px inset, 0 4px 24px -4px rgba(0, 0, 0, 0.3)',
            background: 'rgba(10, 10, 15, 0.7)',
          }}
        >
          <a href="#features" onClick={() => setOpen(false)} className="text-sm text-white/50 hover:text-white/90 transition py-1">Features</a>
          <a href="#how-it-works" onClick={() => setOpen(false)} className="text-sm text-white/50 hover:text-white/90 transition py-1">How it works</a>
          <a href="#download" onClick={() => setOpen(false)} className="text-sm text-white/50 hover:text-white/90 transition py-1">Download</a>
          <a
            href="#download"
            onClick={() => setOpen(false)}
            className="mt-2 px-4 py-2.5 bg-gradient-to-r from-sky-500 to-blue-600 text-white text-sm font-medium rounded-xl text-center shadow-glow-sm"
          >
            Download Volo
          </a>
        </div>
      )}
    </div>
  );
}
