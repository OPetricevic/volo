import type { ExtensionMessage, VoiceState, MicMode } from "@/shared/types";
import { parseCommand, buildExecutionUrl } from "@/shared/parser";
import { getMicMode, setMicMode } from "@/shared/storage";
import { registerDevice, isAuthenticated, sendCommand, healthCheck } from "@/shared/api";
import { enqueue, flush } from "@/shared/retry-queue";

// Global state
let voiceState: VoiceState = "idle";
let micMode: MicMode = "always";
let apiOnline = false;

// Initialize on install
chrome.runtime.onInstalled.addListener(async () => {
  micMode = await getMicMode();
  console.log("[Volo] Extension installed. Mic mode:", micMode);

  // Register device with backend (get JWT)
  const hasToken = await isAuthenticated();
  if (!hasToken) {
    const registered = await registerDevice();
    console.log("[Volo] Device registration:", registered ? "success" : "failed (offline mode)");
  }

  // Check API health
  apiOnline = await healthCheck();
  console.log("[Volo] API status:", apiOnline ? "online" : "offline");
});

// Also check on startup (not just install)
chrome.runtime.onStartup.addListener(async () => {
  micMode = await getMicMode();
  apiOnline = await healthCheck();

  // Retry registration if we don't have a token yet
  const hasToken = await isAuthenticated();
  if (!hasToken) {
    await registerDevice();
  }
});

// Handle messages from content script and popup
chrome.runtime.onMessage.addListener(
  (message: ExtensionMessage, _sender, sendResponse) => {
    switch (message.type) {
      case "GET_STATE":
        sendResponse({ type: "STATE_RESPONSE", state: voiceState, micMode, apiOnline });
        return true;

      case "SET_MIC_MODE":
        micMode = message.mode;
        setMicMode(message.mode);
        broadcastToTabs({ type: "SET_MIC_MODE", mode: message.mode });
        sendResponse({ success: true });
        return true;

      case "VOICE_STATE_CHANGED":
        voiceState = message.state;
        updateBadge(message.state);
        return false;

      case "COMMAND_DETECTED":
        handleCommand(message.transcript, _sender.tab?.url);
        return false;

      case "ACTIVATE_LISTENING":
        broadcastToTabs({ type: "ACTIVATE_LISTENING" });
        return false;
    }
  }
);

// Handle keyboard shortcut
chrome.commands.onCommand.addListener((command) => {
  if (command === "activate-volo") {
    broadcastToTabs({ type: "ACTIVATE_LISTENING" });
  }
});

async function handleCommand(transcript: string, currentUrl?: string) {
  console.log("[Volo] Received transcript:", transcript);

  // 0. Check macros FIRST (user-defined custom commands)
  const { matchMacro, isMacroCacheStale, fetchMacros } = await import("@/shared/api");

  // Refresh cache if stale
  if (await isMacroCacheStale()) {
    console.log("[Volo] Macro cache stale, refreshing...");
    await fetchMacros().catch((err) => console.warn("[Volo] Macro fetch failed:", err));
  }

  const macro = await matchMacro(transcript);
  if (macro) {
    console.log("[Volo] Macro matched:", macro.name, "→", macro.actions.length, "actions");
    for (const action of macro.actions) {
      if (action.type === "navigate" && action.url) {
        await chrome.tabs.create({ url: action.url });
        console.log("[Volo] Macro action: opened", action.url);
      }
    }
    // Still send to API for history tracking
    if (apiOnline) {
      sendCommand(transcript, currentUrl).catch(() => {});
    }
    return;
  }

  // 1. No macro match — parse locally (instant response)
  const localCommand = parseCommand(transcript);
  console.log("[Volo] Local parse:", localCommand);

  // 2. Execute immediately based on local parse (zero latency)
  await executeCommand(localCommand);

  // 3. Send to backend in background (for history + learning)
  //    Don't await — fire and forget so the user isn't blocked
  if (apiOnline) {
    sendCommand(transcript, currentUrl).then((apiResult) => {
      if (apiResult) {
        console.log("[Volo] API response:", apiResult);
      } else {
        // API returned null (network error) — queue for retry
        enqueue(transcript, currentUrl);
        apiOnline = false;
      }
    }).catch(() => {
      // API call failed — queue for retry
      enqueue(transcript, currentUrl);
      apiOnline = false;
    });
  } else {
    // API is offline — queue immediately
    enqueue(transcript, currentUrl);
  }
}

async function executeCommand(command: ReturnType<typeof parseCommand>) {
  // Browser control actions
  if (command.action === "browser-control") {
    const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
    if (!tab?.id) return;

    switch (command.query) {
      case "back":
        chrome.tabs.goBack(tab.id);
        break;
      case "forward":
        chrome.tabs.goForward(tab.id);
        break;
      case "close-tab":
        chrome.tabs.remove(tab.id);
        break;
      case "new-tab":
        chrome.tabs.create({});
        break;
      case "scroll-down":
      case "scroll-up":
        chrome.tabs.sendMessage(tab.id, {
          type: "EXECUTE_COMMAND",
          command,
        });
        break;
    }
    return;
  }

  // Navigation actions — open URL
  const url = buildExecutionUrl(command);
  if (url) {
    const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
    if (tab?.id) {
      chrome.tabs.update(tab.id, { url });
    } else {
      chrome.tabs.create({ url });
    }
  }
}

function updateBadge(state: VoiceState) {
  switch (state) {
    case "listening":
      chrome.action.setBadgeBackgroundColor({ color: "#34d399" });
      chrome.action.setBadgeText({ text: "●" });
      break;
    case "processing":
      chrome.action.setBadgeBackgroundColor({ color: "#fbbf24" });
      chrome.action.setBadgeText({ text: "..." });
      break;
    case "wake-word-detected":
      chrome.action.setBadgeBackgroundColor({ color: "#6366f1" });
      chrome.action.setBadgeText({ text: "!" });
      break;
    default:
      chrome.action.setBadgeText({ text: "" });
  }
}

async function broadcastToTabs(message: ExtensionMessage) {
  const tabs = await chrome.tabs.query({});
  for (const tab of tabs) {
    if (tab.id) {
      chrome.tabs.sendMessage(tab.id, message).catch(() => {
        // Tab might not have content script loaded
      });
    }
  }
}

// Periodic health check (every 5 minutes)
setInterval(async () => {
  const wasOffline = !apiOnline;
  apiOnline = await healthCheck();

  // If API just came back online, flush the retry queue
  if (apiOnline && wasOffline) {
    flush(async (transcript, currentUrl) => {
      const result = await sendCommand(transcript, currentUrl);
      return result !== null;
    });
  }
}, 5 * 60 * 1000);
