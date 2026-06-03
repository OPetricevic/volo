import { useState, useRef, useEffect } from "react";
import { config } from "../config";

interface Message {
  id: string;
  role: "user" | "assistant";
  content: string;
}

export function ChatView() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    scrollRef.current?.scrollTo({ top: scrollRef.current.scrollHeight, behavior: "smooth" });
  }, [messages]);

  useEffect(() => {
    if (!window.volo?.onActivateVoice) return;
    const cleanup = window.volo.onActivateVoice(() => {
      console.log("[Volo] Voice activated via hotkey");
    });
    return cleanup;
  }, []);

  async function handleSend() {
    const text = input.trim();
    if (!text || loading) return;

    const userMsg: Message = { id: crypto.randomUUID(), role: "user", content: text };
    setMessages((prev) => [...prev, userMsg]);
    setInput("");
    setLoading(true);

    try {
      let historyContext = "(No history available)";
      try {
        const histRes = await fetch(`${config.apiUrl}/api/v1/history/context`, { headers: { "Content-Type": "application/json" } });
        if (histRes.ok) {
          const histData = await histRes.json();
          historyContext = histData.data?.context || historyContext;
        }
      } catch { /* API offline */ }

      const prompt = buildPrompt(text, historyContext);
      const ollamaRes = await fetch(`${config.ollamaUrl}/api/generate`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ model: "gemma2:2b", prompt, stream: false, options: { temperature: 0.3, num_predict: 256 } }),
      });

      if (!ollamaRes.ok) throw new Error("Ollama not available");
      const ollamaData = await ollamaRes.json();
      const reply = ollamaData.response || "I couldn't generate a response.";

      setMessages((prev) => [...prev, { id: crypto.randomUUID(), role: "assistant", content: reply.trim() }]);
    } catch {
      setMessages((prev) => [...prev, { id: crypto.randomUUID(), role: "assistant", content: "Volo AI is not available. Make sure Ollama is running locally.\n\nVoice commands still work — say \"Hey Volo\" followed by a command." }]);
    }

    setLoading(false);
  }

  function handleKeyDown(e: React.KeyboardEvent) {
    if (e.key === "Enter" && !e.shiftKey) { e.preventDefault(); handleSend(); }
  }

  return (
    <div className="flex flex-col h-full">
      {/* Messages — centered, max-width like ChatGPT */}
      <div ref={scrollRef} className="flex-1 overflow-y-auto">
        <div className="max-w-[720px] mx-auto px-6 py-6">
          {messages.length === 0 && <EmptyState />}
          {messages.map((msg) => (
            <MessageBlock key={msg.id} message={msg} />
          ))}
          {loading && (
            <div className="py-6">
              <div className="flex items-center gap-2">
                <VoloIcon />
                <div className="flex gap-1 ml-1">
                  <span className="w-1.5 h-1.5 bg-[var(--text-muted)] rounded-full animate-bounce" style={{ animationDelay: "0ms" }} />
                  <span className="w-1.5 h-1.5 bg-[var(--text-muted)] rounded-full animate-bounce" style={{ animationDelay: "150ms" }} />
                  <span className="w-1.5 h-1.5 bg-[var(--text-muted)] rounded-full animate-bounce" style={{ animationDelay: "300ms" }} />
                </div>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Input — centered, prominent */}
      <div className="border-t border-[var(--border)]/20 bg-[var(--chat)]">
        <div className="max-w-[720px] mx-auto px-6 py-4">
          <div className="flex items-end bg-[var(--input)] rounded-2xl px-4 py-3 border border-[var(--border)]/30 focus-within:border-[var(--text-faint)]/50 transition">
            <button className="w-8 h-8 flex items-center justify-center rounded-full hover:bg-[var(--border)]/30 text-[var(--text-faint)] hover:text-[var(--text-muted)] transition shrink-0 mb-0.5" aria-label="Voice">
              <svg width="18" height="18" viewBox="0 0 18 18" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round">
                <rect x="6" y="2" width="6" height="9" rx="3"/><path d="M3.5 8.5a5.5 5.5 0 0 0 11 0"/><line x1="9" y1="14" x2="9" y2="16"/>
              </svg>
            </button>
            <textarea
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Message Volo..."
              rows={1}
              className="flex-1 bg-transparent text-[14px] text-[var(--text)] placeholder:text-[var(--text-faint)] resize-none outline-none max-h-[200px] px-2 py-1 leading-relaxed"
            />
            <button
              onClick={handleSend}
              disabled={!input.trim() || loading}
              className="w-8 h-8 flex items-center justify-center rounded-full bg-[var(--text)] text-[var(--bg)] disabled:opacity-20 hover:opacity-80 transition shrink-0 mb-0.5"
              aria-label="Send"
            >
              <svg width="14" height="14" viewBox="0 0 14 14" fill="currentColor"><path d="M7 1v12M3 5l4-4 4 4"/></svg>
            </button>
          </div>
          <p className="text-[11px] text-[var(--text-faint)] mt-2 text-center">
            Volo can make mistakes. Voice: Ctrl+Shift+V
          </p>
        </div>
      </div>
    </div>
  );
}

function EmptyState() {
  return (
    <div className="flex flex-col items-center justify-center h-full min-h-[400px] text-center">
      <div className="w-12 h-12 rounded-full bg-[var(--accent)] flex items-center justify-center mb-4">
        <span className="text-white text-lg font-bold">V</span>
      </div>
      <h2 className="text-xl font-semibold text-[var(--text)] mb-2">How can I help?</h2>
      <p className="text-[14px] text-[var(--text-muted)] max-w-sm">
        Ask about your search history, or give me a voice command. Try "What did I search today?" or say "Hey Volo, open YouTube".
      </p>
    </div>
  );
}

function MessageBlock({ message }: { message: Message }) {
  const isAssistant = message.role === "assistant";

  return (
    <div className={`py-5 ${isAssistant ? "" : ""}`}>
      {isAssistant && (
        <div className="flex items-start gap-3">
          <VoloIcon />
          <div className="text-[14px] text-[var(--text)] leading-[1.7] whitespace-pre-wrap pt-0.5">
            {message.content}
          </div>
        </div>
      )}
      {!isAssistant && (
        <div className="text-[14px] text-[var(--text)] leading-[1.7] whitespace-pre-wrap">
          {message.content}
        </div>
      )}
    </div>
  );
}

function VoloIcon() {
  return (
    <div className="w-7 h-7 rounded-full bg-[var(--accent)] flex items-center justify-center shrink-0">
      <span className="text-white text-[10px] font-bold">V</span>
    </div>
  );
}

function buildPrompt(message: string, historyContext: string): string {
  return `You are Volo, a voice assistant's history helper. The user will ask questions about their browsing and search history.

Given the user's question and their recent command history below, answer naturally and concisely. Keep responses short and helpful.

If the user asks about something not in their history, say so politely.
If the user gives a command (like "open youtube"), tell them to use voice mode instead.

Recent command history:
${historyContext}

User: ${message}

Volo:`;
}
