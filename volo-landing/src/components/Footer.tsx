export function Footer() {
  return (
    <footer className="border-t border-volo-border/50 py-12">
      <div className="max-w-6xl mx-auto px-6">
        <div className="flex flex-col md:flex-row items-center justify-between gap-6">
          <div className="flex items-center gap-2">
            <div className="w-6 h-6 bg-volo-accent rounded-md flex items-center justify-center">
              <span className="text-white text-[9px] font-bold">V</span>
            </div>
            <span className="text-sm font-semibold">Volo</span>
            <span className="text-xs text-volo-muted ml-2">Voice-first browser assistant</span>
          </div>

          <div className="flex items-center gap-6 text-sm text-volo-muted">
            <a href="https://github.com/volo" className="hover:text-volo-text transition">GitHub</a>
            <a href="#" className="hover:text-volo-text transition">Documentation</a>
            <a href="#" className="hover:text-volo-text transition">Privacy</a>
          </div>
        </div>

        <div className="mt-8 text-center text-xs text-volo-muted/60">
          Bachelor Thesis Project · 2026
        </div>
      </div>
    </footer>
  );
}
