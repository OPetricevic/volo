import { useState, useEffect } from "react";
import { config } from "../config";

interface MacroAction {
  type: string;
  url: string;
}

interface Macro {
  id: string;
  trigger_phrase: string;
  name: string;
  actions: MacroAction[];
  enabled: boolean;
}

const DEMO_MACRO: Macro = {
  id: "default-demo",
  trigger_phrase: "morning routine",
  name: "Morning Routine",
  actions: [
    { type: "navigate", url: "https://mail.google.com" },
    { type: "navigate", url: "https://youtube.com" },
    { type: "navigate", url: "https://weather.com" },
  ],
  enabled: true,
};

export function MacrosView() {
  const [macros, setMacros] = useState<Macro[]>([DEMO_MACRO]);
  const [showCreate, setShowCreate] = useState(false);
  const [loading, setLoading] = useState(true);
  const [offline, setOffline] = useState(false);
  const [error, setError] = useState("");

  // Fetch macros from API on mount
  useEffect(() => {
    fetchMacros();
  }, []);

  async function fetchMacros() {
    setLoading(true);
    setError("");
    try {
      const res = await fetch(`${config.apiUrl}/api/v1/macros`, {
        headers: { "Content-Type": "application/json" },
      });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      const fetched = data.data?.macros ?? [];
      setMacros(fetched.length > 0 ? fetched : [DEMO_MACRO]);
      setOffline(false);
    } catch {
      // API offline — use local demo
      setOffline(true);
      setMacros([DEMO_MACRO]);
    }
    setLoading(false);
  }

  async function handleCreate(macro: Omit<Macro, "id" | "enabled">) {
    setError("");
    try {
      const res = await fetch(`${config.apiUrl}/api/v1/macros`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(macro),
      });
      if (!res.ok) {
        const data = await res.json();
        throw new Error(data.error?.message || "Failed to create macro");
      }
      setShowCreate(false);
      await fetchMacros();
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to create macro");
    }
  }

  async function handleDelete(id: string) {
    try {
      await fetch(`${config.apiUrl}/api/v1/macros/${id}`, { method: "DELETE" });
      await fetchMacros();
    } catch {
      setError("Failed to delete macro. Check your connection.");
    }
  }

  async function handleToggle(id: string) {
    const macro = macros.find((m) => m.id === id);
    if (!macro) return;
    try {
      await fetch(`${config.apiUrl}/api/v1/macros/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ...macro, enabled: !macro.enabled }),
      });
      await fetchMacros();
    } catch {
      setError("Failed to update macro. Check your connection.");
    }
  }

  return (
    <div className="flex-1 overflow-y-auto">
      {/* Header bar */}
      <div className="sticky top-0 z-10 bg-[var(--chat)] border-b border-[var(--border)] px-6 py-4">
        <div className="flex items-center justify-between max-w-3xl">
          <div>
            <h2 className="text-[15px] font-semibold text-[var(--text)]">Macros</h2>
            <p className="text-[12px] text-[var(--text-muted)] mt-0.5">
              Custom voice commands that open multiple sites at once
            </p>
          </div>
          {!showCreate && (
            <button
              onClick={() => setShowCreate(true)}
              disabled={offline}
              className="px-4 py-2 text-[12px] font-medium bg-[var(--accent)] text-white rounded-lg hover:bg-[var(--accent-hover)] disabled:opacity-40 disabled:cursor-not-allowed transition flex items-center gap-1.5"
            >
              <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor"><path d="M8 1.5a.75.75 0 0 1 .75.75v5h5a.75.75 0 0 1 0 1.5h-5v5a.75.75 0 0 1-1.5 0v-5h-5a.75.75 0 0 1 0-1.5h5v-5A.75.75 0 0 1 8 1.5z"/></svg>
              New Macro
            </button>
          )}
        </div>
      </div>

      <div className="px-6 py-5 max-w-3xl">
        {/* Offline / connection notice */}
        {offline && (
          <div className="mb-4 px-4 py-3 rounded-lg bg-[var(--warning)]/10 border border-[var(--warning)]/20 flex items-center gap-2">
            <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor" className="text-[var(--warning)] shrink-0"><path d="M8 1a7 7 0 1 0 0 14A7 7 0 0 0 8 1zm0 10.5a.75.75 0 1 1 0-1.5.75.75 0 0 1 0 1.5zM8.75 4.5v4a.75.75 0 0 1-1.5 0v-4a.75.75 0 0 1 1.5 0z"/></svg>
            <p className="text-[12px] text-[var(--warning)]">
              Offline — showing demo macro. Connect to the internet and sign in to create and sync macros across devices.
            </p>
          </div>
        )}

        {/* Error */}
        {error && (
          <div className="mb-4 px-4 py-3 rounded-lg bg-[var(--danger)]/10 border border-[var(--danger)]/20">
            <p className="text-[12px] text-[var(--danger)]">{error}</p>
          </div>
        )}

        {/* Requires internet notice */}
        {!offline && (
          <div className="mb-4 px-3 py-2 rounded-lg border border-[var(--border)] flex items-center gap-2">
            <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor" className="text-[var(--text-faint)] shrink-0"><path d="M8 0a8 8 0 1 0 0 16A8 8 0 0 0 8 0zm.75 4.5v4a.75.75 0 0 1-1.5 0v-4a.75.75 0 0 1 1.5 0zM8 11.5a.75.75 0 1 1 0 1.5.75.75 0 0 1 0-1.5z"/></svg>
            <p className="text-[11px] text-[var(--text-faint)]">
              Requires internet · Macros sync to your account and are used by the browser extension
            </p>
          </div>
        )}

        {/* How it works — shown when 0-1 macros */}
        {macros.length <= 1 && !showCreate && (
          <div className="mb-6">
            <div className="rounded-xl border border-[var(--border)] bg-[var(--input)] p-5">
              <h3 className="text-[13px] font-medium text-[var(--text)] mb-3 flex items-center gap-2">
                <svg width="14" height="14" viewBox="0 0 16 16" fill="currentColor" className="text-[var(--accent)]"><path d="M8 1.5a6.5 6.5 0 1 0 0 13 6.5 6.5 0 0 0 0-13zM0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8zm7.25-2.75a.75.75 0 1 1 1.5 0 .75.75 0 0 1-1.5 0zM6.75 7.5a.75.75 0 0 1 .75-.75h.5a.75.75 0 0 1 .75.75v3a.5.5 0 0 0 .5.5h.25a.75.75 0 0 1 0 1.5h-.25A2 2 0 0 1 7.25 10.5V8.25a.75.75 0 0 1-.5-.75z"/></svg>
                How macros work
              </h3>
              <div className="space-y-3 text-[12px] text-[var(--text-muted)]">
                <div className="flex items-start gap-3">
                  <span className="text-[var(--text-faint)] font-mono text-[11px] mt-0.5">1.</span>
                  <p>Create a macro with a <span className="text-[var(--text)]">trigger phrase</span> and a list of websites to open.</p>
                </div>
                <div className="flex items-start gap-3">
                  <span className="text-[var(--text-faint)] font-mono text-[11px] mt-0.5">2.</span>
                  <p>Say <span className="text-[var(--accent)]">"Hey Volo, [trigger phrase]"</span> and all sites open instantly in your browser.</p>
                </div>
                <div className="flex items-start gap-3">
                  <span className="text-[var(--text-faint)] font-mono text-[11px] mt-0.5">3.</span>
                  <p>Macros sync to your account — the <span className="text-[var(--text)]">browser extension</span> picks them up automatically.</p>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Create form */}
        {showCreate && (
          <div className="mb-5">
            <MacroForm
              onSave={handleCreate}
              onCancel={() => setShowCreate(false)}
            />
          </div>
        )}

        {/* Macro list */}
        {!loading && macros.length > 0 && (
          <div className="space-y-3">
            {macros.map((macro) => (
              <MacroCard
                key={macro.id}
                macro={macro}
                onDelete={() => handleDelete(macro.id)}
                onToggle={() => handleToggle(macro.id)}
                disabled={offline}
              />
            ))}
          </div>
        )}

        {/* Empty state */}
        {!loading && macros.length === 0 && !showCreate && (
          <div className="text-center py-8">
            <p className="text-[13px] text-[var(--text-faint)]">No macros yet. Click "New Macro" to create one.</p>
          </div>
        )}
      </div>
    </div>
  );
}

function MacroCard({ macro, onDelete, onToggle, disabled }: { macro: Macro; onDelete: () => void; onToggle: () => void; disabled?: boolean }) {
  return (
    <div className={`p-4 rounded-xl border transition ${macro.enabled ? "border-[var(--border)] bg-[var(--input)]" : "border-[var(--border)]/50 bg-[var(--input)]/50 opacity-50"}`}>
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 rounded-lg bg-[var(--accent)]/10 border border-[var(--accent)]/20 flex items-center justify-center">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-[var(--accent)]"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" strokeLinecap="round" strokeLinejoin="round"/></svg>
          </div>
          <div>
            <h4 className="text-[13px] font-medium text-[var(--text)]">{macro.name}</h4>
            <p className="text-[11px] text-[var(--accent)]">"Hey Volo, {macro.trigger_phrase}"</p>
          </div>
        </div>
        <div className="flex items-center gap-3">
          <button
            onClick={onToggle}
            disabled={disabled}
            className={`w-9 h-5 rounded-full transition relative disabled:opacity-40 ${macro.enabled ? "bg-[var(--accent)]" : "bg-[var(--border)]"}`}
          >
            <div className={`w-3.5 h-3.5 bg-white rounded-full absolute top-[3px] transition-all ${macro.enabled ? "left-[18px]" : "left-[3px]"}`} />
          </button>
          <button onClick={onDelete} disabled={disabled} className="text-[var(--text-faint)] hover:text-[var(--danger)] disabled:opacity-40 transition p-1.5 rounded hover:bg-[var(--danger)]/10">
            <svg width="13" height="13" viewBox="0 0 16 16" fill="currentColor"><path d="M5.5 1A.5.5 0 0 1 6 .5h4a.5.5 0 0 1 .5.5v1h3a.5.5 0 0 1 0 1h-.5l-.8 10.4a1.5 1.5 0 0 1-1.5 1.4H5.3a1.5 1.5 0 0 1-1.5-1.4L3 3.5h-.5a.5.5 0 0 1 0-1h3V1zm1 .5v.5h3V1.5h-3z"/></svg>
          </button>
        </div>
      </div>
      <div className="mt-3 ml-11 flex flex-wrap gap-2">
        {macro.actions.map((action, i) => (
          <span key={i} className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-[var(--bg)] border border-[var(--border)] text-[11px] text-[var(--text-muted)]">
            <svg width="10" height="10" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" className="text-[var(--text-faint)]"><path d="M2 8h12M10 4l4 4-4 4"/></svg>
            {action.url.replace(/^https?:\/\//, '')}
          </span>
        ))}
      </div>
    </div>
  );
}

function MacroForm({ onSave, onCancel }: { onSave: (macro: Omit<Macro, "id" | "enabled">) => void; onCancel: () => void }) {
  const [name, setName] = useState("");
  const [trigger, setTrigger] = useState("");
  const [urls, setUrls] = useState([""]);

  function addUrl() {
    if (urls.length < 3) setUrls([...urls, ""]);
  }

  function removeUrl(index: number) {
    if (urls.length > 1) setUrls(urls.filter((_, i) => i !== index));
  }

  function updateUrl(index: number, value: string) {
    setUrls(urls.map((u, i) => (i === index ? value : u)));
  }

  function handleSubmit() {
    const validUrls = urls.filter((u) => u.trim() !== "");
    if (!name.trim() || !trigger.trim() || validUrls.length === 0) return;

    onSave({
      name: name.trim(),
      trigger_phrase: trigger.trim().toLowerCase(),
      actions: validUrls.map((url) => ({
        type: "navigate",
        url: url.startsWith("http") ? url : `https://${url}`,
      })),
    });
  }

  const isValid = name.trim() && trigger.trim() && urls.some((u) => u.trim());

  return (
    <div className="rounded-xl border border-[var(--accent)]/30 bg-[var(--input)] p-5">
      <div className="flex items-center justify-between mb-5">
        <h3 className="text-[14px] font-medium text-[var(--text)]">New Macro</h3>
        <button onClick={onCancel} className="text-[var(--text-faint)] hover:text-[var(--text)] transition text-[12px]">
          Cancel
        </button>
      </div>

      <div className="space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="text-[11px] font-medium text-[var(--text-muted)] uppercase tracking-wider mb-1.5 block">Name</label>
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Morning routine"
              className="w-full px-3 py-2.5 bg-[var(--bg)] border border-[var(--border)] rounded-lg text-[13px] text-[var(--text)] placeholder:text-[var(--text-faint)] outline-none focus:border-[var(--accent)]/50 transition"
            />
          </div>
          <div>
            <label className="text-[11px] font-medium text-[var(--text-muted)] uppercase tracking-wider mb-1.5 block">Trigger phrase</label>
            <input
              value={trigger}
              onChange={(e) => setTrigger(e.target.value)}
              placeholder="morning routine"
              className="w-full px-3 py-2.5 bg-[var(--bg)] border border-[var(--border)] rounded-lg text-[13px] text-[var(--text)] placeholder:text-[var(--text-faint)] outline-none focus:border-[var(--accent)]/50 transition"
            />
          </div>
        </div>

        {trigger && (
          <p className="text-[11px] text-[var(--accent)] -mt-2">
            Say: "Hey Volo, {trigger.toLowerCase()}"
          </p>
        )}

        <div>
          <label className="text-[11px] font-medium text-[var(--text-muted)] uppercase tracking-wider mb-1.5 block">
            Sites to open ({urls.filter(u => u.trim()).length}/3)
          </label>
          <div className="space-y-2">
            {urls.map((url, i) => (
              <div key={i} className="flex items-center gap-2">
                <span className="text-[11px] text-[var(--text-faint)] w-4 text-center">{i + 1}</span>
                <input
                  value={url}
                  onChange={(e) => updateUrl(i, e.target.value)}
                  placeholder="gmail.com"
                  className="flex-1 px-3 py-2.5 bg-[var(--bg)] border border-[var(--border)] rounded-lg text-[13px] text-[var(--text)] placeholder:text-[var(--text-faint)] outline-none focus:border-[var(--accent)]/50 transition"
                />
                {urls.length > 1 && (
                  <button onClick={() => removeUrl(i)} className="text-[var(--text-faint)] hover:text-[var(--danger)] transition p-1.5 rounded hover:bg-[var(--danger)]/10">
                    <svg width="12" height="12" viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round"><line x1="4" y1="8" x2="12" y2="8"/></svg>
                  </button>
                )}
              </div>
            ))}
          </div>
          {urls.length < 3 && (
            <button onClick={addUrl} className="mt-2 text-[12px] text-[var(--accent)] hover:text-[var(--accent-hover)] transition flex items-center gap-1">
              <svg width="10" height="10" viewBox="0 0 16 16" fill="currentColor"><path d="M8 1.5a.75.75 0 0 1 .75.75v5h5a.75.75 0 0 1 0 1.5h-5v5a.75.75 0 0 1-1.5 0v-5h-5a.75.75 0 0 1 0-1.5h5v-5A.75.75 0 0 1 8 1.5z"/></svg>
              Add site
            </button>
          )}
        </div>
      </div>

      <div className="mt-6 flex justify-end">
        <button
          onClick={handleSubmit}
          disabled={!isValid}
          className="px-5 py-2.5 text-[13px] font-medium bg-[var(--accent)] text-white rounded-lg hover:bg-[var(--accent-hover)] disabled:opacity-30 disabled:cursor-not-allowed transition"
        >
          Create Macro
        </button>
      </div>
    </div>
  );
}
