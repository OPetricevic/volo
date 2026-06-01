import { useState } from "react";
import { Sidebar } from "./components/Sidebar";
import { ChatView } from "./views/ChatView";
import { HistoryView } from "./views/HistoryView";
import { SettingsView } from "./views/SettingsView";

export type View = "chat" | "history" | "settings";

export function App() {
  const [activeView, setActiveView] = useState<View>("chat");

  return (
    <div className="flex h-screen bg-[var(--bg)]">
      {/* Sidebar */}
      <Sidebar activeView={activeView} onNavigate={setActiveView} />

      {/* Main content */}
      <main className="flex-1 flex flex-col overflow-hidden bg-[var(--chat)]">
        {/* Titlebar — minimal, just window controls */}
        <div className="drag-region h-9 flex items-center justify-end px-2 shrink-0">
          <div className="no-drag flex items-center">
            <button onClick={() => window.volo.minimize()} className="w-11 h-9 flex items-center justify-center hover:bg-[var(--input)]/40 text-[var(--text-muted)]" aria-label="Minimize">
              <svg width="10" height="1" viewBox="0 0 10 1" fill="currentColor"><rect width="10" height="1"/></svg>
            </button>
            <button onClick={() => window.volo.maximize()} className="w-11 h-9 flex items-center justify-center hover:bg-[var(--input)]/40 text-[var(--text-muted)]" aria-label="Maximize">
              <svg width="9" height="9" viewBox="0 0 9 9" fill="none" stroke="currentColor" strokeWidth="1.2"><rect x="0.5" y="0.5" width="8" height="8" rx="1"/></svg>
            </button>
            <button onClick={() => window.volo.close()} className="w-11 h-9 flex items-center justify-center hover:bg-[#da373c] hover:text-white text-[var(--text-muted)]" aria-label="Close">
              <svg width="10" height="10" viewBox="0 0 10 10" fill="none" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round"><line x1="1" y1="1" x2="9" y2="9"/><line x1="9" y1="1" x2="1" y2="9"/></svg>
            </button>
          </div>
        </div>

        {/* View */}
        {activeView === "chat" && <ChatView />}
        {activeView === "history" && <HistoryView />}
        {activeView === "settings" && <SettingsView />}
      </main>
    </div>
  );
}
