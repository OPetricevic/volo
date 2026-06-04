export function Download() {
  return (
    <section id="download" className="py-20 md:py-32 relative">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 text-center relative">
        <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-white">
          Ready to try?
        </h2>
        <p className="mt-4 text-white/50 text-base md:text-lg max-w-xl mx-auto">
          Free, open source, no account required. Start using Volo in under a minute.
        </p>

        <div className="mt-10 flex flex-col sm:flex-row items-center justify-center gap-3 sm:gap-4">
          {/* Chrome extension — primary CTA */}
          <a
            href="https://github.com/OPetricevic/volo/releases/latest"
            className="w-full sm:w-auto inline-flex items-center justify-center gap-2.5 px-7 py-3.5 bg-white text-[#060609] font-semibold rounded-full hover:shadow-[0_0_30px_rgba(255,255,255,0.15)] transition-all text-center text-[15px]"
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
              <circle cx="8" cy="8" r="7" stroke="currentColor" strokeWidth="1.2"/>
              <circle cx="8" cy="8" r="3" fill="currentColor" opacity="0.9"/>
              <path d="M8 1a7 7 0 0 1 4.95 2.05L8 8" stroke="currentColor" strokeWidth="1.2" fill="none"/>
            </svg>
            Add to Browser
          </a>

          {/* Desktop download — secondary */}
          <a
            href="https://github.com/OPetricevic/volo/releases/latest"
            className="w-full sm:w-auto inline-flex items-center justify-center gap-2.5 px-7 py-3.5 glass-card text-white/90 font-medium rounded-full hover:text-white hover:border-white/[0.15] transition-all text-center text-[15px]"
          >
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <path d="M0 2.3l6.5-.9v6.3H0V2.3zm7.3-1L15.8.1v7.6H7.3V1.3zM15.8 8.4V16l-8.5-1.2V8.4h8.5zM6.5 14.7L0 13.8V8.4h6.5v6.3z"/>
            </svg>
            Download for Windows
          </a>
        </div>

        <p className="mt-6 text-xs text-white/20">
          Windows 10+ · Chrome / Opera / Edge · MIT License
        </p>

        {/* Requirements */}
        <div className="mt-12 max-w-lg mx-auto">
          <div className="grid grid-cols-2 gap-4">
            <div className="glass-card rounded-xl p-4 text-center">
              <div className="w-8 h-8 mx-auto mb-2 rounded-lg bg-sky-500/[0.08] border border-sky-500/20 flex items-center justify-center">
                <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
                  <circle cx="8" cy="8" r="7" stroke="url(#ext-g)" strokeWidth="1.2"/>
                  <circle cx="8" cy="8" r="3" fill="url(#ext-g)" opacity="0.8"/>
                  <path d="M8 1a7 7 0 0 1 4.95 2.05L8 8" stroke="url(#ext-g)" strokeWidth="1.2"/>
                  <defs><linearGradient id="ext-g" x1="1" y1="1" x2="15" y2="15" gradientUnits="userSpaceOnUse"><stop stopColor="#38bdf8"/><stop offset="1" stopColor="#0ea5e9"/></linearGradient></defs>
                </svg>
              </div>
              <p className="text-[13px] font-medium text-white/80">Extension</p>
              <p className="text-[11px] text-white/35 mt-1">Chrome, Opera, or Edge</p>
              <p className="text-[11px] text-white/35">No other requirements</p>
            </div>
            <div className="glass-card rounded-xl p-4 text-center">
              <div className="w-8 h-8 mx-auto mb-2 rounded-lg bg-sky-500/[0.08] border border-sky-500/20 flex items-center justify-center">
                <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" className="text-sky-400">
                  <path d="M0 2.3l6.5-.9v6.3H0V2.3zm7.3-1L15.8.1v7.6H7.3V1.3zM15.8 8.4V16l-8.5-1.2V8.4h8.5zM6.5 14.7L0 13.8V8.4h6.5v6.3z"/>
                </svg>
              </div>
              <p className="text-[13px] font-medium text-white/80">Desktop App</p>
              <p className="text-[11px] text-white/35 mt-1">16GB RAM · Windows 10+</p>
              <p className="text-[11px] text-white/35">~3GB disk for AI model</p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
