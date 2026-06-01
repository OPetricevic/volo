import { useState, useEffect } from "react";

// Demo conversation that plays on loop
const DEMO_SCRIPT = [
  { role: "user" as const, text: "What did I listen to yesterday?" },
  { role: "assistant" as const, text: "Yesterday you played:\n• lofi beats on YouTube — 9:14am\n• jazz playlist on Spotify — 2:30pm\n• ambient music on YouTube — 8:45pm" },
  { role: "user" as const, text: "Play that lofi one again" },
  { role: "assistant" as const, text: "Opening YouTube → lofi beats ✓" },
];

export function ProductShowcase() {
  return (
    <section className="py-20 md:py-32">
      <div className="max-w-6xl mx-auto px-6">
        {/* Desktop App Mockup with live demo */}
        <div className="relative">
          <div className="rounded-xl overflow-hidden shadow-2xl shadow-black/60 border border-white/[0.04]">
            <div className="flex h-[460px] md:h-[520px]">
              {/* Sidebar */}
              <MockSidebar />
              {/* Chat with animation */}
              <MockChat />
            </div>
          </div>

          <div className="absolute -bottom-5 left-1/2 -translate-x-1/2 px-5 py-2 rounded-full text-xs font-medium text-[#949ba4] border border-white/[0.06] backdrop-blur-sm" style={{ background: 'rgba(30,31,34,0.95)' }}>
            Volo Desktop — Live demo
          </div>
        </div>

        {/* Extension section */}
        <div className="mt-32 grid md:grid-cols-2 gap-12 items-center">
          <div>
            <h2 className="text-2xl md:text-3xl font-bold tracking-tight">Works right in your browser.</h2>
            <p className="mt-4 text-volo-muted leading-relaxed">The extension listens for "Hey Volo" on every page. Say a command and it executes instantly.</p>
            <div className="mt-8 space-y-4">
              <CommandExample command="Hey Volo, open YouTube play lofi" result="→ Opens YouTube, plays first result" />
              <CommandExample command="Hey Volo, search React hooks" result="→ Google search" />
              <CommandExample command="Hey Volo, go to GitHub" result="→ Navigates to github.com" />
              <CommandExample command="Hey Volo, close tab" result="→ Closes current tab" />
            </div>
          </div>

          {/* Extension popup */}
          <div className="flex justify-center">
            <div className="w-[280px] rounded-lg overflow-hidden shadow-xl border border-white/[0.06]" style={{ background: '#2b2d31' }}>
              <div className="p-4">
                <div className="flex items-center gap-2 mb-4">
                  <div className="w-6 h-6 bg-[#5865f2] rounded-md flex items-center justify-center">
                    <span className="text-white text-[9px] font-bold">V</span>
                  </div>
                  <span className="text-[13px] font-semibold text-[#dbdee1]">Volo</span>
                  <div className="ml-auto w-2 h-2 rounded-full bg-[#23a559] animate-pulse"></div>
                </div>

                <div className="rounded-md p-3 mb-3" style={{ background: '#313338' }}>
                  <div className="text-[10px] text-[#6d6f78] mb-1">Status</div>
                  <div className="text-[12px] font-medium text-[#23a559]">Listening for "Hey Volo"...</div>
                </div>

                <div className="rounded-md p-3" style={{ background: '#313338' }}>
                  <div className="text-[10px] text-[#6d6f78] mb-2">Microphone</div>
                  <div className="space-y-1.5">
                    <div className="px-2.5 py-1.5 border border-[#5865f2]/50 rounded text-[11px] font-medium text-[#dbdee1]" style={{ background: 'rgba(88,101,242,0.08)' }}>Always Listen</div>
                    <div className="px-2.5 py-1.5 border border-[#3f4147]/50 rounded text-[11px] text-[#949ba4]">Push to Talk</div>
                    <div className="px-2.5 py-1.5 border border-[#3f4147]/50 rounded text-[11px] text-[#949ba4]">Off</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

// ─── Animated Chat Mockup ────────────────────────────────

function MockChat() {
  const [messages, setMessages] = useState<Array<{ role: "user" | "assistant"; text: string }>>([]);
  const [typing, setTyping] = useState("");
  const [showThinking, setShowThinking] = useState(false);
  const [done, setDone] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => playScript(0), 2000);
    return () => clearTimeout(timer);
  }, []);

  function playScript(index: number) {
    if (index >= DEMO_SCRIPT.length) {
      // Demo finished — pre-fill input with a teaser and stop
      setTimeout(() => {
        typeTeaser();
      }, 1500);
      return;
    }

    const step = DEMO_SCRIPT[index];

    if (step.role === "user") {
      let charIndex = 0;
      setTyping("");
      const typeInterval = setInterval(() => {
        charIndex++;
        setTyping(step.text.slice(0, charIndex));
        if (charIndex >= step.text.length) {
          clearInterval(typeInterval);
          setTimeout(() => {
            setTyping("");
            setMessages((prev) => [...prev, { role: "user", text: step.text }]);
            setShowThinking(true);
            setTimeout(() => {
              setShowThinking(false);
              playScript(index + 1);
            }, 1200);
          }, 500);
        }
      }, 45);
    } else {
      setMessages((prev) => [...prev, { role: "assistant", text: step.text }]);
      setTimeout(() => playScript(index + 1), 2500);
    }
  }

  function typeTeaser() {
    const teaser = "What do I usually open in the morning?";
    let charIndex = 0;
    setTyping("");
    const typeInterval = setInterval(() => {
      charIndex++;
      setTyping(teaser.slice(0, charIndex));
      if (charIndex >= teaser.length) {
        clearInterval(typeInterval);
        setDone(true); // Stop — cursor blinks, waiting for "user"
      }
    }, 50);
  }

  return (
    <div className="flex-1 flex flex-col" style={{ background: '#313338' }}>
      {/* Titlebar */}
      <div className="h-9 flex items-center justify-end px-2 shrink-0">
        <div className="flex">
          <div className="w-10 h-9 flex items-center justify-center text-[#949ba4]"><svg width="9" height="1" viewBox="0 0 9 1" fill="currentColor"><rect width="9" height="1"/></svg></div>
          <div className="w-10 h-9 flex items-center justify-center text-[#949ba4]"><svg width="8" height="8" viewBox="0 0 8 8" fill="none" stroke="currentColor" strokeWidth="1"><rect x="0.5" y="0.5" width="7" height="7" rx="0.5"/></svg></div>
          <div className="w-10 h-9 flex items-center justify-center text-[#949ba4]"><svg width="9" height="9" viewBox="0 0 9 9" fill="none" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round"><line x1="1" y1="1" x2="8" y2="8"/><line x1="8" y1="1" x2="1" y2="8"/></svg></div>
        </div>
      </div>

      {/* Messages */}
      <div className="flex-1 overflow-hidden">
        <div className="max-w-[520px] mx-auto px-5 py-4 space-y-5">
          {/* Welcome */}
          <div className="flex items-start gap-3">
            <div className="w-6 h-6 rounded-full bg-[#5865f2] flex items-center justify-center shrink-0 mt-0.5">
              <span className="text-white text-[8px] font-bold">V</span>
            </div>
            <p className="text-[12px] text-[#dbdee1] leading-relaxed">Hey! Ask me about your search history, or give me a voice command.</p>
          </div>

          {/* Animated messages */}
          {messages.map((msg, i) => (
            <div key={i}>
              {msg.role === "user" ? (
                <p className="text-[12px] text-[#dbdee1] leading-relaxed">{msg.text}</p>
              ) : (
                <div className="flex items-start gap-3">
                  <div className="w-6 h-6 rounded-full bg-[#5865f2] flex items-center justify-center shrink-0 mt-0.5">
                    <span className="text-white text-[8px] font-bold">V</span>
                  </div>
                  <p className="text-[12px] text-[#dbdee1] leading-relaxed whitespace-pre-wrap">{msg.text}</p>
                </div>
              )}
            </div>
          ))}

          {/* Thinking indicator */}
          {showThinking && (
            <div className="flex items-center gap-3">
              <div className="w-6 h-6 rounded-full bg-[#5865f2] flex items-center justify-center shrink-0">
                <span className="text-white text-[8px] font-bold">V</span>
              </div>
              <div className="flex gap-1">
                <span className="w-1.5 h-1.5 bg-[#949ba4] rounded-full animate-bounce" style={{ animationDelay: '0ms' }}/>
                <span className="w-1.5 h-1.5 bg-[#949ba4] rounded-full animate-bounce" style={{ animationDelay: '150ms' }}/>
                <span className="w-1.5 h-1.5 bg-[#949ba4] rounded-full animate-bounce" style={{ animationDelay: '300ms' }}/>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Input with typing animation */}
      <div className="px-5 pb-4">
        <div className="max-w-[520px] mx-auto">
          <div className="flex items-center rounded-2xl px-4 py-2.5 border border-[#3f4147]/40" style={{ background: '#383a40' }}>
            <svg width="15" height="15" viewBox="0 0 15 15" fill="none" stroke="#6d6f78" strokeWidth="1.3" strokeLinecap="round" className="shrink-0"><rect x="5" y="1.5" width="5" height="7.5" rx="2.5"/><path d="M2.5 7a5 5 0 0 0 10 0"/><line x1="7.5" y1="12" x2="7.5" y2="13.5"/></svg>
            <span className="flex-1 text-[12px] ml-3 min-h-[18px]">
              {typing ? (
                <span className="text-[#dbdee1]">{typing}<span className="animate-pulse text-[#5865f2]">|</span></span>
              ) : (
                <span className="text-[#6d6f78]">Message Volo...</span>
              )}
            </span>
            <button
              onClick={() => {
                if (done && typing) {
                  // "Send" the teaser message
                  const msg = typing;
                  setTyping("");
                  setMessages((prev) => [...prev, { role: "user", text: msg }]);
                  setShowThinking(true);
                  setDone(false);
                  // Fake response after thinking
                  setTimeout(() => {
                    setShowThinking(false);
                    setMessages((prev) => [...prev, { role: "assistant", text: "In the morning you usually open:\n• Gmail — 72% of the time\n• YouTube — 18%\n• Facebook — 10%\n\nWant me to open Gmail for you?" }]);
                  }, 1500);
                }
              }}
              className={`w-6 h-6 rounded-full flex items-center justify-center transition ${done && typing ? 'bg-[#dbdee1] cursor-pointer hover:scale-110' : typing ? 'bg-[#dbdee1]' : 'bg-[#dbdee1]/20'}`}
            >
              <svg width="11" height="11" viewBox="0 0 11 11" fill={typing ? '#313338' : '#6d6f78'}><path d="M5.5 1v9M3 3.5l2.5-2.5L8 3.5"/></svg>
            </button>
          </div>
          {done && typing && (
            <p className="text-[10px] text-[#5865f2] mt-2 text-center animate-pulse">Click send to try it →</p>
          )}
        </div>
      </div>
    </div>
  );
}

function MockSidebar() {
  return (
    <div className="w-[200px] flex-col hidden md:flex border-r border-white/[0.03]" style={{ background: '#2b2d31' }}>
      <div className="p-3">
        <div className="flex items-center gap-2 px-3 py-2.5 rounded-lg text-[11px] text-[#949ba4]">
          <svg width="13" height="13" viewBox="0 0 13 13" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"><line x1="6.5" y1="2" x2="6.5" y2="11"/><line x1="2" y1="6.5" x2="11" y2="6.5"/></svg>
          New chat
        </div>
      </div>
      <div className="flex-1 px-2 overflow-hidden">
        <div className="px-2 pt-2 pb-1"><span className="text-[9px] font-medium text-[#6d6f78]">Today</span></div>
        <div className="px-3 py-1.5 rounded text-[11px] text-[#dbdee1] bg-[#383a40]/50">What did I listen to...</div>
        <div className="px-3 py-1.5 rounded text-[11px] text-[#949ba4]">React hooks search</div>
        <div className="px-3 py-1.5 rounded text-[11px] text-[#949ba4]">Open YouTube lofi</div>
        <div className="px-2 pt-4 pb-1"><span className="text-[9px] font-medium text-[#6d6f78]">Yesterday</span></div>
        <div className="px-3 py-1.5 rounded text-[11px] text-[#949ba4]">GitHub repos</div>
        <div className="px-3 py-1.5 rounded text-[11px] text-[#949ba4]">Weather check</div>
      </div>
      <div className="p-2 border-t border-white/[0.03]">
        <div className="px-3 py-1.5 text-[11px] text-[#949ba4]">Settings</div>
      </div>
    </div>
  );
}

function CommandExample({ command, result }: { command: string; result: string }) {
  return (
    <div className="flex items-start gap-3 group">
      <div className="w-5 h-5 mt-0.5 rounded-full bg-[#5865f2]/10 flex items-center justify-center flex-shrink-0 group-hover:bg-[#5865f2]/20 transition">
        <svg width="10" height="10" viewBox="0 0 10 10" fill="none" stroke="#5865f2" strokeWidth="1.5"><rect x="3" y="0.5" width="4" height="6" rx="2"/><path d="M1.5 4.5a3.5 3.5 0 0 0 7 0"/></svg>
      </div>
      <div>
        <p className="text-sm text-volo-text font-medium group-hover:text-white transition">{command}</p>
        <p className="text-xs text-volo-muted mt-0.5">{result}</p>
      </div>
    </div>
  );
}
