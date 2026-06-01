export function Download() {
  return (
    <section id="download" className="py-20 md:py-32 border-t border-volo-border/50">
      <div className="max-w-4xl mx-auto px-6 text-center">
        <h2 className="text-3xl md:text-4xl font-bold tracking-tight">
          Ready to try?
        </h2>
        <p className="mt-4 text-volo-muted text-lg max-w-xl mx-auto">
          Free, open source, no account required. Start using Volo in under a minute.
        </p>

        <div className="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4">
          {/* Desktop download */}
          <a
            href="#"
            className="w-full sm:w-auto flex items-center gap-3 px-6 py-4 bg-volo-surface border border-volo-border rounded-xl hover:border-volo-accent/50 transition-all group"
          >
            <div className="w-10 h-10 bg-volo-accent/10 rounded-lg flex items-center justify-center group-hover:bg-volo-accent/20 transition">
              <svg width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" strokeWidth="1.5" className="text-volo-accent">
                <rect x="3" y="2" width="14" height="12" rx="2"/><line x1="6" y1="16" x2="14" y2="16"/><line x1="10" y1="14" x2="10" y2="16"/>
              </svg>
            </div>
            <div className="text-left">
              <div className="text-xs text-volo-muted">Download</div>
              <div className="text-sm font-semibold">Volo for Windows</div>
            </div>
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" className="text-volo-muted ml-4">
              <path d="M8 3v8M5 8l3 3 3-3"/><line x1="3" y1="13" x2="13" y2="13"/>
            </svg>
          </a>

          {/* Extension */}
          <a
            href="#"
            className="w-full sm:w-auto flex items-center gap-3 px-6 py-4 bg-volo-surface border border-volo-border rounded-xl hover:border-volo-accent/50 transition-all group"
          >
            <div className="w-10 h-10 bg-volo-accent/10 rounded-lg flex items-center justify-center group-hover:bg-volo-accent/20 transition">
              <svg width="20" height="20" viewBox="0 0 20 20" fill="none" stroke="currentColor" strokeWidth="1.5" className="text-volo-accent">
                <circle cx="10" cy="10" r="8"/><path d="M10 6v4l2.5 1.5"/>
              </svg>
            </div>
            <div className="text-left">
              <div className="text-xs text-volo-muted">Add to browser</div>
              <div className="text-sm font-semibold">Chrome Extension</div>
            </div>
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" className="text-volo-muted ml-4">
              <path d="M4 8h8M9 5l3 3-3 3"/>
            </svg>
          </a>
        </div>

        <p className="mt-6 text-xs text-volo-muted">
          v0.1.0 · Windows 10+ · Chrome / Opera / Edge
        </p>
      </div>
    </section>
  );
}
