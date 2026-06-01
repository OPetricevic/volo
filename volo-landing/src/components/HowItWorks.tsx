export function HowItWorks() {
  return (
    <section id="how-it-works" className="py-20 md:py-32 border-t border-volo-border/50">
      <div className="max-w-4xl mx-auto px-6">
        <div className="text-center mb-16">
          <h2 className="text-3xl md:text-4xl font-bold tracking-tight">
            Three steps. That's it.
          </h2>
          <p className="mt-4 text-volo-muted text-lg">No setup wizards. No configuration. Just install and talk.</p>
        </div>

        <div className="space-y-0">
          <Step
            number="01"
            title="Install"
            description="Add the Chrome extension or download the desktop app. Grant microphone access when prompted."
          />
          <Step
            number="02"
            title='Say "Hey Volo"'
            description="The wake word activates listening. Then say your command naturally — no specific phrasing required."
          />
          <Step
            number="03"
            title="Done"
            description="Volo executes immediately. Search results open, sites navigate, tabs close. Under 10ms."
          />
        </div>
      </div>
    </section>
  );
}

function Step({ number, title, description }: { number: string; title: string; description: string }) {
  return (
    <div className="flex gap-6 py-8 border-b border-volo-border/50 last:border-0">
      <div className="text-3xl font-bold text-volo-accent/30 font-mono w-12 flex-shrink-0">
        {number}
      </div>
      <div>
        <h3 className="text-xl font-semibold mb-2">{title}</h3>
        <p className="text-volo-muted leading-relaxed">{description}</p>
      </div>
    </div>
  );
}
