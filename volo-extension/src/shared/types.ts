/** Microphone listening mode */
export type MicMode = "always" | "once" | "off";

/** Parsed command result */
export interface CommandResult {
  action: "search" | "navigate" | "open-and-search" | "open-and-play" | "browser-control";
  target?: string;
  query?: string;
  confidence: number;
}

/** Voice recognition state */
export type VoiceState = "idle" | "listening" | "processing" | "wake-word-detected" | "error";

/** Messages between extension parts */
export type ExtensionMessage =
  | { type: "VOICE_STATE_CHANGED"; state: VoiceState }
  | { type: "COMMAND_DETECTED"; transcript: string }
  | { type: "EXECUTE_COMMAND"; command: CommandResult }
  | { type: "GET_STATE" }
  | { type: "STATE_RESPONSE"; state: VoiceState; micMode: MicMode }
  | { type: "SET_MIC_MODE"; mode: MicMode }
  | { type: "ACTIVATE_LISTENING" };

/** Known site mappings */
export const SITE_MAP: Record<string, string> = {
  youtube: "https://www.youtube.com",
  google: "https://www.google.com",
  github: "https://www.github.com",
  twitter: "https://www.x.com",
  x: "https://www.x.com",
  reddit: "https://www.reddit.com",
  spotify: "https://open.spotify.com",
  discord: "https://discord.com",
  gmail: "https://mail.google.com",
  netflix: "https://www.netflix.com",
  twitch: "https://www.twitch.tv",
  linkedin: "https://www.linkedin.com",
  stackoverflow: "https://stackoverflow.com",
};

/** Search URL patterns for sites */
export const SEARCH_MAP: Record<string, (query: string) => string> = {
  youtube: (q) => `https://www.youtube.com/results?search_query=${encodeURIComponent(q)}`,
  google: (q) => `https://www.google.com/search?q=${encodeURIComponent(q)}`,
  github: (q) => `https://github.com/search?q=${encodeURIComponent(q)}`,
  spotify: (q) => `https://open.spotify.com/search/${encodeURIComponent(q)}`,
  reddit: (q) => `https://www.reddit.com/search/?q=${encodeURIComponent(q)}`,
};
