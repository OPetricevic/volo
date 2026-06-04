import type { View } from "../App";

interface SidebarProps {
  activeView: View;
  onNavigate: (view: View) => void;
}

export function Sidebar({ activeView, onNavigate }: SidebarProps) {
  return (
    <aside className="w-[260px] bg-[var(--sidebar)] flex flex-col h-full">
      {/* Workspace header */}
      <div className="px-4 pt-4 pb-3 drag-region">
        <div className="no-drag flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-5 h-5 bg-[var(--accent)] rounded flex items-center justify-center">
              <span className="text-white text-[8px] font-bold">V</span>
            </div>
            <span className="text-[13px] font-semibold text-[var(--text)]">Volo</span>
          </div>
          <div className="flex items-center gap-0.5">
            <button className="w-7 h-7 flex items-center justify-center rounded hover:bg-[var(--border)] text-[var(--text-faint)] hover:text-[var(--text-muted)] transition" aria-label="Search">
              <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M7 2.5a4.5 4.5 0 1 0 0 9 4.5 4.5 0 0 0 0-9zM1 7a6 6 0 1 1 10.7 3.7l3.3 3.3-1 1-3.3-3.3A6 6 0 0 1 1 7z"/></svg>
            </button>
            <button
              onClick={() => onNavigate("chat")}
              className="w-7 h-7 flex items-center justify-center rounded hover:bg-[var(--border)] text-[var(--text-faint)] hover:text-[var(--text-muted)] transition"
              aria-label="New chat"
            >
              <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M8 1.5a.75.75 0 0 1 .75.75v5h5a.75.75 0 0 1 0 1.5h-5v5a.75.75 0 0 1-1.5 0v-5h-5a.75.75 0 0 1 0-1.5h5v-5A.75.75 0 0 1 8 1.5z"/></svg>
            </button>
          </div>
        </div>
      </div>

      {/* Navigation */}
      <div className="px-2 space-y-0.5">
        <NavItem
          icon={<svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M2.5 3A1.5 1.5 0 0 1 4 1.5h8A1.5 1.5 0 0 1 13.5 3v7A1.5 1.5 0 0 1 12 11.5H5.7L3 14V11.5h-.5A1.5 1.5 0 0 1 1 10V3.5 3zm2 0v7h1v1.5l1.7-1.5H12V3H4.5z"/></svg>}
          label="Chat"
          active={activeView === "chat"}
          onClick={() => onNavigate("chat")}
        />
        <NavItem
          icon={<svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M8 1.5a6.5 6.5 0 1 0 0 13 6.5 6.5 0 0 0 0-13zM0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8zm8.5-3.5v3.2l2.4 1.4-.7 1.3L7.5 9V4.5h1z"/></svg>}
          label="History"
          active={activeView === "history"}
          onClick={() => onNavigate("history")}
        />
        <NavItem
          icon={<svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" transform="scale(0.58) translate(3,1)"/></svg>}
          label="Macros"
          active={activeView === "macros"}
          onClick={() => onNavigate("macros")}
        />
        <NavItem
          icon={<svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M8 1a3 3 0 0 0-3 3v4a3 3 0 1 0 6 0V4a3 3 0 0 0-3-3zM3.5 7.5a.75.75 0 0 1 .75.75A3.75 3.75 0 0 0 8 12a3.75 3.75 0 0 0 3.75-3.75.75.75 0 0 1 1.5 0A5.25 5.25 0 0 1 8.75 13.4v1.35h1.5a.75.75 0 0 1 0 1.5h-4.5a.75.75 0 0 1 0-1.5h1.5V13.4A5.25 5.25 0 0 1 2.75 8.25a.75.75 0 0 1 .75-.75z"/></svg>}
          label="Voice Log"
          active={false}
          onClick={() => {}}
        />
      </div>

      {/* Conversation list */}
      <div className="flex-1 overflow-y-auto px-2 mt-4">
        <div className="px-2 pb-1.5">
          <span className="text-[10px] font-semibold text-[var(--text-faint)] uppercase tracking-wider">Recent</span>
        </div>
        <ConvItem label="What did I listen to yesterday" active />
        <ConvItem label="React hooks best practices" />
        <ConvItem label="Open YouTube lofi beats" />
        <ConvItem label="Tailwind vs CSS modules" />

        <div className="px-2 pt-5 pb-1.5">
          <span className="text-[10px] font-semibold text-[var(--text-faint)] uppercase tracking-wider">Yesterday</span>
        </div>
        <ConvItem label="GitHub Actions workflow" />
        <ConvItem label="Docker compose debugging" />
        <ConvItem label="Weather Vienna" />

        <div className="px-2 pt-5 pb-1.5">
          <span className="text-[10px] font-semibold text-[var(--text-faint)] uppercase tracking-wider">Previous 7 days</span>
        </div>
        <ConvItem label="Go error handling patterns" />
        <ConvItem label="Spotify jazz playlist" />
      </div>

      {/* Bottom — settings */}
      <div className="p-2 border-t border-[var(--border)]">
        <NavItem
          icon={<svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M8 5.5a2.5 2.5 0 1 0 0 5 2.5 2.5 0 0 0 0-5zM4 8a4 4 0 1 1 8 0 4 4 0 0 1-8 0z"/><path d="M9.4 1.2a.75.75 0 0 0-1.4-.1L7.1 3H5.5a.75.75 0 0 0-.6.3l-1 1.4-1.8.2a.75.75 0 0 0-.6 1l.6 1.8-.9 1.6a.75.75 0 0 0 .3 1l1.6 1 .2 1.8a.75.75 0 0 0 1 .6l1.8-.6 1.6.9a.75.75 0 0 0 1-.3l1-1.6 1.8.2a.75.75 0 0 0 .6-1l-.6-1.8.9-1.6a.75.75 0 0 0-.3-1l-1.6-1-.2-1.8a.75.75 0 0 0-1-.6L9.5 3.5l-.1-2.3z"/></svg>}
          label="Settings"
          active={activeView === "settings"}
          onClick={() => onNavigate("settings")}
        />
      </div>
    </aside>
  );
}

function NavItem({ icon, label, active, onClick }: { icon: React.ReactNode; label: string; active: boolean; onClick: () => void }) {
  return (
    <button
      onClick={onClick}
      className={`w-full flex items-center gap-2.5 px-3 py-1.5 rounded-md text-[13px] font-medium transition ${
        active
          ? "bg-[var(--border)] text-[var(--text)]"
          : "text-[var(--text-muted)] hover:bg-[rgba(255,255,255,0.03)] hover:text-[var(--text)]"
      }`}
    >
      <span className={active ? "text-[var(--text)]" : "text-[var(--text-faint)]"}>{icon}</span>
      {label}
    </button>
  );
}

function ConvItem({ label, active }: { label: string; active?: boolean }) {
  return (
    <div className={`px-3 py-1.5 rounded-md text-[12px] cursor-pointer transition truncate ${
      active
        ? "bg-[var(--border)] text-[var(--text)]"
        : "text-[var(--text-muted)] hover:bg-[rgba(255,255,255,0.03)] hover:text-[var(--text)]"
    }`}>
      {label}
    </div>
  );
}
