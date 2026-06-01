export function Nav() {
  return (
    <nav className="fixed top-0 inset-x-0 z-50 border-b border-volo-border/50 bg-volo-bg/80 backdrop-blur-xl">
      <div className="max-w-6xl mx-auto px-6 h-14 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="w-7 h-7 bg-volo-accent rounded-lg flex items-center justify-center">
            <span className="text-white text-xs font-bold">V</span>
          </div>
          <span className="font-semibold text-sm">Volo</span>
        </div>

        <div className="hidden md:flex items-center gap-8 text-sm text-volo-muted">
          <a href="#features" className="hover:text-volo-text transition">Features</a>
          <a href="#how-it-works" className="hover:text-volo-text transition">How it works</a>
          <a href="#download" className="hover:text-volo-text transition">Download</a>
        </div>

        <a
          href="#download"
          className="px-4 py-1.5 bg-volo-surface-2 border border-volo-border rounded-full text-xs font-medium text-volo-text hover:bg-volo-accent hover:border-volo-accent hover:text-white transition-all"
        >
          Download
        </a>
      </div>
    </nav>
  );
}
