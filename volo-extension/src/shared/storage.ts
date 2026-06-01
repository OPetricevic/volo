import type { MicMode } from "./types";

const STORAGE_KEYS = {
  micMode: "volo_mic_mode",
  onboarded: "volo_onboarded",
} as const;

export async function getMicMode(): Promise<MicMode> {
  const result = await chrome.storage.local.get(STORAGE_KEYS.micMode);
  return (result[STORAGE_KEYS.micMode] as MicMode) || "always";
}

export async function setMicMode(mode: MicMode): Promise<void> {
  await chrome.storage.local.set({ [STORAGE_KEYS.micMode]: mode });
}

export async function getOnboarded(): Promise<boolean> {
  const result = await chrome.storage.local.get(STORAGE_KEYS.onboarded);
  return result[STORAGE_KEYS.onboarded] === true;
}

export async function setOnboarded(value: boolean): Promise<void> {
  await chrome.storage.local.set({ [STORAGE_KEYS.onboarded]: value });
}
