import { describe, it, expect } from "vitest";

// Test the prompt building logic extracted from ChatView
function buildPrompt(message: string, historyContext: string): string {
  return `You are Volo, a voice assistant's history helper. The user will ask questions about their browsing and search history.

Given the user's question and their recent command history below, answer naturally and concisely. Keep responses short and helpful.

If the user asks about something not in their history, say so politely.
If the user gives a command (like "open youtube"), tell them to use voice mode instead.

Recent command history:
${historyContext}

User: ${message}

Volo:`;
}

describe("buildPrompt", () => {
  it("includes the user message", () => {
    const prompt = buildPrompt("What did I search yesterday?", "(empty)");
    expect(prompt).toContain("What did I search yesterday?");
  });

  it("includes the history context", () => {
    const history = "- [Jun 1 09:14] open-and-play → youtube: lofi beats";
    const prompt = buildPrompt("What did I listen to?", history);
    expect(prompt).toContain("lofi beats");
    expect(prompt).toContain("youtube");
  });

  it("includes system instructions", () => {
    const prompt = buildPrompt("hello", "");
    expect(prompt).toContain("You are Volo");
    expect(prompt).toContain("voice assistant");
  });

  it("ends with Volo: marker for the model to continue", () => {
    const prompt = buildPrompt("test", "");
    expect(prompt.trimEnd().endsWith("Volo:")).toBe(true);
  });

  it("handles empty history gracefully", () => {
    const prompt = buildPrompt("What did I do?", "(No commands in history yet)");
    expect(prompt).toContain("No commands in history yet");
  });

  it("handles empty message", () => {
    const prompt = buildPrompt("", "some history");
    expect(prompt).toContain("User: \n");
  });

  it("handles very long history without breaking", () => {
    const longHistory = Array(50)
      .fill("- [Jun 1 09:00] search: react hooks (search react hooks)")
      .join("\n");
    const prompt = buildPrompt("summarize", longHistory);
    expect(prompt.length).toBeGreaterThan(1000);
    expect(prompt).toContain("summarize");
  });

  it("handles special characters in message", () => {
    const prompt = buildPrompt('What about "react" & <hooks>?', "");
    expect(prompt).toContain('"react"');
    expect(prompt).toContain("&");
    expect(prompt).toContain("<hooks>");
  });
});
