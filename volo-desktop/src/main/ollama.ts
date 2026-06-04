import { app } from "electron";
import { join } from "path";
import { existsSync } from "fs";
import { execFile, ChildProcess } from "child_process";

/**
 * Manages the bundled Ollama instance for Volo AI chat.
 * Ollama is always included with the installer — no separate download needed.
 * On first launch, it checks if the model is pulled and starts the server.
 */
export class OllamaManager {
  private process: ChildProcess | null = null;
  private ollamaPath: string;
  private running = false;

  constructor() {
    // Ollama installs to LocalAppData\Programs\Ollama by default on Windows
    const localAppData = process.env.LOCALAPPDATA || join(app.getPath("home"), "AppData", "Local");
    this.ollamaPath = join(localAppData, "Programs", "Ollama", "ollama.exe");
  }

  /** Check if Ollama binary exists (should always be true in packaged app) */
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

  /** Initialize — start Ollama and ensure model is available */
  async initialize(): Promise<void> {
    const running = await this.isRunning();
    if (running) {
      this.running = true;
      console.log("[Ollama] Already running");
      return;
    }

    if (!this.isInstalled()) {
      console.error("[Ollama] Binary not found at:", this.ollamaPath);
      return;
    }

    await this.start();

    // Ensure model is pulled (handles case where install was interrupted)
    if (this.running) {
      const hasModel = await this.hasModel();
      if (!hasModel) {
        console.log("[Ollama] Model not found, pulling gemma2:2b...");
        await this.pullModel();
      }
    }
  }

  /** Start the local Ollama server */
  async start(): Promise<boolean> {
    if (!this.isInstalled()) return false;

    return new Promise((resolve) => {
      this.process = execFile(this.ollamaPath, ["serve"], (error) => {
        if (error && !error.killed) {
          console.error("[Ollama] Process exited:", error.message);
          this.running = false;
        }
      });

      // Give it a moment to start
      setTimeout(async () => {
        this.running = await this.isRunning();
        if (this.running) {
          console.log("[Ollama] Server started successfully");
        } else {
          console.error("[Ollama] Failed to start within timeout");
        }
        resolve(this.running);
      }, 3000);
    });
  }

  /** Stop the local Ollama server */
  stop(): void {
    if (this.process) {
      this.process.kill();
      this.process = null;
      this.running = false;
      console.log("[Ollama] Server stopped");
    }
  }

  /** Pull the Gemma model (called if model is missing on startup) */
  private async pullModel(): Promise<boolean> {
    return new Promise((resolve) => {
      execFile(this.ollamaPath, ["pull", "gemma2:2b"], (error) => {
        if (error) {
          console.error("[Ollama] Failed to pull model:", error.message);
          resolve(false);
        } else {
          console.log("[Ollama] Model pulled successfully");
          resolve(true);
        }
      });
    });
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

  /** Get current status for the renderer */
  getStatus(): { installed: boolean; running: boolean } {
    return { installed: this.isInstalled(), running: this.running };
  }
}
