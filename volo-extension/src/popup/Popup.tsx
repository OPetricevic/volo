import { useEffect, useState } from "react";
import type { MicMode, VoiceState, ExtensionMessage } from "@/shared/types";
import { getHistory, clearHistory } from "@/shared/api";
import { getOnboarded, setOnboarded } from "@/shared/storage";
import { Onboarding } from "./Onboarding";

interface HistoryItem {
  id: string;
  transcript: string;
  parsed_action: string;
  executed_at: string;
}

export function Popup() {
  const [showOnboarding, setShowOnboarding] = useState<boolean | null>(null);
  const [voiceState, setVoiceState] = useState<VoiceState>("idle");
  const [micMode, setMicMode] = useState<MicMode>("always");
  const [apiOnline, setApiOnline] = useState(false);
  const [tab, setTab] = useState<"status" | "history">("status");
  const [history, setHistory] = useState<HistoryItem[]>([]);
  const [historyLoading, setHistoryLoading] = useState(false);

  // Check if onboarding is needed
  useEffect(() => {
    getOnboarded().then((done) => setShowOnboarding(!done));
  }, []);

  useEffect(() => {
    chrome.runtime.sendMessage(
      { type: "GET_STATE" } satisfies ExtensionMessage,
      (response) => {
        if (response) {
          setVoiceState(response.state);
          setMicMode(response.micMode);
          setApiOnline(response.apiOnline ?? false);
        }
      }
    );
  }, []);

  function handleMicModeChange(mode: MicMode) {
    setMicMode(mode);
    chrome.runtime.sendMessage({ type: "SET_MIC_MODE", mode } satisfies ExtensionMessage);
  }

  function handleActivate() {
    chrome.runtime.sendMessage({ type: "ACTIVATE_LISTENING" } satisfies ExtensionMessage);
  }

  async function loadHistory() {
    setHistoryLoading(true);
    const result = await getHistory(1, 10);
    if (result) setHistory(result.commands);
    setHistoryLoading(false);
  }

  useEffect(() => {
    if (tab === "history" && apiOnline) loadHistory();
  }, [tab, apiOnline]);

  // Onboarding check
  function handleOnboardingComplete(mode: MicMode) {
    setMicMode(mode);
    setShowOnboarding(false);
    setOnboarded(true);
    chrome.runtime.sendMessage({ type: "SET_MIC_MODE", mode } satisfies ExtensionMessage);
  }

  if (showOnboarding === null) return <div className="w-[320px] h-[200px] bg-[#1e1f22]" />;
  if (showOnboarding) return <Onboarding onComplete={handleOnboardingComplete} />;

  return (
    <div className="w-[320px] bg-volo-bg text-volo-text">
      {/* Header */}
      <div className="flex items-center gap-2 px-4 pt-4 pb-3">
        <div className="w-7 h-7 bg-volo-accent rounded-lg flex items-center justify-center">
          <span className="text-white text-[10px] font-bold">V</span>
        </div>
        <span className="text-[14px] font-semibold">Volo</span>
        <div className="ml-auto flex items-center gap-2">
          <StatusDot online={apiOnline} label="API" />
          <VoiceDot state={voiceState} />
        </div>
      </div>

      {/* Tabs */}
      <div className="flex border-b border-volo-border/30 px-4">
        <TabButton label="Status" active={tab === "status"} onClick={() => setTab("status")} />
        <TabButton label="History" active={tab === "history"} onClick={() => setTab("history")} />
      </div>

      {/* Content */}
      <div className="p-4">
        {tab === "status" ? (
          <>
            {/* Voice status */}
            <div className="bg-volo-surface rounded-lg p-3 mb-3">
              <div className="text-[11px] text-volo-faint uppercase tracking-wider mb-1">Voice</div>
              <div className={`text-[13px] font-medium ${voiceState === "listening" ? "text-volo-success" : voiceState === "processing" ? "text-volo-warning" : "text-volo-muted"}`}>
                {getStatusText(voiceState)}
              </div>
            </div>

            {/* Mic mode */}
            <div className="bg-volo-surface rounded-lg p-3 mb-3">
              <div className="text-[11px] text-volo-faint uppercase tracking-wider mb-2">Microphone</div>
              <div className="space-y-1">
                <MicOption label="Always Listen" active={micMode === "always"} onClick={() => handleMicModeChange("always")} />
                <MicOption label="Push to Talk" active={micMode === "once"} onClick={() => handleMicModeChange("once")} />
                <MicOption label="Off" active={micMode === "off"} onClick={() => handleMicModeChange("off")} />
              </div>
            </div>

            {/* Activate (push-to-talk) */}
            {micMode === "once" && voiceState === "idle" && (
              <button onClick={handleActivate} className="w-full py-2.5 bg-volo-accent text-white text-[13px] font-medium rounded-lg hover:bg-volo-accent-hover transition">
                🎤 Activate
              </button>
            )}

            {/* Help */}
            <div className="mt-3 text-[11px] text-volo-faint">
              <p className="text-volo-muted font-medium mb-1">Try saying:</p>
              <p>"Hey Volo, search React hooks"</p>
              <p>"Hey Volo, open YouTube lofi"</p>
            </div>
          </>
        ) : (
          <HistoryTab history={history} loading={historyLoading} apiOnline={apiOnline} onClear={async () => { await clearHistory(); setHistory([]); }} />
        )}
      </div>
    </div>
  );
}

// ─── Sub-components ──────────────────────────────────────

function TabButton({ label, active, onClick }: { label: string; active: boolean; onClick: () => void }) {
  return (
    <button
      onClick={onClick}
      className={`px-3 py-2 text-[12px] font-medium border-b-2 transition ${
        active ? "border-volo-accent text-volo-text" : "border-transparent text-volo-muted hover:text-volo-text"
      }`}
    >
      {label}
    </button>
  );
}

function StatusDot({ online, label }: { online: boolean; label: string }) {
  return (
    <div className="flex items-center gap-1" title={`${label}: ${online ? "connected" : "offline"}`}>
      <div className={`w-2 h-2 rounded-full ${online ? "bg-volo-success" : "bg-volo-faint"}`} />
    </div>
  );
}

function VoiceDot({ state }: { state: VoiceState }) {
  const cls = state === "listening" ? "bg-volo-success animate-pulse" : state === "processing" ? "bg-volo-warning" : state === "wake-word-detected" ? "bg-volo-accent animate-pulse" : "bg-volo-faint";
  return <div className={`w-2 h-2 rounded-full ${cls}`} title={`Voice: ${state}`} />;
}

function MicOption({ label, active, onClick }: { label: string; active: boolean; onClick: () => void }) {
  return (
    <button
      onClick={onClick}
      className={`w-full text-left px-3 py-2 rounded text-[12px] font-medium transition ${
        active ? "bg-volo-accent/10 text-volo-text border border-volo-accent/40" : "text-volo-muted border border-transparent hover:bg-volo-input/30 hover:text-volo-text"
      }`}
    >
      {label}
    </button>
  );
}

function HistoryTab({ history, loading, apiOnline, onClear }: { history: HistoryItem[]; loading: boolean; apiOnline: boolean; onClear: () => void }) {
  if (!apiOnline) return <p className="text-[12px] text-volo-faint text-center py-6">API offline — history unavailable</p>;
  if (loading) return <p className="text-[12px] text-volo-faint text-center py-6">Loading...</p>;
  if (history.length === 0) return <p className="text-[12px] text-volo-faint text-center py-6">No commands yet</p>;

  return (
    <>
      <div className="space-y-1.5 max-h-[200px] overflow-y-auto">
        {history.map((item) => (
          <div key={item.id} className="bg-volo-surface rounded-md px-3 py-2">
            <div className="text-[12px] text-volo-text truncate">{item.transcript}</div>
            <div className="flex items-center justify-between mt-0.5">
              <span className="text-[10px] text-volo-accent">{item.parsed_action}</span>
              <span className="text-[10px] text-volo-faint">{new Date(item.executed_at).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}</span>
            </div>
          </div>
        ))}
      </div>
      <button onClick={onClear} className="w-full mt-3 py-2 text-[11px] text-volo-faint border border-volo-border/30 rounded-md hover:text-volo-danger hover:border-volo-danger/50 transition">
        Clear History
      </button>
    </>
  );
}

function getStatusText(state: VoiceState): string {
  switch (state) {
    case "listening": return "Listening for \"Hey Volo\"...";
    case "wake-word-detected": return "Heard you! Listening...";
    case "processing": return "Processing...";
    default: return "Idle";
  }
}
