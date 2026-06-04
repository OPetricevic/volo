export function Features() {
  return (
    <section id="features" className="py-20 md:py-36 relative ambient-glow">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 relative">
        <div className="text-center mb-14 md:mb-20">
          <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-white/95">
            Everything you need.
            <br />
            <span className="text-white/30">Nothing you don't.</span>
          </h2>
        </div>

        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-5">
          <FeatureCard
            icon={<VoiceIcon />}
            title="Voice Commands"
            description={`Say "Hey Volo" followed by any command. Search, navigate, open apps, control tabs — hands-free.`}
          />
          <FeatureCard
            icon={<HistoryIcon />}
            title="AI History Chat"
            description={`Ask about your browsing history naturally. "What did I listen to yesterday?" — it knows.`}
          />
          <FeatureCard
            icon={<BoltIcon />}
            title="Instant Execution"
            description="Commands execute in under 10ms. No loading spinners, no waiting. Speak and it's done."
          />
          <FeatureCard
            icon={<LockIcon />}
            title="Privacy First"
            description="Voice processing happens locally via Web Speech API. Your data stays on your machine."
          />
          <FeatureCard
            icon={<LayersIcon />}
            title="Cross-Platform"
            description="Browser extension for Chrome, Opera, Edge. Desktop app for Windows. Same experience everywhere."
          />
          <FeatureCard
            icon={<WaveIcon />}
            title="Learns Your Patterns"
            description="Frequently used commands surface faster. The hybrid intent parser adapts to how you speak."
          />
        </div>
      </div>
    </section>
  );
}

function FeatureCard({ icon, title, description }: { icon: React.ReactNode; title: string; description: string }) {
  return (
    <div className="glass-card glass-card-hover glass-highlight rounded-2xl p-6 sm:p-8 transition-all duration-300 group">
      <div className="mb-4 sm:mb-5 opacity-80 group-hover:opacity-100 transition-opacity">
        {icon}
      </div>
      <h3 className="text-[14px] sm:text-[15px] font-semibold mb-2 text-white/85 group-hover:text-white/95 transition-colors">{title}</h3>
      <p className="text-[12px] sm:text-[13px] text-white/35 leading-relaxed group-hover:text-white/45 transition-colors">{description}</p>
    </div>
  );
}

// ─── Premium minimal SVG icons ───────────────────────────

function VoiceIcon() {
  return (
    <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
      <rect x="12" y="4" width="8" height="14" rx="4" stroke="url(#voice-grad)" strokeWidth="1.5"/>
      <path d="M7 15a9 9 0 0 0 18 0" stroke="url(#voice-grad)" strokeWidth="1.5" strokeLinecap="round"/>
      <line x1="16" y1="24" x2="16" y2="28" stroke="url(#voice-grad)" strokeWidth="1.5" strokeLinecap="round"/>
      <line x1="12" y1="28" x2="20" y2="28" stroke="url(#voice-grad)" strokeWidth="1.5" strokeLinecap="round"/>
      <defs><linearGradient id="voice-grad" x1="7" y1="4" x2="25" y2="28" gradientUnits="userSpaceOnUse"><stop stopColor="#38bdf8"/><stop offset="1" stopColor="#0ea5e9"/></linearGradient></defs>
    </svg>
  );
}

function HistoryIcon() {
  return (
    <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
      <circle cx="16" cy="16" r="11" stroke="url(#hist-grad)" strokeWidth="1.5"/>
      <path d="M16 9v7l4.5 2.5" stroke="url(#hist-grad)" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
      <circle cx="16" cy="16" r="2" fill="url(#hist-grad)" opacity="0.4"/>
      <defs><linearGradient id="hist-grad" x1="5" y1="5" x2="27" y2="27" gradientUnits="userSpaceOnUse"><stop stopColor="#38bdf8"/><stop offset="1" stopColor="#0ea5e9"/></linearGradient></defs>
    </svg>
  );
}

function BoltIcon() {
  return (
    <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
      <path d="M18 4L8 18h7l-1 10 10-14h-7l1-10z" stroke="url(#bolt-grad)" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
      <defs><linearGradient id="bolt-grad" x1="8" y1="4" x2="24" y2="28" gradientUnits="userSpaceOnUse"><stop stopColor="#38bdf8"/><stop offset="1" stopColor="#0ea5e9"/></linearGradient></defs>
    </svg>
  );
}

function LockIcon() {
  return (
    <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
      <rect x="8" y="14" width="16" height="12" rx="3" stroke="url(#lock-grad)" strokeWidth="1.5"/>
      <path d="M11 14v-3a5 5 0 0 1 10 0v3" stroke="url(#lock-grad)" strokeWidth="1.5" strokeLinecap="round"/>
      <circle cx="16" cy="20" r="2" fill="url(#lock-grad)" opacity="0.6"/>
      <line x1="16" y1="22" x2="16" y2="24" stroke="url(#lock-grad)" strokeWidth="1.5" strokeLinecap="round"/>
      <defs><linearGradient id="lock-grad" x1="8" y1="6" x2="24" y2="26" gradientUnits="userSpaceOnUse"><stop stopColor="#38bdf8"/><stop offset="1" stopColor="#0ea5e9"/></linearGradient></defs>
    </svg>
  );
}

function LayersIcon() {
  return (
    <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
      <path d="M16 6L4 12l12 6 12-6L16 6z" stroke="url(#layer-grad)" strokeWidth="1.5" strokeLinejoin="round"/>
      <path d="M4 18l12 6 12-6" stroke="url(#layer-grad)" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" opacity="0.6"/>
      <path d="M4 24l12 6 12-6" stroke="url(#layer-grad)" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" opacity="0.3"/>
      <defs><linearGradient id="layer-grad" x1="4" y1="6" x2="28" y2="30" gradientUnits="userSpaceOnUse"><stop stopColor="#38bdf8"/><stop offset="1" stopColor="#0ea5e9"/></linearGradient></defs>
    </svg>
  );
}

function WaveIcon() {
  return (
    <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
      <path d="M4 16c2-4 4-8 6-8s4 16 6 16 4-16 6-16 4 8 6 8" stroke="url(#wave-grad)" strokeWidth="1.5" strokeLinecap="round"/>
      <defs><linearGradient id="wave-grad" x1="4" y1="8" x2="28" y2="24" gradientUnits="userSpaceOnUse"><stop stopColor="#38bdf8"/><stop offset="1" stopColor="#0ea5e9"/></linearGradient></defs>
    </svg>
  );
}
