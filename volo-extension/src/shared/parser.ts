import { CommandResult, SITE_MAP, SEARCH_MAP } from "./types";

const WAKE_WORD = "hey volo";

/**
 * Check if transcript contains the wake word.
 * Returns the command portion after the wake word, or null.
 */
export function extractAfterWakeWord(transcript: string): string | null {
  const lower = transcript.toLowerCase().trim();
  const idx = lower.indexOf(WAKE_WORD);
  if (idx === -1) return null;
  const after = lower.slice(idx + WAKE_WORD.length).trim();
  return after.length > 0 ? after : "";
}

/**
 * Check if a transcript starts with or contains the wake word.
 */
export function containsWakeWord(transcript: string): boolean {
  return transcript.toLowerCase().trim().includes(WAKE_WORD);
}

/**
 * Rule-based intent parser.
 * Handles the fast path for obvious commands.
 */
export function parseCommand(raw: string): CommandResult {
  const text = raw.toLowerCase().trim();

  // Browser control commands
  if (text === "go back" || text === "back") {
    return { action: "browser-control", query: "back", confidence: 0.95 };
  }
  if (text === "go forward" || text === "forward") {
    return { action: "browser-control", query: "forward", confidence: 0.95 };
  }
  if (text === "close tab" || text === "close this tab") {
    return { action: "browser-control", query: "close-tab", confidence: 0.95 };
  }
  if (text === "new tab" || text === "open new tab") {
    return { action: "browser-control", query: "new-tab", confidence: 0.95 };
  }
  if (text.startsWith("scroll down")) {
    return { action: "browser-control", query: "scroll-down", confidence: 0.9 };
  }
  if (text.startsWith("scroll up")) {
    return { action: "browser-control", query: "scroll-up", confidence: 0.9 };
  }

  // "open [site] play [query]" → open-and-play
  const playMatch = text.match(/^open\s+(\w+)\s+play\s+(.+)$/);
  if (playMatch) {
    const [, site, query] = playMatch;
    const target = findSite(site);
    if (target) {
      return { action: "open-and-play", target, query, confidence: 0.92 };
    }
  }

  // "open [site] [query]" → open-and-search
  const openSearchMatch = text.match(/^open\s+(\w+)\s+(.+)$/);
  if (openSearchMatch) {
    const [, site, query] = openSearchMatch;
    const target = findSite(site);
    if (target) {
      return { action: "open-and-search", target, query, confidence: 0.9 };
    }
  }

  // "open [site]" → navigate
  const openMatch = text.match(/^(?:open|go to|goto)\s+(.+)$/);
  if (openMatch) {
    const input = openMatch[1].trim();
    // Check if it looks like a URL (has a dot) before trying site matching
    if (input.includes(".")) {
      const url = input.startsWith("http") ? input : `https://${input}`;
      return { action: "navigate", target: url, confidence: 0.85 };
    }
    const site = findSite(input);
    if (site) {
      return { action: "navigate", target: site, confidence: 0.92 };
    }
  }

  // "search [query]" or "search for [query]"
  const searchMatch = text.match(/^(?:search for|search|look up|find)\s+(.+)$/);
  if (searchMatch) {
    return { action: "search", query: searchMatch[1], confidence: 0.9 };
  }

  // Fallback: treat entire text as a search query
  return { action: "search", query: text, confidence: 0.6 };
}

/**
 * Find a site key from user input (fuzzy).
 */
function findSite(input: string): string | undefined {
  const clean = input.toLowerCase().replace(/[^a-z0-9]/g, "");

  // Direct match
  if (SITE_MAP[clean]) return clean;

  // Fuzzy: "tube" → youtube, "yt" → youtube
  const aliases: Record<string, string> = {
    tube: "youtube",
    yt: "youtube",
    gh: "github",
    tw: "twitter",
    ig: "instagram",
    so: "stackoverflow",
  };
  if (aliases[clean]) return aliases[clean];

  // Partial match
  for (const key of Object.keys(SITE_MAP)) {
    if (key.includes(clean) || clean.includes(key)) return key;
  }

  return undefined;
}

/**
 * Build the execution URL for a command.
 */
export function buildExecutionUrl(command: CommandResult): string | null {
  switch (command.action) {
    case "search":
      return `https://www.google.com/search?q=${encodeURIComponent(command.query || "")}`;

    case "navigate":
      if (command.target && SITE_MAP[command.target]) {
        return SITE_MAP[command.target];
      }
      return command.target || null;

    case "open-and-search":
    case "open-and-play":
      if (command.target && command.query && SEARCH_MAP[command.target]) {
        return SEARCH_MAP[command.target](command.query);
      }
      if (command.target && SITE_MAP[command.target]) {
        return SITE_MAP[command.target];
      }
      return null;

    default:
      return null;
  }
}
