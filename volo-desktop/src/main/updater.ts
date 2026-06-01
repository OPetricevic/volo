import { autoUpdater } from "electron";
import { app } from "electron";

/**
 * Auto-updater using Electron's built-in autoUpdater.
 * Checks GitHub Releases for new versions.
 * Downloads and installs on next app restart.
 */
export function initAutoUpdater() {
  // Only in production (packaged app)
  if (!app.isPackaged) return;

  const server = "https://update.electronjs.org";
  const repo = "volo/volo-desktop"; // GitHub owner/repo
  const url = `${server}/${repo}/${process.platform}-${process.arch}/${app.getVersion()}`;

  autoUpdater.setFeedURL({ url });

  // Check for updates every 4 hours
  setInterval(() => {
    autoUpdater.checkForUpdates();
  }, 4 * 60 * 60 * 1000);

  // Check on startup (after 10s delay to not block launch)
  setTimeout(() => {
    autoUpdater.checkForUpdates();
  }, 10_000);

  autoUpdater.on("update-downloaded", () => {
    // Will install on next restart
    console.log("[Updater] Update downloaded, will install on restart");
  });

  autoUpdater.on("error", (err) => {
    console.error("[Updater] Error:", err.message);
  });
}
