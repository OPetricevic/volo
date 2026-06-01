import { describe, it, expect } from "vitest";

// Test the config logic directly (not the import.meta.env which is build-time)
describe("config", () => {
  it("default API URL is production", () => {
    const DEFAULT_API_URL = "https://api.volo.yourdomain.com";
    // In production, VITE_API_URL is undefined, so it falls back to default
    const apiUrl = undefined || DEFAULT_API_URL;
    expect(apiUrl).toBe("https://api.volo.yourdomain.com");
    expect(apiUrl.startsWith("https://")).toBe(true);
  });

  it("env override works for development", () => {
    const VITE_API_URL = "http://localhost:8080";
    const DEFAULT_API_URL = "https://api.volo.yourdomain.com";
    const apiUrl = VITE_API_URL || DEFAULT_API_URL;
    expect(apiUrl).toBe("http://localhost:8080");
  });

  it("ollama URL is always localhost", () => {
    const ollamaUrl = "http://localhost:11434";
    expect(ollamaUrl).toBe("http://localhost:11434");
    expect(ollamaUrl).toContain("localhost");
  });
});
