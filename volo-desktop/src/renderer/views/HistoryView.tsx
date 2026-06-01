import { useState, useEffect } from "react";
import { config } from "../config";

interface HistoryItem {
  id: string;
  transcript: string;
  parsed_action: string;
  parsed_target?: string;
  parsed_query?: string;
  confidence: number;
  executed_at: string;
}

export function HistoryView() {
  const [history, setHistory] = useState<HistoryItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState("");

  useEffect(() => {
    loadHistory();
  }, []);

  async function loadHistory() {
    setLoading(true);
    try {
      const response = await fetch(`${config.apiUrl}/api/v1/history?page=1&page_size=50`);
      const data = await response.json();
      setHistory(data.data?.commands || []);
    } catch {
      setHistory([]);
    }
    setLoading(false);
  }

  const filtered = history.filter((item) =>
    item.transcript.toLowerCase().includes(filter.toLowerCase())
  );

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="p-4 border-b border-[var(--border)]">
        <h2 className="text-lg font-semibold text-[var(--text)]">Command History</h2>
        <div className="mt-3">
          <input
            type="text"
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
            placeholder="Search history..."
            className="w-full bg-[var(--surface)] border border-[var(--border)] rounded-lg px-3 py-2 text-sm text-[var(--text)] placeholder:text-[var(--text-muted)] outline-none focus:border-[var(--accent)] transition"
          />
        </div>
      </div>

      {/* List */}
      <div className="flex-1 overflow-y-auto p-4 space-y-2">
        {loading ? (
          <p className="text-sm text-[var(--text-muted)] text-center py-8">Loading...</p>
        ) : filtered.length === 0 ? (
          <p className="text-sm text-[var(--text-muted)] text-center py-8">
            {filter ? "No matching commands." : "No commands yet. Say \"Hey Volo\" to get started."}
          </p>
        ) : (
          filtered.map((item) => (
            <HistoryCard key={item.id} item={item} />
          ))
        )}
      </div>
    </div>
  );
}

function HistoryCard({ item }: { item: HistoryItem }) {
  const time = new Date(item.executed_at).toLocaleString([], {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });

  const actionColors: Record<string, string> = {
    search: "text-[var(--accent)]",
    navigate: "text-[var(--teal)]",
    "open-and-search": "text-[var(--accent)]",
    "open-and-play": "text-[var(--success)]",
    "browser-control": "text-[var(--text-muted)]",
  };

  return (
    <div className="bg-[var(--surface)] border border-[var(--border)] rounded-lg p-3 hover:border-[var(--accent)]/50 transition cursor-pointer">
      <div className="flex items-start justify-between gap-3">
        <div className="flex-1 min-w-0">
          <p className="text-sm text-[var(--text)] truncate">{item.transcript}</p>
          <div className="flex items-center gap-2 mt-1">
            <span className={`text-xs font-medium ${actionColors[item.parsed_action] || ""}`}>
              {item.parsed_action}
            </span>
            {item.parsed_target && (
              <span className="text-xs text-[var(--text-muted)]">→ {item.parsed_target}</span>
            )}
          </div>
        </div>
        <span className="text-[10px] text-[var(--text-muted)] whitespace-nowrap">{time}</span>
      </div>
    </div>
  );
}
