import { app } from "electron";
import { join } from "path";
import { existsSync } from "fs";
import { execFile, ChildProcess } from "child_process";

/**
 * Manages the local Ollama instance for Volo AI chat.
 * Ollama is optional — installed via checkbox during setup or later from Settings.
 */
export class OllamaManager {
  private process: ChildProcess | null = null;
  private ollamaPath: string;
  private running = false;

  constructor() {
    // Ollama binary lives inside the app's resources
    this.ollamaPath = join(app.getPath("userData"), "ollama", "ollama.exe");
  }

  /** Check if Ollama is installed locally */
  isInstalled(): boolean {
    return existsSync(this.ollamaPath);
  }

  /** Check if Ollama is running (either our instance or system-wide) */
  async isRunning(): Promise<boolean> {
    try {
      const response = await fetch("http://localhost:11434/api/tags", {
        signal: AbortSignal.timeout(2000),
      });
      return response.ok;
    } catch {
      return false;
    }
  }

  /** Check status and start if installed but not running */
  async checkStatus(): Promise<void> {
    const running = await this.isRunning();
    if (running) {
      this.running = true;
      return;
    }

    if (this.isInstalled()) {
      await this.start();
    }
  }

  /** Start the local Ollama server */
  async start(): Promise<boolean> {
    if (!this.isInstalled()) return false;

    return new Promise((resolve) => {
      this.process = execFile(this.ollamaPath, ["serve"], (error) => {
        if (error) {
          console.error("[Ollama] Failed to start:", error);
          this.running = false;
          resolve(false);
        }
      });

      // Give it a moment to start
      setTimeout(async () => {
        this.running = await this.isRunning();
        resolve(this.running);
      }, 2000);
    });
  }

  /** Stop the local Ollama server */
  stop(): void {
    if (this.process) {
      this.process.kill();
      this.process = null;
      this.running = false;
    }
  }

  /** Download and install Ollama + Gemma model */
  async install(): Promise<{ success: boolean; error?: string }> {
    // This would be handled by the NSIS installer during setup.
    // If installing from Settings post-install, we download the binary
    // and pull the model. For now, return a stub.
    return { success: false, error: "Install from Settings not yet implemented. Reinstall Volo with AI enabled." };
  }

  /** Check if the Gemma model is available */
  async hasModel(): Promise<boolean> {
    try {
      const response = await fetch("http://localhost:11434/api/tags");
      const data = await response.json();
      return data.models?.some((m: { name: string }) => m.name.includes("gemma"));
    } catch {
      return false;
    }
  }
}
