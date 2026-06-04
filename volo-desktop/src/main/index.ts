import {
  app,
  BrowserWindow,
  Tray,
  Menu,
  globalShortcut,
  ipcMain,
  nativeImage,
  shell,
} from "electron";
import { join } from "path";
import { OllamaManager } from "./ollama";
import { initAutoUpdater } from "./updater";

let mainWindow: BrowserWindow | null = null;
let tray: Tray | null = null;
let ollamaManager: OllamaManager;

// Single instance lock
const gotLock = app.requestSingleInstanceLock();
if (!gotLock) {
  app.quit();
}

app.on("second-instance", () => {
  if (mainWindow) {
    if (mainWindow.isMinimized()) mainWindow.restore();
    mainWindow.focus();
  }
});

app.whenReady().then(() => {
  createWindow();
  createTray();
  registerShortcuts();
  initAutoUpdater();
  ollamaManager = new OllamaManager();
  ollamaManager.initialize().catch((err) => {
    console.error("[Ollama] Init failed (non-critical):", err);
  });
});

app.on("window-all-closed", (e: Event) => {
  // Don't quit — stay in tray
  e.preventDefault();
});

app.on("before-quit", () => {
  globalShortcut.unregisterAll();
  ollamaManager?.stop();
});

// ─── Window ──────────────────────────────────────────────

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1100,
    height: 720,
    minWidth: 800,
    minHeight: 500,
    frame: false,
    backgroundColor: "#0a0a0c",
    show: false,
    webPreferences: {
      preload: join(__dirname, "../preload/index.mjs"),
      contextIsolation: true,
      nodeIntegration: false,
      sandbox: false,
    },
  });

  // Load renderer
  if (process.env.ELECTRON_RENDERER_URL) {
    mainWindow.loadURL(process.env.ELECTRON_RENDERER_URL);
  } else {
    mainWindow.loadFile(join(__dirname, "../renderer/index.html"));
  }

  mainWindow.on("ready-to-show", () => {
    mainWindow?.show();
  });

  // Minimize to tray instead of closing
  mainWindow.on("close", (e) => {
    if (!app.isQuitting) {
      e.preventDefault();
      mainWindow?.hide();
    }
  });
}

// ─── Tray ────────────────────────────────────────────────

function createTray() {
  // Placeholder icon — replace with real icon
  const icon = nativeImage.createEmpty();
  tray = new Tray(icon);

  const contextMenu = Menu.buildFromTemplate([
    { label: "Open Volo", click: () => mainWindow?.show() },
    { type: "separator" },
    { label: "Listening: On", enabled: false },
    { type: "separator" },
    {
      label: "Quit",
      click: () => {
        (app as typeof app & { isQuitting: boolean }).isQuitting = true;
        app.quit();
      },
    },
  ]);

  tray.setToolTip("Volo");
  tray.setContextMenu(contextMenu);
  tray.on("click", () => mainWindow?.show());
}

// ─── Global Shortcuts ────────────────────────────────────

function registerShortcuts() {
  globalShortcut.register("CommandOrControl+Shift+V", () => {
    if (mainWindow?.isVisible()) {
      mainWindow.focus();
    } else {
      mainWindow?.show();
    }
    // Notify renderer to activate voice
    mainWindow?.webContents.send("activate-voice");
  });
}

// ─── IPC Handlers ────────────────────────────────────────

// Window controls (custom titlebar)
ipcMain.on("window:minimize", () => {
  console.log("[IPC] minimize");
  mainWindow?.minimize();
});
ipcMain.on("window:maximize", () => {
  console.log("[IPC] maximize");
  if (mainWindow?.isMaximized()) {
    mainWindow.unmaximize();
  } else {
    mainWindow?.maximize();
  }
});
ipcMain.on("window:close", () => {
  console.log("[IPC] close");
  mainWindow?.hide();
});

// Open external links
ipcMain.on("open-external", (_e, url: string) => {
  shell.openExternal(url);
});

// Ollama status
ipcMain.handle("ollama:status", async () => {
  return ollamaManager.isRunning();
});

ipcMain.handle("ollama:install", async () => {
  return ollamaManager.install();
});
