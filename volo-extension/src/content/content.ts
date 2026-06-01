/**
 * Volo Content Script
 * Runs on every page. Handles voice recognition via Web Speech API.
 *
 * NOTE: This file is built as a content script (classic script, no ES imports).
 * All dependencies are inlined here to avoid chunk splitting issues.
 */

// === Inlined types ===
type MicMode = "always" | "once" | "off";
type VoiceState = "idle" | "listening" | "processing" | "wake-word-detected";

// === Inlined wake word detection ===
const WAKE_WORD = "hey volo";

function containsWakeWord(transcript: string): boolean {
  return transcript.toLowerCase().trim().includes(WAKE_WORD);
}

function extractAfterWakeWord(transcript: string): string | null {
  const lower = transcript.toLowerCase().trim();
  const idx = lower.indexOf(WAKE_WORD);
  if (idx === -1) return null;
  const after = lower.slice(idx + WAKE_WORD.length).trim();
  return after.length > 0 ? after : "";
}

// === Inlined storage ===
async function getMicMode(): Promise<MicMode> {
  const result = await chrome.storage.local.get("volo_mic_mode");
  return (result["volo_mic_mode"] as MicMode) || "always";
}

// === State ===
let recognition: SpeechRecognition | null = null;
let voiceState: VoiceState = "idle";
let micMode: MicMode = "always";
let wakeWordDetected = false;
let commandBuffer = "";
let silenceTimer: ReturnType<typeof setTimeout> | null = null;

// === Init ===
async function init() {
  micMode = await getMicMode();

  if (micMode === "always") {
    startListening();
  }

  chrome.runtime.onMessage.addListener((message: { type: string; mode?: MicMode; command?: { action: string; query?: string } }) => {
    switch (message.type) {
      case "SET_MIC_MODE":
        micMode = message.mode!;
        if (micMode === "off") {
          stopListening();
        } else if (micMode === "always") {
          startListening();
        }
        break;
      case "ACTIVATE_LISTENING":
        startListening();
        break;
      case "EXECUTE_COMMAND":
        if (message.command) {
          executeLocalCommand(message.command);
        }
        break;
    }
  });
}

function startListening() {
  if (recognition) return;

  const SpeechRecognitionCtor = window.SpeechRecognition || window.webkitSpeechRecognition;

  if (!SpeechRecognitionCtor) {
    console.warn("[Volo] Web Speech API not supported in this browser.");
    return;
  }

  recognition = new SpeechRecognitionCtor();
  recognition.continuous = true;
  recognition.interimResults = true;
  recognition.lang = "en-US";

  recognition.onstart = () => {
    setState("listening");
    console.log("[Volo] Listening...");
  };

  recognition.onresult = (event: SpeechRecognitionEvent) => {
    const lastResult = event.results[event.results.length - 1];
    const transcript = lastResult[0].transcript;
    const isFinal = lastResult.isFinal;

    if (!wakeWordDetected) {
      if (containsWakeWord(transcript)) {
        wakeWordDetected = true;
        commandBuffer = extractAfterWakeWord(transcript) || "";
        setState("wake-word-detected");
        showOverlay();

        if (isFinal && commandBuffer.length > 0) {
          processCommand(commandBuffer);
        } else {
          resetSilenceTimer();
        }
      }
    } else {
      const afterWake = extractAfterWakeWord(transcript);
      commandBuffer = afterWake || transcript;

      if (isFinal) {
        processCommand(commandBuffer);
      } else {
        resetSilenceTimer();
      }
    }
  };

  recognition.onerror = (event: SpeechRecognitionErrorEvent) => {
    console.warn("[Volo] Speech error:", event.error);
    if (event.error === "not-allowed") {
      micMode = "off";
      setState("idle");
    }
  };

  recognition.onend = () => {
    if (micMode === "always" && voiceState !== "processing") {
      wakeWordDetected = false;
      commandBuffer = "";
      setTimeout(() => {
        if (micMode === "always") {
          recognition = null;
          startListening();
        }
      }, 300);
    } else {
      setState("idle");
      recognition = null;
    }
  };

  try {
    recognition.start();
  } catch (e) {
    console.warn("[Volo] Failed to start recognition:", e);
    recognition = null;
  }
}

function stopListening() {
  if (recognition) {
    recognition.abort();
    recognition = null;
  }
  wakeWordDetected = false;
  commandBuffer = "";
  setState("idle");
  hideOverlay();
}

function processCommand(command: string) {
  if (!command.trim()) {
    wakeWordDetected = false;
    commandBuffer = "";
    hideOverlay();
    return;
  }

  setState("processing");
  console.log("[Volo] Command:", command);

  chrome.runtime.sendMessage({
    type: "COMMAND_DETECTED",
    transcript: command.trim(),
  });

  wakeWordDetected = false;
  commandBuffer = "";
  clearSilenceTimer();
  hideOverlay();

  setTimeout(() => {
    if (micMode === "always") {
      setState("listening");
    } else {
      setState("idle");
    }
  }, 500);
}

function resetSilenceTimer() {
  clearSilenceTimer();
  silenceTimer = setTimeout(() => {
    if (wakeWordDetected && commandBuffer.trim()) {
      processCommand(commandBuffer);
    }
  }, 2000);
}

function clearSilenceTimer() {
  if (silenceTimer) {
    clearTimeout(silenceTimer);
    silenceTimer = null;
  }
}

function setState(state: VoiceState) {
  voiceState = state;
  chrome.runtime.sendMessage({
    type: "VOICE_STATE_CHANGED",
    state,
  }).catch(() => {});
}

function executeLocalCommand(command: { action: string; query?: string }) {
  if (command.query === "scroll-down") {
    window.scrollBy({ top: 400, behavior: "smooth" });
  } else if (command.query === "scroll-up") {
    window.scrollBy({ top: -400, behavior: "smooth" });
  }
}

// === Overlay UI ===
let overlay: HTMLElement | null = null;

function showOverlay() {
  if (overlay) return;

  overlay = document.createElement("div");
  overlay.id = "volo-overlay";
  overlay.innerHTML = `
    <div style="
      position: fixed;
      top: 20px;
      right: 20px;
      z-index: 2147483647;
      background: #2b2d31;
      border: 1px solid #5865f2;
      border-radius: 10px;
      padding: 10px 18px;
      display: flex;
      align-items: center;
      gap: 10px;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4);
      font-family: -apple-system, system-ui, sans-serif;
      animation: volo-fade-in 0.2s ease;
    ">
      <div style="
        width: 8px;
        height: 8px;
        background: #23a559;
        border-radius: 50%;
        animation: volo-pulse 1.2s infinite;
      "></div>
      <span style="color: #dbdee1; font-size: 12px; font-weight: 500;">Listening...</span>
    </div>
    <style>
      @keyframes volo-pulse {
        0%, 100% { opacity: 1; transform: scale(1); }
        50% { opacity: 0.5; transform: scale(1.3); }
      }
      @keyframes volo-fade-in {
        from { opacity: 0; transform: translateY(-8px); }
        to { opacity: 1; transform: translateY(0); }
      }
    </style>
  `;
  document.body.appendChild(overlay);
}

function hideOverlay() {
  if (overlay) {
    overlay.remove();
    overlay = null;
  }
}

// Start
init();
