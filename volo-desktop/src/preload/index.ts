import { contextBridge, ipcRenderer } from "electron";

/**
 * Preload script — exposes safe IPC methods to the renderer.
 * The renderer cannot access Node.js or Electron directly.
 */
contextBridge.exposeInMainWorld("volo", {
  // Window controls
  minimize: () => ipcRenderer.send("window:minimize"),
  maximize: () => ipcRenderer.send("window:maximize"),
  close: () => ipcRenderer.send("window:close"),

  // External links
  openExternal: (url: string) => ipcRenderer.send("open-external", url),

  // Ollama
  getOllamaStatus: () => ipcRenderer.invoke("ollama:status") as Promise<boolean>,
  installOllama: () => ipcRenderer.invoke("ollama:install") as Promise<{ success: boolean; error?: string }>,

  // Voice activation from main process (global hotkey)
  onActivateVoice: (callback: () => void) => {
    ipcRenderer.on("activate-voice", callback);
    return () => ipcRenderer.removeListener("activate-voice", callback);
  },
});

// Type declaration for the renderer
export type VoloAPI = {
  minimize: () => void;
  maximize: () => void;
  close: () => void;
  openExternal: (url: string) => void;
  getOllamaStatus: () => Promise<boolean>;
  installOllama: () => Promise<{ success: boolean; error?: string }>;
  onActivateVoice: (callback: () => void) => () => void;
};
