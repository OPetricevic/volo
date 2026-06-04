import { useState, useEffect } from "react";

// Pre-built conversations for each "recent" item
const CONVERSATIONS: Record<string, Array<{ role: "user" | "assistant"; text: string }>> = {
  "What did I listen to yesterday": [
    { role: "user", text: "What did I listen to yesterday?" },
    { role: "assistant", text: "Yesterday you played:\n• lofi beats on YouTube — 9:14am\n• jazz playlist on Spotify — 2:30pm\n• ambient music on YouTube — 8:45pm" },
  ],
  "React hooks best practices": [
    { role: "user", text: "What did I search about React hooks?" },
    { role: "assistant", text: "You searched for:\n• 'react hooks best practices 2025'\n• 'useEffect cleanup explained'\n• 'custom hooks patterns'\n\nMostly on Google and Stack Overflow." },
  ],
  "Open YouTube lofi beats": [
    { role: "user", text: "Open YouTube lofi beats" },
    { role: "assistant", text: "I can't open pages directly — use voice mode for that!\n\nSay: \"Hey Volo, open YouTube lofi beats\"\n\nOr press Ctrl+Shift+V to activate." },
  ],
  "Tailwind vs CSS modules": [
    { role: "user", text: "What was I reading about Tailwind?" },
    { role: "assistant", text: "You visited 3 articles comparing Tailwind vs CSS Modules:\n• dev.to — 'Why I switched to Tailwind'\n• css-tricks.com — 'Modular CSS in 2025'\n• reddit.com/r/webdev — discussion thread" },
  ],
  "GitHub Actions workflow": [
    { role: "user", text: "Show me my GitHub activity" },
    { role: "assistant", text: "Yesterday you visited:\n• github.com/OPetricevic/volo — 6 times\n• GitHub Actions docs — workflow syntax\n• Stack Overflow — 'github actions ssh deploy'" },
  ],
  "Docker compose debugging": [
    { role: "user", text: "What Docker stuff did I look at?" },
    { role: "assistant", text: "You searched:\n• 'docker compose volume not mounting'\n• 'postgres healthcheck docker'\n\nVisited docs.docker.com and a Stack Overflow answer." },
  ],
  "Weather Vienna": [
    { role: "user", text: "Did I check the weather?" },
    { role: "assistant", text: "Yes — you opened weather.com/Vienna at 7:42am yesterday." },
  ],
};

const RECENT_ITEMS = [
  { label: "What did I listen to yesterday", time: "2m" },
  { label: "React hooks best practices", time: "1h" },
  { label: "Open YouTube lofi beats", time: "3h" },
  { label: "Tailwind vs CSS modules", time: "5h" },
  { label: "GitHub Actions workflow", time: "18h" },
  { label: "Docker compose debugging", time: "20h" },
  { label: "Weather Vienna", time: "1d" },
];

export function ProductShowcase() {
  return (
    <section className="py-16 md:py-32 relative">
      <div className="max-w-6xl mx-auto px-4 sm:px-6">

        {/* ─── EXTENSION FIRST — the main product ─── */}
        <div className="grid md:grid-cols-2 gap-10 md:gap-12 items-center">
          <div>
            <h2 className="text-2xl md:text-3xl font-bold tracking-tight text-white/95">Voice-control your browser.</h2>
            <p className="mt-4 text-white/40 leading-relaxed">The extension listens for "Hey Volo" on every page. Say a command and it executes instantly — no clicks, no typing.</p>
            <div className="mt-8 space-y-4">
              <CommandExample command={`"Hey Volo, open YouTube play lofi"`} result="Opens YouTube, plays first result" />
              <CommandExample command={`"Hey Volo, search React hooks"`} result="Google search in new tab" />
              <CommandExample command={`"Hey Volo, go to GitHub"`} result="Navigates to github.com" />
              <CommandExample command={`"Hey Volo, close tab"`} result="Closes current tab" />
            </div>
          </div>

          {/* Extension popup mockup */}
          <div className="flex justify-center">
            <div className="relative">
              <div className="absolute -inset-6 rounded-2xl pointer-events-none" style={{
                background: 'radial-gradient(70% 70% at 50% 50%, rgba(30, 32, 38, 0.4) 0%, transparent 100%)',
              }} />
              <div className="w-[280px] sm:w-[300px] rounded-xl overflow-hidden relative" style={{
                background: '#0a0a0c',
                boxShadow: '0 0 0 0.5px rgba(255, 255, 255, 0.04)',
              }}>
                <div className="flex items-center gap-2 px-4 pt-4 pb-3">
                  <div className="w-6 h-6 bg-sky-500 rounded-lg flex items-center justify-center">
                    <span className="text-white text-[9px] font-bold">V</span>
                  </div>
                  <span className="text-[13px] font-semibold text-[#d0d6e0]">Volo</span>
                  <div className="ml-auto flex items-center gap-2">
                    <div className="flex items-center gap-1">
                      <div className="w-1.5 h-1.5 rounded-full bg-emerald-500"></div>
                      <span className="text-[9px] text-[#62666d]">API</span>
                    </div>
                    <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse shadow-[0_0_6px_rgba(34,197,94,0.5)]"></div>
                  </div>
                </div>
                <div className="flex border-b border-white/[0.06] px-4">
                  <div className="px-3 py-2 text-[11px] font-medium text-[#d0d6e0] border-b-2 border-sky-500">Status</div>
                  <div className="px-3 py-2 text-[11px] font-medium text-[#62666d]">History</div>
                </div>
                <div className="p-4 space-y-3">
                  <div className="rounded-lg p-3 bg-white/[0.03] border border-white/[0.06]">
                    <div className="text-[9px] text-[#62666d] uppercase tracking-wider mb-1">Voice</div>
                    <div className="text-[12px] font-medium text-emerald-400">Listening for "Hey Volo"...</div>
                  </div>
                  <div className="rounded-lg p-3 bg-white/[0.03] border border-white/[0.06]">
                    <div className="text-[9px] text-[#62666d] uppercase tracking-wider mb-2">Microphone</div>
                    <div className="space-y-1.5">
                      <div className="px-3 py-2 rounded border border-sky-500/30 text-[11px] font-medium text-[#d0d6e0] bg-sky-500/[0.06]">Always Listen</div>
                      <div className="px-3 py-2 rounded border border-white/[0.06] text-[11px] text-[#9c9da1]">Click to Talk</div>
                      <div className="px-3 py-2 rounded border border-white/[0.06] text-[11px] text-[#9c9da1]">Off</div>
                    </div>
                  </div>
                  <div className="text-[10px] text-[#62666d]">
                    <p className="text-[#9c9da1] font-medium mb-1">Try saying:</p>
                    <p>"Hey Volo, search React hooks"</p>
                    <p>"Hey Volo, open YouTube lofi"</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* ─── DESKTOP APP SECOND — the companion ─── */}
        <div className="mt-24 md:mt-32">
          <div className="text-center mb-10">
            <h2 className="text-2xl md:text-3xl font-bold tracking-tight text-white/95">Chat with your browsing history.</h2>
            <p className="mt-4 text-white/40 leading-relaxed max-w-xl mx-auto">The desktop app lets you ask questions about what you've browsed. Powered by a local AI — your data never leaves your machine.</p>
          </div>

          <div className="relative">
            <div
              className="absolute inset-0 -inset-x-6 sm:-inset-x-10 -top-6 -bottom-10 rounded-2xl"
              style={{
                background: 'radial-gradient(70% 70% at 50% 55%, rgba(30, 32, 38, 0.5) 0%, rgba(6, 6, 9, 0) 100%)',
              }}
            />

            <div
              className="relative rounded-xl sm:rounded-2xl overflow-hidden"
              style={{
                background: '#101112',
                boxShadow: '0 0 0 0.5px rgba(255, 255, 255, 0.04)',
              }}
            >
              <div
                className="absolute inset-0 pointer-events-none z-10 opacity-[0.35] mix-blend-overlay"
                style={{
                  backgroundImage: 'url("data:image/svg+xml,%3Csvg viewBox=\'0 0 256 256\' xmlns=\'http://www.w3.org/2000/svg\'%3E%3Cfilter id=\'g\'%3E%3CfeTurbulence type=\'fractalNoise\' baseFrequency=\'0.85\' numOctaves=\'4\' stitchTiles=\'stitch\'/%3E%3C/filter%3E%3Crect width=\'100%25\' height=\'100%25\' filter=\'url(%23g)\' opacity=\'0.08\'/%3E%3C/svg%3E")',
                  backgroundSize: '256px 256px',
                }}
              />

              <div className="flex h-[380px] sm:h-[440px] md:h-[540px] relative z-0">
                <InteractiveMockup />
              </div>
            </div>

            <div className="absolute -bottom-5 left-1/2 -translate-x-1/2 px-4 py-2 glass-card rounded-full text-[10px] sm:text-xs font-medium text-white/40 whitespace-nowrap">
              Volo Desktop — Interactive demo
            </div>
          </div>
        </div>

      </div>
    </section>
  );
}

// ─── Interactive Mockup (sidebar + chat as one component with shared state) ───

function InteractiveMockup() {
  const [activeConv, setActiveConv] = useState("What did I listen to yesterday");
  const [messages, setMessages] = useState(CONVERSATIONS["What did I listen to yesterday"]);

  function handleSelectConv(label: string) {
    setActiveConv(label);
    setMessages(CONVERSATIONS[label] || []);
  }

  return (
    <>
      {/* Sidebar */}
      <div className="w-[200px] lg:w-[240px] flex-col hidden md:flex" style={{ background: '#101112' }}>
        <div className="px-4 pt-4 pb-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <div className="w-5 h-5 bg-sky-500 rounded flex items-center justify-center">
                <span className="text-white text-[8px] font-bold">V</span>
              </div>
              <span className="text-[13px] font-semibold text-[#d0d6e0]">Volo</span>
            </div>
            <div className="flex items-center gap-0.5">
              <div className="w-6 h-6 flex items-center justify-center rounded hover:bg-white/[0.04] text-[#62666d]">
                <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M7 2.5a4.5 4.5 0 1 0 0 9 4.5 4.5 0 0 0 0-9zM1 7a6 6 0 1 1 10.7 3.7l3.3 3.3-1 1-3.3-3.3A6 6 0 0 1 1 7z"/></svg>
              </div>
              <div className="w-6 h-6 flex items-center justify-center rounded hover:bg-white/[0.04] text-[#62666d]">
                <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M8 1.5a.75.75 0 0 1 .75.75v5h5a.75.75 0 0 1 0 1.5h-5v5a.75.75 0 0 1-1.5 0v-5h-5a.75.75 0 0 1 0-1.5h5v-5A.75.75 0 0 1 8 1.5z"/></svg>
              </div>
            </div>
          </div>
        </div>

        <div className="px-2 space-y-0.5">
          <div className="flex items-center gap-2.5 px-3 py-1.5 rounded-md text-[13px] bg-white/[0.06] text-[#e2e4e7] cursor-default">
            <span className="text-[#d0d6e0]"><svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M2.5 3A1.5 1.5 0 0 1 4 1.5h8A1.5 1.5 0 0 1 13.5 3v7A1.5 1.5 0 0 1 12 11.5H5.7L3 14V11.5h-.5A1.5 1.5 0 0 1 1 10V3.5 3zm2 0v7h1v1.5l1.7-1.5H12V3H4.5z"/></svg></span>
            <span className="font-medium">Chat</span>
          </div>
          <div className="flex items-center gap-2.5 px-3 py-1.5 rounded-md text-[13px] text-[#9c9da1] hover:bg-white/[0.03] hover:text-[#d0d6e0] cursor-default transition">
            <span className="text-[#62666d]"><svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M8 1.5a6.5 6.5 0 1 0 0 13 6.5 6.5 0 0 0 0-13zM0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8zm8.5-3.5v3.2l2.4 1.4-.7 1.3L7.5 9V4.5h1z"/></svg></span>
            <span className="font-medium">History</span>
          </div>

          <div className="px-2 pt-5 pb-1.5">
            <span className="text-[11px] font-medium text-[#62666d]">Recent</span>
          </div>
          {RECENT_ITEMS.map((item) => (
            <button
              key={item.label}
              onClick={() => handleSelectConv(item.label)}
              className={`w-full text-left flex items-center justify-between px-3 py-1.5 rounded-md text-[12px] cursor-pointer transition ${
                activeConv === item.label
                  ? "bg-white/[0.06] text-[#d0d6e0]"
                  : "text-[#9c9da1] hover:bg-white/[0.03] hover:text-[#d0d6e0]"
              }`}
            >
              <span className="truncate">{item.label}</span>
              <span className="text-[10px] text-[#62666d] ml-2 shrink-0">{item.time}</span>
            </button>
          ))}
        </div>

        <div className="mt-auto p-2 border-t border-white/[0.04]">
          <div className="flex items-center gap-2.5 px-3 py-2 rounded-md text-[12px] text-[#9c9da1] hover:bg-white/[0.03] cursor-default transition">
            <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M8 5.5a2.5 2.5 0 1 0 0 5 2.5 2.5 0 0 0 0-5zM4 8a4 4 0 1 1 8 0 4 4 0 0 1-8 0z"/></svg>
            Settings
          </div>
        </div>
      </div>

      {/* Chat */}
      <div className="flex-1 flex flex-col border-l border-white/[0.06]" style={{ background: '#121314' }}>
        {/* Titlebar */}
        <div className="h-8 sm:h-9 flex items-center justify-between px-4 shrink-0 border-b border-white/[0.04]">
          <span className="text-[11px] text-[#62666d]">Chat</span>
          <div className="hidden sm:flex items-center gap-1">
            <div className="w-8 h-7 flex items-center justify-center text-[#62666d] hover:text-[#9c9da1] rounded transition"><svg width="10" height="1" viewBox="0 0 10 1" fill="currentColor"><rect width="10" height="1"/></svg></div>
            <div className="w-8 h-7 flex items-center justify-center text-[#62666d] hover:text-[#9c9da1] rounded transition"><svg width="9" height="9" viewBox="0 0 9 9" fill="none" stroke="currentColor" strokeWidth="1.2"><rect x="0.5" y="0.5" width="8" height="8" rx="1"/></svg></div>
            <div className="w-8 h-7 flex items-center justify-center text-[#62666d] hover:text-[#9c9da1] rounded transition"><svg width="10" height="10" viewBox="0 0 10 10" fill="none" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round"><line x1="1" y1="1" x2="9" y2="9"/><line x1="9" y1="1" x2="1" y2="9"/></svg></div>
          </div>
        </div>

        {/* Messages */}
        <div className="flex-1 overflow-y-auto">
          <div className="max-w-[520px] mx-auto px-3 sm:px-5 py-3 sm:py-4 space-y-3 sm:space-y-4">
            {messages.map((msg, i) => (
              <div key={`${activeConv}-${i}`}>
                {msg.role === "user" ? (
                  <div className="flex justify-end">
                    <div className="rounded-xl rounded-tr-sm px-3 py-2 bg-sky-500/[0.12] border border-sky-500/20 max-w-[75%]">
                      <p className="text-[11px] sm:text-[12px] text-[#d0d6e0] leading-relaxed">{msg.text}</p>
                    </div>
                  </div>
                ) : (
                  <div className="flex items-start gap-2 sm:gap-3">
                    <div className="w-5 h-5 sm:w-6 sm:h-6 rounded-full bg-sky-500 flex items-center justify-center shrink-0 mt-0.5">
                      <span className="text-white text-[7px] sm:text-[9px] font-bold">V</span>
                    </div>
                    <div className="rounded-xl rounded-tl-sm px-3 py-2 bg-white/[0.04] max-w-[80%]">
                      <p className="text-[11px] sm:text-[12px] text-[#d0d6e0]/80 leading-relaxed whitespace-pre-wrap">{msg.text}</p>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* Input */}
        <div className="px-3 sm:px-5 pb-3 sm:pb-4 border-t border-white/[0.04]">
          <div className="max-w-[520px] mx-auto pt-3">
            <div className="flex items-center rounded-xl px-3 sm:px-4 py-2.5 bg-white/[0.03] border border-white/[0.06]">
              <div className="w-6 h-6 flex items-center justify-center rounded-full text-[#62666d]">
                <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor"><path d="M8 1a3 3 0 0 0-3 3v4a3 3 0 1 0 6 0V4a3 3 0 0 0-3-3zM3.5 7.5a.75.75 0 0 1 .75.75A3.75 3.75 0 0 0 8 12a3.75 3.75 0 0 0 3.75-3.75.75.75 0 0 1 1.5 0A5.25 5.25 0 0 1 8.75 13.4v1.35h1.5a.75.75 0 0 1 0 1.5h-4.5a.75.75 0 0 1 0-1.5h1.5V13.4A5.25 5.25 0 0 1 2.75 8.25a.75.75 0 0 1 .75-.75z"/></svg>
              </div>
              <span className="flex-1 text-[11px] sm:text-[13px] ml-2 text-[#62666d]">Message Volo...</span>
              <div className="w-6 h-6 rounded-full flex items-center justify-center bg-white/[0.06]">
                <svg width="11" height="11" viewBox="0 0 14 14" fill="#62666d"><path d="M7 1v12M3 5l4-4 4 4"/></svg>
              </div>
            </div>
            <p className="text-[10px] text-[#62666d] mt-2 text-center">Volo can make mistakes · Ctrl+Shift+V for voice</p>
          </div>
        </div>
      </div>
    </>
  );
}

function CommandExample({ command, result }: { command: string; result: string }) {
  return (
    <div className="flex items-start gap-3 group">
      <div className="w-5 h-5 mt-0.5 rounded-full bg-sky-400/[0.08] border border-sky-400/20 flex items-center justify-center flex-shrink-0 group-hover:bg-sky-400/[0.12] group-hover:border-sky-400/30 transition">
        <svg width="10" height="10" viewBox="0 0 10 10" fill="none" stroke="#38bdf8" strokeWidth="1.5"><rect x="3" y="0.5" width="4" height="6" rx="2"/><path d="M1.5 4.5a3.5 3.5 0 0 0 7 0"/></svg>
      </div>
      <div>
        <p className="text-sm text-white/75 font-medium group-hover:text-white/90 transition">{command}</p>
        <p className="text-xs text-white/30 mt-0.5">{result}</p>
      </div>
    </div>
  );
}
