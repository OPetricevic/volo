export function Features() {
  return (
    <section id="features" className="py-24 md:py-36 border-t border-white/[0.04]">
      <div className="max-w-6xl mx-auto px-6">
        <div className="text-center mb-20">
          <h2 className="text-3xl md:text-4xl font-bold tracking-tight">
            Everything you need.
            <br />
            <span className="text-volo-muted">Nothing you don't.</span>
          </h2>
        </div>

        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-px bg-white/[0.04] rounded-2xl overflow-hidden border border-white/[0.04]">
          <FeatureCard
            icon={<VoiceIcon />}
            title="Voice Commands"
            description="Say 'Hey Volo' followed by any command. Search, navigate, open apps, control tabs — all by voice."
          />
          <FeatureCard
            icon={<HistoryIcon />}
            title="AI History Chat"
            description="Ask Volo about your browsing history in natural language. 'What did I listen to yesterday?' — it knows."
          />
          <FeatureCard
            icon={<BoltIcon />}
            title="Instant Execution"
            description="Commands execute in under 10ms. No loading, no waiting. Speak and it's done."
          />
          <FeatureCard
            icon={<LockIcon />}
            title="Privacy First"
            description="Voice processing happens locally. Your data stays on your machine. Optional cloud sync."
          />
          <FeatureCard
            icon={<LayersIcon />}
            title="Cross-Platform"
            description="Browser extension for Chrome, Opera, Edge. Desktop app for Windows. Same experience everywhere."
          />
          <FeatureCard
            icon={<WaveIcon />}
            title="Learns Your Patterns"
            description="The more you use it, the smarter it gets. Frequently used commands surface faster over time."
          />
        </div>
      </div>
    </section>
  );
}

function FeatureCard({ icon, title, description }: { icon: React.ReactNode; title: string; description: string }) {
  return (
    <div className="group p-8 bg-[#0a0a0f] hover:bg-[#0f0f14] transition-colors">
      <div className="mb-5">
        {icon}
      </div>
      <h3 className="text-[15px] font-semibold mb-2 text-white/90">{title}</h3>
      <p className="text-[13px] text-white/40 leading-relaxed">{description}</p>
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
      <defs><linearGradient id="voice-grad" x1="7" y1="4" x2="25" y2="28" gradientUnits="userSpaceOnUse"><stop stopColor="#5865f2"/><stop offset="1" stopColor="#0ea5e9"/></linearGradient></defs>
    </svg>
  );
}

function HistoryIcon() {
  return (
    <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
      <circle cx="16" cy="16" r="11" stroke="url(#hist-grad)" strokeWidth="1.5"/>
      <path d="M16 9v7l4.5 2.5" stroke="url(#hist-grad)" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
      <circle cx="16" cy="16" r="2" fill="url(#hist-grad)" opacity="0.4"/>
      <defs><linearGradient id="hist-grad" x1="5" y1="5" x2="27" y2="27" gradientUnits="userSpaceOnUse"><stop stopColor="#5865f2"/><stop offset="1" stopColor="#0ea5e9"/></linearGradient></defs>
    </svg>
  );
}

function BoltIcon() {
  return (
    <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
      <path d="M18 4L8 18h7l-1 10 10-14h-7l1-10z" stroke="url(#bolt-grad)" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
      <defs><linearGradient id="bolt-grad" x1="8" y1="4" x2="24" y2="28" gradientUnits="userSpaceOnUse"><stop stopColor="#5865f2"/><stop offset="1" stopColor="#0ea5e9"/></linearGradient></defs>
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
      <defs><linearGradient id="lock-grad" x1="8" y1="6" x2="24" y2="26" gradientUnits="userSpaceOnUse"><stop stopColor="#5865f2"/><stop offset="1" stopColor="#0ea5e9"/></linearGradient></defs>
    </svg>
  );
}

function LayersIcon() {
  return (
    <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
      <path d="M16 6L4 12l12 6 12-6L16 6z" stroke="url(#layer-grad)" strokeWidth="1.5" strokeLinejoin="round"/>
      <path d="M4 18l12 6 12-6" stroke="url(#layer-grad)" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" opacity="0.6"/>
      <path d="M4 24l12 6 12-6" stroke="url(#layer-grad)" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" opacity="0.3"/>
      <defs><linearGradient id="layer-grad" x1="4" y1="6" x2="28" y2="30" gradientUnits="userSpaceOnUse"><stop stopColor="#5865f2"/><stop offset="1" stopColor="#0ea5e9"/></linearGradient></defs>
    </svg>
  );
}

function WaveIcon() {
  return (
    <svg width="32" height="32" viewBox="0 0 32 32" fill="none">
      <path d="M4 16c2-4 4-8 6-8s4 16 6 16 4-16 6-16 4 8 6 8" stroke="url(#wave-grad)" strokeWidth="1.5" strokeLinecap="round"/>
      <defs><linearGradient id="wave-grad" x1="4" y1="8" x2="28" y2="24" gradientUnits="userSpaceOnUse"><stop stopColor="#5865f2"/><stop offset="1" stopColor="#0ea5e9"/></linearGradient></defs>
    </svg>
  );
}
