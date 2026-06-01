/**
 * Custom window titlebar (frameless window).
 * Draggable area + window control buttons.
 */
export function Titlebar() {
  return (
    <div className="drag-region flex items-center justify-between h-8 bg-[var(--surface)] border-b border-[var(--border)] px-3 select-none">
      <span className="text-xs text-[var(--text-muted)] font-medium">Volo</span>

      <div className="no-drag flex items-center gap-1">
        <button
          onClick={() => window.volo.minimize()}
          className="w-6 h-6 flex items-center justify-center rounded hover:bg-[var(--surface-2)] text-[var(--text-muted)] hover:text-[var(--text)] transition"
          aria-label="Minimize"
        >
          <svg width="10" height="1" viewBox="0 0 10 1" fill="currentColor">
            <rect width="10" height="1" />
          </svg>
        </button>
        <button
          onClick={() => window.volo.maximize()}
          className="w-6 h-6 flex items-center justify-center rounded hover:bg-[var(--surface-2)] text-[var(--text-muted)] hover:text-[var(--text)] transition"
          aria-label="Maximize"
        >
          <svg width="9" height="9" viewBox="0 0 9 9" fill="none" stroke="currentColor" strokeWidth="1.2">
            <rect x="0.5" y="0.5" width="8" height="8" rx="1" />
          </svg>
        </button>
        <button
          onClick={() => window.volo.close()}
          className="w-6 h-6 flex items-center justify-center rounded hover:bg-[var(--danger)] text-[var(--text-muted)] hover:text-white transition"
          aria-label="Close"
        >
          <svg width="10" height="10" viewBox="0 0 10 10" fill="none" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round">
            <line x1="1" y1="1" x2="9" y2="9" />
            <line x1="9" y1="1" x2="1" y2="9" />
          </svg>
        </button>
      </div>
    </div>
  );
}
