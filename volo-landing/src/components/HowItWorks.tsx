export function HowItWorks() {
  return (
    <section id="how-it-works" className="py-20 md:py-32 relative">
      <div className="max-w-4xl mx-auto px-4 sm:px-6">
        <div className="text-center mb-12 md:mb-16">
          <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-white/95">
            Three steps. That's it.
          </h2>
          <p className="mt-4 text-white/35 text-base md:text-lg">No setup wizards. No configuration. Install and talk.</p>
        </div>

        {/* Glass container for steps */}
        <div className="glass-card glass-highlight rounded-2xl overflow-hidden">
          <Step
            number="01"
            title="Install"
            description="Add the Chrome extension or download the desktop app. Grant microphone access when prompted — that's the only permission needed."
            last={false}
          />
          <Step
            number="02"
            title={`Say "Hey Volo"`}
            description="The wake word activates listening. Then speak naturally — no specific phrasing required. Volo understands intent, not just keywords."
            last={false}
          />
          <Step
            number="03"
            title="Done"
            description="Volo executes immediately. Search results open, sites navigate, tabs close. Rule-based parsing in 1ms, ML fallback in 8ms."
            last={true}
          />
        </div>
      </div>
    </section>
  );
}

function Step({ number, title, description, last }: { number: string; title: string; description: string; last: boolean }) {
  return (
    <div className={`flex gap-4 sm:gap-6 p-6 sm:p-8 ${!last ? 'border-b border-white/[0.06]' : ''} hover:bg-white/[0.02] transition-colors`}>
      <div className="text-2xl sm:text-3xl font-bold text-sky-400/30 font-mono w-10 sm:w-12 flex-shrink-0">
        {number}
      </div>
      <div>
        <h3 className="text-lg sm:text-xl font-semibold mb-2 text-white/85">{title}</h3>
        <p className="text-sm sm:text-base text-white/35 leading-relaxed">{description}</p>
      </div>
    </div>
  );
}
