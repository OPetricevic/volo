import type { View } from "../App";

interface SidebarProps {
  activeView: View;
  onNavigate: (view: View) => void;
}

export function Sidebar({ activeView, onNavigate }: SidebarProps) {
  return (
    <aside className="w-[260px] bg-[var(--sidebar)] flex flex-col h-full border-r border-[var(--border)]/20">
      {/* New chat button */}
      <div className="p-3 drag-region">
        <button
          onClick={() => onNavigate("chat")}
          className="no-drag w-full flex items-center gap-2 px-3 py-2.5 rounded-lg text-[13px] text-[var(--text-muted)] hover:bg-[var(--input)]/40 transition"
        >
          <svg width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"><line x1="8" y1="3" x2="8" y2="13"/><line x1="3" y1="8" x2="13" y2="8"/></svg>
          New chat
        </button>
      </div>

      {/* Conversation list */}
      <div className="flex-1 overflow-y-auto px-2">
        {/* Today */}
        <div className="px-2 pt-4 pb-1">
          <span className="text-[11px] font-medium text-[var(--text-faint)]">Today</span>
        </div>
        <ConvItem label="What did I listen to yesterday" active />
        <ConvItem label="React hooks search" />
        <ConvItem label="Open YouTube lofi beats" />

        {/* Yesterday */}
        <div className="px-2 pt-5 pb-1">
          <span className="text-[11px] font-medium text-[var(--text-faint)]">Yesterday</span>
        </div>
        <ConvItem label="GitHub repos for Go" />
        <ConvItem label="Weather check" />

        {/* Older */}
        <div className="px-2 pt-5 pb-1">
          <span className="text-[11px] font-medium text-[var(--text-faint)]">Previous 7 days</span>
        </div>
        <ConvItem label="Spotify jazz playlist" />
        <ConvItem label="Close all tabs command" />
      </div>

      {/* Bottom actions */}
      <div className="p-2 border-t border-[var(--border)]/20">
        <button
          onClick={() => onNavigate("history")}
          className={`w-full flex items-center gap-2 px-3 py-2 rounded-lg text-[13px] transition ${
            activeView === "history" ? "bg-[var(--input)]/50 text-[var(--text)]" : "text-[var(--text-muted)] hover:bg-[var(--input)]/30"
          }`}
        >
          <svg width="15" height="15" viewBox="0 0 15 15" fill="none" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round"><circle cx="7.5" cy="7.5" r="6"/><path d="M7.5 4v3.5l2.5 1.5"/></svg>
          All history
        </button>
        <button
          onClick={() => onNavigate("settings")}
          className={`w-full flex items-center gap-2 px-3 py-2 rounded-lg text-[13px] transition ${
            activeView === "settings" ? "bg-[var(--input)]/50 text-[var(--text)]" : "text-[var(--text-muted)] hover:bg-[var(--input)]/30"
          }`}
        >
          <svg width="15" height="15" viewBox="0 0 15 15" fill="none" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round"><circle cx="7.5" cy="7.5" r="2"/><path d="M7.5 1.5v2M7.5 11.5v2M1.5 7.5h2M11.5 7.5h2"/></svg>
          Settings
        </button>
      </div>
    </aside>
  );
}

function ConvItem({ label, active }: { label: string; active?: boolean }) {
  return (
    <div className={`px-3 py-2 rounded-lg text-[13px] cursor-pointer transition truncate ${
      active
        ? "bg-[var(--input)]/50 text-[var(--text)]"
        : "text-[var(--text-muted)] hover:bg-[var(--input)]/25 hover:text-[var(--text)]"
    }`}>
      {label}
    </div>
  );
}
