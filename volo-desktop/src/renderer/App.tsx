import { useState } from "react";
import { Sidebar } from "./components/Sidebar";
import { ChatView } from "./views/ChatView";
import { HistoryView } from "./views/HistoryView";
import { MacrosView } from "./views/MacrosView";
import { SettingsView } from "./views/SettingsView";

export type View = "chat" | "history" | "macros" | "settings";

export function App() {
  const [activeView, setActiveView] = useState<View>("chat");

  return (
    <div className="flex h-screen bg-[var(--bg)]">
      {/* Sidebar */}
      <Sidebar activeView={activeView} onNavigate={setActiveView} />

      {/* Main content */}
      <main className="flex-1 flex flex-col overflow-hidden bg-[var(--chat)] border-l border-[var(--border)]">
        {/* Titlebar */}
        <div className="drag-region h-9 flex items-center justify-between px-4 shrink-0 border-b border-[var(--border)]">
          <span className="text-[11px] text-[var(--text-faint)]">
            {activeView === "chat" ? "Chat" : activeView === "history" ? "History" : "Settings"}
          </span>
          <div className="no-drag flex items-center gap-0.5">
            <button onClick={() => window.volo.minimize()} className="w-8 h-7 flex items-center justify-center rounded hover:bg-[var(--border)] text-[var(--text-faint)] hover:text-[var(--text-muted)] transition" aria-label="Minimize">
              <svg width="10" height="1" viewBox="0 0 10 1" fill="currentColor"><rect width="10" height="1"/></svg>
            </button>
            <button onClick={() => window.volo.maximize()} className="w-8 h-7 flex items-center justify-center rounded hover:bg-[var(--border)] text-[var(--text-faint)] hover:text-[var(--text-muted)] transition" aria-label="Maximize">
              <svg width="9" height="9" viewBox="0 0 9 9" fill="none" stroke="currentColor" strokeWidth="1.2"><rect x="0.5" y="0.5" width="8" height="8" rx="1"/></svg>
            </button>
            <button onClick={() => window.volo.close()} className="w-8 h-7 flex items-center justify-center rounded hover:bg-[var(--danger)] hover:text-white text-[var(--text-faint)] transition" aria-label="Close">
              <svg width="10" height="10" viewBox="0 0 10 10" fill="none" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round"><line x1="1" y1="1" x2="9" y2="9"/><line x1="9" y1="1" x2="1" y2="9"/></svg>
            </button>
          </div>
        </div>

        {/* View */}
        {activeView === "chat" && <ChatView />}
        {activeView === "history" && <HistoryView />}
        {activeView === "macros" && <MacrosView />}
        {activeView === "settings" && <SettingsView />}
      </main>
    </div>
  );
}
