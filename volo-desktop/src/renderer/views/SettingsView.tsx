import { useState, useEffect } from "react";
import { config } from "../config";

type MicMode = "always" | "once" | "off";

export function SettingsView() {
  const [micMode, setMicMode] = useState<MicMode>("always");
  const [ollamaStatus, setOllamaStatus] = useState<"checking" | "running" | "stopped" | "not-installed">("checking");
  const [apiStatus, setApiStatus] = useState<"checking" | "connected" | "offline">("checking");

  useEffect(() => {
    checkStatuses();
  }, []);

  async function checkStatuses() {
    // Check API
    try {
      const res = await fetch(`${config.apiUrl}/health`);
      setApiStatus(res.ok ? "connected" : "offline");
    } catch {
      setApiStatus("offline");
    }

    // Check Ollama
    try {
      const running = await window.volo.getOllamaStatus();
      setOllamaStatus(running ? "running" : "stopped");
    } catch {
      setOllamaStatus("not-installed");
    }
  }

  return (
    <div className="flex-1 overflow-y-auto p-6 max-w-2xl">
      <h2 className="text-lg font-semibold text-[var(--text)] mb-6">Settings</h2>

      {/* Microphone */}
      <SettingsSection title="Microphone">
        <div className="space-y-2">
          <RadioOption
            label="Always Listen"
            description="Continuously listens for 'Hey Volo'"
            checked={micMode === "always"}
            onChange={() => setMicMode("always")}
          />
          <RadioOption
            label="Click to Talk"
            description="Activate with Ctrl+Shift+V or click the mic button"
            checked={micMode === "once"}
            onChange={() => setMicMode("once")}
          />
          <RadioOption
            label="Off"
            description="Voice disabled — text input only"
            checked={micMode === "off"}
            onChange={() => setMicMode("off")}
          />
        </div>
      </SettingsSection>

      {/* Connection Status */}
      <SettingsSection title="Connections">
        <div className="space-y-3">
          <StatusRow
            label="Volo API"
            status={apiStatus === "connected" ? "online" : apiStatus === "checking" ? "checking" : "offline"}
          />
          <StatusRow
            label="Volo AI (Ollama)"
            status={ollamaStatus === "running" ? "online" : ollamaStatus === "checking" ? "checking" : "offline"}
          />
          {ollamaStatus === "not-installed" && (
            <p className="text-xs text-[var(--text-muted)] ml-6">
              Volo AI is not installed. Reinstall Volo with the AI option enabled, or install Ollama manually.
            </p>
          )}
        </div>
      </SettingsSection>

      {/* Account */}
      <SettingsSection title="Account">
        <div className="bg-[var(--surface)] border border-[var(--border)] rounded-lg p-4">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-[var(--surface-2)] flex items-center justify-center">
              <span className="text-sm text-[var(--text-muted)]">U</span>
            </div>
            <div>
              <p className="text-sm font-medium text-[var(--text)]">Anonymous User</p>
              <p className="text-xs text-[var(--text-muted)]">Sign in to sync across devices</p>
            </div>
          </div>
          <div className="mt-4 flex gap-2">
            <button className="px-3 py-1.5 text-xs font-medium bg-[var(--accent)] text-white rounded-md hover:bg-[var(--accent-hover)] transition">
              Sign in with Google
            </button>
            <button className="px-3 py-1.5 text-xs font-medium border border-[var(--border)] text-[var(--text-muted)] rounded-md hover:text-[var(--text)] hover:border-[var(--text-muted)] transition">
              Create Account
            </button>
          </div>
        </div>
      </SettingsSection>

      {/* About */}
      <SettingsSection title="About">
        <div className="text-sm text-[var(--text-muted)] space-y-1">
          <p>Volo Desktop v0.1.0</p>
          <p>Voice-first browser assistant</p>
        </div>
      </SettingsSection>
    </div>
  );
}

function SettingsSection({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="mb-8">
      <h3 className="text-xs font-semibold uppercase tracking-wider text-[var(--text-muted)] mb-3">{title}</h3>
      {children}
    </div>
  );
}

function RadioOption({
  label,
  description,
  checked,
  onChange,
}: {
  label: string;
  description: string;
  checked: boolean;
  onChange: () => void;
}) {
  return (
    <label className={`flex items-start gap-3 p-3 rounded-lg border cursor-pointer transition ${
      checked ? "border-[var(--accent)] bg-[var(--accent)]/5" : "border-[var(--border)] hover:border-[var(--text-muted)]"
    }`}>
      <input
        type="radio"
        checked={checked}
        onChange={onChange}
        className="mt-0.5 accent-[var(--accent)]"
      />
      <div>
        <div className="text-sm font-medium text-[var(--text)]">{label}</div>
        <div className="text-xs text-[var(--text-muted)]">{description}</div>
      </div>
    </label>
  );
}

function StatusRow({ label, status }: { label: string; status: "online" | "offline" | "checking" }) {
  const colors = {
    online: "bg-[var(--success)]",
    offline: "bg-[var(--danger)]",
    checking: "bg-[var(--warning)] animate-pulse",
  };
  const labels = {
    online: "Connected",
    offline: "Offline",
    checking: "Checking...",
  };

  return (
    <div className="flex items-center gap-3">
      <div className={`w-2 h-2 rounded-full ${colors[status]}`} />
      <span className="text-sm text-[var(--text)]">{label}</span>
      <span className="text-xs text-[var(--text-muted)] ml-auto">{labels[status]}</span>
    </div>
  );
}
